package nfv5

import (
	"fmt"
	"net"
	"time"

	"github.com/goNfCollector/common"
	"github.com/tehmaze/netflow/netflow5"
)

func Prepare(addr string, p *netflow5.Packet, portMap common.PortMap, portMapErr error) []common.Metric {
	nfExporter, _, _ := net.SplitHostPort(addr)
	var metrics []common.Metric
	var met common.Metric
	for _, r := range p.Records {
		met = common.Metric{OutBytes: "0", InBytes: "0", OutPacket: "0", InPacket: "0", Device: nfExporter}

		met.Time = p.Header.Unix
		met.Uptime = time.Duration(p.Header.SysUptime) * time.Nanosecond

		met.FlowVersion = "Netflow-V5"
		met.Direction = "Unsupported"
		met.First = fmt.Sprintf("%v", r.First)
		met.Last = fmt.Sprintf("%v", r.Last)
		met.Protocol = fmt.Sprintf("%v", r.Protocol)
		met.ProtoName = common.ProtoToName(met.Protocol)
		met.Bytes = fmt.Sprintf("%v", r.Bytes)
		met.Packets = fmt.Sprintf("%v", r.Packets)
		met.TCPFlags = fmt.Sprintf("%v", r.TCPFlags)
		met.SrcAs = fmt.Sprintf("%v", r.SrcAS)
		met.DstAs = fmt.Sprintf("%v", r.DstAS)
		met.SrcMask = fmt.Sprintf("%v", r.SrcMask)
		met.DstMask = fmt.Sprintf("%v", r.DstMask)

		met.NextHop = fmt.Sprintf("%v", r.NextHop)

		met.InEthernet = fmt.Sprintf("%v", r.Input)
		met.OutEthernet = fmt.Sprintf("%v", r.Output)
		met.ToS = fmt.Sprintf("%v", r.ToS)

		met.SrcIP = fmt.Sprintf("%v", r.SrcAddr)
		met.DstIP = fmt.Sprintf("%v", r.DstAddr)

		met.SrcPort = fmt.Sprintf("%v", r.SrcPort)
		met.SrcPortName = common.GetPortName(met.SrcPort, met.ProtoName, portMap, portMapErr)

		met.DstPort = fmt.Sprintf("%v", r.DstPort)
		met.DstPortName = common.GetPortName(met.DstPort, met.ProtoName, portMap, portMapErr)

		metrics = append(metrics, met)

	}
	return metrics
}
