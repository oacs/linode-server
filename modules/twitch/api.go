package twitch

import (
	"context"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"example.com/m/v2/modules/apiClient"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

const TWITCH_TOKEN_URL = "https://id.twitch.tv/oauth2/token"
const TWITCH_API_URL = "https://api.twitch.tv/helix"
const REDIRECT_URL = "https://tucos.dev/api/twitch/oauth/callback"

func fetchTwitchChannelInfo() (channelInfo ChannelInfo, err error) {
	if accessToken == "" {
		err = errors.New("no access token")
		return channelInfo, err
	}
	log.Debug("fetching Twitch channel info\n")
	headers := http.Header{
		"Authorization": {"Bearer " + accessToken},
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

	accessToken = authInfo.AccessToken
}

// HandleLogin is a Handler that redirects the user to Twitch for login, and provides the 'state'
// parameter which protects against login CSRF.
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, oauthSessionName)
	if err != nil {
		log.Error("corrupted session %s -- generated new", err)
		err = nil
	}

	var tokenBytes [255]byte
	if _, err := rand.Read(tokenBytes[:]); err != nil {
		log.Error("error generating random bytes: %s", err)
		return
	}

	state := hex.EncodeToString(tokenBytes[:])

	session.AddFlash(state, stateCallbackKey)

	if err = session.Save(r, w); err != nil {
		log.Error("error saving session: %s", err)
		return
	}

	http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusTemporaryRedirect)

	return
}

// HandleOauth2Callback is a Handler for oauth's 'redirect_uri' endpoint;
// it validates the state token and retrieves an OAuth token from the request parameters.
func HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, oauthSessionName)
	if err != nil {
		log.Error("corrupted session %s -- generated new", err)
		err = nil
	}

	// ensure we flush the csrf challenge even if the request is ultimately unsuccessful
	defer func() {
		if err := session.Save(r, w); err != nil {
			log.Error("error saving session: %s", err)
		}
	}()

	switch stateChallenge, state := session.Flashes(stateCallbackKey), r.FormValue("state"); {
	case state == "", len(stateChallenge) < 1:
		err = errors.New("missing state challenge")
	case state != stateChallenge[0]:
		err = fmt.Errorf("invalid oauth state, expected '%s', got '%s'\n", state, stateChallenge[0])
	}

	if err != nil {
		log.Error("error validating state: %s", err)
		return
	}

	token, err := oauth2Config.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Error("error exchanging code for token: %s", err)
		return
	}

	// add the oauth token to session
	session.Values[oauthTokenKey] = token

	log.Debug("Access token: %s\n", token.AccessToken)
	return
}

func initOAuth() {
	// Gob encoding for gorilla/sessions
	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       strings.Split(scopes, " "),
		Endpoint:     twitch.Endpoint,
		RedirectURL:  REDIRECT_URL,
	}
}
