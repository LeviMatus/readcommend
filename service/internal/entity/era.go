package entity

// Era is a domain entity which indicates the literary timeframe in which
// a Book was published, as described by minimum and maximum years.
type Era struct {
	// ID is the primary identifier for an Era category.
	ID int32 `json:"id"`

	// Title is the name/title of an Era.
	Title string `json:"title"`

	// MinYear is the minimum year of publication for a Book to be considered
	// a part of an Era. If nil, then there is no floor to the considered year.
	MinYear *int16 `json:"minYear,omitempty"`

	// MaxYear is the maximum year of publication for a Book to be considered
	// a part of an Era. If nil, then there is no cap to the considered year.
	MaxYear *int16 `json:"maxYear,omitempty"`
}
