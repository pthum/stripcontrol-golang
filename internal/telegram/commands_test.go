package telegram

import (
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/pthum/stripcontrol-golang/internal/service"
	servicemocks "github.com/pthum/stripcontrol-golang/internal/service/mocks"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type cmdMocks struct {
	lsvc *servicemocks.LEDService
	ch   *cmdHandler
}

const expectedCommandCount = 4

func TestGetCommands(t *testing.T) {
	mocks := createCmdHandlerMocks(t)

	bcs := mocks.ch.GetCommands()
	assert.Len(t, bcs, expectedCommandCount)
}

func TestActionHelp(t *testing.T) {
	mocks := createCmdHandlerMocks(t)

	res := mocks.ch.actionGetCommands(nil)

	assert.NotEmpty(t, res)
	// expect two times of the expectedCommandCount, as we expect that the command is listed and it's usage
	cnt := strings.Count(res, "/")
	assert.Equal(t, 2*expectedCommandCount, cnt)
}

func TestActionGetAll(t *testing.T) {
	ts := []model.LedStrip{
		{
			BaseModel: model.BaseModel{ID: 123},
			Name:      "test1",
		},
		{
			BaseModel: model.BaseModel{ID: 456},
			Name:      "test2",
		},
	}
	mocks := createCmdHandlerMocks(t)
	mocks.lsvc.EXPECT().GetAll().Return(ts, nil)

	res := mocks.ch.actionGetAll(nil)

	assert.NotEmpty(t, res)
	assert.Contains(t, res, ts[0].GetStringID())
	assert.Contains(t, res, ts[0].Name)
	assert.Contains(t, res, ts[1].GetStringID())
	assert.Contains(t, res, ts[1].Name)
}

func TestActionGetAll_Error(t *testing.T) {
	mocks := createCmdHandlerMocks(t)
	mocks.lsvc.EXPECT().GetAll().Return(nil, assert.AnError)

	res := mocks.ch.actionGetAll(nil)

	assert.NotEmpty(t, res)
	assert.Contains(t, res, "Error getting")
}

func TestSetLed_Disable(t *testing.T) {
	id := "12"
	mocks := createCmdHandlerMocks(t)
	mocks.lsvc.
		EXPECT().
		GetLEDStrip(id).
		Return(&model.LedStrip{BaseModel: model.BaseModel{ID: 12}, Name: "test"}, nil)
	mocks.lsvc.
		EXPECT().
		UpdateLEDStrip(id, mock.Anything).
		Run(func(id string, updMdl model.LedStrip) {
			assert.False(t, updMdl.Enabled)
		}).
		Return(nil)
	msg := createTestMessage("/ledoff " + id)

	res := mocks.ch.setLEDState(false, msg)

	assert.NotEmpty(t, res)
	assert.Contains(t, res, "Turning off ")
	assert.Contains(t, res, "Disabled ID "+id)
}

func TestSetLed_Enable(t *testing.T) {
	id := "12"
	mocks := createCmdHandlerMocks(t)
	mocks.lsvc.
		EXPECT().
		GetLEDStrip(id).
		Return(&model.LedStrip{BaseModel: model.BaseModel{ID: 12}, Name: "test"}, nil)
	mocks.lsvc.
		EXPECT().
		UpdateLEDStrip(id, mock.Anything).
		Run(func(id string, updMdl model.LedStrip) {
			assert.True(t, updMdl.Enabled)
		}).
		Return(nil)
	msg := createTestMessage("/ledon " + id)

	res := mocks.ch.setLEDState(true, msg)

	assert.NotEmpty(t, res)
	assert.Contains(t, res, "Turning on ")
	assert.Contains(t, res, "Enabled ID "+id)
}

func TestSetLed_EnableMulti(t *testing.T) {
	id1 := "12"
	id2 := "23"
	mocks := createCmdHandlerMocks(t)
	mocks.lsvc.
		EXPECT().
		GetLEDStrip(id1).
		Return(&model.LedStrip{BaseModel: model.BaseModel{ID: 12}, Name: "test"}, nil)
	mocks.lsvc.
		EXPECT().
		UpdateLEDStrip(id1, mock.Anything).
		Return(nil)
	mocks.lsvc.
		EXPECT().
		GetLEDStrip(id2).
		Return(&model.LedStrip{BaseModel: model.BaseModel{ID: 23}, Name: "test"}, nil)
	mocks.lsvc.
		EXPECT().
		UpdateLEDStrip(id2, mock.Anything).
		Return(nil)
	msg := createTestMessage("/ledon " + id1 + " " + id2)

	res := mocks.ch.setLEDState(true, msg)

	assert.NotEmpty(t, res)
	assert.Contains(t, res, "Turning on ")
	assert.Contains(t, res, "Enabled ID "+id1)
}

func TestSetLed_EnableNoIDs(t *testing.T) {
	mocks := createCmdHandlerMocks(t)
	msg := createTestMessage("/ledon ")

	res := mocks.ch.setLEDState(true, msg)

	assert.NotEmpty(t, res)
	assert.Contains(t, res, "Turning on ")
}

func TestSetLed_EnableGetError(t *testing.T) {
	id := "12"
	mocks := createCmdHandlerMocks(t)
	mocks.lsvc.
		EXPECT().
		GetLEDStrip(id).
		Return(nil, assert.AnError)
	msg := createTestMessage("/ledon " + id)

	res := mocks.ch.setLEDState(true, msg)

	assert.NotEmpty(t, res)
	assert.Contains(t, res, "Turning on ")
	assert.Contains(t, res, "Error getting ID "+id)
}

func TestSetLed_EnableUpdateError(t *testing.T) {
	id := "12"
	mocks := createCmdHandlerMocks(t)
	mocks.lsvc.
		EXPECT().
		GetLEDStrip(id).
		Return(&model.LedStrip{BaseModel: model.BaseModel{ID: 12}, Name: "test"}, nil)
	mocks.lsvc.
		EXPECT().
		UpdateLEDStrip(id, mock.Anything).
		Return(assert.AnError)
	msg := createTestMessage("/ledon " + id)

	res := mocks.ch.setLEDState(true, msg)

	assert.NotEmpty(t, res)
	assert.Contains(t, res, "Turning on ")
	assert.Contains(t, res, "Error updating ID "+id)
}

func TestCommandForMsg(t *testing.T) {
	mocks := createCmdHandlerMocks(t)
	msg := createTestMessage("/help")

	res := mocks.ch.commandForMsg(msg)

	assert.NotNil(t, res)
	assert.Equal(t, "/help", res.Cmd)
}
func TestCommandForMsg_Empty(t *testing.T) {
	mocks := createCmdHandlerMocks(t)
	msg := createTestMessage("")

	res := mocks.ch.commandForMsg(msg)

	assert.Nil(t, res)
}

func TestCallForMsg(t *testing.T) {
	mocks := createCmdHandlerMocks(t)
	msg := createTestMessage("/help")

	res := mocks.ch.callForMsg(msg)

	assert.NotNil(t, res)
	assert.Contains(t, res, "All available commands")
}

func TestCallForMsg_Unknown(t *testing.T) {
	mocks := createCmdHandlerMocks(t)
	msg := createTestMessage("/doesntexist")

	res := mocks.ch.callForMsg(msg)

	assert.NotNil(t, res)
	assert.Empty(t, res)
}

func createTestMessage(text string) *tgbotapi.Message {
	return &tgbotapi.Message{Text: text}
}

func createCmdHandlerMocks(t *testing.T) *cmdMocks {
	i := do.New()
	lsvc := servicemocks.NewLEDService(t)
	do.ProvideValue[service.LEDService](i, lsvc)
	ch := NewCmdHandler(i)
	return &cmdMocks{
		lsvc: lsvc,
		ch:   ch,
	}
}
