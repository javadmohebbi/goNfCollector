package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/updater"
)

var Version = "development"

var BuildTime = "build date time"

func main() {

	ips := flag.Bool("ipsum", false, "Downlaod IPSum")

	ipl := flag.Bool("ip2l", false, "Downlaod IP2Location")
	iplAsn := flag.Bool("ip2l-asn", false, "Downlaod IP2Location ASN")
	iplProx := flag.Bool("ip2l-proxy", false, "Downlaod IP2Location Proxy")

	// version
	version := flag.Bool("version", false, "Print version")

	// parse the flags
	flag.Parse()

	if *version {
		fmt.Printf("\n'%v'\n\tVersion: %v\n", filepath.Base(os.Args[0]), Version)
		fmt.Printf("\tBuild Date: %v\n\n", BuildTime)
		os.Exit(0)
	}

	flag.Parse()

	if *ips {
		// download ip sum
		downloadIPSum()
	}

	// download ip2location
	downloadIP2Location(*ipl, *iplAsn, *iplProx)

}

func downloadIP2Location(ipl, iplASN, iplProx bool) {
	// download IP2location lite db
	i2lConf := getIP2LocationCof()

	if ipl {
		// download lite
		updater.DownloadDatabase(
			filepath.Base(i2lConf.IP),    // file name
			filepath.Dir(i2lConf.IP)+"/", // directory
			"/tmp/",                      // tmp file path
			"https://download.openintelligence24.com/vendors/ip2location/IP2LOCATION-LITE-DB11.IPV6.BIN.ZIP", // url to download
			true, // need unzip
		)

	}

	if iplASN {
		// download asn
		updater.DownloadDatabase(
			filepath.Base(i2lConf.ASN),    // file name
			filepath.Dir(i2lConf.ASN)+"/", // directory
			"/tmp/",                       // tmp file path
			"https://download.openintelligence24.com/vendors/ip2location/IP2LOCATION-LITE-ASN.IPV6.CSV.ZIP", // url to download
			true, // need unzip
		)
	}

	if iplProx {
		// download proxy
		updater.DownloadDatabase(
			filepath.Base(i2lConf.Proxy),    // file name
			filepath.Dir(i2lConf.Proxy)+"/", // directory
			"/tmp/",                         // tmp file path
			"https://download.openintelligence24.com/vendors/ip2location/IP2PROXY-LITE-PX10.IPV6.CSV.ZIP", // url to download
			true, // need unzip
		)
	}

}

func downloadIPSum() {
	// download IPSUM db
	colConf := getCollectorCof()
	updater.DownloadDatabase(
		filepath.Base(colConf.IPReputation.IPSumPath),    // file name
		filepath.Dir(colConf.IPReputation.IPSumPath)+"/", // directory
		"/tmp/", // tmp file path
		"https://download.openintelligence24.com/vendors/ipsum/ipsum.txt", // url to download
		false, // does not need unzip
	)
}

func getIP2LocationCof() *configurations.IP2Location {
	// create new instance of configurations interface
	cfg, err := configurations.New(configurations.CONF_TYPE_IP2LOCATION)
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

	c := conf.(*configurations.IP2Location)

	return c

}

func getCollectorCof() *configurations.Collector {
	// create new instance of configurations interface
	cfg, err := configurations.New(configurations.CONF_TYPE_COLLECTOR)
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

	c := conf.(*configurations.Collector)

	return c
}
