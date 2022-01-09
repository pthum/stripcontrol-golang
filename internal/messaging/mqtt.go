package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/internal/config"
	"github.com/pthum/stripcontrol-golang/internal/database"
	"github.com/pthum/stripcontrol-golang/internal/models"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// StripEvent a StripEvent
type StripEvent struct {
	Type  string   `json:"type,omitempty"`
	ID    null.Int `json:"id,omitempty"`
	Strip OptStrip `json:"state,omitempty"`
}

// OptStrip Optional Strip
type OptStrip struct {
	Valid bool
	Strip struct {
		ID      int64      `json:"id,omitempty"`
		Name    string     `json:"name,omitempty"`
		Enabled bool       `json:"enabled,omitempty"`
		MisoPin int64      `json:"misoPin,omitempty"`
		SclkPin int64      `json:"sclkPin,omitempty"`
		NumLeds int64      `json:"numLeds,omitempty"`
		SpeedHz int64      `json:"speedHz,omitempty"`
		Profile OptProfile `json:"profile,omitempty"`
	}
}

//OptProfile Optional Profile
type OptProfile struct {
	Profile models.ColorProfile
	Valid   bool
}

// ProfileEvent a StripEvent
type ProfileEvent struct {
	Type  string     `json:"type,omitempty"`
	ID    null.Int   `json:"id,omitempty"`
	State OptProfile `json:"state,omitempty"`
}

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

// MQClient the mosquitto client
var (
	mqclient MQTT.Client
)

// Init initialize messaging
func Init() {
	if config.CONFIG.Messaging.Disabled {
		return
	}
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	configString := fmt.Sprintf("tcp://%s:%s", config.CONFIG.Messaging.Host, config.CONFIG.Messaging.Port)
	opts := MQTT.NewClientOptions().AddBroker(configString)
	opts.SetClientID("stripcontrol-go")
	opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
	mqclient = MQTT.NewClient(opts)
	if token := mqclient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Print("initialized messaging")
}

// Close closes connections to message broker
func Close() {
	if config.CONFIG.Messaging.Disabled {
		return
	}
	mqclient.Disconnect(100)
	log.Print("message broker connection gracefully closed")
}

// PublishStripSaveEvent publishes a strip save event
func PublishStripSaveEvent(id null.Int, strip models.LedStrip) (err error) {
	var event = createStripEvent(id, strip)
	err = PublishStripEvent(event)
	return
}

// PublishStripDeleteEvent publishes a strip save event
func PublishStripDeleteEvent(id null.Int) (err error) {
	var event = createDeleteStripEvent(id)
	err = PublishStripEvent(event)
	return
}

// PublishStripEvent publishes a strip event
func PublishStripEvent(event StripEvent) (err error) {
	err = publish(config.CONFIG.Messaging.StripTopic, event)
	return
}
func publish(topic string, event interface{}) (err error) {
	// TODO async
	if mqclient == nil {
		err = fmt.Errorf("Not initialized")
		return
	}
	data, err2 := json.Marshal(event)
	if err2 != nil {
		log.Printf("Error %s", err2.Error())
	}
	log.Printf("sending to topic %s event: %s", topic, string(data))
	if token := mqclient.Publish(topic, 0, false, data); token.Wait() && token.Error() != nil {
		log.Printf("error: %s", token.Error().Error())
		err = fmt.Errorf("Failed to send message")
	}
	return
}

// PublishProfileSaveEvent publishes a profile save event
func PublishProfileSaveEvent(id null.Int, profile models.ColorProfile) (err error) {
	var event = createProfileEvent(id, profile)
	err = publishProfileEvent(event)
	return
}

// PublishProfileDeleteEvent publishes a profile delete event
func PublishProfileDeleteEvent(id null.Int) (err error) {
	var event = createDeleteProfileEvent(id)
	err = publishProfileEvent(event)
	return
}

// PublishProfileEvent publishes a profile event
func publishProfileEvent(event ProfileEvent) (err error) {
	err = publish(config.CONFIG.Messaging.ProfileTopic, event)
	return
}

// CreateStripEvent creates a strip event
func createStripEvent(id null.Int, strip models.LedStrip) (event StripEvent) {
	event = StripEvent{
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
		var prof, _ = database.GetColorProfile(strconv.FormatInt(strip.ProfileID.Int64, 10))
		event.Strip.Strip.Profile.Valid = true
		event.Strip.Strip.Profile.Profile = prof
	}
	return
}

// createDeleteStripEvent creates a strip event
func createDeleteStripEvent(id null.Int) (event StripEvent) {
	event = StripEvent{
		Type: "DELETE",
		ID:   id,
	}
	event.Strip.Valid = false
	return
}

// createProfileEvent creates a profile event
func createProfileEvent(id null.Int, profile models.ColorProfile) (event ProfileEvent) {
	event = ProfileEvent{
		Type: "SAVE",
		ID:   id,
	}
	event.State.Profile = profile
	event.State.Valid = true
	return
}

// createDeleteProfileEvent creates a strip event
func createDeleteProfileEvent(id null.Int) (event ProfileEvent) {
	event = ProfileEvent{
		Type: "DELETE",
		ID:   id,
	}
	event.State.Valid = false
	return
}

// MarshalJSON marshals json for OptProfile
func (profile OptProfile) MarshalJSON() (data []byte, err error) {
	if profile.Valid {
		data, err = json.Marshal(profile.Profile)
		return
	}
	err = nil
	data = []byte("null")
	return
}

// MarshalJSON marshals json for OptProfile
func (strip OptStrip) MarshalJSON() (data []byte, err error) {
	if strip.Valid {
		data, err = json.Marshal(strip.Strip)
		return
	}
	err = nil
	data = []byte("null")
	return
}
