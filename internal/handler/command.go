package handler

import (
	tele "gopkg.in/telebot.v4"
	"log"
)

func (h *Handler) MsgHandler(c tele.Context) error {
	msg := c.Message()
	user := msg.Sender
	mainChatID := h.cfg.ChatId

	if msg.Chat.ID != mainChatID {
		return h.handleUserMessage(c, user, msg, mainChatID)
	}
	return h.handleMainChatMessage(msg, mainChatID)
}

func (h *Handler) handleUserMessage(c tele.Context, user *tele.User, msg *tele.Message, mainChatID int64) error {

	chat := &tele.Chat{ID: mainChatID}
	topicId := h.Repo.GetTopicId(user.ID)
	if topicId == 0 {
		topic := &tele.Topic{
			Name: user.FirstName,
		}

		newTopic, err := h.bot.CreateTopic(chat, topic)
		log.Println("newTopic", newTopic)
		log.Println(user)
		if err != nil {
			log.Printf("Failed to create topic: %v", err)
			return c.Send("Failed to create topic: " + err.Error())
		}

		h.Repo.CreateTopic(user.ID, newTopic.ThreadID)

		topicId = newTopic.ThreadID
	}

	_, err := h.bot.Forward(chat, msg, &tele.SendOptions{
		ThreadID: topicId,
	})

	return err

}

func (h *Handler) handleMainChatMessage(msg *tele.Message, mainChatID int64) error {
	if msg.ThreadID == 0 {
		return nil
	}
	userId := h.Repo.GetUserId(msg.ThreadID)
	if userId == 0 {
		log.Printf("No user found for thread %d", msg.ThreadID)
		return nil
	}
	userCht := &tele.Chat{ID: userId}
	_, err := h.bot.Forward(userCht, msg)
	return err
}
