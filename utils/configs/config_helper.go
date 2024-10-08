package configs

import (
	"errors"
	"fmt"

	"github.com/smartpet/websocket/constant"
	"github.com/smartpet/websocket/models"
)

func GetAWSClientCredentials(name string) (error, models.AWSClientCredentialsConfig) {
	awsConfig, err := Get(constant.AWSConfig)
	if err != nil {
		return err, models.AWSClientCredentialsConfig{}
	}
	AWS_Config := awsConfig.Sub(name)
	config := models.AWSClientCredentialsConfig{
		Region:          GetStringWithEnv(AWS_Config.GetString(constant.AWSRegionKey)),
		AccessId:        GetStringWithEnv(AWS_Config.GetString(constant.AWSAccessKeyID)),
		SecretAccessKey: GetStringWithEnv(AWS_Config.GetString(constant.AWSSecretAccessKey)),
	}
	return nil, config
}

func GetAWSConfig(name string, service string) (string, error) {
	awsConfig, err := Get(constant.AWSConfig)
	if err != nil {
		return "", err
	}
	AWS_Config := awsConfig.Sub(name)
	config_value := GetStringWithEnv(AWS_Config.GetString(service))
	return config_value, nil
}

func GetAppConfig(service string, isSecure bool) (string, error) {
	secretConfig, err := Get(constant.ApplicationConfig)
	if err != nil {
		return "", err
	}
	if isSecure {
		return GetStringWithEnv(secretConfig.GetString(service)), nil
	}
	return secretConfig.GetString(service), nil

}

func GetExternalAPI(appname string, service string, isSecure bool) (string, error) {
	apiConfig, err := Get(constant.ExternalConfig)
	if err != nil {
		return "", err
	}
	API_Config := apiConfig.Sub(appname)
	fmt.Println(API_Config.GetString(service))
	if isSecure {
		return GetStringWithEnv(API_Config.GetString(service)), nil
	}

	return API_Config.GetString(service), nil
}

func GetConfigWithValue(config string, appname string, service string) (string, error) {
	apiConfig, err := GetClient().Get(config, appname)
	if err != nil {
		return "", err
	}
	value := apiConfig.(map[string]interface{})[service]
	if value == nil {
		return "", errors.New(fmt.Sprintf("%s key not found %s config", service, config))
	}
	return value.(string), nil
}
