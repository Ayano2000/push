package router

import (
	"context"
	"fmt"
	"github.com/Ayano2000/push/internal/handlers"
	"github.com/Ayano2000/push/internal/types"
	"net/http"
	"sync"
)

const (
	PatternString                     = "%s %s"
	muxContextKey types.MuxContextKey = "router"
)

type Router struct {
	mutex      sync.RWMutex
	routes     map[string]http.HandlerFunc
	handler    *handlers.Handler
	middleware []func(http.HandlerFunc) http.HandlerFunc
}

func NewDynamicMux(handler *handlers.Handler) *Router {
	return &Router{
		routes:     make(map[string]http.HandlerFunc),
		handler:    handler,
		middleware: make([]func(http.HandlerFunc) http.HandlerFunc, 0),
	}
}

// HandleFunc registers a new route with its handler function
func (dmux *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	dmux.mutex.Lock()
	defer dmux.mutex.Unlock()

	dmux.routes[pattern] = handler
}

// RegisterWebhook adds a new webhook route dynamically
func (dmux *Router) RegisterWebhook(webhook types.Webhook) {
	pattern := fmt.Sprintf(PatternString, webhook.Method, webhook.Path)
	dmux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		dmux.handler.HandleMessage(w, r, webhook)
	})
}

// ServeHTTP implements the http.Handler interface
func (dmux *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dmux.mutex.RLock()
	pattern := fmt.Sprintf(PatternString, r.Method, r.URL.Path)
	handler, exists := dmux.routes[pattern]
	dmux.mutex.RUnlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	ctx := context.WithValue(r.Context(), muxContextKey, dmux)
	r = r.WithContext(ctx)

	// Apply middleware to the handler if any
	handlerWithMiddleware := dmux.applyMiddleware(handler)
	handlerWithMiddleware(w, r)
}

// Use appends the given functions to middleware
func (dmux *Router) Use(middleware ...func(http.HandlerFunc) http.HandlerFunc) {
	dmux.middleware = append(dmux.middleware, middleware...)
}

// applyMiddleware will invoke the dmux.middleware functions in reverse order
// so the first middleware in the chain is the outermost wrapper
func (dmux *Router) applyMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	for i := len(dmux.middleware) - 1; i >= 0; i-- {
		handler = dmux.middleware[i](handler)
	}
	return handler
}

func RegisterRoutes(handler *handlers.Handler) (*Router, error) {
	dmux := NewDynamicMux(handler)

	// Register static router
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

var _ types.WebhookRegistrar = (*Router)(nil)
