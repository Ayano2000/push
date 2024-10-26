package handlers

import (
	"encoding/json"
	"github.com/Ayano2000/push/internal/types"
	"github.com/itchyny/gojq"
	"github.com/minio/minio-go/v7"
	"net/http"
)

func (h *Handler) CreateBucket(w http.ResponseWriter, r *http.Request) {
	var bucket types.Bucket
	err := json.NewDecoder(r.Body).Decode(&bucket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// JQ filter is valid
	if bucket.JQFilter != "" {
		_, err = gojq.Parse(bucket.JQFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Create bucket if the name isn't already taken
	exists, err := h.Services.Minio.BucketExists(r.Context(), bucket.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Bucket already exists", http.StatusConflict)
		return
	}

	err = h.Services.Minio.MakeBucket(r.Context(), bucket.Name, minio.MakeBucketOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add db row
	_, err = h.Services.DB.Exec(r.Context(), `
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
