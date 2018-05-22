package main

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/floodcode/tgbot"
)

const (
	configPath = "config.json"
)

var (
	bot       tgbot.TelegramBot
	botUser   tgbot.User
	botConfig BotConfig
	cmdMatch  *regexp.Regexp
	cmdList   = map[string]func(BotRequest){}
)

// BotConfig contains bot's environment variables
type BotConfig struct {
	Token string `json:"token"`
}

// BotRequest represents bot command
type BotRequest struct {
	msg  tgbot.Message
	cmd  string
	args string
}

// QuickAnswer sends simple text message in reply to origin message
func (req *BotRequest) QuickAnswer(text string) {
	bot.SendMessage(tgbot.SendMessageConfig{
		ChatID:           tgbot.ChatID(req.msg.Chat.ID),
		Text:             text,
		ReplyToMessageID: req.msg.MessageID,
		ParseMode:        tgbot.ParseModeMarkdown(),
	})
}

func main() {
	loadConfig()
	addRoutes()
	startBot()
}

func startBot() {
	var err error
	bot, err = tgbot.New(botConfig.Token)
	checkError(err)

	botUser, err = bot.GetMe()
	checkError(err)

	cmdMatch = regexp.MustCompile(`(?s)^\/([a-zA-Z_]+)(?:@` + botUser.Username + `)?(?:[\s\n]+(.+)|)$`)

	err = bot.Poll(tgbot.PollConfig{
		Delay:    250,
		Callback: updatesCallback,
	})

	checkError(err)
}

func loadConfig() {
	configData, err := ioutil.ReadFile(configPath)
	checkError(err)

	err = json.Unmarshal(configData, &botConfig)
	checkError(err)
}

func getRoute(route string) (callback func(BotRequest), ok bool) {
	callback, ok = cmdList[route]
	return callback, ok
}

func addRoute(route string, callback func(BotRequest)) {
	cmdList[route] = callback
}

func updatesCallback(updates []tgbot.Update) {
	for _, update := range updates {
		if update.Message == nil || len(update.Message.Text) == 0 {
			continue
		}

		processTextMessage(update.Message)
	}
}

func processTextMessage(msg *tgbot.Message) {
	match := cmdMatch.FindStringSubmatch(msg.Text)

	if match == nil {
		return
	}

	processRequest(BotRequest{
		msg:  *msg,
		cmd:  strings.ToLower(match[1]),
		args: match[2],
	})
}

func processRequest(req BotRequest) {
	if callback, ok := getRoute(req.cmd); ok {
		callback(req)
	}
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
