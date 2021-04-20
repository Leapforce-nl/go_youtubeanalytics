package Service

import (
	errortools "github.com/leapforce-libraries/go_errortools"
	google "github.com/leapforce-libraries/go_google"
	bigquery "github.com/leapforce-libraries/go_google/bigquery"
)

const (
	apiName string = "Youtube"
	apiURL  string = "https://youtubeanalytics.googleapis.com/v2"
)

// Service stores Service configuration
//
type Service struct {
	googleService *google.Service
}

type ServiceConfig struct {
	ClientID     string
	ClientSecret string
	Scope        string
}

// methods
//
func NewService(serviceConfig *ServiceConfig, bigQueryService *bigquery.Service) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if serviceConfig.ClientID == "" {
		return nil, errortools.ErrorMessage("ClientID not provided")
	}

	if serviceConfig.ClientSecret == "" {
		return nil, errortools.ErrorMessage("ClientSecret not provided")
	}

	config := google.ServiceConfig{
		APIName:      apiName,
		ClientID:     serviceConfig.ClientID,
		ClientSecret: serviceConfig.ClientSecret,
		Scope:        serviceConfig.Scope,
	}

	googleService := google.NewService(config, bigQueryService)

	return &Service{googleService}, nil
}

func (service *Service) InitToken() *errortools.Error {
	return service.googleService.InitToken()
}
