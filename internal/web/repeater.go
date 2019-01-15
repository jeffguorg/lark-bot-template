package web

import (
	"log"
	"net/http"

	"github.com/jeffguorg/lark-bot-template/internal/background"

	"github.com/jeffguorg/lark-bot-template/lark"
)

func RepeaterHandler(res http.ResponseWriter, event RichTextEvent) {
	log.Println("repeat handler entered. chat type is", event.ChatType)
	if event.ChatType == "group" {
		if err := background.LarkClient.SendMessage(lark.NewSimpleGroupMessage(event.OpenChatID, event.Text)); err != nil {
			log.Println(err)
		}
	} else if event.ChatType == "private" {
		if err := background.LarkClient.SendMessage(lark.NewSimplePrivateMessage(event.UserOpenID, event.Text)); err != nil {
			log.Println(err)
		}
	}
}
