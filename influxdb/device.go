package influxdb

import (
	"fmt"

	"github.com/goNfCollector/common"
)

// write device measurement
func (i *InfluxDBv2) measureDevice(metrics []common.Metric) {
	// define new write api
	wapi := i.client.WriteAPI(i.Org, i.Bucket)

	for _, m := range metrics {
		protoLine := fmt.Sprintf("device,ver=%v,device=%v bytes=%vu,packets=%vu %v",
			m.FlowVersion,
			m.Device,
			m.Bytes, m.Packets,
			m.Time.UnixNano(),
		)
		// write proto line records
		wapi.WriteRecord(protoLine)
	}

	// write to influx
	wapi.Flush()
}
