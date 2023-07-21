package config

import (
	"encoding/json"
	"os"
)

type SettingsModel struct {
	TelegramBotToken  string `json:"telegramBotToken"`
	PicsFolder        string `json:"picsFolder"`
	ReactorFolderName string `json:"reactorFolderName"`
}

const CONFIG_ENV = "FURRYBOT_CONFIG_FILE"
const DEFAULT_CONFIG = "defaultSettings.json"

var Settings = SettingsModel{}

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
func ReadSettingsFromJson(filePath string) error {
	f, _ := os.Open(filePath)

	Settings = SettingsModel{}

	decoder := json.NewDecoder(f)

	if err := decoder.Decode(&Settings); err != nil {
		return err
	}

	return nil
}
