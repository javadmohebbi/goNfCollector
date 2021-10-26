package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/goNfCollector/common"
	"github.com/goNfCollector/database/model"
	"github.com/sirupsen/logrus"
)

// get all threats
func (api *APIServer) getAllThreats(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Extract query string for pagination
	pg := common.ExtractPaginationQueryString(r.URL.Query())

	fields := []string{"threats.source", "threats.kind", "hosts.host", "hosts.info"}
	strWhere := common.PaginationStrWhereBuilder(pg["Filter"].(string), fields)

	// fmt.Println(strWhere)

	strQuery := `
		SELECT
			threats.id as threat_id, threats.source as threat_source, threats.counter as threat_counter,
			threats.reputation as threat_reputation, threats.kind as threat_kind,
			hosts.host as threat_host, hosts.id as threat_host_id, hosts.info as threat_host_info,
			threats.acked as threat_acked, threats.closed as threat_closed,
			threats.false_positive as threat_false_positive,
			threats.created_at as created_at,
			threats.updated_at as updated_at
		FROM
			threats
		JOIN
			hosts
		ON
			hosts.id = threats.host_id
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
		ThreatID            uint      `json:"threat_id" gorm:"threat_id"`
		ThreatSource        string    `json:"threat_source" gorm:"threat_source"`
		ThreatCounter       uint      `json:"threat_counter" gorm:"threat_counter"`
		ThreatReputation    uint      `json:"threat_reputation" gorm:"threat_reputation"`
		ThreatKind          string    `json:"threat_kind" gorm:"threat_kind"`
		ThreatHost          string    `json:"threat_host" gorm:"threat_host"`
		ThreatHostInfo      string    `json:"threat_host_info" gorm:"threat_host_info"`
		ThreatHostID        uint      `json:"threat_host_id" gorm:"threat_host_id"`
		ThreatAcked         bool      `json:"threat_acked" gorm:"threat_acked"`
		ThreatClosed        bool      `json:"threat_closed" gorm:"threat_closed"`
		ThreatFalsePositive bool      `json:"threat_false_positive" gorm:"threat_false_positive"`
		UpdatedAt           time.Time `json:"updated_at" gorm:"updated_at"`
		CreatedAt           time.Time `json:"created_at" gorm:"created_at"`
	}
	var filteredResult []*FilteredResult
	api.db.Raw(strQuery).Scan(&filteredResult)

	type CountStruct struct {
		Count uint `json:"count" gorm:"count"`
	}
	var total CountStruct

	countQuery := "SELECT count(*) as count FROM (SELECT count (*) as count, threats.source,threats.kind,hosts.host, hosts.info FROM threats JOIN hosts "
	countQuery += " ON hosts.id = threats.host_id WHERE "
	countQuery += strWhere + " GROUP BY threats.source,threats.kind, hosts.host, hosts.info) as filtered;"

	api.db.Raw(countQuery).Scan(&total)

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
