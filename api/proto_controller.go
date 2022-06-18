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

// get all protocols
func (api *APIServer) getAllProtocols(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Extract query string for pagination
	pg := common.ExtractPaginationQueryString(r.URL.Query())

	fields := []string{"protocols.protocol_name", "protocols.info"}
	strWhere := common.PaginationStrWhereBuilder(pg["Filter"].(string), fields)

	strQuery := `
		SELECT
			id as protocol_id, protocols.protocol_name as protocol, info as protocol_info, created_at, updated_at
		FROM
			protocols
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
		ProtocolID   uint      `json:"protocol_id" gorm:"protocol_id"`
		Protocol     string    `json:"protocol" gorm:"protocol"`
		ProtocolInfo string    `json:"protocol_info" gorm:"protocol_info"`
		UpdatedAt    time.Time `json:"updated_at" gorm:"updated_at"`
		CreatedAt    time.Time `json:"created_at" gorm:"created_at"`
	}
	var filteredResult []*FilteredResult
	api.db.Raw(strQuery).Scan(&filteredResult)

	type CountStruct struct {
		Count uint `json:"count" gorm:"count"`
	}
	var total CountStruct
	api.db.Raw("SELECT count(*) as count from protocols WHERE " + strWhere + ";").Scan(&total)

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

// this function will get information about the top protos
// on a device based on an intervals like 1m, 2h ...
func (api *APIServer) getTopProtosByDeviceByInterval(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// summary
	type Result struct {
		ProtocolID   uint   `json:"protocol_id"`
		ProtocolName string `json:"protocol_name"`

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

	strQuery := `
		SELECT
			protocol_id,
			protocols.protocol_name as protocol_name,
			sum(flows.byte) as total_bytes,
			sum(flows.packet) as total_packets
		FROM
			flows
		JOIN
			protocols
		ON
			flows.protocol_id = protocols.id
		WHERE
			flows.created_at > NOW() - INTERVAL '` + interval + `'
			AND flows.device_id = ` + deviceID + `
		GROUP BY
			protocol_id, protocols.protocol_name
		ORDER BY
			total_bytes desc
		LIMIT ` + top + `
	;`
	// fmt.Println(strQuery)

	api.db.Raw(strQuery).Scan(&result)

	_ = json.NewEncoder(w).Encode(&result)
	return

}

// get port report only when SOURCE OF DESTINATION report based on interval
func (api *APIServer) getProtocolReportWhenHostSrcOrDst(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	type Result struct {
		Protocol interface{} `json:"protocols"`
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

	// get the report
	result.Protocol = api.protcolRptWhenHostSrcOrDst(host, interval, top, dir, h.ID)

	_ = json.NewEncoder(w).Encode(&result)
	return

}
