package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/goNfCollector/debugger"
	"github.com/goNfCollector/fwsock"
	"github.com/sirupsen/logrus"
)

func main() {

	// conf file path
	confFilePath := flag.String("confPath", "/opt/oi24/netflow-collector/etc/", "Path to conf directory. (trailing slash is needed!)")
	// parse the flags
	flag.Parse()

	// create & configure logrus
	logr := logrus.New()

	// Create new debug
	d := debugger.New(true, logr, "log")

	fws := fwsock.NewClient(d, logr, *confFilePath)

	fws.SetChann(make(chan os.Signal, 1))
	signal.Notify(fws.Ch, syscall.SIGINT, syscall.SIGTERM)

	// handle signals
	go func() {
		<-fws.Ch

		//
		log.Println("CTRL + C recvd")

		close(fws.Ch)

		// close socket listener
		fws.Close()

		os.Exit(0)
	}()

	go fws.Reader(fws.Conn)

	req := fwsock.ClientServerReqResp{
		API:     true,
		Command: fwsock.CMD_INIT,
	}
	bts, err := req.JSONToStringClientServerReqResp()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(bts)

	_, err = fws.Conn.Write([]byte(fmt.Sprintf("%s\n", bts)))
	if err != nil {
		log.Fatalln("write error:", err)
	}

	for {
	}

}
