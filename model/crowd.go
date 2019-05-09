package model

import "time"

type Crowd struct {
	Count int `json:"count"`
	Time  time.Time
}

type TotalSniffed struct {
	Count int `json:"count"`
}
