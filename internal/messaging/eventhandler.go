package messaging

import (
	"github.com/pthum/stripcontrol-golang/internal/model"
)

//go:generate mockery --name=EventHandler --with-expecter=true
type EventHandler interface {
	Shutdown() error
	PublishProfileEvent(event *model.ProfileEvent) error
	PublishStripEvent(event *model.StripEvent) error
}
