package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	UrlKey    = "url"
	HeightKey = "height"
	WidthKey  = "width"
)

var (
	ErrUrlNotFound    = errors.New("url query parameter not found")
	ErrHeightNotFound = errors.New("height query parameter not found")
	ErrWidthNotFound  = errors.New("width query parameter not found")
)

type ImageRequest struct {
	URL    string
	Height int
	Width  int
}

type ImageHandler struct {
	// logger
	// user usecase
}

func New() *ImageHandler {
	return &ImageHandler{}
}

func (i *ImageHandler) GetImage(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		imageReq, err := i.parseQuery(values)
		if err != nil {
			return
		}

		fmt.Printf("%+v", imageReq)
	}
}

func (i *ImageHandler) parseQuery(v url.Values) (*ImageRequest, error) {
	urls, ok := v[UrlKey]
	if !ok || len(urls[0]) == 0 {
		return nil, ErrUrlNotFound
	}

	heights, ok := v[HeightKey]
	if !ok || len(heights[0]) == 0 {
		return nil, ErrHeightNotFound
	}

	widths, ok := v[WidthKey]
	if !ok || len(widths[0]) == 0 {
		return nil, ErrWidthNotFound
	}

	url := urls[0]

	height, err := strconv.Atoi(heights[0])
	if err != nil {
		return nil, err
	}

	width, err := strconv.Atoi(widths[0])
	if err != nil {
		return nil, err
	}

	return &ImageRequest{
		URL:    url,
		Height: height,
		Width:  width,
	}, nil
}
