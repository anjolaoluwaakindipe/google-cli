package appclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleClientService interface {
}

type GoogleClient struct {
}

func (uc GoogleClient) saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credentials to file: %v \n", path)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		log.Fatalf("Unable to catch oauth token: %v", err)
	}

	defer f.Close()

	json.NewEncoder(f).Encode(token)
}

func (uc GoogleClient) getTokenFromCache(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var token *oauth2.Token = &oauth2.Token{}

	err = json.NewDecoder(f).Decode(token)

	if err != nil {
		return nil, err
	}

	return token, nil
}
func (uc GoogleClient) getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

func (uc GoogleClient) GetClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := uc.getTokenFromCache(tokFile)
	if err != nil {
		tok = uc.getTokenFromWeb(config)
		uc.saveToken(tokFile, tok)
	}

	return config.Client(context.Background(), tok)
}

func NewDriveClient() (*drive.Service, error) {
	ctx := context.Background()

	cred, err := os.ReadFile("./credentials.json")

	if err != nil {
		log.Fatalf("Unable to parse client seccret")
		return nil, err
	}

	config, err := google.ConfigFromJSON(cred, drive.DriveScope)

	if err != nil {
		return nil, err 
	}

	client := GoogleClient{}.GetClient(config)

	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))

	return driveService, err
}
