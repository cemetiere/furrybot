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
	AndreyALUCommand,
	FemboyRegisterCommand,
	ChooseTodaysFemboyCommand,
	ShowLeaderboardCommand,
	ShowBalanceCommand,
	FuckCommand,
	GiveCumCommand,
	SpawnCumCommand,
	BalanceLeaderboardCommand,
}

type Bot struct {
	echotron.API
	ChatId          int64
	BotName         string
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
	if bot.BotName == "" {
		res, _ := bot.GetMe()
		if res.Result != nil {
			bot.BotName = res.Result.Username
		}
	}

	if update.Message != nil {
		command, _ := GetCommandFromUpdate(update)
		log.Printf("[%s] %s | %s", update.Message.From.Username, update.Message.Text, command)
		bot.Username2UserId[update.Message.From.Username] = update.Message.From.ID
	}

	for _, command := range commandsList {
		if command.Predicate(bot, update) {
			err := command.Executor(bot, update)
			if err != nil {
				log.Printf("Failed to reply to [%s], error: %s", update.Message.From.Username, err)
			}
			break
		}
	}
}
