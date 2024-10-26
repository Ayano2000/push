package handlers

import (
	"github.com/Ayano2000/push/internal/pkg/transformer"
	"github.com/Ayano2000/push/internal/types"
	"io"
	"log"
	"net/http"
)

// HandleRequest will read and dump the request body in minio: after running it
// through the jq filter for the endpoint (if one is set), before forwarding it to
// the endpoints defined forward_to value (if one is set)
func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request, bucket types.Bucket) {
	preTransform, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if bucket.PreservePayload {
		err = h.Services.Minio.PutObject(r.Context(), bucket, string(preTransform))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	postTransform, err := transformer.Transform(r.Context(), string(preTransform), bucket.JQFilter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.Services.Minio.PutObject(r.Context(), bucket, postTransform)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// todo: request forwarding

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write([]byte(postTransform)); err != nil {
		log.Printf("Error writing response: %v", err)
	}
	return
}
