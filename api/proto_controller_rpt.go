package api

import "fmt"

// get report about host when seen as SOURCE or DESTINATION
func (api *APIServer) protcolRptWhenHostSrcOrDst(host, interval, top, direction string, hostID uint) interface{} {

	type Result struct {
		WhenSrcOrDst []*ProtocolRPTWhenHostSrcOrDstResult `json:"list"`
	}
	var result Result

	// threats
	result.WhenSrcOrDst = api._protoRptWhenHostSrcOrDst(host, interval, top, direction, hostID)

	return result
}

// get report about host when seen as SOURCE or DESTINATION
func (api *APIServer) _protoRptWhenHostSrcOrDst(host, interval, top, direction string, hostID uint) []*ProtocolRPTWhenHostSrcOrDstResult {

	var result []*ProtocolRPTWhenHostSrcOrDstResult

	strQuery := `
		SELECT
		    flows.protocol_id,
			protocols.protocol as protocol,
			protocols.protocol_name as protocol_name,
			sum(flows.byte) as total_bytes,
			sum(flows.packet) as total_packets
		FROM
			flows
		JOIN
			protocols
		ON
			protocols.id = flows.protocol_id
		WHERE
			flows.created_at > NOW() - INTERVAL '` + interval + `'
			AND (
					flows.src_host_id = ` + fmt.Sprintf("%v", hostID) + `
					OR
					flows.dst_host_id = ` + fmt.Sprintf("%v", hostID) + `
			)
		GROUP BY
			protocol_id, protocols.protocol, protocols.protocol_name
		ORDER BY
			total_bytes desc
		LIMIT ` + top + `
	;`

	api.db.Raw(strQuery).Scan(&result)

	return result
}
