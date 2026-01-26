package core

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AccessToken string
	AppVersion  string
	UserAgent   string
	SiteDomain  string
	DBpath      string
}

func Load() (*Config, error) {
	godotenv.Load(".env")

	cfg := &Config{
		AccessToken: os.Getenv("ACCESS_TOKEN"),
		AppVersion:  os.Getenv("APP_VERSION"),
		UserAgent:   os.Getenv("USER_AGENT"),
		SiteDomain:  os.Getenv("SITE_DOMAIN"),
		DBpath:      os.Getenv("DB_PATH"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.AccessToken == "" {
		return fmt.Errorf("ACCESS_TOKEN must be set") //TODO: Figure out how to generate these tokens
	}
	if c.AppVersion == "" {
		return fmt.Errorf("APP_VERSION must be set")
	}
	if c.UserAgent == "" {
		return fmt.Errorf("USER_AGENT must be set")
	}
	if c.SiteDomain == "" {
		return fmt.Errorf("SITE_DOMAIN must be set")
	}
	if c.DBpath == "" {
		return fmt.Errorf("DB_PATH must be set")
	}
	return nil
}

func (c *Config) WebSocketURL() string {
	return fmt.Sprintf(
		"wss://zond.api.2gis.ru/api/1.1/user/ws?appVersion=%s&channels=markers,sharing,routes&token=%s",
		c.AppVersion,
		c.AccessToken,
	)
}
