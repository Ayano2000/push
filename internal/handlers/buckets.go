package handlers

import (
	"encoding/json"
	"github.com/Ayano2000/push/internal/pkg/transformer"
	"github.com/Ayano2000/push/internal/types"
	"log"
	"net/http"
)

// CreateBucket will create a minio Bucket,
// a database row and update the server to listen for requests
// made to http://basepath/<bucket_name>
func (h *Handler) CreateBucket(w http.ResponseWriter, r *http.Request) {
	var bucket types.Bucket
	err := json.NewDecoder(r.Body).Decode(&bucket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// JQ filter is valid
	err = transformer.ValidFilter(bucket.JQFilter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create bucket if the name isn't already taken
	err = h.Services.Minio.CreateBucket(r.Context(), bucket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add db row
	err = h.Services.DB.CreateBucket(r.Context(), bucket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ListBuckets(w http.ResponseWriter, r *http.Request) {
	buckets, err := h.Services.DB.GetBuckets(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(buckets); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (h *Handler) GetBucketContents(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	bucket, err := h.Services.DB.GetBucketByName(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := h.Services.Minio.GetObjects(r.Context(), bucket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(content); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
