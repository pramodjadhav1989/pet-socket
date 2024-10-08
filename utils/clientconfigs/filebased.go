package clientconfigs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"

	"github.com/spf13/cast"
)

type lockedKoanf struct {
	*koanf.Koanf
	mu *sync.RWMutex
}

type fileBasedClient struct {
	options   fileBasedClientOptions
	configs   map[string]*lockedKoanf
	listeners map[string]ChangeListener
	mu        sync.RWMutex
}

type fileBasedClientOptions struct {
	path        string
	configNames []string
	configType  string
}

func getFileBasedClientOptions(options map[string]interface{}) (fileBasedClientOptions, error) {
	clientOptions := fileBasedClientOptions{}
	var val interface{}
	var ok bool
	if val, ok = options["configsDirectory"]; ok {
		if clientOptions.path, ok = val.(string); ok {
			clientOptions.path = filepath.Clean(clientOptions.path)
		} else {
			return clientOptions, errors.New("invalid config directory provided")
		}
	} else {
		return clientOptions, errors.New("config directory not provided")
	}
	if val, ok = options["configNames"]; ok {
		if clientOptions.configNames, ok = val.([]string); !ok {
			return clientOptions, errors.New("invalid config names provided, should be an array of strings")
		}

	} else {
		return clientOptions, errors.New("no configs provided to be used")
	}
	if val, ok = options["configType"]; ok {
		if clientOptions.configType, ok = val.(string); !ok || (clientOptions.configType != jsonType &&
			clientOptions.configType != yamlType &&
			clientOptions.configType != tomlType) {
			return clientOptions, fmt.Errorf("invalid config type provided should be one of %s, %s or %s",
				jsonType, yamlType, tomlType)
		}
	} else {
		return clientOptions, errors.New("no config type provided")
	}
	return clientOptions, nil
}

func (f *fileBasedClient) onConfigChange(name string) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if l, ok := f.listeners[name]; ok {
		l(name)
	}
}

func (f *fileBasedClient) getPath(options fileBasedClientOptions, name string) (string, error) {
	switch options.configType {
	case jsonType:
		path := filepath.Join(options.path, fmt.Sprintf("%s.%s", name, options.configType))
		_, err := ioutil.ReadFile(path)
		return path, err
	case yamlType:
		path := filepath.Join(options.path, fmt.Sprintf("%s.%s", name, options.configType))
		_, err := ioutil.ReadFile(path)
		if err == nil {
			return path, nil
		}
		path = filepath.Join(options.path, fmt.Sprintf("%s.%s", name, "yml"))
		_, err = ioutil.ReadFile(path)
		return path, err
	case tomlType:
		path := filepath.Join(options.path, fmt.Sprintf("%s.%s", name, options.configType))
		_, err := ioutil.ReadFile(path)
		return path, err
	default:
		return "", errors.New("invalid config type")
	}
}

func (f *fileBasedClient) getConfig(options fileBasedClientOptions, name string) (*lockedKoanf, error) {
	k := koanf.New(".")
	var p koanf.Parser
	switch options.configType {
	case jsonType:
		p = json.Parser()
	case yamlType:
		p = yaml.Parser()
	case tomlType:
		p = toml.Parser()
	}
	path, err := f.getPath(options, name)
	if err != nil {
		return nil, err
	}
	fp := file.Provider(path)
	err = k.Load(fp, p)
	if err != nil {
		return nil, err
	}
	mu := sync.RWMutex{}
	err = fp.Watch(func(_ interface{}, err error) {
		if err != nil {
			return
		}
		mu.Lock()
		_ = k.Load(fp, p)
		mu.Unlock()
		f.onConfigChange(name)
	})
	if err != nil {
		return nil, err
	}
	return &lockedKoanf{Koanf: k, mu: &mu}, nil
}

func newFileBasedClient(options map[string]interface{}) (*fileBasedClient, error) {
	clientOptions, err := getFileBasedClientOptions(options)
	if err != nil {
		return nil, err
	}
	client := &fileBasedClient{
		options: clientOptions,
	}
	client.configs = make(map[string]*lockedKoanf)
	for _, name := range clientOptions.configNames {
		v, err := client.getConfig(clientOptions, name)
		if err != nil {
			return nil, err
		}
		client.configs[name] = v
	}
	client.listeners = make(map[string]ChangeListener)
	return client, nil
}

func (f *fileBasedClient) AddChangeListener(config string, listener ChangeListener) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.listeners[config] = listener
	return nil
}

func (f *fileBasedClient) RemoveChangeListener(config string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.listeners, config)
	return nil
}

func (f *fileBasedClient) Get(config, key string) (interface{}, error) {
	if k, ok := f.configs[config]; ok {
		k.mu.RLock()
		defer k.mu.RUnlock()
		d := k.Get(key)
		if d == nil {
			return nil, ErrKeyNotFound
		}
		return d, nil
	} else {
		return nil, ErrConfigNotAdded
	}
}

