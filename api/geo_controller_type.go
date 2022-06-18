package api

// Geo when src or DST
type GeoRPTWhenHostSrcOrDstResult struct {
	GeoID        uint   `json:"geo_id" gorm:"geo_id"`
	CountryShort string `json:"country_short" gorm:"country_short"`
	CountryLong  string `json:"country_long" gorm:"country_long"`
	Region       string `json:"region" gorm:"region"`
	City         string `json:"city" gorm:"city"`
	TotalBytes   uint   `json:"total_bytes" gorm:"total_bytes"`
	TotalPackets uint   `json:"total_packets" gorm:"total_packets"`
}
