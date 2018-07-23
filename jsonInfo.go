package main

type jsonInfo struct {
	User_id                   string `json: "user_id"`
	Post_type                 string `json: "post_type"` 				// either "event" or "identification"
	Event_type                string `json: "event_type"`
	Property_type             string `json: "property_type"`
	Country                   string `json: "country"`
	Time                      string `json: "time"`
	Person_id                 string `json: "person_id"`
	Loan_application_id       string `json: "loan_application_id"`       // available after "Score User"
	Account_validation_status string `json: "account_validation_status"` // available after ""
	Kyc_status                string `json: "KYC_status"`                // available after ""
}
