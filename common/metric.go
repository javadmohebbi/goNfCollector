package common

import "time"

type Metric struct {
	Time                time.Time     `json:"time"`
	Uptime              time.Duration `json:"uptime"`
	FlowVersion         string        `header:"NF Ver" json:"flowVersionName"`
	FlowVersionNumber   uint          `header:"NF Ver Num" json:"flowVersion"`
	Device              string        `header:"Device" json:"device"`
	Last                string        `header:"Last" json:"last"`
	First               string        `header:"First" json:"first"`
	Bytes               string        `header:"Bytes" json:"bytes"`
	Packets             string        `header:"Packets" json:"packets"`
	InBytes             string        `header:"Bytes In" json:"inBytes"`
	InPacket            string        `header:"Packets In" json:"inPackets"`
	OutBytes            string        `header:"Bytes Out" json:"outBytes"`
	OutPacket           string        `header:"Packets Out" json:"outPackets"`
	InEthernet          string        `header:"In Eth" json:"inEthernet"`
	OutEthernet         string        `header:"Out Eth" json:"outEthernet"`
	SrcIP               string        `header:"SrcIP" json:"srcIP"`
	SrcIp2lCountryShort string        `header:"sCountry_S" json:"sCountryShort"`
	SrcIp2lCountryLong  string        `header:"sCountry_L" json:"sCountryLong"`
	SrcIp2lState        string        `header:"sState" json:"sState"`
	SrcIp2lCity         string        `header:"sCity" json:"sCity"`
	SrcIp2lLat          string        `header:"sLat" json:"sLat"`
	SrcIp2lLong         string        `header:"sLong" json:"sLOng"`
	DstIP               string        `header:"DstIP" json:"dstIP"`
	DstIp2lCountryShort string        `header:"dCountry_S" json:"dCountryShort"`
	DstIp2lCountryLong  string        `header:"dCountry_L" json:"dCountryLong"`
	DstIp2lState        string        `header:"dState" json:"dState"`
	DstIp2lCity         string        `header:"dCity" json:"dCity"`
	DstIp2lLat          string        `header:"dLat" json:"dLat"`
	DstIp2lLong         string        `header:"dLong" json:"dLong"`
	Protocol            string        `header:"Proto" json:"proto"`
	ProtoName           string        `header:"ProtoName" json:"protoName"`
	ToS                 string        `header:"ToS" json:"ToS"`
	SrcPort             string        `header:"SrcPort" json:"srcPort"`
	SrcPortName         string        `header:"SrcPortName" json:"srcPortName"`
	DstPort             string        `header:"DstPort" json:"dstPort"`
	DstPortName         string        `header:"DstPortName" json:"dstPortName"`
	FlowSamplerId       string        `header:"FlowSampleId" json:"flowSampleId"`
	VendorPROPRIETARY   string        `header:"VendorPROPRIETARY" json:"vendorPropretary"`
	NextHop             string        `header:"NextHop" json:"nextHop"`
	DstMask             string        `header:"DstMask" json:"dstMask"`
	SrcMask             string        `header:"SrcMask" json:"srcMask"`
	TCPFlags            string        `header:"TCPFlags" json:"tcpFlags"`
	Direction           string        `header:"Direction" json:"direction"`
	DstAs               string        `header:"DstAs" json:"dstAs"`
	SrcAs               string        `header:"SrcAs" json:"srcAs"`

	// this are sent to UI on web socket
	FlagFin bool `json:"fin"`
	FlagSyn bool `json:"syn"`
	FlagRst bool `json:"rst"`
	FlagPsh bool `json:"psh"`
	FlagAck bool `json:"ack"`
	FlagUrg bool `json:"urg"`
	FlagEce bool `json:"ece"`
	FlagCwr bool `json:"cwr"`
}
