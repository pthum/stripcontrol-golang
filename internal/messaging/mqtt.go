package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/internal/config"
	"github.com/pthum/stripcontrol-golang/internal/database"
	"github.com/pthum/stripcontrol-golang/internal/model"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type EventHandler interface {
	Close()
	PublishStripSaveEvent(id null.Int, strip model.LedStrip) (err error)
	PublishStripDeleteEvent(id null.Int) (err error)
	PublishStripEvent(event model.StripEvent) (err error)
	PublishProfileSaveEvent(id null.Int, profile model.ColorProfile) (err error)
	PublishProfileDeleteEvent(id null.Int) (err error)
}

type MQTTHandler struct {
	mqclient MQTT.Client
	dbr      database.DBReader
}

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func New() EventHandler {
	// return NoOp implementation if disabled
	if config.CONFIG.Messaging.Disabled {
		return &NoOpEventHandler{}
	}

	return &MQTTHandler{
		mqclient: initital(),
	}
}

// Init initialize messaging
func initital() MQTT.Client {
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	configString := fmt.Sprintf("tcp://%s:%s", config.CONFIG.Messaging.Host, config.CONFIG.Messaging.Port)
	opts := MQTT.NewClientOptions().AddBroker(configString)
	opts.SetClientID("stripcontrol-go")
	opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
	mqclient := MQTT.NewClient(opts)
	if token := mqclient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Print("initialized messaging")

	return mqclient
}

// Close closes connections to message broker
func (m *MQTTHandler) Close() {
	m.mqclient.Disconnect(100)
	log.Print("message broker connection gracefully closed")
}

// PublishStripSaveEvent publishes a strip save event
func (m *MQTTHandler) PublishStripSaveEvent(id null.Int, strip model.LedStrip) (err error) {
	var event = m.createStripEvent(id, strip)
	err = m.PublishStripEvent(event)
	return
}

// PublishStripDeleteEvent publishes a strip save event
func (m *MQTTHandler) PublishStripDeleteEvent(id null.Int) (err error) {
	var event = m.createDeleteStripEvent(id)
	err = m.PublishStripEvent(event)
	return
}

// PublishStripEvent publishes a strip event
func (m *MQTTHandler) PublishStripEvent(event model.StripEvent) (err error) {
	err = m.publish(config.CONFIG.Messaging.StripTopic, event)
	return
}
func (m *MQTTHandler) publish(topic string, event interface{}) (err error) {
	// TODO async
	if m.mqclient == nil {
		err = fmt.Errorf("Not initialized")
		return
	}
	data, err2 := json.Marshal(event)
	if err2 != nil {
		log.Printf("Error %s", err2.Error())
	}
	log.Printf("sending to topic %s event: %s", topic, string(data))
	if token := m.mqclient.Publish(topic, 0, false, data); token.Wait() && token.Error() != nil {
		log.Printf("error: %s", token.Error().Error())
		err = fmt.Errorf("Failed to send message")
	}
	return
}

// PublishProfileSaveEvent publishes a profile save event
func (m *MQTTHandler) PublishProfileSaveEvent(id null.Int, profile model.ColorProfile) (err error) {
	var event = m.createProfileEvent(id, profile)
	err = m.publishProfileEvent(event)
	return
}

// PublishProfileDeleteEvent publishes a profile delete event
func (m *MQTTHandler) PublishProfileDeleteEvent(id null.Int) (err error) {
	var event = m.createDeleteProfileEvent(id)
	err = m.publishProfileEvent(event)
	return
}

// PublishProfileEvent publishes a profile event
func (m *MQTTHandler) publishProfileEvent(event model.ProfileEvent) (err error) {
	err = m.publish(config.CONFIG.Messaging.ProfileTopic, event)
	return
}

// CreateStripEvent creates a strip event
func (m *MQTTHandler) createStripEvent(id null.Int, strip model.LedStrip) (event model.StripEvent) {
	event = model.StripEvent{
		Type: "SAVE",
		ID:   id,
	}
	event.Strip.Valid = true
	event.Strip.Strip.ID = strip.ID
	event.Strip.Strip.Name = strip.Name
	event.Strip.Strip.Enabled = strip.Enabled
	event.Strip.Strip.MisoPin = strip.MisoPin.Int64
	event.Strip.Strip.SclkPin = strip.SclkPin.Int64
	event.Strip.Strip.NumLeds = strip.NumLeds.Int64
	event.Strip.Strip.SpeedHz = strip.SpeedHz.Int64
	if strip.ProfileID.Valid {
		var prof model.ColorProfile
		var _ = m.dbr.Get(strconv.FormatInt(strip.ProfileID.Int64, 10), &prof)
		event.Strip.Strip.Profile.Valid = true
		event.Strip.Strip.Profile.Profile = prof
	}
	return
}

// createDeleteStripEvent creates a strip event
func (m *MQTTHandler) createDeleteStripEvent(id null.Int) (event model.StripEvent) {
	event = model.StripEvent{
		Type: "DELETE",
		ID:   id,
	}
	event.Strip.Valid = false
	return
}

// createProfileEvent creates a profile event
func (m *MQTTHandler) createProfileEvent(id null.Int, profile model.ColorProfile) (event model.ProfileEvent) {
	event = model.ProfileEvent{
		Type: "SAVE",
		ID:   id,
	}
	event.State.Profile = profile
	event.State.Valid = true
	return
}

// createDeleteProfileEvent creates a strip event
func (m *MQTTHandler) createDeleteProfileEvent(id null.Int) (event model.ProfileEvent) {
	event = model.ProfileEvent{
		Type: "DELETE",
		ID:   id,
	}
	event.State.Valid = false
	return
}
