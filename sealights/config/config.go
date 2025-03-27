package config

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/konflux-ci/release-service/__sealights__/pkg/common"
)

var ConfigFile = `package __sealights__

var ConfigData = %q
`

func GetConfigFile(data []byte) string {
	sEnc := b64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf(ConfigFile, sEnc)
}

var DefaultConfig = &Config{
	InstrumentMap:                  make(map[string]string),
	CustomerId:                     "",
	Technology:                     "",
	AppName:                        "",
	ModuleName:                     "",
	Build:                          "",
	Branch:                         "",
	Generated:                      0,
	AgentId:                        "",
	BuildSessionId:                 "",
	Token:                          "",
	ServerUrl:                      "",
	ProxyUrl:                       "",
	LabId:                          "",
	TestStage:                      "",
	LogLevel:                       "",
	AgentVersion:                   "",
	LightMode:                      false,
	CollectorUrl:                   "",
	NoTests:                        false,
	EventBufferSize:                1000,
	FootprintsBufferSize:           10000,
	ConnectionRetryCount:           3,
	ConnectionRetryWaitTime:        5 * time.Second,
	ConnectionRetryMaxWaitTime:     30 * time.Second,
	ConnectionTimeout:              30 * time.Second,
	DisableOnNoConnection:          false,
	EnableRemoteConfig:             false,
	RemoteConfigIntervalSec:        300,
	TestSelection:                  true,
	TestRecommendationSleepSeconds: 60,
	TestsRunnerMode:                false,
	GinkgoEnabled:                  false,
}

type Config struct {
	InstrumentMap                  map[string]string `json:"instrumentMap"`
	CustomerId                     string            `json:"customerId"`
	Technology                     string            `json:"technology"`
	AppName                        string            `json:"appName"`
	ModuleName                     string            `json:"moduleName"`
	Build                          string            `json:"build"`
	Branch                         string            `json:"branch"`
	Generated                      int64             `json:"generated"`
	AgentId                        string            `json:"agentId"`
	BuildSessionId                 string            `json:"buildSessionId"`
	Token                          string            `json:"token"`
	ServerUrl                      string            `json:"serverUrl"`
	ProxyUrl                       string            `json:"proxyUrl"`
	LabId                          string            `json:"labId"`
	TestStage                      string            `json:"testStage"`
	LogLevel                       string            `json:"logLevel"`
	AgentVersion                   string            `json:"agentVersion"`
	LightMode                      bool              `json:"lightMode"`
	CollectorUrl                   string            `json:"collectorUrl"`
	NoTests                        bool              `json:"noTests"`
	EventBufferSize                int64             `json:"eventBufferSize"`
	FootprintsBufferSize           int64             `json:"footprintsBufferSize"`
	ConnectionRetryCount           int64             `json:"connectionRetryCount"`
	ConnectionRetryWaitTime        time.Duration     `json:"connectionRetryWaitTime"`
	ConnectionRetryMaxWaitTime     time.Duration     `json:"connectionRetryMaxWaitTime"`
	ConnectionTimeout              time.Duration     `json:"connectionTimeout"`
	DisableOnNoConnection          bool              `json:"disableOnNoConnection"`
	EnableRemoteConfig             bool              `json:"enableRemoteConfig"`
	RemoteConfigIntervalSec        int64             `json:"remoteConfigIntervalSec"`
	TestSelection                  bool              `json:"testSelection"`
	DisableSealightsOnInit         bool              `json:"disableSealightsOnInit"`
	TestGroupId                    string            `json:"testGroupId"`
	TestRecommendationSleepSeconds int64             `json:"testRecommendationSleepSeconds"`
	TestsRunnerMode                bool              `json:"testsRunnerMode"`
	GinkgoEnabled                  bool              `json:"ginkgoEnabled"`
}

func NewConfig() *Config {
	return DefaultConfig
}

func (c *Config) Add(key string, value string) {
	c.InstrumentMap[key] = value
}

func (c *Config) Get(key string) (string, bool) {
	val, found := c.InstrumentMap[key]
	return val, found
}

