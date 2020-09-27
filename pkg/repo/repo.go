package repo

import "github.com/jackc/pgx/v4"

func New(conn *pgx.Conn) *Repo {
	return &Repo{
		conn: conn,
	}
}

type Repo struct {
	conn *pgx.Conn
}
