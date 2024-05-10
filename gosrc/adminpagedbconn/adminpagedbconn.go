package adminpagedbconn

import (
	"fmt"
	"log"
	"strings"

	"kiosk_electoric_charge/gosrc/adminjwt"
	"kiosk_electoric_charge/gosrc/db"
	"kiosk_electoric_charge/gosrc/structs"

	"golang.org/x/crypto/bcrypt"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// 로그인
func Login_select(email string, pwd string) structs.LoginDataForm {
	var result structs.LoginDataForm
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	var Email string
	var Name string
	var Password string

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" Login_select db err : %s", err)
	}

	// // 암호화
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), 11) // 10 means difficulty
	// if err != nil {
	// 	log.Printf("failed to GenerateFromPassword: %v", err)
	// 	return result
	// }

	// log.Printf("password : %v", string(hashedPassword))

	if err := conn.QueryRow("SELECT email,name,password FROM admin_user WHERE email = ?", email).Scan(&Email, &Name, &Password); err != nil {
		log.Printf("select user, %v", err)
		return result
	}

	if err := bcrypt.CompareHashAndPassword([]byte(Password), []byte(pwd)); err != nil {
		log.Printf("incorrect password, %v", err)
		return result
	}

	tokenString, err := adminjwt.GenerateToken(Name)
	if err != nil {
		return result
	}

	log.Printf("key is %v", tokenString)

	result = structs.LoginDataForm{Email: Email, Name: tokenString}

	defer conn.Close()

	return result
}