func (f *fileBasedClient) GetD(config, key string, defaultValue interface{}) interface{} {
	val, err := f.Get(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (f *fileBasedClient) GetInt(config, key string) (int64, error) {
	val, err := f.Get(config, key)
	if err != nil {
		return 0, err
	}
	return cast.ToInt64E(val)
}

func (f *fileBasedClient) GetIntD(config, key string, defaultValue int64) int64 {
	val, err := f.GetInt(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (f *fileBasedClient) GetFloat(config, key string) (float64, error) {
	val, err := f.Get(config, key)
	if err != nil {
		return 0, err
	}
	return cast.ToFloat64E(val)
}

func (f *fileBasedClient) GetFloatD(config, key string, defaultValue float64) float64 {
	val, err := f.GetFloat(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (f *fileBasedClient) GetString(config, key string) (string, error) {
	val, err := f.Get(config, key)
	if err != nil {
		return "", err
	}
	return cast.ToStringE(val)
}

func (f *fileBasedClient) GetStringD(config, key string, defaultValue string) string {
	val, err := f.GetString(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (f *fileBasedClient) GetBool(config, key string) (bool, error) {
	val, err := f.Get(config, key)
	if err != nil {
		return false, err
	}
	return cast.ToBoolE(val)
}

func (f *fileBasedClient) GetBoolD(config, key string, defaultValue bool) bool {
	val, err := f.GetBool(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (f *fileBasedClient) GetSlice(config, key string) ([]interface{}, error) {
	val, err := f.Get(config, key)
	if err != nil {
		return nil, err
	}
	return cast.ToSliceE(val)
}

func (f *fileBasedClient) GetSliceD(config, key string, defaultValue []interface{}) []interface{} {
	val, err := f.GetSlice(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (f *fileBasedClient) GetIntSlice(config, key string) ([]int64, error) {
	val, err := f.Get(config, key)
	if err != nil {
		return nil, err
	}
	return toInt64SliceE(val)
}

func (f *fileBasedClient) GetIntSliceD(config, key string, defaultValue []int64) []int64 {
	val, err := f.GetIntSlice(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (f *fileBasedClient) GetFloatSlice(config, key string) ([]float64, error) {
	val, err := f.Get(config, key)
	if err != nil {
		return nil, err
	}
	return toFloat64SliceE(val)
}

func (f *fileBasedClient) GetFloatSliceD(config, key string, defaultValue []float64) []float64 {
	val, err := f.GetFloatSlice(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (f *fileBasedClient) GetStringSlice(config, key string) ([]string, error) {
	val, err := f.Get(config, key)
	if err != nil {
		return nil, err
	}
	return cast.ToStringSliceE(val)
}

func (f *fileBasedClient) GetStringSliceD(config, key string, defaultValue []string) []string {
	val, err := f.GetStringSlice(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (f *fileBasedClient) GetBoolSlice(config, key string) ([]bool, error) {
	val, err := f.Get(config, key)
	if err != nil {
		return nil, err
	}
	return cast.ToBoolSliceE(val)
}

func (f *fileBasedClient) GetBoolSliceD(config, key string, defaultValue []bool) []bool {
	val, err := f.GetBoolSlice(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (f *fileBasedClient) GetMap(config, key string) (map[string]interface{}, error) {
	val, err := f.Get(config, key)
	if err != nil {
		return nil, err
	}
	return cast.ToStringMapE(val)
}

func (f *fileBasedClient) GetMapD(config, key string, defaultValue map[string]interface{}) map[string]interface{} {
	val, err := f.GetMap(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (f *fileBasedClient) GetIntMap(config, key string) (map[string]int64, error) {
	val, err := f.Get(config, key)
	if err != nil {
		return nil, err
	}
	return cast.ToStringMapInt64E(val)
}

func (f *fileBasedClient) GetIntMapD(config, key string, defaultValue map[string]int64) map[string]int64 {
	val, err := f.GetIntMap(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (f *fileBasedClient) GetFloatMap(config, key string) (map[string]float64, error) {
	val, err := f.Get(config, key)
	if err != nil {
		return nil, err
	}
	return toStringMapFloat64E(val)
}

func (f *fileBasedClient) GetFloatMapD(config, key string, defaultValue map[string]float64) map[string]float64 {
	val, err := f.GetFloatMap(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (f *fileBasedClient) GetStringMap(config, key string) (map[string]string, error) {
	val, err := f.Get(config, key)
	if err != nil {
		return nil, err
	}
	return cast.ToStringMapStringE(val)
}

func (f *fileBasedClient) GetStringMapD(config, key string, defaultValue map[string]string) map[string]string {
	val, err := f.GetStringMap(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (f *fileBasedClient) GetBoolMap(config, key string) (map[string]bool, error) {
	val, err := f.Get(config, key)
	if err != nil {
		return nil, err
	}
	return cast.ToStringMapBoolE(val)
}

func (f *fileBasedClient) GetBoolMapD(config, key string, defaultValue map[string]bool) map[string]bool {
	val, err := f.GetBoolMap(config, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (f *fileBasedClient) Unmarshal(config, key string, value interface{}) error {
	val, err := f.Get(config, key)
	if err != nil {
		return err
	}
	return unmarshal(val, value)
}

func (f *fileBasedClient) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for k := range f.configs {
		delete(f.configs, k)
	}
	for k := range f.listeners {
		delete(f.listeners, k)
	}
	return nil
}
