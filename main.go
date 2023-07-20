package main

import (
	"furrybot/commands"
	"furrybot/config"
	"furrybot/images"
	"log"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var commandsMap = map[string]commands.Command{
	"/get_furry":      commands.GetFurryPic,
	"/get_furry_list": commands.ListAvailablePics,
}

func createChatContext(settings *config.Settings, repository images.IImageRepository) commands.ChatContext {
	return commands.ChatContext{
		Settings:        settings,
		ImageRepository: repository,
	}
}

func main() {
	settings, err := config.ReadSettingsFromJson(config.GetSettingsPath())

	if err != nil {
		log.Fatalln("Failed to load configuration:", err)
	}

	defaultRepository, err := images.NewLocalFilesImageRepository(settings.PicsFolder)

	if err != nil {
		log.Fatalln("Failed to create repository:", err)
	}

	log.Printf("Image repository initialized, loaded %v pics\n", len(defaultRepository.GetImages()))

	// TODO: Each chat should have its own context for
	// users to configure, for example, which repository
	// to use
	ctx := createChatContext(&settings, defaultRepository)

	bot, err := tgbotapi.NewBotAPI(settings.TelegramBotToken)

	if err != nil {
		log.Fatalln("Failed to create bot:", err)
	}

	log.Printf("Bot initialized. Authenticated as %v\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updatesChannel, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Fatalln("Failed to create updates channel:", err)
	}

	log.Println("Waiting for messages...")
	for update := range updatesChannel {
		go handleUpdate(update, bot, &ctx)
	}
}

func handleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI, ctx *commands.ChatContext) {
	if update.Message == nil {
		return
	}
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	for k, command := range commandsMap {
		if update.Message.Text == k {
			reply := command(update.Message, ctx)
			bot.Send(reply)
			break
		}
	}
}
