package model

import "gorm.io/gorm"

// Geo Model
// It uses for CRUD operation for 'Geo Locations' (IP2Location)
type Geo struct {
	gorm.Model

	CountryShort string `gorm:"not null"`
	CountryLong  string `gorm:"not null"`

	Region string `gorm:"not null"`

	City string `gorm:"not null"`

	Latitude  float32
	Longitude float32
}
