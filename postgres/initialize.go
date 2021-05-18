package postgres

import (
	"context"
	"fmt"

	"github.com/goNfCollector/debugger"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

var (
	tb_netflow_devices_metrics = "netflow_devices_metrics"
	tb_netflow_version         = "netflow_versions_metrics"
	tb_netflow_host            = "netflow_hosts_metrics"
	tb_netflow_port            = "netflow_ports_metrics"
	tb_netflow_protocol        = "netflow_protocols_metrics"
	tb_netflow_domain          = "netflow_domains_metrics"
	tb_netflow_country         = "netflow_countries_metrics"
	tb_netflow_state           = "netflow_states_metrics"
	tb_netflow_city            = "netflow_cities_metrics"
	tb_netflow_flag            = "netflow_flags_metrics"
	tb_netflow_asn             = "netflow_ans_metrics"
	tb_netflow_threat          = "netflow_threats_metrics"
	tb_netflow_details         = "netflow_details_metrics"
)

// initialize postgres database
func initializeDatabase(ctx context.Context, pool *pgxpool.Pool, d *debugger.Debugger) error {

	var err error

	d.Verbose(fmt.Sprintf("initializing table %s", tb_netflow_devices_metrics), logrus.DebugLevel)
	// initialize device metric table
	err = initTBDevices(ctx, pool)
	if err != nil {
		return err
	}

	d.Verbose(fmt.Sprintf("initializing table %s", tb_netflow_version), logrus.DebugLevel)
	err = initTBVersion(ctx, pool)
	if err != nil {
		return err
	}

	d.Verbose(fmt.Sprintf("initializing table %s", tb_netflow_host), logrus.DebugLevel)
	err = initTBHost(ctx, pool)
	if err != nil {
		return err
	}

	d.Verbose(fmt.Sprintf("initializing table %s", tb_netflow_port), logrus.DebugLevel)
	err = initTBPort(ctx, pool)
	if err != nil {
		return err
	}

	d.Verbose(fmt.Sprintf("initializing table %s", tb_netflow_protocol), logrus.DebugLevel)
	err = initTBProtocol(ctx, pool)
	if err != nil {
		return err
	}

	d.Verbose(fmt.Sprintf("initializing table %s", tb_netflow_domain), logrus.DebugLevel)
	err = initTBDomain(ctx, pool)
	if err != nil {
		return err
	}

	d.Verbose(fmt.Sprintf("initializing table %s", tb_netflow_country), logrus.DebugLevel)
	err = initTBCountry(ctx, pool)
	if err != nil {
		return err
	}

	d.Verbose(fmt.Sprintf("initializing table %s", tb_netflow_state), logrus.DebugLevel)
	err = initTBState(ctx, pool)
	if err != nil {
		return err
	}

	d.Verbose(fmt.Sprintf("initializing table %s", tb_netflow_city), logrus.DebugLevel)
	err = initTBCity(ctx, pool)
	if err != nil {
		return err
	}

	d.Verbose(fmt.Sprintf("initializing table %s", tb_netflow_flag), logrus.DebugLevel)
	err = initTBFlag(ctx, pool)
	if err != nil {
		return err
	}

	d.Verbose(fmt.Sprintf("initializing table %s", tb_netflow_asn), logrus.DebugLevel)
	err = initTBASN(ctx, pool)
	if err != nil {
		return err
	}

	d.Verbose(fmt.Sprintf("initializing table %s", tb_netflow_threat), logrus.DebugLevel)
	err = initTBThreat(ctx, pool)
	if err != nil {
		return err
	}

	d.Verbose(fmt.Sprintf("initializing table %s", tb_netflow_details), logrus.DebugLevel)
	err = initTBDetail(ctx, pool)
	if err != nil {
		return err
	}

	return nil
}

// initialize detail table
func initTBDetail(ctx context.Context, pool *pgxpool.Pool) error {
	q := `
		CREATE TABLE IF NOT EXISTS ` + tb_netflow_details + ` (
			detail_id serial PRIMARY KEY,

			device_id INT NOT NULL,

			flow_version TEXT NOT NULL,

			protocol_id INT NOT NULL,

			src_asn_id INT NOT NULL,
			src_host_id INT NOT NULL,
			src_port_id INT NOT NULL,
			src_country_id INT NOT NULL,
			src_state_id INT NOT NULL,
			src_city_id INT NOT NULL,


			dst_asn_id INT NOT NULL,
			dst_host_id INT NOT NULL,
			dst_port_id INT NOT NULL,
			dst_country_id INT NOT NULL,
			dst_state_id INT NOT NULL,
			dst_city_id INT NOT NULL,

			next_hop TEXT NOT NULL,

			flag_id INT NOT NULL,


			time TIMESTAMPTZ NOT NULL,


			FOREIGN KEY (device_id)
			REFERENCES ` + tb_netflow_devices_metrics + ` (device_id),

			FOREIGN KEY (protocol_id)
			REFERENCES ` + tb_netflow_protocol + ` (protocol_id),

			FOREIGN KEY (src_asn_id)
			REFERENCES ` + tb_netflow_asn + ` (asn_id),

			FOREIGN KEY (src_host_id)
			REFERENCES ` + tb_netflow_host + ` (host_id),

			FOREIGN KEY (src_port_id)
			REFERENCES ` + tb_netflow_port + ` (port_id),

			FOREIGN KEY (src_country_id)
			REFERENCES ` + tb_netflow_country + ` (country_id),

			FOREIGN KEY (src_state_id)
			REFERENCES ` + tb_netflow_state + ` (state_id),

			FOREIGN KEY (src_city_id)
			REFERENCES ` + tb_netflow_city + ` (city_id),



			FOREIGN KEY (dst_asn_id)
			REFERENCES ` + tb_netflow_asn + ` (asn_id),

			FOREIGN KEY (dst_host_id)
			REFERENCES ` + tb_netflow_host + ` (host_id),

			FOREIGN KEY (dst_port_id)
			REFERENCES ` + tb_netflow_port + ` (port_id),

			FOREIGN KEY (dst_country_id)
			REFERENCES ` + tb_netflow_country + ` (country_id),

			FOREIGN KEY (dst_state_id)
			REFERENCES ` + tb_netflow_state + ` (state_id),

			FOREIGN KEY (dst_city_id)
			REFERENCES ` + tb_netflow_city + ` (city_id),

			FOREIGN KEY (flag_id)
			REFERENCES ` + tb_netflow_flag + ` (flag_id)

		);
	`
	_, err := pool.Exec(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

// initialize threat table
func initTBThreat(ctx context.Context, pool *pgxpool.Pool) error {
	q := `
		CREATE TABLE IF NOT EXISTS ` + tb_netflow_threat + ` (
			threat_id serial PRIMARY KEY,

			first_seen TIMESTAMPTZ NOT NULL,
			last_seen TIMESTAMPTZ NOT NULL,

			source TEXT NOT NULL,
			kind TEXT NOT NULL,

			host_id INT NOT NULL,

			FOREIGN KEY (host_id)
      		REFERENCES ` + tb_netflow_host + ` (host_id)
		);
	`
	_, err := pool.Exec(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

// initialize asn table
func initTBASN(ctx context.Context, pool *pgxpool.Pool) error {
	q := `
		CREATE TABLE IF NOT EXISTS ` + tb_netflow_asn + ` (
			asn_id serial PRIMARY KEY,

			first_seen TIMESTAMPTZ NOT NULL,
			last_seen TIMESTAMPTZ NOT NULL,

			asn TEXT NOT NULL
		);
	`
	_, err := pool.Exec(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

// initialize flag table
func initTBFlag(ctx context.Context, pool *pgxpool.Pool) error {
	q := `
		CREATE TABLE IF NOT EXISTS ` + tb_netflow_flag + ` (
			flag_id serial PRIMARY KEY,

			first_seen TIMESTAMPTZ NOT NULL,
			last_seen TIMESTAMPTZ NOT NULL,

			flag TEXT NOT NULL
		);
	`
	_, err := pool.Exec(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

// initialize city table
func initTBCity(ctx context.Context, pool *pgxpool.Pool) error {
	q := `
		CREATE TABLE IF NOT EXISTS ` + tb_netflow_city + ` (
			city_id serial PRIMARY KEY,

			first_seen TIMESTAMPTZ NOT NULL,
			last_seen TIMESTAMPTZ NOT NULL,

			city TEXT NOT NULL
		);
	`
	_, err := pool.Exec(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

// initialize state table
func initTBState(ctx context.Context, pool *pgxpool.Pool) error {
	q := `
		CREATE TABLE IF NOT EXISTS ` + tb_netflow_state + ` (
			state_id serial PRIMARY KEY,

			first_seen TIMESTAMPTZ NOT NULL,
			last_seen TIMESTAMPTZ NOT NULL,

			state TEXT NOT NULL
		);
	`
	_, err := pool.Exec(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

// initialize country table
func initTBCountry(ctx context.Context, pool *pgxpool.Pool) error {
	q := `
		CREATE TABLE IF NOT EXISTS ` + tb_netflow_country + ` (
			country_id serial PRIMARY KEY,

			first_seen TIMESTAMPTZ NOT NULL,
			last_seen TIMESTAMPTZ NOT NULL,

			country_long TEXT NOT NULL,
			country_short TEXT NOT NULL
		);
	`
	_, err := pool.Exec(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

// initialize domain table
func initTBDomain(ctx context.Context, pool *pgxpool.Pool) error {
	q := `
		CREATE TABLE IF NOT EXISTS ` + tb_netflow_domain + ` (
			domain_id serial PRIMARY KEY,

			first_seen TIMESTAMPTZ NOT NULL,
			last_seen TIMESTAMPTZ NOT NULL,

			domain TEXT NOT NULL
		);
	`
	_, err := pool.Exec(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

// initialize protocol table
func initTBProtocol(ctx context.Context, pool *pgxpool.Pool) error {
	q := `
		CREATE TABLE IF NOT EXISTS ` + tb_netflow_protocol + ` (
			protocol_id serial PRIMARY KEY,

			first_seen TIMESTAMPTZ NOT NULL,
			last_seen TIMESTAMPTZ NOT NULL,

			proto TEXT NOT NULL,
			proto_name TEXT NOT NULL
		);
	`
	_, err := pool.Exec(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

// initialize port table
func initTBPort(ctx context.Context, pool *pgxpool.Pool) error {
	q := `
		CREATE TABLE IF NOT EXISTS ` + tb_netflow_port + ` (
			port_id serial PRIMARY KEY,

			first_seen TIMESTAMPTZ NOT NULL,
			last_seen TIMESTAMPTZ NOT NULL,

			port TEXT NOT NULL,
			port_name TEXT NOT NULL
		);
	`
	_, err := pool.Exec(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

// initialize host table
func initTBHost(ctx context.Context, pool *pgxpool.Pool) error {
	q := `
		CREATE TABLE IF NOT EXISTS ` + tb_netflow_host + ` (
			host_id serial PRIMARY KEY,

			first_seen TIMESTAMPTZ NOT NULL,
			last_seen TIMESTAMPTZ NOT NULL,

			host TEXT NOT NULL
		);
	`
	_, err := pool.Exec(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

// initialize src table
func initTBVersion(ctx context.Context, pool *pgxpool.Pool) error {
	q := `
		CREATE TABLE IF NOT EXISTS ` + tb_netflow_version + ` (
			version_id serial PRIMARY KEY,
			time TIMESTAMPTZ NOT NULL,

			version TEXT NOT NULL
		);
	`
	_, err := pool.Exec(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

// initialize device table
func initTBDevices(ctx context.Context, pool *pgxpool.Pool) error {
	q := `
		CREATE TABLE IF NOT EXISTS ` + tb_netflow_devices_metrics + ` (
			device_id serial PRIMARY KEY,

			first_seen TIMESTAMPTZ NOT NULL,
			last_seen TIMESTAMPTZ NOT NULL,

			device_ip TEXT NOT NULL,

			device_name TEXT NULL,
			device_info TEXT NULL

		);

	`
	// TSDB timescale
	// q := `
	// 	CREATE TABLE IF NOT EXISTS ` + tb_netflow_devices_metrics + ` (
	// 		time TIMESTAMPTZ NOT NULL,
	// 		device_ip TEXT NOT NULL,
	// 		version TEXT NOT NULL,
	// 		byte DOUBLE PRECISION,
	// 		packet DOUBLE PRECISION
	// 	);
	// 	SELECT create_hypertable('` + tb_netflow_devices_metrics + `', 'time',chunk_time_interval => 86400000000, if_not_exists => TRUE);
	// `
	_, err := pool.Exec(ctx, q)
	if err != nil {
		return err
	}

	return nil
}
