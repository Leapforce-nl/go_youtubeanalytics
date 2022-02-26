package youtube

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	y_types "github.com/leapforce-libraries/go_youtube/types"
)

type SearchResult struct {
	Kind       string                 `json:"kind"`
	Etag       string                 `json:"etag"`
	Id         SearchResultId         `json:"id"`
	Snippet    SearchResultSnippet    `json:"snippet"`
	Statistics SearchResultStatistics `json:"statistics"`
}

type SearchResultId struct {
	Kind    string `json:"kind"`
	VideoId string `json:"videoId"`
}

type SearchResultSnippet struct {
	PublishedAt          *y_types.DateTimeString `json:"publishedAt"`
	ChannelId            string                  `json:"channelId"`
	Title                string                  `json:"title"`
	Description          string                  `json:"description"`
	Thumbnails           Thumbnails              `json:"thumbnails"`
	ChannelTitle         string                  `json:"channelTitle"`
	LiveBroadcastContent string                  `json:"liveBroadcastContent"`
	PublishTime          *y_types.DateTimeString `json:"publishTime"`
}

type SearchResultStatistics struct {
	ViewCount     uint64 `json:"viewCount,string"`
	LikeCount     uint64 `json:"likeCount,string"`
	DislikeCount  uint64 `json:"dislikeCount,string"`
	FavoriteCount uint64 `json:"favoriteCount,string"`
	CommentCount  uint64 `json:"commentCount,string"`
}

type SearchPart string

const (
	SearchPartSnippet SearchPart = "snippet"
)

type ChannelType string

const (
	ChannelTypeAny  ChannelType = "any"
	ChannelTypeShow ChannelType = "show"
)

type EventType string

const (
	EventTypeCompleted EventType = "completed"
	EventTypeLive      EventType = "live"
	EventTypeUpcoming  EventType = "upcoming"
)

type Order string

const (
	OrderDate       Order = "date"
	OrderRating     Order = "rating"
	OrderRelevance  Order = "relevance"
	OrderTitle      Order = "title"
	OrderVideoCount Order = "videoCount"
	OrderViewCount  Order = "viewCount"
)

type SafeSearch string

const (
	SafeSearchModerate SafeSearch = "moderate"
	SafeSearchNone     SafeSearch = "none"
	SafeSearchStrict   SafeSearch = "strict"
)

type SearchType string

const (
	SearchTypeChannel  SearchType = "channel"
	SearchTypePlaylist SearchType = "playlist"
	SearchTypeVideo    SearchType = "video"
)

type VideoDefinition string

const (
	VideoDefinitionAny      VideoDefinition = "any"
	VideoDefinitionHigh     VideoDefinition = "high"
	VideoDefinitionStandard VideoDefinition = "standard"
)

type VideoCaption string

const (
	VideoCaptionAny           VideoCaption = "any"
	VideoCaptionClosedCaption VideoCaption = "closedCaption"
	VideoCaptionNone          VideoCaption = "none"
)

type VideoDimension string

const (
	VideoDimension2D  VideoDimension = "2d"
	VideoDimension3D  VideoDimension = "3d"
	VideoDimensionAny VideoDimension = "any"
)

type VideoDuration string

const (
	VideoDurationAny    VideoDuration = "any"
	VideoDurationLong   VideoDuration = "long"
	VideoDurationMedium VideoDuration = "medium"
	VideoDurationShort  VideoDuration = "short"
)

type VideoEmbeddable string

const (
	VideoEmbeddableAny  VideoEmbeddable = "any"
	VideoEmbeddableTrue VideoEmbeddable = "true"
)

type VideoLicense string

const (
	VideoLicenseAny            VideoLicense = "any"
	VideoLicenseCreativeCommon VideoLicense = "creativeCommon"
	VideoLicenseYoutube        VideoLicense = "youtube"
)

type VideoSyndicated string

const (
	VideoSyndicatedAny  VideoSyndicated = "any"
	VideoSyndicatedTrue VideoSyndicated = "true"
)

type VideoType string

const (
	VideoTypeAny     VideoType = "any"
	VideoTypeEpisode VideoType = "episode"
	VideoTypeMovie   VideoType = "movie"
)

type SearchConfig struct {
	Part                   []SearchPart
	ForContentOwner        *bool
	ForDeveloper           *bool
	ForMine                *bool
	RelatedToVideoId       *string
	ChannelId              *string
	ChannelType            *ChannelType
	EventType              *EventType
	Location               *string
	LocationRadius         *string
	MaxResults             *uint64
	OnBehalfOfContentOwner *string
	Order                  *Order
	PageToken              *string
	PublishedAfter         *y_types.DateTimeString
	PublishedBefore        *y_types.DateTimeString
	Query                  *string
	RegionCode             *string
	RelevanceLanguage      *string
	SafeSearch             *SafeSearch
	TopicId                *string
	Type                   *SearchType
	VideoCaption           *VideoCaption
	VideoCategoryId        *string
	VideoDefinition        *VideoDefinition
	VideoDimension         *VideoDimension
	VideoDuration          *VideoDuration
	VideoEmbeddable        *VideoEmbeddable
	VideoLicense           *VideoLicense
	VideoSyndicated        *VideoSyndicated
	VideoType              *VideoType
}

