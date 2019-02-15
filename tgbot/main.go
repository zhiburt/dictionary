package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/dictionary/tgbot/commands"
	"github.com/dictionary/tgbot/wheather"

	dict "github.com/dictionary/tgbot/dict"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//-telegrambottoken=503349141:AAFWEFbAnbuqOHomYYpA486qN9JEvT6ezvU
//-chatId=-1001102816482 --dictserviceaddr=:8081 --openwheather=c5781e9363a95ba9aafc4d9820d99a3e
var (
	dictServiceAddr  string = os.Getenv("dictserviceaddr")
	telegramBotToken string = os.Getenv("telegrambottoken")
	apiWheatherKey   string = os.Getenv("openwheather")
	chatID           string = os.Getenv("chatId")
)

func init() {

	dictServiceAddr = os.Getenv("dictserviceaddr")
	telegramBotToken = os.Getenv("telegrambottoken")
	apiWheatherKey = os.Getenv("openwheather")
	chatID = os.Getenv("chatId")
	if dictServiceAddr == "" ||
		telegramBotToken == "" ||
		apiWheatherKey == "" ||
		chatID == "" {
		panic("You need to type ENV variables")
	}
}

// ErrService default error
var ErrService = errors.New("cannot does that thing, try do it later")

func main() {
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Println(telegramBotToken)
		log.Panic(err)
	}
	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	SendTodayWheterToChat(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		go handleCommand(bot, &update)
	}
}

func handleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	var reply string
	commandName := update.Message.Command()
	command := commands.NewCommandByName(commandName)

	var ctx = context.Background()
	dictionary := dict.NewDict(dictServiceAddr)
	log.Printf("dict created %s", dictServiceAddr)
	resp, err := command.Exec(ctx, dictionary, update.Message.Text)
	if err != nil {
		log.Printf("Command [%s] has executed unsuccessully", command.Name())
		reply = ErrService.Error()
	} else {
		log.Printf("Command [%s] has executed successully", command.Name())
		reply = resp
	}

	log.Printf("%s", update.Message.Chat.LastName)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}

// SendTodayWheterToChat gets the wheather today
// and send formatting one to chat
func SendTodayWheterToChat(bot *tgbotapi.BotAPI) {
	welcome := "Hi! Today we have\n"
	w := wheather.WheatherToday(apiWheatherKey)

	msg := tgbotapi.NewMessage(chatIDtoInt64(chatID), welcome+w)
	bot.Send(msg)
}

func chatIDtoInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}
