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

type Video struct {
	Kind       string          `json:"kind"`
	Etag       string          `json:"etag"`
	Id         string          `json:"id"`
	Snippet    VideoSnippet    `json:"snippet"`
	Status     VideoStatus     `json:"status"`
	Statistics VideoStatistics `json:"statistics"`
}

type VideoSnippet struct {
	PublishedAt          *y_types.DateTimeString `json:"publishedAt"`
	ChannelId            string                  `json:"channelId"`
	Title                string                  `json:"title"`
	Description          string                  `json:"description"`
	Thumbnails           Thumbnails              `json:"thumbnails"`
	ChannelTitle         string                  `json:"channelTitle"`
	Tags                 []string                `json:"tags"`
	CategoryId           uint64                  `json:"categoryId,string"`
	LiveBroadcastContent string                  `json:"liveBroadcastContent"`
	DefaultLanguage      string                  `json:"defaultLanguage"`
	Localized            Localized               `json:"localized"`
	DefaultAudioLanguage string                  `json:"defaultAudioLanguage"`
}

type VideoStatus struct {
	UploadStatus            string                  `json:"uploadStatus"`
	FailureReason           string                  `json:"failureReason"`
	RejectionReason         string                  `json:"rejectionReason"`
	PrivacyStatus           string                  `json:"privacyStatus"`
	PublishAt               *y_types.DateTimeString `json:"publishAt"`
	License                 string                  `json:"license"`
	Embeddable              bool                    `json:"embeddable"`
	PublicStatsViewable     bool                    `json:"publicStatsViewable"`
	MadeForKids             bool                    `json:"madeForKids"`
	SelfDeclaredMadeForKids bool                    `json:"selfDeclaredMadeForKids"`
}

type VideoStatistics struct {
	ViewCount     uint64 `json:"viewCount,string"`
	LikeCount     uint64 `json:"likeCount,string"`
	DislikeCount  uint64 `json:"dislikeCount,string"`
	FavoriteCount uint64 `json:"favoriteCount,string"`
	CommentCount  uint64 `json:"commentCount,string"`
}

type VideoPart string

const (
	VideoPartContentDetails       VideoPart = "contentDetails"
	VideoPartFileDetails          VideoPart = "fileDetails"
	VideoPartId                   VideoPart = "id"
	VideoPartLiveStreamingDetails VideoPart = "liveStreamingDetails"
	VideoPartLocalizations        VideoPart = "localizations"
	VideoPartPlayer               VideoPart = "player"
	VideoPartProcessingDetails    VideoPart = "processingDetails"
	VideoPartRecordingDetails     VideoPart = "recordingDetails"
	VideoPartSnippet              VideoPart = "snippet"
	VideoPartStatistics           VideoPart = "statistics"
	VideoPartStatus               VideoPart = "status"
	VideoPartSuggestions          VideoPart = "suggestions"
	VideoPartTopicDetails         VideoPart = "topicDetails"
)

type GetVideosConfig struct {
	Part                   []VideoPart
	Chart                  *string
	Id                     *string
	MyRating               *string
	H1                     *string
	MaxHeight              *uint64
	MaxResults             *uint64
	MaxWidth               *uint64
	OnBehalfOfContentOwner *string
	PageToken              *string
	RegionCode             *string
	VideoCategoryId        *string
}

func (service *Service) GetVideos(getVideosConfig *GetVideosConfig) (*[]Video, *string, *errortools.Error) {
	values := url.Values{}

	var videoParts []string
	for _, videoPart := range getVideosConfig.Part {
		videoParts = append(videoParts, string(videoPart))
	}
	values.Set("part", strings.Join(videoParts, ","))

	if getVideosConfig.Chart != nil {
		values.Set("chart", *getVideosConfig.Chart)
	}

	if getVideosConfig.Id != nil {
		values.Set("id", *getVideosConfig.Id)
	}

	if getVideosConfig.MyRating != nil {
		values.Set("myRating", *getVideosConfig.MyRating)
	}

	if getVideosConfig.H1 != nil {
		values.Set("h1", *getVideosConfig.H1)
	}

	if getVideosConfig.MaxHeight != nil {
		values.Set("maxHeight", fmt.Sprintf("%v", *getVideosConfig.MaxHeight))
	}

	if getVideosConfig.MaxResults != nil {
		values.Set("maxResults", fmt.Sprintf("%v", *getVideosConfig.MaxResults))
	}

	if getVideosConfig.MaxWidth != nil {
		values.Set("maxWidth", fmt.Sprintf("%v", *getVideosConfig.MaxWidth))
	}

	if getVideosConfig.OnBehalfOfContentOwner != nil {
		values.Set("onBehalfOfContentOwner", *getVideosConfig.OnBehalfOfContentOwner)
	}

	if getVideosConfig.PageToken != nil {
		values.Set("pageToken", *getVideosConfig.PageToken)
	}

	if getVideosConfig.RegionCode != nil {
		values.Set("regionCode", *getVideosConfig.RegionCode)
	}

	if getVideosConfig.VideoCategoryId != nil {
		values.Set("videoCategoryId", *getVideosConfig.VideoCategoryId)
	}

	videos := []Video{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlData("videos"),
		Parameters:    &values,
		ResponseModel: &videos,
	}

	_, _, response, e := service.httpRequestWrapped(&requestConfig)
	if e != nil {
		return nil, nil, e
	}

	return &videos, response.NextPageToken, nil
}
