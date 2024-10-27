package routes

import (
	"context"
	"fmt"
	"github.com/Ayano2000/push/internal/handlers"
	"github.com/Ayano2000/push/internal/types"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *handlers.Handler) error {
	mux.HandleFunc("POST /webhooks", handler.CreateWebhook)
	mux.HandleFunc("GET /webhooks", handler.GetWebhooks)
	mux.HandleFunc("GET /webhooks/{name}/content", handler.GetWebhookContent)
	mux.HandleFunc("DELETE /webhooks/{name}", handler.DeleteWebhook)
	mux.HandleFunc("DELETE /webhooks/{name}/content", handler.DeleteWebhookContents)

	// custom endpoints
	var webhooks []types.Webhook
	webhooks, err := handler.Services.DB.GetWebhooks(context.Background())
	if err != nil {
		return err
	}

	for _, webhook := range webhooks {
		pattern := fmt.Sprintf("%s %s", webhook.Method, webhook.Path)
		mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			handler.HandleMessage(w, r, webhook)
		})
	}

	return nil
}
