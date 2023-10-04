package commands

import (
	"fmt"
	"furrybot/config"
	"furrybot/femboy"
	"furrybot/images"
	"log"
	"math/rand"
	"strconv"
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

func GetCommandFromUpdate(update *echotron.Update) (string, string) {
	for _, entity := range update.Message.Entities {
		if entity.Type == "bot_command" {
			command := update.Message.Text[entity.Offset : entity.Offset+entity.Length]
			if index := strings.Index(command, "@"); index >= 0 {
				command = command[:index]
			}
			return command, strings.Trim(update.Message.Text[entity.Offset+entity.Length:], " ")
		}
	}

	return "", ""
}

func CreateMessageFullMatchPredicate(commandName string) CommandExecutionPredicate {
	return func(bot *Bot, u *echotron.Update) bool {
		if u.Message == nil {
			return false
		}
		command, _ := GetCommandFromUpdate(u)
		return command == commandName
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

		_, err = bot.SendMessage(fmt.Sprintf("@%s Ты был выбран фембоем!\nТы получил(а) %d cum(s)", memberResp.Result.User.Username, balanceGift), update.ChatID(), nil)
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
				log.Printf("Failed to get username from id: %d\n", memberResp.ErrorCode)
			}
			if memberResp.Result == nil {
				removed++
				bot.FemboyGame.RemovePlayerByUserId(p.UserId)
			}

			msg += fmt.Sprintf("%d. %s - %d раз\n", i+1-removed, memberResp.Result.User.Username, p.Wins)
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

var FuckCommand = Command{
	func(bot *Bot, update *echotron.Update) bool {
		if update.Message == nil {
			return false
		}

		command, params := GetCommandFromUpdate(update)

		if command != "/fuck" {
			return false
		}

		return !strings.Contains(params, " ") && strings.HasPrefix(params, "@")
	},
	func(bot *Bot, update *echotron.Update) error {
		_, params := GetCommandFromUpdate(update)
		target := params[1:]

		if target == bot.BotName {
			penalty := rand.Int63n(bot.Balance.GetBalance(update.Message.From.ID) + 1)
			bot.Balance.DecreaseBalance(update.Message.From.ID, penalty)

			msg := fmt.Sprintf(
				">ш< бота трахать нельзя! За это ты будешь наказан!\n @%s получил штраф в сумме %d cum(s)",
				update.Message.From.Username, penalty,
			)
			_, err := bot.SendMessage(msg, update.ChatID(), nil)
			return err
		}

		userId, ok := bot.Username2UserId[target]

		if !ok {
			_, err := bot.SendMessage("Трахать можно только тех, кто что-то писал в чате!", update.ChatID(), nil)
			return err
		}

		if !bot.Balance.DecreaseBalance(update.Message.From.ID, config.Settings.TrahCost) {
			_, err := bot.SendMessage(
				fmt.Sprintf("У тебя недостаточно cum(s) чтобы трахнуть! Необходимо %d cum(s)", config.Settings.TrahCost),
				update.ChatID(),
				nil,
			)
			return err
		}

		trahBonus := rand.Int63n(config.Settings.TrahCost + 1)
		bot.Balance.IncreaseBalance(userId, trahBonus)

		var msg string

		if update.Message.From.ID == userId {
			msg = fmt.Sprintf("Самотрах!\n @%s трахнул себя и получил %d cum(s)", target, trahBonus)
		} else {
			msg = fmt.Sprintf("@%s трахнул @%s!\n@%s получил %d cum(s)", update.Message.From.Username, target, target, trahBonus)
		}

		_, err := bot.SendMessage(msg, update.ChatID(), nil)

		return err
	},
}

var GiveCumCommand = Command{
	func(bot *Bot, update *echotron.Update) bool {
		command, params := GetCommandFromUpdate(update)
		if command != "/give" {
			return false
		}

		params = strings.Trim(params, " ")
		params_parts := strings.Split(params, " ")

		if len(params_parts) != 2 {
			return false
		}

		if _, err := strconv.Atoi(params_parts[1]); !strings.HasPrefix(params_parts[0], "@") || err != nil {
			return false
		}

		return true
	},
	func(bot *Bot, update *echotron.Update) error {
		_, params := GetCommandFromUpdate(update)
		params = strings.Trim(params, " ")
		params_parts := strings.Split(params, " ")
		target := params_parts[0][1:]
		amount, _ := strconv.Atoi(params_parts[1])

		userId, ok := bot.Username2UserId[target]

		if !ok {
			_, err := bot.SendMessage("Переводить cum(s) только тем, кто что-то писал в чате!", update.ChatID(), nil)
			return err
		}

		if !bot.Balance.DecreaseBalance(update.Message.From.ID, int64(amount)) {
			_, err := bot.SendMessage("У тебя недостаточно средств на счету!", update.ChatID(), nil)
			return err
		}

		bot.Balance.IncreaseBalance(userId, int64(amount))

		_, err := bot.SendMessage(
			fmt.Sprintf("@%s перевёл @%s %d cum(s)", update.Message.From.Username, target, amount),
			update.ChatID(),
			nil,
		)

		return err
	},
}

var SpawnCumCommand = Command{
	func(bot *Bot, update *echotron.Update) bool {
		res, err := bot.GetChatAdministrators(update.ChatID())

		if err != nil {
			return false
		}

		senderIsAdmin := false
		for _, member := range res.Result {
			if member.User.Username == update.Message.From.Username {
				senderIsAdmin = true
				break
			}
		}
		if !senderIsAdmin {
			return false
		}

		command, params := GetCommandFromUpdate(update)

		if command != "/spawn" {
			return false
		}

		params = strings.Trim(params, " ")
		_, err = strconv.Atoi(params)
		return err == nil
	},
	func(bot *Bot, update *echotron.Update) error {
		_, params := GetCommandFromUpdate(update)
		params = strings.Trim(params, " ")
		amount, _ := strconv.Atoi(params)

		bot.Balance.IncreaseBalance(update.Message.From.ID, int64(amount))

		_, err := bot.SendMessage(
			fmt.Sprintf("@%s начислено на твой счёт %d cum(s)", update.Message.From.Username, amount),
			update.ChatID(),
			nil,
		)

		return err
	},
}

var BalanceLeaderboardCommand = Command{
	CreateMessageFullMatchPredicate("/balance_leaderboard"),
	func(bot *Bot, update *echotron.Update) error {
		balanceSlice := bot.Balance.GetSortedBalanceSlice()

		removed := 0
		msg := "Список cum лидеров: \n"
		for i, userBalance := range balanceSlice {
			memberResp, err := bot.GetChatMember(update.ChatID(), userBalance.UserId)
			if err != nil {
				log.Printf("Failed to get username from id: %d\n", memberResp.ErrorCode)
				removed++
			}

			msg += fmt.Sprintf("%d. %s - %d cum(s)\n", i+1-removed, memberResp.Result.User.Username, userBalance.Balance)
		}

		_, err := bot.SendMessage(msg, update.ChatID(), nil)
		return err
	},
}
