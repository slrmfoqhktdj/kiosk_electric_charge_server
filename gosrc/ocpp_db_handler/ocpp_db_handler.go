package ocpp_db_handler

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// 장치 상태 1 정상 , -1 에러,  2 충전중 , 3 연결끊김

// 기기 부팅
func BootNotification_db_handler(device_name string, device_sw_v string) string {
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	var enabled_row string

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		fmt.Println(err)
		return "2"
	}

	err = conn.QueryRow("SELECT enabled FROM device WHERE id = ?", device_name).Scan(&enabled_row)
	if err != nil {
		fmt.Printf("srows1 : %s", err)
		return "2"
	}

	rows, err := conn.Query("UPDATE device SET s_version = ? WHERE id = ?", device_sw_v, device_name)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return "2"
	}

	log.Printf("BootNotification_db_handler : %s , %v", enabled_row, rows)

	defer conn.Close()

	return enabled_row
}

// 기기 상태 업데이트	--- 에러,충전중,정상 써야할듯. 연결, 비연결은 관리자페이지에서 불러올때 main.go에서 csms에 있는 connector 모음집에서 꺼내와 확인해야할듯
// func Websocket_connector_update(device_name string, device_state string) bool {
// 	dbname := os.Getenv("DBNAME")
// 	dbdata := os.Getenv("DBDATA")

// 	conn, err := sql.Open(dbname, dbdata)
// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}

// 	rows, err := conn.Query("UPDATE device SET state = ?  WHERE id = ?", device_state, device_name)
// 	if err != nil {
// 		fmt.Printf("rows1 : %s", err)
// 		return false
// 	}

// 	log.Printf("Websocket_connector_update : %v", rows)

// 	defer conn.Close()

// 	return true
// }

// 결제정보 insert
func DataTransfer_set_pay(device_id string, charge_Pay string, pay_By string, credit_num string, pay_num string, pay_day string) bool {
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Query("INSERT INTO pay_h(pay_num,d_id,c_pay,pay_by,c_num,pay_day) VALUES(?,?,?,?,?,?)", pay_num, device_id, charge_Pay, pay_By, credit_num, pay_day)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("DataTransfer_set_pay : %v", rows)

	defer conn.Close()

	return true
}

// 결제정보 부분취소 update
func DataTransfer_update_pay(can_pay string, c_pay_day string, c_c_num string, pay_num string) bool {
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Query("UPDATE pay_h SET can_pay = ? , c_pay_day = ? , c_c_num = ? WHERE pay_num = ?) VALUES(?,?,?,?)", can_pay, c_pay_day, c_c_num, pay_num)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("DataTransfer_update_pay : %v", rows)

	defer conn.Close()

	return true
}

// 결제 여부 확인
func Authorize_select_payuser(pay_num string) bool {
	var exist bool

	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows := conn.QueryRow("SELECT EXISTS(SELECT pay_num FROM pay_h WHERE pay_num = ?)", pay_num)

	rows.Scan(&exist)

	log.Printf("Authorize_select_payuser : %v", rows)

	defer conn.Close()

	log.Printf("aa : %v", exist)

	return exist

}

// 충전정보 insert
func StartTransaction_set_start(c_id string, c_conn int, c_amount float32, Timestamp string) bool {

	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Query("INSERT INTO charge_h(c_id,c_conn,c_amount,c_time_s) VALUES(?,?,?,?)", c_id, c_conn, c_amount, Timestamp)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("StartTransaction_set_start : %v", rows)

	defer conn.Close()

	return true
}

// 충전정보 update
func StopTransaction_update_end(id string, c_time_e string, c_amount float32) bool {

	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Query("UPDATE charge_h SET c_amount = ? ,c_time_e = ? WHERE c_id = ?", c_amount, c_time_e, id)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("StopTransaction_update_end : %v", rows)

	defer conn.Close()

	return true
}

func StatusNotification_update(id string, ConnectorId int, errcode string, status string, vendorid string, vendorerrorcode string, info string, timestemp string) {

	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		fmt.Println(err)
	}

	if vendorid > "" {
		rows, err := conn.Query("UPDATE device_status SET state = ? , state_e = ?, v_id = ?, v_error = ? , e_info = ? , u_time = ? WHERE id = ?", status, errcode, vendorid, vendorerrorcode, info, timestemp, id)
		if err != nil {
			fmt.Printf("rows1 : %s", err)
		}
		log.Printf("StatusNotification_update vendor false : %v", rows)

	} else {
		rows, err := conn.Query("UPDATE device_status SET state = ? , state_e = ?, e_info = ? , u_time = ?  WHERE id = ?", status, errcode, info, timestemp, id)
		if err != nil {
			fmt.Printf("rows1 : %s", err)
		}
		log.Printf("StatusNotification_update vendor true : %v", rows)
	}

	defer conn.Close()

}

// 충전량 update
func Metervalue_update(c_id string, t_id int, time string, value string) {

	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		fmt.Println(err)
	}

	rows, err := conn.Query("UPDATE charge_h SET c_time_u = ?, c_amount = ? WHERE c_id = ? and c_conn = ?", time, value, c_id, t_id)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
	}

	log.Printf("Metervalue_update : %v", rows)

	defer conn.Close()
}
