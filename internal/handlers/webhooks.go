package handlers

import (
	"encoding/json"
	"github.com/Ayano2000/push/internal/pkg/logger"
	"github.com/Ayano2000/push/internal/pkg/transformer"
	"github.com/Ayano2000/push/internal/types"
	"github.com/pkg/errors"
	"net/http"
)

// CreateWebhook will create a minio Webhook,
// a database row and update the server to listen for requests
// made to http://basepath/<webhook_name>
func (h *Handler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	log := logger.GetFromContext(r.Context())
	var webhook types.Webhook
	err := json.NewDecoder(r.Body).Decode(&webhook)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, requestBodyDecodingErrorMessage, http.StatusInternalServerError)
		return
	}

	_, err = transformer.IsValidFilter(webhook.JQFilter)
	if err != nil {
		log.Error().Err(err).Msg("Failed to validate JQ filter")
		http.Error(w, invalidJQFilterErrorMessage, http.StatusBadRequest)
		return
	}

	err = h.Services.Minio.CreateBucket(r.Context(), webhook)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create minio bucket")
		http.Error(w, createWebhookErrorMessage, http.StatusInternalServerError)
		return
	}

	err = h.Services.DB.CreateWebhook(r.Context(), webhook)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create webhook row in psql")
		http.Error(w, createWebhookErrorMessage, http.StatusInternalServerError)
		return
	}

	// update router to include this route
	if registrar, ok := r.Context().Value(muxContextKey).(types.WebhookRegistrar); ok {
		registrar.RegisterWebhook(webhook)
	} else {
		err = errors.WithStack(errors.Errorf("failed to retrieve WebhookRegistrar from context"))
		log.Error().Err(err).Msg("Failed to retrieve WebhookRegistrar from context")
		http.Error(w, createWebhookErrorMessage, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetWebhooks(w http.ResponseWriter, r *http.Request) {
	log := logger.GetFromContext(r.Context())

	webhooks, err := h.Services.DB.GetWebhooks(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("")
		http.Error(w, getWebhooksErrorMessage, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(webhooks); err != nil {
		log.Error().Err(err).Msg("Failed to encode response")
	}
}

func (h *Handler) GetWebhookContent(w http.ResponseWriter, r *http.Request) {
	log := logger.GetFromContext(r.Context())

	params, ok := r.Context().Value(urlParamContextKey).(map[string]string)
	if !ok {
		err := errors.WithStack(
			errors.Errorf("failed to retrieve webhook name from context"))
		log.Error().Err(err).Msg("Failed to retrieve webhook name from context")
		http.Error(w, getWebhookContentErrorMessage, http.StatusInternalServerError)
		return
	}

	webhook, err := h.Services.DB.GetWebhookByName(r.Context(), params["name"])
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve webhook from db")
		http.Error(w, getWebhookContentErrorMessage, http.StatusInternalServerError)
		return
	}

	content, err := h.Services.Minio.GetObjects(r.Context(), webhook)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list objects from minio")
		http.Error(w, getWebhookContentErrorMessage, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(content); err != nil {
		log.Error().Err(err).Msg("Failed to encode response")
	}
}

func (h *Handler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	// todo
}

func (h *Handler) DeleteWebhookContents(w http.ResponseWriter, r *http.Request) {
	// todo
}
