package main

import (
	"fmt"

	"github.com/dnpotter001/ig-image-upload/internal/editor"
)

func main() {
	fmt.Println("Ig poster starting")

	// var igAccessToken string = os.Getenv("IG_ACCESS_TOKEN")
	// var igAppSecret string = os.Getenv("IG_APP_SECRET")
	// var googleClientId string = os.Getenv("GOOGLE_CLIENT_EMAIL")
	// var googleClientId string = os.Getenv("GOOGLE_CLIENT_EMAIL")
	// var googleClientSecret string = os.Getenv("GOOGLE_PRIVATE_KEY")
	// var googleRefreshToken string = os.Getenv("GOOGLE_REFRESH_TOKEN")
	// var publicDomain string = os.Getenv("PUBLIC_DOMAIN")

	// checkAllEnvs(igAccessToken, igAppSecret, googleClientId, googleClientSecret, googleRefreshToken, publicDomain)
	// fmt.Println("All envs variables loaded")

	editor.AddWhiteBorder("../tests/000029050020.jpg", "")
}

func checkAllEnvs(igAccessToken, igAppSecret, googleClientId, googleClientSecret, googleRefreshToken, publicDomain string) {

	if igAccessToken == "" {
		panic("IG_ACCESS_TOKEN is missing")
	}
	if igAppSecret == "" {
		panic("IG_APP_SECRET is missing")
	}

	if googleClientId == "" {
		panic("GOOGLE_CLIENT_EMAIL is missing")
	}
	if googleClientSecret == "" {
		panic("GOOGLE_PRIVATE_KEY is missing")
	}
	if googleRefreshToken == "" {
		panic("GOOGLE_REFRESH_TOKEN is missing")
	}
	if publicDomain == "" {
		panic("PUBLIC_DOMAIN is missing")
	}

}
