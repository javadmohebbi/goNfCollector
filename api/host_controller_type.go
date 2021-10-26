package api

import "time"

// Host Report Threat Result
type HostRPTThreatsResult struct {
	ThreatID            uint      `json:"threat_id" gorm:"threat_id"`
	ThreatSource        string    `json:"threat_source" gorm:"threat_source"`
	ThreatCounter       uint      `json:"threat_counter" gorm:"threat_counter"`
	ThreatReputation    uint      `json:"threat_reputation" gorm:"threat_reputation"`
	ThreatKind          string    `json:"threat_kind" gorm:"threat_kind"`
	ThreatHost          string    `json:"threat_host" gorm:"threat_host"`
	ThreatHostInfo      string    `json:"threat_host_info" gorm:"threat_host_info"`
	ThreatHostID        uint      `json:"threat_host_id" gorm:"threat_host_id"`
	ThreatAcked         bool      `json:"threat_acked" gorm:"threat_acked"`
	ThreatClosed        bool      `json:"threat_closed" gorm:"threat_closed"`
	ThreatFalsePositive bool      `json:"threat_false_positive" gorm:"threat_false_positive"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"updated_at"`
	CreatedAt           time.Time `json:"created_at" gorm:"created_at"`
}

// host when src
type HostRPTWhenSrcOrDstResult struct {
	HostID       uint   `json:"host_id" gorm:"host_id"`
	Host         string `json:"host" gorm:"host"`
	HostInfo     string `json:"host_info" gorm:"host_info"`
	TotalBytes   uint   `json:"total_bytes" gorm:"total_bytes"`
	TotalPackets uint   `json:"total_packets" gorm:"total_packets"`
}
