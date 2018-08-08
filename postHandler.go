package main

import (
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"log"
	"errors"
	"github.com/aws/aws-lambda-go/lambda"
)

func postUserProperties(jsonStr string) {
	sendPost("identification", "/identify", jsonStr)
}

func postEvent(jsonStr string) {
	sendPost("event", "/httpapi", jsonStr)
}

// make HTTP post -- find API key in Amplitude project settings
func sendPost(postType string, urlPath string, jsonStr string) {
	var fullUrl = amplitudeUrlHost + urlPath
	fmt.Println("URL:>", fullUrl)
	var jsonBytes = []byte(jsonStr)
	fmt.Println("JSON being sent: ", jsonStr)
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonBytes))
	fmt.Println("Destination of POST: " + fullUrl)
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add(postType, jsonStr)
	q.Add("api_key", "70f8f71d95f5a2910ee56e56f9d24756")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	// Non-2XX response does not cause an error. Error returned if too many redirects or HTTP protocol error.
	// (see article)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Debug("response Status:", resp.Status)
	log.Debug("response Headers:", resp.Header)

	// Handle each categories of responses
	switch e := resp.StatusCode {
		case 200:
			// we need to be careful about how we use ioutil.ReadAll, because of memory constraints
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("response Body:", string(body))	
		case 400:
			log.Error("Malformed JSON error")
			// try the next one
			return nil
		case 500 <= e && e <= 511:
			// e.g. Amp down (network error)
			return errors.New("Amplitude 500 Error: %s", string(resp.Body))
		default:
			return errors.New("Unknown Error %d. %s", e, resp)
		}
}
