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
	servicemocks "github.com/pthum/stripcontrol-golang/internal/service/mocks"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type lhMocks struct {
	*baseMocks
	lsvc *servicemocks.LEDService
	lh   *ledHandlerImpl
}

func TestLedRoutes(t *testing.T) {
	mcks := createLEDHandlerMocks(t)
	routes := mcks.lh.ledRoutes()
	assert.Equal(t, 8, len(routes))
}

func TestGetAllLEDStrips(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	destarr := []model.LedStrip{*createValidDummyStrip()}
	mocks.lsvc.
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

func TestGetAllLEDStrips_Error(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	destarr := []model.LedStrip{}
	mocks.lsvc.
		EXPECT().
		GetAll().
		Return(destarr, assert.AnError).
		Once()
	req, w := prepareHttpTest(http.MethodGet, ledstripPath, nil, nil)

	mocks.lh.GetAllLedStrips(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestGetLEDStrip(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	retObj := createValidDummyStrip()
	reqId := idStr(retObj.ID)
	mocks.lsvc.
		EXPECT().
		GetLEDStrip(reqId).
		Return(retObj, nil).
		Once()
	req, w := prepareHttpTest(http.MethodGet, ledstripIDPath, uv{"id": reqId}, nil)

	mocks.lh.GetLedStrip(w, req)
	res := w.Result()
	defer res.Body.Close()

	var result model.LedStrip
	bodyToObj(t, res, &result)

	assert.Equal(t, *retObj, result)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetLEDStrip_Error(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	reqId := "6000"
	mocks.lsvc.
		EXPECT().
		GetLEDStrip(reqId).
		Return(nil, errors.New("nothing found")).
		Once()
	req, w := prepareHttpTest(http.MethodGet, ledstripIDPath, uv{"id": reqId}, nil)

	mocks.lh.GetLedStrip(w, req)
	res := w.Result()
	defer res.Body.Close()

	var result model.LedStrip
	bodyToObj(t, res, &result)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestCreateLEDStrip(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	reqObj := createValidDummyStrip()

	mocks.lsvc.
		EXPECT().
		CreateLEDStrip(mock.Anything).
		Return(nil).
		Once()

	body := objToReader(t, reqObj)
	req, w := prepareHttpTest(http.MethodPost, ledstripPath, nil, body)

	mocks.lh.CreateLedStrip(w, req)

	res := w.Result()
	defer res.Body.Close()

	expectedObj := reqObj
	var result model.LedStrip
	bodyToObj(t, res, &result)

	assert.Equal(t, *expectedObj, result)
	assert.Contains(t, res.Header["Location"][0], idStr(expectedObj.ID))
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestCreateLEDStrip_MissingBody(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	var body io.Reader
	req, w := prepareHttpTest(http.MethodPost, ledstripPath, nil, body)

	mocks.lh.CreateLedStrip(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestCreateLEDStrip_Error(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	reqObj := createValidDummyStrip()
	mocks.lsvc.
		EXPECT().
		CreateLEDStrip(mock.Anything).
		Return(assert.AnError).
		Once()
	body := objToReader(t, reqObj)
	req, w := prepareHttpTest(http.MethodPost, ledstripPath, nil, body)

	mocks.lh.CreateLedStrip(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestDeleteLEDStrip(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	mocks.lsvc.
		EXPECT().
		DeleteLEDStrip(mock.Anything).
		Return(nil).
		Once()
	req, w := prepareHttpTest(http.MethodDelete, ledstripIDPath, uv{"id": "185"}, nil)

	mocks.lh.DeleteLedStrip(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func TestDeleteLEDStrip_DeleteError(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	mocks.lsvc.
		EXPECT().
		DeleteLEDStrip(mock.Anything).
		Return(assert.AnError).
		Once()
	req, w := prepareHttpTest(http.MethodDelete, ledstripIDPath, uv{"id": "185"}, nil)

	mocks.lh.DeleteLedStrip(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestUpdateLEDStrip(t *testing.T) {
	inputObj := createValidDummyStrip()
	inputObj.ProfileID = null.IntFrom(15)
	fakeProfile := createDummyProfile()
	fakeProfile.ID = inputObj.ProfileID.Int64

	mocks := createLEDHandlerMocks(t)
	body := objToReader(t, inputObj)
	mocks.lsvc.
		EXPECT().
		UpdateLEDStrip(mock.Anything, mock.Anything).
		Return(nil).
		Once()

	req, w := prepareHttpTest(http.MethodPut, ledstripIDPath, uv{"id": "185"}, body)

	mocks.lh.UpdateLedStrip(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestUpdateLEDStrip_MissingBody(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	var body io.Reader
	req, w := prepareHttpTest(http.MethodPut, ledstripIDPath, uv{"id": "185"}, body)

	mocks.lh.UpdateLedStrip(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUpdateLEDStrip_UpdateError(t *testing.T) {
	inputObj := createValidDummyStrip()
	inputObj.ProfileID = null.IntFrom(15)
	mocks := createLEDHandlerMocks(t)
	mocks.lsvc.
		EXPECT().
		UpdateLEDStrip(mock.Anything, mock.Anything).
		Return(assert.AnError).
		Once()
	body := objToReader(t, inputObj)
	req, w := prepareHttpTest(http.MethodPut, ledstripIDPath, uv{"id": "185"}, body)

	mocks.lh.UpdateLedStrip(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestGetProfileForLEDStrip(t *testing.T) {
	fakeProfile := createDummyProfile()
	returnObj := createValidDummyStrip()
	returnObj.ProfileID = null.IntFrom(fakeProfile.ID)
	mocks := createLEDHandlerMocks(t)
	stripIdStr := idStr(returnObj.ID)
	mocks.lsvc.
		EXPECT().
		GetProfileForStrip(stripIdStr).
		Return(fakeProfile, nil).
		Once()

	req, w := prepareHttpTest(http.MethodGet, ledstripIDProfilePath, uv{"id": stripIdStr}, nil)

	mocks.lh.GetProfileForStrip(w, req)
	res := w.Result()
	defer res.Body.Close()

	var result model.ColorProfile
	bodyToObj(t, res, &result)
	assert.Equal(t, *fakeProfile, result)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetProfileForLEDStrip_Error(t *testing.T) {
	returnObj := createValidDummyStrip()
	mocks := createLEDHandlerMocks(t)
	stripIdStr := idStr(returnObj.ID)
	mocks.lsvc.
		EXPECT().
		GetProfileForStrip(stripIdStr).
		Return(nil, assert.AnError).
		Once()
	req, w := prepareHttpTest(http.MethodGet, ledstripIDProfilePath, uv{"id": stripIdStr}, nil)

	mocks.lh.GetProfileForStrip(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
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
	mocks.lsvc.
		EXPECT().
		UpdateProfileForStrip(getStripIdStr, mock.Anything).
		Return(&updateProfile, nil).
		Once()
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

func TestUpdateProfileForLEDStrip_Error(t *testing.T) {
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
	mocks.lsvc.
		EXPECT().
		UpdateProfileForStrip(getStripIdStr, mock.Anything).
		Return(nil, assert.AnError).
		Once()
	body := objToReader(t, updateProfile)
	req, w := prepareHttpTest(http.MethodPut, ledstripIDProfilePath, uv{"id": getStripIdStr}, body)

	mocks.lh.UpdateProfileForStrip(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestUpdateProfileForLEDStrip_MissingBody(t *testing.T) {
	mocks := createLEDHandlerMocks(t)
	getStripIdStr := "6000"
	var body io.Reader
	req, w := prepareHttpTest(http.MethodPut, ledstripIDProfilePath, uv{"id": getStripIdStr}, body)

	mocks.lh.UpdateProfileForStrip(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestRemoveProfileForLEDStrip(t *testing.T) {
	getStrip := createValidDummyStrip()
	mocks := createLEDHandlerMocks(t)
	getStripIdStr := idStr(getStrip.ID)
	mocks.lsvc.
		EXPECT().
		RemoveProfileForStrip(getStripIdStr).
		Return(nil).
		Once()

	req, w := prepareHttpTest(http.MethodDelete, ledstripIDProfilePath, uv{"id": getStripIdStr}, nil)

	mocks.lh.RemoveProfileForStrip(w, req)
	time.Sleep(50 * time.Millisecond)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func TestRemoveProfileForLEDStrip_Error(t *testing.T) {
	getStrip := createValidDummyStrip()
	mocks := createLEDHandlerMocks(t)
	getStripIdStr := idStr(getStrip.ID)
	mocks.lsvc.
		EXPECT().
		RemoveProfileForStrip(getStripIdStr).
		Return(assert.AnError).
		Once()
	req, w := prepareHttpTest(http.MethodDelete, ledstripIDProfilePath, uv{"id": getStripIdStr}, nil)

	mocks.lh.RemoveProfileForStrip(w, req)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
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

func createLEDHandlerMocks(t *testing.T) *lhMocks {
	i := do.New()
	bm := createBaseMocks(i, t)
	ls := servicemocks.NewLEDService(t)
	do.ProvideValue[service.LEDService](i, ls)
	lh, err := NewLEDHandler(i)
	assert.NoError(t, err)
	return &lhMocks{
		baseMocks: bm,
		lsvc:      ls,
		lh:        lh.(*ledHandlerImpl),
	}
}
