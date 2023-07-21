package commands

import (
	"fmt"
	"furrybot/images"
	"log"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type CommandExecutor func(*tgbotapi.Message, *ChatContext) tgbotapi.Chattable

// Checks whether to execute a command or not
type CommandExecutionPredicate func(*tgbotapi.Update, *ChatContext) bool

type Command struct {
	Predicate CommandExecutionPredicate
	Executor  CommandExecutor
}

type ChatContext struct {
	ImageRepository images.IImageRepository
}

func CreateMessageFullMatchPredicate(commandName string) CommandExecutionPredicate {
	return func(u *tgbotapi.Update, ctx *ChatContext) bool {
		return u.Message.Command() == commandName
	}
}

func GetFurryPic(message *tgbotapi.Message, ctx *ChatContext) tgbotapi.Chattable {
	image, err := ctx.ImageRepository.GetRandomImagePath()

	if err != nil {
		log.Printf("Failed to fetch image from repository. Error: %s", err)
		return tgbotapi.NewMessage(message.Chat.ID, "Не удалось получить картинку, попробуйте ещё раз позже")
	}

	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, image)
	return msg
}

var GetFurryPicCommand = Command{
	CreateMessageFullMatchPredicate("get_furry"),
	GetFurryPic,
}

func ListAvailablePics(message *tgbotapi.Message, ctx *ChatContext) tgbotapi.Chattable {
	msg := "List of available images: \n"

	if len(ctx.ImageRepository.GetImages()) == 0 {
		msg += "Empty (This source might not support image listing)"
	} else {
		for _, v := range ctx.ImageRepository.GetImages() {
			msg += v + "\n"
		}
		msg += "Total: " + fmt.Sprint(len(ctx.ImageRepository.GetImages()))
	}

	return tgbotapi.NewMessage(message.Chat.ID, msg)
}

var GetFurryListCommand = Command{
	CreateMessageFullMatchPredicate("get_furry_list"),
	ListAvailablePics,
}
