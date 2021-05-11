package youtube

import (
	"encoding/json"
	"fmt"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_google "github.com/leapforce-libraries/go_google"
	go_http "github.com/leapforce-libraries/go_http"
)

const (
	apiName         string = "Youtube"
	apiURLAnalytics string = "https://youtubeanalytics.googleapis.com/v2"
	apiURLData      string = "https://youtube.googleapis.com/youtube/v3"
)

// Service stores Service configuration
//
type Service struct {
	apiKey      string
	httpService *go_http.Service
	quotaCosts  int64
}

type ServiceConfig struct {
	APIKey string
}

// methods
//
func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
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
		apiKey:      serviceConfig.APIKey,
		httpService: httpService,
	}, nil
}

func (service *Service) httpRequest(httpMethod string, requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *Response, *errortools.Error) {
	// add api key
	requestConfig.SetParameter("key", service.apiKey)

	responseModel := requestConfig.ResponseModel

	// copy request config, add generic Response
	_response := Response{}
	_requestConfig := *requestConfig
	_requestConfig.ResponseModel = &_response

	// add error model
	errorResponse := go_google.ErrorResponse{}
	_requestConfig.ErrorModel = &errorResponse

	request, response, e := service.httpService.HTTPRequest(httpMethod, &_requestConfig)
	if errorResponse.Error.Message != "" {
		e.SetMessage(errorResponse.Error.Message)
	}

	// unmarshal items
	err := json.Unmarshal(_response.Items, responseModel)
	if err != nil {
		if e == nil {
			e = new(errortools.Error)
		}
		e.SetMessage(err)
	}

	return request, response, &_response, e
}

func (service *Service) get(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *Response, *errortools.Error) {
	return service.httpRequest(http.MethodGet, requestConfig)
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

func (service *Service) APIName() string {
	return apiName
}

func (service *Service) APIKey() string {
	return service.apiKey
}

func (service *Service) APICallCount() int64 {
	//return service.httpService.RequestCount()
	return service.quotaCosts
}

func (service *Service) APIReset() {
	service.httpService.ResetRequestCount()
}
