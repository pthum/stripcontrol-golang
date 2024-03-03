package service

import (
	"errors"
	"strconv"

	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/internal/database"
	alog "github.com/pthum/stripcontrol-golang/internal/log"
	"github.com/pthum/stripcontrol-golang/internal/messaging"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/samber/do"
)

//go:generate mockery --name=LEDService --with-expecter=true --outpkg=servicemocks
type LEDService interface {
	GetAll() ([]model.LedStrip, error)
	GetLEDStrip(id string) (*model.LedStrip, error)
	CreateLEDStrip(mdl *model.LedStrip) error
	UpdateLEDStrip(id string, updMdl model.LedStrip) error
	DeleteLEDStrip(id string) error
	UpdateProfileForStrip(id string, updProf model.ColorProfile) (*model.ColorProfile, error)
	GetProfileForStrip(id string) (*model.ColorProfile, error)
	RemoveProfileForStrip(id string) error
}

type ledSvc struct {
	dbh   database.DBHandler[model.LedStrip]
	cpDbh database.DBHandler[model.ColorProfile]
	mh    messaging.EventHandler
	l     alog.Logger
}

func NewLEDService(i *do.Injector) (LEDService, error) {
	lsdb := do.MustInvoke[database.DBHandler[model.LedStrip]](i)
	cpdb := do.MustInvoke[database.DBHandler[model.ColorProfile]](i)
	mh := do.MustInvoke[messaging.EventHandler](i)
	l := alog.NewLogger("ledservice")
	return &ledSvc{
		dbh:   lsdb,
		cpDbh: cpdb,
		mh:    mh,
		l:     l,
	}, nil
}

func (l *ledSvc) GetAll() ([]model.LedStrip, error) {
	return l.dbh.GetAll()
}
func (l *ledSvc) GetLEDStrip(id string) (*model.LedStrip, error) {
	return l.dbh.Get(id)
}
func (l *ledSvc) CreateLEDStrip(mdl *model.LedStrip) error {
	// generate an id
	mdl.GenerateID()
	l.l.Debug("Generated ID %d", mdl.ID)

	if err := l.dbh.Create(mdl); err != nil {
		return err
	}

	go l.publishStripSaveEvent(null.NewInt(0, false), *mdl, nil)
	return nil
}

func (l *ledSvc) UpdateLEDStrip(id string, updMdl model.LedStrip) error {
	// Get model if exist
	strip, err := l.dbh.Get(id)
	if err != nil {
		return model.NewAppErr(404, err)
	}

	// profile shouldn't be updated through this endpoint
	updMdl.ProfileID = strip.ProfileID

	if err := l.dbh.Update(*strip, updMdl); err != nil {
		return model.NewAppErr(400, err)
	}
	// load profile for event
	profile, err := l.cpDbh.Get(strconv.FormatInt(updMdl.ProfileID.Int64, 10))
	if err == nil {
		go l.publishStripSaveEvent(updMdl.GetNullID(), updMdl, profile)
	}
	return nil
}

func (l *ledSvc) DeleteLEDStrip(id string) error {
	// Get model if exist
	strip, err := l.dbh.Get(id)
	if err != nil {
		return model.NewAppErr(404, err)
	}

	if err := l.dbh.Delete(strip); err != nil {
		return model.NewAppErr(400, err)
	}
	var event = model.NewStripEvent(strip.GetNullID(), model.Delete)
	go l.mh.PublishStripEvent(event)
	return nil
}

func (l *ledSvc) UpdateProfileForStrip(id string, updProf model.ColorProfile) (*model.ColorProfile, error) {
	// Get model if exist
	strip, err := l.dbh.Get(id)
	if err != nil {
		return nil, model.NewAppErr(404, err)
	}

	profile, err := l.cpDbh.Get(updProf.GetStringID())
	if err != nil {
		return nil, model.NewAppErr(404, err)
	}

	strip.ProfileID = profile.GetNullID()

	if err := l.dbh.Save(strip); err != nil {
		l.l.Error("Error: %s", err)
		return nil, model.NewAppErr(500, err)
	}

	go l.publishStripSaveEvent(strip.GetNullID(), *strip, profile)
	return profile, nil
}

func (l *ledSvc) GetProfileForStrip(id string) (*model.ColorProfile, error) {
	// Get model if exist
	strip, err := l.dbh.Get(id)
	if err != nil {
		return nil, model.NewAppErr(404, err)
	}

	if !strip.ProfileID.Valid {
		return nil, model.NewAppErr(404, errors.New("Profile not found"))
	}

	profile, err := l.cpDbh.Get(strconv.FormatInt(strip.ProfileID.Int64, 10))
	if err != nil {
		return nil, model.NewAppErr(404, err)
	}
	return profile, nil
}

func (l *ledSvc) publishStripSaveEvent(id null.Int, strip model.LedStrip, profile *model.ColorProfile) {
	var event = model.NewStripEvent(id, model.Save).With(&strip)

	if strip.ProfileID.Valid {
		if profile != nil {
			event.Strip.With(*profile)
		}
	}

	if err := l.mh.PublishStripEvent(event); err != nil {
		l.l.Error("error: %s", err.Error())
		return
	}
}

func (l *ledSvc) RemoveProfileForStrip(id string) error {
	// Get model if exist
	strip, err := l.dbh.Get(id)
	if err != nil {
		return model.NewAppErr(404, err)
	}

	strip.ProfileID.Valid = false

	if err := l.dbh.Save(strip); err != nil {
		return model.NewAppErr(500, err)
	}

	go l.publishStripSaveEvent(strip.GetNullID(), *strip, nil)
	return nil
}
