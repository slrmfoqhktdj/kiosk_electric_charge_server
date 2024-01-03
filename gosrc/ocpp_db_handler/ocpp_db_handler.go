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
func BootNotification_update(device_name string, device_state string, device_sw_v string) bool {
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Query("UPDATE device SET state = ? ,s_version = ? WHERE id = ?", device_state, device_sw_v, device_name)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("BootNotification_update : %v", rows)

	defer conn.Close()

	return true
}

// 기기 상태 업데이트
func Websocket_connector_update(device_name string, device_state string) bool {
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Query("UPDATE device SET state = ?  WHERE id = ?", device_state, device_name)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("Websocket_connector_update : %v", rows)

	defer conn.Close()

	return true
}

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
func StartTransaction_set_start(c_id string, c_conn int, c_amount int, Timestamp string) bool {

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
func StopTransaction_update_end(id string, c_time_e string, c_amount int) bool {

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

// 충전기 커넥터 상태 update
func StatusNotification_connector(id string, ConnectorId int, state string, errcode string) {

	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		fmt.Println(err)
	}

	if ConnectorId == 1 {
		rows, err := conn.Query("UPDATE device SET l_state = ? WHERE id = ?", state, id)
		if err != nil {
			fmt.Printf("rows1 : %s", err)
		}

		log.Printf("StatusNotification_connector left : %v", rows)
	} else if ConnectorId == 2 {
		rows, err := conn.Query("UPDATE device SET r_state = ? WHERE id = ?", state, id)
		if err != nil {
			fmt.Printf("rows1 : %s", err)
		}

		log.Printf("StatusNotification_connector right : %v", rows)
	}

	defer conn.Close()

}

// 충전기 상태 update
func StatusNotification_device(id string, state string, errcode string) {

	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		fmt.Println(err)
	}

	rows, err := conn.Query("UPDATE device SET state = ? WHERE id = ?", id, state)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
	}

	log.Printf("StatusNotification_device : %v", rows)

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
