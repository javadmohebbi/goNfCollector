package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/updater"
)

func main() {

	// download ip sum
	downloadIPSum()

	// download ip2location
	downloadIP2Location()

}

func downloadIP2Location() {
	// download IP2location lite db
	i2lConf := getIP2LocationCof()

	// download asn
	updater.DownloadDatabase(
		filepath.Base(i2lConf.ASN),    // file name
		filepath.Dir(i2lConf.ASN)+"/", // directory
		"/tmp/",                       // tmp file path
		"https://github.com/javadmohebbi/goNfCollector/raw/main/updates/IP2LOCATION-LITE-ASN.IPV6.CSV.ZIP", // url to download
		true, // need unzip
	)

	// download lite
	updater.DownloadDatabase(
		filepath.Base(i2lConf.IP),    // file name
		filepath.Dir(i2lConf.IP)+"/", // directory
		"/tmp/",                      // tmp file path
		"https://github.com/javadmohebbi/goNfCollector/raw/main/updates/IP2LOCATION-LITE-DB11.IPV6.BIN.ZIP", // url to download
		true, // need unzip
	)

	// download proxy
	updater.DownloadDatabase(
		filepath.Base(i2lConf.Proxy),    // file name
		filepath.Dir(i2lConf.Proxy)+"/", // directory
		"/tmp/",                         // tmp file path
		"https://github.com/javadmohebbi/goNfCollector/raw/main/updates/IP2PROXY-LITE-PX10.IPV6.CSV.ZIP", // url to download
		true, // need unzip
	)

}

func downloadIPSum() {
	// download IPSUM db
	colConf := getCollectorCof()
	updater.DownloadDatabase(
		filepath.Base(colConf.IPReputation.IPSumPath),    // file name
		filepath.Dir(colConf.IPReputation.IPSumPath)+"/", // directory
		"/tmp/", // tmp file path
		"https://raw.githubusercontent.com/stamparm/ipsum/master/ipsum.txt", // url to download
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
