package api

import "fmt"

// get report about host when seen as SOURCE or DESTINATION
func (api *APIServer) portRptWhenHostSrcOrDst(host, interval, top, direction string, hostID uint) interface{} {

	type Result struct {
		WhenSrcOrDst []*PortRPTWhenHostSrcOrDstResult `json:"list"`
	}
	var result Result

	// threats
	result.WhenSrcOrDst = api._hostRptWhenHostSrcOrDst(host, interval, top, direction, hostID)

	return result
}

// get report about host when seen as SOURCE or DESTINATION
func (api *APIServer) _hostRptWhenHostSrcOrDst(host, interval, top, direction string, hostID uint) []*PortRPTWhenHostSrcOrDstResult {

	var result []*PortRPTWhenHostSrcOrDstResult

	strQuery := `
		SELECT
		    ` + direction + `_port_id as port_id,
			ports.port_proto as port_proto,
			ports.port_name as port_name,
			sum(flows.byte) as total_bytes,
			sum(flows.packet) as total_packets
		FROM
			flows
		JOIN
			ports
		ON
			ports.id = flows.` + direction + `_port_id
		WHERE
			flows.created_at > NOW() - INTERVAL '` + interval + `'
			AND (
					flows.src_host_id = ` + fmt.Sprintf("%v", hostID) + `
					OR
					flows.dst_host_id = ` + fmt.Sprintf("%v", hostID) + `
			)
		GROUP BY
			` + direction + `_port_id, ports.port_proto, ports.port_name
		ORDER BY
			total_bytes desc
		LIMIT ` + top + `
	;`

	api.db.Raw(strQuery).Scan(&result)

	return result
}
