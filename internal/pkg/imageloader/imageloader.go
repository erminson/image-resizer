package imageloader

import (
	"context"
	"image"
	_ "image/png"
	"net/http"
)

type Image struct {
	Image  image.Image
	Format string
}

type ImageLoader struct {
}

func New() *ImageLoader {
	return &ImageLoader{}
}

func (il *ImageLoader) Load(ctx context.Context, url string) (Image, error) {
	// use http.NewRequestWithContext() instead http.Get
	res, err := http.Get(url)
	if err != nil {
		return Image{}, err
	}
	defer res.Body.Close()

	img, format, err := image.Decode(res.Body)
	if err != nil {
		return Image{}, err
	}

	return Image{
		Image:  img,
		Format: format,
	}, nil
}
