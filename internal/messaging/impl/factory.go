package messagingimpl

import (
	"github.com/pthum/stripcontrol-golang/internal/config"
	"github.com/pthum/stripcontrol-golang/internal/messaging"
)

func New(cfg config.MessagingConfig) messaging.EventHandler {
	// return NoOp implementation if disabled
	if cfg.Disabled {
		return &NoOpEventHandler{}
	}

	return NewMQTT(cfg)
}
