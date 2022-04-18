package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/pthum/null"
	dbm "github.com/pthum/stripcontrol-golang/internal/database/mocks"
	mhm "github.com/pthum/stripcontrol-golang/internal/messaging/mocks"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllColorProfiles(t *testing.T) {
	returnObj := model.ColorProfile{
		ID:         185,
		Red:        null.IntFrom(123),
		Green:      null.IntFrom(234),
		Blue:       null.IntFrom(12),
		Brightness: null.IntFrom(1),
	}

	tests := []struct {
		name           string
		returnError    error
		returnObj      model.ColorProfile
		expectedStatus int
	}{
		{
			name:           "profiles_unavailable",
			returnError:    errors.New("nothing found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "profiles_available",
			returnError:    nil,
			returnObj:      returnObj,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dbh := &dbm.DBHandler{}
			mh := &mhm.EventHandler{}
			cph := CPHandlerImpl{
				dbh: dbh,
				mh:  mh,
			}

			dbh.EXPECT().GetAll(mock.Anything).Run(func(dest interface{}) {
				destarr := dest.(*[]model.ColorProfile)
				if tc.returnError == nil {
					*destarr = append(*destarr, tc.returnObj)
				}
			}).Return(tc.returnError).Once()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			cph.GetAllColorProfiles(w, req)
			res := w.Result()
			defer res.Body.Close()

			if tc.returnError == nil {
				var result []model.ColorProfile
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

func TestGetColorProfile(t *testing.T) {
	returnObj := model.ColorProfile{
		ID:         185,
		Red:        null.IntFrom(123),
		Green:      null.IntFrom(234),
		Blue:       null.IntFrom(12),
		Brightness: null.IntFrom(1),
	}

	tests := []struct {
		name           string
		returnError    error
		returnObj      model.ColorProfile
		expectedStatus int
	}{
		{
			name:           "profile_unavailable",
			returnError:    errors.New("nothing found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "profile_available",
			returnError:    nil,
			returnObj:      returnObj,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dbh := &dbm.DBHandler{}
			mh := &mhm.EventHandler{}
			cph := CPHandlerImpl{
				dbh: dbh,
				mh:  mh,
			}

			dbh.EXPECT().Get(mock.Anything, mock.Anything).Run(func(id string, dest interface{}) {
				destobj := dest.(*model.ColorProfile)
				if tc.returnError == nil {
					*destobj = tc.returnObj
				}
			}).Return(tc.returnError).Once()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			cph.GetColorProfile(w, req)
			res := w.Result()
			defer res.Body.Close()

			if tc.returnError == nil {
				var result model.ColorProfile
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

func TestCreateColorProfile(t *testing.T) {
	returnObj := model.ColorProfile{
		ID:         185,
		Red:        null.IntFrom(123),
		Green:      null.IntFrom(234),
		Blue:       null.IntFrom(12),
		Brightness: null.IntFrom(1),
	}

	tests := []struct {
		name           string
		returnError    error
		body           *model.ColorProfile
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
			cph := CPHandlerImpl{
				dbh: dbh,
				mh:  mh,
			}
			var newId int64
			var body io.Reader
			if tc.body != nil {
				dbh.EXPECT().Create(mock.Anything).Run(func(input interface{}) {
					in := input.(*model.ColorProfile)
					// id should have been generated
					assert.NotEqual(t, tc.body.ID, in.ID)
					newId = in.ID
				}).Return(tc.returnError).Once()
				body, _ = marshHelp(&tc.body)
			}
			req := httptest.NewRequest(http.MethodGet, "/", body)
			w := httptest.NewRecorder()

			cph.CreateColorProfile(w, req)
			res := w.Result()
			defer res.Body.Close()

			if tc.body != nil && tc.returnError == nil {
				expectedObj := tc.body
				expectedObj.ID = newId
				var result model.ColorProfile
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

func TestDeleteColorProfile(t *testing.T) {
	returnObj := model.ColorProfile{
		ID:         185,
		Red:        null.IntFrom(123),
		Green:      null.IntFrom(234),
		Blue:       null.IntFrom(12),
		Brightness: null.IntFrom(1),
	}

	tests := []struct {
		name           string
		getObj         *model.ColorProfile
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
			name:           "missing_profile_to_delete",
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
			cph := CPHandlerImpl{
				dbh: dbh,
				mh:  mh,
			}

			dbh.EXPECT().Get(mock.Anything, mock.Anything).Run(func(id string, dest interface{}) {
				destobj := dest.(*model.ColorProfile)
				if tc.getError == nil {
					*destobj = *tc.getObj
				}
			}).Return(tc.getError).Once()

			if tc.getError == nil {
				dbh.EXPECT().Delete(mock.Anything).Return(tc.deleteError)
				if tc.deleteError == nil {
					mh.EXPECT().PublishProfileDeleteEvent(mock.Anything).Return(nil)
				}
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			cph.DeleteColorProfile(w, req)

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

func TestUpdateColorProfile(t *testing.T) {
	returnObj := model.ColorProfile{
		ID:         185,
		Red:        null.IntFrom(123),
		Green:      null.IntFrom(234),
		Blue:       null.IntFrom(12),
		Brightness: null.IntFrom(1),
	}

	dbObj := model.ColorProfile{
		ID:         185,
		Red:        null.IntFrom(100),
		Green:      null.IntFrom(100),
		Blue:       null.IntFrom(100),
		Brightness: null.IntFrom(2),
	}

	tests := []struct {
		name           string
		body           *model.ColorProfile
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
			name:           "error_on_update",
			body:           &returnObj,
			getError:       nil,
			updateError:    errors.New("update failed"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dbh := &dbm.DBHandler{}
			mh := &mhm.EventHandler{}
			cph := CPHandlerImpl{
				dbh: dbh,
				mh:  mh,
			}

			dbh.EXPECT().Get(mock.Anything, mock.Anything).Run(func(id string, dest interface{}) {
				destobj := dest.(*model.ColorProfile)
				if tc.getError == nil {
					*destobj = *&dbObj
				}
			}).Return(tc.getError).Once()

			var body io.Reader
			if tc.body != nil {
				body, _ = marshHelp(&tc.body)
				dbh.EXPECT().Update(dbObj, *tc.body).Return(tc.updateError)
				if tc.updateError == nil {
					mh.EXPECT().PublishProfileSaveEvent(mock.Anything, mock.Anything).Return(nil)
				}
			}

			req := httptest.NewRequest(http.MethodGet, "/", body)
			w := httptest.NewRecorder()

			cph.UpdateColorProfile(w, req)

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

func unmarshHelp(res *http.Response, obj interface{}) (err error) {
	byteData, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(byteData, &obj)
	return
}

func marshHelp(obj interface{}) (body io.Reader, err error) {
	data, err := json.Marshal(obj)
	body = bytes.NewReader(data)
	return
}
