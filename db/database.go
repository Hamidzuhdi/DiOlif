package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	// Ganti dengan kredensial database Anda
	DB, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/konveksi-app")
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to database")
}
