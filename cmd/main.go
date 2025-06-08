package main

import (
	"log"
	"time"

	"github.com/EviL345/praktika_bot/internal/config"
	"github.com/EviL345/praktika_bot/internal/database"
	"github.com/EviL345/praktika_bot/internal/handler"
	tele "gopkg.in/telebot.v4"
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
	paths := []string{
		tele.OnText,
		tele.OnPhoto,
		tele.OnVideo,
		tele.OnDocument,
		tele.OnAudio,
		tele.OnVoice,
		tele.OnSticker,
		tele.OnLocation,
		tele.OnContact,
		tele.OnPoll,
		tele.OnAnimation,
	}
	hndlr := handler.New(cfg, bot, db)
	for _, p := range paths {
		bot.Handle(p, hndlr.MsgHandler)
	}
	bot.Handle("/start", hndlr.HandleStart)

	log.Println("Бот запущен...")
	bot.Start()
}

