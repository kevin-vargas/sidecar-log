package pubsub

import (
	"fmt"
	"sidecar/configs"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to broker")
}
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

type config struct {
	broker string
	port   string
}
type MQTTI interface {
	Publish(topic string, payload interface{}) error
	Subscribe(topic string) error
	SubscribeWithCB(topic string, cb mqtt.MessageHandler) error
}
type MQTT struct {
	client *mqtt.Client
}

func New() MQTTI {
	client, err := getInstanceMqttclient()
	if err != nil {
		panic(err.Error())
	}
	return &MQTT{
		client: client,
	}
}
func (m *MQTT) Publish(topic string, payload interface{}) error {
	token := (*m.client).Publish(topic, 0, false, payload)
	token.Wait()
	return token.Error()
}
func (m *MQTT) SubscribeWithCB(topic string, cb mqtt.MessageHandler) error {
	token := (*m.client).Subscribe(topic, 0, cb)
	token.Wait()
	return token.Error()
}

func (m *MQTT) Subscribe(topic string) error {
	err := m.SubscribeWithCB(topic, nil)
	return err
}

func getInstanceMqttclient() (*mqtt.Client, error) {
	cfg := configs.Get()
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.MQTT.BROKER, cfg.MQTT.PORT))
	opts.SetClientID(cfg.MQTT.CLIENT.ID)
	opts.SetUsername(cfg.MQTT.CLIENT.USERNAME)
	opts.SetPassword(cfg.MQTT.CLIENT.PASSWORD)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetDefaultPublishHandler(messagePubHandler)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return &client, nil
}
