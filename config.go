package main

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port     int    `yaml:"port" mapstructure:"port"`
	Template string `yaml:"template" mapstructure:"template"`
	Webhook  struct {
		URL    string `yaml:"url" mapstructure:"url"`
		Secret string `yaml:"secret" mapstructure:"secret"`
	} `yaml:"webhook" mapstructure:"webhook"`
}

func (c *Config) Validate() error {
	if c.Port <= 0 {
		return errors.New("invalid port")
	}
	if c.Template == "" {
		return errors.New("template is required")
	}
	if c.Webhook.URL == "" {
		return errors.New("webhook.url is required")
	}
	return nil
}

func init() {
	viper.SetOptions(
		// NOTE: Bind dynamic struct fields from environment variables,
		// even without explicitly letting viper "know" that a key exists via viper.SetDefault() etc.
		// In the future, this feature flag might change: https://github.com/spf13/viper/issues/1851
		viper.ExperimentalBindStruct(),
	)
	// Automatically load from respective environment variables
	viper.AutomaticEnv()
	// Allow getting underscore-delimited environment variables via dot-delimited or hyphen-delimited key values
	// e.g. viper.Get("foo.bar") will lookup "FOO_BAR" environment variable so these can be mapped to structs
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Set defaults
	viper.SetDefault("port", 8080)
	viper.SetDefault("template", "no template set!")
	viper.SetDefault("webhook.url", "")
	viper.SetDefault("webhook.secret", "")
}

func NewConfig(location string) (*Config, error) {
	var c Config
	if location != "" {
		viper.SetConfigFile(location)
		err := viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}
	err := viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}
	err = c.Validate()
	if err != nil {
		return nil, err
	}
	return &c, nil
}
