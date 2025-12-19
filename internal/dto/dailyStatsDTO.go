package dto

import "time"

type DailyStats struct {
	Date  time.Time `json:"date"`
	Count uint64    `json:"count"`
}
