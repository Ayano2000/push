package routes

import (
	"context"
	"fmt"
	"github.com/Ayano2000/push/internal/handlers"
	"github.com/Ayano2000/push/internal/types"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *handlers.Handler) error {
	mux.HandleFunc("POST /buckets", handler.CreateBucket)
	mux.HandleFunc("GET /buckets", handler.ListBuckets)
	mux.HandleFunc("GET /buckets/{name}", handler.GetBucketContents)

	// custom endpoints
	var buckets []types.Bucket
	buckets, err := handler.Services.DB.GetBuckets(context.Background())
	if err != nil {
		return err
	}

	for _, bucket := range buckets {
		pattern := fmt.Sprintf("%s %s", bucket.Method, bucket.Path)
		mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			handler.HandleRequest(w, r, bucket)
		})
	}

	return nil
}
