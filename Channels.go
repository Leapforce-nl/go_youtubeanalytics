package youtube

import (
	"fmt"
	"net/url"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	y_types "github.com/leapforce-libraries/go_youtube/types"
)

type Channel struct {
	Kind       string            `json:"kind"`
	Etag       string            `json:"etag"`
	ID         string            `json:"id"`
	Snippet    ChannelSnippet    `json:"snippet"`
	Statistics ChannelStatistics `json:"statistics"`
}

type ChannelSnippet struct {
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	CustomURL   string                  `json:"customUrl"`
	PublishedAt *y_types.DateTimeString `json:"publishedAt"`
	Thumbnails  Thumbnails              `json:"thumbnails"`
	Localized   Localized               `json:"localized"`
	Country     string                  `json:"country"`
}

type ChannelStatistics struct {
	CommentCount          uint64 `json:"commentCount,string"`
	ViewCount             uint64 `json:"viewCount,string"`
	SubscriberCount       uint64 `json:"subscriberCount,string"`
	HiddenSubscriberCount bool   `json:"hiddenSubscriberCount"`
	VideoCount            uint64 `json:"videoCount,string"`
}

type ChannelPart string

const (
	ChannelPartAuditDetails        ChannelPart = "auditDetails"
	ChannelPartBrandingSettings    ChannelPart = "brandingSettings"
	ChannelPartContentDetails      ChannelPart = "contentDetails"
	ChannelPartContentOwnerDetails ChannelPart = "contentOwnerDetails"
	ChannelPartID                  ChannelPart = "id"
	ChannelPartLocalizations       ChannelPart = "localizations"
	ChannelPartSnippet             ChannelPart = "snippet"
	ChannelPartStatistics          ChannelPart = "statistics"
	ChannelPartStatus              ChannelPart = "status"
	ChannelPartTopicDetails        ChannelPart = "topicDetails"
)

type GetChannelsConfig struct {
	Part                   []ChannelPart
	ForUserName            *string
	ID                     *string
	ManagedByMe            *bool
	Mine                   *bool
	H1                     *string
	MaxResults             *uint64
	OnBehalfOfContentOwner *string
}

func (service *Service) GetChannels(getChannelsConfig *GetChannelsConfig) (*[]Channel, *string, *errortools.Error) {
	values := url.Values{}

	channelParts := []string{}
	for _, channelPart := range getChannelsConfig.Part {
		channelParts = append(channelParts, string(channelPart))
	}
	values.Set("part", strings.Join(channelParts, ","))

	if getChannelsConfig.ForUserName != nil {
		values.Set("forUserName", *getChannelsConfig.ForUserName)
	}

	if getChannelsConfig.ID != nil {
		values.Set("id", *getChannelsConfig.ID)
	}

	if getChannelsConfig.ManagedByMe != nil {
		values.Set("managedByMe", fmt.Sprintf("%v", *getChannelsConfig.ManagedByMe))
	}

	if getChannelsConfig.Mine != nil {
		values.Set("mine", fmt.Sprintf("%v", *getChannelsConfig.Mine))
	}

	if getChannelsConfig.MaxResults != nil {
		values.Set("maxResults", fmt.Sprintf("%v", *getChannelsConfig.MaxResults))
	}

	if getChannelsConfig.OnBehalfOfContentOwner != nil {
		values.Set("onBehalfOfContentOwner", *getChannelsConfig.OnBehalfOfContentOwner)
	}

	channels := []Channel{}

	requestConfig := go_http.RequestConfig{
		URL:           service.urlData("channels"),
		Parameters:    &values,
		ResponseModel: &channels,
	}
	_, _, response, e := service.get(&requestConfig)
	if e != nil {
		return nil, nil, e
	}

	return &channels, response.NextPageToken, nil
}
