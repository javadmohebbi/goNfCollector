package configurations

import (
	"log"

	"github.com/spf13/viper"
)

// this is a struct for trans.yml configuration
type Translations struct {
	confFile confFile

	FlowStartSysUpTime       string `json:"flowStartSysUpTime"`
	FlowEndSysUpTime         string `json:"flowEndSysUpTime"`
	OctetDeltaCount          string `json:"octetDeltaCount"`
	PacketDeltaCount         string `json:"packetDeltaCount"`
	IngressInterface         string `json:"ingressInterface"`
	EgressInterface          string `json:"egressInterface"`
	IpNextHopIPv4Address     string `json:"ipNextHopIPv4Address"`
	SourceIPv4Address        string `json:"sourceIPv4Address"`
	DestinationIPv4Address   string `json:"destinationIPv4Address"`
	ProtocolIdentifier       string `json:"protocolIdentifier"`
	SourceTransportPort      string `json:"sourceTransportPort"`
	DestinationTransportPort string `json:"destinationTransportPort"`
	TcpControlBits           string `json:"tcpControlBits"`
	FlowDirection            string `json:"flowDirection"`

	DestinationIPv4PrefixLength string `json:"destinationIPv4PrefixLength"`
	SourceIPv4PrefixLength      string `json:"sourceIPv4PrefixLength"`
}

// Read configs
func (c *Translations) Read() (interface{}, error) {

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
			log.Println("Unable to unmarshal all configs into struct due to error: ", err)
			log.Println("Default configuration will be set for unavailable fields in trans.yml")
			return c, nil
		}

		log.Fatalln("Unable to unmarshal configs into struct due to error: ", err)
		return c, err
	}

	return c, nil
}

// read minimal configs from OS ENV
func (c *Translations) getMinialConfigsFromOSEnv() bool {

	check := false

	if c.FlowDirection == "" {
		c.FlowStartSysUpTime = "flowStartSysUpTime, flowStartMilliseconds"
		check = true
	}

	if c.FlowEndSysUpTime == "" {
		c.FlowEndSysUpTime = "flowEndSysUpTime, flowEndMilliseconds"
		check = true
	}

	if c.OctetDeltaCount == "" {
		c.OctetDeltaCount = "octetDeltaCount, octetTotalCount"
		check = true
	}

	if c.PacketDeltaCount == "" {
		c.PacketDeltaCount = "packetDeltaCount, packetTotalCount"
		check = true
	}

	if c.IngressInterface == "" {
		c.IngressInterface = "ingressInterface, ingressPhysicalInterface"
		check = true
	}

	if c.EgressInterface == "" {
		c.EgressInterface = "egressInterface, egressPhysicalInterface"
		check = true
	}

	if c.IpNextHopIPv4Address == "" {
		c.IpNextHopIPv4Address = "ipNextHopIPv4Address"
		check = true
	}

	if c.SourceIPv4Address == "" {
		c.SourceIPv4Address = "sourceIPv4Address"
		check = true
	}

	if c.DestinationIPv4Address == "" {
		c.DestinationIPv4Address = "destinationIPv4Address"
		check = true
	}

	if c.ProtocolIdentifier == "" {
		c.ProtocolIdentifier = "protocolIdentifier"
		check = true
	}

	if c.SourceTransportPort == "" {
		check = true
		c.SourceTransportPort = "sourceTransportPort"
	}

	if c.DestinationTransportPort == "" {
		check = true
		c.DestinationTransportPort = "destinationTransportPort"
	}

	if c.TcpControlBits == "" {
		c.TcpControlBits = "tcpControlBits"
		check = true
	}

	if c.FlowDirection == "" {
		check = true
		c.FlowDirection = "flowDirection"
	}

	if c.SourceIPv4PrefixLength == "" {
		check = true
		c.SourceIPv4PrefixLength = "sourceIPv4PrefixLength"
	}

	if c.DestinationIPv4PrefixLength == "" {
		check = true
		c.DestinationIPv4PrefixLength = "destinationIPv4PrefixLength"
	}

	return check

}
