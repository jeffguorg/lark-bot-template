package web

import (
	"encoding/json"
	"net/http"
)

type ChallengeResponse struct {
	Message
	Challenge string
}

func URLVerification(res http.ResponseWriter, msg ChallengeMessage) {
	var response = ChallengeResponse{
		Challenge: msg.Challenge,
	}
	responseBuffer, _ := json.Marshal(&response)
	res.Write(responseBuffer)
}
