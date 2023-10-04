package commands

import (
	"furrybot/balance"
	"furrybot/femboy"
	"furrybot/images"
	"log"

	"github.com/NicoNex/echotron/v3"
)

var commandsList = []Command{
	GetFurryListCommand,
	GetFurryPicCommand,
	ShowRepositorySelectionCommand,
	SelectRepositoryCommand,
	OlegShipulinCommand,
	FemboyRegisterCommand,
	ChooseTodaysFemboyCommand,
	ShowLeaderboardCommand,
	ShowBalanceCommand,
	FuckCommand,
}

type Bot struct {
	echotron.API
	ChatId          int64
	ImageRepository images.IImageRepository
	FemboyGame      *femboy.FemboyGameService
	Balance         *balance.BalanceService
	Username2UserId map[string]int64
}

func CreateBotFactory(token string) echotron.NewBotFn {
	return func(chatId int64) echotron.Bot {
		return &Bot{
			API:             echotron.NewAPI(token),
			ChatId:          chatId,
			ImageRepository: &images.ReactorImageRepository{},
			FemboyGame:      femboy.NewFemboyGameService(),
			Balance:         balance.CreateNewBalanceService(),
			Username2UserId: make(map[string]int64),
		}
	}
}

func (bot *Bot) Update(update *echotron.Update) {
	if update.Message != nil {
		log.Printf("[%s] %s", update.Message.From.Username, update.Message.Text)
		bot.Username2UserId[update.Message.From.Username] = update.Message.From.ID
	}

	for _, command := range commandsList {
		if command.Predicate(bot, update) {
			err := command.Executor(bot, update)
			if err != nil {
				log.Printf("Failed to reply to [%s], error: %s", update.Message.From.Username, err)
				break
			}
		}
	}
}
