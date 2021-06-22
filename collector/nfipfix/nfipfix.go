package nfipfix

import (
	"time"

	"github.com/goNfCollector/common"
	"github.com/tehmaze/netflow/ipfix"

	"net"

	"strings"

	"fmt"
)

func Prepare(addr string, m *ipfix.Message, portMap common.PortMap, portMapErr error) []common.Metric {
	nfExporter, _, _ := net.SplitHostPort(addr)

	var metrics []common.Metric
	var met common.Metric

	for _, ds := range m.DataSets {
		if ds.Records == nil {
			continue
		}

		for _, dr := range ds.Records {

			check := make(map[string]bool)

			// not check true one for now
			check["flowEndSysUpTime"] = true
			check["flowStartSysUpTime"] = true
			check["octetDeltaCount"] = true
			check["packetDeltaCount"] = true
			check["ingressInterface"] = true
			check["egressInterface"] = true
			check["ipNextHopIPv4Address"] = true
			check["sourceIPv4Address"] = false
			check["destinationIPv4Address"] = false
			check["protocolIdentifier"] = false
			check["sourceTransportPort"] = false
			check["destinationTransportPort"] = false
			check["destinationIPv4PrefixLength"] = true
			check["sourceIPv4PrefixLength"] = true
			check["tcpControlBits"] = true
			check["flowDirection"] = false

			met = common.Metric{OutBytes: "0", InBytes: "0", OutPacket: "0", InPacket: "0", Device: nfExporter}

			met.Time = time.Unix(int64(m.Header.ExportTime), 0)
			// no up time

			met.FlowVersion = "IPFIX"
			for _, f := range dr.Fields {

				if f.Translated != nil {
					if f.Translated.Name != "" {

						// fmt.Printf("        NN %s: %v\n", f.Translated.Name, f.Translated.Value)

						switch strings.ToLower(f.Translated.Name) {
						case strings.ToLower("flowEndSysUpTime"):
							check[f.Translated.Name] = true
							met.First = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("flowStartSysUpTime"):
							check[f.Translated.Name] = true
							met.Last = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("octetDeltaCount"):
							check[f.Translated.Name] = true
							met.Bytes = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("packetDeltaCount"):
							check[f.Translated.Name] = true
							met.Packets = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("ingressInterface"):
							check[f.Translated.Name] = true
							met.InEthernet = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("egressInterface"):
							check[f.Translated.Name] = true
							met.OutEthernet = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("sourceIPv4Address"):
							check[f.Translated.Name] = true
							met.SrcIP = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("destinationIPv4Address"):
							check[f.Translated.Name] = true
							met.DstIP = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("protocolIdentifier"):
							check[f.Translated.Name] = true
							met.Protocol = fmt.Sprintf("%v", f.Translated.Value)
							met.ProtoName = common.ProtoToName(met.Protocol)

						case strings.ToLower("sourceTransportPort"):
							check[f.Translated.Name] = true
							met.SrcPort = fmt.Sprintf("%v", f.Translated.Value)
							met.SrcPortName = common.GetPortName(met.SrcPort, met.ProtoName, portMap, portMapErr)

						case strings.ToLower("destinationTransportPort"):
							check[f.Translated.Name] = true
							met.DstPort = fmt.Sprintf("%v", f.Translated.Value)
							met.DstPortName = common.GetPortName(met.DstPort, met.ProtoName, portMap, portMapErr)

						case strings.ToLower("ipNextHopIPv4Address"):
							check[f.Translated.Name] = true
							met.NextHop = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("destinationIPv4PrefixLength"):
							check[f.Translated.Name] = true
							met.DstMask = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("sourceIPv4PrefixLength"):
							check[f.Translated.Name] = true
							met.SrcMask = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("tcpControlBits"):
							check[f.Translated.Name] = true
							met.TCPFlags = fmt.Sprintf("%v", f.Translated.Value)

						case strings.ToLower("flowDirection"):
							check[f.Translated.Name] = true
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

			notAppend := true

			for _, v := range check {
				if !v {
					notAppend = false
					break
				}
			}

			// for k, v := range check {
			// 	fmt.Printf("\n\n=== %v= %v, ", k, v)
			// }
			// fmt.Println("")

			if notAppend {
				metrics = append(metrics, met)
			}
		}
	}

	return metrics
}
