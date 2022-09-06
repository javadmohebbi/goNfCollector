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

// get all ports
func (api *APIServer) getAllPorts(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Extract query string for pagination
	pg := common.ExtractPaginationQueryString(r.URL.Query())

	fields := []string{"ports.port_name", "ports.port_proto", "ports.info"}
	strWhere := common.PaginationStrWhereBuilder(pg["Filter"].(string), fields)

	strQuery := `
		SELECT
			id as port_id, ports.port_name as port, info as port_info, port_proto, created_at, updated_at
		FROM
			ports
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
		PortID    uint      `json:"port_id" gorm:"port_id"`
		Port      string    `json:"port" gorm:"port"`
		PortProto string    `json:"port_proto" gorm:"port_proto"`
		PortInfo  string    `json:"port_info" gorm:"port_info"`
		UpdatedAt time.Time `json:"updated_at" gorm:"updated_at"`
		CreatedAt time.Time `json:"created_at" gorm:"created_at"`
	}
	var filteredResult []*FilteredResult
	api.db.Raw(strQuery).Scan(&filteredResult)

	type CountStruct struct {
		Count uint `json:"count" gorm:"count"`
	}
	var total CountStruct
	api.db.Raw("SELECT count(*) as count from ports WHERE " + strWhere + ";").Scan(&total)

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

// this function will get information about the top source ports
// on a device based on an intervals like 1m, 2h ...
func (api *APIServer) getTopPortsByDeviceByInterval(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// summary
	type Result struct {
		PortID   uint   `json:"port_id"`
		PortName string `json:"port_name"`
		Info     string `json:"info"`

		PortProto string `json:"port_proto"`

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
			` + dir + `_port_id as port_id,
			ports.port_name,
			ports.port_proto,
			ports.info,
			sum(flows.byte) as total_bytes,
			sum(flows.packet) as total_packets
		FROM
			flows
		JOIN
			ports
		ON
			flows.` + dir + `_port_id = ports.id
		WHERE
			flows.created_at > NOW() - INTERVAL '` + interval + `'
			AND flows.device_id = ` + deviceID + `
		GROUP BY
			` + dir + `_port_id, ports.port_name, ports.port_proto, ports.info
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
func (api *APIServer) getPortReportWhenHostSrcOrDst(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	type Result struct {
		Port interface{} `json:"ports"`
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
	result.Port = api.portRptWhenHostSrcOrDst(host, interval, top, dir, h.ID)

	_ = json.NewEncoder(w).Encode(&result)
	return

}

// get a port by id to update Info in UI
func (api *APIServer) getPortByID(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	var obj model.Port

	api.db.Model(&model.Port{}).
		Where("id = ?", mux.Vars(r)["id"]).First(&obj)

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

// save a port new values to db
// by it's Id to update Info in UI side
func (api *APIServer) setPortByID(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	//extract info from json body
	type JSONBodyReq struct {
		ID   uint   `json:"id"`
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

	var obj model.Port

	api.db.Model(&model.Port{}).
		Where("id = ?", mux.Vars(r)["id"]).First(&obj)

	if obj.ID == 0 {
		// not found
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "There are no records in the database regarding your request",
		})
		return
	}

	obj.ID = jbr.ID
	obj.Info = jbr.Info

	if dbc := api.db.Save(&obj); dbc.Error == nil {
		// created
		api.d.Verbose(fmt.Sprint("object updated ", mux.Vars(r)["id"]), logrus.InfoLevel)

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
