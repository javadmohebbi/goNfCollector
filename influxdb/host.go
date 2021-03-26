package influxdb

import (
	"fmt"

	"github.com/goNfCollector/common"
	"github.com/ip2location/ip2location-go"
)

// write host measurement
// kind can be src or dst
func (i *InfluxDBv2) measureHost(metrics []common.Metric, kind string) {

	// check for invalid host
	if kind != "src" && kind != "dst" {
		return
	}

	// define new write api
	wapi := i.client.WriteAPI(i.Org, i.Bucket)

	for _, m := range metrics {
		// ip2location recored
		var i2l *ip2location.IP2Locationrecord

		// host in metrics
		var host string

		// check for src or dst host
		if kind == "src" {
			host = m.SrcIP
		} else {
			host = m.DstIP
		}

		// get location information for host
		i2l = i.getLocation(host)

		protoLine := fmt.Sprintf("%vHost,device=%v,host=%v,countryLong=%v,countryShort=%v,region=%v,city=%v bytes=%vu,packets=%vu %v",
			kind,
			m.Device,
			host,
			i2l.Country_long,
			i2l.Country_short,
			i2l.Region,
			i2l.City,
			m.Bytes, m.Packets,
			m.Time.UnixNano(),
		)

		// write proto line records
		wapi.WriteRecord(protoLine)
	}

	// write to influx
	wapi.Flush()
}
