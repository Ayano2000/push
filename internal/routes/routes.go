package routes

import (
	"context"
	"fmt"
	"github.com/Ayano2000/push/internal/handlers"
	"github.com/Ayano2000/push/internal/types"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *handlers.Handler) error {
	mux.HandleFunc("POST /endpoint", handler.CreateBucket)

	// custom endpoints
	var buckets []types.Bucket
	err := handler.Services.DB.QueryRow(context.Background(), "SELECT * FROM buckets").Scan(buckets)
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
