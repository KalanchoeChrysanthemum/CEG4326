package vision

import (
    "fmt"
    "gocv.io/x/gocv"
    "image"
    "image/color"
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

    gocv.GaussianBlur(cropped, &cropped, image.Pt(11, 11), 0, 0, gocv.BorderDefault)

    // Threshold the cropped image
    gocv.Threshold(cropped, &cropped, 160, 255, gocv.ThresholdBinary)

    // Grid size
    rows, cols := 4, 4

    // Image dimensions
    height := cropped.Rows()
    width := cropped.Cols()

    // Size of each cell
    cellH := height / rows
    cellW := width / cols

    binary := ""

    font := gocv.FontHersheySimplex

    for i := 0; i < rows; i++ {
        for j := 0; j < cols; j++ {
            x := j * cellW
            y := i * cellH

            roi := cropped.Region(imageRect(x, y, cellW, cellH))
            defer roi.Close()

            avgBrightness := computeMeanBrightness(roi)

            var status string
            if avgBrightness > 20 {
                status = "1"
            } else {
                status = "0"
            }
            binary += status

            // Pick rectangle color
            var rectColor color.RGBA
            if avgBrightness > 150 {
                rectColor = color.RGBA{G: 255, A: 255} // Green
            } else {
                rectColor = color.RGBA{R: 255, A: 255} // Red
            }

            // Draw rectangle
            gocv.Rectangle(&cropped, imageRect(x, y, cellW, cellH), color.RGBA{255, 255, 255, 255}, 2)

            // Put text
            text := fmt.Sprintf("%d", avgBrightness)
            textSize := gocv.GetTextSize(text, font, 0.5, 1)
            textX := x + (cellW-textSize.X)/2
            textY := y + (cellH+textSize.Y)/2
            gocv.PutText(&cropped, text, image.Pt(textX, textY), font, 0.5, rectColor, 1)
        }
        binary += "\n"
    }

    // Save processed image
    ok := gocv.IMWrite(outputPath, cropped)
    if !ok {
        return "", fmt.Errorf("error saving image: %s", outputPath)
    }

    return binary, nil
}

// Helper to create a rectangle
func imageRect(x, y, width, height int) image.Rectangle {
    return image.Rect(x, y, x+width, y+height)
}

// computeMeanBrightness calculates the average brightness of the ROI manually
func computeMeanBrightness(roi gocv.Mat) int {
    sum := 0.0
    count := roi.Rows() * roi.Cols()

    for y := 0; y < roi.Rows(); y++ {
        for x := 0; x < roi.Cols(); x++ {
            pixel := roi.GetUCharAt(y, x)
            sum += float64(pixel)
        }
    }

    if count == 0 {
        return 0
    }
    return int(sum / float64(count))
}

