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

func main() {
	readConfig()
	bot, _ := tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Text == "/get_furry" {
			filePath := configuration.PicsFolder + getRandomImage()
			msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, filePath)
			bot.Send(msg)
			log.Printf("Response %s", filePath)
		}
	}
}
func getRandomImage() string {
	files, _ := os.ReadDir(configuration.PicsFolder)
	rand.Seed(time.Now().UTC().UnixNano())
	return files[rand.Intn(len(files))].Name()
}
func readConfig() {
	f, _ := os.Open(configFilePath)

	decoder := json.NewDecoder(f)
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("An error occurred while reading config")
	}
	fmt.Println("Config loaded")
}
