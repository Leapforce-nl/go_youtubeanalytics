package youtube

import (
	"fmt"
	"net/url"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	y_types "github.com/leapforce-libraries/go_youtube/types"
)

type Video struct {
	Kind       string          `json:"kind"`
	Etag       string          `json:"etag"`
	ID         string          `json:"id"`
	Snippet    VideoSnippet    `json:"snippet"`
	Statistics VideoStatistics `json:"statistics"`
}

type VideoSnippet struct {
	PublishedAt          *y_types.DateTimeString `json:"publishedAt"`
	ChannelID            string                  `json:"channelId"`
	Title                string                  `json:"title"`
	Description          string                  `json:"description"`
	Thumbnails           Thumbnails              `json:"thumbnails"`
	ChannelTitle         string                  `json:"channelTitle"`
	Tags                 []string                `json:"tags"`
	CategoryID           uint64                  `json:"categoryId,string"`
	LiveBroadcastContent string                  `json:"liveBroadcastContent"`
	DefaultLanguage      string                  `json:"defaultLanguage"`
	Localized            Localized               `json:"localized"`
	DefaultAudioLanguage string                  `json:"defaultAudioLanguage"`
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
	VideoPartID                   VideoPart = "id"
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
	ID                     *string
	MyRating               *string
	H1                     *string
	MaxHeight              *uint64
	MaxResults             *uint64
	MaxWidth               *uint64
	OnBehalfOfContentOwner *string
	PageToken              *string
	RegionCode             *string
	VideoCategoryID        *string
}

func (service *Service) GetVideos(getVideosConfig *GetVideosConfig) (*[]Video, *string, *errortools.Error) {
	values := url.Values{}

	videoParts := []string{}
	for _, videoPart := range getVideosConfig.Part {
		videoParts = append(videoParts, string(videoPart))
	}
	values.Set("part", strings.Join(videoParts, ","))

	if getVideosConfig.Chart != nil {
		values.Set("chart", *getVideosConfig.Chart)
	}

	if getVideosConfig.ID != nil {
		values.Set("id", *getVideosConfig.ID)
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

	if getVideosConfig.VideoCategoryID != nil {
		values.Set("videoCategoryId", *getVideosConfig.VideoCategoryID)
	}

	videos := []Video{}

	requestConfig := go_http.RequestConfig{
		URL:           service.urlData("videos"),
		Parameters:    &values,
		ResponseModel: &videos,
	}
	_, _, response, e := service.get(&requestConfig)
	if e != nil {
		return nil, nil, e
	}

	return &videos, response.NextPageToken, nil
}
