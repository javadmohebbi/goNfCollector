package nfipfix

import (
	"time"

	"github.com/goNfCollector/common"
	"github.com/tehmaze/netflow/ipfix"

	"net"

	"strings"

	"fmt"
)

func Prepare(addr string, m *ipfix.Message) []common.Metric {
	nfExporter, _, _ := net.SplitHostPort(addr)

	var metrics []common.Metric
	var met common.Metric

	for _, ds := range m.DataSets {
		if ds.Records == nil {
			continue
		}

		for _, dr := range ds.Records {

			// check := make(map[string]bool)
			// check["flowEndSysUpTime"] = false
			// check["flowStartSysUpTime"] = false
			// check["octetDeltaCount"] = false
			// check["packetDeltaCount"] = false
			// check["ingressInterface"] = false
			// check["egressInterface"] = false
			// check["sourceIPv4Address"] = false
			// check["destinationIPv4Address"] = false
			// check["protocolIdentifier"] = false
			// check["sourceTransportPort"] = false
			// check["destinationTransportPort"] = false
			// check["ipNextHopIPv4Address"] = false
			// check["destinationIPv4PrefixLength"] = false
			// check["sourceIPv4PrefixLength"] = false
			// check["tcpControlBits"] = false
			// check["flowDirection"] = false

			met = common.Metric{OutBytes: "0", InBytes: "0", OutPacket: "0", InPacket: "0", Device: nfExporter}

			met.Time = time.Unix(int64(m.Header.ExportTime), 0)
			// no up time

			met.FlowVersion = "IPFIX"
			for _, f := range dr.Fields {

				if f.Translated != nil {
					if f.Translated.Name != "" {

						// fmt.Printf("        NN %s: %v\n", f.Translated.Name, f.Translated.Value)

						// check[f.Translated.Name] = true

						switch strings.ToLower(f.Translated.Name) {
						case strings.ToLower("flowEndSysUpTime"):

							met.First = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("flowStartSysUpTime"):
							met.Last = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("octetDeltaCount"):
							met.Bytes = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("packetDeltaCount"):
							met.Packets = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("ingressInterface"):
							met.InEthernet = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("egressInterface"):
							met.OutEthernet = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("sourceIPv4Address"):
							met.SrcIP = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("destinationIPv4Address"):
							met.DstIP = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("protocolIdentifier"):
							met.Protocol = fmt.Sprintf("%v", f.Translated.Value)
							met.ProtoName = common.ProtoToName(met.Protocol)

						case strings.ToLower("sourceTransportPort"):
							met.SrcPort = fmt.Sprintf("%v", f.Translated.Value)
							met.SrcPortName = common.GetPortName(met.SrcPort, met.ProtoName)

						case strings.ToLower("destinationTransportPort"):
							met.DstPort = fmt.Sprintf("%v", f.Translated.Value)
							met.DstPortName = common.GetPortName(met.DstPort, met.ProtoName)

						case strings.ToLower("ipNextHopIPv4Address"):
							met.NextHop = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("destinationIPv4PrefixLength"):
							met.DstMask = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("sourceIPv4PrefixLength"):
							met.SrcMask = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("tcpControlBits"):
							met.TCPFlags = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("flowDirection"):
							met.Direction = fmt.Sprintf("%v", f.Translated.Value)
							switch met.Direction {
							case "0":
								met.Direction = "Ingress"
							case "1":
								met.Direction = "Egress"
							default:
								met.Direction = "Unsupported"
							}
						}
					} else {
						//fmt.Printf("        TT %d: %v\n", f.Translated.Type, f.Bytes)
						// return nil
						continue
					}
				} else {
					//fmt.Printf("        RR %d: %v (raw)\n", f.Type, f.Bytes)
					continue
				}
			}

			metrics = append(metrics, met)
		}
	}

	return metrics
}
