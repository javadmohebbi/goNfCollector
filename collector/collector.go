package collector

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/goNfCollector/collector/nfipfix"
	"github.com/goNfCollector/collector/nfv1"
	"github.com/goNfCollector/collector/nfv5"
	"github.com/goNfCollector/collector/nfv6"
	"github.com/goNfCollector/collector/nfv7"
	"github.com/goNfCollector/collector/nfv9"
	"github.com/goNfCollector/common"
	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/database"
	"github.com/goNfCollector/debugger"
	"github.com/goNfCollector/exporters"
	"github.com/goNfCollector/influxdb"
	"github.com/goNfCollector/location"
	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
	"github.com/tehmaze/netflow"
	"github.com/tehmaze/netflow/ipfix"
	"github.com/tehmaze/netflow/netflow1"
	"github.com/tehmaze/netflow/netflow5"
	"github.com/tehmaze/netflow/netflow6"
	"github.com/tehmaze/netflow/netflow7"
	"github.com/tehmaze/netflow/netflow9"
	"github.com/tehmaze/netflow/session"
)

const bufferSize int = 8960

const maxQueueSize int = 20480

// Collector
type Collector struct {

	// listen host
	host string

	// listen port
	port int

	// logrus for futture use, not today :-D
	l *logrus.Logger

	// configuration for collector
	c *configurations.Collector

	// translation configuration
	cfTrans *configurations.Translations

	// configuration for ip2location
	iploc *location.IPLocation

	// debugger for verbosing the logs
	d *debugger.Debugger

	// channel
	ch chan os.Signal

	// wait group
	waitGroup *sync.WaitGroup

	exporters []exporters.Exporter

	// bytes
	outgoingMessage outgoingMessage

	isConClosed bool

	// total number of recieved flows
	numberOfRecievedFlows uint64

	// total number of flows sent for export
	numberOfFlowsSentForExport uint64

	// ztdb *zabbix.ZabbixTimeScaleDB

	// portmap for port and protocol
	portmap    common.PortMap
	portmapErr error
}

type outgoingMessage struct {
	recipient *net.UDPAddr
	data      []byte
}

// create new netflow collector
func New(h string, p int, l *logrus.Logger, c *configurations.Collector, d *debugger.Debugger, path string, cfTrans *configurations.Translations) *Collector {

	// getIP2location conf
	// create new instance of configurations interface
	cfg, err := configurations.New(configurations.CONF_TYPE_IP2LOCATION, path)
	if err != nil {
		d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_READ_CONFIG.Int(),
			configurations.ERROR_READ_CONFIG, err),
			logrus.ErrorLevel,
		)
		os.Exit(configurations.ERROR_READ_CONFIG.Int())
	}

	// Read config & return the requested strucut type
	cf, err := cfg.Read()
	if err != nil {
		d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_READ_CONFIG.Int(),
			configurations.ERROR_READ_CONFIG, err),
			logrus.ErrorLevel,
		)
		os.Exit(configurations.ERROR_READ_CONFIG.Int())
	}

	// make new instance of ip2location
	i2l := location.New(cf.(*configurations.IP2Location), d)

	nf := &Collector{
		host: h,
		port: p,
		l:    l,
		c:    c,
		d:    d,

		ch:        make(chan os.Signal, 1),
		waitGroup: &sync.WaitGroup{},

		iploc: i2l,

		cfTrans: cfTrans,
	}

	// portMap definition
	nf.portmap, nf.portmapErr = common.GetServices()

	// extract valid exporters
	nf.exporters = nf.getExporters()

	// grab the signals
	signal.Notify(nf.ch, syscall.SIGINT, syscall.SIGTERM)

	return nf
}

