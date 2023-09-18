package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/k0kubun/pp/v3"
)

var dbWithoutParseTime *sql.DB
var dbWithParseTime *sql.DB

type TimeTypesWithoutParseTime struct {
	ID        int64
	Date      string
	DateTime  string
	Timestamp string
	Year      int
	Time      string
}

type TimeTypesWithParseTime struct {
	ID        int64
	Date      time.Time
	DateTime  time.Time
	Timestamp time.Time
	Year      int
	Time      string
}

func main() {
	err := InitDBs()
	if err != nil {
		panic(err)
	}
	defer CloseDBs()

	var timeTypesWithoutParseTime TimeTypesWithoutParseTime
	if err := dbWithoutParseTime.
		QueryRow("SELECT * FROM time_types WHERE id = ?", 1).
		Scan(
			&timeTypesWithoutParseTime.ID,
			&timeTypesWithoutParseTime.Date,
			&timeTypesWithoutParseTime.DateTime,
			&timeTypesWithoutParseTime.Timestamp,
			&timeTypesWithoutParseTime.Year,
			&timeTypesWithoutParseTime.Time,
		); err != nil {
		panic(err)
	}
	pp.Println(timeTypesWithoutParseTime)

	var timeTypesWithParseTime TimeTypesWithParseTime
	if err := dbWithParseTime.
		QueryRow("SELECT * FROM time_types WHERE id = ?", 1).
		Scan(
			&timeTypesWithParseTime.ID,
			&timeTypesWithParseTime.Date,
			&timeTypesWithParseTime.DateTime,
			&timeTypesWithParseTime.Timestamp,
			&timeTypesWithParseTime.Year,
			&timeTypesWithParseTime.Time,
		); err != nil {
		panic(err)
	}
	pp.Println(timeTypesWithParseTime)
}

func InitDBs() error {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	dbWithoutParseTime, err = sql.Open("mysql", (&mysql.Config{
		DBName:    "test",
		User:      "root",
		Passwd:    "password",
		Addr:      "localhost:3306",
		Net:       "tcp",
		Collation: "utf8mb4_bin",
		ParseTime: false,
		Loc:       jst,
	}).FormatDSN())
	if err != nil {
		return err
	}

	dbWithParseTime, err = sql.Open("mysql", (&mysql.Config{
		DBName:    "test",
		User:      "root",
		Passwd:    "password",
		Addr:      "localhost:3306",
		Net:       "tcp",
		Collation: "utf8mb4_bin",
		ParseTime: true,
		Loc:       jst,
	}).FormatDSN())
	if err != nil {
		return err
	}

	if err := dbWithoutParseTime.Ping(); err != nil {
		return err
	}

	if err := dbWithParseTime.Ping(); err != nil {
		return err
	}

	return nil
}

func CloseDBs() {
	if dbWithoutParseTime != nil {
		if err := dbWithoutParseTime.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}
	if dbWithParseTime != nil {
		if err := dbWithParseTime.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}
}
