package api

import "fmt"

// get report about host when seen as SOURCE or DESTINATION
func (api *APIServer) geoRptWhenHostSrcOrDst(host, interval, top, direction string, hostID uint) interface{} {

	type Result struct {
		WhenSrcOrDst []*GeoRPTWhenHostSrcOrDstResult `json:"list"`
	}
	var result Result

	// threats
	result.WhenSrcOrDst = api._geoRptWhenHostSrcOrDst(host, interval, top, direction, hostID)

	return result
}

// get report about host when seen as SOURCE or DESTINATION
func (api *APIServer) _geoRptWhenHostSrcOrDst(host, interval, top, direction string, hostID uint) []*GeoRPTWhenHostSrcOrDstResult {

	var result []*GeoRPTWhenHostSrcOrDstResult

	strQuery := `
		SELECT
		    ` + direction + `_geo_id as geo_id,
			geos.country_short as country_short,
			geos.country_long as country_long,

			sum(flows.byte) as total_bytes,
			sum(flows.packet) as total_packets
		FROM
			flows
		JOIN
			geos
		ON
			geos.id = flows.` + direction + `_geo_id
		WHERE
			flows.created_at > NOW() - INTERVAL '` + interval + `'
			AND (
					flows.src_host_id = ` + fmt.Sprintf("%v", hostID) + `
					OR
					flows.dst_host_id = ` + fmt.Sprintf("%v", hostID) + `
			)
		GROUP BY
			` + direction + `_geo_id, geos.country_short, geos.country_long
		ORDER BY
			total_bytes desc
		LIMIT ` + top + `
	;`

	api.db.Raw(strQuery).Scan(&result)

	return result
}
