package entity

// Book is a domain entity that contains information about a book.
type Book struct {
	// ID is the primary identifier of a Book.
	ID int32 `json:"id"`

	// Title is the name of the book.
	Title string `json:"title"`

	// YearPublished indicates the year in which the book was published.
	YearPublished int16 `json:"yearPublished"`

	// Rating expresses the popularity/review of the book in a 2-point floating decimal
	// score out of 5.
	Rating float32 `json:"rating"`

	// Pages is the count of pages in the Book.
	Pages int16 `json:"pages"`

	// Genre is the categorical genre of the Book.
	Genre Genre `json:"genre"`

	// Author is the author/writer of the Book.
	Author Author `json:"author"`
}