func (c *Config) SetCustomerId(value string) *Config {
	c.CustomerId = value
	return c
}

func (c *Config) SetTechnology(value string) *Config {
	c.Technology = value
	return c
}

func (c *Config) SetAppName(value string) *Config {
	c.AppName = value
	return c
}

func (c *Config) SetBuild(value string) *Config {
	c.Build = value
	return c
}

func (c *Config) SetBranch(value string) *Config {
	c.Branch = value
	return c
}

func (c *Config) SetGenerated(value int64) *Config {
	c.Generated = value
	return c
}

func (c *Config) SetAgentId(value string) *Config {
	c.AgentId = value
	return c
}

func (c *Config) SetBuildSessionId(value string) *Config {
	c.BuildSessionId = value
	return c
}

func (c *Config) SetToken(value string) *Config {
	c.Token = value
	return c
}

func (c *Config) SetServerUrl(value string) *Config {
	c.ServerUrl = value
	return c
}

func (c *Config) SetProxyUrl(value string) *Config {
	c.ProxyUrl = value
	return c
}

func (c *Config) SetLabId(value string) *Config {
	c.LabId = value
	return c
}

func (c *Config) SetTestStage(value string) *Config {
	c.TestStage = value
	return c
}

func (c *Config) SetLogLevel(value string) *Config {
	c.LogLevel = value
	return c
}

func (c *Config) SetAgentVersion(value string) *Config {
	c.AgentVersion = value
	return c
}

func (c *Config) SetLightMode(value bool) *Config {
	c.LightMode = value
	return c
}

func (c *Config) SetCollectorUrl(value string) *Config {
	c.CollectorUrl = value
	return c
}

func (c *Config) SetModuleName(value string) *Config {
	c.ModuleName = value
	return c
}

func (c *Config) SetNoTests(value bool) *Config {
	c.NoTests = value
	return c
}

func (c *Config) SetEventBufferSize(value int64) *Config {
	c.EventBufferSize = value
	return c
}

func (c *Config) SetFootprintsBufferSize(value int64) *Config {
	c.FootprintsBufferSize = value
	return c
}

func (c *Config) SetConnectionRetryCount(value int64) *Config {
	c.ConnectionRetryCount = value
	return c
}

func (c *Config) SetConnectionRetryWaitTime(value time.Duration) *Config {
	c.ConnectionRetryWaitTime = value
	return c
}

func (c *Config) SetConnectionTimeout(value time.Duration) *Config {
	c.ConnectionTimeout = value
	return c
}

func (c *Config) SetDisableOnNoConnection(value bool) *Config {
	c.DisableOnNoConnection = value
	return c
}

func (c *Config) SetEnableRemoteConfig(value bool) *Config {
	c.EnableRemoteConfig = value
	return c
}

func (c *Config) SetRemoteConfigIntervalSec(value int64) *Config {
	c.RemoteConfigIntervalSec = value
	return c
}

func (c *Config) SetTestSelection(value bool) *Config {
	c.TestSelection = value
	return c
}
func (c *Config) SetDisableSealightsOnInit(value bool) *Config {
	c.DisableSealightsOnInit = value
	return c
}

func (c *Config) SetTestGroupId(value string) *Config {
	c.TestGroupId = value
	return c
}

func (c *Config) SetTestRecommendationSleepSeconds(value int64) *Config {
	c.TestRecommendationSleepSeconds = value
	return c
}

func (c *Config) SetTestsRunnerMode(value bool) *Config {
	c.TestsRunnerMode = value
	return c
}

func (c *Config) SetGinkgoEnabled(value bool) *Config {
	c.GinkgoEnabled = value
	return c
}

func (c *Config) Save(path string) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	if err := enc.Encode(c); err != nil {
		return fmt.Errorf("failed to encode config: %s", err.Error())
	}

	err := ioutil.WriteFile(path+".gob", buffer.Bytes(), 0o600)
	if err != nil {
		return fmt.Errorf("failed to write file: %s with error: %s", path+".gob", err.Error())
	}
	return nil
}

