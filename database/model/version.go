package model

import "gorm.io/gorm"

// Version Model
// It uses for CRUD operation for 'Netflow Versions'
type Version struct {
	gorm.Model

	Version uint `gorm:"not null;unique"`
}
