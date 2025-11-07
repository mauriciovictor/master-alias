package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Shell       string `json:"shell"`
	GithubToken string `json:"github_token"`
}

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".master-alias", "config.json")
}

func LoadConfig() (*Config, error) {
	path := GetConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	dir := filepath.Dir(GetConfigPath())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(GetConfigPath(), data, 0644)
}

func ConfigExists() bool {
	_, err := os.Stat(GetConfigPath())
	return err == nil
}

func PrintConfig() {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Println("⚠️  Nenhuma configuração encontrada.")
		return
	}
	fmt.Printf("💻 Terminal: %s\n", cfg.Shell)
}
