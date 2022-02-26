package youtube

import (
	"fmt"
	"net/http"
	"net/url"

	"cloud.google.com/go/civil"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	go_utilities "github.com/leapforce-libraries/go_utilities"
)

type targetedQueryResult struct {
	Kind          string          `json:"kind"`
	ColumnHeaders []ColumnHeader  `json:"columnHeaders"`
	Rows          [][]interface{} `json:"rows"`
}

type TargetedQueryResult []map[string]interface{}

type ColumnHeader struct {
	Name       string `json:"name"`
	ColumnType string `json:"columnType"`
	DataType   string `json:"dataType"`
}

type DoTargetedQueryConfig struct {
	EndDate                      *civil.Date
	Ids                          *string
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
	/*if service.authorizationMode == go_google.AuthorizationModeAPIKey {
		return nil, errortools.ErrorMessage("OAuth2 authorization required for this endpoint")
	}*/

	values := url.Values{}

	if doTargetedQueryConfig.EndDate != nil {
		values.Set("endDate", doTargetedQueryConfig.EndDate.String())
	}

	if doTargetedQueryConfig.Ids != nil {
		values.Set("ids", *doTargetedQueryConfig.Ids)
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

	_targetedQueryResult := targetedQueryResult{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlAnalytics("reports"),
		Parameters:    &values,
		ResponseModel: &_targetedQueryResult,
	}

	_, _, e := service.googleService().HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	t := TargetedQueryResult{}
	e = t.parse(&_targetedQueryResult)
	if e != nil {
		return nil, e
	}

	return &t, nil
}

func (t *TargetedQueryResult) parse(res *targetedQueryResult) *errortools.Error {
	if res == nil {
		return nil
	}

	for _, row := range res.Rows {
		res1 := make(map[string]interface{})

		for i, columnHeader := range res.ColumnHeaders {
			valueString := fmt.Sprintf("%v", row[i])
			var value interface{}

			switch columnHeader.DataType {
			case "STRING":
				value = valueString
			case "INTEGER":
				valueFloat64, err := go_utilities.ParseFloat(valueString)
				if err != nil {
					errortools.CaptureError(err)
				}
				value = int64(valueFloat64)
			case "FLOAT":
				valueFloat64, err := go_utilities.ParseFloat(valueString)
				if err != nil {
					errortools.CaptureError(err)
				}
				value = valueFloat64
			}

			res1[columnHeader.Name] = value
		}

		*t = append(*t, res1)
	}

	return nil
}
