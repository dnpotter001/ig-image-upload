package gdrive


import (
	"os"
)

const apiBase string = "https://www.googleapis.com/drive/v3/files"

func getPicturesForUpload() []string, error {

	var googleClientId string = os.Getenv("GOOGLE_CLIENT_ID")
    var googleClientSecret string = os.Getenv("GOOGLE_CLIENT_SECRET")

}

func DownloadPicture(fileName string) bytes[], error {

}