// 충전기 위치 정보
func Map_list_select() []structs.MapListItem {
	var result []structs.MapListItem
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Map_list_select err : %s", err)
	}

	rows, err := conn.Query("SELECT id,latitude,longitude,addr FROM device")
	if err != nil {
		log.Printf("Map_list_select err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Latitude string
		var Longitude string
		var Addr string

		rows.Scan(&Id, &Latitude, &Longitude, &Addr)

		d := structs.MapListPositionItem{Latitude: Latitude, Longitude: Longitude}

		m := structs.MapListItem{Id: Id, Position: d, Addr: Addr}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 충전기 설치 정보
func Device_list_select(s_order, e_order string) []structs.DeviceListItem {
	var result []structs.DeviceListItem
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" Device_list_select err : %s", err)
	}

	const sqlString = "SELECT * FROM device ORDER BY i_time DESC limit ?,?"

	rows, err := conn.Query(sqlString, s_order, e_order)
	if err != nil {
		log.Printf(" Device_list_select err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Latitude string
		var Longitude string
		var Addr string
		var Model *string
		var S_Version *string
		var H_C_M string
		var H_P_M string
		var H_B_M string
		var H_W_M string
		var I_Time string
		var U_Time string
		var Enabled string

		rows.Scan(&Id, &Latitude, &Longitude, &Addr, &Model, &S_Version, &H_C_M, &H_P_M, &H_B_M, &H_W_M, &I_Time, &U_Time, &Enabled)

		m := structs.DeviceListItem{Id: Id, Latitude: Latitude, Longitude: Longitude, Addr: Addr, Model: Model, S_Version: S_Version, H_C_M: H_C_M, H_P_M: H_P_M, H_B_M: H_B_M, H_W_M: H_W_M, I_Time: I_Time, U_Time: U_Time, Enabled: Enabled}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 충전기 상태 정보
func Device_state_list_select(dconn []string, s_order, e_order string) []structs.DeviceStateListItem {
	var result []structs.DeviceStateListItem
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" device_state_list_select err : %s", err)
	}

	const sqlString = "SELECT * FROM device_status ORDER BY u_time DESC limit ?,?"

	rows, err := conn.Query(sqlString, s_order, e_order)
	if err != nil {
		log.Printf(" device_state_list_select err : %s", err)
	}

	for rows.Next() {

		var Id string
		// var Conn    string
		var State *string
		var State_e *string
		var V_Id *string
		var V_Error *string
		var E_Info *string
		var U_Time *string
		var F_State *string
		var D_State *string

		rows.Scan(&Id, &State, &State_e, &V_Id, &V_Error, &E_Info, &U_Time, &F_State, &D_State)

		m := structs.DeviceStateListItem{Id: Id, Conn: contains(dconn, Id), State: State, State_e: State_e, V_Id: V_Id, V_Error: V_Error, E_Info: E_Info, U_Time: U_Time, F_State: F_State, D_State: D_State}
		result = append(result, m)
	}

	defer rows.Close()
	defer conn.Close()

	return result
}

// 충전요금 단가 관리
func Price_list_select(s_order, e_order string) []structs.PriceListItem {
	var result []structs.PriceListItem
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" Price_list_select err : %s", err)
	}

	const sqlString = "SELECT * FROM e_price ORDER BY time DESC limit ?,?"

	rows, err := conn.Query(sqlString, s_order, e_order)
	if err != nil {
		log.Printf(" Price_list_select err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Price string
		var Time string
		var Writer string

		rows.Scan(&Id, &Price, &Time, &Writer)

		m := structs.PriceListItem{Id: Id, Price: Price, Time: Time, Writer: Writer}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 충전기 업데이트 리스트 정보
func Update_list_select(s_order, e_order string) []structs.DeviceUpdateListItem {
	var result []structs.DeviceUpdateListItem
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("update_list_select err : %s", err)
	}

	rows, err := conn.Query("SELECT d.id, d.addr,d.model,d.s_version,s.f_state,d.i_time,d.u_time  FROM device AS d LEFT JOIN device_status AS s ON s.id = d.id ORDER BY s.u_time DESC limit ?,?", s_order, e_order)
	if err != nil {
		log.Printf("update_list_select err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Addr string
		var Model *string
		var S_Version *string
		var F_State *string
		var I_Time string
		var U_Time string

		rows.Scan(&Id, &Addr, &Model, &S_Version, &F_State, &I_Time, &U_Time)

		m := structs.DeviceUpdateListItem{Id: Id, Addr: Addr, Model: Model, S_Version: S_Version, F_State: F_State, I_Time: I_Time, U_Time: U_Time}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 충전기 로그 리스트 정보
func Log_list_select(s_order, e_order string) []structs.DeviceLogListItem {
	var result []structs.DeviceLogListItem
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Log_list_select err : %s", err)
	}

	rows, err := conn.Query("SELECT d.id, d.addr,s.d_state  FROM device AS d LEFT JOIN device_status AS s ON s.id = d.id ORDER BY s.d_state DESC  limit ?,?", s_order, e_order)
	if err != nil {
		log.Printf("Log_list_select err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Addr string
		var D_State *string

		rows.Scan(&Id, &Addr, &D_State)

		m := structs.DeviceLogListItem{Id: Id, Addr: Addr, D_State: D_State}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 소프트웨어 업데이트 관리
// func sw_list_select() []structs.SwListItem {
// 	var result []structs.SwListItem
// 	// dbname := os.Getenv("DBNAME")
// 	// dbdata := os.Getenv("DBDATA")

// 	// conn, err := sql.Open(dbname, dbdata)

// 	conn, err := db.OpenDB()
// 	if err != nil {
// 		log.Printf(" User_survey_a err : %s", err)
// 	}

// 	rows, err := conn.Query("select id,version,time from app_vs")
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	for rows.Next() {
// 		var Id string
// 		var Version string
// 		var Time string

// 		rows.Scan(&Id, &Version, &Time)

// 		m := structs.SwListItem{Id: Id, Version: Version, Time: Time}
// 		result = append(result, m)
// 	}
// 	defer rows.Close()
// 	defer conn.Close()

// 	return result
// }

// 충전기별 광고 관리
func Ad_list_select() []structs.AdListItem {
	var result []structs.AdListItem
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("ad_list_select err : %s", err)
	}

	rows, err := conn.Query("select d_id,ad_id,time from device_ad")
	if err != nil {
		log.Printf("ad_list_select err : %s", err)
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
func Charge_list_select(s_order, e_order string) []structs.ChargeListItem {
	var result []structs.ChargeListItem
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("charge_list_select err : %s", err)
	}

	rows, err := conn.Query("SELECT p.pay_num as payid, d_id,c_pay,can_pay,pay_by,c_num,c_c_num,pay_day,c_pay_day,c_conn,c_time,c_amount,c_time_s,c_time_e,c_time_u,c_reason FROM charge_h as c RIGHT JOIN pay_h as p  ON c.c_id = p.pay_num ORDER BY c_time_s DESC limit ?,?", s_order, e_order)
	if err != nil {
		log.Printf("charge_list_select err : %s", err)
	}

	for rows.Next() {
		var Payid string
		var D_id *string
		var C_pay *string
		var Can_pay *string
		var Pay_by *string
		var C_num *string
		var C_c_num *string
		var Pay_day *string
		var C_pay_day *string
		var C_conn *string
		var C_time *string
		var C_amount *string
		var C_time_s *string
		var C_time_e *string
		var C_time_u *string
		var C_reason *string

		rows.Scan(&Payid, &D_id, &C_pay, &Can_pay, &Pay_by, &C_num, &C_c_num, &Pay_day, &C_pay_day, &C_conn, &C_time, &C_amount, &C_time_s, &C_time_e, &C_time_u, &C_reason)

		m := structs.ChargeListItem{Payid: Payid, D_id: D_id, C_pay: C_pay, Can_pay: Can_pay, Pay_by: Pay_by, C_num: C_num, C_c_num: C_c_num, Pay_day: Pay_day, C_pay_day: C_pay_day, C_conn: C_conn, C_time: C_time, C_amount: C_amount, C_time_s: C_time_s, C_time_e: C_time_e, C_time_u: C_time_u, C_reason: C_reason}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	// log.Printf(" charge list : %s", result)

	return result
}

// 일자별 충전량 조회
func Charge_daily_list_select(s_order, e_order string) []structs.ChargeDailyListItem {
	var result []structs.ChargeDailyListItem
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Charge_daily_list_select err : %s", err)
	}

	const sqlString = "SELECT CAST(DATE_ADD(CAST(c.c_time_e as DATETIME),INTERVAL 9 HOUR) as DATE) as day, sum(CAST(c.c_amount as FLOAT)) as amount ,COUNT(c.c_amount) as count, sum(CAST(p.c_pay as UNSIGNED)-CAST(IFNULL(p.can_pay,0) as UNSIGNED)) as pay FROM charge_h as c RIGHT JOIN pay_h as p  ON (c.c_id = p.pay_num) GROUP BY day ORDER BY day DESC limit ?,?"

	rows, err := conn.Query(sqlString, s_order, e_order)
	if err != nil {
		log.Printf("Charge_daily_list_select err : %s", err)
	}

	for rows.Next() {
		var Day string
		var Amount *string
		var Count *string
		var Pay *string

		rows.Scan(&Day, &Amount, &Count, &Pay)

		m := structs.ChargeDailyListItem{Day: Day, Amount: Amount, Count: Count, Pay: Pay}
		if Day != "" {
			result = append(result, m)
		}

	}
	defer rows.Close()
	defer conn.Close()

	// log.Printf(" charge list : %s", result)

	return result
}

// 충전기 그래프 조회
func Charge_daily_graph_list_select() []structs.ChargeDailyGraphListItem {
	var result []structs.ChargeDailyGraphListItem
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("charge_daily_graph_list_select err : %s", err)
	}

	rows, err := conn.Query("SELECT id FROM device")
	if err != nil {
		log.Printf("charge_daily_graph_list_select err : %s", err)
	}

	for rows.Next() {
		var Id string

		rows.Scan(&Id)

		m := structs.ChargeDailyGraphListItem{Id: Id}

		result = append(result, m)

	}
	defer rows.Close()
	defer conn.Close()

	// log.Printf(" charge list : %s", result)

	return result
}

// 충전기 그래프 조회
func Charge_daily_graph_select(starttime, endtime, device_id string) []structs.ChargeDailyGraphItem {
	var result []structs.ChargeDailyGraphItem
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Charge_daily_graph_select err : %s", err)
	}

	//검색하는 기기 이름이 전체 or 개별
	if device_id == "전체" {

		const sqlString = "SELECT A.n AS time , ifnull(B.amount,0) AS amount, ifnull(B.count,0) AS count,ifnull(B.pay,0) AS pay  FROM (SELECT @N := @N +1 AS n FROM charge_h, (select @N:=-1 from DUAL ) NN LIMIT 24) AS A LEFT JOIN (SELECT Hour(DATE_ADD(CAST(c.c_time_e as DATETIME),INTERVAL 9 HOUR)) as hh, sum(CAST(c.c_amount as FLOAT)) as amount ,COUNT(c.c_amount) as count ,sum(CAST(p.c_pay as UNSIGNED)-CAST(IFNULL(p.can_pay,0) as UNSIGNED)) as pay FROM charge_h as c RIGHT JOIN pay_h as p  ON c.c_id = p.pay_num WHERE DATE_ADD(CAST(c.c_time_e as DATETIME),INTERVAL 9 HOUR) BETWEEN ? AND ? GROUP BY hh ORDER BY hh) AS B ON A.n = B.hh"

		rows, err := conn.Query(sqlString, starttime, endtime)
		if err != nil {
			log.Printf("Charge_daily_graph_select err : %s", err)
		}

		for rows.Next() {
			var Time string
			var Amount *string
			var Count *string
			var Pay *string

			rows.Scan(&Time, &Amount, &Count, &Pay)

			m := structs.ChargeDailyGraphItem{Time: Time, Amount: Amount, Count: Count, Pay: Pay}

			result = append(result, m)

		}
		defer rows.Close()
		defer conn.Close()

		// log.Printf(" charge list : %s", result)

		return result
	} else {
		rows, err := conn.Query("SELECT A.n AS time , ifnull(B.amount,0) AS amount, ifnull(B.count,0) AS count,ifnull(B.pay,0) AS pay  FROM (SELECT @N := @N +1 AS n FROM charge_h, (select @N:=-1 from DUAL ) NN LIMIT 24) AS A LEFT JOIN (SELECT Hour(DATE_ADD(CAST(c.c_time_e as DATETIME),INTERVAL 9 HOUR)) as hh, sum(CAST(c.c_amount as FLOAT)) as amount ,COUNT(c.c_amount) as count ,sum(CAST(p.c_pay as UNSIGNED)-CAST(IFNULL(p.can_pay,0) as UNSIGNED)) as pay FROM charge_h as c RIGHT JOIN pay_h as p  ON c.c_id = p.pay_num WHERE DATE_ADD(CAST(c.c_time_e as DATETIME),INTERVAL 9 HOUR) BETWEEN ? AND ? AND p.d_id = ?  GROUP BY hh ORDER BY hh) AS B ON A.n = B.hh", starttime, endtime, device_id)
		if err != nil {
			log.Printf("Charge_daily_graph_select err : %s", err)
		}

		for rows.Next() {
			var Time string
			var Amount *string
			var Count *string
			var Pay *string

			rows.Scan(&Time, &Amount, &Count, &Pay)

			m := structs.ChargeDailyGraphItem{Time: Time, Amount: Amount, Count: Count, Pay: Pay}

			result = append(result, m)

		}
		defer rows.Close()
		defer conn.Close()

		// log.Printf(" charge list : %s", result)

		return result
	}

}

// 충전기 장애 이력
func Err_list_select(s_order, e_order string) []structs.ErrListItem {
	var result []structs.ErrListItem
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" User_survey_a err : %s", err)
	}

	rows, err := conn.Query("select * from device_err order by o_time desc limit ?,?", s_order, e_order)
	if err != nil {
		log.Printf("err_list_select err : %s", err)
	}

	for rows.Next() {
		var Id string
		var D_Id string
		var D_Error string
		var V_Id string
		var V_Error string
		var E_Info string
		var O_Time string

		rows.Scan(&Id, &D_Id, &D_Error, &V_Id, &V_Error, &E_Info, &O_Time)

		m := structs.ErrListItem{Id: Id, D_Id: D_Id, D_Error: D_Error, V_Id: V_Id, V_Error: V_Error, E_Info: E_Info, O_Time: O_Time}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	// log.Printf(" charge list : %s", result)

	return result
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 충전기 기본정보 카운트 반환
func Device_list_count_select() int {
	var result int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Device_list_count_select err : %s", err)
	}

	err = conn.QueryRow("select count(*) from device").Scan(&result)
	if err != nil {
		log.Printf("Device_list_count_select : %s", err)
	}

	defer conn.Close()

	return result
}

// 충전기 상태 정보 리스트 카운트 반환
func DeviceState_list_count_select() int {
	var result int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("DeviceState_list_count_select err : %s", err)
	}

	err = conn.QueryRow("select count(*) from device_status").Scan(&result)
	if err != nil {
		log.Printf("DeviceState_list_count_select : %s", err)
	}

	defer conn.Close()

	return result
}

// 충전기 요금단가관리 카운트 반환
func Price_list_count_select() int {
	var result int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Price_list_count_select err : %s", err)
	}

	err = conn.QueryRow("select count(*) from e_price").Scan(&result)
	if err != nil {
		log.Printf("Price_list_count_select : %s", err)
	}

	defer conn.Close()

	return result
}

// 충전기 업데이트리스트 카운트 반환
func Device_Update_list_count_select() int {
	var result int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Price_list_count_select err : %s", err)
	}

	err = conn.QueryRow("SELECT count(d.id)  FROM device AS d LEFT JOIN device_status AS s ON s.id = d.id").Scan(&result)
	if err != nil {
		log.Printf("Price_list_count_select : %s", err)
	}

	defer conn.Close()

	return result
}

// 충전기 로그 리스트 카운트 반환
func Log_list_count_select() int {
	var result int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Log_list_count_select err : %s", err)
	}

	err = conn.QueryRow("SELECT count(d.id) FROM device AS d LEFT JOIN device_status AS s ON s.id = d.id").Scan(&result)
	if err != nil {
		log.Printf("Log_list_count_select : %s", err)
	}

	defer conn.Close()

	return result
}

