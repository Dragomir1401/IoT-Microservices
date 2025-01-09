package influx

import (
	"context"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxClient struct {
	client influxdb2.Client
	org    string
	bucket string
}

// NewInfluxClient creates a new InfluxDB client
func NewInfluxClient(url, token, org, bucket string) *InfluxClient {
	log.Printf("Using InfluxDB token: %s", token)

	client := influxdb2.NewClientWithOptions(
		url,
		token,
		influxdb2.DefaultOptions(),
	)
	return &InfluxClient{
		client: client,
		org:    org,
		bucket: bucket,
	}
}

// WriteData writes data to InfluxDB
func (i *InfluxClient) WriteData(metric string, tags map[string]string, fields map[string]interface{}, timestamp time.Time) error {
	writeAPI := i.client.WriteAPIBlocking(i.org, i.bucket)
	point := influxdb2.NewPoint(metric, tags, fields, timestamp)
	return writeAPI.WritePoint(context.Background(), point)
}

// Close closes the InfluxDB client
func (i *InfluxClient) Close() {
	i.client.Close()
}
