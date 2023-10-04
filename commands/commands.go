package commands

import (
	"fmt"
	"furrybot/config"
	"furrybot/femboy"
	"furrybot/images"
	"log"
	"math/rand"
	"strings"

	"github.com/NicoNex/echotron/v3"
)

type CommandExecutor func(bot *Bot, update *echotron.Update) error

// Checks whether to execute a command or not
type CommandExecutionPredicate func(bot *Bot, update *echotron.Update) bool

type Command struct {
	Predicate CommandExecutionPredicate
	Executor  CommandExecutor
}

func CreateMessageFullMatchPredicate(commandName string) CommandExecutionPredicate {
	return func(bot *Bot, u *echotron.Update) bool {
		return u.Message != nil && u.Message.Text == commandName
	}
}

var GetFurryPicCommand = Command{
	CreateMessageFullMatchPredicate("/get_furry"),
	func(bot *Bot, update *echotron.Update) error {
		image, err := bot.ImageRepository.GetRandomImagePath()

		if err != nil {
			bot.SendMessage("Не удалось получить картинку, попробуйте ещё раз позже", update.ChatID(), nil)
			return err
		}

		_, err = bot.SendPhoto(echotron.NewInputFilePath(image), update.ChatID(), &echotron.PhotoOptions{
			HasSpoiler: bot.ImageRepository.IsCensored(),
		})
		return err
	},
}

var GetFurryListCommand = Command{
	CreateMessageFullMatchPredicate("/get_furry_list"),
	func(bot *Bot, update *echotron.Update) error {
		msg := "List of available images: \n"

		if len(bot.ImageRepository.GetImages()) == 0 {
			msg += "Empty (This source might not support image listing)"
		} else {
			for _, v := range bot.ImageRepository.GetImages() {
				msg += v + "\n"
			}
			msg += "Total: " + fmt.Sprint(len(bot.ImageRepository.GetImages()))
		}

		_, err := bot.SendMessage(msg, update.ChatID(), nil)
		return err
	},
}

const SELECT_REPOSITORY_PREFIX = "select-repository:"

var ShowRepositorySelectionCommand = Command{
	CreateMessageFullMatchPredicate("/show_repositories"),
	func(bot *Bot, update *echotron.Update) error {
		keys := [][]echotron.InlineKeyboardButton{
			{
				{
					Text:         "Коллекция авторов бота 😈",
					CallbackData: SELECT_REPOSITORY_PREFIX + "local",
				},
			},
			{
				{
					Text:         "Reactor ⚛",
					CallbackData: SELECT_REPOSITORY_PREFIX + "reactor",
				},
			},
			{
				{
					Text:         "FAP Reactor 🍆",
					CallbackData: SELECT_REPOSITORY_PREFIX + "fap_reactor",
				},
			},
		}

		_, err := bot.SendMessage("🐈 Выберите источник картинок", update.ChatID(), &echotron.MessageOptions{
			ReplyMarkup: echotron.InlineKeyboardMarkup{InlineKeyboard: keys},
		})

		return err
	},
}

var SelectRepositoryCommand = Command{
	func(bot *Bot, update *echotron.Update) bool {
		return update.CallbackQuery != nil && strings.HasPrefix(update.CallbackQuery.Data, "select-repository:")
	},
	func(bot *Bot, update *echotron.Update) error {
		repository_name := ""

		switch update.CallbackQuery.Data[len(SELECT_REPOSITORY_PREFIX):] {
		case "local":
			repository, err := images.NewLocalFilesImageRepository(config.Settings.PicsFolder)

			if err != nil {
				bot.AnswerCallbackQuery(update.CallbackQuery.ID, &echotron.CallbackQueryOptions{
					Text: "Что-то пошло не так",
				})
				return nil
			}

			bot.ImageRepository = repository
			repository_name = "коллекция авторов бота"
		case "reactor":
			bot.ImageRepository = &images.ReactorImageRepository{}
			repository_name = "Reactor"
		case "fap_reactor":
			bot.ImageRepository = &images.FapReactorImageRepository{}
			repository_name = "Fap Reactor"
		default:
			bot.AnswerCallbackQuery(update.CallbackQuery.ID, &echotron.CallbackQueryOptions{
				Text: "Что-то пошло не так",
			})
			return nil
		}

		bot.AnswerCallbackQuery(update.CallbackQuery.ID, nil)
		bot.DeleteMessage(update.ChatID(), update.CallbackQuery.Message.ID)
		_, err := bot.SendMessage(fmt.Sprintf("Выбран источник \"%s\"", repository_name), update.ChatID(), nil)
		return err
	},
}