// 충전기 이력조회 카운트 반환
func Charge_list_count_select() int {
	var result int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Charge_list_count_select err : %s", err)
	}

	err = conn.QueryRow("SELECT count(p.pay_num) FROM charge_h as c RIGHT JOIN pay_h as p  ON c.c_id = p.pay_num").Scan(&result)
	if err != nil {
		log.Printf("Charge_list_count_select : %s", err)
	}

	defer conn.Close()

	return result
}

// 일자별 충전량 카운트 반환
func Charge_daily_list_count_select() int {
	var result int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Charge_daily_list_count_select err : %s", err)
	}

	err = conn.QueryRow("SELECT count(c.day) FROM (SELECT CAST(DATE_ADD(CAST(c.c_time_e as DATETIME),INTERVAL 9 HOUR) as DATE) as day FROM charge_h as c RIGHT JOIN pay_h as p  ON (c.c_id = p.pay_num) GROUP BY day) as c").Scan(&result)
	if err != nil {
		log.Printf("Charge_daily_list_count_select : %s", err)
	}

	defer conn.Close()

	return result
}

// 에러리스트 카운트 반환
func Err_list_count_select() int {
	var result int

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Err_list_count_select err : %s", err)
	}

	err = conn.QueryRow("select count(*) from device_err").Scan(&result)
	if err != nil {
		log.Printf("Err_list_count_select : %s", err)
	}

	defer conn.Close()

	return result
}

