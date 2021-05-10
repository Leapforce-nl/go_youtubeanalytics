package youtube

import "encoding/json"

type Response struct {
	Kind          string          `json:"kind"`
	Etag          string          `json:"etag"`
	NextPageToken *string         `json:"nextPageToken"`
	PrevPageToken *string         `json:"prevPageToken"`
	RegionCode    *string         `json:"regionCode"`
	PageInfo      PageInfo        `json:"pageInfo"`
	Items         json.RawMessage `json:"items"`
}
