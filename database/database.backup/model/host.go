package model

import "gorm.io/gorm"

// Host Model
// It uses for CRUD operation for 'Src/Dst hosts'
type Host struct {
	gorm.Model

	Host string `gorm:"not null;unique"`

	Info string
}
