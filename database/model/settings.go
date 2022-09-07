package model

import (
	"time"

	"gorm.io/gorm"
)

// IP2Location Last status custom type
type I2LLastStatType uint

// Convert last update statuses to human-readable strings
func (i2lst I2LLastStatType) String() string {
	return [...]string{
		"Update not started!",
		"Token is invalid!",
		"Can not get access to ip2location web server!",
		"Unknown error!",
		"Update is in progress!",
		"Download completed!",
	}[i2lst-1001]
}

// Constants for the last status of IP2Location update
const (
	IP2L_LAST_STAT_NOT_STARTED   I2LLastStatType = 1001 + iota // 1001
	IP2L_LAST_STAT_TOKEN_ERROR                                 // 1002
	IP2L_LAST_STAT_HTTP_ERROR                                  // 1003
	IP2L_LAST_STAT_UNKNOWN_ERROR                               // 1004
	IP2L_LAST_STAT_IN_PROGRESS                                 // 1005
	IP2L_LAST_STAT_DONE_SUCCESS                                // 1006
)

// Settings model
// for storing settings related to
// netflow analyzer like IP2Location tokens, ...
type Settings struct {
	gorm.Model

	// token for ip2location
	// lite version is free and you could get it from
	// https://lite.ip2location.com
	IP2LToken      *string         `gorm:"ip2location_token"`
	IP2LLastUpdate *time.Time      `gorm:"ip2location_last_update"`
	IP2LLastStat   I2LLastStatType `gorm:"ip2location_last_stat"`
	IP2LProgress   uint            `gorm:"ip2location_progress"`
}
