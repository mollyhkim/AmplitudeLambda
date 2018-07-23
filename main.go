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

var urlHost = "https://api.amplitude.com"

var (
	FACEBOOK_LOGIN           string = "Facebook: Login"
	FACEBOOK_CONNECT         string = "Facebook: Connect"
	KYC_REQUEST              string = "KYC: Request"
	KYC_SUBMIT               string = "KYC: Submit"
	KYC_REVIEW               string = "KYC: Review"
	KYC_DONE                 string = "KYC: Done"
	CASHOUT_REQUEST          string = "Cashout: Request"
	CASHOUT_SUBMIT           string = "Cashout: Submit"
	ACCOUNT_VALIDATION_START string = "Account Validation: Start"
	ACCOUNT_VALIDATION_DONE  string = "Account Validation: Done"
)


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
	Account_validation_status string `json: Account_validation_status`
	Kyc_status                string `json: KYC_status`
}

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

// Parse Json received from Comandante
func parseComandanteJson(cmdJson string) jsonInfo {
	var jsonData []jsonInfo
	err := json.Unmarshal([]byte(cmdJson), &jsonData)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData[0]
}

func createEventJson(inputJson jsonInfo) string { // so far these are events so must post event type
	var newJson string
	eventType := inputJson.Event_type
	switch eventType {
	case "Facebook: Temp":
		newJson = `[{ "user_id":"` + inputJson.User_id + `", "event_type":"updateevent4", "user_properties":
			{"$set": {"Cohort":"Test B"}, "$add": {"friendCount":23}, "testingproperty":"testing", "country":"Mali"}}]`
	case FACEBOOK_CONNECT, FACEBOOK_LOGIN:
		newJson = `[{"user_id":"` + inputJson.User_id + `", "event_type":"` + eventType + `"}]`
	case KYC_REQUEST, KYC_SUBMIT, KYC_REVIEW, CASHOUT_SUBMIT, CASHOUT_REQUEST, ACCOUNT_VALIDATION_START:
		newJson = `[{"user_id":"` + inputJson.User_id + `", "loan_application_id":"` + inputJson.Loan_application_id + `"}]`
	case KYC_DONE:
		newJson = `[{"user_id":"` + inputJson.User_id + `", "loan_application_id":"` + inputJson.Loan_application_id + `",
					"KYC_status":"` + inputJson.Kyc_status + `"}]`
	case ACCOUNT_VALIDATION_DONE:
		newJson = `[{"user_id":"` + inputJson.User_id + `", "loan_application_id":"` + inputJson.Loan_application_id + `",
					"account_validation_status":"` + inputJson.Account_validation_status + `"}]`
	default:
		log.Fatal("ERR: sendCorrespondingInfo: Invalid event type")
	}

	postEvent(newJson)
	return newJson
}

func createPropertiesJson(inputJson jsonInfo) string {
	var newJson string

	switch inputJson.Event_type {
	case "Facebook: Temp":
		newJson = `[{ "user_id":"` + inputJson.User_id + `", "additionalAttribute12":"Just checking for properties update!"}]`
	case FACEBOOK_LOGIN, FACEBOOK_CONNECT:
		newJson = `[{"user_id":"` + inputJson.User_id + `"}]`
	case KYC_REQUEST, KYC_SUBMIT, KYC_REVIEW, CASHOUT_REQUEST, CASHOUT_SUBMIT, ACCOUNT_VALIDATION_START:
		newJson = `[{"user_id":"` + inputJson.User_id + `", "loan_application_id":"` + inputJson.Loan_application_id + `"}]`
	case KYC_DONE:
		newJson = `[{"user_id":"` + inputJson.User_id + `", "loan_application_id":"` + inputJson.Loan_application_id + `",
					"KYC_status":"` + inputJson.Kyc_status + `"}]`
	case ACCOUNT_VALIDATION_DONE:
		newJson = `[{"user_id":"` + inputJson.User_id + `", "loan_application_id":"` + inputJson.Loan_application_id + `",
					"account_validation_status":"` + inputJson.Account_validation_status + `"}]`
	default:
		log.Fatal("ERR: sendCorrespondingInfo: Invalid event type")
	}

	postUserProperties(newJson)
	return newJson
}

func handler() {
	milliseconds := time.Now().UTC().UnixNano() / 1000000
	fmt.Println(milliseconds)
	testJson := `[{"user_id":"test-user-molly-6", "person_id": "2424242", "KYC_status": "APPROVED", 
				"event_type":"Facebook: Connect","language": "ANGLAIS", "country":"ENGLAND",
				"region": "LA County","ip":"188.0.0.1","time": "` + strconv.Itoa(int(milliseconds)) + `"}]`

	parsed := parseComandanteJson(testJson)
	var newJson = createEventJson(parsed)
	//var newJson = createPropertiesJson(parsed)
	fmt.Println("OG JSON")
	fmt.Println(testJson)
	fmt.Println("PARSED JSON")
	fmt.Println(parsed)
	fmt.Println("SENT JSON")
	fmt.Println(newJson)
	fmt.Println("AAAAND WE'RE DONE! :)")

}

func main() {
	handler()
}
