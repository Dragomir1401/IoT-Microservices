package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client mqtt.Client
}

func NewMQTTClient(broker string) (*MQTTClient, error) {
	opts := mqtt.NewClientOptions().AddBroker(broker)
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &MQTTClient{client: client}, nil
}

func (m *MQTTClient) Subscribe(topic string, callback func(topic string, payload []byte)) error {
	token := m.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		callback(msg.Topic(), msg.Payload())
	})
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MQTTClient) Publish(topic string, payload string) error {
	token := m.client.Publish(topic, 0, false, payload)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MQTTClient) Disconnect() {
	m.client.Disconnect(250)
}
