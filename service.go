package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
)

func generateImage() image.Image {
	size := 28
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			с := uint8(rand.Intn(256))
			img.Set(x, y, color.RGBA{с, с, с, 255})
		}
	}

	return img
}

func main() {
	fmt.Println("hello world!")
	img := generateImage()
	file, err := os.Create("go_image.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}
