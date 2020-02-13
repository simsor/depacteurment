package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/disintegration/imaging"
)

func imageWidth(i image.Image) int {
	return i.Bounds().Max.X - i.Bounds().Min.X
}

func imageHeight(i image.Image) int {
	return i.Bounds().Max.Y - i.Bounds().Min.Y
}

func overlayPipe(src io.Reader, dpt string, dst io.Writer) error {
	srcImg, _, err := image.Decode(src)
	if err != nil {
		return err
	}

	dptF, err := os.Open(fmt.Sprintf("dpt_images/dpt%s.png", dpt))
	if err != nil {
		return err
	}

	dptImage, _, err := image.Decode(dptF)
	if err != nil {
		return err
	}

	resizedSrc := imaging.Resize(srcImg, 500, 0, imaging.Lanczos)
	croppedSrc := imaging.CropCenter(resizedSrc, imageWidth(dptImage), imageHeight(dptImage))
	overlayed := imaging.Overlay(croppedSrc, dptImage, image.Point{X: 0, Y: 0}, 1.0)

	png.Encode(dst, overlayed)
	return nil
}
