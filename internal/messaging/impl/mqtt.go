package messagingimpl

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pthum/stripcontrol-golang/internal/config"
	"github.com/pthum/stripcontrol-golang/internal/model"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var once sync.Once

type mqttHandler struct {
	mqclient   mqtt.Client
	opts       *mqtt.ClientOptions
	cfg        config.MessagingConfig
	intialized bool
}

// define a function for the default message handler
var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("TOPIC: %s new message: %s\n", msg.Topic(), msg.Payload())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}

var reconnectHandler mqtt.ReconnectHandler = func(c mqtt.Client, co *mqtt.ClientOptions) {
	log.Printf("Trying to reconnect")
}

func NewMQTT(cfg config.MessagingConfig) *mqttHandler {
	return &mqttHandler{
		opts: buildClientOpts(cfg),
		cfg:  cfg,
	}
}

func (m *mqttHandler) getClient() mqtt.Client {
	if !m.intialized {
		once.Do(func() {
			//create and start a client using the above ClientOptions
			m.mqclient = mqtt.NewClient(m.opts)
			if token := m.mqclient.Connect(); token.Wait() && token.Error() != nil {
				panic(token.Error())
			}
			log.Print("initialized messaging")
			m.intialized = true
		})
	}
	return m.mqclient
}

func buildClientOpts(cfg config.MessagingConfig) *mqtt.ClientOptions {
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	configString := fmt.Sprintf("tcp://%s:%s", cfg.Host, cfg.Port)
	opts := mqtt.NewClientOptions().AddBroker(configString)
	opts.SetClientID("stripcontrol-go")
	opts.SetDefaultPublishHandler(f)
	opts.SetConnectionLostHandler(connectLostHandler)
	opts.SetOnConnectHandler(connectHandler)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(1 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(60 * time.Second)
	opts.SetReconnectingHandler(reconnectHandler)
	return opts
}

// Close closes connections to message broker
func (m *mqttHandler) Shutdown() error {
	m.mqclient.Disconnect(100)
	log.Print("message broker connection gracefully closed")
	return nil
}

// PublishStripEvent publishes a strip event
func (m *mqttHandler) PublishStripEvent(event *model.StripEvent) error {
	return m.publish(m.cfg.StripTopic, event)
}

// PublishProfileEvent publishes a profile event
func (m *mqttHandler) PublishProfileEvent(event *model.ProfileEvent) error {
	return m.publish(m.cfg.ProfileTopic, event)
}

func (m *mqttHandler) publish(topic string, event interface{}) (err error) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error %s", err.Error())
		return
	}
	log.Printf("sending to topic %s event: %s", topic, string(data))
	token := m.getClient().Publish(topic, 0, false, data)
	if token.Wait() && token.Error() != nil {
		log.Printf("error: %s", token.Error().Error())
		err = errors.New("failed to send message")
	}
	return
}
