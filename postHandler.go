package main

import (
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"log"
	"github.com/aws/aws-lambda-go/lambda"
)

func postUserProperties(jsonStr string) {
	sendPost("identification", "/identify", jsonStr)
}

func postEvent(jsonStr string) {
	sendPost("event", "/httpapi", jsonStr)
}

/* TODO: Create a custom Status error type and integrate with Lambda
 */
// Error has methods for HTTP status codes, embeds built-in error interface
type Error interface {
	error
	Status() int
}

// Error with an HTTP status code
type StatusError struct {
	Code int
	Err error
}

// Retrieve the http error from a StatusError
func (s StatusError) Error() string {
	return s.Err.error
}

// Return http code from a StatusError
func (s StatusError) Status() int{
	return s.code
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

/* TODO: Add functionality to differentiate between different HTTP errors (between Lambda<>Amplitude), example:
 * 400 - malformed json (bad request). Log info then return success/null/cary-on
 * 5xx - network error? Means amp is down...
 */
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
