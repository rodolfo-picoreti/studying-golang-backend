package api

type PaginatedItems struct {
	CurrentPage int           `json:"currentPage"`
	TotalPages  int           `json:"totalPages"`
	Items       []interface{} `json:"items"`
}

func NewPaginatedItems(currPage int, totalPages int, itemSize int) *PaginatedItems {
	return &PaginatedItems{
		CurrentPage: currPage,
		TotalPages:  totalPages,
		Items:       make([]interface{}, itemSize),
	}
}
