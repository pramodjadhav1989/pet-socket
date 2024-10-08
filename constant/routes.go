package constant

const (
	Flag           = "ENV"
	ActuatorRoute  = "/actuator/*any"
	LoginWithOTP   = "/login/loginWithOTP"
	VerifyLoginOTP = "/login/VerifyOTP"
	RefreshToken   = "/login/RefreshToken"

	ProviderBusinessDetails = "/register/provider/businessdetails"
	ProviderContact         = "/register/provider/contacts"
	ProviderAddDocuments    = "/register/provider/adddocument"
	ProviderUpdateDocuments = "/register/provider/updatedocument"
	ProviderDeleteDocuments = "/register/provider/deletedocument"
	ProviderMediaLinks      = "/register/provider/medialinks"
	ProviderSessions        = "/register/provider/addsession"
	ProviderSessionsDetails = "/register/provider/addsessiondetails"
	ProviderGetProfile      = "/register/provider/getprofile"
)
