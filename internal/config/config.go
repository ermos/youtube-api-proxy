package config

import (
	"strings"

	"github.com/ermos/dotenv"
)

type Config struct {
	Port                   string
	YouTubeAPIKey          string
	YouTubeAllowedChannels []string
	CacheTTLSeconds        int
}

type rawConfig struct {
	Port                   string `env:"PORT" default:"8080"`
	YouTubeAPIKey          string `env:"YOUTUBE_API_KEY" required:"true"`
	YouTubeAllowedChannels string `env:"YOUTUBE_ALLOWED_CHANNELS" required:"true"`
	CacheTTLSeconds        int    `env:"CACHE_TTL_SECONDS" default:"86400"`
}

func Load() (*Config, error) {
	var raw rawConfig
	if err := dotenv.LoadStruct(&raw); err != nil {
		return nil, err
	}

	return &Config{
		Port:                   raw.Port,
		YouTubeAPIKey:          raw.YouTubeAPIKey,
		YouTubeAllowedChannels: parseChannels(raw.YouTubeAllowedChannels),
		CacheTTLSeconds:        raw.CacheTTLSeconds,
	}, nil
}

func parseChannels(channels string) []string {
	if channels == "" {
		return []string{}
	}

	parts := strings.Split(channels, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

func (c *Config) IsChannelAllowed(channelID string) bool {
	for _, allowed := range c.YouTubeAllowedChannels {
		if allowed == channelID {
			return true
		}
	}
	return false
}
