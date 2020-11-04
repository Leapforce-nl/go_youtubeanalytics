package Youtube

import (
	"net/http"

	bigquerytools "github.com/leapforce-libraries/go_bigquerytools"

	oauth2 "github.com/leapforce-libraries/go_oauth2"
)

const (
	apiName         string = "Youtube"
	apiURL          string = "https://youtubeanalytics.googleapis.com/v2"
	authURL         string = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenURL        string = "https://oauth2.googleapis.com/token"
	tokenHTTPMethod string = http.MethodPost
	redirectURL     string = "http://localhost:8080/oauth/redirect"
)

// Youtube stores Youtube configuration
//
type Youtube struct {
	oAuth2 *oauth2.OAuth2
}

// methods
//
func NewYoutube(clientID string, clientSecret string, scope string, bigQuery *bigquerytools.BigQuery, isLive bool) (*Youtube, error) {
	yt := Youtube{}
	config := oauth2.OAuth2Config{
		ApiName:         apiName,
		ClientID:        clientID,
		ClientSecret:    clientSecret,
		Scope:           scope,
		RedirectURL:     redirectURL,
		AuthURL:         authURL,
		TokenURL:        tokenURL,
		TokenHTTPMethod: tokenHTTPMethod,
	}
	yt.oAuth2 = oauth2.NewOAuth(config, bigQuery, isLive)
	return &yt, nil
}

func (yt *Youtube) ValidateToken() (*oauth2.Token, error) {
	return yt.oAuth2.ValidateToken()
}

func (yt *Youtube) InitToken() error {
	return yt.oAuth2.InitToken()
}

func (yt *Youtube) Get(url string, model interface{}) (*http.Response, error) {
	res, err := yt.oAuth2.Get(url, model)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (yt *Youtube) Patch(url string, model interface{}) (*http.Response, error) {
	res, err := yt.oAuth2.Patch(url, nil, model)

	if err != nil {
		return nil, err
	}

	return res, nil
}
