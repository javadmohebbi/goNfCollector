package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/goNfCollector/common"
	"github.com/goNfCollector/database/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// get all devices
func (api *APIServer) getAllDevices(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Extract query string for pagination
	pg := common.ExtractPaginationQueryString(r.URL.Query())

	fields := []string{"devices.device", "devices.name", "devices.info"}
	strWhere := common.PaginationStrWhereBuilder(pg["Filter"].(string), fields)

	// strQuery := `
	// 	SELECT id as device_id, device, name as device_name, info as device_info,
	// 		(select created_at from flows where device_id = devices.id order by created_at desc limit 1 ) as last_activity,
	// 		(select version from versions join flows on flows.version_id = versions.id where flows.device_id = versions.id order by flows.created_at desc limit 1) as flow_version
	// 	FROM
	// 		devices
	// `

	strQuery := `
		SELECT id as device_id, device, name as device_name, info as device_info,
			(select created_at from flows where device_id = devices.id order by created_at desc limit 1 ) as last_activity
		FROM
			devices
	`

	if strWhere != "" {
		strQuery += ` WHERE ` + strWhere + `

		`
	}

	if pg["Order"].(string) != "" {
		strQuery += ` ORDER BY ` + pg["Order"].(string) + ` ` + pg["OrderType"].(string)
	}

	if !pg["NoPagination"].(bool) {
		strQuery += ` LIMIT ` + fmt.Sprintf("%v", pg["PerPage"].(int)) + `
			OFFSET ` + fmt.Sprintf("%v", pg["PerPage"].(int)*(pg["Page"].(int)-1)) + ` ;`
	} else {
		strQuery += `  ;`
	}

	type FilteredResult struct {
		DeviceID     uint      `json:"device_id" gorm:"device_id"`
		Device       string    `json:"device" gorm:"device"`
		DeviceName   string    `json:"device_name" gorm:"device_name"`
		DeviceInfo   string    `json:"device_info" gorm:"device_info"`
		LastActivity time.Time `json:"last_activity" gorm:"last_activity"`
		FlowVersion  uint      `json:"flow_version" gorm:"flow_version"`
	}
	var filteredResult []*FilteredResult
	api.db.Raw(strQuery).Scan(&filteredResult)
	for _, fr := range filteredResult {
		type tmpVer struct {
			Version uint `json:"version" gorm:"version"`
		}
		strSubQuery := `SELECT
							version
						FROM
							versions
						JOIN flows
							ON versions.id = flows.version_id
						WHERE
							flows.device_id = ` + fmt.Sprintf("%v", fr.DeviceID) + `
						ORDER BY flows.created_at
						LIMIT 1
		`
		var tv tmpVer
		api.db.Raw(strSubQuery).Scan(&tv)
		fr.FlowVersion = tv.Version
	}

	type CountStruct struct {
		Count uint `json:"count" gorm:"count"`
	}
	var total CountStruct

	api.db.Raw("SELECT count(*) as count from devices WHERE " + strWhere + ";").Scan(&total)

	if len(filteredResult) == 0 {
		// not found
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "There are no records in the database regarding your request",
		})
		return
	}

	if pg["NoPagination"].(bool) {
		_ = json.NewEncoder(w).Encode(&filteredResult)
		return
	} else {
		type paginateResult struct {
			Pagination model.PaginationModel
			Result     []*FilteredResult
		}
		result := paginateResult{
			Pagination: model.PaginationModel{
				Page:    pg["Page"].(int),
				PerPage: pg["PerPage"].(int),
				Total:   int(total.Count),
			},
			Result: filteredResult,
		}
		_ = json.NewEncoder(w).Encode(result)
		return

	}

	// model
	// var dvs []GetAllDevice
	// api.db.Model(&model.Device{}).Find(&dvs)

	// return
}

// this function will get information about the last summary of
// devices base on an intervals like 1m, 2h ...
func (api *APIServer) getAllDevicesSummaryByInterval(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// summary
	type Result struct {
		DeviceID uint   `json:"device_id"`
		Device   string `json:"device"`

		FlowCount uint `json:"flow_count"`

		DeviceName string `json:"device_name"`

		DeviceInfo string `json:"device_info"`

		TotalBytes   uint `json:"total_bytes"`
		TotalPackets uint `json:"total_packets"`
	}
	var result []Result

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// extract interval
	interval := common.GetPGInterval(mux.Vars(r)["interval"])

	api.db.Raw(`
		SELECT
			devices.id as device_id,
			devices.device,
			count(*) as flow_count,
			devices.name as device_name,
			devices.info as device_info,
			sum(flows.byte) as total_bytes,
			sum(flows.packet) as total_packets
		FROM
			flows
		JOIN
			devices
		ON
			flows.device_id = devices.id
		WHERE
			flows.created_at > NOW() - INTERVAL '` + interval + `'
		GROUP BY
			devices.id, devices.device, device_name, device_info;
	`).Scan(&result)

	_ = json.NewEncoder(w).Encode(&result)
	return

}

