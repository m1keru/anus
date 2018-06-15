package tools

import (
	"database/sql"
)

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	CheckErr(err)
	if db == nil {
		panic("db nil")
	}
	return db
}

func Migrate(db *sql.DB) {
	sql := `
	CREATE TABLE IF NOT EXISTS ansible_scripts(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name VARCHAR NOT NULL DEFAULT "",
		path VARCHAR(255) NOT NULL UNIQUE,
		description VARCHAR(255) DEFAULT "" 
	);
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		login VARCHAR NOT NULL UNIQUE,
		email VARCHAR DEFAULT "",
		PASSWORD VARCHAR NOT NULL
	);
	`
	_, err := db.Exec(sql)
	CheckErr(err)
}
