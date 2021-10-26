package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/goNfCollector/common"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// this function will get information about the ethernets
// on a device based on an intervals like 1m, 2h ...
func (api *APIServer) getEthernetsByDeviceByInterval(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	type Result struct {
		Time time.Time `json:"_time" gorm:"column:_time"`

		Eth        uint   `json:"eth" gorm:"column:eth"`
		EthUniqKey string `json:"eth_key" gorm:"column:eth_key"`
		EthName    string `json:"eth_name" gorm:"column:eth_name"`

		TotalBytes   uint `json:"total_bytes"`
		TotalPackets uint `json:"total_packets"`
	}
	var resultIn, resultOut []Result

	// // summary
	// type ResultIn struct {
	// 	Time time.Time `json:"_time" gorm:"column:_time"`

	// 	EthIn        uint   `json:"eth" gorm:"column:eth_in"`
	// 	EthInUniqKey string `json:"eth_key" gorm:"column:eth_in_key"`
	// 	EthInName    string `json:"eth_name" gorm:"column:eth_in_name"`

	// 	TotalBytes   uint `json:"total_bytes"`
	// 	TotalPackets uint `json:"total_packets"`
	// }
	// var resultIn []ResultIn

	// // summary
	// type ResultOut struct {
	// 	Time time.Time `json:"_time" gorm:"column:_time"`

	// 	EthOut        uint   `json:"eth" gorm:"column:eth_out"`
	// 	EthOutUniqKey string `json:"eth_key" gorm:"column:eth_out_key"`
	// 	EthOutName    string `json:"eth_name" gorm:"column:eth_out_name"`

	// 	TotalBytes   uint `json:"total_bytes"`
	// 	TotalPackets uint `json:"total_packets"`
	// }
	// var resultOut []ResultOut

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// extract interval
	interval := common.GetPGInterval(mux.Vars(r)["interval"])

	deviceID := mux.Vars(r)["deviceID"]
	if deviceID == "" {
		deviceID = "1"
	}

	timeTruncate := common.GetPGGroupByInterval(mux.Vars(r)["interval"])

	type ethCountStruct struct {
		UniqKey string `gorm:"uniq_key" json:"uniq_key"`
		ID      uint   `gorm:"id" json:"id"`
	}

	type TheResult struct {
		Ingres  bool     `json:"ingres"`
		Outgres bool     `json:"outgres"`
		Data    []Result `json:"data"`
	}
	var theResult []TheResult

	var ethCount []ethCountStruct
	strEthernetCounts := `SELECT uniq_key, id FROM ethernets WHERE device_id = ` + deviceID + ``
	api.db.Raw(strEthernetCounts).Scan(&ethCount)

	// type ethResult struct {
	// 	EthCount []ethCountStruct `json:"keys"`
	// 	In       []Re
	// 	Out      []ResultOut
	// }
	// ethResults := ethResult{
	// 	EthCount: ethCount,
	// }

	for _, eth := range ethCount {
		strQuery := `
		SELECT
			date_trunc('` + timeTruncate + `', flows.created_at) as "_time",
			flows.in_ethernet_id as eth,
			ethernets.uniq_key as eth_key,

			sum(flows.byte) as total_bytes,
			sum(flows.packet) as total_packets
		FROM
			flows
		JOIN
			ethernets
		ON
			flows.in_ethernet_id = ethernets.id
		WHERE
			flows.created_at > NOW() - INTERVAL '` + interval + `'
			AND flows.device_id = ` + deviceID + `
			 AND ethernets.uniq_key = '` + eth.UniqKey + `'
			 AND flows.in_ethernet_id = ` + fmt.Sprintf("%v", eth.ID) + `
		GROUP BY
			_time, flows.in_ethernet_id, ethernets.uniq_key, ethernets.name
		ORDER BY
			_time
	;`
		// fmt.Println(strQuery)
		api.db.Raw(strQuery).Scan(&resultIn)

		// out
		strQuery = `
		SELECT
			date_trunc('` + timeTruncate + `', flows.created_at) as "_time",
			flows.out_ethernet_id as eth,
			ethernets.uniq_key as eth_key,

			sum(flows.byte) as total_bytes,
			sum(flows.packet) as total_packets
		FROM
			flows
		JOIN
			ethernets
		ON
			flows.out_ethernet_id = ethernets.id
		WHERE
			flows.created_at > NOW() - INTERVAL '` + interval + `'
			AND flows.device_id = ` + deviceID + `
			 AND ethernets.uniq_key = '` + eth.UniqKey + `'
			 AND flows.out_ethernet_id = ` + fmt.Sprintf("%v", eth.ID) + `
		GROUP BY
			_time, flows.out_ethernet_id, ethernets.uniq_key, ethernets.name
		ORDER BY
			_time
	;`
		// fmt.Println(strQuery)

		api.db.Raw(strQuery).Scan(&resultOut)

		theResult = append(theResult, TheResult{
			Ingres:  true,
			Outgres: false,
			Data:    resultIn,
		})

		theResult = append(theResult, TheResult{
			Ingres:  false,
			Outgres: true,
			Data:    resultOut,
		})

	}

	_ = json.NewEncoder(w).Encode(&theResult)
	return

}
