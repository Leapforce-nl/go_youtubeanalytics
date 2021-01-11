package Service

import (
	google "github.com/leapforce-libraries/go_google"
)

const (
	APIName string = "Youtube"
	APIURL  string = "https://youtubeanalytics.googleapis.com/v2"
)

// Service stores Service configuration
//
type Service struct {
	googleService *google.Service
}

// methods
//
func NewService(clientID string, clientSecret string, scope string, bigQuery *google.BigQuery) *Service {
	config := google.ServiceConfig{
		APIName:      APIName,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        scope,
	}

	googleService := google.NewService(config, bigQuery)

	return &Service{googleService}
}
