package clientconfigs

import (
	"errors"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/http/cookiejar"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/appconfigdata"
	jsoniter "github.com/json-iterator/go"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cast"
	"golang.org/x/net/context"
	"golang.org/x/net/publicsuffix"
	"gopkg.in/yaml.v3"
)

// credential modes
const (
	AppConfigStaticCredentialsMode = iota
	AppConfigSharedCredentialMode
)

const defaultAppConfigCheckInterval = time.Minute

type appConfigClient struct {
	options   appConfigClientOptions
	cfg       aws.Config
	client    *appconfigdata.Client
	parser    func([]byte, interface{}) error
	signal    chan struct{}
	configs   map[string]*appConfig
	mu        sync.RWMutex
	listeners map[string]ChangeListener
}

type appConfig struct {
	mu    sync.RWMutex
	token string
	data  map[string]interface{}
}

type appConfigClientOptions struct {
	id              string
	credentialsMode int
	region          string
	accessKeyID     string
	secretKey       string
	token           string
	app             string
	env             string
	configType      string
	configNames     []string
	checkInterval   time.Duration
}

type appConfigHTTPClientConfig struct {
	connectTimeout        time.Duration
	keepAliveDuration     time.Duration
	maxIdleConnections    int
	idleConnectionTimeout time.Duration
	tlsHandshakeTimeout   time.Duration
	expectContinueTimeout time.Duration
	timeout               time.Duration
}

func getAppConfigOption(options map[string]interface{}, key string) (string, error) {
	var val interface{}
	var ok bool
	var s string
	if val, ok = options[key]; ok {
		if s, ok = val.(string); !ok {
			return s, fmt.Errorf("invalid %s, must be a string", key)
		}
	} else {
		return s, fmt.Errorf("missing %s", key)
	}
	return s, nil
}

func getAppConfigClientOptions(options map[string]interface{}) (appConfigClientOptions, error) {
	var clientOptions appConfigClientOptions
	var err error
	clientOptions.id, err = getAppConfigOption(options, "id")
	if err != nil {
		return clientOptions, err
	}
	if cm, ok := options["credentialsMode"]; ok && cm == AppConfigSharedCredentialMode {
		clientOptions.credentialsMode = AppConfigSharedCredentialMode
	}
	clientOptions.region, err = getAppConfigOption(options, "region")
	if err != nil {
		return clientOptions, err
	}
	if clientOptions.credentialsMode == AppConfigStaticCredentialsMode {
		clientOptions.accessKeyID, err = getAppConfigOption(options, "accessKeyId")
		if err != nil {
			return clientOptions, err
		}
		clientOptions.secretKey, err = getAppConfigOption(options, "secretKey")
		if err != nil {
			return clientOptions, err
		}
	}
	clientOptions.app, err = getAppConfigOption(options, "app")
	if err != nil {
		return clientOptions, err
	}
	clientOptions.env, err = getAppConfigOption(options, "env")
	if err != nil {
		return clientOptions, err
	}
	clientOptions.configType, err = getAppConfigOption(options, "configType")
	if err != nil {
		return clientOptions, err
	}
	if val, ok := options["configNames"]; ok {
		if clientOptions.configNames, ok = val.([]string); !ok {
			return clientOptions, errors.New("invalid config names provided, should be an array of strings")
		}
	} else {
		return clientOptions, errors.New("missing configs")
	}
	if val, ok := options["checkInterval"]; ok {
		if clientOptions.checkInterval, ok = val.(time.Duration); !ok {
			return clientOptions, errors.New("invalid check interval provided, must be a time duration")
		}
	}
	if clientOptions.checkInterval <= (time.Second * 15) {
		clientOptions.checkInterval = defaultAppConfigCheckInterval
	}
	return clientOptions, nil
}

