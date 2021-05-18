package configurations

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// api server configuration
type APIServer struct {
	confFile confFile

	Debug          bool   `json:"debug"`
	LogFile        string `json:"logFile"`
	Listen         listen `json:"listen"`
	JWTSecretToken string `json:"-"`
	TLS            tls    `json:"tls"`
	HTTP           http   `json:"http"`
}

// http server extra cnfigurations
type http struct {
	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	ReadTimeOut int64 `josn:"readTimeOut"`

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	WriteTimeOut int64 `josn:"writeTimeOut"`

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, there is no timeout.
	IdleTimeOut int64 `josn:"idleTimeOut"`

	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what
	// is considered too slow for the body. If ReadHeaderTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, there is no timeout.
	ReadHeaderTimeOut int64 `josn:"readHeaderTimeOut"`

	CORS httpCors `json:"cors"`
}

// http cors configuraitons
type httpCors struct {

	// AllowAll create a new Cors handler with permissive configuration allowing all
	// origins with all standard methods with any header and credentials.
	AllowAll bool `json:"allowAll"`

	// AllowedHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [] but "Origin" is always appended to the list.
	AllowedHeaders []string `json:"allowedHeaders"`

	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods []string `json:"allowedMethods"`

	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// An origin may contain a wildcard (*) to replace 0 or more characters
	// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penalty.
	// Only one wildcard can be used per origin.
	// Default value is ["*"]
	AllowedOrigins []string `json:"allowedOrigins"`
}

// Read configs
func (c *APIServer) Read() (interface{}, error) {

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
func (c *APIServer) getMinialConfigsFromOSEnv() bool {

	c.Debug, _ = strconv.ParseBool(os.Getenv("NFC_API_DEBUG"))

	c.Listen.Address = os.Getenv("NFC_API_LISTEN_ADDRESS")

	c.Listen.Port, _ = strconv.Atoi(os.Getenv("NFC_API_LISTEN_PORT"))

	c.JWTSecretToken = os.Getenv("NFC_API_JWT_SECRET")

	c.LogFile = os.Getenv("NFC_API_LOG_FILE")

	if c.Listen.Address != "" && (c.Listen.Port > 0 && c.Listen.Port <= 65535) {
		return true
	}

	return false

}
