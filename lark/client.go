package lark

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type CommonResponse struct {
	Code int
}

type Client struct {
	BaseURL   string
	AppID     string
	AppSecret string

	Token struct {
		TenantAccessToken string
		ExpiresAt         time.Time
		NeedRefreshing    time.Time
	}
}

var httpclient = http.Client{
	Timeout: time.Second * 15,
}

func (lark *Client) Post(url string, body io.Reader) []byte {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Println("failed to create request", err)
		return nil
	}
	req.Header.Add("Authorization", "Bearer "+lark.Token.TenantAccessToken)
	res, err := httpclient.Do(req)
	if err != nil {
		log.Println("failed to send post http request", err)
		return nil
	}
	if res.StatusCode/100 != 2 {
		log.Println("unexpected status code", res.StatusCode)
		return nil
	}
	buf, _ := ioutil.ReadAll(res.Body)
	return buf
}