func (c *Config) SaveToGoFile(path string) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	if err := enc.Encode(c); err != nil {
		return fmt.Errorf("failed to encode config: %s", err.Error())
	}
	cfgFileStr := GetConfigFile(buffer.Bytes())
	err := ioutil.WriteFile(path+".go", []byte(cfgFileStr), 0o600)
	if err != nil {
		return fmt.Errorf("failed to write file: %s with error: %s", path+".gob", err.Error())
	}
	return nil
}

func (c *Config) validate() error {
	if c.Token == "" {
		return fmt.Errorf("token is empty")
	}
	if c.ServerUrl == "" {
		return fmt.Errorf("server url is empty")
	}
	if c.AgentId == "" {
		return fmt.Errorf("agent id is empty")
	}
	if c.BuildSessionId == "" {
		return fmt.Errorf("build session id is empty")
	}

	if c.AppName == "" {
		return fmt.Errorf("app name is empty")
	}
	if c.Build == "" {
		return fmt.Errorf("build is empty")
	}
	if c.Branch == "" {
		return fmt.Errorf("branch is empty")
	}

	return nil
}

func (c *Config) Load(strData string) error {
	if err := c.loadFromString(strData); err != nil {
		return err
	}
	c.parseEnvString(&c.Token, "SEALIGHTS_TOKEN", "").
		parseEnvFile(&c.Token, "SEALIGHTS_TOKEN_FILE").
		parseEnvString(&c.ServerUrl, "SEALIGHTS_SERVER_URL", "").
		parseEnvString(&c.AgentId, "SEALIGHTS_AGENT_ID", common.NewUUID()).
		parseEnvString(&c.BuildSessionId, "SEALIGHTS_BUILD_SESSION_ID", common.NewUUID()).
		parseEnvFile(&c.BuildSessionId, "SEALIGHTS_BUILD_SESSION_ID_FILE").
		parseEnvString(&c.AppName, "SEALIGHTS_APP_NAME", "").
		parseEnvString(&c.Build, "SEALIGHTS_BUILD", "").
		parseEnvString(&c.Branch, "SEALIGHTS_BRANCH", "").
		parseEnvString(&c.AgentVersion, "SEALIGHTS_AGENT_VERSION", "").
		parseEnvString(&c.CollectorUrl, "SEALIGHTS_COLLECTOR_URL", "").
		parseEnvString(&c.ModuleName, "SEALIGHTS_MODULE_NAME", "").
		parseEnvBool(&c.NoTests, "SEALIGHTS_NO_TESTS").
		parseEnvInt(&c.EventBufferSize, "SEALIGHTS_EVENT_BUFFER_SIZE", 1000).
		parseEnvInt(&c.FootprintsBufferSize, "SEALIGHTS_FOOTPRINTS_BUFFER_SIZE", 10000).
		parseEnvString(&c.LabId, "SEALIGHTS_LAB_ID", c.BuildSessionId).
		parseEnvString(&c.Technology, "SEALIGHTS_TECHNOLOGY", "").
		parseEnvString(&c.ProxyUrl, "SEALIGHTS_PROXY_URL", "").
		parseEnvString(&c.TestStage, "SEALIGHTS_TEST_STAGE", "").
		parseEnvString(&c.LogLevel, "SEALIGHTS_LOG_LEVEL", "").
		parseEnvBool(&c.LightMode, "SEALIGHTS_LIGHT_MODE").
		parseEnvInt(&c.ConnectionRetryCount, "SEALIGHTS_CONNECTION_RETRY_COUNT", 3).
		parseEnvDuration(&c.ConnectionRetryWaitTime, "SEALIGHTS_CONNECTION_RETRY_WAIT_TIME", time.Second*5).
		parseEnvDuration(&c.ConnectionRetryMaxWaitTime, "SEALIGHTS_CONNECTION_RETRY_MAX_WAIT_TIME", time.Second*30).
		parseEnvDuration(&c.ConnectionTimeout, "SEALIGHTS_CONNECTION_TIMEOUT", time.Second*30).
		parseEnvBool(&c.DisableOnNoConnection, "SEALIGHTS_DISABLE_ON_NO_CONNECTION").
		parseEnvBool(&c.EnableRemoteConfig, "SEALIGHTS_ENABLE_REMOTE_CONFIG").
		parseEnvBool(&c.TestSelection, "SEALIGHTS_TEST_SELECTION").
		parseEnvInt(&c.RemoteConfigIntervalSec, "SEALIGHTS_REMOTE_CONFIG_INTERVAL_SEC", 300).
		parseEnvString(&c.TestGroupId, "SEALIGHTS_TEST_GROUP_ID", "").
		parseEnvInt(&c.TestRecommendationSleepSeconds, "SEALIGHTS_TEST_RECOMMENDATION_SLEEP_SECONDS", 60).
		parseEnvBool(&c.TestsRunnerMode, "SEALIGHTS_FUNCTIONAL_TESTS_MODE").
		parseEnvBool(&c.GinkgoEnabled, "SEALIGHTS_ENABLE_GINKGO")

	if err := c.validate(); err != nil {
		return err
	}
	return nil
}

