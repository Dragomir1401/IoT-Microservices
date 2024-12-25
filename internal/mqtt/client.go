package mqtt

import (
	"github.com/eclipse/paho.mqtt.golang"
)

// MQTTClient wraps the Paho MQTT client
type Client struct {
	client mqtt.Client
}

// NewMQTTClient creates a new MQTT client
func NewMQTTClient(broker string) (*Client, error) {
	opts := mqtt.NewClientOptions().AddBroker(broker)
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &Client{client: client}, nil
}

// Subscribe subscribes to a given topic
func (m *Client) Subscribe(topic string, callback func(topic string, payload []byte)) error {
	if token := m.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		callback(msg.Topic(), msg.Payload())
	}); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Publish sends a message to a topic
func (m *Client) Publish(topic string, payload string) error {
	if token := m.client.Publish(topic, 0, false, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Disconnect disconnects the client
func (m *Client) Disconnect() {
	m.client.Disconnect(250)
}
