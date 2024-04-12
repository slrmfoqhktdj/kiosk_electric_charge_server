package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func OpenDB() (*sql.DB, error) {
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		log.Printf("DB conn err : %v", err)
		return conn, err
	}

	conn.SetMaxIdleConns(20)                  //idle connection pool 의 최대 연결 수를 설정
	conn.SetMaxOpenConns(20)                  //데이터베이스에 열린 연결의 최대 수가 설정
	conn.SetConnMaxLifetime(time.Minute * 10) //연결을 재사용할 수 있는 최대 시간을 설정

	return conn, nil
}
