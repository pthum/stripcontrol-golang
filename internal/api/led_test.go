package api

import (
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type lhMocks struct {
	*baseMocks
	lh *LEDHandlerImpl
}

func TestLedRoutes(t *testing.T) {
	bm := createBaseMocks(t)
	routes := ledRoutes(bm.lsDbh, bm.cpDbh, bm.mh)
	assert.Equal(t, 8, len(routes))
}

func TestGetAllLedStrips(t *testing.T) {
	tests := []getTest[model.LedStrip]{
		{
			name:           "strips_unavailable",
			returnError:    errors.New("nothing found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "strips_available",
			returnError:    nil,
			returnObj:      createValidDummyStrip(),
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mocks := createLEDHandlerMocks(t)
			destarr := []model.LedStrip{}
			if tc.returnError == nil {
				destarr = append(destarr, *tc.returnObj)
			}
			mocks.lsDbh.
				EXPECT().
				GetAll().
				Return(destarr, tc.returnError).
				Once()

			req, w := prepareHttpTest(http.MethodGet, ledstripPath, nil, nil)

			mocks.lh.GetAllLedStrips(w, req)
			res := w.Result()
			defer res.Body.Close()

			if tc.returnError == nil {
				var result []model.LedStrip
				bodyToObj(t, res, &result)

				assert.Equal(t, *tc.returnObj, result[0])
			}

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestGetLedStrips(t *testing.T) {
	tests := []getTest[model.LedStrip]{
		{
			name:           "strip_unavailable",
			returnError:    errors.New("nothing found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "strip_available",
			returnError:    nil,
			returnObj:      createValidDummyStrip(),
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mocks := createLEDHandlerMocks(t)

			reqId := "6000"
			if tc.returnObj != nil {
				reqId = idStr(tc.returnObj.ID)
			}
			mocks.expectDBStripGet(tc.returnObj, tc.returnError)
			req, w := prepareHttpTest(http.MethodGet, ledstripIDPath, uv{"id": reqId}, nil)

			mocks.lh.GetLedStrip(w, req)
			res := w.Result()
			defer res.Body.Close()

			if tc.returnError == nil {
				var result model.LedStrip
				bodyToObj(t, res, &result)

				assert.Equal(t, *tc.returnObj, result)
			}

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestCreateLedStrip(t *testing.T) {
	tests := []createTest[model.LedStrip]{
		{
			name:           "success case",
			returnError:    nil,
			body:           createValidDummyStrip(),
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "missing body",
			returnError:    nil,
			body:           nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "error on save",
			returnError:    errors.New("save failed"),
			body:           createValidDummyStrip(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "error on publish",
			body:           createValidDummyStrip(),
			publishError:   errors.New("publish failed"),
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mocks := createLEDHandlerMocks(t)
			var newId int64
			var body io.Reader
			if tc.body != nil {
				mocks.lsDbh.
					EXPECT().
					Create(mock.Anything).
					Run(func(input *model.LedStrip) {
						// id should have been generated
						assert.NotEqual(t, tc.body.ID, input.ID)
						newId = input.ID
					}).
					Return(tc.returnError).
					Once()
				body = objToReader(t, tc.body)
				if tc.returnError == nil {
					mocks.expectPublishStripEvent(t, model.Save, newId, true, false, tc.publishError)
				}
			}

			req, w := prepareHttpTest(http.MethodPost, ledstripPath, nil, body)

			mocks.lh.CreateLedStrip(w, req)

			// small sleep to have the async routines run
			time.Sleep(50 * time.Millisecond)

			res := w.Result()
			defer res.Body.Close()

			if tc.body != nil && tc.returnError == nil {
				expectedObj := tc.body
				expectedObj.ID = newId
				var result model.LedStrip
				bodyToObj(t, res, &result)

				assert.Equal(t, *expectedObj, result)
				assert.Contains(t, res.Header["Location"][0], idStr(newId))
			}

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestDeleteLedStrip(t *testing.T) {
	tests := []deleteTest[model.LedStrip]{
		{
			name:           "success case",
			getObj:         createValidDummyStrip(),
			getError:       nil,
			deleteError:    nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "missing strip to delete",
			getObj:         nil,
			getError:       errors.New("not found"),
			deleteError:    nil,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "error on delete",
			getObj:         createValidDummyStrip(),
			getError:       nil,
			deleteError:    errors.New("delete failed"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "error on delete",
			getObj:         createValidDummyStrip(),
			getError:       nil,
			publishError:   errors.New("publish failed"),
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mocks := createLEDHandlerMocks(t)

			mocks.expectDBStripGet(tc.getObj, tc.getError)

			if tc.getError == nil {
				mocks.lsDbh.
					EXPECT().
					Delete(mock.Anything).
					Return(tc.deleteError)
				if tc.deleteError == nil {
					mocks.expectPublishStripEvent(t, model.Delete, tc.getObj.ID, false, false, tc.publishError)
				}
			}

			req, w := prepareHttpTest(http.MethodDelete, ledstripIDPath, uv{"id": "185"}, nil)

			mocks.lh.DeleteLedStrip(w, req)

			// small sleep to have the async routines run
			time.Sleep(50 * time.Millisecond)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestUpdateLedStrip(t *testing.T) {

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

	fakeProfile := model.ColorProfile{
		BaseModel:  model.BaseModel{ID: 15},
		Red:        null.IntFrom(100),
		Green:      null.IntFrom(100),
		Blue:       null.IntFrom(100),
		Brightness: null.IntFrom(2),
	}

	tests := []updateTest[model.LedStrip]{
		{
			name:           "success case",
			body:           inputObj,
			getError:       nil,
			updateError:    nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing profile to update",
			body:           nil,
			getError:       errors.New("not found"),
			updateError:    nil,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "missing profile body",
			body:           nil,
			getError:       nil,
			updateError:    nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "error on update",
			body:        inputObj,
			getError:    nil,
			updateError: errors.New("update failed"),
			// we ignore errors on update for the sake of performance, see comment in UpdateLedStrip
			expectedStatus: http.StatusOK,
		},
		{
			name:         "publish error",
			body:         inputObj,
			publishError: errors.New("publish error"),
			// we ignore errors on update for the sake of performance, see comment in UpdateLedStrip
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mocks := createLEDHandlerMocks(t)

			mocks.expectDBStripGet(&dbObj, tc.getError)

			var body io.Reader
			if tc.body != nil {
				body = objToReader(t, tc.body)
				mocks.lsDbh.
					EXPECT().
					Update(dbObj, *tc.body).
					Return(tc.updateError)
				if tc.updateError == nil {
					mocks.expectPublishStripEvent(t, model.Save, tc.body.ID, true, true, tc.publishError)
					mocks.expectDBProfileGet(&fakeProfile, nil)
				}
			}

			req, w := prepareHttpTest(http.MethodPut, ledstripIDPath, uv{"id": "185"}, body)

			mocks.lh.UpdateLedStrip(w, req)

			// small sleep to have the async routines run
			time.Sleep(50 * time.Millisecond)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestGetProfileForLedStrip(t *testing.T) {
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

	fakeProfile := model.ColorProfile{
		BaseModel:  model.BaseModel{ID: 15},
		Red:        null.IntFrom(100),
		Green:      null.IntFrom(100),
		Blue:       null.IntFrom(100),
		Brightness: null.IntFrom(2),
	}

	tests := []struct {
		name            string
		getStripError   error
		getStrip        model.LedStrip
		getProfileError error
		expectedStatus  int
	}{
		{
			name:           "success case",
			getStripError:  nil,
			getStrip:       returnObj,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "strip missing",
			getStripError:  errors.New("strip not found"),
			getStrip:       returnObj,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "no profile referenced",
			getStripError:  nil,
			getStrip:       *createValidDummyStrip(),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:            "profile missing",
			getStripError:   nil,
			getStrip:        returnObj,
			getProfileError: errors.New("profile not found"),
			expectedStatus:  http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mocks := createLEDHandlerMocks(t)

			stripIdStr := idStr(tc.getStrip.ID)
			mocks.expectDBStripGet(&tc.getStrip, tc.getStripError)

			if tc.getStripError == nil && tc.getStrip.ProfileID.Valid == true {
				mocks.expectDBProfileGet(&fakeProfile, tc.getProfileError)
			}

			req, w := prepareHttpTest(http.MethodGet, ledstripIDProfilePath, uv{"id": stripIdStr}, nil)

			mocks.lh.GetProfileForStrip(w, req)
			res := w.Result()
			defer res.Body.Close()

			if tc.getStripError == nil && tc.getStrip.ProfileID.Valid == true && tc.getProfileError == nil {
				var result model.ColorProfile
				bodyToObj(t, res, &result)
				assert.Equal(t, fakeProfile, result)
			}

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestUpdateProfileForLedStrip(t *testing.T) {
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

	tests := []struct {
		name            string
		getStripError   error
		getStrip        model.LedStrip
		body            *model.ColorProfile
		getProfileError error
		saveError       error
		publishError    error
		expectedStatus  int
	}{
		{
			name:           "success case",
			getStripError:  nil,
			getStrip:       returnObj,
			body:           &updateProfile,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "strip missing",
			getStripError:  errors.New("strip not found"),
			getStrip:       returnObj,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:            "body missing",
			getStripError:   nil,
			getStrip:        returnObj,
			getProfileError: nil,
			body:            nil,
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name:            "profile missing",
			getStripError:   nil,
			getStrip:        returnObj,
			getProfileError: errors.New("profile not found"),
			body:            &updateProfile,
			expectedStatus:  http.StatusNotFound,
		},
		{
			name:           "strip save error",
			getStripError:  nil,
			getStrip:       returnObj,
			body:           &updateProfile,
			saveError:      errors.New("save error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "strip publish error",
			getStripError:  nil,
			getStrip:       returnObj,
			body:           &updateProfile,
			publishError:   errors.New("publish error"),
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mocks := createLEDHandlerMocks(t)
			getStripIdStr := idStr(tc.getStrip.ID)
			mocks.expectDBStripGet(&tc.getStrip, tc.getStripError)

			if tc.getStripError == nil && tc.body != nil {
				mocks.expectDBProfileGet(tc.body, tc.getProfileError)

				if tc.getProfileError == nil {
					mocks.expectDBStripSave(tc.saveError)

					if tc.saveError == nil {
						mocks.expectPublishStripEvent(t, model.Save, tc.getStrip.ID, true, true, tc.publishError)
					}
				}
			}

			var body io.Reader
			if tc.body != nil {
				body = objToReader(t, tc.body)
			}

			req, w := prepareHttpTest(http.MethodPut, ledstripIDProfilePath, uv{"id": getStripIdStr}, body)

			mocks.lh.UpdateProfileForStrip(w, req)
			time.Sleep(50 * time.Millisecond)
			res := w.Result()
			defer res.Body.Close()

			if tc.getStripError == nil && tc.getProfileError == nil && tc.body != nil && tc.saveError == nil {
				var result model.ColorProfile
				bodyToObj(t, res, &result)
				assert.Equal(t, *tc.body, result)
			}

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestRemoveProfileForLedStrip(t *testing.T) {
	tests := []struct {
		name           string
		getStripError  error
		getStrip       model.LedStrip
		saveError      error
		publishError   error
		expectedStatus int
	}{
		{
			name:           "success case",
			getStripError:  nil,
			getStrip:       *createValidDummyStrip(),
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "strip missing",
			getStripError:  errors.New("strip not found"),
			getStrip:       *createValidDummyStrip(),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "save error",
			getStripError:  nil,
			getStrip:       *createValidDummyStrip(),
			saveError:      errors.New("save error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "publish error",
			getStripError:  nil,
			getStrip:       *createValidDummyStrip(),
			publishError:   errors.New("publish error"),
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mocks := createLEDHandlerMocks(t)
			getStripIdStr := idStr(tc.getStrip.ID)
			mocks.expectDBStripGet(&tc.getStrip, tc.getStripError)

			if tc.getStripError == nil {
				mocks.expectDBStripSave(tc.saveError)
				if tc.saveError == nil {
					mocks.expectPublishStripEvent(t, model.Save, tc.getStrip.ID, true, false, tc.publishError)
				}
			}

			req, w := prepareHttpTest(http.MethodDelete, ledstripIDProfilePath, uv{"id": getStripIdStr}, nil)

			mocks.lh.RemoveProfileForStrip(w, req)
			time.Sleep(50 * time.Millisecond)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}
func (lhm *lhMocks) expectDBStripGet(getStrip *model.LedStrip, getStripError error) {
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

func (lhm *lhMocks) expectDBStripSave(saveError error) {
	lhm.lsDbh.
		EXPECT().
		Save(mock.Anything).
		Return(saveError)
}

func (lhm *lhMocks) expectPublishStripEvent(t *testing.T, typ model.EventType, id int64, valid bool, expectProfile bool, publishError error) {
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

func createLEDHandlerMocks(t *testing.T) *lhMocks {
	bm := createBaseMocks(t)
	return &lhMocks{
		baseMocks: bm,
		lh: &LEDHandlerImpl{
			dbh:   bm.lsDbh,
			cpDbh: bm.cpDbh,
			mh:    bm.mh,
		},
	}
}
