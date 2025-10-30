package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

var DB *sql.DB

// Connect initializes and verifies the MySQL connection
func Connect() {
	var err error

	// credentials

	username := "root"
	password := "Born2001$"
	host := "127.0.0.1"
	port := "3306"
	dbname := "hardware_shop"

	// Building a  connection 

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, dbname)

	// Open a connection

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf(" Error opening MySQL connection: %v", err)
	}

	// Test the connection

	err = DB.Ping()
	if err != nil {
		log.Fatalf(" Cannot connect to MySQL: %v", err)
	}

	log.Println(" Connected to MySQL database successfully!")
}
