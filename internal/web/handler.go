package web

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Message struct {
	Token     string
	Challenge string
}

type ChallengeResponse struct {
	Challenge string
}

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
	if len(msg.Challenge) > 0 {
		var response = ChallengeResponse{
			Challenge: msg.Challenge,
		}
		responseBuffer, _ := json.Marshal(&response)
		res.Write(responseBuffer)
		return
	} else {
		log.Println(string(body))
	}
}

func init() {
	router.Get("/", webhook)
}
