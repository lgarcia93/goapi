package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DatabaseInfo info
type databaseInfo struct {
	Username     string
	Password     string
	DatabaseName string
	Hostname     string
}

func getDatabaseInfo() databaseInfo {
	user := os.Getenv("databaseUser")
	password := os.Getenv("databasePassword")
	hostname := os.Getenv("databaseHost")

	if user == "" {
		user = "root"
	}

	if password == "" {
		password = "root"
	}

	if hostname == "" {
		hostname = "localhost:3306"
	}

	databaseInfo := databaseInfo{
		Hostname:     hostname,
		DatabaseName: "fitapp",
		Username:     user,
		Password:     password,
	}

	return databaseInfo
}

func getConnectionString() string {
	databaseInfo := getDatabaseInfo()

	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true",
		databaseInfo.Username,
		databaseInfo.Password,
		databaseInfo.Hostname,
		databaseInfo.DatabaseName,
	)
}

// Connection holds pointer to DBs
var Connection *sql.DB

// ConfigDatabase setup mysql database
func ConfigDatabase() {

	db, err := sql.Open("mysql", getConnectionString())

	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	Connection = db
}
