package twitch

import (
	"os"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	stateCallbackKey = "oauth-state-callback"
	oauthSessionName = "oauth-session"
	oauthTokenKey    = "oauth-token"
	scopes           = "channel:read:subscriptions bits:read channel:read:redemptions user:read:follows"
)

var (
	// secrets
	cookieSecret string
	clientID     string
	clientSecret string // idk
	twitchUser   string

	serverType  string // global
	userID      string
	baseUrl     string
	token       string
	Port        string
	channelInfo ChannelInfo
)

func getSecrets() {
	clientID = os.Getenv("TWITCH_CLIENT_ID")
	clientSecret = os.Getenv("TWITCH_CLIENT_SECRET")
	cookieSecret = os.Getenv("COOKIE_SECRET")
	twitchUser = os.Getenv("TWITCH_USER")
}

var (
	// Consider storing the secret in an environment variable or a dedicated storage system.
	oauth2Config *clientcredentials.Config
	cookieStore  *sessions.CookieStore
)
