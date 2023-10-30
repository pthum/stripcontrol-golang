package messagingimpl

import (
	"github.com/pthum/stripcontrol-golang/internal/model"
)

// NoOpEventHandler that does nothing
type NoOpEventHandler struct {
}

func (m *NoOpEventHandler) Shutdown() error {
	//no op
	return nil
}

func (m *NoOpEventHandler) PublishStripEvent(event *model.StripEvent) error {
	return nil
}

func (m *NoOpEventHandler) PublishProfileEvent(event *model.ProfileEvent) error {
	return nil
}
