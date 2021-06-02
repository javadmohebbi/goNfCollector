package model

import "gorm.io/gorm"

// AutonomousSystem Model
// It uses for CRUD operation for 'AutonomousSystem information'
type AutonomousSystem struct {
	gorm.Model

	ASN string `gorm:"not null"`

	Info string
}
