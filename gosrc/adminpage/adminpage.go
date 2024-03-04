package adminpage

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"text/template"

	"kiosk_electoric_charge/gosrc/structs"
)

// html 구획 나누기
var templ = func() *template.Template {

	//template 에 function 넣기
	funcadd := template.FuncMap{
		"add": func(a int, b int) int {
			return a + b
		},
	}

	// t := template.New("")
	t := template.New("").Funcs(funcadd)
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			fmt.Println(path)
			_, err = t.ParseFiles(path)
			if err != nil {
				fmt.Println(err)
			}
		}
		return err
	})

	if err != nil {
		fmt.Print(err)
	}
	return t
}()

// 충전기 설치 정보
func device_list_select() []structs.DeviceListItem {
	var result []structs.DeviceListItem
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		log.Printf(" User_survey_a err : %s", err)
	}

	rows, err := conn.Query("select id,latitude,longitude,addr,s_version,h_c_m,h_p_m,h_b_m,state from device")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for rows.Next() {
		var Id string
		var Latitude string
		var Longitude string
		var Addr string
		var S_Version string
		var H_C_M string
		var H_P_M string
		var H_B_M string
		var State string

		rows.Scan(&Id, &Latitude, &Longitude, &Addr, &S_Version, &H_C_M, &H_P_M, &H_B_M, &State)

		m := structs.DeviceListItem{Id: Id, Latitude: Latitude, Longitude: Longitude, Addr: Addr, S_Version: S_Version, H_C_M: H_C_M, H_P_M: H_P_M, H_B_M: H_B_M, State: State}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 충전요금 단가 관리
func price_list_select() []structs.PriceListItem {
	var result []structs.PriceListItem
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		log.Printf(" User_survey_a err : %s", err)
	}

	rows, err := conn.Query("select id,price,time from e_price")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for rows.Next() {
		var Id string
		var Price string
		var Time string

		rows.Scan(&Id, &Price, &Time)

		m := structs.PriceListItem{Id: Id, Price: Price, Time: Time}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 소프트웨어 업데이트 관리
func sw_list_select() []structs.SwListItem {
	var result []structs.SwListItem
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		log.Printf(" User_survey_a err : %s", err)
	}

	rows, err := conn.Query("select id,version,time from app_vs")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for rows.Next() {
		var Id string
		var Version string
		var Time string

		rows.Scan(&Id, &Version, &Time)

		m := structs.SwListItem{Id: Id, Version: Version, Time: Time}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 충전기별 광고 관리
func ad_list_select() []structs.AdListItem {
	var result []structs.AdListItem
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		log.Printf(" User_survey_a err : %s", err)
	}

	rows, err := conn.Query("select d_id,ad_id,time from device_ad")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for rows.Next() {
		var D_Id string
		var AD_Id string
		var Time string

		rows.Scan(&D_Id, &AD_Id, &Time)

		m := structs.AdListItem{D_Id: D_Id, AD_Id: AD_Id, Time: Time}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 충전량 및 이력조회
func charge_list_select() []structs.ChargeListItem {
	var result []structs.ChargeListItem
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		log.Printf(" User_survey_a err : %s", err)
	}

	rows, err := conn.Query("select pay_num,d_id,c_amount,c_pay,pay_by,c_num,pay_day,c_time_s,c_time_e,c_time from pay_h")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for rows.Next() {
		var Pay_Num string
		var D_Id string
		var C_Amount string
		var C_Pay string
		var Pay_By string
		var C_Num string
		var Pay_Day string
		var C_Time_S string
		var C_Time_E string
		var C_Time string

		rows.Scan(&Pay_Num, &D_Id, &C_Amount, &C_Pay, &Pay_By, &C_Num, &Pay_Day, &C_Time_S, &C_Time_E, &C_Time)

		// log.Printf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s", Pay_Num, D_Id, C_Amount, C_Pay, Pay_By, C_Num, Pay_Day, C_Time_S, C_Time_E, C_Time)

		m := structs.ChargeListItem{Pay_Num: Pay_Num, D_Id: D_Id, C_Amount: C_Amount, C_Pay: C_Pay, Pay_By: Pay_By, C_Num: C_Num, Pay_Day: Pay_Day, C_Time_S: C_Time_S, C_Time_E: C_Time_E, C_Time: C_Time}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	// log.Printf(" charge list : %s", result)

	return result
}

// 충전기 장애 이력
func err_list_select() []structs.ErrListItem {
	var result []structs.ErrListItem
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		log.Printf(" User_survey_a err : %s", err)
	}

	rows, err := conn.Query("select d_id,o_time,cause,r_time from device_err")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for rows.Next() {
		var D_Id string
		var O_Time string
		var Cause string
		var R_Time string

		rows.Scan(&D_Id, &O_Time, &Cause, &R_Time)

		m := structs.ErrListItem{D_Id: D_Id, O_Time: O_Time, Cause: Cause, R_Time: R_Time}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	// log.Printf(" charge list : %s", result)

	return result
}

func LoginPage_Handler(w http.ResponseWriter, r *http.Request) {
	templ.ExecuteTemplate(w, "loginpage", &structs.Page{Title: "LoginPage"})
}

func Login_Handler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/devicelist", http.StatusTemporaryRedirect)

}

// 충전기 설치 정보
func Devicelist_Page_Handler(w http.ResponseWriter, r *http.Request) {
	templ.ExecuteTemplate(w, "main_menu", &structs.DeviceListPage{Title: "기기 목록", Page: 0, Items: device_list_select()})
}

// 충전요금 단가 관리
func Pricemng_Page_Handler(w http.ResponseWriter, r *http.Request) {
	templ.ExecuteTemplate(w, "main_menu", &structs.PricemngPage{Title: "요금 단가 관리", Page: 2, Items: price_list_select()})
}

// 소프트웨어 업데이트 관리
func Swmng_Page_Handler(w http.ResponseWriter, r *http.Request) {
	templ.ExecuteTemplate(w, "main_menu", &structs.SwmngPage{Title: "SW 업데이트 관리", Page: 4, Items: sw_list_select()})
}

//광고물 관리

// 충전기별 광고 관리
func Admng_Page_Handler(w http.ResponseWriter, r *http.Request) {
	templ.ExecuteTemplate(w, "main_menu", &structs.AdmngPage{Title: "충전기별 광고 관리", Page: 3, Items: ad_list_select()})
}

// 충전량 및 이력조회
func Chargelist_Page_Handler(w http.ResponseWriter, r *http.Request) {
	templ.ExecuteTemplate(w, "main_menu", &structs.ChargeListPage{Title: "충전량 및 이력조회", Page: 1, Items: charge_list_select()})
}

//일자별 충전량 조회
//충전기 그래프 조회

// 충전기 장애 이력
func Errlist_Page_Handler(w http.ResponseWriter, r *http.Request) {
	templ.ExecuteTemplate(w, "main_menu", &structs.ErrListPage{Title: "기기 장애 이력", Page: 6, Items: err_list_select()})
}

//매출실적 및 정산 조회
//사용자 관리 권한
//공통코드

///////////////////////////////////////////////////////////			// 기능 핸들러

// 기기 등록
func insert_device(device_name string, device_state string, device_latitude string, device_longitude string, device_addr string, device_sw_v string, device_charger_m string, device_pay_m string, device_hw_m string) {
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rows, err := conn.Query("INSERT INTO device(id,state,latitude,longitude,addr,s_version,h_c_m,h_p_m,h_b_m) VALUES(?,?,?,?,?,?,?,?,?)", device_name, device_state, device_latitude, device_longitude, device_addr, device_sw_v, device_charger_m, device_pay_m, device_hw_m)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		os.Exit(1)
	}

	rows2, err := conn.Query("INSERT INTO device(id) VALUES(?)", device_name)
	if err != nil {
		fmt.Printf("rows2 : %s", err)
		os.Exit(1)
	}

	log.Printf("insert_device : %v , %v", rows, rows2)

	defer conn.Close()
}

// 기기 수정
func update_device(device_name string, device_state string, device_latitude string, device_longitude string, device_addr string, device_sw_v string, device_charger_m string, device_pay_m string, device_hw_m string) {
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rows, err := conn.Query("UPDATE device SET state = ?,latitude = ?,longitude = ?,addr = ?,s_version = ?,h_c_m = ?,h_p_m = ?,h_b_m = ? WHERE id = ?", device_state, device_latitude, device_longitude, device_addr, device_sw_v, device_charger_m, device_pay_m, device_hw_m, device_name)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		os.Exit(1)
	}

	log.Printf("update_device : %v", rows)

	defer conn.Close()
}

// 기기 등록
func Register_Device_Handler(w http.ResponseWriter, r *http.Request) {
	device_name := r.PostFormValue("device-name")
	device_state := r.PostFormValue("device-state")
	device_latitude := r.PostFormValue("device-latitude")
	device_longitude := r.PostFormValue("device-longitude")
	device_addr := r.PostFormValue("device-addr")
	device_sw_v := r.PostFormValue("device-sw-v")
	device_charger_m := r.PostFormValue("device-charger-m")
	device_pay_m := r.PostFormValue("device-pay-m")
	device_hw_m := r.PostFormValue("device-hw-m")

	log.Printf("value %s,%s,%s,%s,%s,%s,%s,%s,%s", device_name, device_state, device_latitude, device_longitude, device_addr, device_sw_v, device_charger_m, device_pay_m, device_hw_m)

	insert_device(device_name, device_state, device_latitude, device_longitude, device_addr, device_sw_v, device_charger_m, device_pay_m, device_hw_m)

	http.Redirect(w, r, "/devicelist", http.StatusMovedPermanently)
}

// 기기 수정
func Updater_Device_Handler(w http.ResponseWriter, r *http.Request) {
	device_name := r.PostFormValue("device-name")
	device_state := r.PostFormValue("device-state")
	device_latitude := r.PostFormValue("device-latitude")
	device_longitude := r.PostFormValue("device-longitude")
	device_addr := r.PostFormValue("device-addr")
	device_sw_v := r.PostFormValue("device-sw-v")
	device_charger_m := r.PostFormValue("device-charger-m")
	device_pay_m := r.PostFormValue("device-pay-m")
	device_hw_m := r.PostFormValue("device-hw-m")

	log.Printf("value %s,%s,%s,%s,%s,%s,%s,%s,%s", device_name, device_state, device_latitude, device_longitude, device_addr, device_sw_v, device_charger_m, device_pay_m, device_hw_m)

	update_device(device_name, device_state, device_latitude, device_longitude, device_addr, device_sw_v, device_charger_m, device_pay_m, device_hw_m)

	http.Redirect(w, r, "/devicelist", http.StatusMovedPermanently)
}

// 기기 조회
func Selecter_Device_Handler(w http.ResponseWriter, r *http.Request) {
	device_name := r.PostFormValue("device-name-select")
	device_state := r.PostFormValue("device-state-select")
	device_addr := r.PostFormValue("device-addr-select")

	log.Printf("value %s,%s,%s", device_name, device_state, device_addr)

	http.Redirect(w, r, "/devicelist", http.StatusMovedPermanently)
}
