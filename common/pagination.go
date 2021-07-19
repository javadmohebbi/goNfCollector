package common

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ExtractPaginationQueryString - Extract pagination from get request querystring
func ExtractPaginationQueryString(qry url.Values) map[string]interface{} {

	pg := map[string]interface{}{
		"Page":         1,
		"PerPage":      15,
		"Filter":       "",
		"Order":        "created_at",
		"OrderType":    "desc",
		"NoPagination": false,
	}

	// Extract noPagination from query
	paramNoPagination := qry.Get("noPagination")
	if paramNoPagination == "0" || paramNoPagination == "" {
		pg["NoPagination"] = false
	} else {
		pg["NoPagination"] = true
	}

	// Extract Page from query
	if paramPage, err := strconv.Atoi(qry.Get("page")); err == nil {
		pg["Page"] = paramPage
	}

	// Extract item per page
	if paramPerPage, err := strconv.Atoi(qry.Get("perPage")); err == nil {
		pg["PerPage"] = paramPerPage
	}

	// Extract needed filter
	if qry.Get("filter") != "" {
		pg["Filter"] = qry.Get("filter")
	}

	// Extract needed filter
	if qry.Get("IsPaginate") == "" {
		pg["IsPaginate"] = false
	} else if qry.Get("IsPaginate") == "false" {
		pg["IsPaginate"] = false
	} else {
		pg["IsPaginate"] = false
	}

	// Extract order by direction
	if qry.Get("order") != "" {
		pg["Order"] = qry.Get("order")
		switch qry.Get("orderType") {
		case "asc":
			pg["OrderType"] = "asc"
		case "desc":
			pg["OrderType"] = "desc"
		}
	}

	return pg
}

func PaginationStrWhereBuilder(filter string, fields []string) (strWhere string) {
	strWhere = ""
	if filter != "" {
		for k, v := range fields {
			// likeFlt := "%" + strings.ReplaceAll(r.Pagination.Filter, " ", "%") + "%"

			// to lower to make a better filter
			likeFlt := strings.ToLower(filter)

			// since PROTOCOL stores as CAPITAL letter
			// it must be Capitalized in WHERE statement
			// if v == "protocols.protocol_name" || v == "ports.port_proto" {
			// 	likeFlt = strings.ToUpper(likeFlt)
			// }

			if k == len(fields)-1 {
				strWhere += fmt.Sprintf("LOWER(%s) LIKE '%%%v%%'", v, likeFlt)
			} else {
				strWhere += fmt.Sprintf("LOWER(%s) LIKE '%%%v%%' OR ", v, likeFlt)
			}
		}
	} else {
		strWhere = " 1 = 1 "
	}

	return strWhere
}
