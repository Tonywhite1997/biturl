package dto

import "time"

type DateFilterReq struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
