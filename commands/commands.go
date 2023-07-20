package commands

import (
	"fmt"
	"furrybot/config"
	"furrybot/images"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Command func(*tgbotapi.Message, *ChatContext) tgbotapi.Chattable

type ChatContext struct {
	ImageRepository images.IImageRepository
	Settings        *config.Settings
}

func GetFurryPic(message *tgbotapi.Message, ctx *ChatContext) tgbotapi.Chattable {
	image := ctx.ImageRepository.GetRandomImagePath()
	msg := tgbotapi.NewPhotoUpload(message.Chat.ID, image)
	return msg
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
