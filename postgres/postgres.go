package postgres

import (
	"fmt"
	"log"
	"strings"

	"github.com/goNfCollector/common"
	"github.com/ip2location/ip2location-go"
	"github.com/sirupsen/logrus"
)

// write to db
func (p *Postgres) Write(metrics []common.Metric) error {

	// // device measurement
	// go p.measureDevice(metrics)

	// go p.measureDetailsRelatedMetrics(metrics)

	// // write src host
	// go p.measureSrcDstRelatedMetrics(metrics, "src")

	// // // write src dst
	// go p.measureSrcDstRelatedMetrics(metrics, "dst")

	log.Println(" ============= WRITE CALLED ============= ")

	return nil
}

// getLocation of ip address
func (p *Postgres) getLocation(ip string) *ip2location.IP2Locationrecord {
	// get public ip
	il, _ := p.iplocation.GetAll(ip)

	if il.Country_short == "-" {
		// maybe a local IP address
		il, _ = p.iplocation.GetAllPrivate(ip)
	}

	//remove -,_ from strings in order to use them as tag in influxDB
	il.Country_long = p.removeInvalidCharFromTags(il.Country_long)
	il.Country_short = p.removeInvalidCharFromTags(il.Country_short)
	il.City = p.removeInvalidCharFromTags(il.City)
	il.Region = p.removeInvalidCharFromTags(il.Region)
	il.Isp = p.removeInvalidCharFromTags(il.Isp)
	il.Domain = p.removeInvalidCharFromTags(il.Domain)
	il.Netspeed = p.removeInvalidCharFromTags(il.Netspeed)
	il.Iddcode = p.removeInvalidCharFromTags(il.Iddcode)
	il.Areacode = p.removeInvalidCharFromTags(il.Areacode)
	il.Weatherstationcode = p.removeInvalidCharFromTags(il.Weatherstationcode)
	il.Weatherstationname = p.removeInvalidCharFromTags(il.Weatherstationname)
	il.Mcc = p.removeInvalidCharFromTags(il.Mcc)
	il.Mnc = p.removeInvalidCharFromTags(il.Mnc)
	il.Mobilebrand = p.removeInvalidCharFromTags(il.Mobilebrand)
	il.Usagetype = p.removeInvalidCharFromTags(il.Usagetype)

	// return ip2location info
	return il
}

func (p *Postgres) removeInvalidCharFromTags(s string) string {
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

// close postgres db connection
func (p *Postgres) Close() error {

	// closing ....
	p.Debuuger.Verbose(fmt.Sprintf("Closing  %v:%v db:%v", p.Host, p.Port, p.DB), logrus.InfoLevel)

	p.Debuuger.Verbose(fmt.Sprintf("Please wait until pending writes finish on %v:%v db:%v", p.Host, p.Port, p.DB), logrus.InfoLevel)

	defer p.Debuuger.Verbose(fmt.Sprintf("%v:%v db:%v closed!", p.Host, p.Port, p.DB), logrus.InfoLevel)

	// wait until all of tasks are done
	p.WaitGroup.Wait()

	// close influxdb client
	p.pool.Close()

	// close location
	p.iplocation.Close()

	// because influx db returns nithing on close :-D
	return nil
}
