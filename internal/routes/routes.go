package routes

import (
	"context"
	"fmt"
	"github.com/Ayano2000/push/internal/handlers"
	"github.com/Ayano2000/push/internal/types"
	"net/http"
	"sync"
)

const PatternString = "%s %s"

type DynamicMux struct {
	mutex      sync.RWMutex
	routes     map[string]http.HandlerFunc
	handler    *handlers.Handler
	middleware []func(http.HandlerFunc) http.HandlerFunc
}

func NewDynamicMux(handler *handlers.Handler) *DynamicMux {
	return &DynamicMux{
		routes:     make(map[string]http.HandlerFunc),
		handler:    handler,
		middleware: make([]func(http.HandlerFunc) http.HandlerFunc, 0),
	}
}

func (dmux *DynamicMux) HandleFunc(pattern string, handler http.HandlerFunc) {
	dmux.mutex.Lock()
	defer dmux.mutex.Unlock()

	dmux.routes[pattern] = handler
}

func (dmux *DynamicMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dmux.mutex.RLock()
	pattern := fmt.Sprintf(PatternString, r.Method, r.URL.Path)
	handler, exists := dmux.routes[pattern]
	dmux.mutex.RUnlock()

	if !exists {
		http.NotFound(w, r)
		return
	}
	handler(w, r)
}

func (dmux *DynamicMux) RegisterWebhook(webhook types.Webhook) {
	pattern := fmt.Sprintf(PatternString, webhook.Method, webhook.Path)
	dmux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		dmux.handler.HandleMessage(w, r, webhook)
	})
}

// Use adds methods to append to middleware
func (dmux *DynamicMux) Use(middleware ...func(http.HandlerFunc) http.HandlerFunc) {
	dmux.middleware = append(dmux.middleware, middleware...)
}

// Apply middleware to a handler
func (dmux *DynamicMux) applyMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	// Reverse order so the first middleware in the chain is the outermost wrapper
	for i := len(dmux.middleware) - 1; i >= 0; i-- {
		handler = dmux.middleware[i](handler)
	}
	return handler
}

func RegisterRoutes(handler *handlers.Handler) (*DynamicMux, error) {
	dmux := NewDynamicMux(handler)

	// Register static routes
	dmux.HandleFunc("POST /webhooks", handler.CreateWebhook)
	dmux.HandleFunc("GET /webhooks", handler.GetWebhooks)
	dmux.HandleFunc("GET /webhooks/{name}/content", handler.GetWebhookContent)
	dmux.HandleFunc("DELETE /webhooks/{name}", handler.DeleteWebhook)
	dmux.HandleFunc("DELETE /webhooks/{name}/content", handler.DeleteWebhookContents)

	// Register existing webhooks
	var webhooks []types.Webhook
	webhooks, err := handler.Services.DB.GetWebhooks(context.Background())
	if err != nil {
		return nil, err
	}

	for _, webhook := range webhooks {
		dmux.RegisterWebhook(webhook)
	}

	return dmux, nil
}