// listen to the provided configuration
func (nf *Collector) listen() (*net.UDPConn, configurations.ErrorCodes, error) {
	// Resolve Address
	sAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%v", nf.host, nf.port))
	if err != nil {
		nf.d.Verbose(fmt.Sprintf("[%d]-%s: (%v)", configurations.ERROR_LISTEN_RESOLVE_UDP_ADDRESS.Int(), configurations.ERROR_LISTEN_RESOLVE_UDP_ADDRESS, err), logrus.ErrorLevel)
		return nil, configurations.ERROR_LISTEN_RESOLVE_UDP_ADDRESS, err
	}

	// liste on provided host:port
	nf.d.Verbose(fmt.Sprintf("listening on %s:%d", nf.host, nf.port), logrus.DebugLevel)
	con, err := net.ListenUDP("udp", sAddr)
	if err != nil {
		nf.d.Verbose(fmt.Sprintf("[%d]-%s: (%v)", configurations.ERROR_LISTEN_ON_UDP.Int(), configurations.ERROR_LISTEN_ON_UDP, err), logrus.ErrorLevel)
		return nil, configurations.ERROR_LISTEN_ON_UDP, err
	}

	// set up connection read buffer
	if err = con.SetReadBuffer(bufferSize); err != nil {
		nf.d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_CAN_T_SET_CONNECTION_READ_BUFFER.Int(),
			configurations.ERROR_CAN_T_SET_CONNECTION_READ_BUFFER, err),
			logrus.ErrorLevel,
		)
		return nil, configurations.ERROR_CAN_T_SET_CONNECTION_READ_BUFFER, err
	}

	// server is listening
	nf.d.Verbose(fmt.Sprintf("Server is now listening on %s:%d (UDP)...!", nf.host, nf.port), logrus.InfoLevel)

	return con, configurations.NO_ERROR, nil
}

// serve the netflow collector service
func (nf *Collector) Serve() {

	udpConn, ec, err := nf.listen()
	if err != nil {
		// error with the custom error codes
		os.Exit(ec.Int())
	}

	// start collecting netflows
	nf.collect(udpConn)

}

// this method will do the collection
func (nf *Collector) collect(conn *net.UDPConn) {
	// close the udp connection
	// defer conn.Close()
	// defer nf.ztdb.Close()

	// set the buffer size
	data := make([]byte, bufferSize)

	// set the decoder
	decoders := make(map[string]*netflow.Decoder)

	// make notify channel
	nf.ch = make(chan os.Signal, 1)
	signal.Notify(nf.ch,
		// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGKILL, // "always fatal", "SIGKILL and SIGSTOP may not be caught by a program"
		syscall.SIGHUP,  // "terminal is disconnected"
	)

	go func() {
		// check if channel signal has notified
		<-nf.ch
		// defer close(nf.ch)

		nf.d.Verbose("Stopping netflow collector ...",
			logrus.InfoLevel,
		)

		nf.isConClosed = true

		// cleaning up things

		// close all open
		for _, e := range nf.exporters {
			// close exporter clients if needed
			e.Close()
		}

		defer nf.d.Verbose(fmt.Sprintf("Please wait until netflow collector finishes pending %v jobs", runtime.NumGoroutine()), logrus.InfoLevel)
		nf.waitGroup.Wait()

		color.Red.Printf("\nApp Exited due to recieved signal from OS or User!\n")

		os.Exit(0)
	}()

	// collect & wait until
	// get the SIGTERM
	for {

		if !nf.isConClosed {
			// read data recieved from exporter device
			length, remote, err := conn.ReadFrom(data)
			if err != nil {
				nf.d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
					configurations.ERROR_CAN_T_READ_DATA.Int(),
					configurations.ERROR_CAN_T_READ_DATA, err),
					logrus.ErrorLevel,
				)

				continue
			}

			// find the decoders
			// or if not, make new
			d, found := decoders[remote.String()]
			if !found {
				s := session.New()
				d = netflow.NewDecoder(s)
				decoders[remote.String()] = d
			}

			// use netflow decoder to decode the recieved netflow,
			// if possible!
			m, err := d.Read(bytes.NewBuffer(data[:length]))
			if err != nil {
				nf.d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
					configurations.ERROR_CAN_T_DECODE_NETFLOW_DATA.Int(),
					configurations.ERROR_CAN_T_DECODE_NETFLOW_DATA, err),
					logrus.ErrorLevel,
				)
				continue
			}

			// write debug info
			// nf.d.Verbose(fmt.Sprintf("received %d bytes from %s\n", length, remote), logrus.DebugLevel)

			// parse netflow
			nf.waitGroup.Add(1)
			go nf.parse(m, remote, data)
		}

	}

}

