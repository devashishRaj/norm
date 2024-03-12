package norm

import (
	"context"
	"encoding/json"
	"fmt"
	chDriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/devashishRaj/norm/CH_conn"
	"github.com/devashishRaj/norm/metricStruct"
	"github.com/hashicorp/nomad/api"
	"log"
	"os"
	"sync"
	"time"
)

type WorkerParams struct {
	AppendData  chan<- []metricStruct.ClickHouseSchema
	FetchedData metricStruct.NomadMetrics
	ChConn      chDriver.Conn
}

// FetchMetrics : fetch nomad metrics in pretty format and returns a json file
// as per metricStruct.NomadMetrics struct
func FetchMetrics() metricStruct.NomadMetrics {

	cfg := api.DefaultConfig()
	cfg.Address = "http://127.0.0.1:4646"
	c, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	op := c.Operator()
	qo := &api.QueryOptions{
		Params: map[string]string{
			"pretty": "1",
		},
	}
	metrics, err := op.Metrics(qo)
	if err != nil {
		panic(err)
	}
	if metrics == nil {
		panic(err)
	}
	//dataToFile(metrics)
	var nomadmetrics metricStruct.NomadMetrics
	err = json.Unmarshal(metrics, &nomadmetrics)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return metricStruct.NomadMetrics{}
	}
	return nomadmetrics

}

// NomadmetricsBulksend : collects metrics at interval of one second
// and after 10th collection it sends them into clickhouse
func NomadmetricsBulksend() int {
	conn, err := CH_conn.ConnectCH()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	bulkDataChannel := make(chan []metricStruct.ClickHouseSchema)
	var wg sync.WaitGroup
	count := 0
	// inside go routine to make other parts reachable
	wg.Add(1)
	go func() {
		for {
			metrics := FetchMetrics()
			params := WorkerParams{
				AppendData:  bulkDataChannel,
				FetchedData: metrics,
				ChConn:      conn,
			}
			go BatchBuild(params)
			count++
			// collect metrics at interval of 1 second and send them in channel
			time.Sleep(1 * time.Second)
			// once metrics has been collected 10 times , they are send
			if count == 10 {
				fmt.Println("In collection")
				var bulkBatch []metricStruct.ClickHouseSchema
				// // creates a single slice by joining all 10 slices
				for ; count > 0; count-- {
					fmt.Println(count)
					singleBatch := <-bulkDataChannel
					bulkBatch = append(bulkBatch, singleBatch...)
				}
				BulkSend(bulkBatch, conn)
				fmt.Println("Bulk sent")
			}

		}
	}()
	wg.Wait()
	return 0
}

// BatchBuild : it creates a slice of a struct metricStruct.ClickHouseSchema
// using metrics collected from FetchMetrics()
func BatchBuild(params WorkerParams) {
	fmt.Println("Creating build batch")
	var bulkBatch []metricStruct.ClickHouseSchema
	for _, counter := range params.FetchedData.Counters {
		bulkBatch = append(bulkBatch, metricStruct.ClickHouseSchema{

			Source:     "HTTP METRIC API",
			MetricType: "Counter",
			StatName:   "Max",
			MetricName: counter.Name,
			Timestamp:  time.Now(),
			StatValue:  counter.Max,
			Labels:     counter.CLabels,
		})
		bulkBatch = append(bulkBatch, metricStruct.ClickHouseSchema{
			Source:     "HTTP METRIC API",
			MetricType: "Counter",
			StatName:   "Min",
			MetricName: counter.Name,
			Timestamp:  time.Now(),
			StatValue:  counter.Min,
			Labels:     counter.CLabels,
		})
		bulkBatch = append(bulkBatch, metricStruct.ClickHouseSchema{
			Source:     "HTTP METRIC API",
			MetricType: "Counter",
			StatName:   "Mean",
			MetricName: counter.Name,
			Timestamp:  time.Now(),
			StatValue:  counter.Mean,
			Labels:     counter.CLabels,
		})

		bulkBatch = append(bulkBatch, metricStruct.ClickHouseSchema{
			Source:     "HTTP METRIC API",
			MetricType: "Counter",
			StatName:   "Stddev",
			MetricName: counter.Name,
			Timestamp:  time.Now(),
			StatValue:  counter.Stddev,
			Labels:     counter.CLabels,
		})

		bulkBatch = append(bulkBatch, metricStruct.ClickHouseSchema{
			Source:     "HTTP METRIC API",
			MetricType: "Counter",
			StatName:   "Rate",
			MetricName: counter.Name,
			Timestamp:  time.Now(),
			StatValue:  counter.Rate,
			Labels:     counter.CLabels,
		})
		bulkBatch = append(bulkBatch, metricStruct.ClickHouseSchema{
			Source:     "HTTP METRIC API",
			MetricType: "Counter",
			StatName:   "Sum",
			MetricName: counter.Name,
			Timestamp:  time.Now(),
			StatValue:  counter.Sum,
			Labels:     counter.CLabels,
		})

	}
	for _, gauge := range params.FetchedData.Gauges {

		bulkBatch = append(bulkBatch, metricStruct.ClickHouseSchema{
			Source:     "HTTP METRIC API",
			MetricType: "gauge",
			StatName:   "Value",
			MetricName: gauge.Name,
			Timestamp:  time.Now(),
			StatValue:  gauge.Value,
			Labels:     gauge.GLabels,
		})
	}
	params.AppendData <- bulkBatch
	fmt.Println("Build done")
	//defer params.WaitGroup.Done()

}

// BulkSend :  insert slice of metricStruct.ClickHouseSchema
func BulkSend(bulkdata []metricStruct.ClickHouseSchema, conn chDriver.Conn) {
	fmt.Println("In bulk Send")

	ctx := context.Background()
	batch, err := conn.PrepareBatch(ctx, "INSERT INTO nomad.metrics")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, data := range bulkdata {
		err = batch.AppendStruct(&metricStruct.ClickHouseSchema{
			Source:     data.Source,
			MetricType: data.MetricType,
			StatName:   data.StatName,
			MetricName: data.MetricName,
			Timestamp:  data.Timestamp,
			StatValue:  data.StatValue,
			Labels:     data.Labels,
		})
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Sending....")
	err = batch.Send()
	if err != nil {
		panic(err)
	}

}

func dataToFile(body []byte) {

	file, err := os.Create("data.json")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	_, err = file.Write(body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("JSON data successfully written to file.")
}