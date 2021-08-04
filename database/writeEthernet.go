package database

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/database/model"
	"github.com/sirupsen/logrus"
)

// write ethernet into database if not exist yet
// otherwise it will update the last seen
func (p *Postgres) writeEthernet(eth string, device string, deviceID uint) (ethernetID int, err error) {
	var ethModel model.Ethernet

	et, err := strconv.Atoi(eth)
	if err != nil {
		p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
			configurations.ERROR_CAN_T_CONVERT_ETH_INDEX.Int(),
			configurations.ERROR_CAN_T_CONVERT_ETH_INDEX, err),
			logrus.ErrorLevel,
		)
		return 0, err
	}

	// object exist in cache
	if v, err := p.getCached("ethernet_" + device + "_" + eth); err == nil {
		ethModel = v.(model.Ethernet)
		// return int(ethModel.ID), nil
	} else {
		p.db.Where("uniq_key = ?", fmt.Sprintf("%v:%v", device, uint(et))).First(&ethModel)
	}

	if ethModel.ID == 0 {
		// not found
		// need to be inserted to db
		ethModel := model.Ethernet{
			Ethernet: uint(et),
			DeviceID: deviceID,
			UniqKey:  fmt.Sprintf("%v:%v", device, et),
		}

		// insert to db
		result := p.db.Create(&ethModel)

		// cache it
		p.cachedIt("ethernet_"+device+"_"+eth, ethModel)

		// check for error
		if result.Error != nil {

			// check if cache not prepared and not resolved
			if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
				return p.writeEthernet(eth, device, deviceID)
			}

			p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
				configurations.ERROR_CAN_T_INSERT_ETH_INFO.Int(),
				configurations.ERROR_CAN_T_INSERT_ETH_INFO, result.Error),
				logrus.ErrorLevel,
			)
			return int(-1), result.Error
		}

		return int(ethModel.ID), nil

	} else {
		return int(ethModel.ID), nil
	}
}
