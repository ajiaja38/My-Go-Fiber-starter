package model

type ResponseEntity[T any] struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    T               `json:"data,omitempty"`
	Meta    *MetaPagination `json:"meta,omitempty"`
}

type MetaPagination struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"totalPage"`
	TotalData int `json:"totalData"`
}

type ResponseError[T any] struct {
	ResponseEntity[T]
	Path string `json:"path"`
}
