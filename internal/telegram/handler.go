package telegram

import (
	"log"

	"slices"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pthum/stripcontrol-golang/internal/config"
)

type tgHandler struct {
	cfg config.TelegramConfig
	bot *tgbotapi.BotAPI
}

func NewHandler(cfg config.TelegramConfig) *tgHandler {
	bot, err := tgbotapi.NewBotAPI(cfg.BotKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	return &tgHandler{
		cfg: cfg,
		bot: bot,
	}
}

func (h *tgHandler) Handle() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			h.handleUpdate(update)
		}
	}
}

func (h *tgHandler) handleUpdate(update tgbotapi.Update) {
	if !slices.Contains(h.cfg.AllowedUserIDs, update.Message.From.ID) {
		log.Printf("message from user %v with ID %v isn't contained in allowed list", update.Message.From.UserName, update.Message.From.ID)
		return
	}
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	resMsg := callForMsg(update.Message)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, resMsg)
	msg.ReplyToMessageID = update.Message.MessageID

	h.bot.Send(msg)
}
