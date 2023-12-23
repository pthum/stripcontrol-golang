package service

import (
	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/internal/database"
	"github.com/pthum/stripcontrol-golang/internal/messaging"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/samber/do"
)

type CPService interface {
	GetAll() ([]model.ColorProfile, error)
	GetColorProfile(id string) (*model.ColorProfile, error)
	CreateColorProfile(mdl *model.ColorProfile) error
	UpdateColorProfile(id string, updMdl model.ColorProfile) error
	DeleteColorProfile(id string) error
}

type cpService struct {
	dbh database.DBHandler[model.ColorProfile]
	mh  messaging.EventHandler
}

func NewCPService(i *do.Injector) (CPService, error) {
	dbh := do.MustInvoke[database.DBHandler[model.ColorProfile]](i)
	mh := do.MustInvoke[messaging.EventHandler](i)
	return &cpService{
		dbh: dbh,
		mh:  mh,
	}, nil
}

func (s *cpService) GetAll() ([]model.ColorProfile, error) {
	return s.dbh.GetAll()
}

func (s *cpService) GetColorProfile(id string) (*model.ColorProfile, error) {
	return s.dbh.Get(id)
}

func (s *cpService) CreateColorProfile(mdl *model.ColorProfile) error {
	// generate an id
	mdl.GenerateID()

	return s.dbh.Create(mdl)
}

func (s *cpService) UpdateColorProfile(id string, updMdl model.ColorProfile) error {
	// Get model if exist
	profile, err := s.dbh.Get(id)
	if err != nil {
		return model.NewAppErr(404, err)
	}

	if err = s.dbh.Update(*profile, updMdl); err != nil {
		return model.NewAppErr(400, err)
	}

	var event = model.NewProfileEvent(null.NewInt(updMdl.ID, true), model.Save).With(updMdl)
	go s.mh.PublishProfileEvent(event)
	return nil
}

func (s *cpService) DeleteColorProfile(id string) error {
	// Get model if exist
	profile, err := s.dbh.Get(id)
	if err != nil {
		return model.NewAppErr(404, err)
	}
	if err := s.dbh.Delete(profile); err != nil {
		return model.NewAppErr(400, err)
	}

	var event = model.NewProfileEvent(null.NewInt(profile.ID, true), model.Delete)
	go s.mh.PublishProfileEvent(event)
	return nil
}
