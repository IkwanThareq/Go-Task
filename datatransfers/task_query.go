package datatransfers

type TaskQueryParams struct {
	Status   string `form:"status"`
	Priority string `form:"priority"`
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
}

type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
}
