package commands

import (
	"fmt"
	"furrybot/config"
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
