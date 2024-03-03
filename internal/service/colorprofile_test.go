package service

import (
	"errors"
	"strconv"
	"sync"
	"testing"

	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/internal/database"
	dbm "github.com/pthum/stripcontrol-golang/internal/database/mocks"
	"github.com/pthum/stripcontrol-golang/internal/messaging"
	mhm "github.com/pthum/stripcontrol-golang/internal/messaging/mocks"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type cphMocks struct {
	*baseMocks
	cps *cpService
}
type baseMocks struct {
	cpDbh *dbm.DBHandler[model.ColorProfile]
	lsDbh *dbm.DBHandler[model.LedStrip]
	mh    *mhm.EventHandler
}

func TestGetAllColorProfiles(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	expRet := createDummyProfile()
	destarr := []model.ColorProfile{*expRet}
	mocks.cpDbh.
		EXPECT().
		GetAll().
		Return(destarr, nil).
		Once()

	result, err := mocks.cps.GetAll()

	assert.NoError(t, err)
	assert.Equal(t, *expRet, result[0])
}

func TestGetAllColorProfiles_GetError(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	destarr := []model.ColorProfile{}
	mocks.cpDbh.
		EXPECT().
		GetAll().
		Return(destarr, errors.New("get error")).
		Once()

	res, err := mocks.cps.GetAll()
	assert.Error(t, err)
	assert.Len(t, res, 0)
}

func TestGetColorProfile(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	retObj := createDummyProfile()
	idS := idStr(retObj.ID)
	mocks.expectDBProfileGet(retObj, nil)

	res, err := mocks.cps.GetColorProfile(idS)

	assert.NoError(t, err)
	assert.Equal(t, retObj, res)
}

func TestGetColorProfile_GetError(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	retObj := createDummyProfile()
	idS := idStr(retObj.ID)
	mocks.expectDBProfileGet(nil, errors.New("not found"))

	res, err := mocks.cps.GetColorProfile(idS)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCreateColorProfile(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	inBody := *createDummyProfile()
	input := inBody
	var newId int64
	mocks.cpDbh.
		EXPECT().
		Create(mock.Anything).
		Run(func(input *model.ColorProfile) {
			// id should have been generated
			assert.NotEqual(t, inBody.ID, input.ID)
			newId = input.ID
		}).
		Return(nil).
		Once()

	expectedObj := inBody
	err := mocks.cps.CreateColorProfile(&input)

	expectedObj.ID = newId

	assert.NoError(t, err)
	assert.Equal(t, expectedObj, input)
}

func TestCreateColorProfile_SaveError(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	inBody := *createDummyProfile()
	input := inBody
	mocks.cpDbh.
		EXPECT().
		Create(mock.Anything).
		Run(func(input *model.ColorProfile) {
			// id should have been generated
			assert.NotEqual(t, inBody.ID, input.ID)
		}).
		Return(errors.New("save failed")).
		Once()

	err := mocks.cps.CreateColorProfile(&input)
	assert.Error(t, err)
}

func TestDeleteColorProfile(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	getObj := createDummyProfile()
	mocks.expectDBProfileGet(getObj, nil)

	mocks.cpDbh.
		EXPECT().
		Delete(mock.Anything).
		Return(nil)
	wg := mocks.expectPublishProfileEvent(t, model.Delete, getObj.ID, nil)
	idS := idStr(getObj.ID)

	err := mocks.cps.DeleteColorProfile(idS)
	wg.Wait()
	assert.NoError(t, err)
}

func TestDeleteColorProfile_MissingDBProfile(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	getObj := createDummyProfile()
	mocks.expectDBProfileGet(nil, errors.New("not found"))
	idS := idStr(getObj.ID)

	err := mocks.cps.DeleteColorProfile(idS)

	assert.Error(t, err)
}

func TestDeleteColorProfile_DeleteError(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	getObj := createDummyProfile()
	mocks.expectDBProfileGet(getObj, nil)

	mocks.cpDbh.
		EXPECT().
		Delete(mock.Anything).
		Return(errors.New("delete error"))
	idS := idStr(getObj.ID)

	err := mocks.cps.DeleteColorProfile(idS)

	assert.Error(t, err)
}

func TestUpdateColorProfile(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	inBody := createDummyProfile()
	dbO := *createProfile(105, 100, 100, 100, 2)
	mocks.expectDBProfileGet(&dbO, nil)

	mocks.cpDbh.
		EXPECT().
		Update(dbO, *inBody).
		Return(nil)
	wg := mocks.expectPublishProfileEvent(t, model.Save, inBody.ID, inBody)
	idS := idStr(dbO.ID)

	err := mocks.cps.UpdateColorProfile(idS, *inBody)
	wg.Wait()
	assert.NoError(t, err)
}

func TestUpdateColorProfile_MissingDBProfile(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	inBody := createDummyProfile()
	dbO := *createProfile(105, 100, 100, 100, 2)
	mocks.expectDBProfileGet(nil, errors.New("not found"))
	idS := idStr(dbO.ID)

	err := mocks.cps.UpdateColorProfile(idS, *inBody)

	assert.Error(t, err)
}

func TestUpdateColorProfile_UpdateError(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	inBody := createDummyProfile()
	dbO := *createProfile(105, 100, 100, 100, 2)
	mocks.expectDBProfileGet(&dbO, nil)

	mocks.cpDbh.
		EXPECT().
		Update(dbO, *inBody).
		Return(errors.New("update error"))
	idS := idStr(dbO.ID)

	err := mocks.cps.UpdateColorProfile(idS, *inBody)

	assert.Error(t, err)
}

func (chm *cphMocks) expectPublishProfileEvent(t *testing.T, typ model.EventType, id int64, body *model.ColorProfile) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	chm.mh.
		EXPECT().
		PublishProfileEvent(mock.Anything).
		Run(func(event *model.ProfileEvent) {
			assert.Equal(t, typ, event.Type)
			assert.Equal(t, id, event.ID.Int64)
			if body != nil {
				assert.Equal(t, *body, event.State.Profile)
			}
			assert.Equal(t, body != nil, event.State.Valid)
			wg.Done()
		}).
		Return(nil)
	return &wg
}

func createDummyProfile() *model.ColorProfile {
	return createProfile(185, 123, 234, 12, 1)
}

func createProfile(id, red, green, blue, brightness int64) *model.ColorProfile {
	return &model.ColorProfile{
		BaseModel:  model.BaseModel{ID: id},
		Red:        null.IntFrom(red),
		Green:      null.IntFrom(green),
		Blue:       null.IntFrom(blue),
		Brightness: null.IntFrom(brightness),
	}
}

func createCPHandlerMocks(t *testing.T) *cphMocks {
	i := do.New()
	bm := createBaseMocks(i, t)
	cps, err := NewCPService(i)
	assert.NoError(t, err)
	return &cphMocks{
		baseMocks: bm,
		cps:       cps.(*cpService),
	}
}

func createBaseMocks(i *do.Injector, t *testing.T) *baseMocks {
	cpDbh := dbm.NewDBHandler[model.ColorProfile](t)
	lsDbh := dbm.NewDBHandler[model.LedStrip](t)
	do.ProvideValue[database.DBHandler[model.ColorProfile]](i, cpDbh)
	do.ProvideValue[database.DBHandler[model.LedStrip]](i, lsDbh)
	mh := mhm.NewEventHandler(t)
	do.ProvideValue[messaging.EventHandler](i, mh)
	return &baseMocks{
		cpDbh: cpDbh,
		lsDbh: lsDbh,
		mh:    mh,
	}
}

func (bm *baseMocks) expectDBProfileGet(getStrip *model.ColorProfile, getError error) {
	getStripIdStr := mock.Anything
	if getStrip != nil {
		getStripIdStr = idStr(getStrip.ID)
	}
	bm.cpDbh.
		EXPECT().
		Get(getStripIdStr).
		Return(getStrip, getError).
		Once()
}
func idStr(id int64) string {
	return strconv.FormatInt(id, 10)
}
