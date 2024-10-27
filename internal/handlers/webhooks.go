package handlers

import (
	"encoding/json"
	"github.com/Ayano2000/push/internal/pkg/transformer"
	"github.com/Ayano2000/push/internal/types"
	"github.com/rs/zerolog/log"
	"net/http"
)

// CreateWebhook will create a minio Webhook,
// a database row and update the server to listen for requests
// made to http://basepath/<webhook_name>
func (h *Handler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	var webhook types.Webhook
	err := json.NewDecoder(r.Body).Decode(&webhook)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// JQ filter is valid
	err = transformer.ValidFilter(webhook.JQFilter)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create webhook if the name isn't already taken
	err = h.Services.Minio.CreateBucket(r.Context(), webhook)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add db row
	err = h.Services.DB.CreateWebhook(r.Context(), webhook)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetWebhooks(w http.ResponseWriter, r *http.Request) {
	webhooks, err := h.Services.DB.GetWebhooks(r.Context())
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(webhooks); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (h *Handler) GetWebhookContent(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	webhook, err := h.Services.DB.GetWebhookByName(r.Context(), name)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := h.Services.Minio.GetObjects(r.Context(), webhook)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(content); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (h *Handler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	// todo
}

func (h *Handler) DeleteWebhookContents(w http.ResponseWriter, r *http.Request) {
	// todo
}
