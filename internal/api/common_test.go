package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	dbm "github.com/pthum/stripcontrol-golang/internal/database/mocks"
	mhm "github.com/pthum/stripcontrol-golang/internal/messaging/mocks"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type getTest[T any] struct {
	name           string
	returnError    error
	returnObj      *T
	expectedStatus int
}

type createTest[T any] struct {
	name           string
	returnError    error
	body           *T
	publishError   error
	expectedStatus int
}

type deleteTest[T any] struct {
	name           string
	getObj         *T
	getError       error
	deleteError    error
	publishError   error
	expectedStatus int
}

type updateTest[T any] struct {
	name           string
	body           *T
	getError       error
	updateError    error
	publishError   error
	expectedStatus int
}

type uv map[string]string

type baseMocks struct {
	cpDbh *dbm.DBHandler[model.ColorProfile]
	lsDbh *dbm.DBHandler[model.LedStrip]
	mh    *mhm.EventHandler
}

func createBaseMocks(t *testing.T) *baseMocks {
	cpDbh := dbm.NewDBHandler[model.ColorProfile](t)
	lsDbh := dbm.NewDBHandler[model.LedStrip](t)
	mh := mhm.NewEventHandler(t)
	return &baseMocks{
		cpDbh: cpDbh,
		lsDbh: lsDbh,
		mh:    mh,
	}
}

func (bm *baseMocks) expectDBProfileGet(getStrip *model.ColorProfile, getError error) {
	getStripIdStr := mock.Anything
	if getStrip != nil {
		getStripIdStr = idStr(getStrip.ID)
	}
	bm.cpDbh.
		EXPECT().
		Get(getStripIdStr).
		Return(getStrip, getError).
		Once()
}

func prepareHttpTest(method string, path string, uv map[string]string, body io.Reader) (req *http.Request, w *httptest.ResponseRecorder) {
	req = httptest.NewRequest(method, path, body)
	w = httptest.NewRecorder()
	if len(uv) > 0 {
		req = mux.SetURLVars(req, uv)
	}
	return
}

func bodyToObj(t *testing.T, res *http.Response, obj interface{}) {
	byteData, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(byteData, &obj)
	assert.Nil(t, err)
}

func objToReader(t *testing.T, obj interface{}) (body io.Reader) {
	data, err := json.Marshal(obj)
	assert.Nil(t, err)
	body = bytes.NewReader(data)
	return
}

func idStr(id int64) string {
	return strconv.FormatInt(id, 10)
}
