package entity

type Era struct {
	ID      int32  `json:"id"`
	Title   string `json:"title"`
	MinYear *int16 `json:"minYear"`
	MaxYear *int16 `json:"maxYear"`
}
