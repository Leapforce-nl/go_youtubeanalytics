package YoutubeAnalytics

import (
	"net/http"

	bigquerytools "github.com/Leapforce-nl/go_bigquerytools"

	oauth2 "github.com/Leapforce-nl/go_oauth2"
)

const (
	apiName         string = "YoutubeAnalytics"
	apiURL          string = "https://youtubeanalytics.googleapis.com/v2"
	authURL         string = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenURL        string = "https://oauth2.googleapis.com/token"
	tokenHTTPMethod string = http.MethodPost
	redirectURL     string = "http://localhost:8080/oauth/redirect"
)

// YoutubeAnalytics stores YoutubeAnalytics configuration
//
type YoutubeAnalytics struct {
	oAuth2 *oauth2.OAuth2
}

// methods
//
func NewYoutube(clientID string, clientSecret string, scope string, bigQuery *bigquerytools.BigQuery, isLive bool) (*YoutubeAnalytics, error) {
	yt := YoutubeAnalytics{}
	yt.oAuth2 = oauth2.NewOAuth(apiName, clientID, clientSecret, scope, redirectURL, authURL, tokenURL, tokenHTTPMethod, bigQuery, isLive)
	return &yt, nil
}

func (yt *YoutubeAnalytics) ValidateToken() (*oauth2.Token, error) {
	return yt.oAuth2.ValidateToken()
}

func (yt *YoutubeAnalytics) InitToken() error {
	return yt.oAuth2.InitToken()
}

func (yt *YoutubeAnalytics) Get(url string, model interface{}) (*http.Response, error) {
	res, err := yt.oAuth2.Get(url, model)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (yt *YoutubeAnalytics) Patch(url string, model interface{}) (*http.Response, error) {
	res, err := yt.oAuth2.Patch(url, nil, model)

	if err != nil {
		return nil, err
	}

	return res, nil
}
