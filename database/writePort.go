package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/database/model"
	"github.com/sirupsen/logrus"
)

// write port into database if not exist yet
// otherwise it will update the last seen
func (p *Postgres) writePort(PortName, protoName, portNumber string, retry int) (portID uint, err error) {

	var portModel model.Port

	// object exist in cache
	if v, err := p.getCached("port_" + PortName); err == nil {
		portModel = v.(model.Port)
		// return portModel.ID, nil
	} else {
		p.db.Where("port_name = ?", PortName).First(&portModel)
	}

	if portModel.ID == 0 {
		// not found
		// need to be inserted to db
		portModel := model.Port{
			PortName:  PortName,
			PortProto: protoName + "/" + portNumber,
		}

		// insert to db
		result := p.db.Create(&portModel)

		// cache it
		p.cachedIt("port_"+PortName, portModel)

		// check for error
		if result.Error != nil {

			// color.Yellow.Printf("======= Retries: %v\n", retry)

			if retry < 3 {
				// check if cache not prepared and not resolved
				if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
					// color.Green.Printf("======= Retries (Duplicate key Error): %v\n", retry)
					return p.writePort(PortName, protoName, portNumber, retry+1)
				}

				// check if too mamy socket error
				if strings.Contains(result.Error.Error(), "socket: too many open files") {
					// color.Blue.Printf("======= Retries (Socket Error): %v\n", retry)
					return p.writePort(PortName, protoName, portNumber, retry+1)
				}

			}

			p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
				configurations.ERROR_CAN_T_INSERT_PORT_INFO.Int(),
				configurations.ERROR_CAN_T_INSERT_PORT_INFO, result.Error),
				logrus.ErrorLevel,
			)
			return 0, result.Error
		}

		return portModel.ID, nil

		// } else {
		// 	return portModel.ID, nil
		// }
	} else {
		// found and updated_at date/time must be updated
		result := p.db.Model(&portModel).Update("updated_at", time.Now())

		// check for error
		// since we want to update just one
		// field in the database (updated_at)
		// we will continue with no error
		// but logs must be generated to the checked to
		// the log file for future investigations
		if result.Error != nil {
			p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
				configurations.ERROR_CAN_T_UPDATE_PORT_INFO.Int(),
				configurations.ERROR_CAN_T_UPDATE_PORT_INFO, result.Error),
				logrus.ErrorLevel,
			)
		}

		// cache it
		p.cachedIt("port_"+PortName, portModel)

		return portModel.ID, nil
	}
}
