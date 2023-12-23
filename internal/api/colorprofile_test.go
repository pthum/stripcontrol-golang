package api

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/pthum/stripcontrol-golang/internal/service"
	servicemocks "github.com/pthum/stripcontrol-golang/internal/service/mocks"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type cphMocks struct {
	cps *servicemocks.CPService
	cph *CPHandlerImpl
}

func TestCPRoutes(t *testing.T) {
	mcks := createCPHandlerMocks(t)
	routes := mcks.cph.colorProfileRoutes()
	assert.Equal(t, 5, len(routes))
}

func TestGetAllColorProfiles(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	expRet := createDummyProfile()
	destarr := []model.ColorProfile{*expRet}
	mocks.cps.
		EXPECT().
		GetAll().
		Return(destarr, nil).
		Once()
	req, w := prepareHttpTest(http.MethodGet, profilePath, nil, nil)

	mocks.cph.GetAllColorProfiles(w, req)

	res := w.Result()
	defer res.Body.Close()
	var result []model.ColorProfile
	bodyToObj(t, res, &result)
	assert.Equal(t, *expRet, result[0])
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetAllColorProfiles_GetError(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	destarr := []model.ColorProfile{}
	mocks.cps.
		EXPECT().
		GetAll().
		Return(destarr, errors.New("get error")).
		Once()

	req, w := prepareHttpTest(http.MethodGet, profilePath, nil, nil)

	mocks.cph.GetAllColorProfiles(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestGetColorProfile(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	retObj := createDummyProfile()
	idS := idStringOrDefault(retObj, "9000")
	mocks.cps.
		EXPECT().
		GetColorProfile(idS).
		Return(retObj, nil)
	req, w := prepareHttpTest(http.MethodGet, profileIDPath, uv{"id": idS}, nil)

	mocks.cph.GetColorProfile(w, req)

	res := w.Result()
	defer res.Body.Close()
	var result model.ColorProfile
	bodyToObj(t, res, &result)
	assert.Equal(t, *retObj, result)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetColorProfile_GetError(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	retObj := createDummyProfile()
	idS := idStringOrDefault(retObj, "9000")
	mocks.cps.
		EXPECT().
		GetColorProfile(idS).
		Return(retObj, errors.New("not found"))
	req, w := prepareHttpTest(http.MethodGet, profileIDPath, uv{"id": idS}, nil)

	mocks.cph.GetColorProfile(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestCreateColorProfile(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	inBody := createDummyProfile()
	mocks.cps.
		EXPECT().
		CreateColorProfile(mock.Anything).
		Return(nil).
		Once()
	body := objToReader(t, inBody)
	req, w := prepareHttpTest(http.MethodPost, profilePath, nil, body)

	mocks.cph.CreateColorProfile(w, req)

	res := w.Result()
	defer res.Body.Close()

	expectedObj := inBody
	var result model.ColorProfile
	bodyToObj(t, res, &result)

	assert.Equal(t, *expectedObj, result)
	assert.Contains(t, res.Header["Location"][0], idStr(expectedObj.ID))

	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestCreateColorProfile_MissingBody(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	var body io.Reader
	req, w := prepareHttpTest(http.MethodPost, profilePath, nil, body)

	mocks.cph.CreateColorProfile(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestCreateColorProfile_SaveError(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	inBody := createDummyProfile()
	mocks.cps.
		EXPECT().
		CreateColorProfile(mock.Anything).
		Return(errors.New("save failed")).
		Once()
	body := objToReader(t, inBody)
	req, w := prepareHttpTest(http.MethodPost, profilePath, nil, body)

	mocks.cph.CreateColorProfile(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}
func TestDeleteColorProfile(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	getObj := createDummyProfile()
	mocks.cps.
		EXPECT().
		DeleteColorProfile(mock.Anything).
		Return(nil)

	idS := idStringOrDefault(getObj, "9000")
	req, w := prepareHttpTest(http.MethodDelete, profileIDPath, uv{"id": idS}, nil)

	mocks.cph.DeleteColorProfile(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func TestDeleteColorProfile_DeleteError(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	getObj := createDummyProfile()

	mocks.cps.
		EXPECT().
		DeleteColorProfile(mock.Anything).
		Return(model.NewAppErr(400, errors.New("delete error")))

	idS := idStringOrDefault(getObj, "9000")
	req, w := prepareHttpTest(http.MethodDelete, profileIDPath, uv{"id": idS}, nil)

	mocks.cph.DeleteColorProfile(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUpdateColorProfile(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	inBody := createDummyProfile()
	body := objToReader(t, inBody)
	dbO := *createProfile(105, 100, 100, 100, 2)

	mocks.cps.
		EXPECT().
		UpdateColorProfile(mock.Anything, mock.Anything).
		Return(nil)

	idS := idStr(dbO.ID)
	req, w := prepareHttpTest(http.MethodPut, profileIDPath, uv{"id": idS}, body)

	mocks.cph.UpdateColorProfile(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestUpdateColorProfile_MissingBody(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	var body io.Reader
	dbO := *createProfile(105, 100, 100, 100, 2)

	idS := idStr(dbO.ID)
	req, w := prepareHttpTest(http.MethodPut, profileIDPath, uv{"id": idS}, body)

	mocks.cph.UpdateColorProfile(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUpdateColorProfile_UpdateError(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	inBody := createDummyProfile()
	body := objToReader(t, inBody)
	dbO := *createProfile(105, 100, 100, 100, 2)
	mocks.cps.
		EXPECT().
		UpdateColorProfile(mock.Anything, mock.Anything).
		Return(model.NewAppErr(400, errors.New("update failed")))

	idS := idStr(dbO.ID)
	req, w := prepareHttpTest(http.MethodPut, profileIDPath, uv{"id": idS}, body)

	mocks.cph.UpdateColorProfile(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func createCPHandlerMocks(t *testing.T) *cphMocks {
	i := do.New()
	cps := servicemocks.NewCPService(t)
	do.ProvideValue[service.CPService](i, cps)
	cph, err := NewCPHandler(i)
	assert.NoError(t, err)
	return &cphMocks{
		cps: cps,
		cph: cph.(*CPHandlerImpl),
	}
}
