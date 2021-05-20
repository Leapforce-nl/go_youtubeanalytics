package youtube

import (
	"fmt"
	"net/url"

	"cloud.google.com/go/civil"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type TargetedQueryResult struct {
	Kind          string          `json:"kind"`
	ColumnHeaders []ColumnHeader  `json:"columnHeaders"`
	Rows          [][]interface{} `json:"rows"`
}

type ColumnHeader struct {
	Name       string `json:"name"`
	ColumnType string `json:"columnType"`
	DataType   string `json:"dataType"`
}

type DoTargetedQueryConfig struct {
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

func (service *Service) DoTargetedQuery(doTargetedQueryConfig *DoTargetedQueryConfig) (*TargetedQueryResult, *errortools.Error) {
	if service.authorizationMode == AuthorizationModeAPIKey {
		return nil, errortools.ErrorMessage("OAuth2 authorization required for this endpoint")
	}

	values := url.Values{}

	if doTargetedQueryConfig.EndDate != nil {
		values.Set("endDate", doTargetedQueryConfig.EndDate.String())
	}

	if doTargetedQueryConfig.IDs != nil {
		values.Set("ids", *doTargetedQueryConfig.IDs)
	}

	if doTargetedQueryConfig.Metrics != nil {
		values.Set("metrics", *doTargetedQueryConfig.Metrics)
	}

	if doTargetedQueryConfig.StartDate != nil {
		values.Set("startDate", doTargetedQueryConfig.StartDate.String())
	}

	if doTargetedQueryConfig.Currency != nil {
		values.Set("currency", *doTargetedQueryConfig.Currency)
	}

	if doTargetedQueryConfig.Dimensions != nil {
		values.Set("dimensions", *doTargetedQueryConfig.Dimensions)
	}

	if doTargetedQueryConfig.Filters != nil {
		values.Set("filters", *doTargetedQueryConfig.Filters)
	}

	if doTargetedQueryConfig.IncludeHistoricalChannelData != nil {
		values.Set("includeHistoricalChannelData", fmt.Sprintf("%v", *doTargetedQueryConfig.IncludeHistoricalChannelData))
	}

	if doTargetedQueryConfig.MaxResults != nil {
		values.Set("maxResults", fmt.Sprintf("%v", *doTargetedQueryConfig.MaxResults))
	}

	if doTargetedQueryConfig.Sort != nil {
		values.Set("sort", *doTargetedQueryConfig.Sort)
	}

	if doTargetedQueryConfig.StartIndex != nil {
		values.Set("startIndex", fmt.Sprintf("%v", *doTargetedQueryConfig.StartIndex))
	}

	targetedQueryResult := TargetedQueryResult{}

	requestConfig := go_http.RequestConfig{
		URL:           service.urlAnalytics("reports"),
		Parameters:    &values,
		ResponseModel: &targetedQueryResult,
	}
	service.pay(1)
	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &targetedQueryResult, nil
}
