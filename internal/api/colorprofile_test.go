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
				destarr := dest.(*model.ColorProfile)
				if tc.returnError == nil {
					*destarr = tc.returnObj
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
