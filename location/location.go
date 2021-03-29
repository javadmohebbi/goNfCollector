package location

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"strconv"

	"github.com/goNfCollector/configurations"
	"github.com/ip2location/ip2location-go"
	"github.com/sirupsen/logrus"
)

// Get IP location & Other info
func (i *IPLocation) GetAll(addr string) (*ip2location.IP2Locationrecord, error) {

	if i.IP == "" {
		i.d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_OPEN_IP2LOCATION_DB.Int(),
			configurations.ERROR_OPEN_IP2LOCATION_DB, "empty path for the DB"),
			logrus.ErrorLevel,
		)
		return nil, errors.New(fmt.Sprintf("%v", configurations.ERROR_OPEN_IP2LOCATION_DB))
	}

	// open db
	db, err := ip2location.OpenDB(i.IP)

	// check for error
	if err != nil {
		i.d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_OPEN_IP2LOCATION_DB.Int(),
			configurations.ERROR_OPEN_IP2LOCATION_DB, err),
			logrus.ErrorLevel,
		)
		return nil, err
	}

	// get information
	lr, err := db.Get_all(addr)

	// check for error
	if err != nil {
		i.d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_GET_IP2LOCATION_INFO.Int(),
			configurations.ERROR_GET_IP2LOCATION_INFO, err),
			logrus.ErrorLevel,
		)
		return nil, err
	}

	// return record
	return &lr, nil
}

// GetPrivateIPAddressInfo - If Private CSV files is provided, It will use this file show information
func (i *IPLocation) GetAllPrivate(ip string) (*ip2location.IP2Locationrecord, bool) {

	CountryShort, CountryLong, Region, City, Timezone, Lati, Loni, StartIP, EndIP := 0, 1, 2, 3, 4, 5, 6, 7, 8

	ip2lRec := &ip2location.IP2Locationrecord{
		Country_short:      "-",
		Country_long:       "-",
		Region:             "-",
		City:               "-",
		Isp:                "-",
		Latitude:           0,
		Longitude:          0,
		Domain:             "-",
		Zipcode:            "-",
		Timezone:           "-",
		Netspeed:           "-",
		Iddcode:            "-",
		Areacode:           "-",
		Weatherstationcode: "-",
		Weatherstationname: "-",
		Mcc:                "-",
		Mnc:                "-",
		Mobilebrand:        "-",
		Elevation:          0,
		Usagetype:          "-",
	}

	csvFile, err := os.Open(i.Local)
	if err != nil {
		fmt.Println("Couldn't open the csv file", err)
		return ip2lRec, false
		// os.Exit(1)
	}
	r := csv.NewReader(csvFile)
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			return ip2lRec, false
		}
		if err != nil {
			// fmt.Printf("Can not read %s file", i.Local)
			return ip2lRec, false
		}

		// check if the provided IP is in the range
		if i.isItInTheRangeIPv4(ip, record[StartIP], record[EndIP]) {
			tmpLat, _ := strconv.ParseFloat(record[Lati], 32)
			tmpLon, _ := strconv.ParseFloat(record[Loni], 32)

			ip2lRec.Country_short = record[CountryShort]
			ip2lRec.Country_long = record[CountryLong]
			ip2lRec.City = record[City]
			ip2lRec.Region = record[Region]
			ip2lRec.Timezone = record[Timezone]
			ip2lRec.Latitude = float32(tmpLat)
			ip2lRec.Longitude = float32(tmpLon)

			return ip2lRec, true
		}
	}
}

// IsItInTheRangeIPv4 return true if its in the range
func (i *IPLocation) isItInTheRangeIPv4(ip string, startIP string, endIP string) bool {
	trial := net.ParseIP(ip)

	// ip is NOT an IPv4
	if trial.To4() == nil {
		// fmt.Println("NOT v4")
		return false
	}

	// ip is in the range
	if bytes.Compare(trial, net.ParseIP(startIP)) >= 0 && bytes.Compare(trial, net.ParseIP(endIP)) <= 0 {
		return true
	}

	// ip is not in the range
	return false
}

// GetAsnName with AS Number
func (i *IPLocation) GetAsnName(asNumber string) string {

	// returning asn name
	asn := "NA"

	// find asn name
	for _, record := range i.allASN {

		if "AS"+record[ASN_DB_NUMBER_INDEX] == asNumber {
			asn = record[ASN_DB_NAME_INDEX]
			break
		}

	}

	return asn
}

// GetProxy Info
func (i *IPLocation) GetProxyInfo(adr string) (bool, ProxyInfo) {

	// returning prx name
	hasPrx := false
	prx := ProxyInfo{}

	host := net.ParseIP(adr)
	nh := big.NewInt(0)
	nh.SetBytes(host)

	// find asn name
	for _, record := range i.allProxies {

		// GET FROM RANGE
		nf := new(big.Int)
		nf, ok := nf.SetString(record[PRX_IP_FROM], 10)
		if !ok {
			continue
		}

		// GET FROM RANGE
		nt := new(big.Int)
		nt, ok = nf.SetString(record[PRX_IP_TO], 10)
		if !ok {
			continue
		}

		inRange := false

		compareFrom := nh.Cmp(nf)

		// means nh >= nf
		if compareFrom >= 0 {
			compareTo := nh.Cmp(nt)

			// means nh <= nt
			if compareTo <= 0 {
				inRange = true
			}
		}

		if inRange {

			prx.ISP = record[PRX_ISP]
			prx.Doamin = record[PRX_DOMAIN]
			prx.Threat = record[PRX_THREAT]
			prx.Type = record[PRX_TYPE]
			prx.UsageType = record[PRX_USAGE_TYPE]

			hasPrx = true
			break

		}

	}

	return hasPrx, prx

}
