package database

import (
	"fmt"
	"strings"

	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/database/model"
	"github.com/sirupsen/logrus"
)

// write host into database if not exist yet
// otherwise it will update the last seen
func (p *Postgres) writeHost(host string) (hostID uint, err error) {
	var hostModel model.Host

	// object exist in cache
	if v, err := p.getCached("host_" + host); err == nil {
		hostModel = v.(model.Host)
		return hostModel.ID, nil
	} else {
		p.db.Where("host = ?", host).First(&hostModel)
	}

	if hostModel.ID == 0 {
		// not found
		// need to be inserted to db
		hostModel := model.Host{
			Host: host,
		}

		// insert to db
		result := p.db.Create(&hostModel)

		// cache it
		p.cachedIt("host_"+host, hostModel)

		// check for error
		if result.Error != nil {

			// check if cache not prepared and not resolved
			if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
				return p.writeHost(host)
			}

			p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
				configurations.ERROR_CAN_T_INSERT_HOST_INFO.Int(),
				configurations.ERROR_CAN_T_INSERT_HOST_INFO, result.Error),
				logrus.ErrorLevel,
			)
			return 0, result.Error
		}

		return hostModel.ID, nil

	} else {
		return hostModel.ID, nil
	}
	// else {
	// 	// found and updated_at date/time must be updated
	// 	result := p.db.Model(&hostModel).Update("updated_at", time.Now())

	// 	// check for error
	// 	// since we want to update just one
	// 	// field in the database (updated_at)
	// 	// we will continue with no error
	// 	// but logs must be generated to the checked to
	// 	// the log file for future investigations
	// 	if result.Error != nil {
	// 		p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
	// 			configurations.ERROR_CAN_T_UPDATE_HOST_INFO.Int(),
	// 			configurations.ERROR_CAN_T_UPDATE_HOST_INFO, result.Error),
	// 			logrus.ErrorLevel,
	// 		)
	// 	}

	// 	// cache it
	// 	p.cachedIt("host_"+host, hostModel)

	// 	return hostModel.ID, nil
	// }
}
