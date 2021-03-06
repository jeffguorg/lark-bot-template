package beebot

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/jeffguorg/lark-bot-template/bot"
)

type BeeBot struct {
	AccessKeyID     string
	AccessKeySecret string

	InstanceID string
}

func must(o interface{}, err error) interface{} {
	return o
}

var regexReplacer = []*regexp.Regexp{
	must(regexp.Compile("<at.*>.*</at>")).(*regexp.Regexp),
}

func (b *BeeBot) Regularization(msg string) string {
	for _, replacer := range regexReplacer {
		msg = replacer.ReplaceAllString(msg, "")
	}
	return msg
}

func (b *BeeBot) Chat(req bot.ChatRequest) (*bot.ChatResponse, error) {
	var response *ChatResponse
	var err error

	for response == nil || response.Code == "SignatureDoesNotMatch" {
		chatReq := b.NewRequest()

		(*chatReq)["Action"] = "Chat"
		(*chatReq)["InstanceId"] = b.InstanceID
		(*chatReq)["Utterance"] = req.Message
		if req.SessionID != "" {
			(*chatReq)["SessionId"] = req.SessionID
		}
		response, err = b.Send(chatReq)
		if err != nil {
			return nil, err
		}
	}

	chatResponse := &bot.ChatResponse{
		SessionID: response.SessionID,
	}

	for _, msg := range response.Messages {
		if matched, err := regexp.MatchString("\n$", chatResponse.Message); err == nil {
			if matched && len(chatResponse.Message) > 0 {
				chatResponse.Message += "\n\n"
			}
		} else {
			log.Println("error processing response", err)
			return nil, err
		}
		switch msg.Type {
		case "Text":
			chatResponse.Message += msg.Text.Content
			break
		case "Knowledge":
			chatResponse.Message += msg.Knowledge.Summary
			break
		case "Recommend":
			if strings.Index(chatResponse.Message, "Try this: \n") == -1 {
				chatResponse.Message += "Try this: \n\n"
			}
			for ind, rec := range msg.Recommends {
				chatResponse.Message += fmt.Sprintf("\t%d: %s\n", ind, rec.Title)
			}
			break
		default:
			chatResponse.Message += msg.Type
		}
	}
	return chatResponse, nil
}
