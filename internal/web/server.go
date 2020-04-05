package web

import (
	"follis.net/internal/config"
	"follis.net/internal/database"
)


type WebConfig struct {
	DB database.CalorDB
}

func Init(config config.BoundConfig) WebConfig {
	return WebConfig{
		DB: config.Database,
	}
}