package influxdb

import (
	"fmt"
	"time"

	"github.com/goNfCollector/common"
)

// write device measurement
func (i *InfluxDBv2) measureDevice(metrics []common.Metric) {

	i.WaitGroup.Add(1)
	defer i.WaitGroup.Done()

	// define new write api
	wapi := i.client.WriteAPI(i.Org, i.Bucket)

	// errorsCh := wapi.Errors()
	// // Create go proc for reading and logging errors
	// go func() {
	// 	for err := range errorsCh {
	// 		log.Printf("influxDB write error: %s\n", err.Error())
	// 	}
	// }()

	for _, m := range metrics {
		protoLine := fmt.Sprintf("device,ver=%v,device=%v bytes=%vu,packets=%vu %v",
			m.FlowVersion,
			m.Device,
			m.Bytes,
			m.Packets,
			time.Now().Add(-time.Duration(m.Time.Second())).UnixNano(),
		)
		// fmt.Println("====", protoLine)

		// write proto line records
		wapi.WriteRecord(protoLine)

	}

	// write to influx
	wapi.Flush()
}
