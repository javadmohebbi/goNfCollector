package api

import (
	"fmt"
)

// get report about threats
// based on host and interval
func (api *APIServer) hostRptThreat(host, interval string, hostID uint) interface{} {

	type Result struct {
		Threats []*HostRPTThreatsResult `json:"list"`
	}
	var result Result

	// threats
	result.Threats = api._hostRptThreats(host, interval, hostID)

	return result
}

// get report if host is part of a threat
func (api *APIServer) _hostRptThreats(host, interval string, hostID uint) []*HostRPTThreatsResult {

	var result []*HostRPTThreatsResult

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
		WHERE
			threats.host_id = ` + fmt.Sprintf("%v", hostID) + `
			AND threats.updated_at > NOW() - INTERVAL '` + interval + `'
	`

	api.db.Raw(strQuery).Scan(&result)

	return result
}
