package driver

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	// database/sql 只需要和stdlib交互就可以
	// _ "github.com/jackx/pgconn"
	// _ "github.com/jackx/pgx/v5"
	"log"
	"time"
)

// dirver.go 调用ConnextSQL返回一个DB结构体，持有postgresql的db

// DB holds the database connection pool.
// is a datastruct
type DB struct {
	SQL *sql.DB
}

// dbConn is a DB
var dbConn = &DB{}

// some set
const maxOpenConns = 10
const maxIdleConns = 5
const connMaxLifetime = 5 * time.Minute // 5 minutes

// ConnectSQL holds connect set return a driver.DB
func ConnectSQL(dsn string) (*DB, error) {
	db, err := NewDatabase(dsn)
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	err = testDB(db)
	if err != nil {
		return nil, err
	}

	dbConn.SQL = db
	return dbConn, nil
}

// testDB receive a *sql.DB and test the db connection
func testDB(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}

// NewDatabase receive a dsn and retrun a sql.DB and error
func NewDatabase(dsn string) (*sql.DB, error) {

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
