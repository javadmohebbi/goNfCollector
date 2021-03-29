package location

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/debugger"
	"github.com/ip2location/ip2location-go"
	"github.com/sirupsen/logrus"
)

const (
	ASN_DB_NUMBER_INDEX int = 10 // index of as number in csv db
	ASN_DB_NAME_INDEX   int = 11 // index of as name in csv db

	PRX_IP_FROM    int = 0  // index of IP_FROM in csv db
	PRX_IP_TO      int = 1  // index of IP_TO in csv db
	PRX_TYPE       int = 2  // index of PROXY_TYPE in csv db
	PRX_ISP        int = 7  // index of ISP in csv db
	PRX_DOMAIN     int = 8  // index of DOMAIN in csv db
	PRX_USAGE_TYPE int = 9  // index of USAGE_TYPE in csv db
	PRX_THREAT     int = 13 // index of USAGE_TYPE in csv db
)

// IP 2 Location struct
type IPLocation struct {
	// ASN DB Path
	ASN string `json:"asn"`

	// IP DB Path
	IP string `json:"ip"`

	// Proxy DB Path
	Proxy string `json:"proxy"`

	// ip2location db
	ip2lDB *ip2location.DB

	// read huge proxy db
	allProxies [][]string

	// asn records
	allASN [][]string

	// local records
	locals [][]string

	d *debugger.Debugger
}

// create new IP2location
func New(c *configurations.IP2Location, debugger *debugger.Debugger) *IPLocation {

	// read proxy db
	file, err := os.Open(c.Proxy)
	if err != nil {
		debugger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_CAN_T_OPEN_PROXY_DB.Int(),
			configurations.ERROR_CAN_T_OPEN_PROXY_DB, err),
			logrus.ErrorLevel,
		)
	} else {
		defer file.Close()
	}
	parser := csv.NewReader(file)
	records, err := parser.ReadAll()
	if err != nil {
		debugger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_CAN_T_READ_PROXY_DB.Int(),
			configurations.ERROR_CAN_T_READ_PROXY_DB, err),
			logrus.ErrorLevel,
		)
	}

	// read proxy db
	fileASN, err := os.Open(c.Proxy)
	if err != nil {
		debugger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_CAN_T_OPEN_ASN_DB.Int(),
			configurations.ERROR_CAN_T_OPEN_ASN_DB, err),
			logrus.ErrorLevel,
		)
	} else {
		defer fileASN.Close()
	}
	parserASN := csv.NewReader(fileASN)
	recordsASN, err := parserASN.ReadAll()
	if err != nil {
		debugger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_CAN_T_READ_ASN_DB.Int(),
			configurations.ERROR_CAN_T_READ_ASN_DB, err),
			logrus.ErrorLevel,
		)
	}

	// read csv file
	csvFile, err := os.Open(c.Local)
	if err != nil {
		debugger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_CAN_T_OPEN_LOCAL_DB.Int(),
			configurations.ERROR_CAN_T_OPEN_LOCAL_DB, err),
			logrus.ErrorLevel,
		)
	} else {
		defer csvFile.Close()
	}
	r := csv.NewReader(csvFile)
	recordsLocal, err := r.ReadAll()
	if err != nil {
		debugger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_CAN_T_READ_LOCAL_DB.Int(),
			configurations.ERROR_CAN_T_READ_LOCAL_DB, err),
			logrus.ErrorLevel,
		)
	}

	// open db
	db, err := ip2location.OpenDB(c.IP)
	// check for error
	if err != nil {
		debugger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_OPEN_IP2LOCATION_DB.Int(),
			configurations.ERROR_OPEN_IP2LOCATION_DB, err),
			logrus.ErrorLevel,
		)
		os.Exit(configurations.ERROR_OPEN_IP2LOCATION_DB.Int())
	}

	return &IPLocation{
		ASN:   c.ASN,
		IP:    c.IP,
		Proxy: c.Proxy,

		d:          debugger,
		allProxies: records,
		allASN:     recordsASN,
		locals:     recordsLocal,
		ip2lDB:     db,
	}
}

func (l *IPLocation) Close() {
	// close ip 2 location db
	l.ip2lDB.Close()
}

// proxy info struct, needed info
type ProxyInfo struct {
	Type      string `json:"type"`
	ISP       string `json:"isp"`
	Doamin    string `json:"domain"`
	UsageType string `json:"userType"`
	Threat    string `json:"threat"`
}
