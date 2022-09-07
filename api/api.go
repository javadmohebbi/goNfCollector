package api

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	socketio "github.com/googollee/go-socket.io"

	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/database"
	"github.com/goNfCollector/debugger"
	"github.com/goNfCollector/location"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// API Server struct to prepare more
// integrity functions to goNfCollector App
// This API server requires POSTGRES exporter
// to be enabled & at least one postgres server
// MUST be defined in collector.yml
type APIServer struct {
	// API server host
	host string

	// API server port
	port int

	// logrus for futture use, not today :-D
	l *logrus.Logger

	// configuration for collector
	c *configurations.Collector

	apiConf *configurations.APIServer

	// debugger for verbosing the logs
	d *debugger.Debugger

	// postgres db
	pgdb database.Postgres

	// gorm DB
	db *gorm.DB

	ip2l *location.IPLocation

	// httpLogfiles
	httpAccessHasError bool
	httpAccessLog      *os.File

	httpErrorHasError bool
	httpErrorLog      *os.File

	// channel
	ch chan os.Signal

	// wait group
	waitGroup *sync.WaitGroup

	// httpServer
	httpSrv *http.Server

	apiSocketServer *socketio.Server
}

// Create new HTTP server
// it will read host:port from configuration file
func New(l *logrus.Logger, c *configurations.Collector, d *debugger.Debugger, path string) *APIServer {

	// create new instance of configurations interface
	cfg, err := configurations.New(configurations.CONF_TYPE_API_SERVER, path)
	if err != nil {
		d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_READ_CONFIG.Int(),
			configurations.ERROR_READ_CONFIG, err),
			logrus.ErrorLevel,
		)
		os.Exit(configurations.ERROR_READ_CONFIG.Int())
	}

	// Read config & return the requested struct type
	cf, err := cfg.Read()
	if err != nil {
		d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_READ_CONFIG.Int(),
			configurations.ERROR_READ_CONFIG, err),
			logrus.ErrorLevel,
		)
		os.Exit(configurations.ERROR_READ_CONFIG.Int())
	}

	configs := cf.(*configurations.APIServer)

	// open http access log file
	alfHasError := false
	accessLogfile, err := os.OpenFile(configs.AccessLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		alfHasError = true
		d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_CAN_T_OPEN_OR_CREATE_HTTP_ACCESS_LOG_FILE.Int(),
			configurations.ERROR_CAN_T_OPEN_OR_CREATE_HTTP_ACCESS_LOG_FILE, err),
			logrus.ErrorLevel,
		)
	}

	// open http error log file
	elfHasError := false
	errorLogfile, err := os.OpenFile(configs.ErrorLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		elfHasError = true
		d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_CAN_T_OPEN_OR_CREATE_HTTP_ERROR_LOG_FILE.Int(),
			configurations.ERROR_CAN_T_OPEN_OR_CREATE_HTTP_ERROR_LOG_FILE, err),
			logrus.ErrorLevel,
		)
	}

	// getIP2location conf
	// create new instance of configurations interface
	cfgLoc, err := configurations.New(configurations.CONF_TYPE_IP2LOCATION, path)
	if err != nil {
		d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_READ_CONFIG.Int(),
			configurations.ERROR_READ_CONFIG, err),
			logrus.ErrorLevel,
		)
		os.Exit(configurations.ERROR_READ_CONFIG.Int())
	}

	// Read config & return the requested strucut type
	cfl, err := cfgLoc.Read()
	if err != nil {
		d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_READ_CONFIG.Int(),
			configurations.ERROR_READ_CONFIG, err),
			logrus.ErrorLevel,
		)
		os.Exit(configurations.ERROR_READ_CONFIG.Int())
	}

	// make new instance of ip2location
	i2l := location.New(cfl.(*configurations.IP2Location), d)

	// collector will connect to pg-bouncer
	// and API server will connect to postgres db directly
	// so Exporter configs will be ommited and commented and
	// after that we will initialize our pgsql

	// var pgDB database.Postgres
	// for _, ex := range c.Exporter.Postgres {
	// 	// create new Postgres
	// 	pgDB = database.New(ex.Host, ex.User, ex.Password, ex.DB, c.IPReputation, ex.Port, d, i2l, 20, 50, 1*time.Hour)

	// 	// just get the very first one inside api server
	// 	break
	// }

	var pgDB database.Postgres
	for _, pg := range configs.Postgres {
		pgDB = database.New(pg.Host, pg.User, pg.Password, pg.DB, c.IPReputation, pg.Port, d, i2l, 20, 50, 1*time.Hour)

		// just get the very first one inside api server
		break
	}

	api := &APIServer{
		host:    configs.Listen.Address,
		port:    configs.Listen.Port,
		l:       l,
		c:       c,
		apiConf: configs,
		d:       d,

		pgdb: pgDB,
		db:   pgDB.GetDB(),

		ip2l: i2l,

		httpAccessHasError: alfHasError,
		httpErrorHasError:  elfHasError,
		httpAccessLog:      accessLogfile,
		httpErrorLog:       errorLogfile,

		ch:        make(chan os.Signal, 1),
		waitGroup: &sync.WaitGroup{},
	}

	api.ch = make(chan os.Signal, 1)
	signal.Notify(api.ch,
		// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGKILL, // "always fatal", "SIGKILL and SIGSTOP may not be caught by a program"
		syscall.SIGHUP,  // "terminal is disconnected"
	)

	// catch signals
	go func() {
		<-api.ch

		api.d.Verbose("Stopping API HTTP server...!", logrus.InfoLevel)

		// close http log files
		defer api.closeLogFiles()

		if err := api.httpSrv.Close(); err != nil {
			api.d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
				configurations.ERROR_CAN_T_STOP_API_HTTP_SERVER.Int(),
				configurations.ERROR_CAN_T_STOP_API_HTTP_SERVER, err),
				logrus.ErrorLevel,
			)
			os.Exit(configurations.ERROR_CAN_T_STOP_API_HTTP_SERVER.Int())
		} else {
			api.d.Verbose("API Server has stopped!", logrus.InfoLevel)
			os.Exit(0)
		}
	}()

	return api
}

