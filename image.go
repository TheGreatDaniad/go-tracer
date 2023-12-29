package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func CreateSampleImage() {
	width := 1024
	height := 1024
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill the image with the specified color
	for y := 0; y < width; y++ {
		for x := 0; x < height; x++ {
			r := uint8((x / 4) % 255)
			g := uint8((y / 4) % 255)
			b := uint8((x * y / 16) % 255)
			c := color.RGBA{r, g, b, 255}
			img.Set(x, y, c)
		}
	}

}

func SaveImage(img image.Image, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
