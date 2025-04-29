package vision

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"strings"
)

// ProcessImage loads an image, processes it, and saves the output.
// It returns the binary string result.
func ProcessImage(inputPath, outputPath string) (string, error) {
	// Load the image in grayscale
	img := gocv.IMRead(inputPath, gocv.IMReadGrayScale)
	if img.Empty() {
		return "", fmt.Errorf("error reading image: %s", inputPath)
	}
	defer img.Close()

	cropX := 310
	cropY := 280
	cropW := 25
	cropH := 25

	// Crop the image (x=600, y=500, width=1000, height=850)
	cropped := img.Region(imageRect(cropX, cropY, cropX+cropW, cropY+cropH))
	defer cropped.Close()

	gocv.GaussianBlur(cropped, &cropped, image.Pt(15, 15), 0, 0, gocv.BorderDefault)

	// Threshold the cropped image
	gocv.Threshold(cropped, &cropped, 200, 255, gocv.ThresholdBinary)

	// Grid size
	rows, cols := 4, 4

	// Image dimensions
	height := cropped.Rows()
	width := cropped.Cols()

	// Size of each cell
	cellH := height / rows
	cellW := width / cols

	var builder strings.Builder
	builder.Grow(16)

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
				x := j * cellW
				y := i * cellH

				roi := cropped.Region(imageRect(x, y, cellW, cellH))
				avgBrightness := computeMeanBrightness(roi)
				roi.Close()

				if avgBrightness > 15 {
					builder.WriteByte('1')
				} else {
					builder.WriteByte('0')
				}

				gocv.Rectangle(&cropped, imageRect(x, y, cellW, cellH), color.RGBA{255, 255, 255, 255}, 2)
		}
	}

	//  Save processed image
	ok := gocv.IMWrite(outputPath, cropped)
	if !ok {
		return "", fmt.Errorf("error saving image: %s", outputPath)
	}

	return builder.String(), nil
}

// Helper to create a rectangle
func imageRect(x, y, width, height int) image.Rectangle {
	return image.Rect(x, y, x+width, y+height)
}

// computeMeanBrightness calculates the average brightness of the ROI manually
func computeMeanBrightness(roi gocv.Mat) int {
	return int(roi.Mean().Val1)
}
