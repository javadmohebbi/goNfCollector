package configurations

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Socket struct {
	confFile confFile

	Socket string `json:"socket"`
}

// Read configs
func (c *Socket) Read() (interface{}, error) {

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
			log.Fatalf("Configuration file %v.%v not found in path %v", c.confFile.path, c.confFile.file, c.confFile.ext)
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
func (c *Socket) getMinialConfigsFromOSEnv() bool {

	c.Socket = os.Getenv("NFC_SOCK_PATH")

	if c.Socket != "" {
		return true
	}

	return false

}
