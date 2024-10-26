package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/itchyny/gojq"
	"github.com/jackc/pgx/v5"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"net/http"
	"os"
)

type Bucket struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Method      string `json:"method"`
	JQFilter    string `json:"jq_filter,omitempty"`
	ForwardTo   string `json:"forward_to,omitempty"`
}

func main() {
	mux := http.NewServeMux()

	services := NewServices()
	defer services.Cleanup()

	var buckets []Bucket
	err := services.DB.QueryRow(context.Background(), "SELECT * FROM buckets").Scan(buckets)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to query endpoints: %v\n", err)
		os.Exit(1)
	}

	mux.HandleFunc("POST /endpoint", services.CreateBucket)

	for _, bucket := range buckets {
		pattern := fmt.Sprintf("%s %s", bucket.Method, bucket.Path)
		mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			services.HandleRequest(w, r, bucket)
		})
	}
}

// HandleRequest will read and dump the request body in minio: after running it
// through the jq filter for the endpoint (if one is set), before forwarding it to
// the endpoints defined forward_to value (if one is set)
func (s *Services) HandleRequest(w http.ResponseWriter, r *http.Request, bucket Bucket) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uid, err := uuid.NewRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// todo: an option to dump the pre transform payload as well
	payload := body
	var results []interface{}
	if bucket.JQFilter != "" {
		query, err := gojq.Parse(bucket.JQFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		iter := query.RunWithContext(r.Context(), body)
		for {
			value, hasNextValue := iter.Next()
			if !hasNextValue {
				break
			}

			if err, ok := value.(error); ok {
				var haltError *gojq.HaltError
				if errors.As(err, &haltError) && haltError.Value() == nil {
					break
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			results = append(results, value)
		}
	}

	if len(results) > 0 {
		payload, err = json.Marshal(results)
		if err != nil {
			http.Error(w, "JQ Filter produced invalid result: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	_, err = s.Minio.PutObject(
		r.Context(),
		r.RequestURI,
		uid.String(),
		io.NopCloser(bytes.NewReader(payload)),
		int64(len(payload)),
		minio.PutObjectOptions{
			ContentType: "application/json",
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// todo: request forwarding
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(payload); err != nil {
		log.Printf("Error writing response: %v", err)
	}
	return
}

func (s *Services) CreateBucket(w http.ResponseWriter, r *http.Request) {
	var bucket Bucket
	err := json.NewDecoder(r.Body).Decode(&bucket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// JQ filter is valid
	if bucket.JQFilter != "" {
		_, err = gojq.Parse(bucket.JQFilter)
		if err != nil {
			http.Error(w, "Invalid JQ filter provided", http.StatusInternalServerError)
			return
		}
	}

	// Create bucket if the name isn't already taken
	exists, err := s.Minio.BucketExists(r.Context(), bucket.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Bucket already exists", http.StatusConflict)
		return
	}

	err = s.Minio.MakeBucket(r.Context(), bucket.Name, minio.MakeBucketOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add db row
	_, err = s.DB.Exec(r.Context(), `
		INSERT INTO buckets (name, path, method, description, jq_filter, forward_to) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		bucket.Name, bucket.Path, bucket.Method, bucket.Description, bucket.JQFilter, bucket.ForwardTo,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type Services struct {
	DB    *pgx.Conn
	Minio *minio.Client
}

func NewServices() *Services {
	pgsqlConn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	minioConn, err := minio.New(os.Getenv("MINIO_HOST"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), ""),
		Secure: os.Getenv("MINIO_USE_SSL") == "true",
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &Services{
		DB:    pgsqlConn,
		Minio: minioConn,
	}
}

func (s *Services) Cleanup() error {
	return s.DB.Close(context.Background())
}
