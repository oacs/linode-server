package server

import (
	// "fmt"
	// "io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"example.com/m/v2/modules/twitch"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func Main() {
	log.SetLevel(log.DebugLevel)
	twitch.Init()
	http.HandleFunc("/", HelloServer)

	// twitch
	http.HandleFunc("/twitch/oauth/login", twitch.HandleLogin)
	http.HandleFunc("/twitch/oauth/callback", twitch.HandleOAuth2Callback)

	err := http.ListenAndServeTLS(":4443", "/etc/letsencrypt/live/tucos.dev/fullchain.pem", "/etc/letsencrypt/live/tucos.dev/privkey.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
