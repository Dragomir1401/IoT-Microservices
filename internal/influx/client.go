package influx

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxClient struct {
	client influxdb2.Client
	org    string
	bucket string
}

func NewInfluxClient(url, token, org, bucket string) *InfluxClient {
	client := influxdb2.NewClient(url, token)
	return &InfluxClient{
		client: client,
		org:    org,
		bucket: bucket,
	}
}

func (i *InfluxClient) WriteData(metric string, tags map[string]string, fields map[string]interface{}, timestamp time.Time) error {
	writeAPI := i.client.WriteAPIBlocking(i.org, i.bucket)
	point := influxdb2.NewPoint(metric, tags, fields, timestamp)
	return writeAPI.WritePoint(context.Background(), point)
}

func (i *InfluxClient) Close() {
	i.client.Close()
}
