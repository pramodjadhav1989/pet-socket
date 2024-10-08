package constant

var ErrorCodeMap = map[string]string{
	"ABP11000": "Internal Server Error",
	"ABP11001": "Invalid Parameters",
	"ABP11002": "Unknown Party Code",
	"ABP11003": "Unauthorized Access",
	"ABP11004": "Request validation failed",
	"ABP11005": "Data not available",
	"ABP11006": "Failed to set the otp to inactive ",
	"ABP11007": "failed to insert the generated otp to db ",
	"ABP11008": "Invalid Session ID",
}

var SMSErrorCodeMap = map[string]string{
	"ES1001":      "ES1001 Authentication Failed (invalid username/password)",
	"ES1004":      "ES1004 Invalid Senderid",
	"ES1009":      "ES1009 Sorry unable to process request",
	"ES1013":      "ES1013 Template id is invalid",
	"ES1002":      "ES1002 Unauthorized Usage - insufficient privilege",
	"ES1007":      "ES1007 Account Deactivated",
	"MSGBLANK":    "Message is blank",
	"ACEXPIRED":   "Account is Expire",
	"LIMITEXCEED": "You have Exceeded your SMS Limit.",
}
