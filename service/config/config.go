package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gookit/slog"
)

const (
	configSecret = "d52f2419b66d792e100448ffbd9a4b8893ccb87448339d7fc76ae58762b719f1" // totally insecure, sue me
)

type Config struct {
	TotalAmount   int    `json:"total_amount"`
	WalletConnect string `json:"wallet_connect"`
	Smtp          string `json:"smtp"`
	EmailTemplate string `json:"email_template"`
	DefaultTotal  *int   `json:"default_total"`
	DefaultEach   *int   `json:"default_each"`
}

func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".zap")
	configPath := filepath.Join(configDir, "config.json")
	return configPath
}

func SetConfig(config Config) error {
	configPath := getConfigPath()
	if config.WalletConnect != "" {
		config.WalletConnect = encrypt(configSecret, config.WalletConnect)
	}
	if config.Smtp != "" {
		config.Smtp = encrypt(configSecret, config.Smtp)
	}
	configData, _ := json.MarshalIndent(config, "", "  ")
	slog.Debug("Writing config", configData)

	os.MkdirAll(filepath.Dir(configPath), 0700)
	return ioutil.WriteFile(configPath, configData, 0600)
}

func LoadConfig() (*Config, error) {
	slog.Debug("Loading config")
	data, err := ioutil.ReadFile(getConfigPath())
	if err != nil {
		slog.Warn("Failed fetching config", err)
		return &Config{}, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.WalletConnect != "" {
		config.WalletConnect = decrypt(configSecret, config.WalletConnect)
	}
	if config.Smtp != "" {
		config.Smtp = decrypt(configSecret, config.Smtp)
	}
	if config.DefaultEach == nil {
		amount := 100
		config.DefaultEach = &amount
	}
	slog.Debug("Fetched config", config)

	return &config, nil
}
