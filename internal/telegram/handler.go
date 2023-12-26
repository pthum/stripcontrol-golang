package telegram

import (
	"log"

	"slices"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pthum/stripcontrol-golang/internal/config"
	"github.com/samber/do"
)

type tgHandler struct {
	cfg config.TelegramConfig
	bot *tgbotapi.BotAPI
	ch  *cmdHandler
}

func NewHandler(i *do.Injector, cfg config.TelegramConfig) *tgHandler {
	var bot *tgbotapi.BotAPI
	if cfg.Enable {
		var err error
		bot, err = tgbotapi.NewBotAPI(cfg.BotKey)
		if err != nil {
			log.Panic(err)
		}
		if cfg.EnableDebug {
			bot.Debug = true
		}
		log.Printf("Authorized on account %s", bot.Self.UserName)
	}
	s := NewCmdHandler(i)

	return &tgHandler{
		cfg: cfg,
		bot: bot,
		ch:  s,
	}
}

func (h *tgHandler) Handle() {
	if h.bot == nil {
		return
	}
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

	resMsg := h.ch.callForMsg(update.Message)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, resMsg)
	msg.ReplyToMessageID = update.Message.MessageID

	_, err := h.bot.Send(msg)
	if err != nil {
		log.Printf("Error sending telegram message: %v", err)
	}
}
