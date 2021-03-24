package collector

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/goNfCollector/collector/nfipfix"
	"github.com/goNfCollector/collector/nfv1"
	"github.com/goNfCollector/collector/nfv5"
	"github.com/goNfCollector/collector/nfv6"
	"github.com/goNfCollector/collector/nfv7"
	"github.com/goNfCollector/collector/nfv9"
	"github.com/goNfCollector/common"
	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/debugger"
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

// Collector
type Collector struct {

	// listen host
	host string

	// listen port
	port int

	// logrus for futture use, not today :-D
	l *logrus.Logger

	// configuration for futture use, not today :-D
	c *configurations.Collector

	// debugger for verbosing the logs
	d *debugger.Debugger

	// channel
	ch chan bool

	// wait group
	waitGroup *sync.WaitGroup
}

// create new netflow collector
func New(h string, p int, l *logrus.Logger, c *configurations.Collector, d *debugger.Debugger) *Collector {
	return &Collector{
		host: h,
		port: p,
		l:    l,
		c:    c,
		d:    d,

		ch:        make(chan bool),
		waitGroup: &sync.WaitGroup{},
	}
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
	defer conn.Close()

	// done the wait group
	defer nf.waitGroup.Done()

	// set the buffer size
	data := make([]byte, bufferSize)

	// set the decoder
	decoders := make(map[string]*netflow.Decoder)

	// collect & wait until
	// get the SIGTERM
	for {
		select {
		case <-nf.ch:
			nf.d.Verbose("Stopping netflow collector ...",
				logrus.InfoLevel,
			)
			return
		default:
			// nothing to do
		}

		// read data recieved from exporter device
		length, remote, err := conn.ReadFrom(data)
		if err != nil {
			nf.d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
				configurations.ERROR_CAN_T_READ_DATA.Int(),
				configurations.ERROR_CAN_T_READ_DATA, err),
				logrus.DebugLevel,
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
				logrus.DebugLevel,
			)
			continue
		}

		// write debug info
		nf.d.Verbose(fmt.Sprintf("received %d bytes from %s\n", length, remote), logrus.DebugLevel)

		// parse netflow
		go nf.parse(m, remote, data)

	}
}

// parse netflow from traffic
func (nf *Collector) parse(m interface{}, remote net.Addr, data []byte) {
	defer nf.waitGroup.Done()
	nf.waitGroup.Add(1)

	var metrics []common.Metric

	switch p := m.(type) {
	case *netflow1.Packet:
		metrics = nfv1.Prepare(remote.String(), p)

	case *netflow5.Packet:
		metrics = nfv5.Prepare(remote.String(), p)

	case *netflow6.Packet:
		metrics = nfv6.Prepare(remote.String(), p)

	case *netflow7.Packet:
		metrics = nfv7.Prepare(remote.String(), p)

	case *netflow9.Packet:
		metrics = nfv9.Prepare(remote.String(), p)

	case *ipfix.Message:
		metrics = nfipfix.Prepare(remote.String(), p)

	}

	log.Println(metrics)

}