func getAppConfigHTTPClientConfig(options map[string]interface{}) appConfigHTTPClientConfig {
	// providing the defaults to the http client config
	cfg := appConfigHTTPClientConfig{
		connectTimeout:        time.Second * 10,
		keepAliveDuration:     time.Second * 30,
		maxIdleConnections:    100,
		idleConnectionTimeout: time.Second * 90,
		tlsHandshakeTimeout:   time.Second * 10,
		expectContinueTimeout: time.Second,
		timeout:               time.Second * 15,
	}
	// now checking for overrides
	if val, ok := options["connectTimeout"]; ok {
		if d, ok := val.(time.Duration); ok {
			cfg.connectTimeout = d
		}
	}
	if val, ok := options["keepAliveDuration"]; ok {
		if d, ok := val.(time.Duration); ok {
			cfg.keepAliveDuration = d
		}
	}
	if val, ok := options["maxIdleConnections"]; ok {
		if d, ok := val.(int); ok {
			cfg.maxIdleConnections = d
		}
	}
	if val, ok := options["idleConnectionTimeout"]; ok {
		if d, ok := val.(time.Duration); ok {
			cfg.idleConnectionTimeout = d
		}
	}
	if val, ok := options["tlsHandshakeTimeout"]; ok {
		if d, ok := val.(time.Duration); ok {
			cfg.tlsHandshakeTimeout = d
		}
	}
	if val, ok := options["expectContinueTimeout"]; ok {
		if d, ok := val.(time.Duration); ok {
			cfg.expectContinueTimeout = d
		}
	}
	if val, ok := options["timeout"]; ok {
		if d, ok := val.(time.Duration); ok {
			cfg.timeout = d
		}
	}
	return cfg
}

func (a *appConfigClient) watchConfig(ctx context.Context, name string, config *appConfig) {
	ticker := time.NewTicker(a.options.checkInterval)
	for {
		select {
		case <-ticker.C:
			// try to fetch the configurations again
			result, err := a.client.GetLatestConfiguration(ctx, &appconfigdata.GetLatestConfigurationInput{
				ConfigurationToken: aws.String(config.token)})
			if err != nil {
				// it might be possible that the configuration is deleted
				// stop the watch
				return
			}
			config.token = *result.NextPollConfigurationToken
			if len(result.Configuration) > 0 {
				// something has changed
				var data map[string]interface{}
				err = a.parser(result.Configuration, &data)
				if err != nil {
					// someone has added incorrect configurations
					// ignore the change for now then
				} else {
					config.mu.Lock()
					config.data = data
					config.mu.Unlock()
					// now we also need to notify the listener if any
					a.mu.RLock()
					if l, ok := a.listeners[name]; ok {
						l(config.data)
					}
					a.mu.RUnlock()
				}
			}
		case <-a.signal:
			// close the watch
			return
		}
	}
}

func (a *appConfigClient) fetchAndWatchConfigs(ctx context.Context) {
	app := aws.String(a.options.app)
	env := aws.String(a.options.env)
	for _, c := range a.options.configNames {
		s, err := a.client.StartConfigurationSession(ctx, &appconfigdata.StartConfigurationSessionInput{
			ApplicationIdentifier:                app,
			ConfigurationProfileIdentifier:       aws.String(c),
			EnvironmentIdentifier:                env,
			RequiredMinimumPollIntervalInSeconds: aws.Int32(int32(math.Floor(a.options.checkInterval.Seconds()))),
		})
		if err == nil {
			result, err := a.client.GetLatestConfiguration(ctx, &appconfigdata.GetLatestConfigurationInput{
				ConfigurationToken: s.InitialConfigurationToken})
			if err == nil {
				// now it means that the result exists
				// try to parse the result now
				var data map[string]interface{}
				err = a.parser(result.Configuration, &data)
				if err == nil {
					cfg := appConfig{
						token: *result.NextPollConfigurationToken,
						data:  data,
					}
					a.configs[c] = &cfg
					go a.watchConfig(ctx, c, &cfg)
				}
			}
		}
	}
}

func getParser(configType string) func([]byte, interface{}) error {
	switch configType {
	case jsonType:
		return jsoniter.Unmarshal
	case yamlType:
		return yaml.Unmarshal
	case tomlType:
		return toml.Unmarshal
	}
	return nil
}

func getAppConfigHTTPClient(options map[string]interface{}) *http.Client {
	if val, ok := options["httpClient"]; ok {
		if c, ok := val.(*http.Client); ok {
			return c
		}
	}
	c := getAppConfigHTTPClientConfig(options)
	cookieJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return http.DefaultClient
	}
	dialer := &net.Dialer{
		Timeout:   c.connectTimeout,
		KeepAlive: c.keepAliveDuration,
	}

	return &http.Client{
		Jar: cookieJar,
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          c.maxIdleConnections,
			IdleConnTimeout:       c.idleConnectionTimeout,
			TLSHandshakeTimeout:   c.tlsHandshakeTimeout,
			ExpectContinueTimeout: c.expectContinueTimeout,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
		// global timeout value for all requests
		Timeout: c.timeout,
	}
}

