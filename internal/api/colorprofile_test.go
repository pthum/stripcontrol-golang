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

type cphMocks struct {
	*baseMocks
	cph *CPHandlerImpl
}

func TestCPRoutes(t *testing.T) {
	bm := createBaseMocks(t)
	routes := colorProfileRoutes(bm.cpDbh, bm.mh)
	assert.Equal(t, 5, len(routes))
}

func TestGetAllColorProfiles(t *testing.T) {
	tests := []getTest[model.ColorProfile]{
		{
			name:           "profiles_unavailable",
			returnError:    errors.New("nothing found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "profiles_available",
			returnError:    nil,
			returnObj:      createDummyProfile(),
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mocks := createCPHandlerMocks(t)
			destarr := []model.ColorProfile{}
			if tc.returnError == nil {
				destarr = append(destarr, *tc.returnObj)
			}
			mocks.cpDbh.
				EXPECT().
				GetAll().
				Return(destarr, tc.returnError).
				Once()

			req, w := prepareHttpTest(http.MethodGet, profilePath, nil, nil)

			mocks.cph.GetAllColorProfiles(w, req)
			res := w.Result()
			defer res.Body.Close()

			if tc.returnError == nil {
				var result []model.ColorProfile
				bodyToObj(t, res, &result)

				assert.Equal(t, *tc.returnObj, result[0])
			}

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestGetColorProfile(t *testing.T) {
	tests := []getTest[model.ColorProfile]{
		{
			name:           "profile_unavailable",
			returnError:    errors.New("nothing found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "profile_available",
			returnError:    nil,
			returnObj:      createDummyProfile(),
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mocks := createCPHandlerMocks(t)
			idS := idStringOrDefault(tc.returnObj, "9000")
			mocks.expectDBProfileGet(tc.returnObj, tc.returnError)
			req, w := prepareHttpTest(http.MethodGet, profileIDPath, uv{"id": idS}, nil)

			mocks.cph.GetColorProfile(w, req)

			res := w.Result()
			defer res.Body.Close()

			if tc.returnError == nil {
				var result model.ColorProfile
				bodyToObj(t, res, &result)
				assert.Equal(t, *tc.returnObj, result)
			}

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestCreateColorProfile(t *testing.T) {
	tests := []createTest[model.ColorProfile]{
		{
			name:           "success_case",
			returnError:    nil,
			body:           createDummyProfile(),
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
			body:           createDummyProfile(),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mocks := createCPHandlerMocks(t)
			var newId int64
			var body io.Reader
			if tc.body != nil {
				mocks.cpDbh.
					EXPECT().
					Create(mock.Anything).
					Run(func(input *model.ColorProfile) {
						// id should have been generated
						assert.NotEqual(t, tc.body.ID, input.ID)
						newId = input.ID
					}).
					Return(tc.returnError).
					Once()
				body = objToReader(t, tc.body)
			}
			req, w := prepareHttpTest(http.MethodPost, profilePath, nil, body)

			mocks.cph.CreateColorProfile(w, req)
			res := w.Result()
			defer res.Body.Close()

			if tc.body != nil && tc.returnError == nil {
				expectedObj := tc.body
				expectedObj.ID = newId
				var result model.ColorProfile
				bodyToObj(t, res, &result)

				assert.Equal(t, *expectedObj, result)
				assert.Contains(t, res.Header["Location"][0], idStr(newId))
			}

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestDeleteColorProfile(t *testing.T) {
	tests := []deleteTest[model.ColorProfile]{
		{
			name:           "success_case",
			getObj:         createDummyProfile(),
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
			getObj:         createDummyProfile(),
			getError:       nil,
			deleteError:    errors.New("delete failed"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mocks := createCPHandlerMocks(t)

			mocks.expectDBProfileGet(tc.getObj, tc.getError)

			if tc.getError == nil {
				mocks.cpDbh.
					EXPECT().
					Delete(mock.Anything).
					Return(tc.deleteError)
				if tc.deleteError == nil {
					mocks.expectPublishProfileEvent(t, model.Delete, tc.getObj.ID, nil)
				}
			}
			idS := idStringOrDefault(tc.getObj, "9000")
			req, w := prepareHttpTest(http.MethodDelete, profileIDPath, uv{"id": idS}, nil)

			mocks.cph.DeleteColorProfile(w, req)

			// small sleep to have the async routines run
			time.Sleep(50 * time.Millisecond)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestUpdateColorProfile(t *testing.T) {
	tests := []updateTest[model.ColorProfile]{
		{
			name:           "success_case",
			body:           createDummyProfile(),
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
			body:           createDummyProfile(),
			getError:       nil,
			updateError:    errors.New("update failed"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mocks := createCPHandlerMocks(t)

			dbO := *createProfile(105, 100, 100, 100, 2)

			mocks.expectDBProfileGet(&dbO, tc.getError)

			var body io.Reader
			if tc.body != nil {
				body = objToReader(t, tc.body)
				mocks.cpDbh.
					EXPECT().
					Update(dbO, *tc.body).
					Return(tc.updateError)

				if tc.updateError == nil {
					mocks.expectPublishProfileEvent(t, model.Save, tc.body.ID, tc.body)
				}
			}

			idS := idStr(dbO.ID)
			req, w := prepareHttpTest(http.MethodPut, profileIDPath, uv{"id": idS}, body)

			mocks.cph.UpdateColorProfile(w, req)

			// small sleep to have the async routines run
			time.Sleep(50 * time.Millisecond)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
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
	bm := createBaseMocks(t)
	return &cphMocks{
		baseMocks: bm,
		cph: &CPHandlerImpl{
			dbh: bm.cpDbh,
			mh:  bm.mh,
		},
	}
}
