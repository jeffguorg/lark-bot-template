package web

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jeffguorg/lark-bot-template/bot/beebot"
)

type Message struct {
	Type  string
	Token string
}

type ChallengeMessage struct {
	Type      string
	Token     string
	Challenge string
}

type RichTextEvent struct {
	Type     string
	MsgType  string
	ChatType string `json:"chat_type"`

	UserOpenID    string `json:"user_open_id"`
	OpenChatID    string `json:"open_chat_id"`
	OpenMessageID string `json:"open_message_id"`

	Text  string
	Title string
}

type EventMessage struct {
	Message
	Timestamp float64 `json:"ts,string"`
	Event     RichTextEvent
}

var bbot = &beebot.BeeBot{}
var BotHandler = NewBotHandler(bbot)

func webhook(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(500)
		return
	}
	var msg Message
	if err = json.Unmarshal(body, &msg); err != nil {
		res.WriteHeader(400)
		return
	}
	log.Println("message type is", msg.Type)
	if msg.Type == "url_verification" {
		var challengeMessage ChallengeMessage
		if err := json.Unmarshal(body, &challengeMessage); err != nil {
			log.Println("failed to unmarshal challenge message", err)
			res.WriteHeader(400)
			return
		}
		URLVerification(res, challengeMessage)
		return
	} else if msg.Type == "event_callback" {
		log.Println("entering event callback")
		var textMsg EventMessage
		log.Println("body is", string(body))
		if err := json.Unmarshal(body, &textMsg); err != nil {
			log.Println("failed to unmarshal text message", err)
			res.WriteHeader(400)
			return
		}
		log.Println("handle the request to repeat handler")
		BotHandler(res, textMsg.Event)
	}
}

func init() {
	router.Post("/", webhook)
}
