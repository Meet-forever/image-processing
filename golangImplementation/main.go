package main

import (
	"fmt"
	"image"
	"image/png"
	_ "image/png"
	"log"
	"math"
	"os"
	"sync"
)

const CHANNELS = 4

func main() {
	// img := loadImage("test.png")

	// Gray Scaling
	// grayScaleLumi(img)
	// fastGrayScale(img)

	// Convolution
	// fastGrayScaleLumi(img)
	// edgex := []float64{
	// 	-1, -2, -1,
	// 	0, 0, 0,
	// 	1, 2, 1}

	// conv_0_border(img, edgex, 3, 3, 0)
	// conv_0_border(img, edgex, 3, 3, 1)
	// conv_0_border(img, edgex, 3, 3, 2)

	// edgey := []float64{
	// 	-1, 0, 1,
	// 	-2, 0, 2,
	// 	-1, 0, 1}

	// conv_0_border(img, edgey, 3, 3, 0)
	// conv_0_border(img, edgey, 3, 3, 1)
	// conv_0_border(img, edgey, 3, 3, 2)

	// saveImage("temp.png", img)
}

func loadImage(filePath string) *image.RGBA {
	fmt.Printf("Loading %s\n", filePath)
	imgFile, err := os.Open(filePath)
	defer imgFile.Close()
	if err != nil {
		log.Println("Cannot read the file: ", err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Println("Cannot decode the file: ", err)
	}

	return img.(*image.RGBA)
}

func saveImage(filePath string, img *image.RGBA) {
	fmt.Printf("Saving %s\n", filePath)
	imgFile, err := os.Create(filePath)
	defer imgFile.Close()
	if err != nil {
		log.Println("Failed to create file: ", err)
	}
	png.Encode(imgFile, img)
}

func grayScale(img *image.RGBA) {
	rect := img.Rect
	var gray uint8 = 0
	for i := 0; i < rect.Dx()*rect.Dy()*CHANNELS; i += CHANNELS {
		gray = (img.Pix[i] + img.Pix[i+1] + img.Pix[i+2]) / 3
		img.Pix[i] = gray
		img.Pix[i+1] = gray
		img.Pix[i+2] = gray
	}
}

func grayScaleLumi(img *image.RGBA) {
	rect := img.Rect
	var lumi uint8 = 0
	for i := 0; i < rect.Dx()*rect.Dy()*CHANNELS; i += CHANNELS {
		lumi = uint8(float64(img.Pix[i])*0.2126 + float64(img.Pix[i+1])*0.7152 + float64(img.Pix[i+2])*0.0722)
		img.Pix[i] = lumi
		img.Pix[i+1] = lumi
		img.Pix[i+2] = lumi
	}
}

type factorCord struct {
	x, y, z float64
}

func updateGrayPixelChunk(arr []uint8, start int, end int, factor factorCord, wg *sync.WaitGroup) {
	defer wg.Done()
	var eta uint8 = 0
	for i := start; i < end; i += 4 {
		eta = uint8(float64(arr[i])*factor.x + float64(arr[i+1])*factor.y + float64(arr[i+2])*factor.z)
		arr[i] = eta
		arr[i+1] = eta
		arr[i+2] = eta
	}
}

// Concurrent gray scaling (average)
func fastGrayScale(img *image.RGBA) {
	rect := img.Rect
	size := rect.Dx() * rect.Dy() * CHANNELS
	splits := 5
	chunk := size / splits
	wg := new(sync.WaitGroup)
	wg.Add(splits)
	for i := 0; i < size; i += chunk {
		go updateGrayPixelChunk(img.Pix, i, int(math.Min(float64(i+chunk), float64(size))), factorCord{1 / 3, 1 / 3, 1 / 3}, wg)
	}
	wg.Wait()
}

// Concurrent gray scaling (lumination)
func fastGrayScaleLumi(img *image.RGBA) {
	rect := img.Rect
	size := rect.Dx() * rect.Dy() * CHANNELS
	splits := 5
	chunk := size / splits
	wg := new(sync.WaitGroup)
	wg.Add(splits)
	for i := 0; i < size; i += chunk {
		go updateGrayPixelChunk(img.Pix, i, int(math.Min(float64(i+chunk), float64(size))), factorCord{0.2126, 0.7152, 0.0722}, wg)
	}
	wg.Wait()
}

// Convolution Function
func conv_0_border(img *image.RGBA, ker []float64, ker_w uint32, ker_h uint32, channel uint8) {
	size := img.Bounds().Dx() * img.Bounds().Dy() * CHANNELS
	w := int(img.Bounds().Dx())
	h := int(img.Bounds().Dy())
	new_data := make([]uint8, size)
	cr, cc := ker_h/2, ker_w/2
	center := cr*ker_w + cc
	for i := uint64(channel); i < uint64(size); i += CHANNELS {
		var c float64 = 0
		for j := -int64(cr); j < int64(ker_h)-int64(cr); j += 1 {
			var row int64 = (int64(i) / CHANNELS / int64(w)) - j
			if row < 0 || row > int64(h)-1 {
				continue
			}
			for k := -int64(cc); k < int64(ker_w)-int64(cc); k += 1 {
				var col int64 = int64(i)/CHANNELS%int64(w) - k
				if col < 0 || col > int64(w)-1 {
					continue
				}
				c += ker[int64(center)+j*int64(ker_w)+k] * float64(img.Pix[(row*int64(w)+col)*CHANNELS+int64(channel)])
			}
		}
		new_data[i/CHANNELS] = uint8(c)
	}
	for i := uint64(channel); i < uint64(size); i += CHANNELS {
		img.Pix[i] = new_data[i/CHANNELS]
	}
}
