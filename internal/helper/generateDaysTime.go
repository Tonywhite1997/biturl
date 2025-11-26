package helper

import "time"

func GenerateDate(inDaysTime uint) *time.Time {
	expires := time.Now().Add(24 * time.Duration(inDaysTime) * time.Hour)
	return &expires
}
