package v1

import (
	"context"
	"net/http"
)

func NewRouter(ctx context.Context, router *http.ServeMux) {
	imageHandler := New()

	router.HandleFunc("/", imageHandler.GetImage(ctx))
}
