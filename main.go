package main

import (
	"furrybot/commands"
	"furrybot/config"
	"log"

	"github.com/NicoNex/echotron/v3"
)

func main() {
	err := config.ReadSettingsFromJson(config.GetSettingsPath())
	if err != nil {
		log.Fatalln("Failed to load configuration:", err)
	}
	log.Println("Settings loaded")

	botFactory := commands.CreateBotFactory(config.Settings.TelegramBotToken)

	dsp := echotron.NewDispatcher(config.Settings.TelegramBotToken, botFactory)
	log.Fatalln(dsp.Poll())
}