func (service *Service) Search(searchConfig *SearchConfig) (*[]SearchResult, *string, *errortools.Error) {
	values := url.Values{}

	searchResultParts := []string{}
	for _, searchResultPart := range searchConfig.Part {
		searchResultParts = append(searchResultParts, string(searchResultPart))
	}
	values.Set("part", strings.Join(searchResultParts, ","))

	if searchConfig.ForContentOwner != nil {
		values.Set("forContentOwner", fmt.Sprintf("%v", *searchConfig.ForContentOwner))
	}

	if searchConfig.ForDeveloper != nil {
		values.Set("forDeveloper", fmt.Sprintf("%v", *searchConfig.ForDeveloper))
	}

	if searchConfig.ForMine != nil {
		values.Set("forMine", fmt.Sprintf("%v", *searchConfig.ForMine))
	}

	if searchConfig.RelatedToVideoId != nil {
		values.Set("relatedToVideoId", *searchConfig.RelatedToVideoId)
	}

	if searchConfig.ChannelId != nil {
		values.Set("channelId", *searchConfig.ChannelId)
	}

	if searchConfig.ChannelType != nil {
		values.Set("channelType", fmt.Sprintf("%v", *searchConfig.ChannelType))
	}

	if searchConfig.EventType != nil {
		values.Set("eventType", fmt.Sprintf("%v", *searchConfig.EventType))
	}

	if searchConfig.LocationRadius != nil {
		values.Set("locationRadius", *searchConfig.LocationRadius)
	}

	if searchConfig.MaxResults != nil {
		values.Set("maxResults", fmt.Sprintf("%v", *searchConfig.MaxResults))
	}

	if searchConfig.OnBehalfOfContentOwner != nil {
		values.Set("onBehalfOfContentOwner", *searchConfig.OnBehalfOfContentOwner)
	}

	if searchConfig.Order != nil {
		values.Set("order", fmt.Sprintf("%v", *searchConfig.Order))
	}

	if searchConfig.PageToken != nil {
		values.Set("pageToken", *searchConfig.PageToken)
	}

	if searchConfig.PublishedAfter != nil {
		values.Set("publishedAfter", searchConfig.PublishedAfter.String())
	}

	if searchConfig.PublishedBefore != nil {
		values.Set("publishedBefore", searchConfig.PublishedBefore.String())
	}

	if searchConfig.Query != nil {
		values.Set("q", *searchConfig.Query)
	}

	if searchConfig.RegionCode != nil {
		values.Set("regionCode", *searchConfig.RegionCode)
	}

	if searchConfig.SafeSearch != nil {
		values.Set("safeSearch", fmt.Sprintf("%v", *searchConfig.SafeSearch))
	}

	if searchConfig.TopicId != nil {
		values.Set("topicId", *searchConfig.TopicId)
	}

	if searchConfig.Type != nil {
		values.Set("type", fmt.Sprintf("%v", *searchConfig.Type))
	}

	if searchConfig.VideoCaption != nil {
		values.Set("videoCaption", fmt.Sprintf("%v", *searchConfig.VideoCaption))
	}

	if searchConfig.VideoCategoryId != nil {
		values.Set("videoCategoryId", *searchConfig.VideoCategoryId)
	}

	if searchConfig.VideoDefinition != nil {
		values.Set("videoDefinition", fmt.Sprintf("%v", *searchConfig.VideoDefinition))
	}

	if searchConfig.VideoDimension != nil {
		values.Set("videoDimension", fmt.Sprintf("%v", *searchConfig.VideoDimension))
	}

	if searchConfig.VideoDuration != nil {
		values.Set("videoDuration", fmt.Sprintf("%v", *searchConfig.VideoDuration))
	}

	if searchConfig.VideoEmbeddable != nil {
		values.Set("videoEmbeddable", fmt.Sprintf("%v", *searchConfig.VideoEmbeddable))
	}

	if searchConfig.VideoLicense != nil {
		values.Set("videoLicense", fmt.Sprintf("%v", *searchConfig.VideoLicense))
	}

	if searchConfig.VideoSyndicated != nil {
		values.Set("videoSyndicated", fmt.Sprintf("%v", *searchConfig.VideoSyndicated))
	}

	if searchConfig.VideoType != nil {
		values.Set("videoType", fmt.Sprintf("%v", *searchConfig.VideoType))
	}

	searchResults := []SearchResult{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlData("search"),
		Parameters:    &values,
		ResponseModel: &searchResults,
	}

	_, _, response, e := service.httpRequestWrapped(&requestConfig)
	if e != nil {
		return nil, nil, e
	}

	return &searchResults, response.NextPageToken, nil
}
