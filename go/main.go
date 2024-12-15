package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	svgo "github.com/ajstarks/svgo"
)

const (
	REGION_SIZE      = 10
	IMAGE_PATH       = "../imagep/image.jpg"
	OUTPUT_SVG_PATH  = "output.svg"
	GRAYSCALE_PATH   = "grayscale.png" // Path to save the grayscale image
	MAX_DOT_SIZE     = REGION_SIZE
	THRESHOLD_NO_DOT = MAX_DOT_SIZE / 10
	A5_WIDTH         = 148.5 // in mm
	A5_HEIGHT        = 210.0 // in mm
)

func main() {
	// Load the image
	img, err := loadImage(IMAGE_PATH)
	if err != nil {
		log.Fatalf("Failed to load image: %v", err)
	}

	// Convert image to grayscale
	grayImg := convertToGrayscale(img)

	// Save the grayscale image to a file
	err = saveGrayscaleImage(grayImg, GRAYSCALE_PATH)
	if err != nil {
		log.Fatalf("Failed to save grayscale image: %v", err)
	}
	log.Printf("Grayscale image saved as %s.", GRAYSCALE_PATH)

	// Generate the dot array
	dotArray := imageToVaryingDots(grayImg, REGION_SIZE)

	// Generate SVG with varying dot sizes
	plotAndSaveVaryingDotsSVG(dotArray, REGION_SIZE, OUTPUT_SVG_PATH, img.Bounds().Size())
}

// loadImage loads an image from a file path.
func loadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// convertToGrayscale converts an image to grayscale.
func convertToGrayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayImg.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}

	return grayImg
}

// saveGrayscaleImage saves a grayscale image to a file.
func saveGrayscaleImage(img *image.Gray, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}

// imageToVaryingDots converts a grayscale image to an array of dots with varying sizes.
func imageToVaryingDots(img *image.Gray, regionSize int) [][]float64 {
	bounds := img.Bounds()
	width := bounds.Dx() / regionSize
	height := bounds.Dy() / regionSize
	dotArray := make([][]float64, height)

	for i := 0; i < height; i++ {
		dotArray[i] = make([]float64, width)
		for j := 0; j < width; j++ {
			var sum int
			for y := 0; y < regionSize; y++ {
				for x := 0; x < regionSize; x++ {
					sum += int(img.GrayAt(j*regionSize+x, i*regionSize+y).Y)
				}
			}
			dotArray[i][j] = float64(sum) / float64(regionSize*regionSize)
		}
	}

	return dotArray
}

// plotAndSaveVaryingDotsSVG generates and saves an SVG with varying dot sizes.
func plotAndSaveVaryingDotsSVG(dotArray [][]float64, regionSize int, savePath string, imageSize image.Point) {
	file, err := os.Create(savePath)
	if err != nil {
		log.Fatalf("Failed to create SVG file: %v", err)
	}
	defer file.Close()

	canvas := svgo.New(file)
	canvas.Start(int(A5_WIDTH*10), int(A5_HEIGHT*10)) // SVGO uses units in pixels

	// Set white background
	canvas.Rect(0, 0, int(A5_WIDTH*10), int(A5_HEIGHT*10), "fill:white")

	originalWidth := float64(imageSize.X)
	originalHeight := float64(imageSize.Y)
	aspectRatio := originalWidth / originalHeight

	var figWidth, figHeight float64
	if aspectRatio > (A5_WIDTH / A5_HEIGHT) {
		figWidth = A5_WIDTH * 10
		figHeight = (A5_WIDTH * 10) / aspectRatio
	} else {
		figHeight = A5_HEIGHT * 10
		figWidth = (A5_HEIGHT * 10) * aspectRatio
	}

	xOffset := (A5_WIDTH*10 - figWidth) / 2
	yOffset := (A5_HEIGHT*10 - figHeight) / 2

	for i := range dotArray {
		for j := range dotArray[i] {
			dotSize := (255 - dotArray[i][j]) / 255 * MAX_DOT_SIZE
			if dotSize < THRESHOLD_NO_DOT {
				continue
			}
			x := float64(j*regionSize)*figWidth/originalWidth + xOffset
			y := float64(i*regionSize)*figHeight/originalHeight + yOffset
			radius := dotSize / 2 * figWidth / originalWidth
			canvas.Circle(int(x), int(y), int(radius), "fill:none;stroke:black")
		}
	}

	canvas.End()
	log.Printf("SVG saved as %s.", savePath)
}
