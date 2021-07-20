package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/database/model"
	"github.com/sirupsen/logrus"
)

// write protocol into database if not exist yet
// otherwise it will update the last seen
func (p *Postgres) writeProtocol(protocol, protoName string) (protocolID uint, err error) {
	var protoModel model.Protocol

	// object exist in cache
	if v, err := p.getCached("proto_" + protoName); err == nil {
		protoModel = v.(model.Protocol)
		return protoModel.ID, nil
	} else {
		p.db.Where("protocol = ?", protocol).First(&protoModel)
	}

	if protoModel.ID == 0 {
		// not found
		// need to be inserted to db
		protoModel := model.Protocol{
			Protocol:     protocol,
			ProtocolName: protoName,
		}

		// insert to db
		result := p.db.Create(&protoModel)

		// cache it
		p.cachedIt("proto_"+protoName, protoModel)

		// check for error
		if result.Error != nil {

			// check if cache not prepared and not resolved
			if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
				return p.writeProtocol(protocol, protoName)
			}

			p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
				configurations.ERROR_CAN_T_INSERT_PROTOCOL_INFO.Int(),
				configurations.ERROR_CAN_T_INSERT_PROTOCOL_INFO, result.Error),
				logrus.ErrorLevel,
			)
			return 0, result.Error
		}

		return protoModel.ID, nil

		// } else {
		// 	return protoModel.ID, nil
		// }

	} else {
		// found and updated_at date/time must be updated
		result := p.db.Model(&protoModel).Update("updated_at", time.Now())

		// check for error
		// since we want to update just one
		// field in the database (updated_at)
		// we will continue with no error
		// but logs must be generated to the checked to
		// the log file for future investigations
		if result.Error != nil {
			p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
				configurations.ERROR_CAN_T_UPDATE_PROTOCOL_INFO.Int(),
				configurations.ERROR_CAN_T_UPDATE_PROTOCOL_INFO, result.Error),
				logrus.ErrorLevel,
			)
		}

		// cache it
		p.cachedIt("proto_"+protoName, protoModel)

		return protoModel.ID, nil
	}
}
