package main

import (
	"bytes"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var dbLogger *log.Logger
var db *sqlx.DB

func init() {
	dbLogger = log.New(os.Stdout, "SRS MySQL: ", log.LstdFlags)
}

/**
 * GetConnectionString() - David James
 * Description: Does as named, in a prettier way than one long line of code.
 * Reference: "user:password@tcp(localhost:5555)/dbname?charset=utf8"
**/
func GetDBConnectionString() string {
	var buffer bytes.Buffer

	buffer.WriteString(os.Getenv("MYSQL_USER"))
	buffer.WriteString(":")
	buffer.WriteString(os.Getenv("MYSQL_PASS"))
	buffer.WriteString("@tcp(")
	buffer.WriteString(os.Getenv("MYSQL_HOST"))
	buffer.WriteString(":")
	buffer.WriteString(os.Getenv("MYSQL_PORT"))
	buffer.WriteString(")/")
	buffer.WriteString(os.Getenv("MYSQL_DB"))

	// Any additional MySQL Driver settings should be appended here.
	buffer.WriteString("?charset=utf8")

	return buffer.String()
}

/**
 * ConnectDB() - David James
 * Description: Connects to MySQL Database, and pings to ensure
 *  an open connection pool exists and functions.
**/
func ConnectDB() {
	db = sqlx.MustConnect("mysql", GetDBConnectionString())
	err := db.Ping()
	if err != nil {
		dbLogger.Fatal("Unable to ping MySQL Database.")
	}
	dbLogger.Println("Connected to MySQL Database.")
}
