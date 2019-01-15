package background

import (
	"net/http"
	"time"

	"github.com/jeffguorg/lark-bot-template/internal/config"

	"github.com/jeffguorg/lark-bot-template/lark"
)

var httpclient = http.Client{
	Timeout: time.Second * 15,
}

var LarkClient *lark.Client

func OnCobraInitialized() {
	LarkClient = &lark.Client{
		AppID:     config.Configuration.App.ID,
		AppSecret: config.Configuration.App.Secret,
		BaseURL:   config.Configuration.App.Baseurl,
	}
	LarkClient.Token.NeedRefreshing = time.Unix(0, 0)
	go ContinousRefreshToken()
}
