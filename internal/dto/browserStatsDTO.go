package dto

type BrowserStats struct {
	Browser string `json:"browser_type"`
	Count   uint64 `json:"count"`
}
