package influx

import (
	"testing"
	"time"
)

func TestWriteData(t *testing.T) {
	client := NewInfluxClient("http://localhost:8086", "test-token", "test-org", "test-bucket")
	defer client.Close()

	err := client.WriteData("test_metric", map[string]string{"tag1": "value1"}, map[string]interface{}{"field1": 42.0}, time.Now())
	if err != nil {
		t.Fatalf("Error writing data: %v", err)
	}
}
