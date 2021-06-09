package model

import "gorm.io/gorm"

// Protocol Model
// It uses for CRUD operation for 'Protocol'
type Protocol struct {
	gorm.Model

	Protocol string `gorm:"not null"`

	ProtocolName string `gorm:"not null;unique"`

	Info string
}
