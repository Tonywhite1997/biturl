package domain

import "time"

type Stat struct {
	ID         string    `json:"id" gorm:"PrimaryKey"`
	URLShortID string    `json:"url_short_id"`
	UserIP     string    `json:"user_ip"`
	UserAgent  string    `json:"user_agent"`
	Referer    string    `json:"referer"`
	Country    string    `json:"country"`
	City       string    `json:"city"`
	Device     string    `json:"device"`
	OS         string    `json:"os"`
	Browser    string    `json:"browser"`
	Timestamp  time.Time `json:"timestamp"`
}
