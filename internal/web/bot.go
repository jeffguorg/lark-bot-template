package web

import (
	"log"
	"net/http"

	"github.com/jeffguorg/lark-bot-template/internal/config"

	"github.com/jeffguorg/lark-bot-template/bot"
	"github.com/jeffguorg/lark-bot-template/internal/background"
	"github.com/jeffguorg/lark-bot-template/lark"
)

var userSessions = make(map[string]string)

func NewBotHandler(b bot.Bot) func(res http.ResponseWriter, event RichTextEvent) {
	return func(res http.ResponseWriter, event RichTextEvent) {
		var (
			chatResponse *bot.ChatResponse
			err          error
		)
		if ses, ok := userSessions[event.UserOpenID+event.OpenChatID]; ok {
			var chatRequest = bot.ChatRequest{
				SessionID: ses,
				Message:   b.Regularization(event.Text),
			}
			chatResponse, err = b.Chat(chatRequest)
			if err != nil {
				log.Println("failed to get response from bot", err)
				return
			}
			if chatResponse.SessionID != "" {
				userSessions[event.UserOpenID+event.OpenChatID] = chatResponse.SessionID
			}

		} else {
			var chatRequest = bot.ChatRequest{
				Message: event.Text,
			}
			chatResponse, err = b.Chat(chatRequest)
			if err != nil {
				log.Println("failed to get response from bot", err)
				return
			}
			if chatResponse.SessionID != "" {
				userSessions[event.UserOpenID+event.OpenChatID] = chatResponse.SessionID
			}
		}

		if event.ChatType == "group" {
			if err := background.LarkClient.SendMessage(lark.NewSimpleGroupMessage(event.OpenChatID, chatResponse.Message)); err != nil {
				log.Println("failed to send message", err)
			}
		} else if event.ChatType == "private" {
			if err := background.LarkClient.SendMessage(lark.NewSimplePrivateMessage(event.UserOpenID, chatResponse.Message)); err != nil {
				log.Println("failed to send message", err)
			}
		}

	}
}
func BotOnInitialized() {
	bbot.AccessKeyID = config.Configuration.Bot.Beebot.ID
	bbot.AccessKeySecret = config.Configuration.Bot.Beebot.Secret
	bbot.InstanceID = config.Configuration.Bot.Beebot.InstanceID
}
