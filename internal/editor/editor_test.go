package editor

import (
	"os/exec"
	"testing"
)

func TestIsImageMagicInstalled(t *testing.T) {
	_, err := exec.Command("magick", "--version").CombinedOutput()

	if err != nil {
		t.Errorf("Check if image magick installation")
	}
}

func TestAddedBorderRuns(t *testing.T) {
	_, err := AddWhiteBorder("tests/input.jpg", "tests/output.jpg")

	if err != nil {
		t.Fatalf("AddWhiteBorder returned an error: %v", err)
	}

}
