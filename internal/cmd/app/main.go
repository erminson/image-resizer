package main

import (
	"context"
	"github.com/erminson/image-resizer/internal/controller/http/v1"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	r := &http.ServeMux{}
	v1.NewRouter(ctx, r)

	log.Fatal(http.ListenAndServe(":80", r))
}
