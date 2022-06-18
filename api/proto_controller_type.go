package api

// Port when src or DST
type ProtocolRPTWhenHostSrcOrDstResult struct {
	ProtocolID   uint   `json:"protocol_id" gorm:"protocol_id"`
	Protocol     string `json:"protocol" gorm:"protocol"`
	ProtocolName string `json:"protocol_name" gorm:"protocol_name"`
	TotalBytes   uint   `json:"total_bytes" gorm:"total_bytes"`
	TotalPackets uint   `json:"total_packets" gorm:"total_packets"`
}
