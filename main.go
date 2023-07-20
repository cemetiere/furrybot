package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"log"
	"math/rand"
	"os"
	"time"
)

type Config struct {
	TelegramBotToken string `json:"TelegramBotToken"`
	PicsFolder       string `json:"picsFolder"`
}

var configuration = Config{}
var configFilePath = "config.json"
var files []os.DirEntry

func main() {
	readConfig()
	files, _ = os.ReadDir(configuration.PicsFolder)
	bot, _ := tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)
	msgCh := make(chan tgbotapi.Update)

	go func() {
		for update := range updates {
			msgCh <- update
		}
	}()

	for update := range msgCh {
		go handleMessage(update, bot)
	}
}

func handleMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message == nil {
		return
	}
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	if update.Message.Text == "/get_furry" {
		filePath := configuration.PicsFolder + getRandomImage()
		msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, filePath)
		bot.Send(msg)
		log.Printf("Response %s", filePath)
	}
}
func getRandomImage() string {
	rand.Seed(time.Now().UTC().UnixNano())
	return files[rand.Intn(len(files))].Name()
}
func readConfig() {
	f, err := os.Open(configFilePath)

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("An error occurred while reading config")
	}
	fmt.Println("Config loaded")
}
