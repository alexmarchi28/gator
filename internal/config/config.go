package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

const defaultDBURL = "connection_string_goes_here"
const defaultCurrentUserName = "username_goes_here"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() Config {
	cfg, err := read()
	if err != nil {
		panic(err)
	}

	return cfg
}

func ReadOrCreate() (Config, bool, error) {
	cfg, err := read()
	if err == nil {
		return cfg, false, nil
	}

	if !os.IsNotExist(err) {
		return Config{}, false, err
	}

	cfg = Config{
		DBURL:           defaultDBURL,
		CurrentUserName: defaultCurrentUserName,
	}
	if err := write(cfg); err != nil {
		return Config{}, false, err
	}

	return cfg, true, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(*c)
}

func read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, configFileName), nil
}

func write(cfg Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFilePath, data, 0644)
}
