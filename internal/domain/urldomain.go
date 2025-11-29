package domain

import "time"

const (
	MaxLength = 16
)

type URL struct {
	ID             uint       `json:"id" gorm:"PrimaryKey"`
	ShortCode      string     `json:"short_code" gorm:"size:16; uniqueIndex; not null"`
	StatsAccessKey string     `json:"stats_access_key" gorm:"unique"`
	OriginalURL    string     `json:"original_url" gorm:"type:text; not null"`
	ExpiresAt      *time.Time `json:"expires_at"`
	CreatedAt      time.Time  `json:"created_at"  gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at"  gorm:"autoUpdateTime"`
}
