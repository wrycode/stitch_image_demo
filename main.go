package main

import (
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/rasterizer"	
	"image"
	"image/draw"
	"image/png"
	"os"
	// "fmt"
)
func convertCanvasToImage(c *canvas.Canvas) image.Image {
	var res canvas.Resolution
	res = canvas.Resolution(1.0)
	col := canvas.LinearColorSpace{}
	img := rasterizer.Draw(c, res, col)
	return img
}

func main() {
	images := make([]image.Image, 0)	

	for i := 0; i < 12; i++ {
		// Create new canvas of dimension 100x100 mm
		c := canvas.New(100, 100)

		// Create a canvas context used to keep drawing state
		ctx := canvas.NewContext(c)

		// Create a triangle path from an SVG path and draw it to the canvas
		triangle, err := canvas.ParseSVGPath("L60 0L30 60z")
		if err != nil {
			panic(err)
		}
		ctx.SetFillColor(canvas.Mediumseagreen)
		ctx.DrawPath(20, 20, triangle)
		var image image.Image
		image = convertCanvasToImage(c)
		images = append(images, image)
	}

	// start stitching
	upperLeftX := images[0].Bounds().Max.X
	upperLeftY := images[0].Bounds().Max.Y
	imgWidth := upperLeftX * len(images)
	imgHeight := upperLeftY

	// create new blank image with a size that depends on number of images
	newImage := image.NewNRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Start drawing images from the slice to new blank image
	for i, img := range images {
		rect := image.Rect(upperLeftX*i, 0, upperLeftX*i+upperLeftX, imgHeight)
		draw.Draw(newImage, rect, img, img.Bounds().Min, draw.Src)
	}

	// Create resulting image file on disk.
	imgFile, err := os.Create("tiled.png")
	if err != nil {
		panic(err)
	}
	defer imgFile.Close()

	// Encode writes the Image m to w in PNG format.
	err = png.Encode(imgFile, newImage)
	if err != nil {
		panic(err)
	}
	
}
