package handlers

import (
	"encoding/json"
	"github.com/Ayano2000/push/internal/pkg/transformer"
	"github.com/Ayano2000/push/internal/types"
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

func (h *Handler) ListBuckets(w http.ResponseWriter, r *http.Request) {

}
