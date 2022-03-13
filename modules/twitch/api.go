package twitch

import (
	"encoding/gob"
	"encoding/json"
	"net/http"
	"net/url"

	"example.com/m/v2/modules/apiClient"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const TWITCH_TOKEN_URL = "https://id.twitch.tv/oauth2/token"
const TWITCH_API_URL = "https://api.twitch.tv/helix"

func fetchTwitchChannelInfo(token string) (channelInfo ChannelInfo, err error) {
	log.Debug("fetching Twitch channel info\n")
	headers := http.Header{
		"Authorization": {"Bearer " + token},
		"Client-ID":     {clientID},
	}
	responseData, err := apiClient.Get(TWITCH_API_URL+"/users?login="+twitchUser, headers)
	if err != nil {
		return channelInfo, err
	}

	log.Debug("response from channel info ", string(responseData))

	var channelArrayInfo struct {
		Data []ChannelInfo `json:"data"`
	}

	json.Unmarshal(responseData, &channelArrayInfo)

	return channelArrayInfo.Data[0], nil
}

func fetchApiToken() (string, error) {
	cookieStore = sessions.NewCookieStore([]byte(cookieSecret))
	gob.Register(&oauth2.Token{})

	log.Debug("fetching Twitch API token from ", TWITCH_TOKEN_URL)

	// create url with query parameters
	base, err := url.Parse(TWITCH_TOKEN_URL)
	if err != nil {
		return "", err
	}

	// Query params
	params := url.Values{}
	params.Add("client_secret", clientSecret)
	params.Add("client_id", clientID)
	params.Add("grant_type", "client_credentials")
	params.Add("scopes", scopes)
	base.RawQuery = params.Encode()

	headers := http.Header{
		"Content-Type": {"application/json"},
	}

	responseData, err := apiClient.Post(base.String(), nil, headers)
	if err != nil {
		return "", err
	}

	log.Debug("Response body: ", string(responseData))

	var authInfo struct {
		AccessToken string `json:"access_token"`
	}

	json.Unmarshal(responseData, &authInfo)

	return authInfo.AccessToken, nil

	// 	oauth2Config = &clientcredentials.Config{
	// 		ClientID:     clientID,
	// 		ClientSecret: clientSecret,
	// 		TokenURL:     TWITCH_TOKEN_URL,
	// 	}

	// 	token, err := oauth2Config.Token(context.Background())
	// 	if err == nil {
	// 		log.Debug("Access token: %s\n", token.AccessToken)
	// 	} else {
	// 		log.Fatal(err)
	// 	}
	// 	return token.AccessToken, err
}
