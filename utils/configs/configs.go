package configs

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/smartpet/websocket/constant"
	config "github.com/smartpet/websocket/utils/clientconfigs"
	"github.com/smartpet/websocket/utils/flags"
	"github.com/spf13/viper"
)

type providers struct {
	providers map[string]*viper.Viper
	mu        sync.Mutex
}

var baseConfigPath string
var p *providers

// Init is used to initialize the configurations
func Init(path string) {
	baseConfigPath = path
	p = &providers{
		providers: make(map[string]*viper.Viper),
	}
}

// Get is used to get the instance to the config provider for the
// configuration name
func Get(name string) (*viper.Viper, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// see for an existing provider
	if provider, ok := p.providers[name]; ok {
		// provider already exists
		// use it
		return provider, nil
	}

	// try to get the provider
	provider := viper.New()
	provider.SetConfigName(name)
	provider.AddConfigPath(baseConfigPath)
	err := provider.ReadInConfig()
	if err != nil {
		// config not found or some other parsing errors
		return nil, fmt.Errorf("config %s error : %w", name, err)
	}

	// add a watcher for this provider
	provider.WatchConfig()

	// successfully found config, store it for future use
	p.providers[name] = provider

	return provider, nil
}

// GetStringWithEnv is used to get the value from environment variables
func GetStringWithEnv(key string) string {
	s := os.Getenv(key[2 : len(key)-1])
	return s
}

type Client struct {
	config.Client
	env map[string]string
}

// appConfigClient is the instance of the config client to be used by the application
var client *Client

// InitTestModeConfigs is used to initialize the configs
func InitTestModeConfigs(directory string, configNames ...string) error {
	c, err := config.New(config.Options{
		Provider: config.FileBased,
		Params: map[string]interface{}{
			"configsDirectory":      directory,
			constant.ConfigNamesKey: configNames,
			constant.ConfigTypeKey:  "yaml",
		},
	})
	if err != nil {
		return err
	}
	client = getClient(c)
	return nil
}

func InitReleaseModeConfigs(configNames ...string) error {
	c, err := config.New(config.Options{
		Provider: config.AWSAppConfig,
		Params: map[string]interface{}{
			constant.ConfigIDKey:       constant.Application,
			constant.ConfigRegionKey:   flags.AWSRegion(),
			constant.ConfigAccessKeyID: flags.AWSAccessKeyID(),
			constant.ConfigSecretKey:   flags.AWSSecretAccessKey(),
			constant.ConfigAppKey:      constant.Application,
			constant.ConfigEnvKey:      flags.Env(),
			constant.ConfigTypeKey:     "yaml",
			constant.ConfigNamesKey:    configNames,
		},
	})
	if err != nil {
		return err
	}
	client = getClient(c)
	return nil
}

func GetClient() *Client {
	return client
}

func getEnvironment() map[string]string {
	env := os.Environ()
	result := make(map[string]string)
	for _, e := range env {
		s := strings.Split(e, "=")
		if len(s) == 2 {
			result[s[0]] = s[1]
		}
	}
	return result
}

func getClient(c config.Client) *Client {
	return &Client{Client: c, env: getEnvironment()}
}

func (client *Client) GetStringWithEnv(config, key string) (string, error) {
	// Fetch the config value
	s, err := client.GetString(config, key)
	// If error no pointing moving ahead
	if err != nil {
		return s, err
	}

	// Fill value of config value with environment variable.
	for k, v := range client.env {
		s = strings.ReplaceAll(s, fmt.Sprintf("${%s}", k), v)
	}
	return s, nil
}

// For variables in appConfig in form of key: "${ENV_VAR}"
func (client *Client) GetStringWithEnvD(config, key, defaultValue string) string {
	// Fetch the config value
	s, err := client.GetString(config, key)
	if err != nil {
		return defaultValue
	}

	// Fill value of config value with environment variable.
	for k, v := range client.env {
		s = strings.ReplaceAll(s, fmt.Sprintf("${%s}", k), v)
	}
	return s
}
