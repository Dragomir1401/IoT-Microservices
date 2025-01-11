package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strings"
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

// isZeroOrNaN checks if a float64 value is zero or NaN
func isZeroOrNaN(value float64) bool {
	return value == 0 || value != value
}

func main() {
	// Read environment variables
	mqttBroker := os.Getenv("MQTT_BROKER")
	influxURL := os.Getenv("INFLUXDB_URL")
	//influxToken := os.Getenv("INFLUXDB_TOKEN")
	//influxOrg := os.Getenv("INFLUXDB_ORG")
	//influxBucket := os.Getenv("INFLUXDB_BUCKET")
	influxDatabase := os.Getenv("INFLUXDB_DATABASE")
	debug := os.Getenv("DEBUG_DATA_FLOW") == "true"

	// Force required environment variables
	//if mqttBroker == "" || influxURL == "" || influxToken == "" || influxOrg == "" || influxBucket == "" {
	//	log.Printf("Environment Variables: MQTT_BROKER=%s, INFLUXDB_URL=%s, INFLUXDB_ORG=%s, INFLUXDB_BUCKET=%s",
	//		mqttBroker, influxURL, influxOrg, influxBucket)
	//
	//	log.Fatal("Missing required environment variables")
	//}

	if mqttBroker == "" || influxURL == "" || influxDatabase == "" {
		log.Printf("Environment Variables: MQTT_BROKER=%s, INFLUXDB_URL=%s, INFLUXDB_DATABASE=%s",
			mqttBroker, influxURL, influxDatabase)

		log.Fatal("Missing required environment variables")
	}

	// Create MQTT and InfluxDB clients
	mqttClient, err := mqtt.NewMQTTClient(mqttBroker)
	if err != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", err)
	}
	defer mqttClient.Disconnect()
	influxClient := influx.NewInfluxClient(influxURL, influxDatabase)
	defer influxClient.Close()

	// Define the message handler
	messageHandler := func(topic string, payload []byte) {
		if debug {
			log.Printf("Received message on topic [%s]: %s", topic, payload)
		}

		var data SensorData
		if err := json.Unmarshal(payload, &data); err != nil {
			log.Printf("Failed to parse message: %v", err)
			return
		}

		// Split the topic into measurement and tags
		topicParts := strings.Split(topic, "/")
		if len(topicParts) != 2 {
			log.Printf("Invalid topic format: %s", topic)
			return
		}

		measurement := topicParts[0] // location - "UPB"
		station := topicParts[1]     // sensor - "RPi_1"

		// Parse timestamp or use the current time if missing
		timestamp := time.Now()
		if data.Timestamp != "" {
			parsedTime, err := time.Parse(time.RFC3339, data.Timestamp)
			if err == nil {
				timestamp = parsedTime
			} else {
				log.Printf("Invalid timestamp in message, using current time: %v", err)
			}
		}

		// Tag with the station name
		tags := map[string]string{
			"station": station,
		}
		fields := map[string]interface{}{}

		// Validate and skip data with all invalid or missing values
		if !isZeroOrNaN(data.BAT) {
			fields["BAT"] = data.BAT
		} else {
			log.Printf("Invalid or missing numeric value for BAT, skipping...")
		}

		if !isZeroOrNaN(data.HUMID) {
			fields["HUMID"] = data.HUMID
		} else {
			log.Printf("Invalid or missing numeric value for HUMID, skipping...")
		}

		if !isZeroOrNaN(data.TMP) {
			fields["TMP"] = data.TMP
		} else {
			log.Printf("Invalid or missing numeric value for TMP, skipping...")
		}

		// If no valid fields, skip writing to InfluxDB
		if len(fields) == 0 {
			log.Printf("No valid numeric fields to write for topic [%s]", topic)
			return
		}

		// Write the data to InfluxDB
		if err := influxClient.WriteData(measurement, tags, fields, timestamp); err != nil {
			log.Printf("Failed to write data to InfluxDB: %v", err)
		} else if debug {
			log.Printf("Data written to InfluxDB: Measurement=%s, Tags=%v, Fields=%v", measurement, tags, fields)
		}
	}

	// Subscribe to all topics
	if err := mqttClient.Subscribe("#", messageHandler); err != nil {
		log.Fatalf("Failed to subscribe to MQTT topics: %v", err)
	}

	log.Println("MQTT adaptor is running...")

	// Wait for termination signal using a Go primitive channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down MQTT adaptor...")
}
