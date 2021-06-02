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

	// Can not open or create http access log file
	ERROR_CAN_T_OPEN_OR_CREATE_HTTP_ACCESS_LOG_FILE

	// Can not open or create http error log file
	ERROR_CAN_T_OPEN_OR_CREATE_HTTP_ERROR_LOG_FILE

	// Can not listen & serve HTTP server
	ERROR_CAN_T_LISTEN_AND_SERVE_HTTP_SERVER

	// Can not stop api HTTP server
	ERROR_CAN_T_STOP_API_HTTP_SERVER

	// Can not insert device info
	ERROR_CAN_T_INSERT_DEVICE_INFO

	// Can not update device info
	ERROR_CAN_T_UPDATE_DEVICE_INFO

	// Can not insert metrics to postgres DB
	ERROR_CAN_T_INSERT_METRICS_TO_POSTGRES_DB

	// No metrics in the array
	ERROR_NO_METRICS_IN_THE_ARRAY

	// Can not insert version info
	ERROR_CAN_T_INSERT_VERSION_INFO
	// Can not update version info
	ERROR_CAN_T_UPDATE_VERSION_INFO

	// Can not insert protocol info
	ERROR_CAN_T_INSERT_PROTOCOL_INFO
	// Can not update protocol info
	ERROR_CAN_T_UPDATE_PROTOCOL_INFO

	// Can not insert autonomous system info
	ERROR_CAN_T_INSERT_AUTONOMOUS_INFO
	// Can not update autonomous system info
	ERROR_CAN_T_UPDATE_AUTONOMOUS_INFO

	// Can not insert host info
	ERROR_CAN_T_INSERT_HOST_INFO
	// Can not update host info
	ERROR_CAN_T_UPDATE_HOST_INFO

	// Can not insert port info
	ERROR_CAN_T_INSERT_PORT_INFO
	// Can not update port info
	ERROR_CAN_T_UPDATE_PORT_INFO

	// Can not insert GEO info
	ERROR_CAN_T_INSERT_GEO_INFO
	// Can not update GEO info
	ERROR_CAN_T_UPDATE_GEO_INFO

	// Can not insert flag info
	ERROR_CAN_T_INSERT_FLAG_INFO
	// Can not update flag info
	ERROR_CAN_T_UPDATE_FLAG_INFO

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

		"Can not open or create http access log file",
		"Can not open or create http error log file",

		"Can not listen & serve API HTTP server",
		"Can not stop API HTTP server",

		"Can not insert device info",
		"Can not update device info",
		"Can not insert metrics to postgres DB",
		"No metrics in the array",

		"Can not insert version info",
		"Can not update version info",

		"Can not insert protocol info",
		"Can not update protocol info",

		"Can not insert autonomous system info",
		"Can not update autonomous system info",

		"Can not insert host info",
		"Can not update host info",

		"Can not insert port info",
		"Can not update port info",

		"Can not insert geo info",
		"Can not update geo info",

		"Can not insert flag info",
		"Can not update flag info",

		"Unknown error",
	}[e]
}
