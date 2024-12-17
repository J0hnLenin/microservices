package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func generateImage(size int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	var wg sync.WaitGroup

	for y := 0; y < size; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := 0; x < size; x++ {
				r := uint8(rand.Intn(256))
				g := uint8(rand.Intn(256))
				b := uint8(rand.Intn(256))
				img.Set(x, y, color.RGBA{r, g, b, 255})
			}
		}(y)
	}
	wg.Wait()
	return img
}

func imageHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		message := "Invalid http method"
		http.Error(writer, message, http.StatusMethodNotAllowed)
		timenow := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println(timenow, "ERROR:", message)
		return
	}

	paramValue := request.Header.Get("size")
	if paramValue == "" {
		message := "Size parameter not found"
		http.Error(writer, message, http.StatusBadRequest)
		timenow := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println(timenow, "ERROR:", message)
		return
	}

	size, err := strconv.Atoi(paramValue)
	if err != nil || size <= 0 {
		message := "Invalid size parameter value"
		http.Error(writer, message, http.StatusBadRequest)
		timenow := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println(timenow, "ERROR:", message)
		return
	}

	imageData := generateImage(size)

	message := "created image with size"
	timenow := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(timenow, "INFO:", message, paramValue)

	writer.Header().Set("Content-Type", "image/png")
	png.Encode(writer, imageData)
}

func main() {
	http.HandleFunc("/generate-image", imageHandler)
	timenow := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(timenow, "INFO:", "server started at port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		timenow := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println(timenow, "ERROR:", err)
	}
}
