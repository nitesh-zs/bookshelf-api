package book

const (
	createBook = `insert into "book"(id, title, author, summary, genre, year, reg_num, publisher, language, image_uri) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
)
