package model

import "gorm.io/gorm"

// Device Model
// It uses for CRUD operation for 'Netflow Exporter' devices
type Device struct {
	gorm.Model

	Device string `gorm:"not null"`

	Name string

	Info string
}
