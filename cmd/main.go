package main

import (
	"github.com/EviL345/praktika_bot/internal/config"
	"github.com/EviL345/praktika_bot/internal/database"
	"github.com/EviL345/praktika_bot/internal/handler"
	tele "gopkg.in/telebot.v4"
	"log"
	"time"
)

func main() {
	cfg := config.New()

	db := database.New(&cfg.Db)

	defer db.Close()

	bot, err := tele.NewBot(tele.Settings{
		Token:  cfg.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	hndlr := handler.New(cfg, bot, db)

	bot.Handle(tele.OnText, hndlr.TextHandler)

	log.Println("Бот запущен...")
	bot.Start()
}
