package youtube

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"reflect"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	y_types "github.com/leapforce-libraries/go_youtube/types"
)

// Job //
//
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

// Report //
//
type Report struct {
	ID            string                  `json:"id"`
	JobID         string                  `json:"jobId"`
	StartTime     y_types.DateTimeString  `json:"startTime"`
	EndTime       y_types.DateTimeString  `json:"endTime"`
	CreateTime    y_types.DateTimeString  `json:"createTime"`
	JobExpireTime *y_types.DateTimeString `json:"jobExpireTime"`
	DownloadURL   string                  `json:"downloadUrl"`
}

type GetReportsConfig struct {
	JobID        string
	CreatedAfter *y_types.DateTimeString
}

type GetReportsResponse struct {
	Reports       *[]Report `json:"reports"`
	NextPageToken *string   `json:"nextPageToken"`
}

func (service *Service) GetReports(getReportsConfig *GetReportsConfig) (*[]Report, *errortools.Error) {
	if service.authorizationMode == AuthorizationModeAPIKey {
		return nil, errortools.ErrorMessage("OAuth2 authorization required for this endpoint")
	}

	if getReportsConfig == nil {
		return nil, errortools.ErrorMessage("GetReportsConfig is nil")
	}

	reports := []Report{}
	values := url.Values{}

	if getReportsConfig.CreatedAfter != nil {
		values.Set("createdAfter", getReportsConfig.CreatedAfter.String())
	}

	for true {
		getReportsResponse := GetReportsResponse{}

		requestConfig := go_http.RequestConfig{
			URL:           service.apiURLReporting(fmt.Sprintf("jobs/%s/reports", getReportsConfig.JobID)),
			Parameters:    &values,
			ResponseModel: &getReportsResponse,
		}
		service.pay(1)
		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		if getReportsResponse.Reports == nil {
			break
		}

		reports = append(reports, *getReportsResponse.Reports...)

		if getReportsResponse.NextPageToken == nil {
			break
		}

		values.Set("pageToken", *getReportsResponse.NextPageToken)
	}

	return &reports, nil
}

type BulkReportChannelBasicA2 struct {
	Date                           string `csv:"date"`
	ChannelID                      string `csv:"channelId"`
	VideoID                        string `csv:"videoId"`
	LiveOrOnDemand                 string `csv:"live_or_on_demand"`
	SubscribedStatus               string `csv:"subscribed_status"`
	CountryCode                    string `csv:"countryCode"`
	Views                          string `csv:"views"`
	Comments                       string `csv:"comments"`
	Likes                          string `csv:"likes"`
	Dislikes                       string `csv:"dislikes"`
	Shares                         string `csv:"shares"`
	WatchTimeMinutes               string `csv:"watchTime_minutes"`
	AverageViewDurationSeconds     string `csv:"average_view_duration_seconds"`
	AverageViewDurationPercentage  string `csv:"average_view_duration_percentage"`
	AnnotationImpressions          string `csv:"annotationImpressions"`
	AnnotationClickableImpressions string `csv:"annotationClickableImpressions"`
	AnnotationClicks               string `csv:"annotationClicks"`
	AnnotationClickThroughRate     string `csv:"annotationClickThroughRate"`
	AnnotationClosableImpressions  string `csv:"annotationClosableImpressions"`
	AnnotationCloses               string `csv:"annotationCloses"`
	AnnotationCloseRate            string `csv:"annotationCloseRate"`
	CardTeaserImpressions          string `csv:"cardTeaserImpressions"`
	CardTeaserClickRate            string `csv:"cardTeaserClickRate"`
	CardTeaserClicks               string `csv:"cardTeaserClicks"`
	CardImpressions                string `csv:"cardImpressions"`
	CardClicks                     string `csv:"cardClicks"`
	CardClickRate                  string `csv:"cardClickRate"`
	SubscribersGained              string `csv:"subscribers_gained"`
	SubscribersLost                string `csv:"subscribers_lost"`
	VideosAddedToPlaylists         string `csv:"videos_addedTo_playlists"`
	VideosRemovedFromPlaylists     string `csv:"videosRemoved_from_playlists"`
	RedViews                       string `csv:"red_views"`
	RedWatchTimeMinutes            string `csv:"red_watchTime_minutes"`
}

func (service *Service) DownloadReport(url string) (*[]BulkReportChannelBasicA2, *errortools.Error) {
	requestConfig := go_http.RequestConfig{
		URL: url,
	}
	_, res, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	skipRows := 1

	data := []BulkReportChannelBasicA2{}
	fieldCountStruct := reflect.TypeOf(BulkReportChannelBasicA2{}).NumField()

	defer res.Body.Close()

	r := csv.NewReader(res.Body)
	r.Comma = []rune(";")[0]

	row := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, errortools.ErrorMessagef("Error while reading response of %s, row: %v, message: %s", url, row, err.Error())
		}

		if row >= skipRows {
			fieldCountRecord := len(record)

			if fieldCountRecord == fieldCountStruct-1 {

				dr := BulkReportChannelBasicA2{
					//FileName: file.Name,
				}

				fieldIndexStruct := 1 //first field contains filename
				for fieldIndexStruct < fieldCountStruct {
					fieldIndexRecord := fieldIndexStruct - 1
					val := reflect.ValueOf(&dr)
					s := val.Elem()
					f := s.Field(fieldIndexStruct)

					if f.IsValid() {
						if f.CanSet() {
							if f.Kind() == reflect.String {
								f.SetString(record[fieldIndexRecord])
							}
						}
					}

					fieldIndexStruct++
				}

				data = append(data, dr)
			} else {
				return nil, errortools.ErrorMessagef("Row %v in response of %s contains %v fields instead of the expect %v fields.", row, url, fieldCountRecord, fieldCountStruct-1)
			}
		}

		row++
	}

	return &data, nil
}
