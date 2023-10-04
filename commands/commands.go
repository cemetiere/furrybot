package commands

import (
	"fmt"
	"furrybot/config"
	"furrybot/femboy"
	"furrybot/images"
	"log"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type CommandExecutor func(u *tgbotapi.Update, ctx *ChatContext, bot *tgbotapi.BotAPI) tgbotapi.Chattable

// Checks whether to execute a command or not
type CommandExecutionPredicate func(u *tgbotapi.Update, ctx *ChatContext) bool

type Command struct {
	Predicate CommandExecutionPredicate
	Executor  CommandExecutor
}

type ChatContext struct {
	ImageRepository images.IImageRepository
	FemboyGame      *femboy.FemboyGameService
}

func CreateMessageFullMatchPredicate(commandName string) CommandExecutionPredicate {
	return func(u *tgbotapi.Update, ctx *ChatContext) bool {
		return u.Message != nil && u.Message.Command() == commandName
	}
}

var GetFurryPicCommand = Command{
	CreateMessageFullMatchPredicate("get_furry"),
	func(u *tgbotapi.Update, ctx *ChatContext, bot *tgbotapi.BotAPI) tgbotapi.Chattable {
		image, err := ctx.ImageRepository.GetRandomImagePath()

		if err != nil {
			log.Printf("Failed to fetch image from repository. Error: %s", err)
			return tgbotapi.NewMessage(u.Message.Chat.ID, "Не удалось получить картинку, попробуйте ещё раз позже")
		}

		msg := tgbotapi.NewPhotoUpload(u.Message.Chat.ID, image)
		return msg
	},
}

var GetFurryListCommand = Command{
	CreateMessageFullMatchPredicate("get_furry_list"),
	func(u *tgbotapi.Update, ctx *ChatContext, bot *tgbotapi.BotAPI) tgbotapi.Chattable {
		msg := "List of available images: \n"

		if len(ctx.ImageRepository.GetImages()) == 0 {
			msg += "Empty (This source might not support image listing)"
		} else {
			for _, v := range ctx.ImageRepository.GetImages() {
				msg += v + "\n"
			}
			msg += "Total: " + fmt.Sprint(len(ctx.ImageRepository.GetImages()))
		}

		return tgbotapi.NewMessage(u.Message.Chat.ID, msg)
	},
}

const SELECT_REPOSITORY_PREFIX = "select-repository:"

var ShowRepositorySelectionCommand = Command{
	CreateMessageFullMatchPredicate("show_repositories"),
	func(u *tgbotapi.Update, ctx *ChatContext, bot *tgbotapi.BotAPI) tgbotapi.Chattable {
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "🐈 Выберите источник картинок")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Коллекция авторов бота 😈", SELECT_REPOSITORY_PREFIX+"local"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Reactor ⚛", SELECT_REPOSITORY_PREFIX+"reactor"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("FAP Reactor 🍆", SELECT_REPOSITORY_PREFIX+"fap_reactor"),
			),
		)

		return msg
	},
}

var SelectRepositoryCommand = Command{
	func(u *tgbotapi.Update, ctx *ChatContext) bool {
		return u.CallbackQuery != nil && strings.HasPrefix(u.CallbackQuery.Data, "select-repository:")
	},
	func(u *tgbotapi.Update, ctx *ChatContext, bot *tgbotapi.BotAPI) tgbotapi.Chattable {
		repository_name := ""

		switch u.CallbackQuery.Data[len(SELECT_REPOSITORY_PREFIX):] {
		case "local":
			repository, err := images.NewLocalFilesImageRepository(config.Settings.PicsFolder)

			if err != nil {
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(u.CallbackQuery.ID, "Что-то пошло не так"))
				return nil
			}

			ctx.ImageRepository = repository
			repository_name = "коллекция авторов бота"
		case "reactor":
			ctx.ImageRepository = &images.ReactorImageRepository{}
			repository_name = "Reactor"
		case "fap_reactor":
			ctx.ImageRepository = &images.FapReactorImageRepository{}
			repository_name = "Fap Reactor"
		default:
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(u.CallbackQuery.ID, "Что-то пошло не так"))
			return nil
		}

		bot.AnswerCallbackQuery(tgbotapi.NewCallback(u.CallbackQuery.ID, ""))
		bot.DeleteMessage(tgbotapi.NewDeleteMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID))
		return tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Выбран источник \"%s\"", repository_name))
	},
}

