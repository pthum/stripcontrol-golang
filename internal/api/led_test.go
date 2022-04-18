package api

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/pthum/null"
	dbm "github.com/pthum/stripcontrol-golang/internal/database/mocks"
	mhm "github.com/pthum/stripcontrol-golang/internal/messaging/mocks"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllLedStrips(t *testing.T) {
	returnObj := model.LedStrip{
		ID:          185,
		Description: "Test",
		Enabled:     false,
		MisoPin:     null.IntFrom(12),
		Name:        "Test",
		NumLeds:     null.IntFrom(5),
		SclkPin:     null.IntFrom(13),
		SpeedHz:     null.IntFrom(80000),
	}

	tests := []struct {
		name           string
		returnError    error
		returnObj      model.LedStrip
		expectedStatus int
	}{
		{
			name:           "strips_unavailable",
			returnError:    errors.New("nothing found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "strips_available",
			returnError:    nil,
			returnObj:      returnObj,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dbh := &dbm.DBHandler{}
			mh := &mhm.EventHandler{}
			cph := LEDHandlerImpl{
				dbh: dbh,
				mh:  mh,
			}

			dbh.EXPECT().GetAll(mock.Anything).Run(func(dest interface{}) {
				destarr := dest.(*[]model.LedStrip)
				if tc.returnError == nil {
					*destarr = append(*destarr, tc.returnObj)
				}
			}).Return(tc.returnError).Once()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			cph.GetAllLedStrips(w, req)
			res := w.Result()
			defer res.Body.Close()

			if tc.returnError == nil {
				var result []model.LedStrip
				err := unmarshHelp(res, &result)
				assert.Nil(t, err)
				assert.Equal(t, tc.returnObj, result[0])
			}

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			dbh.AssertExpectations(t)
			mh.AssertExpectations(t)
		})
	}
}

func TestGetLedStrips(t *testing.T) {
	returnObj := model.LedStrip{
		ID:          185,
		Description: "Test",
		Enabled:     false,
		MisoPin:     null.IntFrom(12),
		Name:        "Test",
		NumLeds:     null.IntFrom(5),
		SclkPin:     null.IntFrom(13),
		SpeedHz:     null.IntFrom(80000),
	}

	tests := []struct {
		name           string
		returnError    error
		returnObj      model.LedStrip
		expectedStatus int
	}{
		{
			name:           "strip_unavailable",
			returnError:    errors.New("nothing found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "strip_available",
			returnError:    nil,
			returnObj:      returnObj,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dbh := &dbm.DBHandler{}
			mh := &mhm.EventHandler{}
			lh := LEDHandlerImpl{
				dbh: dbh,
				mh:  mh,
			}

			dbh.EXPECT().Get(mock.Anything, mock.Anything).Run(func(id string, dest interface{}) {
				destobj := dest.(*model.LedStrip)
				if tc.returnError == nil {
					*destobj = tc.returnObj
				}
			}).Return(tc.returnError).Once()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			lh.GetLedStrip(w, req)
			res := w.Result()
			defer res.Body.Close()

			if tc.returnError == nil {
				var result model.LedStrip
				err := unmarshHelp(res, &result)
				assert.Nil(t, err)
				assert.Equal(t, tc.returnObj, result)
			}

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			dbh.AssertExpectations(t)
			mh.AssertExpectations(t)
		})
	}
}

