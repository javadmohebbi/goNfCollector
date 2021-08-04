package influxdb

import (
	"fmt"
	"time"

	"github.com/goNfCollector/common"
	"github.com/ip2location/ip2location-go"
)

// write host, port & protocols, ... measurement
// kind can be src or dst
func (i *InfluxDBv2) measureSrcDstRelatedMetrics(metrics []common.Metric, kind string) {

	i.WaitGroup.Add(1)
	defer i.WaitGroup.Done()

	// check for invalid host
	if kind != "src" && kind != "dst" {
		return
	}

	// define new write api
	wapi := i.client.WriteAPI(i.Org, i.Bucket)

	for _, m := range metrics {

		t := time.Now().Add(-time.Duration(m.Time.Second())).UnixNano()

		// ip2location recored
		var i2l *ip2location.IP2Locationrecord

		// host in metrics
		var host, port string

		// check for src or dst host
		if kind == "src" {
			host = m.SrcIP
			port = m.SrcPortName

		} else {
			host = m.DstIP
			port = m.DstPortName

		}

		// get location information for host
		i2l = i.getLocation(host)

		// protoLine := fmt.Sprintf("%vHost,device=%v,ASN=%v,host=%v,countryLong=%v,countryShort=%v,region=%v,city=%v bytes=%vu,packets=%vu %v",
		// 	kind,
		// 	m.Device,
		// 	i.asnLookup(host),
		// 	host,
		// 	i2l.Country_long,
		// 	i2l.Country_short,
		// 	i2l.Region,
		// 	i2l.City,
		// 	m.Bytes, m.Packets,
		// 	t,
		// )

		protoLine := fmt.Sprintf("%vHost,device=%v,host=%v,countryLong=%v,countryShort=%v,region=%v,city=%v bytes=%vu,packets=%vu %v",
			kind,
			m.Device,
			host,
			i2l.Country_long,
			i2l.Country_short,
			i2l.Region,
			i2l.City,
			m.Bytes, m.Packets,
			t,
		)

		// src & dst port protoline
		protoLinePort := fmt.Sprintf("%vPort,device=%v,port=%v,countryLong=%v,countryShort=%v,region=%v,city=%v bytes=%vu,packets=%vu %v",
			kind,
			m.Device,
			port,
			i2l.Country_long,
			i2l.Country_short,
			i2l.Region,
			i2l.City,
			m.Bytes, m.Packets,
			t,
		)

		// protocol protoline
		protoLineProtocol := fmt.Sprintf("%vProtocol,device=%v,proto=%v,countryLong=%v,countryShort=%v,region=%v,city=%v bytes=%vu,packets=%vu %v",
			kind,
			m.Device,
			m.ProtoName,
			i2l.Country_long,
			i2l.Country_short,
			i2l.Region,
			i2l.City,
			m.Bytes, m.Packets,

			t,
		)

		// DNS REverse Looup
		protoLineReverseLookup := fmt.Sprintf("%vDnsLookup,device=%v,host=%v,domain=%v,countryLong=%v,countryShort=%v,region=%v,city=%v bytes=%vu,packets=%vu %v",
			kind,
			m.Device,
			host,
			i.revereseDNS(host),
			i2l.Country_long,
			i2l.Country_short,
			i2l.Region,
			i2l.City,
			m.Bytes, m.Packets,
			t,
		)

		// ASN Name
		// protoLineASN := fmt.Sprintf("%vAS,device=%v,as=%v,countryLong=%v,countryShort=%v,region=%v,city=%v bytes=%vu,packets=%vu %v",
		// 	kind,
		// 	m.Device,
		// 	i.asnLookup(host),
		// 	i2l.Country_long,
		// 	i2l.Country_short,
		// 	i2l.Region,
		// 	i2l.City,
		// 	m.Bytes, m.Packets,
		// 	t,
		// )

		// write proto line records
		wapi.WriteRecord(protoLine)

		// for ports
		wapi.WriteRecord(protoLinePort)

		// for protocols
		wapi.WriteRecord(protoLineProtocol)

		// reverse dns lookup
		wapi.WriteRecord(protoLineReverseLookup)

		// asn
		// wapi.WriteRecord(protoLineASN)

		// host reputation
		for _, rpu := range i.reputations {
			if resp := rpu.Get(host); resp.Current > 0 {
				protoLineReputation := fmt.Sprintf("reputation,device=%v,host=%v,source=%v,countryLong=%v,countryShort=%v,region=%v,city=%v score=%v %v",
					m.Device,
					host,
					rpu.GetType(),
					i2l.Country_long,
					i2l.Country_short,
					i2l.Region,
					i2l.City,
					resp.Current,
					time.Now().Add(-time.Duration(m.Time.Second())).UnixNano(),
				)
				// write reputation
				wapi.WriteRecord(protoLineReputation)
			}
		}

	}

	// write to influx
	wapi.Flush()
}
