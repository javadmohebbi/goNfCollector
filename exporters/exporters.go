package exporters

import (
	"errors"
	"fmt"

	"github.com/goNfCollector/common"
	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/debugger"
	"github.com/goNfCollector/influxdb"
	"github.com/sirupsen/logrus"
)

type Exporter interface {

	// write to exporter
	Write([]common.Metric) error

	// close exporter
	Close() error
}

// create new Exporter
// in case of error will return nil
func New(exporter interface{}, d *debugger.Debugger) (*Exporter, error) {

	// find the correct type of exporter
	switch exporter.(type) {

	// if it's influxdb exporter type
	case influxdb.InfluxDBv2:

		// type assertion from interface to the struct
		exp := exporter.(influxdb.InfluxDBv2)

		// create new exporter for influxdb
		e := Exporter(&exp)

		// return it
		return &e, nil

	default:
		err := errors.New(configurations.ERROR_NO_VALID_EXPORTER_FOUND.String())
		// can not find any valid exporter type
		d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_NO_VALID_EXPORTER_FOUND.Int(),
			configurations.ERROR_NO_VALID_EXPORTER_FOUND, err),
			logrus.ErrorLevel,
		)
		return nil, err
	}

}
