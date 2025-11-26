package dto

import "time"

type URLdto struct {
	OriginalURL string     `json:"original_url" validate:"required,url"`
	ShortCode   string     `json:"short_code" validate:"omitempty, max=16"`
	ExpiresAt   *time.Time `json:"expires_at"`
}
