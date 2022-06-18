package api

// Port when src or DST
type PortRPTWhenHostSrcOrDstResult struct {
	PortID       uint   `json:"port_id" gorm:"port_id"`
	PortProto    string `json:"port_proto" gorm:"port_proto"`
	PortName     string `json:"port_name" gorm:"port_name"`
	TotalBytes   uint   `json:"total_bytes" gorm:"total_bytes"`
	TotalPackets uint   `json:"total_packets" gorm:"total_packets"`
}