///////////////////////////////////////////////////////////			// 기능 핸들러

// 기기 등록
func Insert_device(Id, Latitude, Longitude, Addr, Model, S_Version, H_C_M, H_P_M, H_B_M, H_W_M, Enabled string) bool {
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("INSERT INTO device(id,latitude,longitude,addr,model,s_version,h_c_m,h_p_m,h_b_m,h_w_m,enabled) VALUES(?,?,?,?,?,?,?,?,?,?,?)", Id, Latitude, Longitude, Addr, Model, S_Version, H_C_M, H_P_M, H_B_M, H_W_M, Enabled)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	//상태 등록
	rows2, err := conn.Exec("INSERT INTO device_status(id) VALUES(?)", Id)
	if err != nil {
		fmt.Printf("rows2 : %s", err)
		return false
	}

	log.Printf("insert_device : %v , %v ", rows, rows2)

	defer conn.Close()

	return true
}

// 기기 수정
func Update_device(Id, Latitude, Longitude, Addr, Model, S_Version, H_C_M, H_P_M, H_B_M, H_W_M, Enabled string) bool {
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("UPDATE device SET id =? , latitude =? , longitude =? , addr =? , model = ?,s_version =? , h_c_m =? , h_p_m =? , h_b_m =? , h_w_m =? , enabled =? WHERE id = ?", Id, Latitude, Longitude, Addr, Model, S_Version, H_C_M, H_P_M, H_B_M, H_W_M, Enabled, Id)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("update_device : %v", rows)

	defer conn.Close()

	return true
}