func TestCreateLedStrip(t *testing.T) {
	returnObj := model.LedStrip{
		ID:          185,
		Description: "Test",
		Enabled:     false,
		MisoPin:     null.IntFrom(12),
		Name:        "Test",
		NumLeds:     null.IntFrom(5),
		SclkPin:     null.IntFrom(13),
		SpeedHz:     null.IntFrom(80000),
	}
	tests := []struct {
		name           string
		returnError    error
		body           *model.LedStrip
		expectedStatus int
	}{
		{
			name:           "success_case",
			returnError:    nil,
			body:           &returnObj,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "missing_body",
			returnError:    nil,
			body:           nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "error_on_save",
			returnError:    errors.New("save failed"),
			body:           &returnObj,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dbh := &dbm.DBHandler{}
			mh := &mhm.EventHandler{}
			lh := LEDHandlerImpl{
				dbh: dbh,
				mh:  mh,
			}
			var newId int64
			var body io.Reader
			if tc.body != nil {
				dbh.EXPECT().Create(mock.Anything).Run(func(input interface{}) {
					in := input.(*model.LedStrip)
					// id should have been generated
					assert.NotEqual(t, tc.body.ID, in.ID)
					newId = in.ID
				}).Return(tc.returnError).Once()
				body, _ = marshHelp(&tc.body)
				if tc.returnError == nil {
					mh.EXPECT().PublishStripSaveEvent(mock.Anything, mock.Anything).Return(nil)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/", body)
			w := httptest.NewRecorder()

			lh.CreateLedStrip(w, req)

			// small sleep to have the async routines run
			time.Sleep(50 * time.Millisecond)

			res := w.Result()
			defer res.Body.Close()

			if tc.body != nil && tc.returnError == nil {
				expectedObj := tc.body
				expectedObj.ID = newId
				var result model.LedStrip
				err := unmarshHelp(res, &result)
				assert.Nil(t, err)
				assert.Equal(t, *expectedObj, result)
				assert.Contains(t, res.Header["Location"][0], strconv.FormatInt(newId, 10))
			}

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			dbh.AssertExpectations(t)
			mh.AssertExpectations(t)
		})
	}
}

func TestDeleteLedStrip(t *testing.T) {
	returnObj := model.LedStrip{
		ID:          185,
		Description: "Test",
		Enabled:     false,
		MisoPin:     null.IntFrom(12),
		Name:        "Test",
		NumLeds:     null.IntFrom(5),
		SclkPin:     null.IntFrom(13),
		SpeedHz:     null.IntFrom(80000),
	}

	tests := []struct {
		name           string
		getObj         *model.LedStrip
		getError       error
		deleteError    error
		expectedStatus int
	}{
		{
			name:           "success_case",
			getObj:         &returnObj,
			getError:       nil,
			deleteError:    nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "missing_strip_to_delete",
			getObj:         nil,
			getError:       errors.New("not found"),
			deleteError:    nil,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "error_on_delete",
			getObj:         &returnObj,
			getError:       nil,
			deleteError:    errors.New("delete failed"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dbh := &dbm.DBHandler{}
			mh := &mhm.EventHandler{}
			lh := LEDHandlerImpl{
				dbh: dbh,
				mh:  mh,
			}

			dbh.EXPECT().Get(mock.Anything, mock.Anything).Run(func(id string, dest interface{}) {
				destobj := dest.(*model.LedStrip)
				if tc.getError == nil {
					*destobj = *tc.getObj
				}
			}).Return(tc.getError).Once()

			if tc.getError == nil {
				dbh.EXPECT().Delete(mock.Anything).Return(tc.deleteError)
				if tc.deleteError == nil {
					mh.EXPECT().PublishStripDeleteEvent(mock.Anything).Return(nil)
				}
			}

			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			w := httptest.NewRecorder()

			lh.DeleteLedStrip(w, req)

			// small sleep to have the async routines run
			time.Sleep(50 * time.Millisecond)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			dbh.AssertExpectations(t)
			mh.AssertExpectations(t)
		})
	}
}

func TestUpdateLedStrip(t *testing.T) {
	returnObj := model.LedStrip{
		ID:          185,
		Description: "Test",
		Enabled:     false,
		MisoPin:     null.IntFrom(12),
		Name:        "Test",
		NumLeds:     null.IntFrom(5),
		SclkPin:     null.IntFrom(13),
		SpeedHz:     null.IntFrom(80000),
	}

	dbObj := model.LedStrip{
		ID:          185,
		Description: "TestFromDb",
		Enabled:     false,
		MisoPin:     null.IntFrom(100),
		Name:        "TestFromDb",
		NumLeds:     null.IntFrom(99),
		SclkPin:     null.IntFrom(99),
		SpeedHz:     null.IntFrom(80001),
	}

	tests := []struct {
		name           string
		body           *model.LedStrip
		getError       error
		updateError    error
		expectedStatus int
	}{
		{
			name:           "success_case",
			body:           &returnObj,
			getError:       nil,
			updateError:    nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing_profile_to_update",
			body:           nil,
			getError:       errors.New("not found"),
			updateError:    nil,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "missing_profile_body",
			body:           nil,
			getError:       nil,
			updateError:    nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "error_on_update",
			body:        &returnObj,
			getError:    nil,
			updateError: errors.New("update failed"),
			// we ignore errors on update for the sake of performance, see comment in UpdateLedStrip
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dbh := &dbm.DBHandler{}
			mh := &mhm.EventHandler{}
			lh := LEDHandlerImpl{
				dbh: dbh,
				mh:  mh,
			}

			dbh.EXPECT().Get(mock.Anything, mock.Anything).Run(func(id string, dest interface{}) {
				destobj := dest.(*model.LedStrip)
				if tc.getError == nil {
					*destobj = *&dbObj
				}
			}).Return(tc.getError).Once()

			var body io.Reader
			if tc.body != nil {
				body, _ = marshHelp(&tc.body)
				dbh.EXPECT().Update(dbObj, *tc.body).Return(tc.updateError)
				if tc.updateError == nil {
					mh.EXPECT().PublishStripSaveEvent(mock.Anything, mock.Anything).Return(nil)
				}
			}

			req := httptest.NewRequest(http.MethodPut, "/", body)
			w := httptest.NewRecorder()

			lh.UpdateLedStrip(w, req)

			// small sleep to have the async routines run
			time.Sleep(50 * time.Millisecond)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			dbh.AssertExpectations(t)
			mh.AssertExpectations(t)
		})
	}
}

