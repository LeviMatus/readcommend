package entity

// Genre is a domain entity that described the genre of a Book or
// literary work.
type Genre struct {
	// ID is the primary identifier of a Genre.
	ID int32 `json:"id"`

	// Title is the name of a Genre, such as "Childrens", "Fiction", of "Fantasy/SciFy."
	Title string `json:"title"`
}
