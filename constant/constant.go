package constant

import "time"

const (
	ApplicationName = "login"
	Port            = 8000
	TestMode        = "test"
	ReleaseMode     = "release"

	AWSAppConfig      = "appConfig"
	Application       = "login"
	ConfigIDKey       = "id"
	ConfigRegionKey   = "region"
	ConfigAccessKeyID = "accessKeyId"
	ConfigSecretKey   = "secretKey"
	ConfigAppKey      = "app"
	ConfigEnvKey      = "env"
	ConfigTypeKey     = "configType"
	ConfigNamesKey    = "configNames"

	DD_MM_YYYY                 = "02-01-2006"
	YYYYMMDD                   = "20060102"
	DDMMYYYY                   = "02012006"
	YYYY_MM_DD                 = "2006-01-02"
	YYYY_MM_DD_HH_MM_SS        = "2006-01-02 15:04:05"
	DD_MMM_YYYY                = "02-Jan-2006"
	YYYYMM                     = "2006-01"
	OTPTEMPLATE                = "resources/template/otp.html"
	SENDEREMAIL                = "support@advikapetworld.com"
	StatusRefreshTimeInSeconds = 2
)

const (
	SMSSvcTimeout     = 2 * time.Second
	SMSSvcRetry       = 2
	SMSSvcWaitTime    = 1 * time.Second
	SMSSvcWaitTimeMax = 2 * time.Second
)

const (
	USERID      = "userid"
	ACCESSTOKEN = "AccessToken"
)
