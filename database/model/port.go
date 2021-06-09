package model

import "gorm.io/gorm"

// Port Model
// It uses for CRUD operation for 'Src/Dst port'
type Port struct {
	gorm.Model

	PortName string `gorm:"not null;unique"`

	// removed, because it was not usefull
	// Port uint `gorm:"not null"`

	PortProto string

	Info string
}
