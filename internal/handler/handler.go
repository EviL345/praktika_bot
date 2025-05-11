package handler

import (
	"github.com/EviL345/praktika_bot/internal/config"
	"gopkg.in/telebot.v4"
)

type Repository interface {
	GetTopicId(userId int64) int
	CreateTopic(userId int64, topicId int)
}

type Handler struct {
	cfg  *config.Config
	bot  *telebot.Bot
	Repo Repository
}

func New(cfg *config.Config, bot *telebot.Bot, repo Repository) *Handler {
	return &Handler{
		cfg:  cfg,
		bot:  bot,
		Repo: repo,
	}
}