// only for on device
// this function will get information about the last summary of
// devices base on an intervals like 1m, 2h ...
func (api *APIServer) getAllDevicesSummaryByIntervalByDevice(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// summary
	type Result struct {
		DeviceID uint   `json:"device_id"`
		Device   string `json:"device"`

		FlowCount uint `json:"flow_count"`

		DeviceName string `json:"device_name"`

		DeviceInfo string `json:"device_info"`

		TotalBytes   uint `json:"total_bytes"`
		TotalPackets uint `json:"total_packets"`
	}
	var result []Result

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// extract interval
	interval := common.GetPGInterval(mux.Vars(r)["interval"])

	api.db.Raw(`
		SELECT
			devices.id as device_id,
			devices.device,
			count(*) as flow_count,
			devices.name as device_name,
			devices.info as device_info,
			sum(flows.byte) as total_bytes,
			sum(flows.packet) as total_packets
		FROM
			flows
		JOIN
			devices
		ON
			flows.device_id = devices.id
		WHERE
			devices.device = '` + mux.Vars(r)["device"] + `'
			AND flows.created_at > NOW() - INTERVAL '` + interval + `'
		GROUP BY
			devices.id, devices.device, device_name, device_info;
	`).Scan(&result)

	_ = json.NewEncoder(w).Encode(&result)
	return

}

// get a device
// by it's IP address
// to update Name, Info in UI
func (api *APIServer) getByDevice(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	var obj model.Device

	api.db.Model(&model.Device{}).
		Where("device = ?", mux.Vars(r)["device"]).First(&obj)

	if obj.ID == 0 {
		// not found
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "There are no records in the database regarding your request",
		})
		return
	}

	_ = json.NewEncoder(w).Encode(&obj)
	return

}

// save a device new values to db
// by it's IP address
// to update Name, Info in UI side
func (api *APIServer) setByDevice(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	//extract info from json body
	type JSONBodyReq struct {
		Name string `json:"Name"`
		Info string `json:"Info"`
	}

	var jbr JSONBodyReq
	err := json.NewDecoder(r.Body).Decode(&jbr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"result":  true,
			"message": err,
		})
		return
	}

	var obj model.Device

	api.db.Model(&model.Device{}).
		Where("device = ?", mux.Vars(r)["device"]).First(&obj)

	if obj.ID == 0 {
		// not found
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "There are no records in the database regarding your request",
		})
		return
	}

	obj.Name = jbr.Name
	obj.Info = jbr.Info

	if dbc := api.db.Save(&obj); dbc.Error == nil {
		// created
		api.d.Verbose(fmt.Sprint("object updated ", mux.Vars(r)["device"]), logrus.InfoLevel)

	} else {
		// not created
		api.d.Verbose(fmt.Sprint("can not update object due to error", dbc.Error), logrus.ErrorLevel)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": fmt.Sprint("can not update object due to error", dbc.Error),
		})
		return
	}

	_ = json.NewEncoder(w).Encode(&obj)
	return

}

// // this function will return an array of
// // packets, bytes .... group by the interval that a user provides
// // for example if user provide 15m
// // it will group the total_byes and total_packets
// // into 1 minute and if provide 1w will group by them by day
// func (api *APIServer) getDeviceSummaryGroupByIntervalByDeviceID(w http.ResponseWriter, r *http.Request) {
// 	// add this line for debugging time
// 	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

// 	// summary
// 	type Result struct {
// 		Time time.Time `json:"_time" gorm:"column:_time"`

// 		DeviceID uint   `json:"device_id"`
// 		Device   string `json:"device"`

// 		FlowCount uint `json:"flow_count"`

// 		DeviceName string `json:"device_name"`

// 		DeviceInfo string `json:"device_info"`

// 		TotalBytes   uint `json:"total_bytes"`
// 		TotalPackets uint `json:"total_packets"`
// 	}
// 	var result []Result

// 	// Set header content type to application/json
// 	w.Header().Set("Content-Type", "application/json")

// 	// extract interval
// 	interval := common.GetPGInterval(mux.Vars(r)["interval"])

// 	timeTruncate := common.GetPGGroupByInterval(mux.Vars(r)["interval"])

// 	api.db.Raw(`
// 		SELECT
// 			date_trunc('`+timeTruncate+`', flows.created_at) as "_time",
// 			devices.id as device_id,
// 			devices.device,
// 			devices.name as device_name,
// 			devices.info as device_info,
// 			count(*) as flow_count,
// 			sum(flows.byte) as total_bytes,
// 			sum(flows.packet) as total_packets
// 		FROM
// 			flows
// 		JOIN
// 			devices
// 		ON
// 			flows.device_id = devices.id
// 		WHERE
// 			flows.created_at > NOW() - INTERVAL '`+interval+`'
// 		AND
// 			flows.device_id = ?
// 		GROUP BY
// 			_time, devices.id, devices.device, device_name, device_info
// 		ORDER BY
// 			_time;
// 	`, mux.Vars(r)["deviceID"]).Scan(&result)

// 	_ = json.NewEncoder(w).Encode(&result)
// 	return

// }
