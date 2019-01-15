package lark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type MsgContent struct {
	Text string `json:"text"`
}

type Message struct {
	MsgType string     `json:"msg_type"`
	Content MsgContent `json:"content"`
}

type PrivateMessage struct {
	Message
	OpenID string `json:"open_id"`
}

type GroupMessage struct {
	Message
	ChatID string `json:"chat_id"`
}

func NewSimplePrivateMessage(openid, msg string) *PrivateMessage {
	return &PrivateMessage{
		OpenID: openid,
		Message: Message{
			MsgType: "text",
			Content: MsgContent{
				Text: msg,
			},
		},
	}
}
func NewSimpleGroupMessage(chatid, msg string) *GroupMessage {
	return &GroupMessage{
		ChatID: chatid,
		Message: Message{
			MsgType: "text",
			Content: MsgContent{
				Text: msg,
			},
		},
	}
}

func (lark *Client) SendMessage(msg interface{}) error {
	buf, _ := json.Marshal(msg)
	bufReader := bytes.NewReader(buf)

	baseURL, _ := url.Parse(lark.BaseURL)
	endpoint := "/open-apis/v3/message/send/"
	finalURL, _ := baseURL.Parse(endpoint)

	res := lark.Post(finalURL.String(), bufReader)
	if res != nil {
		var resp CommonResponse
		if err := json.Unmarshal(res, &resp); err != nil {
			return err
		}
		if resp.Code == 0 {
			return nil
		}
		return fmt.Errorf("Error code %d", resp.Code)
	}
	return fmt.Errorf(string(res))
}
