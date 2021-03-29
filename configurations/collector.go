package configurations

import (
	"log"

	"github.com/spf13/viper"
)

// this is a struct for collector.yml configuration
type Collector struct {
	confFile confFile

	Debug     bool      `json:"debug"`
	Listen    listen    `json:"listen"`
	Forwarder forwarder `json:"forwarder"`
	LogFile   string    `json:"logFile"`

	AlienVault alienVault `json:"alienVault"`

	// netflow collector exporter
	Exporter exporter `json:"exporter"`
}

// Read configs
func (c *Collector) Read() (interface{}, error) {

	// name of config file (without extension)
	viper.SetConfigName(c.confFile.file)

	// REQUIRED if the config file does not have the extension in the name
	viper.SetConfigType(c.confFile.ext)

	// path to look for the config file in
	viper.AddConfigPath(c.confFile.path)

	// verbose to console
	log.Printf("Reading config from \"%v\" \n", c.confFile.path+c.confFile.file+"."+c.confFile.ext)

	// Read config file and check for errors
	if err := viper.ReadInConfig(); err != nil {
		// if file not found set default fallback values
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Logfile %v.%v not found in path %v", c.confFile.path, c.confFile.file, c.confFile.ext)
			return c, nil
		}
		log.Fatalln("Could not read config file due to error: ", err)
		return c, err
	}

	// Unmarshal configs
	err := viper.Unmarshal(&c)

	// check foe unmarshal errors
	if err != nil {
		log.Fatalln("Unable to unmarshal configs into struct due to error: ", err)
		return c, err
	}

	return c, nil
}

// Server listen configuration
type listen struct {

	// Listen Host, usually 0.0.0.0: all IP addresses on this host
	Address string `json:"address"`

	// UDP Port to listen on
	Port int `json:"port"`
}

// Forwarder configuration
type forwarder struct {

	// Enable/Disable forwareder
	Enabled bool `json:"enabled"`

	// host to forward
	Hosts []string `json:"hosts"`

	// UPD port to forward
	Port int `json:"port"`
}

// exporter struct
// will be used to export flow metrics to
type exporter struct {
	InfluxDBs []influxDB `json:"influxDBs"`
}

// influxDB exporter struct
type influxDB struct {
	// Influx DB Host
	Host string `json:"host"`

	// Influx DB Port
	Port int `json:"port"`

	// Influx DB Token
	Token string `json:"token"`

	// Influx DB Bucket
	Bucket string `json:"bucket"`

	// Influx DB Database Org
	Org string `json:"org"`
}

// alienVault conf struct
type alienVault struct {
	APIToken string `json:"apiToken"`
}
