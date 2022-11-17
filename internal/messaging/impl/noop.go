package messagingimpl

import (
	"github.com/pthum/stripcontrol-golang/internal/model"
)

// NoOpEventHandler that does nothing
type NoOpEventHandler struct {
}

func (m *NoOpEventHandler) Close() {
	//no op
}

func (m *NoOpEventHandler) PublishStripEvent(event *model.StripEvent) error {
	return nil
}

func (m *NoOpEventHandler) PublishProfileEvent(event *model.ProfileEvent) error {
	return nil
}
