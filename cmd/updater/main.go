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

	// conf file path
	confFilePath := flag.String("confPath", "/opt/nfcollector/etc/", "Path to conf directory. (trailing slash is needed!)")

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
		downloadIPSum(*confFilePath)
	}

	// download ip2location
	downloadIP2Location(*ipl, *iplAsn, *iplProx, *confFilePath)

}

func downloadIP2Location(ipl, iplASN, iplProx bool, path string) {
	// download IP2location lite db
	i2lConf := getIP2LocationCof(path)

	if ipl {
		// download lite
		updater.DownloadDatabase(
			filepath.Base(i2lConf.IP),    // file name
			filepath.Dir(i2lConf.IP)+"/", // directory
			"/tmp/",                      // tmp file path
			"https://github.com/javadmohebbi/goNfCollector/raw/main/updates/IP2LOCATION-LITE-DB11.IPV6.BIN.ZIP", // url to download
			true, // need unzip
		)

	}

	if iplASN {
		// download asn
		updater.DownloadDatabase(
			filepath.Base(i2lConf.ASN),    // file name
			filepath.Dir(i2lConf.ASN)+"/", // directory
			"/tmp/",                       // tmp file path
			"https://github.com/javadmohebbi/goNfCollector/raw/main/updates/IP2LOCATION-LITE-ASN.IPV6.CSV.ZIP", // url to download
			true, // need unzip
		)
	}

	if iplProx {
		// download proxy
		updater.DownloadDatabase(
			filepath.Base(i2lConf.Proxy),    // file name
			filepath.Dir(i2lConf.Proxy)+"/", // directory
			"/tmp/",                         // tmp file path
			"https://github.com/javadmohebbi/goNfCollector/raw/main/updates/IP2PROXY-LITE-PX10.IPV6.CSV.ZIP", // url to download
			true, // need unzip
		)
	}

}

func downloadIPSum(path string) {
	// download IPSUM db
	colConf := getCollectorCof(path)
	updater.DownloadDatabase(
		filepath.Base(colConf.IPReputation.IPSumPath),    // file name
		filepath.Dir(colConf.IPReputation.IPSumPath)+"/", // directory
		"/tmp/", // tmp file path
		"https://raw.githubusercontent.com/stamparm/ipsum/master/ipsum.txt", // url to download
		false, // does not need unzip
	)
}

func getIP2LocationCof(path string) *configurations.IP2Location {
	// create new instance of configurations interface
	cfg, err := configurations.New(configurations.CONF_TYPE_IP2LOCATION, path)
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

func getCollectorCof(path string) *configurations.Collector {
	// create new instance of configurations interface
	cfg, err := configurations.New(configurations.CONF_TYPE_COLLECTOR, path)
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
