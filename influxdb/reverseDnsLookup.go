package influxdb

import "net"

// get ip addr & try to get the reverse dns lookup
func (i *InfluxDBv2) revereseDNS(ip string) string {
	addr, err := net.LookupAddr(ip)

	if err != nil {
		return "NA-" + ip
	}

	// return first resolved
	return addr[0]
}
