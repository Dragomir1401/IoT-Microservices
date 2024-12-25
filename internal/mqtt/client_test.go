package mqtt

import (
	"testing"
)

func TestMQTTClient(t *testing.T) {
	client, err := NewMQTTClient("tcp://localhost:1883")
	if err != nil {
		t.Fatalf("Failed to connect to MQTT broker: %v", err)
	}
	defer client.Disconnect()

	// Test subscribe and publish
	err = client.Subscribe("test/topic", func(topic string, payload []byte) {
		t.Logf("Received message: %s from topic: %s", payload, topic)
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	err = client.Publish("test/topic", "test message")
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}
}
