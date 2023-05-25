package imageresizer

import (
	"bytes"
	"fmt"
	"github.com/erminson/image-resizer/internal/pkg/imageloader"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
)

type ImageResizer struct {
}

func New() *ImageResizer {
	return &ImageResizer{}
}
func (ir *ImageResizer) Resize(img imageloader.Image, height, width uint) ([]byte, error) {
	resizedImage := resize.Resize(width, height, img.Image, resize.Lanczos3)

	return ir.encodeImage(resizedImage, img.Format)
}

func (ir *ImageResizer) encodeImage(img image.Image, format string) ([]byte, error) {
	buf := new(bytes.Buffer)

	switch format {
	case "jpeg", "jpg":
		err := jpeg.Encode(buf, img, nil)
		if err != nil {
			return nil, err
		}
	case "png":
		err := png.Encode(buf, img)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknow image format: %s", format)
	}

	return buf.Bytes(), nil
}
