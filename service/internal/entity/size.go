package entity

// Size is a domain entity that expresses the size
// of a Book or literary work in terms of a range of pages.
type Size struct {
	// ID describes the unique identifier to a size category.
	ID int32 `json:"id"`

	// Title is the title or name of a size category.
	Title string `json:"title"`

	// MinPages is the minimum number of pages for a book to be considered
	// to fall within this size category. If nil, there is no floor to the pages required.
	MinPages *int16 `json:"minPages,omitempty"`

	// MaxPages is the maximum number of pages required for a book to be considered
	// to fall within this size category. If nil, there is no cap to the pages required.
	MaxPages *int16 `json:"maxPages,omitempty"`
}
