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

	"kiosk_electoric_charge/gosrc/adminpage"
	"kiosk_electoric_charge/gosrc/updateapk"
)

//ST = struct

type ST_AD struct {
	ADId string
}

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

//광고 리스트 반환 db
func get_ADlist_DB(device_id string) []ST_AD {
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	var r []ST_AD

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		log.Printf(" get_pay_Handler : %s", err)
	}

	rows, err := conn.Query("SELECT ad_id FROM device_ad WHERE d_id = ?", device_id)

	if err != nil {
		fmt.Printf("get_ADlist_DB func : %s", err)
	}

	for rows.Next() {
		var ad_id string
		rows.Scan(&ad_id)

		s := ST_AD{ad_id}
		r = append(r, s)
	}

	defer rows.Close()
	defer conn.Close()

	return r

}

//////////////////////////////////////////////////////////////////////////////////////////////////

// 광고 다운 핸들러
func access_ad_Handler(w http.ResponseWriter, r *http.Request) {

	a_id := r.PostFormValue("ad_id")
	log.Printf("access %s", a_id)

	http.ServeFile(w, r, "./advertisement/"+a_id)
}

// 광고 리스트 반환 핸들러
func get_ad_Handler(w http.ResponseWriter, r *http.Request) {
	a_id := r.PostFormValue("device_id")

	log.Printf("id = %s", a_id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(get_ADlist_DB(a_id))
	fmt.Fprint(w, string(data))

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

	//////////////////////////////////////////////////////////////////////////////////////////////

	//관리자 페이지
	//로그인
	mux.Path("/").HandlerFunc(adminpage.Login_Handler)

	//////////////////////////////////////////////////////////////////////////////////////////////

	//앱 업데이트

	//버전 확인
	mux.Path("/api/apk/version").HandlerFunc(updateapk.Get_apk_version_Handler).Methods("POST")

	//앱 apk 파일 다운
	mux.Path("/file/apk/operation/1").HandlerFunc(updateapk.App_apk_download_Handler).Methods("POST")

	//앱 업데이트 완료
	mux.Path("/file/apk/consequnce").HandlerFunc(updateapk.App_update_complete_Handler).Methods("POST")

	//////////////////////////////////////////////////////////////////////////////////////////////

	//광고 리스트
	mux.Path("/api/ad").HandlerFunc(get_ad_Handler).Methods("POST")

	//광고 접근 후 다운
	//경로끝에 /를 붙여줘야 html상에서 보임
	mux.Path("/file/ad").HandlerFunc(access_ad_Handler).Methods("POST")

	//////////////////////////////////////////////////////////////////////////////////////////////

	//가격 리턴
	mux.Path("/api/pay").HandlerFunc(get_pay_Handler).Methods("POST")

	//////////////////////////////////////////////////////////////////////////////////////////////

	n.UseHandler(mux)
	http.ListenAndServe(":8080", n)
}
