package configurations

import (
	"log"
	"os"

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

		// get from os.GetEnv in case of error
		if ok := c.getMinialConfigsFromOSEnv(); ok {
			log.Println("can not read config from: ", c.confFile.path+c.confFile.file+"."+c.confFile.ext)
			log.Println("Minimal configuration will be read using OS environment")
			return c, nil
		}

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
func (c *IP2Location) getMinialConfigsFromOSEnv() bool {

	c.ASN = os.Getenv("NFC_IP2L_ASN")
	c.IP = os.Getenv("NFC_IP2L_IP")
	c.Proxy = os.Getenv("NFC_IP2L_PROXY")
	c.Local = os.Getenv("NFC_IP2L_LOCAL")

	if c.ASN != "" &&
		c.IP != "" &&
		c.Proxy != "" &&
		c.Local != "" {

		return true
	}

	return false

}
