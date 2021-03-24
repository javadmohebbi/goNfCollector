package configurations

import (
	"errors"
	"fmt"
)

// conf file struct to set path, file & extension
type confFile struct {
	path string
	file string
	ext  string
}

// define new configuraion interface
type Configuration interface {
	// Read Configuration file
	Read() (interface{}, error)
}

// Type of congfiguration file
type ConfType int

// enum const for type of configuration files
const (
	CONF_TYPE_COLLECTOR ConfType = iota
	CONF_TYPE_DATABASE
	CONF_TYPE_EXPORTER
	CONF_TYPE_API
	CONF_TYPE_WEB
)

// return filename related to requested configuration
func (ct ConfType) String() string {
	return [...]string{
		"collector",
		"database",
		"exporter",
		"api",
		"web",
	}[ct]
}

// create new configuration
func New(ct ConfType) (Configuration, error) {

	// define default configuration file path, name, ext
	cf := confFile{
		path: "/opt/nfcollector/etc/",
		file: "collector",
		ext:  "yml",
	}

	switch ct {
	case CONF_TYPE_COLLECTOR:
		return Configuration(
			&Collector{
				confFile: cf,
			},
		), nil
	default:
		return Configuration(&Collector{}), errors.New(fmt.Sprintf("No valid configuration type found"))
	}
}
