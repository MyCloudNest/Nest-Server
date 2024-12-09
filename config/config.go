package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type RateLimitConfig struct {
	LimitBody   int  `toml:"limit_body"`
	Enabled     bool `toml:"enabled"`
	MaxRequests int  `toml:"max_requests"`
	ExpireTime  int  `toml:"expire_time"`
}

type WhitelistConfig struct {
	WhitelistedIPs []string `toml:"whitelisted_ips"`
	Enabled        bool     `toml:"enabled"`
}

type PerformanceConfig struct {
	PerFork     bool `toml:"perfork"`
	Concurrency int  `toml:"concurrency"`
}

type CacheConfig struct {
	Enabled bool `toml:"enabled"`
}

type Config struct {
	Server      ServerConfig      `toml:"server"`
	Whitelist   WhitelistConfig   `toml:"whitelist"`
	RateLimit   RateLimitConfig   `toml:"rate_limit"`
	Performance PerformanceConfig `toml:"perfomance"`
	Cache       CacheConfig       `toml:"cache"`
}

const defaultConfigPath = "~/.cloudnest/config.toml"

func LoadConfig() (*Config, error) {
	configPath := defaultConfigPath
	if configPath[:2] == "~/" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, configPath[2:])
	}

	var cfg Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
	}
	return &cfg, nil
}
