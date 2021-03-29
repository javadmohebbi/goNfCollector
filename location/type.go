package location

import (
	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/debugger"
)

const (
	ASN_DB_NUMBER_INDEX int = 3 // index of as number in csv db
	ASN_DB_NAME_INDEX   int = 4 // index of as name in csv db
)

// IP 2 Location struct
type IPLocation struct {
	// ASN DB Path
	ASN string `json:"asn"`

	// IP DB Path
	IP string `json:"ip"`

	// Proxy DB Path
	Proxy string `json:"proxy"`

	// LOCAL CSV DB Path
	Local string `json:"local"`

	d *debugger.Debugger
}

// create new IP2location
func New(c *configurations.IP2Location, debugger *debugger.Debugger) *IPLocation {
	return &IPLocation{
		ASN:   c.ASN,
		IP:    c.IP,
		Proxy: c.IP,
		Local: c.Local,

		d: debugger,
	}
}
