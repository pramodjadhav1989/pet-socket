package flags

import (
	"os"

	"github.com/smartpet/websocket/constant"
	flag "github.com/spf13/pflag"
)

var (
	// env = flag.String(constant.EnvKey, constant.EnvDefaultValue,
	// 	constant.EnvUsage)
	port = flag.Int(constant.PortKey, constant.PortDefaultValue,
		constant.PortUsage)
	baseConfigPath = flag.String(constant.BaseConfigPathKey,
		constant.BaseConfigPathDefaultValue,
		constant.BaseConfigPathUsage)
	applicationMode = flag.String(constant.ModeKey, constant.ModeDefaultValue, constant.ModeUsage)
)

func init() {
	flag.Parse()
}

// Env is the application.yml runtime environment
func Env() string {
	env := os.Getenv(constant.Env)
	if env == "" {
		return constant.EnvDefaultValue
	}
	return env
}

// Port is the application.yml port number where the process will be started
func Port() int {
	return *port
}

// BaseConfigPath is the path that holds the configuration files
func BaseConfigPath() string {
	return *baseConfigPath
}

func ApplicationMode() string {
	return *applicationMode
}

func AWSRegion() string {
	region := os.Getenv(constant.AWSRegionKey)
	if region == "" {
		return constant.AWSRegionDefaultValue
	}
	return region
}

// AWSAccessKeyID is the access key id for aws
func AWSAccessKeyID() string {
	return os.Getenv(constant.AWSAccessKeyID)
}

// AWSSecretAccessKey is the secret access key for aws
func AWSSecretAccessKey() string {
	return os.Getenv(constant.AWSSecretAccessKey)
}

func AWSSessionToken() string {
	return os.Getenv(constant.AWSSessionToken)
}
