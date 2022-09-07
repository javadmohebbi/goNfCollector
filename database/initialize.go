package database

import "github.com/goNfCollector/database/model"

// autoMigrate database
func (p *Postgres) autoMigrate() error {

	return p.db.AutoMigrate(
		&model.AutonomousSystem{},
		&model.Device{},
		&model.Domain{},
		&model.Flag{},
		&model.Geo{},
		&model.Host{},
		&model.Port{},
		&model.Protocol{},
		&model.Threat{},
		&model.Version{},

		&model.Flow{},

		&model.Ethernet{},

		// ip2l settings
		&model.Settings{},
	)

}
