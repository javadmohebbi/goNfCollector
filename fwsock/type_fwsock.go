package fwsock

import (
	"fmt"
	"net"
	"os"

	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/debugger"
	"github.com/sirupsen/logrus"
)

// forwarding all netflow collector to
// new Unix Sockets
type FwSock struct {
	// debugger for verbosing the logs
	d *debugger.Debugger

	// logruse for future use, not today :-D
	l *logrus.Logger

	// Path to unix socket
	SocketAddr string

	// channel
	Ch chan os.Signal

	// socket listener
	lstn net.Listener

	// socekt client list
	socketClientList *socketClientList
}

// create new instance of FwSocket
func New(dd *debugger.Debugger, ll *logrus.Logger, path string) *FwSock {
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

	fws := &FwSock{
		d:          dd,
		l:          ll,
		Ch:         make(chan os.Signal, 1),
		SocketAddr: fwsCfg.Socket,
	}

	return fws
}
