package model

import (
	"time"
)

// all foregin keys as removed to test
// if insert performance has improved
// also ID and deleted_at has removed for preventing
// GORM automatic index

// Flow Model
// It uses for CRUD operation for 'Flow details'
type Flow struct {

	// removed from this model
	// gorm.Model

	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt gorm.DeletedAt `gorm:"index"`

	Device string
	// DeviceID uint
	// Device   Device

	Versio string
	// VersionID uint
	// Version   Version

	Protocol     string
	ProtocolName string
	// ProtocolID uint
	// Protocol   Protocol

	InEthernetIndex string
	// -1 for undefined
	// InEthernetID int
	// InEthernet   Ethernet `gorm:"foreignKey:InEthernetID"`

	OutEthernetIndex string
	// -1 for undefined
	// OutEthernetID int
	// OutEthernet   Ethernet `gorm:"foreignKey:OutEthernetID"`

	SrcASN string
	// SrcASNID uint
	// SrcASN    AutonomousSystem `gorm:"foreignKey:SrcASNID"`

	SrcMac string

	SrcHost string
	// SrcHostID uint
	// SrcHost   Host `gorm:"foreignKey:SrcHostID"`

	SrcPortName  string
	SrcPortProto string
	// SrcPortID uint
	// SrcPort   Port `gorm:"foreignKey:SrcPortID"`

	SrcCountryShort string
	SrcCountryLong  string
	SrcRegion       string
	SrcCity         string
	SrcLatitude     string
	SrcLongitude    string
	// SrcGeoID uint
	// SrcGeo    Geo `gorm:"foreignKey:SrcGeoID"`

	SrcIsThreat     bool
	SrcThreatSource string
	// SrcThreatID *uint `gotm:"null"`
	// SrcThreat   Threat `gorm:"foreignKey:SrcThreatID"`

	DstASN string
	// DstASNID uint
	// DstASN    AutonomousSystem `gorm:"foreignKey:DstASNID"`

	DstMac string

	DstHost string
	// DstHostID uint
	// DstHost   Host `gorm:"foreignKey:DstHostID"`

	DstPortName  string
	DstPortProto string
	// DstPortID uint
	// DstPort   Port `gorm:"foreignKey:DstPortID"`

	DstCountryShort string
	DstCountryLong  string
	DstRegion       string
	DstCity         string
	DstLatitude     string
	DstLongitude    string
	// DstGeoID uint
	// DstGeo    Geo `gorm:"foreignKey:DstGeoID"`

	DstIsThreat     bool
	DstThreatSource string
	// DstThreatID *uint `gotm:"null"`
	// DstThreat   Threat `gorm:"foreignKey:DstThreatID"`

	NextHop string
	// NextHopID uint
	// NextHop      Host `gorm:"foreignKey:NextHopID"`
	NextHopCountryShort string
	NextHopCountryLong  string
	NextHopRegion       string
	NextHopCity         string
	NextHopLatitude     string
	NextHopLongitude    string
	// NextHopGeoID uint
	// NextHopGeo   Geo `gorm:"foreignKey:NextHopGeoID"`

	NextHopIsThreat     bool
	NextHopThreatSource string
	// NextHopThreatID *uint `gotm:"null"`
	// NextHopThreat   Threat `gorm:"foreignKey:NextHopThreatID"`

	FlagID uint
	// Flag   Flag `gorm:"foreignKey:FlagID"`

	FlagFin bool
	FlagSyn bool
	FlagRst bool
	FlagPsh bool
	FlagAck bool
	FlagUrg bool
	FlagEce bool
	FlagCwr bool

	Byte   uint
	Packet uint
}
