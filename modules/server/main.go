package server

import (
	// "fmt"
	// "io"
	"net/http"

	"example.com/m/v2/modules/twitch"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	baseURL string
	routes  map[string]func(http.ResponseWriter, *http.Request)
	status  string
	env     string
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("This is an example server.\n"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func NewServer(baseURL string, env string) *Server {
	s := &Server{
		env:     env,
		baseURL: baseURL,
		routes:  make(map[string]func(http.ResponseWriter, *http.Request)),
	}
	s.routes["/"] = HelloServer
	s.routes["/twitch/ouath/login"] = twitch.HandleLogin
	s.routes["/twitch/ouath/callback"] = twitch.HandleOAuth2Callback
	return s
}

func (s *Server) Start() {
	log.SetLevel(log.DebugLevel)

	// modules init
	twitch.Init()

	switch s.env {
	case "dev":
		log.Info("Running in dev mode")
		log.Info("Starting server on port 8080")
		log.Fatal(http.ListenAndServe(":8081", nil))
	case "prod":
		log.Info("Running in prod mode")
		log.Info("Starting server on port 8080")
		log.Fatal(http.ListenAndServeTLS(":4443", "/etc/letsencrypt/live/tucos.dev/fullchain.pem", "/etc/letsencrypt/live/tucos.dev/privkey.pem", nil))
	}
}
