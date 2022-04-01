package user

const (
	getUserByEmail = `select * from "user" where email = $1`
	createUser     = `insert into "user"(id, email, name, type) values($1, $2, $3, $4)`
)
