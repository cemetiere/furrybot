package main

import (
	"furrybot/balance"
	"furrybot/commands"
	"furrybot/config"
	"furrybot/femboy"
	"furrybot/images"
	"log"

	"github.com/NicoNex/echotron/v3"
)

func createBotFactory(token string) echotron.NewBotFn {
	return func(chatId int64) echotron.Bot {
		return &commands.Bot{
			API:             echotron.NewAPI(token),
			ChatId:          chatId,
			ImageRepository: &images.ReactorImageRepository{},
			FemboyGame:      femboy.NewFemboyGameService(),
			Balance:         balance.CreateNewBalanceService(),
		}
	}
}

func main() {
	err := config.ReadSettingsFromJson(config.GetSettingsPath())
	if err != nil {
		log.Fatalln("Failed to load configuration:", err)
	}
	log.Println("Settings loaded")

	botFactory := createBotFactory(config.Settings.TelegramBotToken)

	dsp := echotron.NewDispatcher(config.Settings.TelegramBotToken, botFactory)
	log.Fatalln(dsp.Poll())
}
