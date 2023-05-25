package v1

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/erminson/image-resizer/internal/pkg/imageloader"
	"github.com/erminson/image-resizer/internal/pkg/imageresizer"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"
)

const (
	UrlKey    = "url"
	HeightKey = "height"
	WidthKey  = "width"
)

var (
	ErrUrlNotFound    = errors.New("url query parameter not found")
	ErrUrlInvalid     = errors.New("url query parameter invalid")
	ErrHeightNotFound = errors.New("height query parameter not found")
	ErrHeightInvalid  = errors.New("height query parameter must be an integer greater than zero")
	ErrWidthNotFound  = errors.New("width query parameter not found")
	ErrWidthInvalid   = errors.New("width query parameter must be an integer greater than zero")
)

type ApiError struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
}

type JsonErrorResponse struct {
	Error *ApiError `json:"error"`
}

type ImageRequest struct {
	URL    string
	Height uint
	Width  uint
}

type ImageHandler struct {
	il       *imageloader.ImageLoader
	ir       *imageresizer.ImageResizer
	maxConn  int64
	currConn atomic.Int64
}

func New(loader *imageloader.ImageLoader, resizer *imageresizer.ImageResizer, maxConn int) *ImageHandler {
	return &ImageHandler{
		il:      loader,
		ir:      resizer,
		maxConn: int64(maxConn),
	}
}

func (i *ImageHandler) GetImage(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if i.maxConn <= i.currConn.Load() {
			apiErr := &ApiError{Status: http.StatusBadRequest, Title: "max connection"}
			i.errorResponse(w, apiErr)
			return
		}

		i.currConn.Add(1)
		defer i.currConn.Add(-1)

		values := r.URL.Query()
		imageReq, err := i.parseQuery(values)
		if err != nil {
			apiErr := &ApiError{Status: http.StatusBadRequest, Title: err.Error()}
			i.errorResponse(w, apiErr)
			return
		}

		log.Println("Loading...")
		img, err := i.il.Load(ctx, imageReq.URL)
		if err != nil {
			apiErr := &ApiError{Status: http.StatusBadRequest, Title: err.Error()}
			i.errorResponse(w, apiErr)
			return
		}
		log.Println("Loaded.")

		log.Println("Resizing...")
		resizesData, err := i.ir.Resize(img, imageReq.Height, imageReq.Width)
		if err != nil {
			apiErr := &ApiError{Status: http.StatusBadRequest, Title: err.Error()}
			i.errorResponse(w, apiErr)
			return
		}
		log.Println("Resized.")

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(resizesData)
	}
}

func (i *ImageHandler) errorResponse(w http.ResponseWriter, apiErr *ApiError) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(apiErr.Status)

	response := JsonErrorResponse{Error: apiErr}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatal(err)
	}
}

func (i *ImageHandler) parseQuery(v url.Values) (ImageRequest, error) {
	urls, ok := v[UrlKey]
	if !ok || len(urls[0]) == 0 {
		return ImageRequest{}, ErrUrlNotFound
	}

	heights, ok := v[HeightKey]
	if !ok || len(heights[0]) == 0 {
		return ImageRequest{}, ErrHeightNotFound
	}

	widths, ok := v[WidthKey]
	if !ok || len(widths[0]) == 0 {
		return ImageRequest{}, ErrWidthNotFound
	}

	imageUrlEncode := urls[0]
	imageUrl, err := url.QueryUnescape(imageUrlEncode)
	if err != nil {
		return ImageRequest{}, ErrUrlInvalid
	}

	height, err := strconv.Atoi(heights[0])
	if err != nil || height <= 0 {
		return ImageRequest{}, ErrHeightInvalid
	}

	width, err := strconv.Atoi(widths[0])
	if err != nil || width <= 0 {
		return ImageRequest{}, ErrWidthInvalid
	}

	return ImageRequest{
		URL:    imageUrl,
		Height: uint(height),
		Width:  uint(width),
	}, nil
}
