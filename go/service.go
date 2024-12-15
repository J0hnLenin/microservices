package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
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
	logFile, err := os.OpenFile("golang.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	if request.Method != http.MethodGet {
		message := "Invalid http method"
		http.Error(writer, message, http.StatusMethodNotAllowed)
		fmt.Println(message)
		logger.Println(message)
		return
	}

	paramValue := request.Header.Get("size")
	if paramValue == "" {
		message := "Size parameter not found"
		http.Error(writer, message, http.StatusBadRequest)
		fmt.Println(message)
		logger.Println(message)
		return
	}

	size, err := strconv.Atoi(paramValue)
	if err != nil || size <= 0 {
		message := "Invalid size parameter value"
		http.Error(writer, message, http.StatusBadRequest)
		fmt.Println(message)
		logger.Println(message)
		return
	}

	imageData := generateImage(size)

	message := "image with size %s created\n"
	fmt.Printf(message, paramValue)
	logger.Printf(message, paramValue)

	writer.Header().Set("Content-Type", "image/png")
	png.Encode(writer, imageData)
}

func main() {
	http.HandleFunc("/generate-image", imageHandler)

	fmt.Println("server started at port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
