package nfipfix

import (
	"time"

	"github.com/goNfCollector/common"
	"github.com/goNfCollector/configurations"
	"github.com/tehmaze/netflow/ipfix"

	"net"

	"fmt"
)

func Prepare(addr string, m *ipfix.Message, portMap common.PortMap, portMapErr error, cfTrans *configurations.Translations) []common.Metric {
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

						switch f.Translated.Name {
						case common.CheckTranslationField("flowEndSysUpTime", cfTrans.FlowEndSysUpTime):
							met.First = fmt.Sprintf("%v", f.Translated.Value)

						case common.CheckTranslationField("flowStartSysUpTime", cfTrans.FlowStartSysUpTime):
							met.Last = fmt.Sprintf("%v", f.Translated.Value)

						case common.CheckTranslationField("octetDeltaCount", cfTrans.OctetDeltaCount):
							met.Bytes = fmt.Sprintf("%v", f.Translated.Value)

						case common.CheckTranslationField("packetDeltaCount", cfTrans.PacketDeltaCount):
							met.Packets = fmt.Sprintf("%v", f.Translated.Value)

						case common.CheckTranslationField("ingressInterface", cfTrans.IngressInterface):
							met.InEthernet = fmt.Sprintf("%v", f.Translated.Value)

						case common.CheckTranslationField("egressInterface", cfTrans.EgressInterface):
							met.OutEthernet = fmt.Sprintf("%v", f.Translated.Value)

						case common.CheckTranslationField("sourceIPv4Address", cfTrans.SourceIPv4Address):
							met.SrcIP = fmt.Sprintf("%v", f.Translated.Value)

						case common.CheckTranslationField("destinationIPv4Address", cfTrans.DestinationIPv4Address):
							met.DstIP = fmt.Sprintf("%v", f.Translated.Value)

						case common.CheckTranslationField("protocolIdentifier", cfTrans.ProtocolIdentifier):
							met.Protocol = fmt.Sprintf("%v", f.Translated.Value)
							met.ProtoName = common.ProtoToName(met.Protocol)

						case common.CheckTranslationField("sourceTransportPort", cfTrans.SourceTransportPort):
							met.SrcPort = fmt.Sprintf("%v", f.Translated.Value)
							met.SrcPortName = common.GetPortName(met.SrcPort, met.ProtoName, portMap, portMapErr)

						case common.CheckTranslationField("destinationTransportPort", cfTrans.DestinationTransportPort):
							met.DstPort = fmt.Sprintf("%v", f.Translated.Value)
							met.DstPortName = common.GetPortName(met.DstPort, met.ProtoName, portMap, portMapErr)

						case common.CheckTranslationField("ipNextHopIPv4Address", cfTrans.IpNextHopIPv4Address):
							met.NextHop = fmt.Sprintf("%v", f.Translated.Value)

						case common.CheckTranslationField("destinationIPv4PrefixLength", cfTrans.DestinationIPv4PrefixLength):
							met.DstMask = fmt.Sprintf("%v", f.Translated.Value)

						case common.CheckTranslationField("sourceIPv4PrefixLength", cfTrans.SourceIPv4PrefixLength):
							met.SrcMask = fmt.Sprintf("%v", f.Translated.Value)

						case common.CheckTranslationField("tcpControlBits", cfTrans.TcpControlBits):
							met.TCPFlags = fmt.Sprintf("%v", f.Translated.Value)

						case common.CheckTranslationField("flowDirection", cfTrans.FlowDirection):
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
