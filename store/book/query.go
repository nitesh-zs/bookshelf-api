package book

const (
	getByID    = `select title, author, summary, genre, year, publisher, image_uri from book where id = $1;`
	getFilters = `select distinct $1 from book;`
)
