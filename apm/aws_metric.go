package apm

import (
	"encoding/json"
	"fmt"
	"time"
)

type Metric struct {
	// Name is the name of this metric.
	Name string
	// Attributes is a map of attributes for this metric.
	Attributes map[string]interface{}
	// AttributesJSON is a json.RawMessage of attributes for this metric. It
	// will only be sent if Attributes is nil.
	AttributesJSON json.RawMessage
	// Value is the value of this metric.
	Value float64
	// Timestamp is the time at which this metric was gathered.  If
	// Timestamp is unset then the Harvester's period start will be used.
	Timestamp time.Time
}

type metricGroup struct {
	// Metrics is the slice of metrics to send with this MetricGroup.
	Metrics []Metric
}

func test() {
	fmt.Println("test")

}