// parse netflow from traffic
func (nf *Collector) parse(m interface{}, remote net.Addr, data []byte) {
	defer nf.waitGroup.Done()

	// metrics to collect
	var metrics []common.Metric

	mustAccept := false

	if nf.c.AcceptFrom != "any" && nf.c.AcceptFrom != "" {
		spl := strings.Split(nf.c.AcceptFrom, ",")

		for _, s := range spl {
			hst, _, _ := net.SplitHostPort(remote.String())
			if s == hst {
				mustAccept = true
				continue
			}
		}

	} else {
		mustAccept = true
	}

	// not listed
	if !mustAccept {
		nf.d.Verbose(fmt.Sprintf("'%v' device is not defined to be accepted", remote.String()), logrus.DebugLevel)
		return
	}

	// check the netflow version
	switch p := m.(type) {
	case *netflow1.Packet:
		metrics = nfv1.Prepare(remote.String(), p, nf.portmap, nf.portmapErr)

	case *netflow5.Packet:
		metrics = nfv5.Prepare(remote.String(), p, nf.portmap, nf.portmapErr)

	case *netflow6.Packet:
		metrics = nfv6.Prepare(remote.String(), p, nf.portmap, nf.portmapErr)

	case *netflow7.Packet:
		metrics = nfv7.Prepare(remote.String(), p, nf.portmap, nf.portmapErr)

	case *netflow9.Packet:
		metrics = nfv9.Prepare(remote.String(), p, nf.portmap, nf.portmapErr, nf.cfTrans)

	case *ipfix.Message:
		metrics = nfipfix.Prepare(remote.String(), p, nf.portmap, nf.portmapErr, nf.cfTrans)

	}

	// export metrics if neededs
	if len(metrics) > 0 {

		if nf.c.Debug {
			nf.d.Verbose(fmt.Sprintf("'%v' flows recieved from '%v'", len(metrics), remote.String()), logrus.DebugLevel)
		}

		nf.numberOfRecievedFlows += uint64(len(metrics))

		nf.export(metrics)
	}

}

// find valid netflow exporters and return them
func (nf *Collector) getExporters() []exporters.Exporter {

	// array to return at the end
	var exps []exporters.Exporter

	// Loop through InfluxDB Exporters
	for _, ex := range nf.c.Exporter.InfluxDBs {

		// create new influxDB
		ifl := influxdb.New(ex.Token, ex.Bucket, ex.Org, ex.Host, nf.c.IPReputation, ex.Port, nf.d, nf.iploc)

		// create new influxDB exporter
		influxExporter, err := exporters.New(ifl, ifl.Debuuger)
		if err != nil {
			// errors handled in the exporter new package
			continue
		}

		// if no error, append it to exporters
		exps = append(exps, *influxExporter)
	}

	// Loop through Postgres Exporters
	for _, ex := range nf.c.Exporter.Postgres {

		// create new Postgres
		ifl := database.New(ex.Host, ex.User, ex.Password, ex.DB, nf.c.IPReputation, ex.Port, nf.d, nf.iploc, ex.MaxIdleConnection, ex.MaxOpenConnection, 1*time.Hour)

		// create new Postgres exporter
		postgresExporter, err := exporters.New(ifl, ifl.Debuuger)
		if err != nil {
			// errors handled in the exporter new package
			continue
		}

		// if no error, append it to exporters
		exps = append(exps, *postgresExporter)
	}

	return exps
}

// export if needed
func (nf *Collector) export(metrics []common.Metric) {

	nf.waitGroup.Add(1)

	// check if there are valid exporters
	if len(nf.exporters) > 0 {
		// loop through exporters
		for _, e := range nf.exporters {
			// write metrics
			// go e.Write(metrics)
			go func(e exporters.Exporter, metrics []common.Metric) {
				nf.numberOfFlowsSentForExport += uint64(len(metrics))
				e.Write(metrics)
				defer nf.waitGroup.Done()
			}(e, metrics)
		}
	}

	// go nf.ztdb.Store(metrics)

}
