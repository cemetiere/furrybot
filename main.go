package main

import (
	"furrybot/commands"
	"furrybot/config"
	"furrybot/images"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"log"
)

var commandsList = []commands.Command{
	commands.GetFurryListCommand,
	commands.GetFurryPicCommand,
	commands.ShowRepositorySelectionCommand,
	commands.SelectRepositoryCommand,
	commands.OlegShipulinCommand,
	commands.FemboyRegisterCommand,
	commands.ChooseTodaysFemboyCommand,
	commands.ShowLeaderboardCommand,
}

func createChatContext(repository images.IImageRepository) commands.ChatContext {
	return commands.ChatContext{
		ImageRepository: repository,
	}
}

func main() {
	err := config.ReadSettingsFromJson(config.GetSettingsPath())
	if err != nil {
		log.Fatalln("Failed to load configuration:", err)
	}
	log.Println("Settings loaded")

	// defaultRepository, err := images.NewLocalFilesImageRepository(settings.PicsFolder)

	// if err != nil {
	// 	log.Fatalln("Failed to create repository:", err)
	// }

	// log.Printf("Image repository initialized, loaded %v pics\n", len(defaultRepository.GetImages()))

	defaultRepository := &images.ReactorImageRepository{}

	// TODO: Each chat should have its own context for
	// users to configure, for example, which repository
	// to use
	ctx := createChatContext(defaultRepository)

	bot, err := tgbotapi.NewBotAPI(config.Settings.TelegramBotToken)

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
	if update.Message != nil {
		log.Printf("[%s] %s | %s", update.Message.From.UserName, update.Message.Text, update.Message.Command())
	}

	for _, command := range commandsList {
		if command.Predicate(&update, ctx) {
			reply := command.Executor(&update, ctx, bot)
			if reply != nil {
				_, err := bot.Send(reply)
				if err != nil {
					log.Printf("Failed to reply to [%s], error: %s", update.Message.From.UserName, err)
				}
				break
			}
		}
	}
}
