package entity

type Size struct {
	ID       int32  `json:"id"`
	Title    string `json:"title"`
	MinPages *int16 `json:"minPages,omitempty"`
	MaxPages *int16 `json:"maxPages,omitempty"`
}
