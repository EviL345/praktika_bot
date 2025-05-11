package handler

import (
	tele "gopkg.in/telebot.v4"
	"log"
)

func (h *Handler) TextHandler(c tele.Context) error {
	msg := c.Message()
	user := msg.Sender

	cht := &tele.Chat{ID: h.cfg.ChatId}
	topicId := h.Repo.GetTopicId(user.ID)
	if topicId == 0 {
		topic := &tele.Topic{
			Name: user.FirstName,
		}

		newTopic, err := h.bot.CreateTopic(cht, topic)
		log.Println("newTopic", newTopic)
		log.Println(user)
		if err != nil {
			log.Printf("Failed to create topic: %v", err)
			return c.Send("Failed to create topic: " + err.Error())
		}

		h.Repo.CreateTopic(user.ID, newTopic.ThreadID)

		topicId = newTopic.ThreadID
	}

	_, err := h.bot.Forward(cht, msg, &tele.SendOptions{
		ThreadID: topicId,
	})

	return err
}
