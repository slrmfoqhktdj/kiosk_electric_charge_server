package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/pat"
	"github.com/joho/godotenv"
	"github.com/urfave/negroni"
)

type ST_GetPay struct {
	Pay int
}

// 가격 반환 db
func get_pay_DB() int {
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		log.Printf(" get_pay_Handler : %s", err)
	}

	var _price int

	err = conn.QueryRow("SELECT price FROM e_price ORDER BY time DESC LIMIT 1").Scan(&_price)

	if err != nil {
		log.Printf(" get_pay_Handler get db : %s", err.Error())
	}

	defer conn.Close()

	return _price
}

// 광고 열람 핸들러
func get_ad_Handler(w http.ResponseWriter, r *http.Request) {

}

// 가격 반환 핸들러
func get_pay_Handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(ST_GetPay{get_pay_DB()})
	fmt.Fprint(w, string(data))
}

func main() {
	err := godotenv.Load("./info.env")
	if err != nil {
		log.Println(err)
		return
	}

	mux := pat.New()

	//negroni를 쓰면 따로 명시하지 않을시 "/"에 접근하면 index.html or default.html 파일에 액세스함
	n := negroni.Classic()

	mux.Path("/api/getpay").HandlerFunc(get_pay_Handler)

	mux.Path("/api/getad").HandlerFunc(get_ad_Handler)

	n.UseHandler(mux)
	http.ListenAndServe(":8080", n)
}
