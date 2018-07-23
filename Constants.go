package main

const amplitudeUrlHost = "https://api.amplitude.com"

// Event type constants
const (
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

// Property type constants
const (
	LOAN_APPLICATION_ID     string = "Loan Application ID"
	LOAN_APPLICATION_NUMBER string = "Loan Application Number"
	LOAN_NUMBER             string = "Loan Number"
	ACQUISITION_SOURCE      string = "Acquisition Source"
)
