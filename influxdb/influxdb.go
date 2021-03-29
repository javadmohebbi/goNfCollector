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

	// device measurement
	go i.measureDevice(metrics)

	// write src host
	go i.measureSrcDstRelatedMetrics(metrics, "src")

	// // write src dst
	go i.measureSrcDstRelatedMetrics(metrics, "dst")

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

	// close location
	i.iplocation.Close()

	i.Debuuger.Verbose(fmt.Sprintf("Closing  %v:%v bucket:%v org:%v", i.Host, i.Port, i.Bucket, i.Org), logrus.DebugLevel)

	// because influx db returns nithing on close :-D
	return nil
}

// fix - & not availale fields in localdb
func (i *InfluxDBv2) fixNotAvailableFileds(il *ip2location.IP2Locationrecord) {
	ip2location.Printrecord(*il)
}
