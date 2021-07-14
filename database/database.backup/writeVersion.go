package database

import (
	"fmt"
	"strings"

	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/database/model"
	"github.com/sirupsen/logrus"
)

// write version into database if not exist yet
// otherwise it will update the last seen
func (p *Postgres) writeVersion(version uint) (versionID uint, err error) {
	var verModel model.Version

	// object exist in cache
	if v, err := p.getCached("version_" + fmt.Sprintf("%v", version)); err == nil {
		verModel = v.(model.Version)
		return verModel.ID, nil
	} else {
		p.db.Where("version = ?", version).First(&verModel)
	}

	if verModel.ID == 0 {
		// not found
		// need to be inserted to db
		verModel := model.Version{
			Version: version,
		}

		// insert to db
		result := p.db.Create(&verModel)

		// cache it
		p.cachedIt("version_"+fmt.Sprintf("%v", version), verModel)

		// check for error
		if result.Error != nil {

			// check if cache not prepared and not resolved
			if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
				return p.writeVersion(version)
			}

			p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
				configurations.ERROR_CAN_T_INSERT_VERSION_INFO.Int(),
				configurations.ERROR_CAN_T_INSERT_VERSION_INFO, result.Error),
				logrus.ErrorLevel,
			)
			return 0, result.Error
		}

		// cache it
		p.cachedIt("version_"+fmt.Sprintf("%v", version), verModel)

		return verModel.ID, nil
	} else {
		return verModel.ID, nil
	}
	// else {
	// 	// found and updated_at date/time must be updated
	// 	result := p.db.Model(&verModel).Update("updated_at", time.Now())

	// 	// check for error
	// 	// since we want to update just one
	// 	// field in the database (updated_at)
	// 	// we will continue with no error
	// 	// but logs must be generated to the checked to
	// 	// the log file for future investigations
	// 	if result.Error != nil {
	// 		p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
	// 			configurations.ERROR_CAN_T_UPDATE_VERSION_INFO.Int(),
	// 			configurations.ERROR_CAN_T_UPDATE_VERSION_INFO, result.Error),
	// 			logrus.ErrorLevel,
	// 		)
	// 	}
	// 	return verModel.ID, nil
	// }
}

// chaneg flow version to uint
func (p *Postgres) versionToUint(ver string) uint {
	switch ver {
	case "Netflow-V1":
		return 1
	case "Netflow-V5":
		return 5
	case "Netflow-V6":
		return 6
	case "Netflow-V7":
		return 7
	case "Netflow-V9":
		return 9
	case "IPFIX":
		return 10
	default:
		// return 0 for not specified version
		return 0
	}
}

// extract version
func (p *Postgres) _getVersion(s string) (uint, error) {
	verID, err := p.writeVersion(p.versionToUint(s))
	if err != nil {
		p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (version: %v)",
			configurations.ERROR_CAN_T_INSERT_METRICS_TO_POSTGRES_DB.Int(),
			configurations.ERROR_CAN_T_INSERT_METRICS_TO_POSTGRES_DB, err),
			logrus.ErrorLevel,
		)
	}
	return verID, err
}
