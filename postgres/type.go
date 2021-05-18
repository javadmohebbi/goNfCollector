package postgres

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/debugger"
	"github.com/goNfCollector/location"
	"github.com/goNfCollector/reputation"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

// postgres database struct to
// access to postgres DB
type Postgres struct {
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

	// debugger
	Debuuger *debugger.Debugger

	// postgres context
	ctx context.Context

	// postgres pool
	pool *pgxpool.Pool

	// reputation
	reputations []reputation.Reputation

	// IP2locaion instance
	iplocation *location.IPLocation

	// channel
	ch chan os.Signal

	WaitGroup *sync.WaitGroup
}

// return exporter info
func (p Postgres) String() string {
	return fmt.Sprintf("Postgress %s:%d user:%s db:%s", p.Host, p.Port, p.User, p.DB)
}

// create new instance of influxDB
func New(host, user, pass, db string, ipReputationConf configurations.IpReputation, port int, d *debugger.Debugger, ip2location *location.IPLocation) Postgres {

	ctx := context.Background()

	d.Verbose(fmt.Sprintf("connecting to postgres db '%s' on %s:%v using username '%s' ", db, host, port, user), logrus.DebugLevel)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		user, pass, host, port, db,
	)

	conn, err := pgxpool.Connect(ctx, connStr)

	if err != nil {
		// fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		// os.Exit(1)
		d.Verbose(fmt.Sprintf("can not connect to postgres db '%s' on %s:%v using username '%s' due to error: %s", db, host, port, user, err.Error()), logrus.ErrorLevel)
		os.Exit(configurations.ERROR_CAN_T_CONNECT_TO_POSTGRES_DB.Int())
	}

	d.Verbose(fmt.Sprintf("initializing postgres db '%s' on %s:%v using username '%s'", db, host, port, user), logrus.DebugLevel)
	err = initializeDatabase(ctx, conn, d)
	if err != nil {
		d.Verbose(fmt.Sprintf("can not initialize postgres db '%s' on %s:%v using username '%s' due to error: %s", db, host, port, user, err.Error()), logrus.ErrorLevel)
		os.Exit(configurations.ERROR_CAN_T_INIT_POSTGRES_DB.Int())
	}

	// add reputation kind to reputation array
	var reputs []reputation.Reputation
	rptIpSum, err := reputation.NewIPSum(ipReputationConf.IPSumPath)
	if err == nil {
		reput, err := reputation.New(rptIpSum, d)

		if err == nil {
			reputs = append(reputs, *reput)
		}
	}

	d.Verbose(fmt.Sprintf("new postgres exporter %v:%v db:%v user:%v is created", host, port, db, user), logrus.DebugLevel)

	// retun influxDB
	return Postgres{
		Host:     host,
		Port:     port,
		User:     user,
		Password: pass,
		DB:       db,
		Debuuger: d,
		ctx:      ctx,
		pool:     conn,

		iplocation: ip2location,

		reputations: reputs,

		WaitGroup: &sync.WaitGroup{},
	}

}
