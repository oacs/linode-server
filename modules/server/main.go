package server

import (
	// "fmt"
	// "io"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func Main() {
	http.HandleFunc("/", HelloServer)
	err := http.ListenAndServeTLS(":4443", "/etc/letsencrypt/live/tucos.dev/fullchain.pem", "/etc/letsencrypt/live/tucos.dev/privkey.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
