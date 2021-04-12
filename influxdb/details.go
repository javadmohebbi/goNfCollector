package influxdb

import (
	"fmt"
	"time"

	"github.com/goNfCollector/common"
	"github.com/ip2location/ip2location-go"
)

// write host, port & protocols, ... measurement

func (i *InfluxDBv2) measureDetailsRelatedMetrics(metrics []common.Metric) {

	i.WaitGroup.Add(1)
	defer i.WaitGroup.Done()

	// define new write api
	wapi := i.client.WriteAPI(i.Org, i.Bucket)

	for _, m := range metrics {

		t := time.Now().Add(-time.Duration(m.Time.Second())).UnixNano()

		// ip2location recored
		var dstI2l, srcI2l *ip2location.IP2Locationrecord

		// host in metrics
		var srcHost, dstHost, srcPort, dstPort string

		// src
		srcHost = m.SrcIP
		srcPort = m.SrcPortName
		srcI2l = i.getLocation(srcHost)

		// dst
		dstHost = m.DstIP
		dstPort = m.DstPortName
		dstI2l = i.getLocation(dstHost)

		// all in and out
		protoLineHostDetail := fmt.Sprintf("detail,device=%v,proto=%v,sASN=%v,shost=%v,sport=%v,scountryLong=%v,scountryShort=%v,sregion=%v,scity=%v,"+
			"dASN=%v,dhost=%v,dport=%v,dcountryLong=%v,dcountryShort=%v,dregion=%v,dcity=%v "+
			"bytes=%vu,packets=%vu %v",
			m.Device,
			m.ProtoName,

			i.asnLookup(m.SrcIP),
			srcHost,
			srcPort,
			srcI2l.Country_long,
			srcI2l.Country_short,
			srcI2l.Region,
			srcI2l.City,

			i.asnLookup(m.DstIP),
			dstHost,
			dstPort,
			dstI2l.Country_long,
			dstI2l.Country_short,
			dstI2l.Region,
			dstI2l.City,

			m.Bytes, m.Packets,

			t,
		)

		// write proto line records for details
		wapi.WriteRecord(protoLineHostDetail)

	}

	// write to influx
	wapi.Flush()
}
