// Package config 承载网关的运行配置。
package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Upstream 是一个上游模型服务。
type Upstream struct {
	Name     string `yaml:"name"`
	BaseURL  string `yaml:"base_url"`
	Protocol string `yaml:"protocol"` // openai | anthropic
	APIKey   string `yaml:"api_key"`
	Weight   int    `yaml:"weight"`
}

// Config 是网关配置。
type Config struct {
	Listen    string     `yaml:"listen"`
	Upstreams []Upstream `yaml:"upstreams"`

	// 限流：先放行、后付费（见 docs/decisions.md ADR-001）
	RateLimit RateLimitConfig `yaml:"rate_limit"`
}

// RateLimitConfig 是 token 级限流的配置。
type RateLimitConfig struct {
	Enabled bool          `yaml:"enabled"`
	Window  time.Duration `yaml:"window"`
	Quota   int           `yaml:"quota"` // 每窗口 token 配额
	Store   string        `yaml:"store"` // memory | redis
}

// Default 返回可直接运行的默认配置。
func Default() *Config {
	return &Config{
		Listen: ":8080",
		RateLimit: RateLimitConfig{
			Enabled: false,
			Window:  time.Minute,
			Quota:   100000,
			Store:   "memory",
		},
	}
}

// Load 读取配置文件；path 为空时返回默认配置。
func Load(path string) (*Config, error) {
	if path == "" {
		return Default(), nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := Default()
	if err := yaml.Unmarshal(b, cfg); err != nil {
		return nil, err
	}
	if cfg.Listen == "" {
		return nil, errors.New("listen 不能为空")
	}
	return cfg, nil
}
