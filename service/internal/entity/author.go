package entity

// Author is the writer of a Book.
type Author struct {
	// ID is the primary identifier of a Book's published Author.
	ID int32 `json:"id"`

	// The FirstName of a Book's published Author.
	FirstName string `json:"firstName"`

	// The LastName of a Book's published Author.
	LastName string `json:"lastName"`
}
