package youtube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_google "github.com/leapforce-libraries/go_google"
	go_http "github.com/leapforce-libraries/go_http"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
)

const (
	apiName         string = "Youtube"
	apiURLAnalytics string = "https://youtubeanalytics.googleapis.com/v2"
	apiURLData      string = "https://youtube.googleapis.com/youtube/v3"
)

// Service stores Service configuration
//
type Service struct {
	id            string
	apiKey        *string
	accessToken   *string
	httpService   *go_http.Service
	googleService *go_google.Service
	quotaCosts    int64
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
		id:          serviceConfig.APIKey,
		apiKey:      &serviceConfig.APIKey,
		httpService: httpService,
	}, nil
}

type ServiceConfigWithAccessToken struct {
	ClientID    string
	AccessToken string
}

func NewServiceWithAccessToken(serviceConfig *ServiceConfigWithAccessToken) (*Service, *errortools.Error) {
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
		accessToken: &serviceConfig.AccessToken,
		id:          IDFromClientID(serviceConfig.ClientID),
		httpService: httpService,
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
		id:            IDFromClientID(serviceConfig.ClientID),
		googleService: googleService,
	}, nil
}

func (service *Service) httpRequest(httpMethod string, requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {

	var request *http.Request
	var response *http.Response
	var e *errortools.Error

	if service.httpService != nil {
		// add error model
		errorResponse := go_google.ErrorResponse{}
		requestConfig.ErrorModel = &errorResponse

		if service.apiKey != nil {
			// add api key
			requestConfig.SetParameter("key", *service.apiKey)
		}
		if service.accessToken != nil {
			// add accesstoken to header
			header := http.Header{}
			header.Set("Authorization", fmt.Sprintf("Bearer %s", *service.accessToken))
			requestConfig.NonDefaultHeaders = &header
		}

		request, response, e = service.httpService.HTTPRequest(httpMethod, requestConfig)

		if e != nil {
			if errorResponse.Error.Message != "" {
				e.SetMessage(errorResponse.Error.Message)
			}
		}
	} else if service.googleService != nil {
		request, response, e = service.googleService.HTTPRequest(httpMethod, requestConfig)
	}

	if e != nil {
		return request, response, e
	}

	return request, response, nil
}

func (service *Service) httpRequestWrapped(httpMethod string, requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *Response, *errortools.Error) {

	responseModel := requestConfig.ResponseModel

	// copy request config, add generic Response
	_response := Response{}
	_requestConfig := *requestConfig
	_requestConfig.ResponseModel = &_response

	request, response, e := service.httpRequest(httpMethod, &_requestConfig)
	if e != nil {
		return request, response, nil, e
	}

	// unmarshal items
	err := json.Unmarshal(_response.Items, responseModel)
	if err != nil {
		return request, response, nil, errortools.ErrorMessage(err)
	}

	return request, response, &_response, nil
}

func (service *Service) get(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	return service.httpRequest(http.MethodGet, requestConfig)
}

func (service *Service) getWrapped(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *Response, *errortools.Error) {
	return service.httpRequestWrapped(http.MethodGet, requestConfig)
}

func (service *Service) urlData(path string) string {
	return fmt.Sprintf("%s/%s", apiURLData, path)
}

func (service *Service) urlAnalytics(path string) string {
	return fmt.Sprintf("%s/%s", apiURLAnalytics, path)
}

func (service *Service) pay(quotaCosts int64) {
	service.quotaCosts += quotaCosts
}

func IDFromClientID(clientID string) string {
	return strings.Split(clientID, ".")[0]
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
