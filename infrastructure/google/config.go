package google

import (
	"fmt"
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google-photos-sync/domain/model"
	"io/ioutil"
)

func GetAuthUrl(accountType model.AccountType) (string, error) {
	authConfig, err := GetAuthConfig(accountType)
	if err != nil {
		return "", err
	}

	return authConfig.AuthCodeURL(string(accountType), oauth2.AccessTypeOffline), err
}

func GetAuthConfig(accountType model.AccountType) (*oauth2.Config, error) {
	var scope = getScopeByAccountType(accountType)

	credentialsJson, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %w", err)
	}

	config, err := google.ConfigFromJSON(credentialsJson, scope...)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %w", err)
	}

	return config, err
}

func getScopeByAccountType(accountType model.AccountType) []string {
	switch accountType {
	case model.To:
		return []string{gphotos.PhotoslibraryScope}
	default:
		return []string{gphotos.PhotoslibraryReadonlyScope}
	}
}
