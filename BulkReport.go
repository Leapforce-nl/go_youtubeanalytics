package youtube

import (
	"encoding/csv"
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	fileio "github.com/leapforce-libraries/go_fileio"
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
	Date                           y_types.BulkReportDateString `csv:"date"`
	ChannelID                      string                       `csv:"channel_id"`
	VideoID                        string                       `csv:"video_id"`
	LiveOrOnDemand                 string                       `csv:"live_or_on_demand"`
	SubscribedStatus               string                       `csv:"subscribed_status"`
	CountryCode                    string                       `csv:"country_code"`
	Views                          int64                        `csv:"views"`
	Comments                       int64                        `csv:"comments"`
	Likes                          int64                        `csv:"likes"`
	Dislikes                       int64                        `csv:"dislikes"`
	Shares                         int64                        `csv:"shares"`
	WatchTimeMinutes               float64                      `csv:"watch_time_minutes"`
	AverageViewDurationSeconds     float64                      `csv:"average_view_duration_seconds"`
	AverageViewDurationPercentage  float64                      `csv:"average_view_duration_percentage"`
	AnnotationImpressions          int64                        `csv:"annotation_impressions"`
	AnnotationClickableImpressions int64                        `csv:"annotation_clickable_impressions"`
	AnnotationClicks               int64                        `csv:"annotation_clicks"`
	AnnotationClickThroughRate     float64                      `csv:"annotation_click_through_rate"`
	AnnotationClosableImpressions  int64                        `csv:"annotation_closable_impressions"`
	AnnotationCloses               int64                        `csv:"annotation_closes"`
	AnnotationCloseRate            float64                      `csv:"annotation_close_rate"`
	CardTeaserImpressions          int64                        `csv:"card_teaser_impressions"`
	CardTeaserClicks               int64                        `csv:"card_teaser_clicks"`
	CardTeaserClickRate            float64                      `csv:"card_teaser_click_rate"`
	CardImpressions                int64                        `csv:"card_impressions"`
	CardClicks                     int64                        `csv:"card_clicks"`
	CardClickRate                  float64                      `csv:"card_click_rate"`
	SubscribersGained              int64                        `csv:"subscribers_gained"`
	SubscribersLost                int64                        `csv:"subscribers_lost"`
	VideosAddedToPlaylists         int64                        `csv:"videos_added_to_playlists"`
	VideosRemovedFromPlaylists     int64                        `csv:"videos_removed_from_playlists"`
	RedViews                       int64                        `csv:"red_views"`
	RedWatchTimeMinutes            float64                      `csv:"red_watch_time_minutes"`
}

func (service *Service) DownloadReport(url string) (*[]BulkReportChannelBasicA2, *errortools.Error) {
	requestConfig := go_http.RequestConfig{
		URL: url,
	}
	_, res, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	data := []BulkReportChannelBasicA2{}

	defer res.Body.Close()

	r := csv.NewReader(res.Body)
	r.Comma = []rune(",")[0]

	e = fileio.GetCSVFromCSVReader(r, &data)
	if e != nil {
		return nil, e
	}

	return &data, nil
}
