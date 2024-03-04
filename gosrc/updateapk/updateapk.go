package updateapk

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"kiosk_electoric_charge/gosrc/structs"
)

// 업데이트 완료
func set_update_version(id string, version string) {
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rows, err := conn.Query("UPDATE device(s_version) VALUES (?) Where id = ?", version, id)
	if err != nil {
		fmt.Printf("rows1 : %s", err)
		os.Exit(1)
	}

	log.Printf("User_survey_t : %v", rows)

	defer conn.Close()

}

// 버전 반환 db
func get_version_DB() string {
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
	if err != nil {
		log.Printf(" get_version_Handler : %s", err)
	}

	var _version string

	err = conn.QueryRow("SELECT version FROM app_vs ORDER BY id DESC LIMIT 1").Scan(&_version)

	if err != nil {
		log.Printf(" get_version_Handler get db : %s", err)
	}

	defer conn.Close()

	log.Printf("get_version_db : %s", _version)

	return _version
}

//////////////////////////////////////////////////////////////////////////////////////////////////

// 앱 apk 파일 다운
func App_apk_download_Handler(w http.ResponseWriter, r *http.Request) {

	a_id := r.PostFormValue("device_id")
	log.Printf("App_apk_download %s", a_id)

	http.ServeFile(w, r, "./apkmanager/kiosk_electric_charger.zip")
}

// 앱 apk 파일 다운
func App_pole_apk_download_Handler(w http.ResponseWriter, r *http.Request) {

	a_id := r.PostFormValue("device_id")
	log.Printf("App_apk_download %s", a_id)

	http.ServeFile(w, r, "./apkmanager/kiosk_electric_charger_pole.zip")
}

// 버전 확인
func Get_apk_version_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Get_apk_version")

	data, _ := json.Marshal(structs.ST_GetVersion{Version: get_version_DB()})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data))
}

// 앱 업데이트 완료
func App_update_complete_Handler(w http.ResponseWriter, r *http.Request) {

	a_id := r.PostFormValue("d_id")
	_v := r.PostFormValue("d_v")
	log.Printf("update version %s - %s", a_id, _v)

	set_update_version(a_id, _v)

	log.Printf("app_update_complete %s", a_id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, http.StatusOK)
}
