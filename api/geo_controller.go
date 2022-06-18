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

// get all Geos
func (api *APIServer) getAllGeos(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Extract query string for pagination
	pg := common.ExtractPaginationQueryString(r.URL.Query())

	fields := []string{"geos.country_short", "geos.country_long", "geos.region", "geos.city"}
	strWhere := common.PaginationStrWhereBuilder(pg["Filter"].(string), fields)

	strQuery := `
		SELECT
			id as geo_id, country_short, country_long, region, city, latitude, longitude, created_at, updated_at
		FROM
			geos
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
		GeoID        uint      `json:"geo_id" gorm:"geo_id"`
		CountryShort string    `json:"country_short" gorm:"country_short"`
		CountryLong  string    `json:"country_long" gorm:"country_long"`
		Region       string    `json:"region" gorm:"region"`
		City         string    `json:"city" gorm:"city"`
		UpdatedAt    time.Time `json:"updated_at" gorm:"updated_at"`

		Latitude  float32 `json:"latitude" gorm:"latitude"`
		Longitude float32 `json:"longitude" gorm:"longitude"`

		CreatedAt time.Time `json:"created_at" gorm:"created_at"`
	}
	var filteredResult []*FilteredResult
	api.db.Raw(strQuery).Scan(&filteredResult)

	type CountStruct struct {
		Count uint `json:"count" gorm:"count"`
	}
	var total CountStruct
	api.db.Raw("SELECT count(*) as count from geos WHERE " + strWhere + ";").Scan(&total)

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

// this function will get information about the top source/dst geo countries
// on a device based on an intervals like 1m, 2h ...
func (api *APIServer) getTopGeoCountryByDeviceByInterval(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// summary
	type Result struct {
		CountryLong  string `json:"country_long"`
		CountryShort string `json:"country_short"`

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
			geos.country_short as country_short,
			geos.country_long as country_long,
			sum(flows.byte) as total_bytes,
			sum(flows.packet) as total_packets
		FROM
			flows
		JOIN
			geos
		ON
			flows.` + dir + `_geo_id = geos.id
		WHERE
			flows.created_at > NOW() - INTERVAL '` + interval + `'
			AND flows.device_id = ` + deviceID + `
		GROUP BY
			geos.country_short, geos.country_long
		ORDER BY
			total_bytes desc
		LIMIT ` + top + `
	;`
	// fmt.Println(strQuery)

	api.db.Raw(strQuery).Scan(&result)

	_ = json.NewEncoder(w).Encode(&result)
	return

}

// get geo report only when SOURCE OF DESTINATION report based on interval
func (api *APIServer) getGeoReportWhenHostSrcOrDst(w http.ResponseWriter, r *http.Request) {
	// add this line for debugging time
	api.d.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	type Result struct {
		Geo interface{} `json:"geo"`
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
	result.Geo = api.geoRptWhenHostSrcOrDst(host, interval, top, dir, h.ID)

	_ = json.NewEncoder(w).Encode(&result)
	return

}
