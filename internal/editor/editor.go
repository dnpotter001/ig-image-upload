package editor

import (
	"fmt"
	"os/exec"
)

const outerHeightAndWidth int = 1080
const innerHeightAndWidth int = 1020

func AddWhiteBorder(imagePath string, outputPath string) (string, error) {
	innerDimension := fmt.Sprintf("%dx%d", innerHeightAndWidth, innerHeightAndWidth)
	outerDimension := fmt.Sprintf("%dx%d", outerHeightAndWidth, outerHeightAndWidth)

	result, err := exec.Command(
		"magick",
		imagePath,
		"-resize", innerDimension, // scale to fit, keeping aspect ratio
		"-background", "white",
		"-gravity", "center", // center the image on the canvas
		"-extent", outerDimension, // pad out to a white square
		outputPath,
	).CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("magick failed: %w: %s", err, result)
	}

	return outputPath, nil
}
