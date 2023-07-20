package config

import (
	"encoding/json"
	"os"
)

type Settings struct {
	TelegramBotToken string `json:"telegramBotToken"`
	PicsFolder       string `json:"picsFolder"`
}

const CONFIG_ENV = "FURRYBOT_CONFIG_FILE"
const DEFAULT_CONFIG = "defaultSettings.json"

// Returns either a predefined file name or
// reads config file path from environment variable
// named "FURRYBOT_CONFIG_FILE"
func GetSettingsPath() string {
	if pathFromEnv := os.Getenv(CONFIG_ENV); pathFromEnv != "" {
		return pathFromEnv
	}

	return DEFAULT_CONFIG
}

// Reads and parses configs
func ReadSettingsFromJson(filePath string) (Settings, error) {
	f, _ := os.Open(filePath)

	settings := Settings{}

	decoder := json.NewDecoder(f)

	if err := decoder.Decode(&settings); err != nil {
		return settings, err
	}

	return settings, nil
}
