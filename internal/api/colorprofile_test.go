package api

import (
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type cphMocks struct {
	*baseMocks
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
	mocks.cpDbh.
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
	mocks.cpDbh.
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
	mocks.expectDBProfileGet(retObj, nil)
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
	mocks.expectDBProfileGet(retObj, errors.New("not found"))
	req, w := prepareHttpTest(http.MethodGet, profileIDPath, uv{"id": idS}, nil)

	mocks.cph.GetColorProfile(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestCreateColorProfile(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	inBody := createDummyProfile()
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
	body := objToReader(t, inBody)
	req, w := prepareHttpTest(http.MethodPost, profilePath, nil, body)

	mocks.cph.CreateColorProfile(w, req)

	res := w.Result()
	defer res.Body.Close()

	expectedObj := inBody
	expectedObj.ID = newId
	var result model.ColorProfile
	bodyToObj(t, res, &result)

	assert.Equal(t, *expectedObj, result)
	assert.Contains(t, res.Header["Location"][0], idStr(newId))

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
	mocks.cpDbh.
		EXPECT().
		Create(mock.Anything).
		Run(func(input *model.ColorProfile) {
			// id should have been generated
			assert.NotEqual(t, inBody.ID, input.ID)
		}).
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
	mocks.expectDBProfileGet(getObj, nil)

	mocks.cpDbh.
		EXPECT().
		Delete(mock.Anything).
		Return(nil)
	mocks.expectPublishProfileEvent(t, model.Delete, getObj.ID, nil)

	idS := idStringOrDefault(getObj, "9000")
	req, w := prepareHttpTest(http.MethodDelete, profileIDPath, uv{"id": idS}, nil)

	mocks.cph.DeleteColorProfile(w, req)

	// small sleep to have the async routines run
	time.Sleep(50 * time.Millisecond)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func TestDeleteColorProfile_MissingDBProfile(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	getObj := createDummyProfile()
	mocks.expectDBProfileGet(nil, errors.New("not found"))

	idS := idStringOrDefault(getObj, "9000")
	req, w := prepareHttpTest(http.MethodDelete, profileIDPath, uv{"id": idS}, nil)

	mocks.cph.DeleteColorProfile(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestDeleteColorProfile_DeleteError(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	getObj := createDummyProfile()
	mocks.expectDBProfileGet(getObj, nil)

	mocks.cpDbh.
		EXPECT().
		Delete(mock.Anything).
		Return(errors.New("delete error"))

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
	mocks.expectDBProfileGet(&dbO, nil)

	mocks.cpDbh.
		EXPECT().
		Update(dbO, *inBody).
		Return(nil)
	mocks.expectPublishProfileEvent(t, model.Save, inBody.ID, inBody)

	idS := idStr(dbO.ID)
	req, w := prepareHttpTest(http.MethodPut, profileIDPath, uv{"id": idS}, body)

	mocks.cph.UpdateColorProfile(w, req)

	// small sleep to have the async routines run
	time.Sleep(50 * time.Millisecond)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestUpdateColorProfile_MissingDBProfile(t *testing.T) {
	mocks := createCPHandlerMocks(t)
	inBody := createDummyProfile()
	body := objToReader(t, inBody)
	dbO := *createProfile(105, 100, 100, 100, 2)
	mocks.expectDBProfileGet(nil, errors.New("not found"))

	idS := idStr(dbO.ID)
	req, w := prepareHttpTest(http.MethodPut, profileIDPath, uv{"id": idS}, body)

	mocks.cph.UpdateColorProfile(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
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
	mocks.expectDBProfileGet(&dbO, nil)

	mocks.cpDbh.
		EXPECT().
		Update(dbO, *inBody).
		Return(errors.New("update failed"))

	idS := idStr(dbO.ID)
	req, w := prepareHttpTest(http.MethodPut, profileIDPath, uv{"id": idS}, body)

	mocks.cph.UpdateColorProfile(w, req)

	// small sleep to have the async routines run
	time.Sleep(50 * time.Millisecond)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func (chm *cphMocks) expectPublishProfileEvent(t *testing.T, typ model.EventType, id int64, body *model.ColorProfile) {
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
		}).
		Return(nil)
}

func idStringOrDefault(obj *model.ColorProfile, def string) string {
	idS := def
	if obj != nil {
		idS = idStr(obj.ID)
	}
	return idS
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
	cph, err := NewCPHandler(i)
	assert.NoError(t, err)
	return &cphMocks{
		baseMocks: bm,
		cph:       cph.(*CPHandlerImpl),
	}
}