func getAppConfigSessionConfig(ctx context.Context, options map[string]interface{}, clientOptions appConfigClientOptions) (cfg aws.Config, err error) {
	if clientOptions.credentialsMode == AppConfigSharedCredentialMode {
		cfg, err = config.LoadDefaultConfig(ctx)
	} else {
		cred := credentials.NewStaticCredentialsProvider(clientOptions.accessKeyID,
			clientOptions.secretKey, "")
		cfg, err = config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(cred), config.WithRegion(clientOptions.region))
	}
	if err != nil {
		return
	}
	cfg.HTTPClient = getAppConfigHTTPClient(options)
	return
}

func newAppConfigClient(options map[string]interface{}) (*appConfigClient, error) {
	ctx := context.Background()
	clientOptions, err := getAppConfigClientOptions(options)
	if err != nil {
		return nil, err
	}
	cfg, err := getAppConfigSessionConfig(ctx, options, clientOptions)
	if err != nil {
		return nil, err
	}
	client := &appConfigClient{
		options:   clientOptions,
		cfg:       cfg,
		client:    appconfigdata.NewFromConfig(cfg),
		parser:    getParser(clientOptions.configType),
		signal:    make(chan struct{}, 1),
		configs:   make(map[string]*appConfig),
		listeners: make(map[string]ChangeListener),
	}
	client.fetchAndWatchConfigs(ctx)
	return client, nil
}

func (a *appConfigClient) AddChangeListener(config string, listener ChangeListener) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.configs[config]; !ok {
		return ErrConfigNotAdded
	}
	a.listeners[config] = listener
	return nil
}

func (a *appConfigClient) RemoveChangeListener(config string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.configs[config]; !ok {
		return ErrConfigNotAdded
	}
	delete(a.listeners, config)
	return nil
}

func get(kList []string, key string, val interface{}) interface{} {
	if len(kList) == 0 {
		if key == "" || key == "." {
			return val
		}
		return nil
	}
	// now need a map
	if m, ok := val.(map[string]interface{}); ok {
		newKey := key + kList[0]
		// for example a.b.c.d, then first try with a
		if v, ok := m[newKey]; ok {
			result := get(kList[1:], "", v)
			if result != nil {
				return result
			}
		}
		// otherwise, proceed with a.[something]
		return get(kList[1:], newKey+".", m)
	}
	// not a map
	return nil
}

func (a *appConfigClient) Get(config, key string) (interface{}, error) {
	if result, ok := a.configs[config]; ok {
		result.mu.RLock()
		defer result.mu.RUnlock()
		if key == "" {
			return result.data, nil
		}
		val := get(strings.Split(key, "."), "", result.data)
		if val == nil {
			return nil, ErrKeyNotFound
		}
		return val, nil
	}
	return nil, ErrConfigNotAdded
}

