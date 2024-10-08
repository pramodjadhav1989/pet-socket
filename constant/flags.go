package constant

// Flag constants
const (
	EnvDefaultValue            = "uat"
	EnvUsage                   = "application.yml runtime environment"
	PortKey                    = "port"
	PortDefaultValue           = 8001
	PortUsage                  = "application.yml port"
	BaseConfigPathKey          = "base-config-path"
	BaseConfigPathDefaultValue = "resources"
	BaseConfigPathUsage        = "path to folder that stores your configurations"
	ModeKey                    = "mode"
	ModeDefaultValue           = "test"
	ModeUsage                  = "run mode of the application, can be test or release"
	Env                        = "env"
	AWSRegionKey               = "s3-region"
	AWSRegionDefaultValue      = "ap-south-1"
	AWSAccessKeyID             = "ACCESS_KEY_ID"
	AWSSecretAccessKey         = "SECRET_KEY"
	AWSSessionToken            = "SESSION_TOCKEN"
	AWSSBucketName             = "s3-bucket"
	SNSTopic                   = "sns-topic"
	AWSEnvName                 = "APP_ENV"
	AppType                    = "app-type"
	SecrateName                = "secrateName"
)