// server api http server
// for providing more functionalities
func (api *APIServer) Serve() {

	// create new mux router
	r := mux.NewRouter()

	// create path prefix for all api routes
	apiRoutes := r.PathPrefix("/v1/api").Subrouter()

	// default / route
	apiRoutes.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("This is Netflow Collector API. For more information visit: https://openintelligence24.com"))
		return
	})

	// SOCKET_IO
	api.IP2LocationUpdate(r)
	defer api.apiSocketServer.Close() //close socket

	// routes for devices
	dr := apiRoutes.PathPrefix("/device").Subrouter()
	api.deviceRoutes(dr)

	// routes for ports
	pr := apiRoutes.PathPrefix("/port").Subrouter()
	api.portRoutes(pr)

	// routes for hosts
	hr := apiRoutes.PathPrefix("/host").Subrouter()
	api.hostRoutes(hr)

	// routes for protocols
	prr := apiRoutes.PathPrefix("/protocol").Subrouter()
	api.protoRoutes(prr)

	// routes for geos
	gr := apiRoutes.PathPrefix("/geo").Subrouter()
	api.geoRoutes(gr)

	// routes for ethernets
	er := apiRoutes.PathPrefix("/eth").Subrouter()
	api.ethernetRoutes(er)

	// routes for threats
	tr := apiRoutes.PathPrefix("/threat").Subrouter()
	api.threatRoutes(tr)

	// routes for flows
	flr := apiRoutes.PathPrefix("/flows").Subrouter()
	api.flowsRoutes(flr)

	// CORS definitions
	c := api.prepareCors()

	// start listening
	api.d.Verbose(fmt.Sprintf("API server is starting on %s:%d TLS=%v", api.apiConf.Listen.Address, api.apiConf.Listen.Port, api.apiConf.TLS.Enable), logrus.DebugLevel)
	var err error
	httpSrv := &http.Server{
		ReadTimeout:       time.Duration(api.apiConf.HTTP.ReadTimeOut) * time.Second,
		WriteTimeout:      time.Duration(api.apiConf.HTTP.WriteTimeOut) * time.Second,
		IdleTimeout:       time.Duration(api.apiConf.HTTP.IdleTimeOut) * time.Second,
		ReadHeaderTimeout: time.Duration(api.apiConf.HTTP.ReadHeaderTimeOut) * time.Second,

		Addr:    fmt.Sprintf("%s:%d", api.apiConf.Listen.Address, api.apiConf.Listen.Port),
		Handler: c.Handler(r),
	}

	// add an instance of http server to struct for future perpose
	api.httpSrv = httpSrv

	api.d.Verbose(fmt.Sprintf("API Server is listening on: '%v:%v' TLS:'%v'", api.apiConf.Listen.Address, api.apiConf.Listen.Port, api.apiConf.TLS.Enable), logrus.InfoLevel)

	if api.apiConf.TLS.Enable {
		// serve HTTPS
		err = api.httpSrv.ListenAndServeTLS(api.apiConf.TLS.Cert, api.apiConf.TLS.Key)
	} else {
		// serve HTTP
		err = api.httpSrv.ListenAndServe()
	}

	// can not listen & serve HTTP server
	if err != nil {
		api.d.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_CAN_T_LISTEN_AND_SERVE_HTTP_SERVER.Int(),
			configurations.ERROR_CAN_T_LISTEN_AND_SERVE_HTTP_SERVER, err),
			logrus.ErrorLevel,
		)
	}

}

// Close HTTP server
func (api *APIServer) Close() (err error) {
	return api.httpSrv.Close()
}

// close log files
func (api *APIServer) closeLogFiles() {
	// close log file after error
	if !api.httpAccessHasError {
		api.httpAccessLog.Close()
	}
	if !api.httpErrorHasError {
		api.httpErrorLog.Close()
	}
}

// prepare CORS
func (api *APIServer) prepareCors() *cors.Cors {

	if api.apiConf.HTTP.CORS.AllowAll {
		return cors.AllowAll()
	}

	var hdrs, mthds, orgs []string

	// prepare headers for CORS
	for _, h := range api.apiConf.HTTP.CORS.AllowedHeaders {
		hdrs = append(hdrs, h)
	}

	// prepare methods for CORS
	for _, m := range api.apiConf.HTTP.CORS.AllowedMethods {
		mthds = append(mthds, m)
	}

	// prepare origins for CORS
	for _, o := range api.apiConf.HTTP.CORS.AllowedOrigins {
		orgs = append(orgs, o)
	}

	// create new CORS definition
	return cors.New(
		cors.Options{
			AllowedHeaders: hdrs,
			AllowedOrigins: orgs,
			AllowedMethods: mthds,

			AllowCredentials:   true,
			OptionsPassthrough: true,
			Debug:              api.apiConf.Debug,
		},
	)

}