func (a *appConfigClient) GetD(config, key string, defaultValue interface{}) interface{} {
	val, err := a.Get(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (a *appConfigClient) GetInt(config, key string) (int64, error) {
	val, err := a.Get(config, key)
	if err != nil {
		return 0, err
	}
	return cast.ToInt64E(val)
}

func (a *appConfigClient) GetIntD(config, key string, defaultValue int64) int64 {
	val, err := a.GetInt(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (a *appConfigClient) GetFloat(config, key string) (float64, error) {
	val, err := a.Get(config, key)
	if err != nil {
		return 0, err
	}
	return cast.ToFloat64E(val)
}

func (a *appConfigClient) GetFloatD(config, key string, defaultValue float64) float64 {
	val, err := a.GetFloat(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (a *appConfigClient) GetString(config, key string) (string, error) {
	val, err := a.Get(config, key)
	if err != nil {
		return "", err
	}
	return cast.ToStringE(val)
}

func (a *appConfigClient) GetStringD(config, key string, defaultValue string) string {
	val, err := a.GetString(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (a *appConfigClient) GetBool(config, key string) (bool, error) {
	val, err := a.Get(config, key)
	if err != nil {
		return false, err
	}
	return cast.ToBoolE(val)
}

func (a *appConfigClient) GetBoolD(config, key string, defaultValue bool) bool {
	val, err := a.GetBool(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (a *appConfigClient) GetSlice(config, key string) ([]interface{}, error) {
	val, err := a.Get(config, key)
	if err != nil {
		return nil, err
	}
	return cast.ToSliceE(val)
}

func (a *appConfigClient) GetSliceD(config, key string, defaultValue []interface{}) []interface{} {
	val, err := a.GetSlice(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (a *appConfigClient) GetIntSlice(config, key string) ([]int64, error) {
	val, err := a.Get(config, key)
	if err != nil {
		return nil, err
	}
	return toInt64SliceE(val)
}

func (a *appConfigClient) GetIntSliceD(config, key string, defaultValue []int64) []int64 {
	val, err := a.GetIntSlice(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (a *appConfigClient) GetFloatSlice(config, key string) ([]float64, error) {
	val, err := a.Get(config, key)
	if err != nil {
		return nil, err
	}
	return toFloat64SliceE(val)
}

func (a *appConfigClient) GetFloatSliceD(config, key string, defaultValue []float64) []float64 {
	val, err := a.GetFloatSlice(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (a *appConfigClient) GetStringSlice(config, key string) ([]string, error) {
	val, err := a.Get(config, key)
	if err != nil {
		return nil, err
	}
	return cast.ToStringSliceE(val)
}

func (a *appConfigClient) GetStringSliceD(config, key string, defaultValue []string) []string {
	val, err := a.GetStringSlice(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (a *appConfigClient) GetBoolSlice(config, key string) ([]bool, error) {
	val, err := a.Get(config, key)
	if err != nil {
		return nil, err
	}
	return cast.ToBoolSliceE(val)
}

func (a *appConfigClient) GetBoolSliceD(config, key string, defaultValue []bool) []bool {
	val, err := a.GetBoolSlice(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (a *appConfigClient) GetMap(config, key string) (map[string]interface{}, error) {
	val, err := a.Get(config, key)
	if err != nil {
		return nil, err
	}
	return cast.ToStringMapE(val)
}

func (a *appConfigClient) GetMapD(config, key string, defaultValue map[string]interface{}) map[string]interface{} {
	val, err := a.GetMap(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (a *appConfigClient) GetIntMap(config, key string) (map[string]int64, error) {
	val, err := a.Get(config, key)
	if err != nil {
		return nil, err
	}
	return cast.ToStringMapInt64E(val)
}

func (a *appConfigClient) GetIntMapD(config, key string, defaultValue map[string]int64) map[string]int64 {
	val, err := a.GetIntMap(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (a *appConfigClient) GetFloatMap(config, key string) (map[string]float64, error) {
	val, err := a.Get(config, key)
	if err != nil {
		return nil, err
	}
	return toStringMapFloat64E(val)
}

func (a *appConfigClient) GetFloatMapD(config, key string, defaultValue map[string]float64) map[string]float64 {
	val, err := a.GetFloatMap(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (a *appConfigClient) GetStringMap(config, key string) (map[string]string, error) {
	val, err := a.Get(config, key)
	if err != nil {
		return nil, err
	}
	return cast.ToStringMapStringE(val)
}

func (a *appConfigClient) GetStringMapD(config, key string, defaultValue map[string]string) map[string]string {
	val, err := a.GetStringMap(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (a *appConfigClient) GetBoolMap(config, key string) (map[string]bool, error) {
	val, err := a.Get(config, key)
	if err != nil {
		return nil, err
	}
	return cast.ToStringMapBoolE(val)
}

func (a *appConfigClient) GetBoolMapD(config, key string, defaultValue map[string]bool) map[string]bool {
	val, err := a.GetBoolMap(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (a *appConfigClient) Unmarshal(config, key string, value interface{}) error {
	val, err := a.Get(config, key)
	if err != nil {
		return err
	}
	return unmarshal(val, value)
}

func (a *appConfigClient) Close() error {
	a.signal <- struct{}{}
	return nil
}
