package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/goNfCollector/common"
	"github.com/goNfCollector/database/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// get all hosts
func (api *APIServer) getAllHosts(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Extract query string for pagination
	pg := common.ExtractPaginationQueryString(r.URL.Query())

	fields := []string{"hosts.host", "hosts.info"}
	strWhere := common.PaginationStrWhereBuilder(pg["Filter"].(string), fields)

	strQuery := `
		SELECT
			id as host_id, host, info as host_info, created_at, updated_at
		FROM
			hosts
	`

	if strWhere != "" {
		strQuery += ` WHERE ` + strWhere + `
		`
	}

	if pg["Order"].(string) != "" {
		strQuery += ` ORDER BY ` + pg["Order"].(string) + ` ` + pg["OrderType"].(string)
	}

	// if !pg["NoPagination"].(bool) {
	// only pagination will be allowed since
	// there are many hosts in the database
	// usually and without pagination
	// its not possible to get more results
	strQuery += ` LIMIT ` + fmt.Sprintf("%v", pg["PerPage"].(int)) + `
			OFFSET ` + fmt.Sprintf("%v", pg["PerPage"].(int)*(pg["Page"].(int)-1)) + ` ;`
	// } else {
	// 	strQuery += `  ;`
	// }

	type FilteredResult struct {
		HostID    uint      `json:"host_id" gorm:"host_id"`
		Host      string    `json:"host" gorm:"host"`
		HostInfo  string    `json:"host_info" gorm:"host_info"`
		UpdatedAt time.Time `json:"updated_at" gorm:"updated_at"`
		CreatedAt time.Time `json:"created_at" gorm:"created_at"`
	}
	var filteredResult []*FilteredResult
	api.db.Raw(strQuery).Scan(&filteredResult)

	type CountStruct struct {
		Count uint `json:"count" gorm:"count"`
	}
	var total CountStruct
	api.db.Raw("SELECT count(*) as count from hosts WHERE " + strWhere + ";").Scan(&total)

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

}

// this function will get information about the top source hosts
// on a device based on an intervals like 1m, 2h ...
func (api *APIServer) getTopHostsByDeviceByInterval(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// summary
	type Result struct {
		HostID   uint   `json:"host_id"`
		Host     string `json:"host"`
		HostInfo string `json:"host_info"`

		TotalBytes   uint `json:"total_bytes"`
		TotalPackets uint `json:"total_packets"`
	}
	var result []Result

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// extract interval
	interval := common.GetPGInterval(mux.Vars(r)["interval"])

	deviceID := mux.Vars(r)["deviceID"]
	if deviceID == "" {
		deviceID = "1"
	}

	top := mux.Vars(r)["top"]
	if _, err := strconv.Atoi(top); err != nil {
		top = "10"
	}

	dir := mux.Vars(r)["direction"]
	if dir != "src" && dir != "dst" {
		dir = "dst"
	}

	strQuery := `
		SELECT
			` + dir + `_host_id as host_id,
			hosts.host as host,
			hosts.info as host_info,
			sum(flows.byte) as total_bytes,
			sum(flows.packet) as total_packets
		FROM
			flows
		JOIN
			hosts
		ON
			flows.` + dir + `_host_id = hosts.id
		WHERE
			flows.created_at > NOW() - INTERVAL '` + interval + `'
			AND flows.device_id = ` + deviceID + `
		GROUP BY
			` + dir + `_host_id, hosts.host, hosts.info
		ORDER BY
			total_bytes desc
		LIMIT ` + top + `
	;`
	// fmt.Println(strQuery)

	api.db.Raw(strQuery).Scan(&result)

	_ = json.NewEncoder(w).Encode(&result)
	return

}

// this function will get information about the top source hosts
// on a device based on an intervals like 1m, 2h ...
func (api *APIServer) getHostReport(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	type Result struct {
		Threat interface{} `json:"threats"`
	}
	var result Result

	// extract interval
	interval := common.GetPGInterval(mux.Vars(r)["interval"])

	// extract top
	top := mux.Vars(r)["top"]
	if _, err := strconv.Atoi(top); err != nil {
		top = "10"
	}

	// extract host
	host := mux.Vars(r)["host"]

	// fetch host from db
	var h model.Host
	api.db.Model(&model.Host{}).Where("host = ?", host).First(&h)

	// get threat report
	result.Threat = api.hostRptThreat(host, interval, h.ID)

	_ = json.NewEncoder(w).Encode(&result)
	return

}

// this function will a report the provided host
// is a SOURCE or DESTINATION in flow recordset based on the time interval
func (api *APIServer) getHostReportWhenSrcOrDst(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	type Result struct {
		Host interface{} `json:"hosts"`
	}
	var result Result

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// extract interval
	interval := common.GetPGInterval(mux.Vars(r)["interval"])

	// extract top
	top := mux.Vars(r)["top"]
	if _, err := strconv.Atoi(top); err != nil {
		top = "10"
	}

	// extract host
	host := mux.Vars(r)["host"]

	dir := mux.Vars(r)["direction"]
	if dir != "src" && dir != "dst" {
		dir = "dst"
	}

	// fetch host from db
	var h model.Host
	api.db.Model(&model.Host{}).Where("host = ?", host).First(&h)

	// get threat report
	result.Host = api.hostRptWhenSrcOrDst(host, interval, top, dir, h.ID)

	_ = json.NewEncoder(w).Encode(&result)
	return

}
