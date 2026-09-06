package gdrive

import "testing"

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

}
