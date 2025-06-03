package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DbUrl    string `json:"db_url"`
	UserName string `json:"current_user_name"`
}

const configFilename = ".gatorconfig.json"

func Read() (Config, error) {
	filename, err := getConfigDir()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c *Config) SetUser(username string) error {
	c.UserName = username
	return write(*c)
}

func write(c Config) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	filename, err := getConfigDir()
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, data, 0644)
	return err
}

func getConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(home, configFilename)
	return fullPath, nil
}
