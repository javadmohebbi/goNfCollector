package api

import (
	"strings"
	"sync"

	"github.com/goNfCollector/common"
)

type FilterLiveFlow struct {
	mu             sync.Mutex
	IsFilterEnable bool                 `json:"isFilterEnable"`
	Device         string               `json:"device"`
	IP             string               `json:"ip"`
	Port           string               `json:"port"`
	SrcOrDst       string               `json:"srcOrDst"` // src, dst, both
	Proto          string               `json:"proto"`
	Country        string               `json:"country"`
	Region         string               `json:"region"`
	City           string               `json:"city"`
	Flags          FilterLiveFlowFlags  `json:"flags"`
	Threat         FilterLiveFlowThreat `json:"threat"`
	FlowVersion    uint                 `json:"flowVersion"`
}
type FilterLiveFlowFlags struct {
	Filtered bool `json:"filtered"`
	List     struct {
		Fin bool `json:"fin"`
		Syn bool `json:"syn"`
		Rst bool `json:"rst"`
		Psh bool `json:"psh"`
		Ack bool `json:"ack"`
		Urg bool `json:"urg"`
		Ece bool `json:"ece"`
		Cwr bool `json:"cwr"`
	} `json:"list"`
}
type FilterLiveFlowThreat struct {
	Filtered bool `json:"filtered"`
	IsThreat bool `json:"isThreat"`
}

// Filter metrics based on provided filter
func (f *FilterLiveFlow) FilterMetrics(metrics *[]common.Metric) []common.Metric {
	var metricsToSend []common.Metric
	if f.IsFilterEnable {

		_metrics, _ := f.filterIP(metrics)
		// log.Println("=== LEN:", len(metrics), "   FIL LEN:", len(*_metrics), "   COUNT=", count)
		_metrics, _ = f.filterPort(_metrics)

		for _, m := range *_metrics {
			if m.IncludeFilterIP || m.IncludeFilterPort {
				metricsToSend = append(metricsToSend, m)
			}
		}

		// for _, m := range metrics {
		// 	// it will use OR instead of AND
		// 	// and thats why it will other IF
		// 	// in case prev. one is false
		// 	shouldInclude := false

		// 	// add to output if included in OR filter
		// 	if (shouldInclude) {
		// 		metricsToSend = append(metricsToSend, m)
		// 	}

		// }
		return metricsToSend
	} else {
		return *metrics
	}
}

// filter based on Provided Port
func (f *FilterLiveFlow) filterPort(metrics *[]common.Metric) (*[]common.Metric, uint) {
	count := uint(0)
	if f.Port != "" {
		var _metrics []common.Metric
		for _, m := range *metrics {
			_m := m
			switch strings.ToLower(f.SrcOrDst) {
			case "src":
				if strings.Contains(m.SrcPort, f.Port) || strings.Contains(m.SrcPortName, f.Port) {
					_m.IncludeFilterPort = true
					count++
				}
			case "dst":
				if strings.Contains(m.DstPort, f.Port) || strings.Contains(m.DstPortName, f.Port) {
					_m.IncludeFilterPort = true
					count++
				}
			case "both":
				if (strings.Contains(m.SrcPort, f.Port) || strings.Contains(m.SrcPortName, f.Port)) ||
					(strings.Contains(m.DstPort, f.Port) || strings.Contains(m.DstPortName, f.Port)) {
					_m.IncludeFilterPort = true
					count++
				}
			default:
				// nothing to do
			}
			_metrics = append(_metrics, _m)
		}
		return &_metrics, count
	}

	return metrics, count

}

// filter based on Provided IP
func (f *FilterLiveFlow) filterIP(metrics *[]common.Metric) (*[]common.Metric, uint) {
	count := uint(0)
	if f.IP != "" {
		var _metrics []common.Metric
		for _, m := range *metrics {
			_m := m
			switch strings.ToLower(f.SrcOrDst) {
			case "src":
				if strings.Contains(m.SrcIP, f.IP) {
					_m.IncludeFilterIP = true
					count++
				}
			case "dst":
				if strings.Contains(m.DstIP, f.IP) {
					_m.IncludeFilterIP = true
					count++
				}
			case "both":
				if strings.Contains(m.SrcIP, f.IP) ||
					strings.Contains(m.DstIP, f.IP) {
					_m.IncludeFilterIP = true
					count++
				}
			default:
				// nothing to do
			}
			_metrics = append(_metrics, _m)
		}
		return &_metrics, count
	}

	return metrics, count
}

// enbale filter
func (f *FilterLiveFlow) EnableFilter() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.IsFilterEnable = true
}

// disable filter
func (f *FilterLiveFlow) DisableFilter() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.IsFilterEnable = false
}

// change device filter
func (f *FilterLiveFlow) SetDeviceFilter(device string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Device = device
}

// change ip filter
func (f *FilterLiveFlow) SetIPFilter(ip string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.IP = ip
}

// change port filter
func (f *FilterLiveFlow) SetPortFilter(port string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Port = port
}

// change proto filter
func (f *FilterLiveFlow) SetProtoFilter(proto string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Proto = proto
}

// change SrcOrDst filter
// for filtering IP/PORT src or dst
func (f *FilterLiveFlow) SetSrcOrDst(srcOrDst string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.SrcOrDst = srcOrDst
}

// change country filter
func (f *FilterLiveFlow) SetCountry(country string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Country = country
}

// change region filter
func (f *FilterLiveFlow) SetRegion(region string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Region = region
}

// change city filter
func (f *FilterLiveFlow) SetCity(city string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.City = city
}

// change flags filter
func (f *FilterLiveFlow) SetFlag(fin, syn, rst, psh, ack, urg, ece, cwr, filtered bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Flags.Filtered = filtered
	f.Flags.List.Fin = fin
	f.Flags.List.Syn = syn
	f.Flags.List.Rst = rst
	f.Flags.List.Psh = psh
	f.Flags.List.Ack = ack
	f.Flags.List.Urg = urg
	f.Flags.List.Ece = ece
	f.Flags.List.Cwr = cwr
}

// change threats filter
func (f *FilterLiveFlow) SetThreat(isthreat, filtered bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Threat.Filtered = filtered
	f.Threat.IsThreat = isthreat
}

// change Flow Version filter
func (f *FilterLiveFlow) SetFlowVersion(flowVersion uint) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.FlowVersion = flowVersion
}
