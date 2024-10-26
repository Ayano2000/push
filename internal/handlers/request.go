package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Ayano2000/push/internal/types"
	"github.com/google/uuid"
	"github.com/itchyny/gojq"
	"github.com/minio/minio-go/v7"
	"io"
	"log"
	"net/http"
)

// HandleRequest will read and dump the request body in minio: after running it
// through the jq filter for the endpoint (if one is set), before forwarding it to
// the endpoints defined forward_to value (if one is set)
func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request, bucket types.Bucket) {
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

	_, err = h.Services.Minio.PutObject(
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
