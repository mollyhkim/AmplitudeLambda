package main

import (
	"encoding/json"
	"log"
	"time"
	"strconv"
)

func parsingInputFromCMD(cmdJson string) jsonInfo {
	var jsonData []jsonInfo
	err := json.Unmarshal([]byte(cmdJson), &jsonData)
	if err != nil {
		log.Fatal(err)
	}

	// check for USER ID
	if len(jsonData[0].User_id) <= 0 {
		log.Fatal("User id cannot be null")
	}

	return jsonData[0]
}

// creates Json for EVENT post
func createEventJson(inputJson jsonInfo) string { // so far these are events so must post event type
	milliseconds := strconv.Itoa(int(time.Now().UTC().UnixNano() / 1000000))

	// check for POST TYPE
	if inputJson.Post_type != "event" {
		log.Fatal(`createPropertiesJson incompatible for json with postType:` + inputJson.Post_type)
	}

	var newJson string
	// check for EVENT TYPE
	eventType := inputJson.Event_type
	switch eventType {
	case FACEBOOK_CONNECT, FACEBOOK_LOGIN:
		newJson = `[{"user_id":"` + inputJson.User_id + `", "event_type":"` + eventType + `", "time": "` + milliseconds + `"}]`
	case KYC_REQUEST, KYC_SUBMIT, KYC_REVIEW, CASHOUT_SUBMIT, CASHOUT_REQUEST, ACCOUNT_VALIDATION_START:
		newJson = `[{"user_id":"` + inputJson.User_id + `", "event_type":"` + eventType + `",
					"loan_application_id":"` + inputJson.Loan_application_id + `", "time": "` + milliseconds + `"}]`
	case KYC_DONE:
		// check for KYC STATUS when appropriate
		if len(inputJson.Kyc_status) <= 0 {
			log.Fatal("Invalid KYC Status for event type KYC: Done")
		}
		newJson = `[{"user_id":"` + inputJson.User_id + `",  "event_type":"` + eventType + `",
					"loan_application_id":"` + inputJson.Loan_application_id + `", 
					"KYC_status":"` + inputJson.Kyc_status + `", "time": "` + milliseconds + `"}]`
	case ACCOUNT_VALIDATION_DONE:
		// check for ACCOUNT VALIDATION STATUS when appropriate
		if len(inputJson.Account_validation_status) <= 0 {
			log.Fatal("Invalid KYC Status for event type Account Validation: Done")
		}
		newJson = `[{"user_id":"` + inputJson.User_id + `",  "event_type":"` + eventType + `",
					"loan_application_id":"` + inputJson.Loan_application_id + `",
					"account_validation_status":"` + inputJson.Account_validation_status + `", "time": "` + milliseconds + `"}]`
	default:
		log.Fatal("Invalid / null event type for createEventJson")
	}

	postEvent(newJson)
	return newJson
}

// creates Json for IDENTIFICATION post
func createPropertiesJson(inputJson jsonInfo) string {
	// check for POST TYPE
	if inputJson.Post_type != "identification" {
		log.Fatal(`createPropertiesJson incompatible for json with postType:` + inputJson.Post_type)
	}

	var newJson string
	// check for PROPERTY TYPE
	switch inputJson.Property_type {
	case LOAN_APPLICATION_ID:
		newJson = `[{"user_id":"` + inputJson.User_id + `"}]`
	case LOAN_APPLICATION_NUMBER:
		newJson = `[{"user_id":"` + inputJson.User_id + `", "loan_application_id":"` + inputJson.Loan_application_id + `"}]`
	case LOAN_NUMBER:
		newJson = `[{"user_id":"` + inputJson.User_id + `", "loan_application_id":"` + inputJson.Loan_application_id + `",
					"KYC_status":"` + inputJson.Kyc_status + `"}]`
	case ACQUISITION_SOURCE:
		newJson = `[{"user_id":"` + inputJson.User_id + `", "loan_application_id":"` + inputJson.Loan_application_id + `",
					"account_validation_status":"` + inputJson.Account_validation_status + `"}]`
	default:
		log.Fatal(`Invalid / null event type:` + inputJson.Event_type + `for createPropertiesJson`)
	}

	postUserProperties(newJson)
	return newJson
}
