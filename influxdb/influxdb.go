package influxdb

import (
	"fmt"
	"strings"

	"github.com/goNfCollector/common"
	"github.com/ip2location/ip2location-go"
	"github.com/sirupsen/logrus"
)

// write to db
func (i *InfluxDBv2) Write(metrics []common.Metric) error {
	// define new write api
	// wapi := i.client.WriteAPI(i.Org, i.Bucket)

	// for _, m := range metrics {
	// 	// ilNextHop := i.getLocation(m.NextHop)
	// 	ilSrc, ilDst := i.getLocation(m.SrcIP), i.getLocation(m.DstIP)

	// 	// prepare protoline
	// 	protoLine := fmt.Sprintf("flows,ver=%v,device=%v,proto=%v,"+
	// 		"sHost=%v,sPort=%v,dHost=%v,dPort=%v,"+
	// 		"sCountSh=%v,sCountLo=%v,sRegion=%v,sCity=%v,"+
	// 		"dCountSh=%v,dCountLo=%v,dRegion=%v,dCity=%v"+
	// 		" bytes=%vi,packets=%vi,version=\"%v\","+
	// 		"ddLat=%f,ddLon=%f,ssLat=%f,ssLon=%f %v",
	// 		m.FlowVersion, m.Device, m.ProtoName,
	// 		m.SrcIP, m.SrcPortName, m.DstIP, m.DstPortName,
	// 		ilSrc.Country_short, ilSrc.Country_long, ilSrc.Region, ilSrc.City,
	// 		ilDst.Country_short, ilDst.Country_long, ilDst.Region, ilDst.City,
	// 		m.Bytes, m.Packets, m.FlowVersion,
	// 		ilDst.Latitude, ilDst.Longitude, ilSrc.Latitude, ilSrc.Longitude, m.Time.UnixNano(),
	// 	)

	// 	// write proto line records
	// 	wapi.WriteRecord(protoLine)
	// }

	// i.Debuuger.Verbose(fmt.Sprintf("writing %v record(s) to %v:%v bucket:%v org:%v ", len(metrics), i.Host, i.Port, i.Bucket, i.Org), logrus.DebugLevel)

	// // write to influx
	// wapi.Flush()

	// device measurement
	go i.measureDevice(metrics)

	return nil
}

// getLocation of ip address
func (i *InfluxDBv2) getLocation(ip string) *ip2location.IP2Locationrecord {
	// get public ip
	il, _ := i.iplocation.GetAll(ip)

	if il.Country_short == "-" {
		// maybe a local IP address
		il, _ = i.iplocation.GetAllPrivate(ip)
	}

	//remove -,_ from strings in order to use them as tag in influxDB
	il.Country_long = i.removeInvalidCharFromTags(il.Country_long)
	il.Country_short = i.removeInvalidCharFromTags(il.Country_short)
	il.City = i.removeInvalidCharFromTags(il.City)
	il.Region = i.removeInvalidCharFromTags(il.Region)
	il.Isp = i.removeInvalidCharFromTags(il.Isp)
	il.Domain = i.removeInvalidCharFromTags(il.Domain)
	il.Netspeed = i.removeInvalidCharFromTags(il.Netspeed)
	il.Iddcode = i.removeInvalidCharFromTags(il.Iddcode)
	il.Areacode = i.removeInvalidCharFromTags(il.Areacode)
	il.Weatherstationcode = i.removeInvalidCharFromTags(il.Weatherstationcode)
	il.Weatherstationname = i.removeInvalidCharFromTags(il.Weatherstationname)
	il.Mcc = i.removeInvalidCharFromTags(il.Mcc)
	il.Mnc = i.removeInvalidCharFromTags(il.Mnc)
	il.Mobilebrand = i.removeInvalidCharFromTags(il.Mobilebrand)
	il.Usagetype = i.removeInvalidCharFromTags(il.Usagetype)

	// return ip2location info
	return il
}

func (i *InfluxDBv2) removeInvalidCharFromTags(s string) string {
	if s == "-" {
		return "NA"
	}
	if strings.Contains(s, "Please upgrade the data file") {
		return "NA"
	}

	rs := strings.Replace(s, ",", " ", -1)
	rs = strings.Replace(rs, "'", " ", -1)
	rs = strings.Replace(rs, " ", "_", -1)

	return rs
}

// close influx db connection
func (i *InfluxDBv2) Close() error {

	// close influxdb client
	i.client.Close()

	i.Debuuger.Verbose(fmt.Sprintf("Closing  %v:%v bucket:%v org:%v", i.Host, i.Port, i.Bucket, i.Org), logrus.DebugLevel)

	// because influx db returns nithing on close :-D
	return nil
}

// fix - & not availale fields in localdb
func (i *InfluxDBv2) fixNotAvailableFileds(il *ip2location.IP2Locationrecord) {
	ip2location.Printrecord(*il)
}
