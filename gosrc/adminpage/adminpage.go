package adminpage

import (
	"encoding/json"
	// "fmt"
	"log"
	"net/http"

	// "os"
	// "path/filepath"
	// "strings"

	// "time"

	// _ "github.com/go-sql-driver/mysql"

	// "text/template"

	"kiosk_electoric_charge/gosrc/adminpagedbconn"
	"kiosk_electoric_charge/gosrc/structs"
)

// // // html 구획 나누기
// var templ = func() *template.Template {

// 	//template 에 function 넣기
// 	funcadd := template.FuncMap{
// 		"add": func(a int, b int) int {
// 			return a + b
// 		},
// 	}

// 	// t := template.New("")
// 	t := template.New("").Funcs(funcadd)
// 	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
// 		if strings.Contains(path, ".html") {
// 			fmt.Println(path)
// 			_, err = t.ParseFiles(path)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 		}
// 		return err
// 	})

// 	if err != nil {
// 		fmt.Print(err)
// 	}
// 	return t
// }()

func Login_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("Login_Handler")

	msg := structs.LoginData{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Login_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.LoginPage{Items: adminpagedbconn.Login_select(msg.Email, msg.PWD)})
	w.Write(json)
}

// // 대쉬보드 정보
// func Dashboard_Page_Handler(w http.ResponseWriter, r *http.Request, conn []string) {

// 	log.Print("Dashboard_Page_Handler")

// 	// 한국 시간대 로드
// 	loc, err := time.LoadLocation("Asia/Seoul")
// 	if err != nil {
// 		fmt.Println("시간대 로딩 에러:", err)
// 		return
// 	}

// 	// 현재 시간을 한국 시간대로 설정
// 	now := time.Now().In(loc)

// 	Device_Id := "전체"
// 	StartTime := fmt.Sprintf("%s %s", now.AddDate(0, 0, -3).Format("2006-01-02"), "00:00:00")
// 	EndTime := fmt.Sprintf("%s %s", now.Format("2006-01-02"), "23:59:59")

// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "application/json")
// 	json, _ := json.Marshal(&structs.DashboardPage{List: Device_list_select(), State: Device_state_list_select(conn), Price: price_list_select(), Daily: charge_daily_list_select(), Maps: map_list_select(), Graph: charge_daily_graph_select(StartTime, EndTime, Device_Id)})
// 	w.Write(json)
// }

// 충전기 위치 정보
func MapList_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("MapList_Page_Handler")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.MapListPage{Items: adminpagedbconn.Map_list_select()})
	w.Write(json)
}

// 충전기 설치 정보
func DeviceList_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("DeviceList_Page_Handler %v,%v", r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceListPage{Count: adminpagedbconn.Device_list_count_select(), Items: adminpagedbconn.Device_list_select(r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))})
	w.Write(json)
}

// 충전기 설치 정보 검색
func DeviceList_Page_Search_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("DeviceList_Page_Search_Handler %v,%v", r.URL.Query().Get("column"), r.URL.Query().Get("data"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceListPage{Items: adminpagedbconn.Device_list_select_search(r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	w.Write(json)
}

// // 충전기 상태 정보
// func DeviceStateList_Page_Handler(w http.ResponseWriter, r *http.Request, conn []string) {

// 	log.Print("DeviceStateList_Page_Handler")

// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "application/json")
// 	json, _ := json.Marshal(&structs.DeviceStateListPage{Items: Device_state_list_select(conn)})
// 	w.Write(json)
// }

// 충전요금 단가 관리
func Pricemng_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Pricemng_Page_Handler %v,%v", r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.PricemngPage{Count: adminpagedbconn.Price_list_count_select(), Items: adminpagedbconn.Price_list_select(r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))})
	w.Write(json)
}

// 충전요금 단가 관리 검색
func Pricemng_Page_Search_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Pricemng_Page_Search_Handler %v,%v", r.URL.Query().Get("column"), r.URL.Query().Get("data"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.PricemngPage{Items: adminpagedbconn.Price_list_select_search(r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	w.Write(json)
}

// 충전기 업데이트 관리
func Updatemng_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Updatemng_Page_Handler %v,%v", r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.UpdatemngPage{Count: adminpagedbconn.Device_Update_list_count_select(), Items: adminpagedbconn.Update_list_select(r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))})
	w.Write(json)
}

// 충전기 업데이트 관리 검색
func Updatemng_Page_Search_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Updatemng_Page_Search_Handler %v,%v", r.URL.Query().Get("column"), r.URL.Query().Get("data"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.UpdatemngPage{Items: adminpagedbconn.Update_list_select_search(r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	w.Write(json)
}

// 충전기 로그 관리
func Logmng_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Logmng_Page_Handler %v,%v", r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.LogmngPage{Count: adminpagedbconn.Log_list_count_select(), Items: adminpagedbconn.Log_list_select(r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))})
	w.Write(json)
}

// 충전기 로그 관리 검색
func Logmng_Page_Search_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Logmng_Page_Search_Handler %v,%v", r.URL.Query().Get("column"), r.URL.Query().Get("data"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.LogmngPage{Items: adminpagedbconn.Log_list_select_search(r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	w.Write(json)
}

