package configurations

import (
	"log"

	"github.com/spf13/viper"
)

// this is a struct for ip2location.yml configuration
type IP2Location struct {
	confFile confFile

	// ASN DB
	ASN string `json:"asn"`

	// IP DB
	IP string `json:"ip"`

	// Proxy DB
	Proxy string `json:"proxy"`

	// local CSV db
	Local string `json:"local"`
}

// Read configs
func (c *IP2Location) Read() (interface{}, error) {

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
