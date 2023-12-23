package service

import (
	"errors"
	"testing"
	"time"

	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type lsMocks struct {
	*baseMocks
	lh *ledSvc
}

func TestGetAllLEDStrips(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	destarr := []model.LedStrip{*createValidDummyStrip()}
	mocks.lsDbh.
		EXPECT().
		GetAll().
		Return(destarr, nil).
		Once()

	res, err := mocks.lh.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, destarr[0], res[0])
}

func TestGetAllLEDStrips_Error(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	destarr := []model.LedStrip{}
	mocks.lsDbh.
		EXPECT().
		GetAll().
		Return(destarr, errors.New("get error")).
		Once()
	res, err := mocks.lh.GetAll()

	assert.Error(t, err)
	assert.Equal(t, destarr, res)
}

func TestGetLEDStrip(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	retObj := createValidDummyStrip()
	reqId := idStr(retObj.ID)
	mocks.expectDBStripGet(retObj, nil)

	res, err := mocks.lh.GetLEDStrip(reqId)
	assert.NoError(t, err)
	assert.Equal(t, retObj, res)
}

func TestGetLEDStrip_Error(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	reqId := "6000"
	mocks.expectDBStripGet(nil, errors.New("nothing found"))

	res, err := mocks.lh.GetLEDStrip(reqId)
	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestCreateLEDStrip(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	reqObj := createValidDummyStrip()
	origReqObj := *reqObj
	var newId int64
	mocks.lsDbh.
		EXPECT().
		Create(mock.Anything).
		Run(func(input *model.LedStrip) {
			// id should have been generated
			assert.NotEqual(t, origReqObj.ID, input.ID)
			newId = input.ID
		}).
		Return(nil).
		Once()
	mocks.expectPublishStripEvent(t, model.Save, newId, true, false, nil)

	err := mocks.lh.CreateLEDStrip(reqObj)

	// small sleep to have the async routines run
	time.Sleep(50 * time.Millisecond)

	assert.NoError(t, err)
}

func TestCreateLEDStrip_SaveError(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	reqObj := createValidDummyStrip()
	mocks.lsDbh.
		EXPECT().
		Create(mock.Anything).
		Return(errors.New("save failed")).
		Once()
	err := mocks.lh.CreateLEDStrip(reqObj)
	assert.Error(t, err)
}

func TestCreateLEDStrip_PublishError(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	reqObj := createValidDummyStrip()
	origReqObj := *reqObj
	var newId int64
	mocks.lsDbh.
		EXPECT().
		Create(mock.Anything).
		Run(func(input *model.LedStrip) {
			// id should have been generated
			assert.NotEqual(t, origReqObj.ID, input.ID)
			newId = input.ID
		}).
		Return(nil).
		Once()
	mocks.expectPublishStripEvent(t, model.Save, newId, true, false, errors.New("publish failed"))

	err := mocks.lh.CreateLEDStrip(reqObj)

	// small sleep to have the async routines run
	time.Sleep(50 * time.Millisecond)
	assert.NoError(t, err)
}

func TestDeleteLEDStrip(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	getObj := createValidDummyStrip()
	mocks.expectDBStripGet(getObj, nil)
	mocks.lsDbh.
		EXPECT().
		Delete(mock.Anything).
		Return(nil)
	mocks.expectPublishStripEvent(t, model.Delete, getObj.ID, false, false, nil)

	err := mocks.lh.DeleteLEDStrip("185")

	// small sleep to have the async routines run
	time.Sleep(50 * time.Millisecond)
	assert.NoError(t, err)
}

func TestDeleteLEDStrip_MissingDBStrip(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	mocks.expectDBStripGet(nil, errors.New("not found"))

	err := mocks.lh.DeleteLEDStrip("185")
	assert.Error(t, err)
}

func TestDeleteLEDStrip_DeleteError(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	getObj := createValidDummyStrip()
	mocks.expectDBStripGet(getObj, nil)
	mocks.lsDbh.
		EXPECT().
		Delete(mock.Anything).
		Return(errors.New("delete error"))

	err := mocks.lh.DeleteLEDStrip("185")
	assert.Error(t, err)
}

func TestDeleteLEDStrip_PublishError(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	getObj := createValidDummyStrip()
	mocks.expectDBStripGet(getObj, nil)
	mocks.lsDbh.
		EXPECT().
		Delete(mock.Anything).
		Return(nil)
	mocks.expectPublishStripEvent(t, model.Delete, getObj.ID, false, false, errors.New("publish failed"))

	err := mocks.lh.DeleteLEDStrip("185")

	// small sleep to have the async routines run
	time.Sleep(50 * time.Millisecond)
	assert.NoError(t, err)
}

func TestUpdateLEDStrip(t *testing.T) {
	dbObj := model.LedStrip{
		BaseModel:   model.BaseModel{ID: 185},
		Description: "TestFromDb",
		Enabled:     false,
		MisoPin:     null.IntFrom(100),
		Name:        "TestFromDb",
		NumLeds:     null.IntFrom(99),
		SclkPin:     null.IntFrom(99),
		SpeedHz:     null.IntFrom(80001),
		ProfileID:   null.IntFrom(15),
	}
	inputObj := createValidDummyStrip()
	inputObj.ProfileID = null.IntFrom(15)
	fakeProfile := createDummyProfile()
	fakeProfile.ID = inputObj.ProfileID.Int64

	mocks := createLEDHandlerMocks(t)
	mocks.expectDBStripGet(&dbObj, nil)
	mocks.lsDbh.
		EXPECT().
		Update(dbObj, *inputObj).
		Return(nil)
	mocks.expectPublishStripEvent(t, model.Save, inputObj.ID, true, true, nil)
	mocks.expectDBProfileGet(fakeProfile, nil)

	err := mocks.lh.UpdateLEDStrip("185", *inputObj)

	// small sleep to have the async routines run
	time.Sleep(50 * time.Millisecond)
	assert.NoError(t, err)
}

func TestUpdateLEDStrip_MissingDBProfile(t *testing.T) {
	inputObj := createValidDummyStrip()
	inputObj.ProfileID = null.IntFrom(15)
	mocks := createLEDHandlerMocks(t)
	mocks.expectDBStripGet(nil, errors.New("not found"))

	err := mocks.lh.UpdateLEDStrip("185", *inputObj)
	assert.Error(t, err)
}

func TestUpdateLEDStrip_UpdateError(t *testing.T) {
	dbObj := model.LedStrip{
		BaseModel:   model.BaseModel{ID: 185},
		Description: "TestFromDb",
		Enabled:     false,
		MisoPin:     null.IntFrom(100),
		Name:        "TestFromDb",
		NumLeds:     null.IntFrom(99),
		SclkPin:     null.IntFrom(99),
		SpeedHz:     null.IntFrom(80001),
		ProfileID:   null.IntFrom(15),
	}
	inputObj := createValidDummyStrip()
	inputObj.ProfileID = null.IntFrom(15)

	mocks := createLEDHandlerMocks(t)
	mocks.expectDBStripGet(&dbObj, nil)
	mocks.lsDbh.
		EXPECT().
		Update(dbObj, *inputObj).
		Return(errors.New("update failed"))

	err := mocks.lh.UpdateLEDStrip("185", *inputObj)
	assert.Error(t, err)
}

func TestUpdateLEDStrip_PublishError(t *testing.T) {
	dbObj := model.LedStrip{
		BaseModel:   model.BaseModel{ID: 185},
		Description: "TestFromDb",
		Enabled:     false,
		MisoPin:     null.IntFrom(100),
		Name:        "TestFromDb",
		NumLeds:     null.IntFrom(99),
		SclkPin:     null.IntFrom(99),
		SpeedHz:     null.IntFrom(80001),
		ProfileID:   null.IntFrom(15),
	}
	inputObj := createValidDummyStrip()
	inputObj.ProfileID = null.IntFrom(15)
	fakeProfile := createDummyProfile()
	fakeProfile.ID = inputObj.ProfileID.Int64
	mocks := createLEDHandlerMocks(t)
	mocks.expectDBStripGet(&dbObj, nil)
	mocks.lsDbh.
		EXPECT().
		Update(dbObj, *inputObj).
		Return(nil)
	mocks.expectPublishStripEvent(t, model.Save, inputObj.ID, true, true, errors.New("publish error"))
	mocks.expectDBProfileGet(fakeProfile, nil)

	err := mocks.lh.UpdateLEDStrip("185", *inputObj)

	// small sleep to have the async routines run
	time.Sleep(50 * time.Millisecond)
	assert.NoError(t, err)
}

func TestGetProfileForLEDStrip(t *testing.T) {
	fakeProfile := *createDummyProfile()
	returnObj := createValidDummyStrip()
	returnObj.ProfileID = null.IntFrom(fakeProfile.ID)
	mocks := createLEDHandlerMocks(t)
	stripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(returnObj, nil)
	mocks.expectDBProfileGet(&fakeProfile, nil)

	result, err := mocks.lh.GetProfileForStrip(stripIdStr)
	assert.NoError(t, err)
	assert.Equal(t, fakeProfile, *result)
}

func TestGetProfileForLEDStrip_MissingStrip(t *testing.T) {
	returnObj := createValidDummyStrip()
	mocks := createLEDHandlerMocks(t)
	stripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(returnObj, errors.New("strip not found"))

	res, err := mocks.lh.GetProfileForStrip(stripIdStr)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestGetProfileForLEDStrip_NoProfile(t *testing.T) {
	returnObj := createValidDummyStrip()
	returnObj.ProfileID = null.NewInt(0, false)
	mocks := createLEDHandlerMocks(t)
	stripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(returnObj, nil)

	res, err := mocks.lh.GetProfileForStrip(stripIdStr)

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestGetProfileForLEDStrip_MissingProfile(t *testing.T) {
	fakeProfile := createDummyProfile()
	returnObj := createValidDummyStrip()
	returnObj.ProfileID = null.IntFrom(fakeProfile.ID)
	mocks := createLEDHandlerMocks(t)
	stripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(returnObj, nil)
	mocks.expectDBProfileGet(nil, errors.New("not found"))

	res, err := mocks.lh.GetProfileForStrip(stripIdStr)

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestUpdateProfileForLEDStrip(t *testing.T) {
	returnObj := model.LedStrip{
		BaseModel:   model.BaseModel{ID: 185},
		Description: "Test",
		Enabled:     false,
		MisoPin:     null.IntFrom(12),
		Name:        "Test",
		NumLeds:     null.IntFrom(5),
		SclkPin:     null.IntFrom(13),
		SpeedHz:     null.IntFrom(80000),
		ProfileID:   null.IntFrom(15),
	}
	updateProfile := model.ColorProfile{
		BaseModel:  model.BaseModel{ID: 16},
		Red:        null.IntFrom(123),
		Green:      null.IntFrom(123),
		Blue:       null.IntFrom(123),
		Brightness: null.IntFrom(2),
	}
	mocks := createLEDHandlerMocks(t)
	getStripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(&returnObj, nil)
	mocks.expectDBProfileGet(&updateProfile, nil)
	mocks.expectDBStripSave(nil)
	mocks.expectPublishStripEvent(t, model.Save, returnObj.ID, true, true, nil)

	res, err := mocks.lh.UpdateProfileForStrip(getStripIdStr, updateProfile)
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, updateProfile, *res)
	assert.NoError(t, err)
}

func TestUpdateProfileForLEDStrip_MissingStrip(t *testing.T) {
	returnObj := model.LedStrip{
		BaseModel:   model.BaseModel{ID: 185},
		Description: "Test",
		Enabled:     false,
		MisoPin:     null.IntFrom(12),
		Name:        "Test",
		NumLeds:     null.IntFrom(5),
		SclkPin:     null.IntFrom(13),
		SpeedHz:     null.IntFrom(80000),
		ProfileID:   null.IntFrom(15),
	}

	updateProfile := model.ColorProfile{
		BaseModel:  model.BaseModel{ID: 16},
		Red:        null.IntFrom(123),
		Green:      null.IntFrom(123),
		Blue:       null.IntFrom(123),
		Brightness: null.IntFrom(2),
	}

	mocks := createLEDHandlerMocks(t)
	getStripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(nil, errors.New("strip not found"))

	res, err := mocks.lh.UpdateProfileForStrip(getStripIdStr, updateProfile)

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestUpdateProfileForLEDStrip_MissingProfile(t *testing.T) {
	returnObj := model.LedStrip{
		BaseModel:   model.BaseModel{ID: 185},
		Description: "Test",
		Enabled:     false,
		MisoPin:     null.IntFrom(12),
		Name:        "Test",
		NumLeds:     null.IntFrom(5),
		SclkPin:     null.IntFrom(13),
		SpeedHz:     null.IntFrom(80000),
		ProfileID:   null.IntFrom(15),
	}
	updateProfile := model.ColorProfile{
		BaseModel:  model.BaseModel{ID: 16},
		Red:        null.IntFrom(123),
		Green:      null.IntFrom(123),
		Blue:       null.IntFrom(123),
		Brightness: null.IntFrom(2),
	}
	mocks := createLEDHandlerMocks(t)
	getStripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(&returnObj, nil)
	mocks.expectDBProfileGet(nil, errors.New("missing profile"))

	res, err := mocks.lh.UpdateProfileForStrip(getStripIdStr, updateProfile)

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestUpdateProfileForLEDStrip_SaveError(t *testing.T) {
	returnObj := model.LedStrip{
		BaseModel:   model.BaseModel{ID: 185},
		Description: "Test",
		Enabled:     false,
		MisoPin:     null.IntFrom(12),
		Name:        "Test",
		NumLeds:     null.IntFrom(5),
		SclkPin:     null.IntFrom(13),
		SpeedHz:     null.IntFrom(80000),
		ProfileID:   null.IntFrom(15),
	}
	updateProfile := model.ColorProfile{
		BaseModel:  model.BaseModel{ID: 16},
		Red:        null.IntFrom(123),
		Green:      null.IntFrom(123),
		Blue:       null.IntFrom(123),
		Brightness: null.IntFrom(2),
	}
	mocks := createLEDHandlerMocks(t)
	getStripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(&returnObj, nil)
	mocks.expectDBProfileGet(&updateProfile, nil)
	mocks.expectDBStripSave(errors.New("save failed"))

	res, err := mocks.lh.UpdateProfileForStrip(getStripIdStr, updateProfile)
	time.Sleep(50 * time.Millisecond)

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestUpdateProfileForLEDStrip_PublishError(t *testing.T) {
	returnObj := model.LedStrip{
		BaseModel:   model.BaseModel{ID: 185},
		Description: "Test",
		Enabled:     false,
		MisoPin:     null.IntFrom(12),
		Name:        "Test",
		NumLeds:     null.IntFrom(5),
		SclkPin:     null.IntFrom(13),
		SpeedHz:     null.IntFrom(80000),
		ProfileID:   null.IntFrom(15),
	}
	updateProfile := model.ColorProfile{
		BaseModel:  model.BaseModel{ID: 16},
		Red:        null.IntFrom(123),
		Green:      null.IntFrom(123),
		Blue:       null.IntFrom(123),
		Brightness: null.IntFrom(2),
	}
	mocks := createLEDHandlerMocks(t)
	getStripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(&returnObj, nil)
	mocks.expectDBProfileGet(&updateProfile, nil)
	mocks.expectDBStripSave(nil)
	mocks.expectPublishStripEvent(t, model.Save, returnObj.ID, true, true, errors.New("publish error"))

	res, err := mocks.lh.UpdateProfileForStrip(getStripIdStr, updateProfile)
	time.Sleep(50 * time.Millisecond)

	assert.NotNil(t, res)
	assert.Equal(t, updateProfile, *res)
	assert.NoError(t, err)
}

func TestRemoveProfileForLEDStrip(t *testing.T) {
	getStrip := createValidDummyStrip()
	mocks := createLEDHandlerMocks(t)
	getStripIdStr := idStr(getStrip.ID)
	mocks.expectDBStripGet(getStrip, nil)

	mocks.expectDBStripSave(nil)
	mocks.expectPublishStripEvent(t, model.Save, getStrip.ID, true, false, nil)

	err := mocks.lh.RemoveProfileForStrip(getStripIdStr)
	time.Sleep(50 * time.Millisecond)

	assert.NoError(t, err)
}

func TestRemoveProfileForLEDStrip_MissingStrip(t *testing.T) {
	getStrip := createValidDummyStrip()
	mocks := createLEDHandlerMocks(t)
	getStripIdStr := idStr(getStrip.ID)
	mocks.expectDBStripGet(nil, errors.New("not found"))

	err := mocks.lh.RemoveProfileForStrip(getStripIdStr)

	assert.Error(t, err)
}

func TestRemoveProfileForLEDStrip_SaveError(t *testing.T) {
	getStrip := createValidDummyStrip()
	mocks := createLEDHandlerMocks(t)
	getStripIdStr := idStr(getStrip.ID)
	mocks.expectDBStripGet(getStrip, nil)
	mocks.expectDBStripSave(errors.New("save error"))

	err := mocks.lh.RemoveProfileForStrip(getStripIdStr)

	assert.Error(t, err)
}

func TestRemoveProfileForLEDStrip_PublishError(t *testing.T) {
	getStrip := createValidDummyStrip()
	mocks := createLEDHandlerMocks(t)
	getStripIdStr := idStr(getStrip.ID)
	mocks.expectDBStripGet(getStrip, nil)

	mocks.expectDBStripSave(nil)
	mocks.expectPublishStripEvent(t, model.Save, getStrip.ID, true, false, errors.New("publish error"))

	err := mocks.lh.RemoveProfileForStrip(getStripIdStr)
	time.Sleep(50 * time.Millisecond)

	assert.NoError(t, err)
}

func (lhm *lsMocks) expectDBStripGet(getStrip *model.LedStrip, getStripError error) {
	getStripIdStr := mock.Anything
	if getStrip != nil {
		getStripIdStr = idStr(getStrip.ID)
	}
	lhm.lsDbh.
		EXPECT().
		Get(getStripIdStr).
		Return(getStrip, getStripError).
		Once()
}

func (lhm *lsMocks) expectDBStripSave(saveError error) {
	lhm.lsDbh.
		EXPECT().
		Save(mock.Anything).
		Return(saveError)
}

func (lhm *lsMocks) expectPublishStripEvent(t *testing.T, typ model.EventType, id int64, valid bool, expectProfile bool, publishError error) {
	lhm.mh.
		EXPECT().
		PublishStripEvent(mock.Anything).
		Run(func(event *model.StripEvent) {
			assert.Equal(t, typ, event.Type)
			assert.Equal(t, id, event.ID.Int64)
			assert.Equal(t, valid, event.Strip.Valid)
			assert.Equal(t, expectProfile, event.Strip.Strip.Profile.Valid)
		}).
		Return(publishError)
}

func createValidDummyStrip() *model.LedStrip {
	return &model.LedStrip{
		BaseModel:   model.BaseModel{ID: 185},
		Description: "Test",
		Enabled:     false,
		MisoPin:     null.IntFrom(12),
		Name:        "Test",
		NumLeds:     null.IntFrom(5),
		SclkPin:     null.IntFrom(13),
		SpeedHz:     null.IntFrom(80000),
	}
}

func createLEDHandlerMocks(t *testing.T) *lsMocks {
	i := do.New()
	bm := createBaseMocks(i, t)
	lh, err := NewLEDService(i)
	assert.NoError(t, err)
	return &lsMocks{
		baseMocks: bm,
		lh:        lh.(*ledSvc),
	}
}
