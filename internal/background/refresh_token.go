package background

import (
	"log"
	"time"
)

func ContinousRefreshToken() {
	for {
		if LarkClient == nil {
			continue
		}
		if time.Now().After(LarkClient.Token.NeedRefreshing) {
			LarkClient.RefreshTenantAccessToken()
			log.Println(LarkClient.Token.TenantAccessToken, LarkClient.Token.ExpiresAt)
		}
		time.Sleep(time.Second)
	}
}
