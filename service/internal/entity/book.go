package entity

type Book struct {
	ID            int32   `json:"id"`
	Title         string  `json:"title"`
	YearPublished int16   `json:"yearPublished"`
	Rating        float32 `json:"rating"`
	Pages         int16   `json:"pages"`
	Genre         Genre   `json:"genre"`
	Author        Author  `json:"author"`
}
