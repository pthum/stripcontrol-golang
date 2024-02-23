package telegram

import (
	"slices"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pthum/stripcontrol-golang/internal/config"
	alog "github.com/pthum/stripcontrol-golang/internal/log"
	"github.com/samber/do"
)

type tgHandler struct {
	cfg config.TelegramConfig
	bot *tgbotapi.BotAPI
	ch  *cmdHandler
	l   alog.Logger
}

func NewHandler(i *do.Injector, cfg config.TelegramConfig) *tgHandler {
	var bot *tgbotapi.BotAPI
	l := alog.NewLogger("telegram")
	if cfg.Enable {
		var err error
		tgbotapi.SetLogger(alog.NewLogLogger("telegram-api"))
		bot, err = tgbotapi.NewBotAPI(cfg.BotKey)
		if err != nil {
			l.Error("%v", err)
			return nil
		}
		if cfg.EnableDebug {
			bot.Debug = true
		}
		l.Info("Authorized on account %s", bot.Self.UserName)
	}
	s := NewCmdHandler(i)

	return &tgHandler{
		cfg: cfg,
		bot: bot,
		ch:  s,
		l:   l,
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
		h.l.Warn("message from user %v with ID %v isn't contained in allowed list", update.Message.From.UserName, update.Message.From.ID)
		return
	}
	h.l.Debug("[%s] %s", update.Message.From.UserName, update.Message.Text)

	resMsg := h.ch.callForMsg(update.Message)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, resMsg)
	msg.ReplyToMessageID = update.Message.MessageID

	_, err := h.bot.Send(msg)
	if err != nil {
		h.l.Info("Error sending telegram message: %v", err)
	}
}
