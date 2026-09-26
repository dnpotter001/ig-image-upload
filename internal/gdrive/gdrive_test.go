package gdrive

import (
	"image/jpeg"
	"os"
	"testing"
)

func TestGetPicturesForUpload(t *testing.T) {
	files, err := GetPicturesForUpload()

	if err != nil {
		t.Fatalf("Error getting files: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No files returned")
	}

	folderExists := false
	for i := range files {
		if files[i].MimeType == "application/vnd.google-apps.folder" &&
			files[i].Name == "ig-image-upload" {
			folderExists = true
		}
	}

	if !folderExists {
		t.Fatal("ig-image-upload folder does not exist")
	}

}

func TestDownloadPicture(t *testing.T) {
	files, err := GetPicturesForUpload()

	if err != nil {
		t.Fatal(err)
	}

	download, err := DownloadImage(files[len(files)-1])

	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Create("tests/output.jpg")

	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	jpeg.Encode(f, download, &jpeg.Options{Quality: 90})

}
