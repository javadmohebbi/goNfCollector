package database

import (
	"fmt"
	"time"

	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/database/model"
	"github.com/sirupsen/logrus"
)

// write device into database if not exist yet
// otherwise it will update the last seen
func (p *Postgres) writeDevice(device string) (deviceID uint, err error) {
	var dev model.Device
	p.db.Where("device = ?", device).First(&dev)

	if dev.ID == 0 {
		// not found
		// need to be inserted to db
		dev := model.Device{
			Device: device,
		}

		// insert to db
		result := p.db.Create(&dev)

		// check for error
		if result.Error != nil {
			p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
				configurations.ERROR_CAN_T_INSERT_DEVICE_INFO.Int(),
				configurations.ERROR_CAN_T_INSERT_DEVICE_INFO, err),
				logrus.ErrorLevel,
			)
			return 0, result.Error
		}

		return dev.ID, nil

	} else {
		// found and updated_at date/time must be updated
		result := p.db.Model(&dev).Update("updated_at", time.Now())

		// check for error
		// since we want to update just one
		// field in the database (updated_at)
		// we will continue with no error
		// but logs must be generated to the checked to
		// the log file for future investigations
		if result.Error != nil {
			p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
				configurations.ERROR_CAN_T_UPDATE_DEVICE_INFO.Int(),
				configurations.ERROR_CAN_T_UPDATE_DEVICE_INFO, err),
				logrus.ErrorLevel,
			)
		}
		return dev.ID, nil
	}
}
