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

// get report about host when seen as SOURCE or DESTINATION
func (api *APIServer) hostRptWhenSrcOrDst(host, interval, top, direction string, hostID uint) interface{} {

	type Result struct {
		WhenSrc []*HostRPTWhenSrcOrDstResult `json:"list"`
	}
	var result Result

	// threats
	result.WhenSrc = api._hostRptWhenSrcOrDst(host, interval, top, direction, hostID)

	return result
}

// get report about host when seen as SOURCE or DESTINATION
func (api *APIServer) _hostRptWhenSrcOrDst(host, interval, top, direction string, hostID uint) []*HostRPTWhenSrcOrDstResult {

	var result []*HostRPTWhenSrcOrDstResult

	opos_dir := "dst"
	if direction == "src" {
		opos_dir = "dst"
	} else {
		opos_dir = "src"
	}

	strQuery := `
		SELECT
			` + direction + `_host_id as host_id,
			hosts.host as host,
			hosts.info as host_info,
			sum(flows.byte) as total_bytes,
			sum(flows.packet) as total_packets
		FROM
			flows
		JOIN
			hosts
		ON
			hosts.id = dst_host_id
		WHERE
			flows.created_at > NOW() - INTERVAL '` + interval + `'
			AND flows.` + opos_dir + `_host_id = '` + fmt.Sprintf("%v", hostID) + `'
		GROUP BY
			` + direction + `_host_id, hosts.host, hosts.info
		ORDER BY
			total_bytes desc
		LIMIT ` + top + `
	;`

	api.db.Raw(strQuery).Scan(&result)

	return result
}
