package lark

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type TenantAccessTokenRequest struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}
type TenantAccessTokenResponse struct {
	CommonResponse
	Token    string `json:"tenant_access_token"`
	Duration int    `json:"expire"`
}

func (lark *Client) RefreshTenantAccessToken() {
	urlBase, _ := url.Parse(lark.BaseURL)
	refreshURL, _ := urlBase.Parse("/open-apis/v3/auth/tenant_access_token/internal/")

	requestForm := TenantAccessTokenRequest{
		AppID:     lark.AppID,
		AppSecret: lark.AppSecret,
	}
	requestBuf, _ := json.Marshal(requestForm)
	requestBufIO := bytes.NewReader(requestBuf)

	res, err := http.Post(refreshURL.String(), "", requestBufIO)
	if err != nil {
		log.Println(err)
		return
	}
	responseBuf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var tokenResponse = &TenantAccessTokenResponse{}
	if err = json.Unmarshal(responseBuf, tokenResponse); err != nil {
		log.Println(err)
		return
	}
	if tokenResponse.Code != 0 {
		log.Println(string(responseBuf))
		return
	}
	lark.Token.TenantAccessToken = tokenResponse.Token
	lark.Token.ExpiresAt = time.Now().Add(time.Duration(tokenResponse.Duration) * time.Second)
	lark.Token.NeedRefreshing = time.Now().Add(time.Duration(tokenResponse.Duration/2) * time.Second)
}
