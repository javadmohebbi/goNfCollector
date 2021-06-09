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

	DeviceID uint
	// Device   Device

	VersionID uint
	// Version   Version

	ProtocolID uint
	// Protocol   Protocol

	// -1 for undefined
	InEthernetID int
	// InEthernet   Ethernet `gorm:"foreignKey:InEthernetID"`

	// -1 for undefined
	OutEthernetID int
	// OutEthernet   Ethernet `gorm:"foreignKey:OutEthernetID"`

	SrcASNID uint
	// SrcASN    AutonomousSystem `gorm:"foreignKey:SrcASNID"`
	SrcHostID uint
	// SrcHost   Host `gorm:"foreignKey:SrcHostID"`
	SrcPortID uint
	// SrcPort   Port `gorm:"foreignKey:SrcPortID"`
	SrcGeoID uint
	// SrcGeo    Geo `gorm:"foreignKey:SrcGeoID"`

	SrcIsThreat bool
	SrcThreatID *uint `gotm:"null"`
	// SrcThreat   Threat `gorm:"foreignKey:SrcThreatID"`

	DstASNID uint
	// DstASN    AutonomousSystem `gorm:"foreignKey:DstASNID"`
	DstHostID uint
	// DstHost   Host `gorm:"foreignKey:DstHostID"`
	DstPortID uint
	// DstPort   Port `gorm:"foreignKey:DstPortID"`
	DstGeoID uint
	// DstGeo    Geo `gorm:"foreignKey:DstGeoID"`

	DstIsThreat bool
	DstThreatID *uint `gotm:"null"`
	// DstThreat   Threat `gorm:"foreignKey:DstThreatID"`

	NextHopID uint
	// NextHop      Host `gorm:"foreignKey:NextHopID"`
	NextHopGeoID uint
	// NextHopGeo   Geo `gorm:"foreignKey:NextHopGeoID"`

	NextHopIsThreat bool
	NextHopThreatID *uint `gotm:"null"`
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
