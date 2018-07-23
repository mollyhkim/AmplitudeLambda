package main

import (
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
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
	fmt.Println("WHAT IM SENDING")
	fmt.Println(jsonStr)
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonBytes))
	fmt.Println("posting to " + fullUrl)
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add(postType, jsonStr)
	q.Add("api_key", "70f8f71d95f5a2910ee56e56f9d24756")
	req.URL.RawQuery = q.Encode()

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
