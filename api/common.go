package api

import (
	"fmt"
	"net/http"
)

type PaginatedRequest struct {
	Page      int `form:"page,default=0"`
	PageLimit int `form:"pageLimit,default=20"`
}

type PaginatedResponse struct {
	CurrentPage int           `json:"currentPage"`
	TotalPages  int           `json:"totalPages"`
	Items       []interface{} `json:"items"`
}

func NewPaginatedResponse(currPage int, totalPages int, itemSize int) *PaginatedResponse {
	return &PaginatedResponse{
		CurrentPage: currPage,
		TotalPages:  totalPages,
		Items:       make([]interface{}, itemSize),
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewAlreadyExistsError() (int, *ErrorResponse) {
	return http.StatusBadRequest, &ErrorResponse{
		Message: "Resource already exists",
	}
}

func NewNotFoundError() (int, *ErrorResponse) {
	return http.StatusNotFound, &ErrorResponse{
		Message: fmt.Sprintf("Resource not found"),
	}
}

func NewBadRequestError(err error) (int, *ErrorResponse) {
	return http.StatusBadRequest, &ErrorResponse{
		Message: fmt.Sprintf("Invalid request %s", err),
	}
}

func NewVersionError() (int, *ErrorResponse) {
	return http.StatusConflict, &ErrorResponse{
		Message: fmt.Sprintf("Resource version do not match"),
	}
}
