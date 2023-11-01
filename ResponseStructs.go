package main

import "time"

type EndpointResponse struct {
	SystemTime time.Time `json:"systemtime"`
	Status     int       `json:"status"`
}

type AllResponse struct {
	TableName   string `json:"table"`
	RecordCount int64  `json:"recordCount"`
}