func (c *Config) IsDebug() bool {
	return c.LogLevel == "debug"
}

func (c *Config) load(path string) error {
	data, err := ioutil.ReadFile(path + ".gob")
	if err != nil {
		return fmt.Errorf("failed to read file: %s with error: %s", path+".gob", err.Error())
	}
	buffer := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buffer)
	if err := dec.Decode(c); err != nil {
		return fmt.Errorf("failed to decode config: %s", err.Error())
	}
	return nil
}

func (c *Config) loadFromString(strData string) error {
	sDec, err := b64.StdEncoding.DecodeString(strData)
	if err != nil {
		return fmt.Errorf("failed to decode config: %s", err.Error())
	}
	buffer := bytes.NewBuffer(sDec)
	dec := gob.NewDecoder(buffer)
	if err := dec.Decode(c); err != nil {
		return fmt.Errorf("failed to decode config: %s", err.Error())
	}
	return nil
}

func (c *Config) loadStringFromFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %s with error: %s", path, err.Error())
	}
	return string(data), nil
}

func (c *Config) parseEnvString(field *string, envVar string, defaultVal string) *Config {
	envVal := os.Getenv(envVar)
	if envVal != "" {
		*field = envVal
	} else {
		if *field == "" && defaultVal != "" {
			*field = defaultVal
		}
	}
	return c
}

func (c *Config) parseEnvInt(field *int64, envVar string, defaultVal int64) *Config {
	envVal := os.Getenv(envVar)
	if envVal != "" {
		nVal, err := strconv.ParseInt(envVal, 10, 64)
		if err == nil {
			*field = nVal
		}

	} else {
		if *field == 0 && defaultVal != 0 {
			*field = defaultVal
		}
	}
	return c
}

func (c *Config) parseEnvBool(field *bool, envVar string) *Config {
	envVal, exist := os.LookupEnv(envVar)
	if !exist {
		return c
	}
	if envVal != "" {
		if envVal == "true" {
			*field = true
		} else {
			*field = false
		}
	}
	return c
}

func (c *Config) parseEnvFile(field *string, envVar string) *Config {
	envVal := os.Getenv(envVar)
	if envVal != "" {
		data, err := c.loadStringFromFile(envVal)
		if err != nil {
			return c
		}
		if data != "" {
			*field = data
		}
	}
	return c
}

func (c *Config) parseEnvDuration(t *time.Duration, envVar string, defaultVal time.Duration) *Config {
	envVal := os.Getenv(envVar)
	if envVal != "" {
		d, err := time.ParseDuration(envVal)
		if err == nil {
			*t = d
		}
	} else {
		if *t == 0 && defaultVal != 0 {
			*t = defaultVal
		}
	}
	return c
}

func (c *Config) GetAllSealightsEnv() string {
	var env []string
	envVar := os.Environ()
	for _, v := range envVar {
		if strings.HasPrefix(v, "SEALIGHTS_") {
			env = append(env, v)
		}
	}
	if len(env) == 0 {
		return "no sealights environment variables were found"
	}

	return strings.Join(env, ", ")
}

func (c *Config) String() string {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("error marshalling config: %s", err.Error())
	}
	return string(data)
}
