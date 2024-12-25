package mqtt

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMQTTClient(t *testing.T) {
	// Start a test MQTT broker
	broker, err := mosquitto.RunContainer(context.Background())
	assert.NoError(t, err)
	defer broker.Terminate(context.Background())

	brokerHost, err := broker.Host(context.Background())
	assert.NoError(t, err)
	brokerPort, err := broker.MappedPort(context.Background(), "1883")
	assert.NoError(t, err)

	brokerURL := "tcp://" + brokerHost + ":" + brokerPort.Port()

	// Create the MQTT client
	client, err := NewMQTTClient(brokerURL)
	assert.NoError(t, err)
	defer client.Disconnect()

	// Test publishing and subscribing
	messageReceived := make(chan string, 1)
	err = client.Subscribe("test/topic", func(topic string, payload []byte) {
		messageReceived <- string(payload)
	})
	assert.NoError(t, err)

	err = client.Publish("test/topic", "hello, world")
	assert.NoError(t, err)

	select {
	case msg := <-messageReceived:
		assert.Equal(t, "hello, world", msg)
	case <-time.After(2 * time.Second):
		t.Fatal("Message not received")
	}
}
