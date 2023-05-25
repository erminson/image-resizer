package main

import (
	"context"
	"flag"
	"github.com/erminson/image-resizer/internal/controller/http/v1"
	"log"
	"net/http"
)

func main() {
	var maxConn int
	flag.IntVar(&maxConn, "n", 5, "number of concurrent requests")
	flag.Parse()

	log.Printf("max conn: %d", maxConn)

	ctx := context.Background()
	router := http.NewServeMux()

	v1.NewRouter(ctx, router, maxConn)

	log.Fatal(http.ListenAndServe(":8080", router))
}
