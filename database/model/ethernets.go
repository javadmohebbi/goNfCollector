package model

import "gorm.io/gorm"

// Ethernet Model
// It uses for CRUD operation for 'Ethernet' devices
type Ethernet struct {
	gorm.Model

	DeviceID uint
	// Device   Device `gorm:"foreignKey:DeviceID"`

	// Just to prevent duplicate entries
	UniqKey string `gorm:"not null;unique"`

	Ethernet uint `gorm:"not null"`

	Name string

	Info string
}
