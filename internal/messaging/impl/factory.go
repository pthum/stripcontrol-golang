package messagingimpl

import (
	"github.com/pthum/stripcontrol-golang/internal/config"
	"github.com/pthum/stripcontrol-golang/internal/messaging"
	"github.com/samber/do"
)

func New(inj *do.Injector) (messaging.EventHandler, error) {
	acfg := do.MustInvoke[*config.Config](inj)
	// return NoOp implementation if disabled
	if acfg.Messaging.Disabled {
		return &NoOpEventHandler{}, nil
	}

	return NewMQTT(acfg.Messaging), nil
}
