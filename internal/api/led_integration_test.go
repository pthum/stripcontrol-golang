package api

import (
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/pthum/stripcontrol-golang/internal/service"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type lhITMocks struct {
	*baseMocks
	lh *ledHandlerImpl
}

func TestGetAllLEDStripsIT(t *testing.T) {
	mocks := createLEDHandlerITMocks(t)
	destarr := []model.LedStrip{*createValidDummyStrip()}
	mocks.lsDbh.
		EXPECT().
		GetAll().
		Return(destarr, nil).
		Once()

	req, w := prepareHttpTest(http.MethodGet, ledstripPath, nil, nil)

	mocks.lh.GetAllLedStrips(w, req)
	res := w.Result()
	defer res.Body.Close()

	var result []model.LedStrip
	bodyToObj(t, res, &result)

	assert.Equal(t, destarr[0], result[0])
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetAllLEDStripsIT_Error(t *testing.T) {
	mocks := createLEDHandlerITMocks(t)
	destarr := []model.LedStrip{}
	mocks.lsDbh.
		EXPECT().
		GetAll().
		Return(destarr, errors.New("get error")).
		Once()

	req, w := prepareHttpTest(http.MethodGet, ledstripPath, nil, nil)

	mocks.lh.GetAllLedStrips(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestGetLEDStripIT(t *testing.T) {
	mocks := createLEDHandlerITMocks(t)
	retObj := createValidDummyStrip()
	reqId := idStr(retObj.ID)
	mocks.expectDBStripGet(retObj, nil)
	req, w := prepareHttpTest(http.MethodGet, ledstripIDPath, uv{"id": reqId}, nil)

	mocks.lh.GetLedStrip(w, req)
	res := w.Result()
	defer res.Body.Close()

	var result model.LedStrip
	bodyToObj(t, res, &result)

	assert.Equal(t, *retObj, result)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetLEDStripIT_Error(t *testing.T) {
	mocks := createLEDHandlerITMocks(t)
	reqId := "6000"
	mocks.expectDBStripGet(nil, errors.New("nothing found"))
	req, w := prepareHttpTest(http.MethodGet, ledstripIDPath, uv{"id": reqId}, nil)

	mocks.lh.GetLedStrip(w, req)
	res := w.Result()
	defer res.Body.Close()

	var result model.LedStrip
	bodyToObj(t, res, &result)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestCreateLEDStripIT(t *testing.T) {
	mocks := createLEDHandlerITMocks(t)
	reqObj := createValidDummyStrip()
	var newId int64
	mocks.lsDbh.
		EXPECT().
		Create(mock.Anything).
		Run(func(input *model.LedStrip) {
			// id should have been generated
			assert.NotEqual(t, reqObj.ID, input.ID)
			newId = input.ID
		}).
		Return(nil).
		Once()
	body := objToReader(t, reqObj)
	mocks.expectPublishStripEvent(t, model.Save, newId, true, false, nil)
	req, w := prepareHttpTest(http.MethodPost, ledstripPath, nil, body)

	mocks.lh.CreateLedStrip(w, req)

	// small sleep to have the async routines run
	time.Sleep(50 * time.Millisecond)

	res := w.Result()
	defer res.Body.Close()

	expectedObj := reqObj
	expectedObj.ID = newId
	var result model.LedStrip
	bodyToObj(t, res, &result)

	assert.Equal(t, *expectedObj, result)
	assert.Contains(t, res.Header["Location"][0], idStr(newId))
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestCreateLEDStripIT_SaveError(t *testing.T) {
	mocks := createLEDHandlerITMocks(t)
	reqObj := createValidDummyStrip()
	mocks.lsDbh.
		EXPECT().
		Create(mock.Anything).
		Return(errors.New("save failed")).
		Once()
	body := objToReader(t, reqObj)
	req, w := prepareHttpTest(http.MethodPost, ledstripPath, nil, body)

	mocks.lh.CreateLedStrip(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestCreateLEDStripIT_PublishError(t *testing.T) {
	mocks := createLEDHandlerITMocks(t)
	reqObj := createValidDummyStrip()
	var newId int64
	mocks.lsDbh.
		EXPECT().
		Create(mock.Anything).
		Run(func(input *model.LedStrip) {
			// id should have been generated
			assert.NotEqual(t, reqObj.ID, input.ID)
			newId = input.ID
		}).
		Return(nil).
		Once()
	body := objToReader(t, reqObj)
	mocks.expectPublishStripEvent(t, model.Save, newId, true, false, errors.New("publish failed"))
	req, w := prepareHttpTest(http.MethodPost, ledstripPath, nil, body)

	mocks.lh.CreateLedStrip(w, req)

	// small sleep to have the async routines run
	time.Sleep(50 * time.Millisecond)

	res := w.Result()
	defer res.Body.Close()

	expectedObj := reqObj
	expectedObj.ID = newId
	var result model.LedStrip
	bodyToObj(t, res, &result)

	assert.Equal(t, *expectedObj, result)
	assert.Contains(t, res.Header["Location"][0], idStr(newId))
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestDeleteLEDStripIT(t *testing.T) {
	mocks := createLEDHandlerITMocks(t)
	getObj := createValidDummyStrip()
	mocks.expectDBStripGet(getObj, nil)
	mocks.lsDbh.
		EXPECT().
		Delete(mock.Anything).
		Return(nil)
	mocks.expectPublishStripEvent(t, model.Delete, getObj.ID, false, false, nil)
	req, w := prepareHttpTest(http.MethodDelete, ledstripIDPath, uv{"id": "185"}, nil)

	mocks.lh.DeleteLedStrip(w, req)

	// small sleep to have the async routines run
	time.Sleep(50 * time.Millisecond)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func TestDeleteLEDStripIT_MissingDBStrip(t *testing.T) {
	mocks := createLEDHandlerITMocks(t)
	mocks.expectDBStripGet(nil, errors.New("not found"))
	req, w := prepareHttpTest(http.MethodDelete, ledstripIDPath, uv{"id": "185"}, nil)

	mocks.lh.DeleteLedStrip(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestDeleteLEDStripIT_DeleteError(t *testing.T) {
	mocks := createLEDHandlerITMocks(t)
	getObj := createValidDummyStrip()
	mocks.expectDBStripGet(getObj, nil)
	mocks.lsDbh.
		EXPECT().
		Delete(mock.Anything).
		Return(errors.New("delete error"))
	req, w := prepareHttpTest(http.MethodDelete, ledstripIDPath, uv{"id": "185"}, nil)

	mocks.lh.DeleteLedStrip(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestDeleteLEDStripIT_PublishError(t *testing.T) {
	mocks := createLEDHandlerITMocks(t)
	getObj := createValidDummyStrip()
	mocks.expectDBStripGet(getObj, nil)
	mocks.lsDbh.
		EXPECT().
		Delete(mock.Anything).
		Return(nil)
	mocks.expectPublishStripEvent(t, model.Delete, getObj.ID, false, false, errors.New("publish failed"))
	req, w := prepareHttpTest(http.MethodDelete, ledstripIDPath, uv{"id": "185"}, nil)

	mocks.lh.DeleteLedStrip(w, req)

	// small sleep to have the async routines run
	time.Sleep(50 * time.Millisecond)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func TestUpdateLEDStripIT(t *testing.T) {
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

	mocks := createLEDHandlerITMocks(t)
	mocks.expectDBStripGet(&dbObj, nil)
	body := objToReader(t, inputObj)
	mocks.lsDbh.
		EXPECT().
		Update(dbObj, *inputObj).
		Return(nil)
	mocks.expectPublishStripEvent(t, model.Save, inputObj.ID, true, true, nil)
	mocks.expectDBProfileGet(fakeProfile, nil)

	req, w := prepareHttpTest(http.MethodPut, ledstripIDPath, uv{"id": "185"}, body)

	mocks.lh.UpdateLedStrip(w, req)

	// small sleep to have the async routines run
	time.Sleep(50 * time.Millisecond)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestUpdateLEDStripIT_MissingDBProfile(t *testing.T) {
	inputObj := createValidDummyStrip()
	inputObj.ProfileID = null.IntFrom(15)
	mocks := createLEDHandlerITMocks(t)
	mocks.expectDBStripGet(nil, errors.New("not found"))
	body := objToReader(t, inputObj)
	req, w := prepareHttpTest(http.MethodPut, ledstripIDPath, uv{"id": "185"}, body)

	mocks.lh.UpdateLedStrip(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestUpdateLEDStripIT_MissingBody(t *testing.T) {
	mocks := createLEDHandlerITMocks(t)
	var body io.Reader
	req, w := prepareHttpTest(http.MethodPut, ledstripIDPath, uv{"id": "185"}, body)

	mocks.lh.UpdateLedStrip(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUpdateLEDStripIT_UpdateError(t *testing.T) {
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

	mocks := createLEDHandlerITMocks(t)
	mocks.expectDBStripGet(&dbObj, nil)
	body := objToReader(t, inputObj)
	mocks.lsDbh.
		EXPECT().
		Update(dbObj, *inputObj).
		Return(errors.New("update failed"))

	req, w := prepareHttpTest(http.MethodPut, ledstripIDPath, uv{"id": "185"}, body)

	mocks.lh.UpdateLedStrip(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUpdateLEDStripIT_PublishError(t *testing.T) {
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
	mocks := createLEDHandlerITMocks(t)
	mocks.expectDBStripGet(&dbObj, nil)
	body := objToReader(t, inputObj)
	mocks.lsDbh.
		EXPECT().
		Update(dbObj, *inputObj).
		Return(nil)
	mocks.expectPublishStripEvent(t, model.Save, inputObj.ID, true, true, errors.New("publish error"))
	mocks.expectDBProfileGet(fakeProfile, nil)

	req, w := prepareHttpTest(http.MethodPut, ledstripIDPath, uv{"id": "185"}, body)

	mocks.lh.UpdateLedStrip(w, req)

	// small sleep to have the async routines run
	time.Sleep(50 * time.Millisecond)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetProfileForLEDStripIT(t *testing.T) {
	fakeProfile := createDummyProfile()
	returnObj := createValidDummyStrip()
	returnObj.ProfileID = null.IntFrom(fakeProfile.ID)
	mocks := createLEDHandlerITMocks(t)
	stripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(returnObj, nil)
	mocks.expectDBProfileGet(fakeProfile, nil)

	req, w := prepareHttpTest(http.MethodGet, ledstripIDProfilePath, uv{"id": stripIdStr}, nil)

	mocks.lh.GetProfileForStrip(w, req)
	res := w.Result()
	defer res.Body.Close()

	var result model.ColorProfile
	bodyToObj(t, res, &result)
	assert.Equal(t, *fakeProfile, result)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetProfileForLEDStripIT_MissingStrip(t *testing.T) {
	returnObj := createValidDummyStrip()
	mocks := createLEDHandlerITMocks(t)
	stripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(returnObj, errors.New("strip not found"))

	req, w := prepareHttpTest(http.MethodGet, ledstripIDProfilePath, uv{"id": stripIdStr}, nil)

	mocks.lh.GetProfileForStrip(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestGetProfileForLEDStripIT_NoProfile(t *testing.T) {
	returnObj := createValidDummyStrip()
	returnObj.ProfileID = null.NewInt(0, false)
	mocks := createLEDHandlerITMocks(t)
	stripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(returnObj, nil)

	req, w := prepareHttpTest(http.MethodGet, ledstripIDProfilePath, uv{"id": stripIdStr}, nil)

	mocks.lh.GetProfileForStrip(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestGetProfileForLEDStripIT_MissingProfile(t *testing.T) {
	fakeProfile := createDummyProfile()
	returnObj := createValidDummyStrip()
	returnObj.ProfileID = null.IntFrom(fakeProfile.ID)
	mocks := createLEDHandlerITMocks(t)
	stripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(returnObj, nil)
	mocks.expectDBProfileGet(nil, errors.New("not found"))

	req, w := prepareHttpTest(http.MethodGet, ledstripIDProfilePath, uv{"id": stripIdStr}, nil)

	mocks.lh.GetProfileForStrip(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestUpdateProfileForLEDStripIT(t *testing.T) {
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
	mocks := createLEDHandlerITMocks(t)
	getStripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(&returnObj, nil)
	mocks.expectDBProfileGet(&updateProfile, nil)
	mocks.expectDBStripSave(nil)
	mocks.expectPublishStripEvent(t, model.Save, returnObj.ID, true, true, nil)
	body := objToReader(t, updateProfile)
	req, w := prepareHttpTest(http.MethodPut, ledstripIDProfilePath, uv{"id": getStripIdStr}, body)

	mocks.lh.UpdateProfileForStrip(w, req)
	time.Sleep(50 * time.Millisecond)
	res := w.Result()
	defer res.Body.Close()

	var result model.ColorProfile
	bodyToObj(t, res, &result)
	assert.Equal(t, updateProfile, result)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestUpdateProfileForLEDStripIT_MissingStrip(t *testing.T) {
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

	mocks := createLEDHandlerITMocks(t)
	getStripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(nil, errors.New("strip not found"))
	body := objToReader(t, updateProfile)
	req, w := prepareHttpTest(http.MethodPut, ledstripIDProfilePath, uv{"id": getStripIdStr}, body)

	mocks.lh.UpdateProfileForStrip(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestUpdateProfileForLEDStripIT_MissingBody(t *testing.T) {
	mocks := createLEDHandlerITMocks(t)
	getStripIdStr := "6000"
	var body io.Reader
	req, w := prepareHttpTest(http.MethodPut, ledstripIDProfilePath, uv{"id": getStripIdStr}, body)

	mocks.lh.UpdateProfileForStrip(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUpdateProfileForLEDStripIT_MissingProfile(t *testing.T) {
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
	mocks := createLEDHandlerITMocks(t)
	getStripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(&returnObj, nil)
	mocks.expectDBProfileGet(nil, errors.New("missing profile"))
	body := objToReader(t, updateProfile)
	req, w := prepareHttpTest(http.MethodPut, ledstripIDProfilePath, uv{"id": getStripIdStr}, body)

	mocks.lh.UpdateProfileForStrip(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestUpdateProfileForLEDStripIT_SaveError(t *testing.T) {
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
	mocks := createLEDHandlerITMocks(t)
	getStripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(&returnObj, nil)
	mocks.expectDBProfileGet(&updateProfile, nil)
	mocks.expectDBStripSave(errors.New("save failed"))
	body := objToReader(t, updateProfile)
	req, w := prepareHttpTest(http.MethodPut, ledstripIDProfilePath, uv{"id": getStripIdStr}, body)

	mocks.lh.UpdateProfileForStrip(w, req)
	time.Sleep(50 * time.Millisecond)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestUpdateProfileForLEDStripIT_PublishError(t *testing.T) {
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
	mocks := createLEDHandlerITMocks(t)
	getStripIdStr := idStr(returnObj.ID)
	mocks.expectDBStripGet(&returnObj, nil)
	mocks.expectDBProfileGet(&updateProfile, nil)
	mocks.expectDBStripSave(nil)
	mocks.expectPublishStripEvent(t, model.Save, returnObj.ID, true, true, errors.New("publish error"))
	body := objToReader(t, updateProfile)
	req, w := prepareHttpTest(http.MethodPut, ledstripIDProfilePath, uv{"id": getStripIdStr}, body)

	mocks.lh.UpdateProfileForStrip(w, req)
	time.Sleep(50 * time.Millisecond)
	res := w.Result()
	defer res.Body.Close()

	var result model.ColorProfile
	bodyToObj(t, res, &result)
	assert.Equal(t, updateProfile, result)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestRemoveProfileForLEDStripIT(t *testing.T) {
	getStrip := createValidDummyStrip()
	mocks := createLEDHandlerITMocks(t)
	getStripIdStr := idStr(getStrip.ID)
	mocks.expectDBStripGet(getStrip, nil)

	mocks.expectDBStripSave(nil)
	mocks.expectPublishStripEvent(t, model.Save, getStrip.ID, true, false, nil)

	req, w := prepareHttpTest(http.MethodDelete, ledstripIDProfilePath, uv{"id": getStripIdStr}, nil)

	mocks.lh.RemoveProfileForStrip(w, req)
	time.Sleep(50 * time.Millisecond)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func TestRemoveProfileForLEDStripIT_MissingStrip(t *testing.T) {
	getStrip := createValidDummyStrip()
	mocks := createLEDHandlerITMocks(t)
	getStripIdStr := idStr(getStrip.ID)
	mocks.expectDBStripGet(nil, errors.New("not found"))
	req, w := prepareHttpTest(http.MethodDelete, ledstripIDProfilePath, uv{"id": getStripIdStr}, nil)

	mocks.lh.RemoveProfileForStrip(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestRemoveProfileForLEDStripIT_SaveError(t *testing.T) {
	getStrip := createValidDummyStrip()
	mocks := createLEDHandlerITMocks(t)
	getStripIdStr := idStr(getStrip.ID)
	mocks.expectDBStripGet(getStrip, nil)
	mocks.expectDBStripSave(errors.New("save error"))
	req, w := prepareHttpTest(http.MethodDelete, ledstripIDProfilePath, uv{"id": getStripIdStr}, nil)

	mocks.lh.RemoveProfileForStrip(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestRemoveProfileForLEDStripIT_PublishError(t *testing.T) {
	getStrip := createValidDummyStrip()
	mocks := createLEDHandlerITMocks(t)
	getStripIdStr := idStr(getStrip.ID)
	mocks.expectDBStripGet(getStrip, nil)

	mocks.expectDBStripSave(nil)
	mocks.expectPublishStripEvent(t, model.Save, getStrip.ID, true, false, errors.New("publish error"))

	req, w := prepareHttpTest(http.MethodDelete, ledstripIDProfilePath, uv{"id": getStripIdStr}, nil)

	mocks.lh.RemoveProfileForStrip(w, req)
	time.Sleep(50 * time.Millisecond)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func (lhm *lhITMocks) expectDBStripGet(getStrip *model.LedStrip, getStripError error) {
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

func (lhm *lhITMocks) expectDBStripSave(saveError error) {
	lhm.lsDbh.
		EXPECT().
		Save(mock.Anything).
		Return(saveError)
}

func (lhm *lhITMocks) expectPublishStripEvent(t *testing.T, typ model.EventType, id int64, valid bool, expectProfile bool, publishError error) {
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

func createLEDHandlerITMocks(t *testing.T) *lhITMocks {
	i := do.New()
	bm := createBaseMocks(i, t)
	ls, err := service.NewLEDService(i)
	assert.NoError(t, err)
	do.ProvideValue(i, ls)
	lh, err := NewLEDHandler(i)
	assert.NoError(t, err)
	return &lhITMocks{
		baseMocks: bm,
		lh:        lh.(*ledHandlerImpl),
	}
}
