package influxdb

import (
	"net"

	"github.com/ammario/ipisp"
)

// get asn lookup
func (i *InfluxDBv2) asnLookup(ip string) string {
	client, err := ipisp.NewDNSClient()

	// if client has error, returns NA
	if err != nil {
		return "NA"
	}

	// close ipisp client
	defer client.Close()

	resp, err := client.LookupIP(net.ParseIP(ip))
	// if no lookup, returns NA
	if err != nil {
		return "NA"
	}

	asnNumber := resp.ASN.String()

	return i.removeInvalidCharFromTags(i.iplocation.GetAsnName(asnNumber) + "_" + asnNumber)

	// // return first resolved
	// return resp.ASN.String()
}
