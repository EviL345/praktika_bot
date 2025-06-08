package handler

import (
	"log"

	tele "gopkg.in/telebot.v4"
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

func (h *Handler) HandleStart(c tele.Context) error {
	startText := "📝Этот бот создан для сбора предложений новостей к публикации в телеграм-канале МГТУ «СТАНКИН»."+
	"\nЕсли у тебя есть интересное событие, достижение или важная информация — просто отправь сообщение с подробностями."+
	"\n\n❗️Пожалуйста, не редактируй уже отправленные сообщения — если хочешь изменить материал, просто пришли заново."
	return c.Send(startText)
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

	if err != nil {
		log.Printf("Failed to forward message: %v", err)
		return err
	}

	// Отправляем подтверждение пользователю
	confirmationMsg := "Спасибо! Мы получили ваше сообщение и в ближайшее время рассмотрим его."
	return c.Send(confirmationMsg)
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

	userChat := &tele.Chat{ID: userId}

	var media any

	switch {
	case msg.Text != "":
		media = msg.Text

	case msg.Photo != nil:
		photo := *msg.Photo
		photo.Caption = msg.Caption
		media = &photo

	case msg.Video != nil:
		video := *msg.Video
		video.Caption = msg.Caption
		media = &video

	case msg.Document != nil:
		doc := *msg.Document
		doc.Caption = msg.Caption
		media = &doc

	case msg.Audio != nil:
		audio := *msg.Audio
		audio.Caption = msg.Caption
		media = &audio

	case msg.Voice != nil:
		media = msg.Voice

	case msg.Sticker != nil:
		media = msg.Sticker

	case msg.Location != nil:
		media = msg.Location

	case msg.Contact != nil:
		media = msg.Contact

	case msg.Poll != nil:
		media = msg.Poll

	case msg.Animation != nil:
		anim := *msg.Animation
		anim.Caption = msg.Caption
		media = &anim

	default:
		log.Printf("Unsupported message type from thread %d", msg.ThreadID)
		return nil
	}

	_, err := h.bot.Send(userChat, media)
	return err
}
