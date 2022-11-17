package messagingimpl

import (
	"testing"

	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNoOpClose(t *testing.T) {
	handler := getTestInstance()
	handler.Close()
}

func TestPublishStripEvent(t *testing.T) {
	handler := getTestInstance()
	assert.Nil(t, handler.PublishStripEvent(&model.StripEvent{}))
}

func TestPublishProfileEvent(t *testing.T) {
	handler := getTestInstance()
	assert.Nil(t, handler.PublishProfileEvent(&model.ProfileEvent{}))
}

func getTestInstance() *NoOpEventHandler {
	return &NoOpEventHandler{}
}