var OlegShipulinCommand = Command{
	CreateMessageFullMatchPredicate("/oleg_shipulin"),
	func(bot *Bot, update *echotron.Update) error {
		if update.Message.From.Username == "real_chilll" {
			_, err := bot.SendMessage("ТЫ ОЛЕГ ШИПУЛИН 🔥🔥🔥🔥🔥", update.ChatID(), nil)
			return err
		} else {
			_, err := bot.SendMessage("ты не олег шипулин 😿", update.ChatID(), nil)
			return err
		}
	},
}

// var Fuck = Command{
// 	CreateMessageFullMatchPredicate("fuck"),
// 	func(bot *Bot, update *echotron.Update) error {
// 		victim := u.Message.CommandArguments()
// 		fmt.Println(victim)
// 		return tgbotapi.NewMessage(u.Message.Chat.ID, "Ты трахнул "+victim)
// 	},
// }

var FemboyRegisterCommand = Command{
	CreateMessageFullMatchPredicate("/femboy_register"),
	func(bot *Bot, update *echotron.Update) error {

		if bot.FemboyGame.RegisterPlayer(update.Message.From.ID) {
			_, err := bot.SendMessage("Теперь ты играешь в фембоев!", update.ChatID(), nil)
			return err
		} else {
			_, err := bot.SendMessage("Ты уже играешь в фембоев!", update.ChatID(), nil)
			return err
		}
	},
}

// TODO: Users who aren't registered shouldn't be able to execute this command
var ChooseTodaysFemboyCommand = Command{
	CreateMessageFullMatchPredicate("/femboy"),
	func(bot *Bot, update *echotron.Update) error {
		winnerId, err := bot.FemboyGame.PickWinner()

		if rlerr, ok := err.(*femboy.RateLimitError); ok {
			_, err := bot.SendMessage(
				fmt.Sprintf("Вы слишком часто вызываете фембоя~\nДайте ботику отдохнуть ещё %d секунд -w-", (rlerr.TimeLeftMs)/1000),
				update.ChatID(),
				nil,
			)
			return err
		}

		if _, ok := err.(*femboy.NoPlayersError); ok {
			_, err := bot.SendMessage(
				"Ещё никто не играет в фембоя! Присоединись к игре с помощью команды /femboy_register",
				update.ChatID(),
				nil,
			)
			return err
		}

		memberResp, err := bot.GetChatMember(update.ChatID(), winnerId)
		if err != nil {
			return err
		}
		if memberResp.Result == nil {
			_, err := bot.SendMessage("Фембойчик был выбран, но похоже, что он уже вышел из чата, попробуйте ещё раз!", update.ChatID(), nil)
			bot.FemboyGame.RemovePlayerByUserId(winnerId)
			return err
		}

		balanceGift := rand.Int63n(config.Settings.MaxfemboyBonus - config.Settings.MinFemboyBonus + 1)
		balanceGift += config.Settings.MinFemboyBonus
		bot.Balance.IncreaseBalance(winnerId, balanceGift)

		_, err = bot.SendMessage(fmt.Sprintf("@%s Ты был выбран фембоем!\nОн(а) получил(а) %d cum(s)", memberResp.Result.User.Username, balanceGift), update.ChatID(), nil)
		return err
	},
}

var ShowLeaderboardCommand = Command{
	CreateMessageFullMatchPredicate("/femboy_leaderboard"),
	func(bot *Bot, update *echotron.Update) error {
		players := bot.FemboyGame.GetSortedPlayerSlice()

		if len(players) == 0 {
			_, err := bot.SendMessage("Список фембоев пуст 😿", update.ChatID(), nil)
			return err
		}

		removed := 0
		msg := "Список фембой лидеров: \n"
		for i, p := range players {
			memberResp, err := bot.GetChatMember(update.ChatID(), p.UserId)
			if err != nil {
				log.Printf("Failed to get username from id: %s\n", memberResp.ErrorCode)
			}
			if memberResp.Result == nil {
				removed++
				bot.FemboyGame.RemovePlayerByUserId(p.UserId)
			}

			msg += fmt.Sprintf("%d. %s - %d раз\n", i+1, memberResp.Result.User.Username, p.Wins)
		}

		_, err := bot.SendMessage(msg, update.ChatID(), nil)

		return err
	},
}

var ShowBalanceCommand = Command{
	CreateMessageFullMatchPredicate("/balance"),
	func(bot *Bot, update *echotron.Update) error {
		balance := bot.Balance.GetBalance(update.Message.From.ID)

		_, err := bot.SendMessage(fmt.Sprintf("На твоём счету %d cum(s)", balance), update.ChatID(), &echotron.MessageOptions{
			ReplyToMessageID: update.Message.ID,
		})

		return err
	},
}
