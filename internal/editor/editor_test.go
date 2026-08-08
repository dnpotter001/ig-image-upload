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
