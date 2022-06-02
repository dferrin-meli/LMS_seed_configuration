package repository

import (
	"database/sql"
	"time"
)

//IDBClient DBClient Interface
type IDBClient interface {
	openDB() *sql.DB
	closeDB(db *sql.DB)
}

//DBClient struct
type DBClient struct {
	connectionString string
}

func (r *DBClient) openDB() *sql.DB {
	db, _ := sql.Open("mysql", r.connectionString)

	db.SetMaxIdleConns(0)
	db.SetConnMaxLifetime(time.Second)
	return db
}

func (r *DBClient) closeDB(db *sql.DB) {
	_ = db.Close()
}