func TestGetProfileForLedStrip(t *testing.T) {
	returnObj := model.LedStrip{
		ID:          185,
		Description: "Test",
		Enabled:     false,
		MisoPin:     null.IntFrom(12),
		Name:        "Test",
		NumLeds:     null.IntFrom(5),
		SclkPin:     null.IntFrom(13),
		SpeedHz:     null.IntFrom(80000),
		ProfileID:   null.IntFrom(15),
	}
	unreferencedObj := model.LedStrip{
		ID:          185,
		Description: "Test",
		Enabled:     false,
		MisoPin:     null.IntFrom(12),
		Name:        "Test",
		NumLeds:     null.IntFrom(5),
		SclkPin:     null.IntFrom(13),
		SpeedHz:     null.IntFrom(80000),
	}

	fakeProfile := model.ColorProfile{
		ID:         15,
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
			name:           "success_case",
			getStripError:  nil,
			getStrip:       returnObj,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "strip_missing",
			getStripError:  errors.New("strip not found"),
			getStrip:       returnObj,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "no_profile_referenced",
			getStripError:  nil,
			getStrip:       unreferencedObj,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:            "profile_missing",
			getStripError:   nil,
			getStrip:        returnObj,
			getProfileError: errors.New("profile not found"),
			expectedStatus:  http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dbh := &dbm.DBHandler{}
			mh := &mhm.EventHandler{}
			lh := LEDHandlerImpl{
				dbh: dbh,
				mh:  mh,
			}

			dbh.EXPECT().Get(strconv.FormatInt(tc.getStrip.ID, 10), mock.Anything).Run(func(id string, dest interface{}) {
				destobj := dest.(*model.LedStrip)
				if tc.getStripError == nil {
					*destobj = tc.getStrip
				}
			}).Return(tc.getStripError).Once()
			if tc.getStripError == nil && tc.getStrip.ProfileID.Valid == true {
				dbh.EXPECT().Get(strconv.FormatInt(tc.getStrip.ProfileID.Int64, 10), mock.Anything).Run(func(id string, dest interface{}) {
					destobj := dest.(*model.ColorProfile)
					if tc.getProfileError == nil {
						*destobj = fakeProfile
					}
				}).Return(tc.getProfileError).Once()
			}
			req := httptest.NewRequest(http.MethodGet, ledstripPath+"/"+strconv.FormatInt(tc.getStrip.ID, 10)+"/profile", nil)
			w := httptest.NewRecorder()
			vars := map[string]string{
				"id": strconv.FormatInt(tc.getStrip.ID, 10),
			}
			req = mux.SetURLVars(req, vars)

			lh.GetProfileForStrip(w, req)
			res := w.Result()
			defer res.Body.Close()

			if tc.getStripError == nil && tc.getStrip.ProfileID.Valid == true && tc.getProfileError == nil {
				var result model.ColorProfile
				err := unmarshHelp(res, &result)
				assert.Nil(t, err)
				assert.Equal(t, fakeProfile, result)
			}

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			dbh.AssertExpectations(t)
			mh.AssertExpectations(t)
		})
	}
}

