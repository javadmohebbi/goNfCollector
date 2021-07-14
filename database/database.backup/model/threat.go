package model

import "gorm.io/gorm"

// Threat Model
// It uses for CRUD operation for 'Threats information'
type Threat struct {
	gorm.Model

	Source string `gorm:"not null"`

	Kind string

	Acked         bool
	Closed        bool
	FalsePositive bool

	// FK on host
	HostID uint
	Host   Host
}
