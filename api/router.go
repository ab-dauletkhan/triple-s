package api

import (
	"net/http"

	"github.com/ab-dauletkhan/triple-s/api/handlers"
)

func Routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Bucket handling
	mux.HandleFunc("PUT /{BucketName}", handlers.CreateBucket)
	mux.HandleFunc("GET /", handlers.ListBuckets)
	mux.HandleFunc("DELETE /{BucketName}", handlers.DeleteBucket)

	// Object handling
	mux.HandleFunc("PUT /{BucketName}/{ObjectKey}", handlers.CreateObject)
	mux.HandleFunc("GET /{BucketName}/{ObjectKey}", handlers.GetObject)
	mux.HandleFunc("DELETE /{BucketName}/{ObjectKey}", handlers.DeleteObject)
	mux.HandleFunc("GET /{BucketName}", handlers.ListObjects)

	return mux
}