// 기기 삭제
func Delete_device(Id string, Latitude string, Longitude string) bool {
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("DELETE FROM device WHERE id =? AND latitude =? AND longitude =?", Id, Latitude, Longitude)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	//상태 삭제
	rows2, err := conn.Exec("DELETE FROM device_status WHERE id =?", Id)
	if err != nil {
		fmt.Printf("rows2 : %s", err)
		return false
	}

	log.Printf("delete_device : %v , %v ", rows, rows2)

	defer conn.Close()

	return true
}

// 단가 등록
func Insert_price(price string, writer string) bool {
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("INSERT INTO e_price(price,writer) VALUES(?,?)", price, writer)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("insert_price : %v ", rows)

	defer conn.Close()

	return true
}

// 단가 삭제
func Delete_price(Id string, writer string) bool {
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		fmt.Println(err)
		return false
	}

	rows, err := conn.Exec("DELETE FROM e_price WHERE id =? AND writer =?", Id, writer)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		return false
	}

	log.Printf("delete_price : %v", rows)
	defer conn.Close()

	return true
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 충전기 설치 정보 검색
func Device_list_select_search(column, data string) []structs.DeviceListItem {
	var result []structs.DeviceListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" Device_list_select_search err : %s", err)
	}

	//컬럼이름이 존재하지 않는지 체크	-- sql injection 공격 차단
	c := map[string]bool{
		"id":        true,
		"latitude":  true,
		"longitude": true,
		"addr":      true,
		"model":     true,
		"s_version": true,
		"h_c_m":     true,
		"h_p_m":     true,
		"h_b_m":     true,
		"h_w_m":     true,
		"i_time":    true,
		"u_time":    true,
		"enabled":   true,
	}

	if !isSafeColumnName(c, column) {
		return result
	}

	queryStirng := fmt.Sprintf("SELECT * FROM device where %s like ? ORDER BY i_time DESC", column)

	rows, err := conn.Query(queryStirng, "%"+strings.Trim(data, " ")+"%")
	if err != nil {
		log.Printf(" Device_list_select_search err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Latitude string
		var Longitude string
		var Addr string
		var Model *string
		var S_Version *string
		var H_C_M string
		var H_P_M string
		var H_B_M string
		var H_W_M string
		var I_Time string
		var U_Time string
		var Enabled string

		rows.Scan(&Id, &Latitude, &Longitude, &Addr, &Model, &S_Version, &H_C_M, &H_P_M, &H_B_M, &H_W_M, &I_Time, &U_Time, &Enabled)

		m := structs.DeviceListItem{Id: Id, Latitude: Latitude, Longitude: Longitude, Addr: Addr, Model: Model, S_Version: S_Version, H_C_M: H_C_M, H_P_M: H_P_M, H_B_M: H_B_M, H_W_M: H_W_M, I_Time: I_Time, U_Time: U_Time, Enabled: Enabled}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 충전기 상태 정보 검색
func Device_state_list_select_search(dconn []string, column, data string) []structs.DeviceStateListItem {
	var result []structs.DeviceStateListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" device_state_list_select_search err : %s", err)
	}

	//컬럼이름이 존재하지 않는지 체크	-- sql injection 공격 차단
	c := map[string]bool{
		"id":      true,
		"state":   true,
		"state_e": true,
		"v_id":    true,
		"v_error": true,
		"e_info":  true,
		"u_time":  true,
		"f_state": true,
		"d_state": true,
	}

	if !isSafeColumnName(c, column) {
		return result
	}

	queryStirng := fmt.Sprintf("SELECT * FROM device_status where %s like ? ORDER BY u_time DESC", column)

	rows, err := conn.Query(queryStirng, "%"+strings.Trim(data, " ")+"%")
	if err != nil {
		log.Printf(" device_state_list_select_search err : %s", err)
	}

	for rows.Next() {

		var Id string
		// var Conn    string
		var State *string
		var State_e *string
		var V_Id *string
		var V_Error *string
		var E_Info *string
		var U_Time *string
		var F_State *string
		var D_State *string

		rows.Scan(&Id, &State, &State_e, &V_Id, &V_Error, &E_Info, &U_Time, &F_State, &D_State)

		m := structs.DeviceStateListItem{Id: Id, Conn: contains(dconn, Id), State: State, State_e: State_e, V_Id: V_Id, V_Error: V_Error, E_Info: E_Info, U_Time: U_Time, F_State: F_State, D_State: D_State}
		result = append(result, m)
	}

	defer rows.Close()
	defer conn.Close()

	return result
}

