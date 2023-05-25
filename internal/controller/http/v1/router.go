package v1

import (
	"context"
	"github.com/erminson/image-resizer/internal/pkg/imageloader"
	"github.com/erminson/image-resizer/internal/pkg/imageresizer"
	"net/http"
)

func NewRouter(ctx context.Context, router *http.ServeMux, maxConn int) {
	loader := imageloader.New()
	resizer := imageresizer.New()
	imageHandler := New(loader, resizer, maxConn)

	router.HandleFunc("/", imageHandler.GetImage(ctx))
}