// 충전량 및 이력조회
func Chargelist_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Chargelist_Page_Handler %v,%v", r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.ChargeListPage{Count: adminpagedbconn.Charge_list_count_select(), Items: adminpagedbconn.Charge_list_select(r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))})
	w.Write(json)

}

// 충전량 및 이력조회 검색
func Chargelist_Page_Search_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Chargelist_Page_Search_Handler %v,%v", r.URL.Query().Get("column"), r.URL.Query().Get("data"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.ChargeListPage{Items: adminpagedbconn.Charge_list_select_search(r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	w.Write(json)

}

// 일자별 충전량 조회
func ChargeDailylist_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("ChargeDailylist_Page_Handler %v,%v", r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.ChargeDailyListPage{Count: adminpagedbconn.Charge_daily_list_count_select(), Items: adminpagedbconn.Charge_daily_list_select(r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))})
	w.Write(json)
}

// 일자별 충전량 검색
func ChargeDailylist_Page_Search_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("ChargeDailylist_Page_Search_Handler %v,%v", r.URL.Query().Get("column"), r.URL.Query().Get("data"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.ChargeDailyListPage{Items: adminpagedbconn.Charge_daily_list_select_search(r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	w.Write(json)
}

// 충전기 그래프 조회
func ChargeDailyGraph_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("ChargeDailyGraph_Page_Handler")

	msg := structs.DeviceGraphTimeData{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("ChargeDailyGraph_Page_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.ChargeDailyGraphPage{ListItem: adminpagedbconn.Charge_daily_graph_list_select(), Items: adminpagedbconn.Charge_daily_graph_select(msg.StartTime, msg.EndTime, msg.DeviceID)})
	w.Write(json)
}

// 충전기 장애 이력
func Errlist_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Errlist_Page_Handler %v,%v", r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.ErrListPage{Count: adminpagedbconn.Err_list_count_select(), Items: adminpagedbconn.Err_list_select(r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))})
	w.Write(json)

}

// 충전기 장애 이력 검색
func Errlist_Page_Search_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Errlist_Page_Search_Handler %v,%v", r.URL.Query().Get("column"), r.URL.Query().Get("data"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.ErrListPage{Items: adminpagedbconn.Err_list_select_search(r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	w.Write(json)
}

// // 소프트웨어 업데이트 관리
// func Swmng_Page_Handler(w http.ResponseWriter, r *http.Request) {
// 	templ.ExecuteTemplate(w, "main_menu", &structs.SwmngPage{Title: "SW 업데이트 관리", Page: 4, Items: sw_list_select()})
// }

//광고물 관리

// // 충전기별 광고 관리
// func Admng_Page_Handler(w http.ResponseWriter, r *http.Request) {
// 	templ.ExecuteTemplate(w, "main_menu", &structs.AdmngPage{Title: "충전기별 광고 관리", Page: 3, Items: adminpagedbconn.Ad_list_select()})
// }

//매출실적 및 정산 조회
//사용자 관리 권한
//공통코드

///////////////////////////////////////////////////////////			// 기능 핸들러

// 기기 등록
func Register_Device_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("Register_Device_Handler")

	msg := structs.DeviceListItem{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Register_Device_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceIUDDataForm{Items: adminpagedbconn.Insert_device(msg.Id, msg.Latitude, msg.Longitude, msg.Addr, *msg.Model, *msg.S_Version, msg.H_C_M, msg.H_P_M, msg.H_B_M, msg.H_W_M, msg.Enabled)})
	w.Write(json)

}

// 기기 수정
func Updater_Device_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("Updater_Device_Handler")

	msg := structs.DeviceListItem{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Updater_Device_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceIUDDataForm{Items: adminpagedbconn.Update_device(msg.Id, msg.Latitude, msg.Longitude, msg.Addr, *msg.Model, *msg.S_Version, msg.H_C_M, msg.H_P_M, msg.H_B_M, msg.H_W_M, msg.Enabled)})
	w.Write(json)
}

// 기기 삭제
func Deleter_Device_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("Deleter_Device_Handler")

	msg := structs.DeviceListItem{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Deleter_Device_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceIUDDataForm{Items: adminpagedbconn.Delete_device(msg.Id, msg.Latitude, msg.Longitude)})
	w.Write(json)
}

// 단가 등록
func Register_Price_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("Register_Price_Handler")

	msg := structs.PriceListItem{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Register_Price_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceIUDDataForm{Items: adminpagedbconn.Insert_price(msg.Price, msg.Writer)})
	w.Write(json)

}

// 단가 삭제
func Deleter_Price_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("Deleter_Price_Handler")

	msg := structs.PriceListItem{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Deleter_Price_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceIUDDataForm{Items: adminpagedbconn.Delete_price(msg.Id, msg.Writer)})
	w.Write(json)
}

// 기기 조회
// func Selecter_Device_Handler(w http.ResponseWriter, r *http.Request) {
// 	device_name := r.PostFormValue("device-name-select")
// 	device_state := r.PostFormValue("device-state-select")
// 	device_addr := r.PostFormValue("device-addr-select")

// 	log.Printf("value %s,%s,%s", device_name, device_state, device_addr)

// 	http.Redirect(w, r, "/devicelist", http.StatusMovedPermanently)
// }
