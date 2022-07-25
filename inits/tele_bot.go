package inits

import (
	telegram "gopkg.in/telebot.v3"
	"log"
)

func InitBot(pref telegram.Settings) (*telegram.Bot, bool) {
	b, err := telegram.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return nil, true
	}
	return b, false
}
