package config

import (
	"log"

	"github.com/spf13/viper"
)

var Configuration struct {
	App struct {
		ID      string
		Secret  string
		Baseurl string
	}
	Bot struct {
		Beebot struct {
			ID         string
			Secret     string
			InstanceID string
		}
	}
}

func init() {
}

func OnCobraInitialized() {
	Configuration.App.Baseurl = "https://oapi.zjurl.cn/open-apis/api/"
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	log.Println("Unmarshalling viper configuration from", viper.ConfigFileUsed())
	viper.Unmarshal(&Configuration)
	Configuration.Bot.Beebot.InstanceID = viper.GetString("bot.beebot.instance_id")
}
