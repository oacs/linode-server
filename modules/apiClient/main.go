package apiClient

import (
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var httpClient = &http.Client{}

func Get(url string, headers http.Header) ([]byte, error) {
	return request("GET", url, headers, nil)
}

func Post(url string, body io.Reader, headers http.Header) ([]byte, error) {
	return request("POST", url, headers, body)
}

func request(verbose string, url string, headers http.Header, body io.Reader) (responseData []byte, err error) {
	req, err := http.NewRequest(verbose, url, body)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal("error calling get url "+url, err)
		return nil, err
	}

	responseData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("error reading response body", err)
		return nil, err
	}

	return responseData, nil
}
