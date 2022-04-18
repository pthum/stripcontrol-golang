package api

import (
	"errors"
	"io"
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

			req := httptest.NewRequest(http.MethodGet, "/", body)

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

			req := httptest.NewRequest(http.MethodGet, "/", nil)

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
