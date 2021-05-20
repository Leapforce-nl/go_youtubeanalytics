package youtube

import (
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	y_types "github.com/leapforce-libraries/go_youtube/types"
)

type Job struct {
	ID           string                 `json:"id"`
	ReportTypeID string                 `json:"reportTypeId"`
	Name         string                 `json:"name"`
	CreateTime   y_types.DateTimeString `json:"createTime"`
}

type CreateJobConfig struct {
	ReportTypeID string `json:"reportTypeId"`
	Name         string `json:"name"`
}

func (service *Service) CreateJob(createJobConfig *CreateJobConfig) (*Job, *errortools.Error) {
	if service.authorizationMode == AuthorizationModeAPIKey {
		return nil, errortools.ErrorMessage("OAuth2 authorization required for this endpoint")
	}

	if createJobConfig == nil {
		return nil, errortools.ErrorMessage("CreateJobConfig is nil")
	}

	job := Job{}

	requestConfig := go_http.RequestConfig{
		URL:           service.apiURLReporting("jobs"),
		BodyModel:     createJobConfig,
		ResponseModel: &job,
	}
	service.pay(1)
	_, _, e := service.post(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &job, nil
}

type GetJobsResponse struct {
	Jobs          *[]Job  `json:"jobs"`
	NextPageToken *string `json:"nextPageToken"`
}

func (service *Service) GetJobs() (*[]Job, *errortools.Error) {
	if service.authorizationMode == AuthorizationModeAPIKey {
		return nil, errortools.ErrorMessage("OAuth2 authorization required for this endpoint")
	}

	jobs := []Job{}
	values := url.Values{}

	for true {
		getJobsResponse := GetJobsResponse{}

		requestConfig := go_http.RequestConfig{
			URL:           service.apiURLReporting("jobs"),
			Parameters:    &values,
			ResponseModel: &getJobsResponse,
		}
		service.pay(1)
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		if getJobsResponse.Jobs == nil {
			break
		}

		jobs = append(jobs, *getJobsResponse.Jobs...)

		if getJobsResponse.NextPageToken == nil {
			break
		}

		values.Set("pageToken", *getJobsResponse.NextPageToken)
	}

	return &jobs, nil
}
