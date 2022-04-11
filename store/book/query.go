package book

const (
	getByID = `select title, author, summary, genre, year, publisher, image_uri from book where id = $1;`

	createBook = `insert into "book" (id, title, author, summary, genre, year, reg_num, publisher, language,
					image_uri) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	updateBook = `update "book" set title = $1,author = $2,summary = $3,genre = $4,year = $5,
					reg_num = $6,publisher = $7,language = $8,image_uri = $9 where id = $10;`
)
