package database

import (
	"fmt"
	"net"
	"strings"

	"github.com/ammario/ipisp"
	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/database/model"
	"github.com/sirupsen/logrus"
)

// write autonomous into database if not exist yet
// otherwise it will update the last seen
func (p *Postgres) writeAutonomous(ip string) (asn string, autonomousID uint, err error) {

	// look for AS number
	asn = p.asnLookup(ip)

	var asnModel model.AutonomousSystem

	// object exist in cache
	if v, err := p.getCached("asn_" + ip); err == nil {
		asnModel = v.(model.AutonomousSystem)
		return asn, asnModel.ID, nil
	} else {
		p.db.Where("asn = ?", asn).First(&asnModel)
	}

	if asnModel.ID == 0 {
		// not found
		// need to be inserted to db
		asnModel := model.AutonomousSystem{
			ASN: asn,
		}

		// insert to db
		result := p.db.Create(&asnModel)

		// cache it
		p.cachedIt("asn_"+ip, asnModel)

		// check for error
		if result.Error != nil {

			// check if cache not prepared and not resolved
			if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
				return p.writeAutonomous(ip)
			}

			p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
				configurations.ERROR_CAN_T_INSERT_AUTONOMOUS_INFO.Int(),
				configurations.ERROR_CAN_T_INSERT_AUTONOMOUS_INFO, result.Error),
				logrus.ErrorLevel,
			)
			return asn, 0, result.Error
		}

		return asn, asnModel.ID, nil

	} else {
		return asn, asnModel.ID, nil
	}
	// else {
	// 	// found and updated_at date/time must be updated
	// 	result := p.db.Model(&asnModel).Update("updated_at", time.Now())

	// 	// check for error
	// 	// since we want to update just one
	// 	// field in the database (updated_at)
	// 	// we will continue with no error
	// 	// but logs must be generated to the checked to
	// 	// the log file for future investigations
	// 	if result.Error != nil {
	// 		p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
	// 			configurations	.ERROR_CAN_T_UPDATE_AUTONOMOUS_INFO.Int(),
	// 			configurations.ERROR_CAN_T_UPDATE_AUTONOMOUS_INFO, result.Error),
	// 			logrus.ErrorLevel,
	// 		)
	// 	}

	// 	// cache it
	// 	p.cachedIt("asn_"+ip, asnModel)

	// 	return asn, asnModel.ID, nil
	// }
}

// get asn lookup
func (p *Postgres) asnLookup(ip string) string {
	client, err := ipisp.NewDNSClient()

	// if client has error, returns NA
	if err != nil {
		return "NA"
	}

	// close ipisp client
	defer client.Close()

	resp, err := client.LookupIP(net.ParseIP(ip))
	// if no lookup, returns NA
	if err != nil {
		return "NA"
	}

	asnNumber := resp.ASN.String()

	return p.iplocation.GetAsnName(asnNumber) + "_" + asnNumber

}
