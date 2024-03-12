package metricStruct

import (
	"time"
)

type NomadLabels map[string]string

type NomadCounter struct {
	Count   int         `json:"Count"`
	CLabels NomadLabels `json:"Labels"`
	Max     float64     `json:"Max"`
	Mean    float64     `json:"Mean"`
	Min     float64     `json:"Min"`
	Name    string      `json:"Name"`
	Rate    float64     `json:"Rate"`
	Stddev  float64     `json:"Stddev"`
	Sum     float64     `json:"Sum"`
}

type NomadGauges struct {
	GLabels NomadLabels `json:"Labels"`
	Name    string      `json:"Name"`
	Value   float64     `json:"Value"`
}

type NomadMetrics struct {
	Counters  []NomadCounter `json:"Counters"`
	Gauges    []NomadGauges  `json:"Gauges"`
	Timestamp string         `json:"Timestamp"`
}

// Clickhouse schema based struct for bulk insert into db
type ClickHouseSchema struct {
	Source     string
	MetricType string
	StatName   string
	MetricName string
	Timestamp  time.Time
	StatValue  float64
	Labels     map[string]string
}