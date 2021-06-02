package model

import "gorm.io/gorm"

// Flag Model
// It uses for CRUD operation for 'Flags'
type Flag struct {
	gorm.Model

	Flags string `gorm:"not null"`

	Info string
}
