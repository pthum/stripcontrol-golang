package messaging

import (
	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/internal/model"
)

// NoOpEventHandler that does nothing
type NoOpEventHandler struct {
}

func (m *NoOpEventHandler) Close() {
}

func (m *NoOpEventHandler) PublishStripSaveEvent(id null.Int, strip model.LedStrip) (err error) {
	return nil
}

func (m *NoOpEventHandler) PublishStripDeleteEvent(id null.Int) (err error) {
	return nil
}

func (m *NoOpEventHandler) PublishProfileSaveEvent(id null.Int, profile model.ColorProfile) (err error) {
	return nil
}

func (m *NoOpEventHandler) PublishProfileDeleteEvent(id null.Int) (err error) {
	return nil
}
