package fwsock

import (
	"fmt"
	"net"
	"os"

	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/debugger"
	"github.com/sirupsen/logrus"
)

type FwSockClient struct {

	// debugger for verbosing the logs
	d *debugger.Debugger

	// logruse for future use, not today :-D
	l *logrus.Logger

	// Path to unix socket
	SocketAddr string

	// channel
	Ch chan os.Signal

	// OutChannel
	OutCh chan ClientServerReqResp

	// socket listener
	Conn net.Conn
}

func NewClient(dd *debugger.Debugger, ll *logrus.Logger, path string) *FwSockClient {
	// socekt conf
	// create new instance of configurations interface
	cfg, err := configurations.New(configurations.CONF_TYPE_SOCKET, path)
	if err != nil {
		dd.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_READ_CONFIG.Int(),
			configurations.ERROR_READ_CONFIG, err),
			logrus.ErrorLevel,
		)
		os.Exit(configurations.ERROR_READ_CONFIG.Int())
	}

	// Read config & return the requested strucut type
	fwsConfig, err := cfg.Read()
	if err != nil {
		dd.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_READ_CONFIG.Int(),
			configurations.ERROR_READ_CONFIG, err),
			logrus.ErrorLevel,
		)
		os.Exit(configurations.ERROR_READ_CONFIG.Int())
	}

	fwsCfg := fwsConfig.(*configurations.Socket)

	fwsClient := &FwSockClient{
		d:          dd,
		l:          ll,
		Ch:         make(chan os.Signal, 1),
		OutCh:      make(chan ClientServerReqResp),
		SocketAddr: fwsCfg.Socket,
	}

	c, err := net.Dial("unix", fwsClient.SocketAddr)
	if err != nil {
		dd.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_CAN_T_DIAL_LINUX_SOCKET_ADDRESS.Int(),
			configurations.ERROR_CAN_T_DIAL_LINUX_SOCKET_ADDRESS, err),
			logrus.ErrorLevel,
		)
		os.Exit(configurations.ERROR_CAN_T_DIAL_LINUX_SOCKET_ADDRESS.Int())
	}

	fwsClient.Conn = c

	return fwsClient
}