// 충전요금 단가 관리 검색
func Price_list_select_search(column, data string) []structs.PriceListItem {
	var result []structs.PriceListItem
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" Price_list_select_search err : %s", err)
	}

	//컬럼이름이 존재하지 않는지 체크	-- sql injection 공격 차단
	c := map[string]bool{
		"id":     true,
		"price":  true,
		"time":   true,
		"writer": true,
	}

	if !isSafeColumnName(c, column) {
		return result
	}

	queryStirng := fmt.Sprintf("SELECT * FROM e_price where %s like ?  ORDER BY time DESC", column)

	rows, err := conn.Query(queryStirng, "%"+strings.Trim(data, " ")+"%")
	if err != nil {
		log.Printf(" Price_list_select_search err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Price string
		var Time string
		var Writer string

		rows.Scan(&Id, &Price, &Time, &Writer)

		m := structs.PriceListItem{Id: Id, Price: Price, Time: Time, Writer: Writer}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 충전기 로그 리스트 정보 검색
func Log_list_select_search(column, data string) []structs.DeviceLogListItem {
	var result []structs.DeviceLogListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Log_list_select_search err : %s", err)
	}

	//컬럼이름이 존재하지 않는지 체크	-- sql injection 공격 차단
	c := map[string]bool{
		"id":      true,
		"addr":    true,
		"d_state": true,
	}

	if !isSafeColumnName(c, column) {
		return result
	} else {
		if strings.Compare(column, "d_state") == 0 {
			column = "s." + column
		} else {
			column = "d." + column
		}
	}

	queryStirng := fmt.Sprintf("SELECT d.id, d.addr,s.d_state  FROM device AS d LEFT JOIN device_status AS s ON s.id = d.id where %s like ?  ORDER BY s.d_state DESC ", column)

	rows, err := conn.Query(queryStirng, "%"+strings.Trim(data, " ")+"%")
	if err != nil {
		log.Printf("Log_list_select_search err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Addr string
		var D_State *string

		rows.Scan(&Id, &Addr, &D_State)

		m := structs.DeviceLogListItem{Id: Id, Addr: Addr, D_State: D_State}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

// 충전량 및 이력조회 검색
func Charge_list_select_search(column, data string) []structs.ChargeListItem {
	var result []structs.ChargeListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Charge_list_select_search err : %s", err)
	}

	//컬럼이름이 존재하지 않는지 체크	-- sql injection 공격 차단
	c := map[string]bool{
		"payid":     true,
		"d_id":      true,
		"c_pay":     true,
		"can_pay":   true,
		"pay_by":    true,
		"c_num":     true,
		"c_c_num":   true,
		"pay_day":   true,
		"c_pay_day": true,
		"c_conn":    true,
		"c_time":    true,
		"c_amount":  true,
		"c_time_s":  true,
		"c_time_e":  true,
		"c_time_u":  true,
	}

	if !isSafeColumnName(c, column) {
		return result
	} else {
		if strings.Compare(column, "payid") == 0 {
			column = "p.pay_num"
		}
	}

	queryStirng := fmt.Sprintf("SELECT p.pay_num as payid, d_id,c_pay,can_pay,pay_by,c_num,c_c_num,pay_day,c_pay_day,c_conn,c_time,c_amount,c_time_s,c_time_e,c_time_u FROM charge_h as c RIGHT JOIN pay_h as p  ON c.c_id = p.pay_num where %s like ?  ORDER BY c_time_s DESC", column)

	rows, err := conn.Query(queryStirng, "%"+strings.Trim(data, " ")+"%")
	if err != nil {
		log.Printf("cCharge_list_select_search err : %s", err)
	}

	for rows.Next() {
		var Payid string
		var D_id *string
		var C_pay *string
		var Can_pay *string
		var Pay_by *string
		var C_num *string
		var C_c_num *string
		var Pay_day *string
		var C_pay_day *string
		var C_conn *string
		var C_time *string
		var C_amount *string
		var C_time_s *string
		var C_time_e *string
		var C_time_u *string

		rows.Scan(&Payid, &D_id, &C_pay, &Can_pay, &Pay_by, &C_num, &C_c_num, &Pay_day, &C_pay_day, &C_conn, &C_time, &C_amount, &C_time_s, &C_time_e, &C_time_u)

		m := structs.ChargeListItem{Payid: Payid, D_id: D_id, C_pay: C_pay, Can_pay: Can_pay, Pay_by: Pay_by, C_num: C_num, C_c_num: C_c_num, Pay_day: Pay_day, C_pay_day: C_pay_day, C_conn: C_conn, C_time: C_time, C_amount: C_amount, C_time_s: C_time_s, C_time_e: C_time_e, C_time_u: C_time_u}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	// log.Printf(" charge list : %s", result)

	return result
}

// 충전기 업데이트 리스트 정보 검색
func Update_list_select_search(column, data string) []structs.DeviceUpdateListItem {
	var result []structs.DeviceUpdateListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Update_list_select_search err : %s", err)
	}

	//컬럼이름이 존재하지 않는지 체크	-- sql injection 공격 차단
	c := map[string]bool{
		"id":        true,
		"addr":      true,
		"s_version": true,
		"f_state":   true,
		"i_time":    true,
		"u_time":    true,
	}

	if !isSafeColumnName(c, column) {
		return result
	} else {
		if strings.Compare(column, "f_state") == 0 {
			column = "s." + column
		} else {
			column = "d." + column
		}
	}

	queryStirng := fmt.Sprintf("SELECT d.id, d.addr,d.model,d.s_version,s.f_state,d.i_time,d.u_time  FROM device AS d LEFT JOIN device_status AS s ON s.id = d.id  where %s like ? ORDER BY s.u_time DESC", column)

	rows, err := conn.Query(queryStirng, "%"+strings.Trim(data, " ")+"%")
	if err != nil {
		log.Printf("Update_list_select_search err : %s", err)
	}

	for rows.Next() {
		var Id string
		var Addr string
		var Model *string
		var S_Version *string
		var F_State *string
		var I_Time string
		var U_Time string

		rows.Scan(&Id, &Addr, &Model, &S_Version, &F_State, &I_Time, &U_Time)

		m := structs.DeviceUpdateListItem{Id: Id, Addr: Addr, Model: Model, S_Version: S_Version, F_State: F_State, I_Time: I_Time, U_Time: U_Time}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	return result
}

func Compare(column, s string) {
	panic("unimplemented")
}

// 일자별 충전량 검색
func Charge_daily_list_select_search(column, data string) []structs.ChargeDailyListItem {
	var result []structs.ChargeDailyListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf("Charge_daily_list_select_search err : %s", err)
	}

	//컬럼이름이 존재하지 않는지 체크	-- sql injection 공격 차단
	c := map[string]bool{
		"day":    true,
		"amount": true,
		"count":  true,
		"pay":    true,
	}

	if !isSafeColumnName(c, column) {
		return result
	}

	queryStirng := fmt.Sprintf("SELECT CAST(DATE_ADD(CAST(c.c_time_e as DATETIME),INTERVAL 9 HOUR) as DATE) as day, sum(CAST(c.c_amount as FLOAT)) as amount ,COUNT(c.c_amount) as count, sum(CAST(p.c_pay as UNSIGNED)-CAST(IFNULL(p.can_pay,0) as UNSIGNED)) as pay FROM charge_h as c RIGHT JOIN pay_h as p  ON (c.c_id = p.pay_num) GROUP BY day HAVING %s like ? ORDER BY day DESC", column)

	rows, err := conn.Query(queryStirng, "%"+strings.Trim(data, " ")+"%")
	if err != nil {
		log.Printf("Charge_daily_list_select_search err : %s", err)
	}

	for rows.Next() {
		var Day string
		var Amount *string
		var Count *string
		var Pay *string

		rows.Scan(&Day, &Amount, &Count, &Pay)

		m := structs.ChargeDailyListItem{Day: Day, Amount: Amount, Count: Count, Pay: Pay}
		if Day != "" {
			result = append(result, m)
		}

	}
	defer rows.Close()
	defer conn.Close()

	// log.Printf(" charge list : %s", result)

	return result
}

// 충전기 장애 이력 검색
func Err_list_select_search(column, data string) []structs.ErrListItem {
	var result []structs.ErrListItem

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" Err_list_select_search err : %s", err)
	}

	//컬럼이름이 존재하지 않는지 체크	-- sql injection 공격 차단
	c := map[string]bool{
		"id":      true,
		"d_id":    true,
		"d_error": true,
		"v_id":    true,
		"v_error": true,
		"e_info":  true,
		"o_time":  true,
	}

	if !isSafeColumnName(c, column) {
		return result
	}

	queryStirng := fmt.Sprintf("select * from device_err where %s like ? order by o_time desc", column)

	rows, err := conn.Query(queryStirng, "%"+strings.Trim(data, " ")+"%")
	if err != nil {
		log.Printf("err_list_select_search err : %s", err)
	}

	for rows.Next() {
		var Id string
		var D_Id string
		var D_Error string
		var V_Id string
		var V_Error string
		var E_Info string
		var O_Time string

		rows.Scan(&Id, &D_Id, &D_Error, &V_Id, &V_Error, &E_Info, &O_Time)

		m := structs.ErrListItem{Id: Id, D_Id: D_Id, D_Error: D_Error, V_Id: V_Id, V_Error: V_Error, E_Info: E_Info, O_Time: O_Time}
		result = append(result, m)
	}
	defer rows.Close()
	defer conn.Close()

	// log.Printf(" charge list : %s", result)

	return result
}

func isSafeColumnName(c map[string]bool, column string) bool {
	// 안전한 컬럼 이름 목록

	_, ok := c[column]
	return ok
}
