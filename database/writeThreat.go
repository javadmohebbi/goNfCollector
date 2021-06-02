package database

import (
	"fmt"

	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/database/model"
	"github.com/sirupsen/logrus"
)

// write threat into database if IP is part of
// any threat
func (p *Postgres) writeThreat(ip string, hostID uint) (threatID uint, hasThreat bool, err error) {

	// check for host reputation
	for _, rpu := range p.reputations {
		if resp := rpu.Get(ip); resp.Current > 0 {

			var threatModel model.Threat
			p.db.Where("source = ? AND kind = ? AND host_id = ?", rpu.GetType(), rpu.GetKind(), hostID).First(&threatModel)

			if threatModel.ID == 0 {
				// not found

				// need to be inserted to db
				threatModel := model.Threat{
					Source:        rpu.GetType(),
					Kind:          rpu.GetKind(),
					Acked:         false,
					Closed:        false,
					FalsePositive: false,

					HostID: hostID,
				}

				// insert to db
				result := p.db.Create(&threatModel)

				// check for error
				if result.Error != nil {
					p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
						configurations.ERROR_CAN_T_INSERT_THREAT_INFO.Int(),
						configurations.ERROR_CAN_T_INSERT_THREAT_INFO, result.Error),
						logrus.ErrorLevel,
					)
					return 0, true, result.Error
				}

				return threatModel.ID, true, nil
			} else {
				return threatModel.ID, true, nil
			}
			// else {
			// 	// found and updated_at date/time must be updated
			// 	result := p.db.Model(&threatModel).Update("updated_at", time.Now())

			// 	// check for error
			// 	// since we want to update just one
			// 	// field in the database (updated_at)
			// 	// we will continue with no error
			// 	// but logs must be generated to the checked to
			// 	// the log file for future investigations
			// 	if result.Error != nil {
			// 		p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			// 			configurations.ERROR_CAN_T_UPDATE_THREAT_INFO.Int(),
			// 			configurations.ERROR_CAN_T_UPDATE_THREAT_INFO, result.Error),
			// 			logrus.ErrorLevel,
			// 		)
			// 	}
			// 	return threatModel.ID, true, nil
			// }

		}
	}

	return 0, false, nil

}
