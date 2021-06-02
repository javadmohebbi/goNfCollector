package database

import (
	"fmt"
	"strconv"
	"time"

	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/database/model"
	"github.com/sirupsen/logrus"
)

// write flag into database if not exist yet
// otherwise it will update the last seen
func (p *Postgres) writeFlag(flgs string) (flagID uint, fin, syn, rst, psh, ack, urg, ece, cwr bool, err error) {

	humanReadableFlags, fin, syn, rst, psh, ack, urg, ece, cwr := p._tcpFlags(flgs)

	var flagModel model.Flag
	p.db.Where("flags = ?", flgs).First(&flagModel)

	if flagModel.ID == 0 {
		// not found
		// need to be inserted to db
		flagModel := model.Flag{
			Flags: flgs,
			Info:  humanReadableFlags,
		}

		// insert to db
		result := p.db.Create(&flagModel)

		// check for error
		if result.Error != nil {
			p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
				configurations.ERROR_CAN_T_INSERT_FLAG_INFO.Int(),
				configurations.ERROR_CAN_T_INSERT_FLAG_INFO, err),
				logrus.ErrorLevel,
			)
			return 0, fin, syn, rst, psh, ack, urg, ece, cwr, result.Error
		}

		return flagModel.ID, fin, syn, rst, psh, ack, urg, ece, cwr, nil

	} else {
		// found and updated_at date/time must be updated
		result := p.db.Model(&flagModel).Update("updated_at", time.Now())

		// check for error
		// since we want to update just one
		// field in the database (updated_at)
		// we will continue with no error
		// but logs must be generated to the checked to
		// the log file for future investigations
		if result.Error != nil {
			p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
				configurations.ERROR_CAN_T_UPDATE_FLAG_INFO.Int(),
				configurations.ERROR_CAN_T_UPDATE_FLAG_INFO, err),
				logrus.ErrorLevel,
			)
		}
		return flagModel.ID, fin, syn, rst, psh, ack, urg, ece, cwr, nil
	}
}

// get flags and return human readable string
// with comma delimited
func (p *Postgres) _tcpFlags(f string) (str string, fin, syn, rst, psh, ack, urg, ece, cwr bool) {

	FIN := 0x01
	SYN := 0x02
	RST := 0x04
	PSH := 0x08
	ACK := 0x10
	URG := 0x20
	ECE := 0x40
	CWR := 0x80

	i, _ := strconv.Atoi(f)

	// hx := fmt.Sprintf("%x", i)

	if FIN&i == 0x01 {
		str += " FIN "
		fin = true
	}
	if SYN&i == 0x02 {
		str += " SYN "
		syn = true
	}
	if RST&i == 0x04 {
		str += " RST "
		rst = true
	}
	if PSH&i == 0x08 {
		str += " PSH "
		psh = true
	}
	if ACK&i == 0x10 {
		str += " ACK "
		ack = true
	}
	if URG&i == 0x20 {
		str += " URG "
		urg = true
	}
	if ECE&i == 0x40 {
		str += " ECE "
		ece = true
	}
	if CWR&i == 0x80 {
		str += " CWR "
		cwr = true
	}

	return str, fin, syn, rst, psh, ack, urg, ece, cwr

}
