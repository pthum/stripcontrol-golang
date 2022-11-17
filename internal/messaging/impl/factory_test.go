package messagingimpl

import (
	"testing"

	"github.com/pthum/stripcontrol-golang/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewNoOp(t *testing.T) {
	cfg := config.MessagingConfig{
		Disabled: true,
	}
	mh := New(cfg)
	_, ok := mh.(*NoOpEventHandler)
	assert.True(t, ok)
}

func TestNewMQTT(t *testing.T) {
	cfg := config.MessagingConfig{
		Disabled: false,
	}
	mh := New(cfg)
	_, ok := mh.(*mqttHandler)
	assert.True(t, ok)
}
