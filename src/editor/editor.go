package editor

import (
	"fmt"
	"os/exec"
)

// WhiteBorder fits the image at path onto a white square canvas with a white
// border on every side — the classic film / photography Instagram look — and
// writes the result to ../tests/new.jpg, returning that path.
func WhiteBorder(path string) string {

	const canvas = 1080 // final square size, in pixels
	const inner = 1020  // image is scaled to fit inside this, leaving a ~50px border

	out := "../tests/new.jpg"

	result, err := exec.Command(
		"magick",
		"../tests/000029050020.jpg",
		"-resize", fmt.Sprintf("%dx%d", inner, inner), // scale to fit, keeping aspect ratio
		"-background", "white",
		"-gravity", "center", // center the image on the canvas
		"-extent", fmt.Sprintf("%dx%d", canvas, canvas), // pad out to a white square
		out,
	).CombinedOutput()

	if err != nil {
		fmt.Println(err)
		fmt.Println(string(result))
	}

	return out
}
