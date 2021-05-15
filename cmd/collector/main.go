package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"

	"github.com/goNfCollector/collector"
	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/debugger"
)

var Version = "development"

var BuildTime = "build date time"

func main() {

	// runtime.GOMAXPROCS(runtime.NumCPU())

	// conf file path
	confFilePath := flag.String("confPath", "/opt/oi24/netflow-collector/etc/", "Path to conf directory. (trailing slash is needed!)")

	// listen address
	addr := flag.String("address", "0.0.0.0", "Collector listen address")

	// listen UDP port
	port := flag.Int("port", 6859, "Collector listen UDP port")

	// debug
	debug := flag.Bool("debug", false, "Enable/Disable debug mode")

	// version
	version := flag.Bool("version", false, "Print version")

	// parse the flags
	flag.Parse()

	if *version {
		fmt.Printf("\n'%v'\n\tVersion: %v\n", filepath.Base(os.Args[0]), Version)
		fmt.Printf("\tBuild Date: %v\n\n", BuildTime)
		os.Exit(0)
	}

	// create new instance of configurations interface
	cfg, err := configurations.New(configurations.CONF_TYPE_COLLECTOR, *confFilePath)
	if err != nil {
		log.Println("Can not create new instance of configuration due to error ", err)
		os.Exit(configurations.ERROR_READ_CONFIG.Int())
	}

	// Read config & return the requested strucut type
	conf, err := cfg.Read()
	if err != nil {
		log.Println("Can not read config file due to error ", err)
		os.Exit(configurations.ERROR_READ_CONFIG.Int())
	}

	// cast cfg to Collector configuration
	collectorConf := conf.(*configurations.Collector)

	// check fo debug in command line argument
	if *debug {
		collectorConf.Debug = *debug
	}

	// check listen address
	if *addr != "" {
		collectorConf.Listen.Address = *addr
	}

	// check listen UDP port
	if *port >= 1 && *port <= 65535 {
		collectorConf.Listen.Port = *port

	}

	// set the number of cpu this app can use
	if collectorConf.CPUNum > 0 && collectorConf.CPUNum <= runtime.NumCPU() {
		color.Yellow.Printf("set max CPU num to: %v\n", collectorConf.CPUNum)
		runtime.GOMAXPROCS(collectorConf.CPUNum)
	} else {
		color.Yellow.Printf("set max CPU num to all available CPUs: %v\n", runtime.NumCPU())
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	// create & configure logrus
	logr := logrus.New()

	// variable for multiwriter
	mw := io.MultiWriter(os.Stdout)

	// set log to file and also logfile
	logr.SetOutput(mw)

	// opening log file for write and append and pass it to io.multiwriter function
	lfn := collectorConf.LogFile
	if *debug {
		color.Blue.Printf("opening log file %v for extra logging\n", lfn)
	}

	// open log file
	lf, err := os.OpenFile(lfn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// close file at the end of statement!
	defer lf.Close()

	// check for log file error
	if err != nil {
		color.Red.Printf("Can not open log file: %v. Logs will be displayed ONLY on standard output (stdout)\n", lfn)
		color.Red.Printf("\t%v\n", err)
	} else {
		// set log file + stdout as log writer
		color.Blue.Printf("set log file %v for extra logging\n", lfn)
		mw = io.MultiWriter(lf, os.Stdout)
	}

	// set log to file and also logfile
	logr.SetOutput(mw)

	// print to console if debuging is enabled
	if collectorConf.Debug {
		color.Yellow.Printf("--- DEBUGGING IS ENABLED ---\n")
	}

	// Create new debug
	d := debugger.New(collectorConf.Debug, logr, "log")

	// create new instance of nfcollector
	nfcol := collector.New(collectorConf.Listen.Address,
		collectorConf.Listen.Port,
		logr, collectorConf, d,
		*confFilePath,
	)

	// serve netflow collector service
	nfcol.Serve()

}
