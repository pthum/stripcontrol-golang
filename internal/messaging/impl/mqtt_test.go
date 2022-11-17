package messagingimpl

import (
	"errors"
	"fmt"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/internal/config"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/pthum/stripcontrol-golang/internal/testutils"
	"github.com/stretchr/testify/assert"
)

type publishFunc func(topic string, payload interface{}) error

type mqttClientFake struct {
	mqtt.Client
	PublishCheck publishFunc
	disconnected bool
}

type mqttTokenFake struct {
	mqtt.Token
	err error
}

var testConfig = config.MessagingConfig{
	Host:         "localhost",
	Port:         "1234",
	ProfileTopic: "TestProfile",
	StripTopic:   "TestStrip",
}
var errReturn = errors.New("some error")
var errExpectedPublish = errors.New("failed to send message")

func TestMqttPublishProfileEvent(t *testing.T) {
	profileEvent := model.NewProfileEvent(null.IntFrom(123), model.Save)
	expectedPayload := testutils.JsonEncode(t, profileEvent)

	testFunc := createTestFunc(t, testConfig.ProfileTopic, expectedPayload, nil)
	handler := createMqttMocks(t, testFunc)

	err := handler.PublishProfileEvent(profileEvent)
	assert.Nil(t, err)
}

func TestMqttPublishProfileEventWithError(t *testing.T) {
	profileEvent := model.NewProfileEvent(null.IntFrom(123), model.Save)
	expectedPayload := testutils.JsonEncode(t, profileEvent)

	testFunc := createTestFunc(t, testConfig.ProfileTopic, expectedPayload, errReturn)
	handler := createMqttMocks(t, testFunc)

	actualError := handler.PublishProfileEvent(profileEvent)
	assert.Equal(t, errExpectedPublish, actualError)
}

func TestMqttPublishStripEvent(t *testing.T) {
	stripEvent := model.NewStripEvent(null.IntFrom(123), model.Delete)
	expectedPayload := testutils.JsonEncode(t, stripEvent)

	testFunc := createTestFunc(t, testConfig.StripTopic, expectedPayload, nil)
	handler := createMqttMocks(t, testFunc)

	err := handler.PublishStripEvent(stripEvent)
	assert.Nil(t, err)
}

func TestMqttPublishStripEventWithError(t *testing.T) {
	stripEvent := model.NewStripEvent(null.IntFrom(123), model.Delete)
	expectedPayload := testutils.JsonEncode(t, stripEvent)
	testFunc := createTestFunc(t, testConfig.StripTopic, expectedPayload, errReturn)
	handler := createMqttMocks(t, testFunc)

	actualError := handler.PublishStripEvent(stripEvent)
	assert.Equal(t, errExpectedPublish, actualError)
}

func TestMqttClose(t *testing.T) {
	handler := createMqttMocks(t, nil)
	handler.Close()
	assert.False(t, handler.mqclient.IsConnected())
}

func createTestFunc(t *testing.T, expectedTopic string, expectedPayload string, returnError error) publishFunc {
	return func(topic string, payload interface{}) error {
		assert.Equal(t, expectedTopic, topic)
		actualPayloadBytes := payload.([]byte)
		assert.Equal(t, expectedPayload, string(actualPayloadBytes))
		return returnError
	}
}

func createMqttMocks(t *testing.T, pubFunc publishFunc) *mqttHandler {
	fake := &mqttClientFake{
		PublishCheck: pubFunc,
	}
	handler := NewMQTT(testConfig)
	handler.intialized = true
	handler.mqclient = fake
	return handler
}

func (c *mqttClientFake) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	fmt.Println("PUBLISH HANDLER CALLED")
	err := c.PublishCheck(topic, payload)
	return mqttTokenFake{
		err: err,
	}
}

func (c *mqttClientFake) Disconnect(quiesce uint) {
	c.disconnected = true
}
func (c *mqttClientFake) IsConnected() bool {
	return !c.disconnected
}

func (t mqttTokenFake) Wait() bool {
	return t.err != nil
}

func (t mqttTokenFake) Error() error {
	return t.err
}
