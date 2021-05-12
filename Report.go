package youtube

import (
	"fmt"
	"net/url"

	"cloud.google.com/go/civil"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type Report struct {
	Kind          string          `json:"kind"`
	ColumnHeaders []ColumnHeader  `json:"columnHeaders"`
	Rows          [][]interface{} `json:"rows"`
}

type ColumnHeader struct {
	Name       string `json:"name"`
	ColumnType string `json:"columnType"`
	DataType   string `json:"dataType"`
}

type GetReportConfig struct {
	EndDate                      *civil.Date
	IDs                          *string
	Metrics                      *string
	StartDate                    *civil.Date
	Currency                     *string
	Dimensions                   *string
	Filters                      *string
	IncludeHistoricalChannelData *bool
	MaxResults                   *uint64
	Sort                         *string
	StartIndex                   *uint64
}

func (service *Service) GetReport(getReportConfig *GetReportConfig) (*Report, *errortools.Error) {
	values := url.Values{}

	if getReportConfig.EndDate != nil {
		values.Set("endDate", getReportConfig.EndDate.String())
	}

	if getReportConfig.IDs != nil {
		values.Set("ids", *getReportConfig.IDs)
	}

	if getReportConfig.Metrics != nil {
		values.Set("metrics", *getReportConfig.Metrics)
	}

	if getReportConfig.StartDate != nil {
		values.Set("startDate", getReportConfig.StartDate.String())
	}

	if getReportConfig.Currency != nil {
		values.Set("currency", *getReportConfig.Currency)
	}

	if getReportConfig.Dimensions != nil {
		values.Set("dimensions", *getReportConfig.Dimensions)
	}

	if getReportConfig.Filters != nil {
		values.Set("filters", *getReportConfig.Filters)
	}

	if getReportConfig.IncludeHistoricalChannelData != nil {
		values.Set("includeHistoricalChannelData", fmt.Sprintf("%v", *getReportConfig.IncludeHistoricalChannelData))
	}

	if getReportConfig.MaxResults != nil {
		values.Set("maxResults", fmt.Sprintf("%v", *getReportConfig.MaxResults))
	}

	if getReportConfig.Sort != nil {
		values.Set("sort", *getReportConfig.Sort)
	}

	if getReportConfig.StartIndex != nil {
		values.Set("startIndex", fmt.Sprintf("%v", *getReportConfig.StartIndex))
	}

	report := Report{}

	requestConfig := go_http.RequestConfig{
		URL:           service.urlAnalytics("reports"),
		Parameters:    &values,
		ResponseModel: &report,
	}
	service.pay(1)
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &report, nil
}
