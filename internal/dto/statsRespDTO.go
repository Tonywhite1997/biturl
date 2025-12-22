package dto

import "time"

type StatsResponse struct {
	TotalClicks    uint64         `json:"total_clicks"`
	UniqueVisitors uint64         `json:"unique_visitors"`
	Browsers       []BrowserStats `json:"browsers"`
	Devices        []DeviceStats  `json:"devices"`
	Countries      []CountryStats `json:"countries"`
	DailyStats     []DailyStats   `json:"daily_stats"`
	OriginalURL    string         `json:"original_url"`
	ExpiresAt      *time.Time     `json:"expires_at"`
}
