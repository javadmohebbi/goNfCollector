package configurations

type ErrorCodes int

const (

	// no error
	NO_ERROR ErrorCodes = iota

	// READ CONFIG ERROR
	ERROR_READ_CONFIG

	// can not resolve the provided configuration
	// to udpAddress
	ERROR_LISTEN_RESOLVE_UDP_ADDRESS

	// Can not listen to the provided host:port
	ERROR_LISTEN_ON_UDP

	// can not set connection read buffer
	ERROR_CAN_T_SET_CONNECTION_READ_BUFFER

	// can not read recieved data from exporter device
	ERROR_CAN_T_READ_DATA

	// can not decode netflow data
	ERROR_CAN_T_DECODE_NETFLOW_DATA

	// It's not a valid exporter
	ERROR_NO_VALID_EXPORTER_FOUND

	// Unable to open ip2location database
	ERROR_OPEN_IP2LOCATION_DB

	// Can't get IP Location information due to error
	ERROR_GET_IP2LOCATION_INFO

	// Can not open ASN CSV DB
	ERROR_CAN_T_OPEN_ASN_DB

	// Can not Read ASN CSV DB
	ERROR_CAN_T_READ_ASN_DB

	// Can not open proxy CSV DB
	ERROR_CAN_T_OPEN_PROXY_DB

	// Can not Read proxy CSV DB
	ERROR_CAN_T_READ_PROXY_DB

	// Can not open local CSV DB
	ERROR_CAN_T_OPEN_LOCAL_DB

	// Can not Read Local CSV DB
	ERROR_CAN_T_READ_LOCAL_DB

	// Can not connect to Postgres DB
	ERROR_CAN_T_CONNECT_TO_POSTGRES_DB

	// Can not initialize postgres DB
	ERROR_CAN_T_INIT_POSTGRES_DB

	// UNKOWN ERORR
	ERROR_UKNOWN
)

func (e ErrorCodes) Int() int {
	return int(e)
}

// error codes to string
func (e ErrorCodes) String() string {
	return [...]string{
		"No Error!",
		"Can not read config file",
		"Can not resolve the provided configuration host:port to UDP address",
		"Can not listen to the provided host:port",
		"Can not set connection read buffer",
		"Can not read recieved data from exporter device",
		"Can not decode netflow data",

		"It's not a valid exporter",

		"Unable to open IP2Location DB",
		"Can't get IP Location information due to error",
		"Can not open ASN CSV DB",
		"Can not read ASN CSV DB",
		"Can not open PROXY CSV DB",
		"Can not read PROXY CSV DB",

		"Can not open Local CSV DB",
		"Can not read Local CSV DB",

		"Can not connect to Postgres DB",
		"Can not initialize postgres DB",

		"Unknown error",
	}[e]
}
