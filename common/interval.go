package common

import (
	"fmt"
	"strconv"
)

// GetPGInterval will get intervals like
// 15s 30s 1m 2m 1h 6h 24h 1w 4d 1M 3M 1y
// and convert it to intervals that postgresDB needs
// s = seconds
// m = minutes
// h = hours
// w = weeks
// M = months
// y = years
// In case of any error it will show last 15minutes
func GetPGInterval(interval string) string {

	// check if if interger is more than 1 and add s to string
	plural := "s"

	// default interval string
	intvl := "minute"

	// default interval integer
	i := 15

	// default returning string
	ret := fmt.Sprintf("%d %s%s", i, intvl, plural)

	// convert string numericals to integer
	i, err := strconv.Atoi(interval[:len(interval)-1])

	// in case of any error return the default
	if err != nil {
		return ret
	}

	// if integer interval is 0 return the default
	if i == 0 {
		return ret
	} else if i == 1 {
		// if interger interval is 1 we don't need to pluralize the string
		plural = ""
	}

	// switch case for check the last character of requested interval
	// and fill the 'intvl' variable with correct string
	switch interval[len(interval)-1:] {
	case "s":
		intvl = "second"
	case "m":
		intvl = "minute"
	case "h":
		intvl = "hour"
	case "d":
		intvl = "day"
	case "w":
		intvl = "week"
	case "M":
		intvl = "month"
	case "y":
		intvl = "year"
	default:
		intvl = "minute"
		i = 15
		plural = "s"
	}

	return fmt.Sprintf("%d %s%s", i, intvl, plural)
}

// this function will get interval and
// return the time_trunc function needed variable
// for group by
func GetPGGroupByInterval(interval string) string {

	var intvl string

	// convert string numericals to integer
	// no need for error check
	i, _ := strconv.Atoi(interval[:len(interval)-1])
	if i == 0 {
		// default 15 minute
		// read doc for GetPGInterval
		i = 15
	}

	// switch case for check the last character of requested interval
	// and fill the 'intvl' variable with correct string
	switch interval[len(interval)-1:] {
	case "s":
		intvl = "second"
		if i > 60 {
			intvl = "minute"
		}
	case "m":
		intvl = "minute"
		if i > 60 {
			intvl = "hour"
		}
	case "h":
		if i > 24 {
			intvl = "day"
		}
		if i == 1 {
			intvl = "minute"
		}
		intvl = "hour"
	case "d":
		if i == 1 {
			intvl = "hour"
		} else {
			intvl = "day"
		}
	case "w":
		if i == 1 {
			intvl = "day"
		} else {
			intvl = "week"
		}
	case "M":
		if i == 1 {
			intvl = "week"
		} else {
			intvl = "month"
		}
	case "y":
		if i == 1 {
			intvl = "month"
		} else {
			intvl = "year"
		}
	default:
		intvl = "minute"
	}

	return intvl

}
