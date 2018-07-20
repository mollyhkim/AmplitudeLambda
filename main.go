package main

import (
	"time"
	"fmt"
	"strconv"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"log"
)

var milliseconds int32
var testJson string
var urlHost = "https://api.amplitude.com"

func postUserProperties(jsonStr string) {
	sendPost("identification", "/identify", jsonStr)
}

func postEvent(jsonStr string) {
	sendPost("event", "/httpapi", jsonStr)
}

// make HTTP post -- find API key in Amplitude project settings
func sendPost(postType string, urlPath string, jsonStr string) {
	var fullUrl = urlHost + urlPath
	fmt.Println("URL:>", fullUrl)
	var jsonBytes = []byte(jsonStr)
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

type jsonInfo struct {
	User_id                   string `json: "user_id"`
	Event_type                string `json: "event_type"`
	Language                  string `json: "language"`
	Country                   string `json: "country"`
	Region                    string `json: "region"`
	Ip                        string `json: "ip"`
	Time                      string `json: "time"`
	Person_id                 string `json: person_id`
	Loan_application_id       string `json: loan_application_id`
	Kyc_status                string `json: KYC_status`
	Account_validation_status string `json: Account_validation_status`
}


// Parse Json received from Comandante
func parseComandanteJson(cmdJson string) jsonInfo {
	var jsonData []jsonInfo
	err := json.Unmarshal([]byte(cmdJson), &jsonData)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData[0]
}

func createAmplitudeJson(originalJson jsonInfo) string {
	var newJson string

	switch originalJson.Event_type {
	case "Facebook":
		newJson = `[{"user_id":"` + originalJson.User_id + `"}]`
	case "KYC":
		newJson = `[{"user_id":"` + originalJson.User_id + `", "loan_application_id":"` + originalJson.Loan_application_id + `"}]`
	case "KYC Review":
		newJson = `[{"user_id":"` + originalJson.User_id + `", "loan_application_id":"` + originalJson.Loan_application_id + `",
					"KYC_status":"` + originalJson.Kyc_status + `"}]`
	case "Cashout":
		newJson = `[{"user_id":"` + originalJson.User_id + `", "loan_application_id":"` + originalJson.Loan_application_id + `"}]`
	case "Account Validation":
		newJson = `[{"user_id":"` + originalJson.User_id + `", "loan_application_id":"` + originalJson.Loan_application_id + `",
					"account_validation_status":"` + originalJson.Account_validation_status + `"}]`
	default:
		log.Fatal("ERR: sendCorrespondingInfo: Invalid event type")
	}

	postUserProperties(newJson)
	return newJson
}

func main() {
	milliseconds := time.Now().UTC().UnixNano() / 1000000
	fmt.Println(milliseconds)
	testJson := `[{"user_id":"test-user-molly", "person_id": "2424242", "KYC Status": "APPROVED", 
				"event_type":"Facebook","language": "ANGLAIS", "country":"DENMARK",
				"region": "LA County","ip":"188.0.0.1","time": "` + strconv.Itoa(int(milliseconds)) + `"}]`
	postUserProperties(testJson)
	postEvent(testJson)
	parsed := parseComandanteJson(testJson)
	fmt.Println(parsed)
	var newJson = createAmplitudeJson(parsed)

	fmt.Println("OG JSON")
	fmt.Println(testJson)
	fmt.Println("PARSED JSON")
	fmt.Println(parsed)
	fmt.Println("SENT JSON")
	fmt.Println(newJson)

	fmt.Println("AAAAND WE'RE DONE! :)")

}
