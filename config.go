package main

import (
	"os"
	"fmt"
	"errors"
	"strings"
)

type Config struct {
	LogLevel    int
	ListenHost  string
	ListenPort  int
	GitLabUrl   string
	GitLabToken string
	LabelPrefix string
	LabelColor  string
	IgnoreUser  string
}

func (config *Config) loadDefault() {
	config.LogLevel = LOG_DEBUG
	config.ListenHost = "localhost"
	config.ListenPort = 8081
	config.GitLabUrl = "https://gitlab.com/"
	config.LabelPrefix = "-"
	config.LabelColor = "#DDDDDD"
	config.IgnoreUser = "~"
}

func (config *Config) populate() error {
	config.loadDefault()

	config.LogLevel = getEnvInt("LOG_LEVEL", config.LogLevel)
	config.ListenHost = getEnvString("LISTEN_HOST", config.ListenHost)
	config.ListenPort = getEnvInt("LISTEN_PORT", config.ListenPort)
	config.GitLabUrl = strings.TrimRight(getEnvString("GITLAB_URL", config.GitLabUrl), "/")
	config.GitLabToken = getEnvString("GITLAB_TOKEN", config.GitLabToken)
	config.LabelPrefix = getEnvString("LABEL_PREFIX", config.LabelPrefix)
	config.LabelColor = getEnvString("LABEL_COLOR", config.LabelColor)
	config.IgnoreUser = getEnvString("IGNORE_USER", config.IgnoreUser)

	err := config.validate()

	return err
}

func (config *Config) validate() error {
	if config.LogLevel < 1 || config.LogLevel > 3 {
		return errors.New("invalid LOG_LEVEL. Should be a number between 1 and 3")
	}

	if config.ListenPort == 0 {
		return errors.New("invalid LISTEN_PORT. Should be a positive number")
	}

	if config.GitLabUrl == "" || (!strings.HasPrefix(config.GitLabUrl, "http://") && !strings.HasPrefix(config.GitLabUrl, "https://")) {
		return errors.New("invalid GITLAB_URL. Should be a URL starting with http:// or https://")
	}

	if config.GitLabToken == "" {
		return errors.New("empty GITLAB_TOKEN")
	}

	return nil
}

func getEnvString(name string, def string) string {
	envVal := os.Getenv(name)
	if envVal == "" {
		return def
	} else
	{
		return envVal
	}
}

func getEnvInt(name string, def int) int {
	envVal := os.Getenv(name)
	if envVal == "" {
		return def
	}

	var result int
	fmt.Sscan(envVal, &result)
	return result
}
