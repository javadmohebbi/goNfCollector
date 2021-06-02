package model

import "gorm.io/gorm"

// Domain Model
// It uses for CRUD operation for 'Domain'
type Domain struct {
	gorm.Model

	Domain string `gorm:"not null"`
}
