package youtube

import (
	"encoding/json"
	"fmt"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_google "github.com/leapforce-libraries/go_google"
	go_http "github.com/leapforce-libraries/go_http"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
)

const (
	apiName         string = "Youtube"
	apiURLAnalytics string = "https://youtubeanalytics.googleapis.com/v2"
	apiURLReporting string = "https://youtubereporting.googleapis.com/v1"
	apiURLData      string = "https://youtube.googleapis.com/youtube/v3"
)

// Service stores Service configuration
//
type Service struct {
	authorizationMode go_google.AuthorizationMode
	id                string
	apiKey            *string
	accessToken       *string
	httpService       *go_http.Service
	googleService     *go_google.Service
	quotaCosts        int64
}

type ServiceConfigWithAPIKey struct {
	APIKey string
}

func NewServiceWithAPIKey(serviceConfig *ServiceConfigWithAPIKey) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if serviceConfig.APIKey == "" {
		return nil, errortools.ErrorMessage("APIKey not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &Service{
		authorizationMode: go_google.AuthorizationModeAPIKey,
		id:                serviceConfig.APIKey,
		apiKey:            &serviceConfig.APIKey,
		httpService:       httpService,
	}, nil
}

type ServiceWithAccessTokenConfig struct {
	ClientID    string
	AccessToken string
}

func NewServiceWithAccessToken(serviceConfig *ServiceWithAccessTokenConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if serviceConfig.AccessToken == "" {
		return nil, errortools.ErrorMessage("AccessToken not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &Service{
		authorizationMode: go_google.AuthorizationModeAccessToken,
		accessToken:       &serviceConfig.AccessToken,
		id:                go_google.ClientIDShort(serviceConfig.ClientID),
		httpService:       httpService,
	}, nil
}

type ServiceConfigOAuth2 struct {
	ClientID          string
	ClientSecret      string
	GetTokenFunction  *func() (*oauth2.Token, *errortools.Error)
	SaveTokenFunction *func(token *oauth2.Token) *errortools.Error
}

func NewServiceOAuth2(serviceConfig *ServiceConfigOAuth2) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if serviceConfig.ClientID == "" {
		return nil, errortools.ErrorMessage("ClientID not provided")
	}

	if serviceConfig.ClientSecret == "" {
		return nil, errortools.ErrorMessage("ClientSecret not provided")
	}
	/*
		getTokenFunction := func() (*oauth2.Token, *errortools.Error) {
			return GetToken(serviceConfig.ClientID, serviceConfig.ChannelID)
		}

		saveTokenFunction := func(token *oauth2.Token) *errortools.Error {
			return SaveToken(serviceConfig.ClientID, serviceConfig.ChannelID, token)
		}*/

	googleServiceConfig := go_google.ServiceConfig{
		APIName:           apiName,
		ClientID:          serviceConfig.ClientID,
		ClientSecret:      serviceConfig.ClientSecret,
		GetTokenFunction:  serviceConfig.GetTokenFunction,
		SaveTokenFunction: serviceConfig.SaveTokenFunction,
	}

	googleService, e := go_google.NewService(&googleServiceConfig, nil)
	if e != nil {
		return nil, e
	}

	return &Service{
		authorizationMode: go_google.AuthorizationModeOAuth2,
		id:                go_google.ClientIDShort(serviceConfig.ClientID),
		googleService:     googleService,
	}, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	var request *http.Request
	var response *http.Response
	var e *errortools.Error

	if service.authorizationMode == go_google.AuthorizationModeOAuth2 {
		request, response, e = service.googleService.HTTPRequest(requestConfig)
	} else {
		// add error model
		errorResponse := go_google.ErrorResponse{}
		requestConfig.ErrorModel = &errorResponse

		if service.authorizationMode == go_google.AuthorizationModeAPIKey {
			// add api key
			requestConfig.SetParameter("key", *service.apiKey)
		}
		if service.accessToken != nil {
			// add accesstoken to header
			header := http.Header{}
			header.Set("Authorization", fmt.Sprintf("Bearer %s", *service.accessToken))
			requestConfig.NonDefaultHeaders = &header
		}

		request, response, e = service.httpService.HTTPRequest(requestConfig)

		if e != nil {
			if errorResponse.Error.Message != "" {
				e.SetMessage(errorResponse.Error.Message)
			}
		}
	}

	if e != nil {
		return request, response, e
	}

	return request, response, nil
}

func (service *Service) httpRequestWrapped(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *Response, *errortools.Error) {

	responseModel := requestConfig.ResponseModel

	// copy request config, add generic Response
	_response := Response{}
	_requestConfig := *requestConfig
	_requestConfig.ResponseModel = &_response

	request, response, e := service.httpRequest(&_requestConfig)
	if e != nil {
		return request, response, nil, e
	}

	if _response.Items == nil {
		return request, response, nil, errortools.ErrorMessage("Response does not contain any items")
	}

	// unmarshal items
	err := json.Unmarshal(_response.Items, responseModel)
	if err != nil {
		return request, response, nil, errortools.ErrorMessage(err)
	}

	return request, response, &_response, nil
}

func (service *Service) urlData(path string) string {
	return fmt.Sprintf("%s/%s", apiURLData, path)
}

func (service *Service) urlAnalytics(path string) string {
	return fmt.Sprintf("%s/%s", apiURLAnalytics, path)
}

func (service *Service) apiURLReporting(path string) string {
	return fmt.Sprintf("%s/%s", apiURLReporting, path)
}

func (service *Service) pay(quotaCosts int64) {
	service.quotaCosts += quotaCosts
}

func (service *Service) InitToken(scope string, accessType *string, prompt *string, state *string) *errortools.Error {
	if service.googleService == nil {
		return nil
	}
	return service.googleService.InitToken(scope, accessType, prompt, state)
}

func (service *Service) APIName() string {
	return apiName
}

func (service *Service) APIKey() string {
	return service.id
}

func (service *Service) APICallCount() int64 {
	//return service.httpService.RequestCount()
	return service.quotaCosts
}

func (service *Service) APIReset() {
	service.httpService.ResetRequestCount()
}
