package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"

	svgo "github.com/ajstarks/svgo"
)

const (
	REGION_SIZE      = 30
	IMAGE_PATH       = "../imagep/zoom.jpg"
	OUTPUT_SVG_PATH  = "../imagep/output.svg"
	MAX_DOT_SIZE     = REGION_SIZE / 5 * 4
	THRESHOLD_NO_DOT = MAX_DOT_SIZE / 50
	A5_WIDTH         = 148.5 // in mm
	A5_HEIGHT        = 210.0 // in mm
)

func main() {
	// Load the image
	img, err := loadImage(IMAGE_PATH)
	if err != nil {
		log.Fatalf("Failed to load image: %v", err)
	}

	// Ensure the image is in portrait orientation
	img = ensurePortrait(img)

	// Convert image to grayscale
	grayImg := convertToGrayscale(img)

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

// rotateImage90 rotates an image by 90 degrees clockwise.
func rotateImage90(img image.Image) image.Image {
	bounds := img.Bounds()
	rotatedImg := image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rotatedImg.Set(bounds.Dy()-y-1, x, img.At(x, y))
		}
	}

	return rotatedImg
}

// ensurePortrait ensures the image is in portrait orientation.
func ensurePortrait(img image.Image) image.Image {
	bounds := img.Bounds()
	if bounds.Dx() > bounds.Dy() {
		return rotateImage90(img)
	}
	return img
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

	new_path := strings.Replace(savePath, ".svg", "_filled.svg", -1)
	file_filled, err := os.Create(new_path)
	if err != nil {
		log.Fatalf("Failed to create SVG file: %v", err)
	}
	defer file_filled.Close()

	canvasHollow := svgo.New(file)
	canvasFilled := svgo.New(file_filled)
	canvasHollow.Start(int(A5_WIDTH*10), int(A5_HEIGHT*10)) // SVGO uses units in pixels
	canvasFilled.Start(int(A5_WIDTH*10), int(A5_HEIGHT*10)) // SVGO uses units in pixels

	// Set white background
	canvasHollow.Rect(0, 0, int(A5_WIDTH*10), int(A5_HEIGHT*10), "fill:white")
	canvasFilled.Rect(0, 0, int(A5_WIDTH*10), int(A5_HEIGHT*10), "fill:white")

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
			canvasHollow.Circle(int(x), int(y), int(radius), "fill:black;stroke:black")
			canvasFilled.Circle(int(x), int(y), int(radius), "fill:none;stroke:black")
		}
	}

	canvasHollow.End()
	canvasFilled.End()
	log.Printf("SVG saved as %s.", savePath)
}
