package model

type PaginationModel struct {
	Page    int `json:"page,omitempty"`
	PerPage int `json:"perPage,omitempty"`
	Total   int `json:"total,ommitempty"`
}
