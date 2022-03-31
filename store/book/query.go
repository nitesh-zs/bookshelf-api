package book

const (
	getByID = `select id, title, author, summary, genre, year, publisher, image_uri from book where id = $1;`
)
