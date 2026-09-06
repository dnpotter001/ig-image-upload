package gdrive

import (
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const JWTTokenURL = "https://oauth2.googleapis.com/token"

func GetPicturesForUpload() ([]*drive.File, error) {
	service, err := getDriveService()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	fileList, err := service.Files.List().PageSize(100).Do()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return fileList.Files, nil
}

func getDriveService() (*drive.Service, error) {
	clientEmail := os.Getenv("GOOGLE_CLIENT_EMAIL")
	if len(clientEmail) == 0 {
		return nil, fmt.Errorf("GOOGLE_CLIENT_EMAIL is empty")
	}

	privateKey := os.Getenv("GOOGLE_PRIVATE_KEY")
	if len(privateKey) == 0 {
		return nil, fmt.Errorf("GOOGLE_PRIVATE_KEY is empty")
	}

	config := &jwt.Config{
		Email:      clientEmail,
		PrivateKey: []byte(strings.ReplaceAll(privateKey, "\\n", "\n")),
		Scopes:     []string{"https://www.googleapis.com/auth/drive.readonly"},
		TokenURL:   JWTTokenURL,
	}

	ctx := context.Background()

	driveService, err := drive.NewService(ctx, option.WithHTTPClient(config.Client(ctx)))
	if err != nil {
		return nil, err
	}

	return driveService, nil
}

func DownloadImage(fileName string) ([]byte, error) {
	return nil, nil
}
