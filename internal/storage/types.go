package storage

import 	"database/sql"

type DataBase struct {
	DB            *sql.DB
	MigrationsDir string
}

