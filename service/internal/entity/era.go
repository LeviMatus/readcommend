package entity

type Era struct {
	ID      int32  `json:"id"`
	Title   string `json:"title"`
	MinYear *int16 `json:"minYear,omitempty"`
	MaxYear *int16 `json:"maxYear,omitempty"`
}
