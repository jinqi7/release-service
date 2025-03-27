package http

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	Token                string
	ServerUrl            string
	ProxyUrl             string
	RetryCount           int
	RetryWaitTime        time.Duration
	RetryMaxWaitTime     time.Duration
	Timeout              time.Duration
	LogLevel             string
	CollectorUrl         string
	Headers              map[string]string
	UseCollectorIfExists bool
}

func NewConfig() *Config {
	logLevel := os.Getenv("SEALIGHTS_LOG_LEVEL")
	return &Config{
		Token:                "",
		ServerUrl:            "",
		ProxyUrl:             "",
		RetryCount:           3,
		RetryWaitTime:        time.Second,
		RetryMaxWaitTime:     time.Second * 60,
		Timeout:              time.Second * 60,
		CollectorUrl:         "",
		LogLevel:             logLevel,
		Headers:              make(map[string]string),
		UseCollectorIfExists: false,
	}
}

func (c *Config) SetToken(token string) *Config {
	c.Token = token
	return c
}

func (c *Config) SetServerUrl(url string) *Config {
	c.ServerUrl = url
	return c
}

func (c *Config) SetProxyUrl(url string) *Config {
	c.ProxyUrl = url
	return c
}

func (c *Config) SetRetryCount(count int) *Config {
	c.RetryCount = count
	return c
}

func (c *Config) SetRetryWaitTime(waitTime time.Duration) *Config {
	c.RetryWaitTime = waitTime
	return c
}

func (c *Config) SetRetryMaxWaitTime(maxWaitTime time.Duration) *Config {
	c.RetryMaxWaitTime = maxWaitTime
	return c
}

func (c *Config) SetTimeout(timeout time.Duration) *Config {
	c.Timeout = timeout
	return c
}

func (c *Config) SetLogLevel(level string) *Config {
	c.LogLevel = level
	return c
}

func (c *Config) SetCollectorUrl(url string) *Config {
	c.CollectorUrl = url
	return c
}

func (c *Config) SetHeaders(headers map[string]string) *Config {
	c.Headers = headers
	return c
}

func (c *Config) SetUseCollectorIfExists(value bool) *Config {
	c.UseCollectorIfExists = value
	return c
}

func (c *Config) getBaseUrl(endpoint string) string {
	if c.UseCollectorIfExists && c.CollectorUrl != "" {
		if !strings.HasSuffix(c.CollectorUrl, "/api") {
			return c.CollectorUrl + "/api" + endpoint
		}
		return c.CollectorUrl + endpoint
	} else {
		return c.ServerUrl + endpoint
	}
}

func (c *Config) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("token is empty")
	}

	if c.ServerUrl == "" && c.CollectorUrl == "" {
		return fmt.Errorf("server url is empty")
	}
	return nil
}

func (c *Config) isDebugLevel() bool {
	return c.LogLevel == "debug"
}
