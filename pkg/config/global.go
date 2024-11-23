package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type KeycloakConfig struct {
	BaseURL      string `yaml:"baseURL"`
	Realm        string `yaml:"realm"`
	Id           string `yaml:"id"`
	ClientId     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
}

type ServerAddress struct {
	Host string `yaml:"host" json:"host"`
	Port string `yaml:"port" json:"port"`
}

type BaseConfig struct {
	Server              ServerAddress  `yaml:"server"`
	Keycloak            KeycloakConfig `yaml:"keycloak"`
	Mode                ServerMode     `yaml:"mode"`
	LogLevel            logrus.Level   `yaml:"logLevel"`
	OpenTelemtryAddress string         `yaml:"openTelemetryAddress"`
}

func (s *ServerAddress) Address() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}