func TestUpdateProfileForLedStrip(t *testing.T) {
	returnObj := model.LedStrip{
		ID:          185,
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
		ID:         16,
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
		expectedStatus  int
	}{
		{
			name:           "success_case",
			getStripError:  nil,
			getStrip:       returnObj,
			body:           &updateProfile,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "strip_missing",
			getStripError:  errors.New("strip not found"),
			getStrip:       returnObj,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:            "body_missing",
			getStripError:   nil,
			getStrip:        returnObj,
			getProfileError: nil,
			body:            nil,
			expectedStatus:  http.StatusBadRequest,
		},
		{
			name:            "profile_missing",
			getStripError:   nil,
			getStrip:        returnObj,
			getProfileError: errors.New("profile not found"),
			body:            &updateProfile,
			expectedStatus:  http.StatusNotFound,
		},
		{
			name:           "strip_save_error",
			getStripError:  nil,
			getStrip:       returnObj,
			body:           &updateProfile,
			saveError:      errors.New("save error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dbh := &dbm.DBHandler{}
			mh := &mhm.EventHandler{}
			lh := LEDHandlerImpl{
				dbh: dbh,
				mh:  mh,
			}

			dbh.EXPECT().Get(strconv.FormatInt(tc.getStrip.ID, 10), mock.Anything).Run(func(id string, dest interface{}) {
				destobj := dest.(*model.LedStrip)
				if tc.getStripError == nil {
					*destobj = tc.getStrip
				}
			}).Return(tc.getStripError).Once()

			if tc.getStripError == nil && tc.body != nil {
				dbh.EXPECT().Get(strconv.FormatInt(tc.body.ID, 10), mock.Anything).Run(func(id string, dest interface{}) {
					destobj := dest.(*model.ColorProfile)
					if tc.getProfileError == nil {
						*destobj = *tc.body
					}
				}).Return(tc.getProfileError).Once()

				if tc.getProfileError == nil {
					dbh.EXPECT().Save(mock.Anything).Return(tc.saveError)

					if tc.saveError == nil {
						mh.EXPECT().PublishStripSaveEvent(mock.Anything, mock.Anything).Return(nil)
					}
				}
			}

			var body io.Reader
			if tc.body != nil {
				body, _ = marshHelp(tc.body)
			}

			req := httptest.NewRequest(http.MethodPut, ledstripPath+"/"+strconv.FormatInt(tc.getStrip.ID, 10)+"/profile", body)
			w := httptest.NewRecorder()
			vars := map[string]string{
				"id": strconv.FormatInt(tc.getStrip.ID, 10),
			}
			req = mux.SetURLVars(req, vars)

			lh.UpdateProfileForStrip(w, req)
			time.Sleep(50 * time.Millisecond)
			res := w.Result()
			defer res.Body.Close()

			if tc.getStripError == nil && tc.getProfileError == nil && tc.body != nil && tc.saveError == nil {
				var result model.ColorProfile
				err := unmarshHelp(res, &result)
				assert.Nil(t, err)
				assert.Equal(t, *tc.body, result)
			}

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			dbh.AssertExpectations(t)
			mh.AssertExpectations(t)
		})
	}
}

func TestRemoveProfileForLedStrip(t *testing.T) {
	returnObj := model.LedStrip{
		ID:          185,
		Description: "Test",
		Enabled:     false,
		MisoPin:     null.IntFrom(12),
		Name:        "Test",
		NumLeds:     null.IntFrom(5),
		SclkPin:     null.IntFrom(13),
		SpeedHz:     null.IntFrom(80000),
		ProfileID:   null.IntFrom(15),
	}

	tests := []struct {
		name           string
		getStripError  error
		getStrip       model.LedStrip
		expectedStatus int
	}{
		{
			name:           "success_case",
			getStripError:  nil,
			getStrip:       returnObj,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "strip_missing",
			getStripError:  errors.New("strip not found"),
			getStrip:       returnObj,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dbh := &dbm.DBHandler{}
			mh := &mhm.EventHandler{}
			lh := LEDHandlerImpl{
				dbh: dbh,
				mh:  mh,
			}

			dbh.EXPECT().Get(strconv.FormatInt(tc.getStrip.ID, 10), mock.Anything).Run(func(id string, dest interface{}) {
				destobj := dest.(*model.LedStrip)
				if tc.getStripError == nil {
					*destobj = tc.getStrip
				}
			}).Return(tc.getStripError).Once()

			if tc.getStripError == nil {
				dbh.EXPECT().Save(mock.Anything).Return(nil)
				mh.EXPECT().PublishStripSaveEvent(mock.Anything, mock.Anything).Return(nil)
			}

			req := httptest.NewRequest(http.MethodDelete, ledstripPath+"/"+strconv.FormatInt(tc.getStrip.ID, 10)+"/profile", nil)
			w := httptest.NewRecorder()
			vars := map[string]string{
				"id": strconv.FormatInt(tc.getStrip.ID, 10),
			}
			req = mux.SetURLVars(req, vars)

			lh.RemoveProfileForStrip(w, req)
			time.Sleep(50 * time.Millisecond)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			dbh.AssertExpectations(t)
			mh.AssertExpectations(t)
		})
	}
}
