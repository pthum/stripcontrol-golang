package messagingimpl

import (
	"testing"

	"github.com/pthum/stripcontrol-golang/internal/config"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
)

func TestNewNoOp(t *testing.T) {
	cfg := config.MessagingConfig{
		Disabled: true,
	}
	inj := provideCfg(cfg)
	mh, err := New(inj)
	assert.NoError(t, err)
	_, ok := mh.(*NoOpEventHandler)
	assert.True(t, ok)
}

func TestNewMQTT(t *testing.T) {
	cfg := config.MessagingConfig{
		Disabled: false,
	}
	inj := provideCfg(cfg)
	mh, err := New(inj)
	assert.NoError(t, err)
	_, ok := mh.(*mqttHandler)
	assert.True(t, ok)
}

func provideCfg(cfg config.MessagingConfig) *do.Injector {
	acfg := config.Config{
		Messaging: cfg,
	}
	inj := do.New()
	do.ProvideValue(inj, &acfg)
	return inj
}
