package influxdb

import (
	"fmt"

	"github.com/goNfCollector/common"
	"github.com/gookit/color"
	"github.com/ip2location/ip2location-go"
	"github.com/sirupsen/logrus"
)

func (i *InfluxDBv2) Write(metrics []common.Metric) error {

	// // define new write api
	// // wapi := i.client.WriteAPI(i.Org, i.Bucket)

	// // wapi.WriteRecord(fmt.Sprintf("test,unit=temp avg=%f,max=%f", 33.4, 44.2))
	// // wapi.WriteRecord(fmt.Sprintf("test,unit=temp avg=%f,max=%f", 23.4, 22.2))
	// // wapi.WriteRecord(fmt.Sprintf("test,unit=temp avg=%f,max=%f", 43.4, 56.2))
	// // wapi.WriteRecord(fmt.Sprintf("test,unit=temp avg=%f,max=%f", 13.4, 14.2))
	// // wapi.WriteRecord(fmt.Sprintf("test,unit=temp avg=%f,max=%f", 32.4, 42.2))
	// // wapi.WriteRecord(fmt.Sprintf("test,unit=temp avg=%f,max=%f", 53.4, 54.2))
	// // wapi.WriteRecord(fmt.Sprintf("test,unit=temp avg=%f,max=%f", 34.4, 45.2))

	// // i.Debuuger.Verbose(fmt.Sprintf("writing %v record(s) to %v:%v bucket:%v org:%v ", len(metrics), i.Host, i.Port, i.Bucket, i.Org), logrus.DebugLevel)

	// // wapi.Flush()

	// // no error will be returned yet
	// // but wapi.Errors() channel will be used
	// // in the future
	// return nil

	var il *ip2location.IP2Locationrecord
	for _, m := range metrics {
		il, _ = i.iplocation.GetAll(m.DstIP)

		if il.Country_short == "-" {
			// maybe a local IP address
			il, _ = i.iplocation.GetAllPrivate(m.DstIP)
		}

		color.Yellow.Printf("\nDstIP: %v from %v\n", m.DstIP, il.Country_short)
	}

	return nil

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
