package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chasinglogic/appdirs"
	"github.com/config-source/cdb/pkg/client"
	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseURL string

	Token string `json:"-" yaml:"-"`
}

var Client *client.Client
var Current Config

func ConfigFile() string {
	app := appdirs.New("cdb")
	configDir := app.UserConfig()
	configFile := filepath.Join(configDir, "config.yaml")
	return configFile
}

func DefaultConfig() Config {
	return Config{
		Token:   os.Getenv("CDB_TOKEN"),
		BaseURL: os.Getenv("CDB_BASE_URL"),
	}
}

func LoadConfig() error {
	cfgFile := ConfigFile()
	fh, err := os.Open(cfgFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to open config file %s: %w", cfgFile, err)
	} else if err != nil {
		cfgDir := appdirs.New("cdb").UserConfig()
		err := os.MkdirAll(cfgDir, 0700)
		if err != nil {
			return fmt.Errorf("failed to create config directories %s: %w", cfgDir, err)
		}

		Current = DefaultConfig()
	} else {
		err := yaml.NewDecoder(fh).Decode(&Current)
		if err != nil {
			return fmt.Errorf("failed to parse config file %s: %w", cfgFile, err)
		}
	}

	if os.Getenv("CDB_TOKEN") != "" {
		Current.Token = os.Getenv("CDB_TOKEN")
	}

	if Current.BaseURL == "" {
		return errors.New(
			"Unable to determine base URL for CDB instance." +
				"Try setting $CDB_BASE_URL or setting up a config file in " +
				cfgFile,
		)
	}

	Client = client.New(Current.Token, Current.BaseURL)
	return nil
}
