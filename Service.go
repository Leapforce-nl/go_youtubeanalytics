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
	apiUrlAnalytics string = "https://youtubeanalytics.googleapis.com/v2"
	apiUrlReporting string = "https://youtubereporting.googleapis.com/v1"
	apiUrlData      string = "https://youtube.googleapis.com/youtube/v3"
)

type Service go_google.Service

func NewServiceWithAccessToken(cfg *go_google.ServiceWithAccessTokenConfig) (*Service, *errortools.Error) {
	googleService, e := go_google.NewServiceWithAccessToken(cfg)
	if e != nil {
		return nil, e
	}
	service := Service(*googleService)
	return &service, nil
}

func NewServiceWithApiKey(cfg *go_google.ServiceWithApiKeyConfig) (*Service, *errortools.Error) {
	googleService, e := go_google.NewServiceWithApiKey(cfg)
	if e != nil {
		return nil, e
	}
	service := Service(*googleService)
	return &service, nil
}

func NewServiceWithOAuth2(cfg *go_google.ServiceWithOAuth2Config) (*Service, *errortools.Error) {
	googleService, e := go_google.NewServiceWithOAuth2(cfg)
	if e != nil {
		return nil, e
	}
	service := Service(*googleService)
	return &service, nil
}

func (service *Service) httpRequestWrapped(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *Response, *errortools.Error) {

	responseModel := requestConfig.ResponseModel

	// copy request config, add generic Response
	_response := Response{}
	_requestConfig := *requestConfig
	_requestConfig.ResponseModel = &_response

	request, response, e := service.googleService().HttpRequest(&_requestConfig)
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
	return fmt.Sprintf("%s/%s", apiUrlData, path)
}

func (service *Service) urlAnalytics(path string) string {
	return fmt.Sprintf("%s/%s", apiUrlAnalytics, path)
}

func (service *Service) apiUrlReporting(path string) string {
	return fmt.Sprintf("%s/%s", apiUrlReporting, path)
}

func (service *Service) ApiName() string {
	return apiName
}

func (service *Service) ApiKey() string {
	return service.googleService().ApiKey()
}

func (service *Service) ApiCallCount() int64 {
	return service.googleService().ApiCallCount()
}

func (service *Service) ApiReset() {
	service.googleService().ApiReset()
}

func (service *Service) googleService() *go_google.Service {
	googleService := go_google.Service(*service)
	return &googleService
}
