package configurations

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// this is a struct for collector.yml configuration
type Collector struct {
	confFile confFile

	Debug      bool      `json:"debug"`
	CPUNum     int       `json:"cpuNum"`
	AcceptFrom string    `json:"acceptFrom"`
	Listen     listen    `json:"listen"`
	Forwarder  forwarder `json:"forwarder"`
	LogFile    string    `json:"logFile"`

	IPReputation IpReputation `json:"ipReputation"`

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

		// get from os.GetEnv in case of error
		if ok := c.getMinialConfigsFromOSEnv(); ok {
			log.Println("can not read config from: ", c.confFile.path+c.confFile.file+"."+c.confFile.ext)
			log.Println("Minimal configuration will be read using OS environment")
			return c, nil
		}

		// if file not found set default fallback values
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Configuration %v.%v not found in path %v", c.confFile.path, c.confFile.file, c.confFile.ext)
			return c, nil
		}
		log.Fatalln("Could not read config file due to error: ", err)
		return c, err
	}

	// Unmarshal configs
	err := viper.Unmarshal(&c)

	// check foe unmarshal errors
	if err != nil {

		// get from os.GetEnv in case of error
		if ok := c.getMinialConfigsFromOSEnv(); ok {
			log.Println("Unable to unmarshal configs into struct due to error: ", err)
			log.Println("Minimal configuration will be read using OS environment")
			return c, nil
		}

		log.Fatalln("Unable to unmarshal configs into struct due to error: ", err)
		return c, err
	}

	return c, nil
}

// read minimal configs from OS ENV
func (c *Collector) getMinialConfigsFromOSEnv() bool {

	c.Debug, _ = strconv.ParseBool(os.Getenv("NFC_DEBUG"))

	c.Listen.Address = os.Getenv("NFC_LISTEN_ADDRESS")

	c.Listen.Port, _ = strconv.Atoi(os.Getenv("NFC_LISTEN_PORT"))

	c.IPReputation.IPSumPath = os.Getenv("NFC_IP_REPTATION_IPSUM")

	c.LogFile = os.Getenv("NFC_LOG_FILE")

	c.CPUNum, _ = strconv.Atoi(os.Getenv("NFC_CPU_NUM"))

	if len(c.Exporter.InfluxDBs) == 0 {

		tp, _ := strconv.Atoi(os.Getenv("NFC_INFLUXDB_PORT"))

		c.Exporter.InfluxDBs = append(c.Exporter.InfluxDBs, influxDB{
			Host:   os.Getenv("NFC_INFLUXDB_HOST"),
			Port:   tp,
			Token:  os.Getenv("NFC_INFLUXDB_TOKEN"),
			Bucket: os.Getenv("NFC_INFLUXDB_BUCKET"),
			Org:    os.Getenv("NFC_INFLUXDB_ORG"),
		})
	}

	if c.Listen.Address != "" && (c.Listen.Port > 0 && c.Listen.Port <= 65535) &&
		c.Exporter.InfluxDBs[0].Host != "" &&
		(c.Exporter.InfluxDBs[0].Port > 0 && c.Exporter.InfluxDBs[0].Port <= 65535) &&
		c.Exporter.InfluxDBs[0].Token != "" && c.Exporter.InfluxDBs[0].Bucket != "" &&
		c.Exporter.InfluxDBs[0].Org != "" {

		return true
	}

	return false

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
	Postgres  []postgres `json:"postgres"`
}

// postgres exporter struct
type postgres struct {
	// postgre host
	Host string `json:"host"`

	// postgres port
	Port int `json:"port"`

	// postgres user
	User string `json:"user"`

	// postgres password
	Password string `json:"password"`

	// postgres db
	DB string `json:"db"`
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

// ipReputation conf struct
type IpReputation struct {
	IPSumPath string `json:"ipSumPath"`
}
