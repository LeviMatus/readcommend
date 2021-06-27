package domain

type Author struct {
	// TODO: find appropriate int category
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
