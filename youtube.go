package Youtube

import (
	google "github.com/leapforce-libraries/go_google"
)

const (
	APIName string = "Youtube"
	APIURL  string = "https://youtubeanalytics.googleapis.com/v2"
)

// Youtube stores Youtube configuration
//
type Youtube struct {
	Client *google.GoogleClient
}

// methods
//
func NewYoutube(clientID string, clientSecret string, scope string, bigQuery *google.BigQuery) *Youtube {
	config := google.GoogleClientConfig{
		APIName:      APIName,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        scope,
	}

	googleClient := google.NewGoogleClient(config, bigQuery)

	return &Youtube{googleClient}
}
