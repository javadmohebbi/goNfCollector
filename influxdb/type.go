package influxdb

import (
	"fmt"

	"github.com/goNfCollector/debugger"
	"github.com/goNfCollector/location"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/sirupsen/logrus"
)

// influx db v2 struct to
// get access to influxdb version 2.x
type InfluxDBv2 struct {
	// influx db host
	Host string `json:"host"`

	// influx db port
	Port int `json:"port"`

	// auth token
	Token string `json:"token"`

	// bucket name
	Bucket string `json:"bucket"`

	// organization
	Org string `json:"org"`

	// debugger
	Debuuger *debugger.Debugger

	// influxdb client
	client influxdb2.Client

	// IP2locaion instance
	iplocation *location.IPLocation
}

// return exporter info
func (i InfluxDBv2) String() string {
	return fmt.Sprintf("%s:%d bucket:%s org:%s", i.Host, i.Port, i.Bucket, i.Org)
}

// create new instance of influxDB
func New(token, bucket, org, host string, port int, d *debugger.Debugger, ip2location *location.IPLocation) InfluxDBv2 {

	// create new influx db client
	client := influxdb2.NewClient(
		fmt.Sprintf("http://%v:%v", host, port),
		token,
	)

	d.Verbose(fmt.Sprintf("new influxDB exporter %v:%v bucket:%v org:%v is created", host, port, bucket, org), logrus.DebugLevel)

	// retun influxDB
	return InfluxDBv2{
		Token:    token,
		Bucket:   bucket,
		Org:      org,
		Host:     host,
		Port:     port,
		Debuuger: d,
		client:   client,

		iplocation: ip2location,
	}

}
