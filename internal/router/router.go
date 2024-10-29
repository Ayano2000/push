package router

import (
	"context"
	"fmt"
	"github.com/Ayano2000/push/internal/handlers"
	"github.com/Ayano2000/push/internal/types"
	"net/http"
	"strings"
	"sync"
)

const (
	PatternString                               = "%s %s"
	muxContextKey      types.MuxContextKey      = "router"
	urlParamContextKey types.UrlParamContextKey = "parameters"
)

type Router struct {
	mutex         sync.RWMutex
	staticRoutes  map[string]*Route
	dynamicRoutes []*Route
	handler       *handlers.Handler
	middleware    []func(http.HandlerFunc) http.HandlerFunc
}

type Route struct {
	pattern    string
	parameters []string
	segments   []string
	handler    http.HandlerFunc
	isDynamic  bool
}

func NewDynamicMux(handler *handlers.Handler) *Router {
	return &Router{
		staticRoutes:  make(map[string]*Route),
		dynamicRoutes: make([]*Route, 0),
		handler:       handler,
		middleware:    make([]func(http.HandlerFunc) http.HandlerFunc, 0),
	}
}

// parsePattern extracts parameter names and segments from a URL pattern
func parsePattern(pattern string) (parameters []string, segments []string, isDynamic bool) {
	parts := strings.Split(pattern, "/")

	for _, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			name := strings.Trim(part, "{}")
			parameters = append(parameters, name)
			segments = append(segments, "*")
			isDynamic = true
		} else {
			segments = append(segments, part)
		}
	}

	return parameters, segments, isDynamic
}

func matchRoute(route *Route, method string, path string) (map[string]string, bool) {
	pathParts := strings.Split(path, "/")

	if len(pathParts) != len(route.segments) {
		return nil, false
	}

	params := make(map[string]string)
	paramIndex := 0

	for i, segment := range route.segments {
		if segment == "*" {
			params[route.parameters[paramIndex]] = pathParts[i]
			paramIndex++
		} else if segment != pathParts[i] {
			return nil, false
		}
	}

	routeMethod := strings.Split(route.pattern, " ")[0]
	if method != routeMethod {
		return nil, false
	}

	return params, true
}

// HandleFunc registers a new route with its handler function
func (dmux *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	dmux.mutex.Lock()
	defer dmux.mutex.Unlock()

	parts := strings.Split(pattern, " ")
	if len(parts) != 2 {
		panic("invalid pattern: " + pattern)
	}
	parameters, segments, isDynamic := parsePattern(parts[1])
	route := &Route{
		pattern:    pattern,
		parameters: parameters,
		segments:   segments,
		handler:    handler,
		isDynamic:  isDynamic,
	}

	if isDynamic {
		dmux.dynamicRoutes = append(dmux.dynamicRoutes, route)
	} else {
		dmux.staticRoutes[pattern] = route
	}
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
	dmux.mutex.RUnlock()

	pattern := fmt.Sprintf(PatternString, r.Method, r.URL.Path)
	if route, exists := dmux.staticRoutes[pattern]; exists {
		ctx := context.WithValue(r.Context(), muxContextKey, dmux)
		r = r.WithContext(ctx)
		dmux.applyMiddleware(route.handler)(w, r)
		return
	}

	for _, route := range dmux.dynamicRoutes {
		if params, ok := matchRoute(route, r.Method, r.URL.Path); ok {
			ctx := context.WithValue(r.Context(), muxContextKey, dmux)
			ctx = context.WithValue(ctx, urlParamContextKey, params)
			r = r.WithContext(ctx)
			dmux.applyMiddleware(route.handler)(w, r)
			return
		}
	}

	http.NotFound(w, r)
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
