package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	impactFont *truetype.Font
	drawer     font.Drawer
)

func init() {
	impactData, err := ioutil.ReadFile("impact.ttf")
	if err != nil {
		panic(err)
	}

	impactFont, err = freetype.ParseFont(impactData)
	if err != nil {
		panic(err)
	}
}

func imageWidth(i image.Image) int {
	return i.Bounds().Max.X - i.Bounds().Min.X
}

func imageHeight(i image.Image) int {
	return i.Bounds().Max.Y - i.Bounds().Min.Y
}

func overlayPipe(src io.Reader, dpt string, txt string, dst io.Writer) error {
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
	var finalImage *image.NRGBA

	if txt != "" {
		finalImage = imaging.New(imageWidth(overlayed), imageHeight(overlayed)+100, color.White)
		finalImage = imaging.Overlay(finalImage, overlayed, image.Point{X: 0, Y: 0}, 1.0)

		pos := fixed.Point26_6{X: fixed.Int26_6(0), Y: fixed.Int26_6((imageHeight(finalImage) - 30) * 64)}

		d := font.Drawer{
			Dst: finalImage,
			Src: image.Black,
			Face: truetype.NewFace(impactFont, &truetype.Options{
				Size: 40,
				DPI:  72,
			}),
			Dot: pos,
		}

		textSize := d.MeasureString(txt)
		d.Dot.X = (fixed.Int26_6(imageWidth(dptImage)) / 2 * 64) - (textSize / 2)

		d.DrawString(txt)
	} else {
		finalImage = overlayed
	}

	png.Encode(dst, finalImage)
	return nil
}
