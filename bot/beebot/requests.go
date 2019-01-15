package beebot

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const endpoint = "https://chatbot.cn-shanghai.aliyuncs.com"

var replacer = strings.NewReplacer(
	"+", "%20",
	"*", "%2A",
	"%7E", "~",
)

type ChatRequest map[string]string

type ResponseMessage struct {
	Type string
	Text struct {
		Content string
		Source  string
	}
	Knowledge struct {
		ID      string
		Title   string
		Summary string
		Content string
	}
	Recommends []struct {
		KnowledgeID  string
		Title        string
		AnswerSource string
	}
}

type ChatResponse struct {
	Code    string
	Message string

	MessageID string
	SessionID string
	Messages  []ResponseMessage
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (bot *BeeBot) NewRequest() *ChatRequest {
	return &ChatRequest{
		"Format":           "JSON",
		"Version":          "2017-10-11",
		"AccessKeyId":      bot.AccessKeyID,
		"SignatureMethod":  "HMAC-SHA1",
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"SignatureVersion": "1.0",
		"SignatureNonce":   RandStringRunes(10),
	}
}

var lstSignStr string

func (bot BeeBot) Sign(req *ChatRequest) {
	strToSign := "GET&%2F&"
	queryString := ""
	keys := []string{}
	for k := range *req {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for ind, k := range keys {
		if ind > 0 {
			queryString += "&"
		}
		queryString += k + "=" + url.QueryEscape((*req)[k])
	}
	strToSign += replacer.Replace(url.QueryEscape(queryString))
	lstSignStr = strToSign

	signer := hmac.New(sha1.New, []byte(bot.AccessKeySecret+"&"))
	signer.Write([]byte(strToSign))
	(*req)["Signature"] = base64.URLEncoding.EncodeToString(signer.Sum(nil))
}

func (bot BeeBot) Send(req *ChatRequest) (*ChatResponse, error) {
	bot.Sign(req)

	reqURL, err := url.Parse(endpoint)
	if err != nil {
		log.Println("error when parsing url", err)
	}
	query := reqURL.Query()
	for k, v := range *req {
		query.Add(k, v)
	}
	reqURL.RawQuery = query.Encode()

	response, err := http.DefaultClient.Get(reqURL.String())
	if err != nil {
		return nil, err
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	chatResponse := &ChatResponse{}
	json.Unmarshal(responseBody, chatResponse)
	log.Println(string(responseBody))
	if ind := strings.IndexRune(chatResponse.Message, ':'); ind >= 0 {
		log.Println(chatResponse.Message[ind+1:], lstSignStr)
		log.Println(lstSignStr == chatResponse.Message[ind+1:])
	}
	return chatResponse, nil
}
