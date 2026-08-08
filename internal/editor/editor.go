package editor

import (
	"fmt"
	"os/exec"
)

const outerHeightAndWidth int = 1080
const innerHeightAndWidth int = 1020

func AddWhiteBorder(path string) string {
	out := "../tests/new.jpg"

	innerDimension := fmt.Sprintf("%dx%d", innerHeightAndWidth, innerHeightAndWidth)
	outerDimension := fmt.Sprintf("%dx%d", outerHeightAndWidth, outerHeightAndWidth)

	result, err := exec.Command(
		"magick",
		path,
		"-resize", innerDimension, // scale to fit, keeping aspect ratio
		"-background", "white",
		"-gravity", "center", // center the image on the canvas
		"-extent", outerDimension, // pad out to a white square
		out,
	).CombinedOutput()

	if err != nil {
		fmt.Println(err)
		fmt.Println(string(result))
	}

	return out
}
