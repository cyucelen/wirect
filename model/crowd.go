package model

import "time"

type Crowd struct {
	Count int `json:"count"`
	Time  time.Time
}
