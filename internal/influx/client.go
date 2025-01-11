package influx

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type InfluxClient struct {
	url      string
	database string
}

// NewInfluxClient creates a new InfluxDB client for InfluxDB 1.x
func NewInfluxClient(url, database string) *InfluxClient {
	return &InfluxClient{
		url:      url,
		database: database,
	}
}

// WriteData writes data to InfluxDB 1.x
func (i *InfluxClient) WriteData(metric string, tags map[string]string, fields map[string]interface{}, timestamp time.Time) error {
	// Construct the line protocol format: "measurement,tagKey=tagValue fieldKey=fieldValue timestamp"
	var tagStrings []string
	for k, v := range tags {
		tagStrings = append(tagStrings, fmt.Sprintf("%s=%s", k, strings.ReplaceAll(v, " ", "\\ ")))
	}
	var fieldStrings []string
	for k, v := range fields {
		fieldStrings = append(fieldStrings, fmt.Sprintf("%s=%v", k, v))
	}
	line := fmt.Sprintf("%s,%s %s %d",
		metric,                          // measurement
		strings.Join(tagStrings, ","),   // tags
		strings.Join(fieldStrings, ","), // fields
		timestamp.UnixNano())            // timestamp

	log.Printf("[DEBUG] Constructed Line Protocol: %s", line)

	// Send the line to InfluxDB
	writeURL := fmt.Sprintf("%s/write?db=%s", i.url, i.database)
	resp, err := http.Post(writeURL, "text/plain", bytes.NewBufferString(line))
	if err != nil {
		return fmt.Errorf("failed to write data to InfluxDB: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusNoContent {
		log.Printf("[ERROR] Failed to write data to InfluxDB. Status Code: %d", resp.StatusCode)
		return fmt.Errorf("failed to write data to InfluxDB: status code %d", resp.StatusCode)
	}

	log.Printf("Data written to InfluxDB: %s", line)
	return nil
}

// Close is a no-op for InfluxDB 1.x
func (i *InfluxClient) Close() {
	log.Println("InfluxDB 1.x client does not require explicit close")
}
