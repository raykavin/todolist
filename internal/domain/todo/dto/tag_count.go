package dto

// TagCount represents tag usage count
type TagCount struct {
	Tag   string `json:"tag"`
	Count int64  `json:"count"`
}
