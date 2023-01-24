package initializers

import (
	"database/sql"
	"log"
)

func InitializingDB() *sql.DB {

	db, err := sql.Open("sqlite", "todo.db")
	if err != nil {
		panic(err)
	}
	sql := `create table if not exists todo (id integer primary key autoincrement,title text,body text,iscomplete Boolean DEFAULT 0);`
	_, err = db.Exec(sql)
	if err != nil {
		log.Printf(" %q %s", err, sql)
		panic(err)
	}
	return db

}
