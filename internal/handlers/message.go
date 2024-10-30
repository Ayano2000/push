package handlers

import (
	"github.com/Ayano2000/push/internal/pkg/logger"
	"github.com/Ayano2000/push/internal/pkg/types"
	"github.com/Ayano2000/push/pkg/transformer"
	"io"
	"net/http"
)

// HandleMessage will read and dump the request body in minio: after running it
// through the jq filter for the endpoint (if one is set), before forwarding it to
// the endpoints defined forward_to value (if one is set)
func (h *Handler) HandleMessage(w http.ResponseWriter, r *http.Request, wh types.Webhook) {
	log := logger.GetFromContext(r.Context())

	preTransform, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request body")
		http.Error(w, requestBodyDecodingErrorMessage, http.StatusInternalServerError)
		return
	}

	if wh.PreservePayload {
		err = h.Services.Minio.PutObject(r.Context(), wh.Name, string(preTransform))
		if err != nil {
			log.Error().Err(err).Msg("failed to upload object to minio")
			http.Error(w, minioUploadErrorMessage, http.StatusInternalServerError)
			return
		}
	}

	postTransform, err := transformer.Transform(r.Context(), string(preTransform), wh.JQFilter)
	if err != nil {
		log.Error().Err(err).Msg("failed to process JQ transform")
		http.Error(w, jqTransformErrorMessage, http.StatusInternalServerError)
		return
	}

	err = h.Services.Minio.PutObject(r.Context(), wh.Name, postTransform)
	if err != nil {
		log.Error().Err(err).Msg("failed to upload object to minio")
		http.Error(w, minioUploadErrorMessage, http.StatusInternalServerError)
		return
	}

	// todo: request forwarding

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write([]byte(postTransform)); err != nil {
		log.Error().Err(err).Msg("failed to write response")
	}
}
