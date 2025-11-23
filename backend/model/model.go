package model

import (
	"time"
)

type DatabaseHealth struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Database  string    `json:"database"`
	Uptime    string    `json:"uptime"`
}

type Meta struct {
	Total  int64 `json:"total"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

type APIResponse[T any] struct {
	Data  T      `json:"data,omitempty"`
	Meta  *Meta  `json:"meta,omitempty"`
	Error string `json:"error,omitempty"`
}

type Offset struct {
	Limit  int64 `query:"limit,default:20"`
	Offset int64 `query:"offset,default:0"`
}
