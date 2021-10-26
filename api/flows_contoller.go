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

// get all flows
// by threat ID
func (api *APIServer) getAllFlowsByThreatID(w http.ResponseWriter, r *http.Request) {

	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Extract query string for pagination
	pg := common.ExtractPaginationQueryString(r.URL.Query())

	fields := []string{"f.src_port_name", "f.dst_port_name",
		"f.src_host", "f.dst_host", "f.next_hop_host", "f.protocol"}
	strWhere := common.PaginationStrWhereBuilder(pg["Filter"].(string), fields)

	// extract interval
	interval := common.GetPGInterval(mux.Vars(r)["interval"])

	// extract threat id
	threatID := -1
	if tr, err := strconv.Atoi(mux.Vars(r)["threatID"]); err != nil {
		threatID = -1
	} else {
		threatID = tr
	}

	// struct to return
	type Result struct {
		SrcIsThreat bool   `json:"src_is_threat"`
		SrcHost     string `json:"src_host"`

		DstIsThreat bool   `json:"dst_is_threat"`
		DstHost     string `json:"dst_host"`

		NextHopIsThreat bool   `json:"next_hop_is_threat"`
		NextHopHost     string `json:"next_hop_host"`

		SrcPortName string `json:"src_port_name"`
		DstPortName string `json:"dst_port_name"`

		Protocol string `json:"protocol"`

		FlagFin bool `json:"flag_fin"`
		FlagSyn bool `json:"flag_syn"`
		FlagRst bool `json:"flag_rst"`
		FlagPsh bool `json:"flag_psh"`
		FlagAck bool `json:"flag_ack"`
		FlagUrg bool `json:"flag_urg"`
		FlagEce bool `json:"flag_ece"`
		FlagCwr bool `json:"flag_cwr"`

		Byte   uint `json:"byte"`
		Packet uint `json:"packet"`

		CreatedAt time.Time `json:"created_at" gorm:"created_at"`
	}

	strQuerySelectFrom := `
	(SELECT

		src_threat_id,
		flows.src_is_threat,
		(SELECT hosts.host FROM hosts WHERE flows.src_host_id = hosts.id) as src_host,

		dst_threat_id,
		flows.dst_is_threat,
		(SELECT hosts.host FROM hosts WHERE flows.dst_host_id = hosts.id) as dst_host,

		next_hop_threat_id,
		flows.next_hop_is_threat,
		(SELECT hosts.host FROM hosts WHERE flows.next_hop_id = hosts.id) as next_hop_host,

		(SELECT ports.port_name FROM ports WHERE flows.src_port_id = ports.id) as src_port_name,

		(SELECT ports.port_name FROM ports WHERE flows.dst_port_id = ports.id) as dst_port_name,

		(SELECT protocols.protocol_name FROM protocols WHERE flows.protocol_id = protocols.id) as protocol,

		flows.flag_fin,
		flows.flag_syn,
		flows.flag_rst,
		flows.flag_psh,
		flows.flag_ack,
		flows.flag_urg,
		flows.flag_ece,
		flows.flag_cwr,

		flows.byte,
		flows.packet,

		flows.created_at
	FROM
		flows
	) f
	`

	strQuery := `
		SELECT * FROM ` + strQuerySelectFrom + `

		WHERE
		(
			(
				f.src_threat_id = ` + fmt.Sprintf("%v", threatID) + `
				OR
				f.dst_threat_id = ` + fmt.Sprintf("%v", threatID) + `
				OR
				f.next_hop_threat_id = ` + fmt.Sprintf("%v", threatID) + `
			)
			AND
			f.created_at > NOW() - INTERVAL '` + interval + `'
		)
		AND
	`
	if strWhere != "" {
		strQuery += ` ( ` + strWhere + ` )
		`
	}

	if pg["Order"].(string) != "" {
		strQuery += ` ORDER BY ` + pg["Order"].(string) + ` ` + pg["OrderType"].(string)
	} else {
		strQuery += ` ORDER BY ` + pg["Order"].(string) + ` ` + pg["OrderType"].(string)
	}

	strQuery += ` LIMIT ` + fmt.Sprintf("%v", pg["PerPage"].(int)) + `
			OFFSET ` + fmt.Sprintf("%v", pg["PerPage"].(int)*(pg["Page"].(int)-1)) + ` ;`

	var filteredResult []*Result
	api.db.Raw(strQuery).Scan(&filteredResult)

	type CountStruct struct {
		Count uint `json:"count" gorm:"count"`
	}
	var total CountStruct
	api.db.Raw("SELECT count(*) as count from " + strQuerySelectFrom + " WHERE " + `
	(
		(
			f.src_threat_id = ` + fmt.Sprintf("%v", threatID) + `
			OR
			f.dst_threat_id = ` + fmt.Sprintf("%v", threatID) + `
			OR
			f.next_hop_threat_id = ` + fmt.Sprintf("%v", threatID) + `
		)
		AND
		f.created_at > NOW() - INTERVAL '` + interval + `'
	)
	AND	( ` + strWhere + ") ;").Scan(&total)

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
			Result     []*Result
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
