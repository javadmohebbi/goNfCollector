package common

type Metric struct {
	FlowVersion         string `header:"NF Ver"`
	Device              string `header:"Device"`
	Last                string //`header:"Last"`
	First               string //`header:"First"`
	Bytes               string `header:"Bytes"`
	Packets             string `header:"Packets"`
	InBytes             string //`header:"Bytes In"`
	InPacket            string //`header:"Packets In"`
	OutBytes            string //`header:"Bytes Out"`
	OutPacket           string //`header:"Packets Out"`
	InEthernet          string //`header:"In Eth"`
	OutEthernet         string //`header:"Out Eth"`
	SrcIP               string `header:"SrcIP"`
	SrcIp2lCountryShort string //`header:"sCountry_S"`
	SrcIp2lCountryLong  string //`header:"sCountry_L"`
	SrcIp2lState        string //`header:"sState"`
	SrcIp2lCity         string //`header:"sCity"`
	SrcIp2lLat          string //`header:"sLat"`
	SrcIp2lLong         string //`header:"sLong"`
	DstIP               string `header:"DstIP"`
	DstIp2lCountryShort string //`header:"dCountry_S"`
	DstIp2lCountryLong  string //`header:"dCountry_L"`
	DstIp2lState        string //`header:"dState"`
	DstIp2lCity         string //`header:"dCity"`
	DstIp2lLat          string //`header:"dLat"`
	DstIp2lLong         string //`header:"dLong"`
	Protocol            string //`header:"Proto"`
	ProtoName           string `header:"ProtoName"`
	SrcToS              string //`header:"SrcToS"`
	SrcPort             string `header:"SrcPort"`
	SrcPortName         string `header:"SrcPortName"`
	DstPort             string `header:"DstPort"`
	DstPortName         string `header:"DstPortName"`
	FlowSamplerId       string //`header:"FlowSampleId"`
	VendorPROPRIETARY   string //`header:"VendorPROPRIETARY"`
	NextHop             string `header:"NextHop"`
	DstMask             string //`header:"DstMask"`
	SrcMask             string //`header:"SrcMask"`
	TCPFlags            string `header:"TCPFlags"`
	Direction           string //`header:"Direction"`
	DstAs               string //`header:"DstAs"`
	SrcAs               string //`header:"SrcAs"`
}
