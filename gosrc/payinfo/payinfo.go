package payinfo

import (
	"fmt"
	"log"
	"net/http"
	"os"

	// _ "github.com/go-sql-driver/mysql"

	"kiosk_electoric_charge/gosrc/db"
)

// 결제정보 insert
func set_pay(device_id string, charge_amount string, charge_Pay string, pay_By string, credit_num string, pay_num string, pay_day string) {
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rows, err := conn.Query("INSERT INTO pay_h(pay_num,d_id,c_amount,c_pay,pay_by,c_num,pay_day) VALUES(?,?,?,?,?,?,?)", pay_num, device_id, charge_amount, charge_Pay, pay_By, credit_num, pay_day)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		os.Exit(1)
	}

	log.Printf("set_pay : %v", rows)

	defer conn.Close()
}

// 충전시작 update
func update_c_start(pay_num string, time string) {
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rows, err := conn.Query("UPDATE pay_h SET c_time_s = ? where pay_num = ?", time, pay_num)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		os.Exit(1)
	}

	log.Printf("update_c_start : %v", rows)

	defer conn.Close()
}

// 충전완료 update
func update_c_end(pay_num string, time string) {
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rows, err := conn.Query("UPDATE pay_h SET c_time_e = ? where pay_num = ?", time, pay_num)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		os.Exit(1)
	}

	log.Printf("update_c_end : %v", rows)

	defer conn.Close()
}

// 결제 정보
func Set_Pay_Handler(w http.ResponseWriter, r *http.Request) {

	device_id := r.PostFormValue("Device_id")
	charge_amount := r.PostFormValue("Charge_amount")
	charge_Pay := r.PostFormValue("Charge_Pay")
	pay_By := r.PostFormValue("Pay_By")
	credit_num := r.PostFormValue("Credit_num")
	pay_num := r.PostFormValue("Pay_num")
	pay_day := r.PostFormValue("Pay_day")

	log.Printf("set pay handler %s , %s ,%s ,%s ,%s ,%s ,%s ", device_id, charge_amount, charge_Pay, pay_By, credit_num, pay_num, pay_day)

	set_pay(device_id, charge_amount, charge_Pay, pay_By, credit_num, pay_num, pay_day)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, http.StatusOK)
}

// 충전 시작
func Set_charging_start_Handler(w http.ResponseWriter, r *http.Request) {

	pay_num := r.PostFormValue("Pay_num")
	start_time := r.PostFormValue("C_time_s")

	update_c_start(pay_num, start_time)

	log.Printf("Set_charging_start_Handler %s , %s ", pay_num, start_time)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, http.StatusOK)
}

// 충전 완료
func Set_charging_end_Handler(w http.ResponseWriter, r *http.Request) {

	pay_num := r.PostFormValue("Pay_num")
	end_time := r.PostFormValue("C_time_e")

	update_c_end(pay_num, end_time)

	log.Printf("Set_charging_start_Handler %s , %s ", pay_num, end_time)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, http.StatusOK)
}