var OlegShipulinCommand = Command{
	CreateMessageFullMatchPredicate("oleg_shipulin"),
	func(u *tgbotapi.Update, ctx *ChatContext, bot *tgbotapi.BotAPI) tgbotapi.Chattable {
		if u.Message.From.UserName == "real_chilll" {
			return tgbotapi.NewMessage(u.Message.Chat.ID, "ТЫ ОЛЕГ ШИПУЛИН 🔥🔥🔥🔥🔥")
		} else {
			return tgbotapi.NewMessage(u.Message.Chat.ID, "ты не олег шипулин 😿")
		}
	},
}

var Fuck = Command{
	CreateMessageFullMatchPredicate("fuck"),
	func(u *tgbotapi.Update, ctx *ChatContext, bot *tgbotapi.BotAPI) tgbotapi.Chattable {
		victim := u.Message.CommandArguments()
		fmt.Println(victim)
		return tgbotapi.NewMessage(u.Message.Chat.ID, "Ты трахнул "+victim)
	},
}

var FemboyRegisterCommand = Command{
	CreateMessageFullMatchPredicate("femboy_register"),
	func(u *tgbotapi.Update, ctx *ChatContext, bot *tgbotapi.BotAPI) tgbotapi.Chattable {

		if ctx.FemboyGame.RegisterPlayer(u.Message.From.UserName) {
			return tgbotapi.NewMessage(u.Message.Chat.ID, "Теперь ты играешь в фембоев!")
		} else {
			return tgbotapi.NewMessage(u.Message.Chat.ID, "Ты уже играешь в фембоев!")
		}
	},
}

// TODO: Users who aren't registered shouldn't be able to execute this command
var ChooseTodaysFemboyCommand = Command{
	CreateMessageFullMatchPredicate("femboy"),
	func(u *tgbotapi.Update, ctx *ChatContext, bot *tgbotapi.BotAPI) tgbotapi.Chattable {
		winnerUsername, err := ctx.FemboyGame.PickWinner()

		if rlerr, ok := err.(*femboy.RateLimitError); ok {
			return tgbotapi.NewMessage(
				u.Message.Chat.ID,
				fmt.Sprintf("Вы слишком часто вызываете фембоя~\nДайте ботику отдохнуть ещё %d секунд -w-", (rlerr.TimeLeftMs)/1000),
			)
		}

		if _, ok := err.(*femboy.NoPlayersError); ok {
			return tgbotapi.NewMessage(
				u.Message.Chat.ID,
				"Ещё никто не играет в фембоя! Присоединись к игре с помощью команды /femboy_register",
			)
		}

		return tgbotapi.NewMessage(u.Message.Chat.ID, fmt.Sprintf("@%s Ты был выбран фембоем!", winnerUsername))
	},
}

var ShowLeaderboardCommand = Command{
	CreateMessageFullMatchPredicate("femboy_leaderboard"),
	func(u *tgbotapi.Update, ctx *ChatContext, bot *tgbotapi.BotAPI) tgbotapi.Chattable {
		players := ctx.FemboyGame.GetSortedPlayerSlice()

		if len(players) == 0 {
			return tgbotapi.NewMessage(u.Message.Chat.ID, "Список фембоев пуст 😿")
		}

		msg := "Список фембой лидеров: \n"
		for i, p := range players {
			msg += fmt.Sprintf("%d. %s - %d раз\n", i+1, p.Username, p.Wins)
		}

		return tgbotapi.NewMessage(u.Message.Chat.ID, msg)
	},
}
