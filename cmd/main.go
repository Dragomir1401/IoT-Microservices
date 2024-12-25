package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/mqtt-adaptor/internal/influx"
	"example.com/mqtt-adaptor/internal/mqtt"
)

type SensorData struct {
	BAT       float64 `json:"BAT"`
	HUMID     float64 `json:"HUMID"`
	TMP       float64 `json:"TMP"`
	Timestamp string  `json:"timestamp"`
}

func main() {
	mqttBroker := os.Getenv("MQTT_BROKER")
	influxURL := os.Getenv("INFLUXDB_URL")
	influxToken := os.Getenv("INFLUXDB_TOKEN")
	influxOrg := os.Getenv("INFLUXDB_ORG")
	influxBucket := os.Getenv("INFLUXDB_BUCKET")
	debug := os.Getenv("DEBUG_DATA_FLOW") == "true"

	if mqttBroker == "" || influxURL == "" || influxToken == "" || influxOrg == "" || influxBucket == "" {
		log.Fatal("Missing required environment variables")
	}

	mqttClient, err := mqtt.NewMQTTClient(mqttBroker)
	if err != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", err)
	}
	defer mqttClient.Disconnect()

	influxClient := influx.NewInfluxClient(influxURL, influxToken, influxOrg, influxBucket)
	defer influxClient.Close()

	messageHandler := func(topic string, payload []byte) {
		if debug {
			log.Printf("Received message on topic [%s]: %s", topic, payload)
		}

		var data SensorData
		if err := json.Unmarshal(payload, &data); err != nil {
			log.Printf("Failed to parse message: %v", err)
			return
		}

		timestamp := time.Now()
		if data.Timestamp != "" {
			parsedTime, err := time.Parse(time.RFC3339, data.Timestamp)
			if err == nil {
				timestamp = parsedTime
			} else {
				log.Printf("Invalid timestamp in message, using current time: %v", err)
			}
		}

		tags := map[string]string{"topic": topic}
		fields := map[string]interface{}{
			"BAT":   data.BAT,
			"HUMID": data.HUMID,
			"TMP":   data.TMP,
		}

		if err := influxClient.WriteData("sensor_data", tags, fields, timestamp); err != nil {
			log.Printf("Failed to write data to InfluxDB: %v", err)
		} else if debug {
			log.Printf("Data written to InfluxDB: %v", fields)
		}
	}

	if err := mqttClient.Subscribe("#", messageHandler); err != nil {
		log.Fatalf("Failed to subscribe to MQTT topics: %v", err)
	}

	log.Println("MQTT adaptor is running...")

	// Așteaptă semnal de oprire (Ctrl+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down MQTT adaptor...")
}
