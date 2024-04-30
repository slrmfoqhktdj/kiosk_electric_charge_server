package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"

	// "io"
	"log"
	"math/rand"
	"net/http"

	// "os"
	"strconv"

	"time"

	v16 "github.com/aliml92/ocpp/v16"
	"github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"kiosk_electoric_charge/gosrc/adminjwt"
	"kiosk_electoric_charge/gosrc/adminpage"
	"kiosk_electoric_charge/gosrc/adminpagedbconn"
	"kiosk_electoric_charge/gosrc/db"
	"kiosk_electoric_charge/gosrc/ocpp_db_handler"
	"kiosk_electoric_charge/gosrc/payinfo"
	"kiosk_electoric_charge/gosrc/structs"
	"kiosk_electoric_charge/gosrc/updateapk"

	// "github.com/rs/cors"
	"go.uber.org/zap"
)

// 가격 반환 db
func get_pay_DB() float64 {
	// dbname := os.Getenv("DBNAME")
	// dbdata := os.Getenv("DBDATA")

	// conn, err := sql.Open(dbname, dbdata)

	conn, err := db.OpenDB()
	if err != nil {
		log.Printf(" get_pay_Handler : %s", err)
		return 0
	}

	var _price float64

	err = conn.QueryRow("SELECT price FROM e_price ORDER BY time DESC LIMIT 1").Scan(&_price)

	if err != nil {
		log.Printf(" get_pay_Handler get db : %s", err.Error())
		return 0
	}

	defer conn.Close()

	return _price
}

// 광고 리스트 반환 db
// func get_ADlist_DB(device_id string) []structs.ST_AD {
// 	var r []structs.ST_AD
// 	// dbname := os.Getenv("DBNAME")
// 	// dbdata := os.Getenv("DBDATA")

// 	// conn, err := sql.Open(dbname, dbdata)

// 	conn, err := db.OpenDB()
// 	if err != nil {
// 		log.Printf(" get_pay_Handler : %s", err)
// 	}

// 	rows, err := conn.Query("SELECT ad_id FROM device_ad WHERE d_id = ?", device_id)

// 	if err != nil {
// 		fmt.Printf("get_ADlist_DB func : %s", err)
// 	}

// 	for rows.Next() {
// 		var ad_id string
// 		rows.Scan(&ad_id)

// 		s := structs.ST_AD{ADId: ad_id}
// 		r = append(r, s)
// 	}

// 	defer rows.Close()
// 	defer conn.Close()

// 	return r

// }

//////////////////////////////////////////////////////////////////////////////////////////////////

// 광고 다운 핸들러
// func access_ad_Handler(w http.ResponseWriter, r *http.Request) {

// 	a_id := r.PostFormValue("ad_id")
// 	log.Printf("access %s", a_id)

// 	http.ServeFile(w, r, "./advertisement/"+a_id)
// }

// 광고 리스트 반환 핸들러
// func get_ad_Handler(w http.ResponseWriter, r *http.Request) {
// 	a_id := r.PostFormValue("device_id")

// 	log.Printf("id = %s", a_id)

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	data, _ := json.Marshal(get_ADlist_DB(a_id))

// 	log.Printf("list = %s", get_ADlist_DB(a_id))

// 	fmt.Fprint(w, string(data))

// }

// 가격 반환 핸들러
func get_pay_Handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(structs.ST_GetPay{Pay: get_pay_DB()})
	fmt.Fprint(w, string(data))
}

var logd *zap.SugaredLogger

var csms *Server

// 원격 선 명령
// clear cache - 충전 지역 인증 캐쉬(18p) 지우기	- (51p)
// ChangeAvailability - 가용성 변경(50p)	- 기기 상태 원격 명령
// ChangeConfiguration - 구성 변경 (51p)
// GetConfiguration	- 구성 가져오기(53p)
// RemoteStartTransaction - 원격 시작 트랜잭션(54p)
// RemoteStopTransaction - 원격 중지 트랜잭션(54p) -
// Reset	- 초기화 (57p)	- 진행중이면 stoptransaction 보내야함 -> 응용프로그램 or 하드웨어 다시 시작
// UnlockConnector	- 커넥터 잠금 해제(61p)	- 커넥터 케이블 고정장치가 오작동 하는 경우 분리
// DataTransfer	--- 커스텀 정보값

// GetDiagnostics	--- 진단확인 (53p)
// UpdateFirmware	--- 펌웨어 업데이트 (64p)

func Command_Handler(w http.ResponseWriter, r *http.Request) {
	if csms != nil {

		cp, b := csms.Load("device_02")
		if b {
			//////////////// ChangeConfiguration
			// var req Payload = ChangeConfigurationReqTT{
			// 	Key:   "WebSocketPingInterval",
			// 	Value: "50",
			// }
			// res, err := cp.Call("ChangeConfiguration", req)
			// if err != nil {
			// 	log.Printf("aa %s", err)
			// }
			// log.Print(res)
			// -----------------
			// var req Payload = ChangeConfigurationReqTT{
			// 	Key:   "MeterValuesSampledData",
			// 	Value: "Energy.Active.Import.Register,SoC",
			// }
			// res, err := cp.Call("ChangeConfiguration", req)
			// if err != nil {
			// 	log.Printf("aa %s", err)
			// }
			// log.Print(res)

			////////////////RemoteStartTransaction			54,59p 설명		-- 안보내지고 설정해야할게 많음
			// randomIdTag := strconv.Itoa(rand.Intn(1000000000)) // 랜덤 IdTag 만들고 - db에 거래 추가하기

			// ocpp_db_handler.DataTransfer_set_pay("device_02", "", "commend", "", randomIdTag, time.Now().Format("060102150405"))

			// var req Payload = RemoteStartTransactionReqTT{
			// 	IdTag:       randomIdTag,
			// 	ConnectorId: 1,
			// 	// ChargingProfile: v16.ChargingProfile{ChargingProfileId: 1},
			// }
			// res, err := cp.Call("RemoteStartTransaction", req)
			// if err != nil {
			// 	log.Printf("aa %s", err)
			// }
			// log.Print(res)

			////////////////RemoteStopTransaction		59p 설명

			// var req Payload = v16.RemoteStopTransactionReq{
			// 	TransactionId: 456789, // 검색 선택에 의해 유동적으로 넣어야함
			// }
			// res, err := cp.Call("RemoteStartTransaction", req)
			// if err != nil {
			// 	log.Printf("aa %s", err)
			// }
			// log.Print(res)

			////////////////clear cache
			// var req Payload = v16.ClearCacheReq{}
			// res, err := cp.Call("ClearCache", req)
			// if err != nil {
			// 	log.Printf("aa %s", err)
			// }
			// log.Print(res)

			////////////////ChangeAvailability - 81p Inoperative = 불가능. Operative = 가능		// 기기의 메인일때의 일회용 명령. 재부팅하면 초기화되버림. db에 device테이블에서 enable을 바꿔줘야 로딩에서 재제가 된다
			// var req Payload = v16.ChangeAvailabilityReq{
			// 	ConnectorId: 1,
			// 	Type:        "Inoperative",
			// }
			// res, err := cp.Call("ChangeAvailability", req)
			// if err != nil {
			// 	log.Printf("aa %s", err)
			// }
			// log.Print(res)

			////////////////Reset	- 97p
			//- Soft : 진행중인 충전 정상적으로 중지, stopTransaction을 보냄. 그다음 응용소프트웨어 다시시작
			//- Hard : 하드웨어 다시 시작
			// var req Payload = v16.ResetReq{
			// 	Type: "Soft",
			// }
			// res, err := cp.Call("Reset", req)
			// if err != nil {
			// 	log.Printf("aa %s", err)
			// }
			// log.Print(res)

			////////////////UnlockConnector

			////////////////GetConfiguration
			//- 없으면 모든 설정 반환
			//- AllowOfflineTxForUnknownId	-- 오프라인 일때 알수없는 id가 충전을 할수있는지.
			//- AuthorizationCacheEnabled  -- (옵션) - 인증캐쉬가 구현되어있는 경우, 인증캐쉬를 사용할건지 아닌지
			//- AuthorizeRemoteTxRequests	-- 원격 충전을 위해 승인을 받아야하는지.
			//- ClockAlignedDataInterval  -- 충전 스테이
			//- ConnectionTimeOut	-- 시작부터 충전 전까지의 취소 시간
			//- ConnectorPhaseRotation	-- 해당기기엔 해당없음 - NotApplicable로 답하면됨
			//- GetConfigurationMaxKeys -- GetConfiguration에서 요청된 구성키의 최대수
			//- HeartbeatInterval	-- heartbeat간격
			//- LocalAuthorizeOffline  -오프라인 일 때 충전 포인트가 인증 캐시 및/또는 로컬 인증 목록을 사용하여 사용자에게 인증할지 여부를 제어합니다 .
			//- LocalPreAuthorize		- 청구 지점이 중앙 시스템의 승인 응답을 기다리지 않고 트랜잭션을 시작하기 위해 승인 캐시 및/또는 로컬 승인 목록을 사용할지 여부를 제어합니다
			//- MeterValuesAlignedData - ClockAlignedDataInterval 초 마다 MeterValues.req PDU에 포함될 측정량 세트를 규정하는 쉼표로 구분된 목록입니다. (35p)
			//- MeterValuesSampledData	- MeterValueSampleInterval 초 마다 MeterValues.req PDU에 포함될 측정량 세트를 규정하는 쉼표로 구분된 목록입니다 . (35p)
			//- MeterValueSampleInterval	-- "MeterValues" PDU에 의해 전송되는 계량(또는 기타) 데이터 샘플링 사이의 시간(초)입니다. 서로 다른 시간 간격으로 데이터를 보내는 대신, 일관된 간격으로 데이터를 전송하도록 설정할 때 사용됩니다. 이것은 데이터의 정확성과 타임스탬프 기반의 동기화를 보장합니다.
			//- NumberOfConnectors	-- 충전포인터의 물리적 충전 커넥터 수
			//- ResetRetries	--실패한 충전 포인트 재설정을 재시도하는 횟수입니다
			//- StopTransactionOnEVSideDisconnect --true 로 설정되면 충전 포인트는 케이블이 EV에서 분리될 때 관리상 트랜잭션을 중지해야 합니다.	(49p)  false 로 설정된 경우 케이블이 EV에서 분리될 때 트랜잭션이 중지되지 않아야 합니다. EV가 다시 연결되면 에너지 전송이 다시 허용됩니다.
			//- StopTransactionOnInvalidId --승인되지 않은 승인 상태(Invalid)를 수신한 경우 청구 지점이 진행 중인 거래를 중지할지 여부,이 트랜잭션에 대한 StartTransaction.conf	(20p)  -True/False
			//- StopTxnAlignedData	-- StopTransaction.req MeterValues.req 의 TransactionData 요소에 포함될 시계 정렬 주기 측정량, 트랜잭션의 모든 ClockAlignedDataInterval 에 대한 PDU
			//- StopTxnSampledData -- StopTransaction.req PDU 의 TransactionData 요소에 포함될 샘플링된 측정 량,충전 세션 시작부터 MeterValueSampleInterval 초
			//- SupportedFeatureProfiles --지원되는 기능 프로필 목록입니다 ,Core, FirmwareManagement, LocalAuthListManagement, Reservation, SmartCharging and RemoteTrigger.
			//- TransactionMessageAttempts -- 중앙 시스템이 거래 관련 메시지를 처리하지 못할 때 청구 지점이 거래 관련 메시지 제출을 시도해야 하는 빈도입니다.
			//- TransactionMessageRetryInterval -- 중앙 시스템이 실패한 거래 관련 메시지를 다시 제출하기 전에 충전 지점이 기다려야 하는 시간
			//- UnlockConnectorOnEVSideDisconnect -- true 로 설정된 경우 충전 포인트는 케이블이 EV에서 분리될 때 충전 포인트 측 케이블의 잠금을 해제해야 합니다. (49p)
			//- WebSocketPingInterval -- 웹소켓 구현에만 해당됩니다. 0은 클라이언트 측 웹소켓 Ping/Pong을 비활성화합니다

			// var req Payload = v16.GetConfigurationReq{
			// 	Key: []string{"a", "b"},
			// }
			// res, err := cp.Call("GetConfiguration", req)
			// if err != nil {
			// 	log.Printf("aa %s", err)
			// }
			// log.Print(res)

			////////////////DataTransfer
			// var req Payload = v16.DataTransferReq{ // nvcat 옵션 변경
			// 	VendorId:  "device_02",
			// 	MessageId: "NVCATSetup",
			// 	Data:      "{'Catid':'7283730001','Bizno':'8158100527','Serverip':'211.33.136.19','ServerPort':'9701'}",
			// 	// Value: "{'Catid':'2393300001','Bizno':'2208115770','Serverip':'211.33.136.2','ServerPort':'9701'}",
			// }
			// res, err := cp.Call("DataTransfer", req)
			// if err != nil {
			// 	log.Printf("aa %s", err)
			// }
			// log.Print(res)

			// var req Payload = v16.DataTransferReq{ //nvcat 재시작 - 뭔가 오류났을때
			// 	VendorId:  "NVCAT",
			// 	MessageId: "NVCATRestart",
			// }
			// res, err := cp.Call("DataTransfer", req)
			// if err != nil {
			// 	log.Printf("aa %s", err)
			// }
			// log.Print(res)

			////////////////GetDiagnostics	// 진단 파일 가져오기 (로그 등등)
			// startTime이 없으면 최근/stoptime 기준 3일전까지의 데이터
			// stoptime이 없으면 최근까지의 데이터
			var req Payload = v16.GetDiagnosticsReq{
				Location:      string("file/log/operation"), // 두개짜리
				Retries:       2,
				RetryInterval: 10,
				StartTime:     "2024-03-03T15:04:05.999Z",
				StopTime:      "2024-03-06T15:04:05.999Z",
			}
			res, err := cp.Call("GetDiagnostics", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res) // 리턴값 파일명. 없으면 오류가 났거나 해당 조건에 맞는 진단이 없어 업로드 안함

			////////////////UpdateFirmware
			// var req Payload = v16.UpdateFirmwareReq{
			// //	Location: "1.1" + "$$" + string("file/apk/operation/1"), // 두개짜리
			//  Location:     "1.1" + "$$" + string("file/apk/operation/2"),		//스마트 폴
			// 	RetrieveDate: time.Now().Format("2006-01-02T15:04:05.000Z"),
			// }
			// res, err := cp.Call("UpdateFirmware", req)
			// if err != nil {
			// 	log.Printf("aa %s", err)
			// }
			// log.Print(res)

			////////////////TriggerMessage
			//BootNotification
			//DiagnosticsStatusNotification
			//FirmwareStatusNotification
			//Heartbeat
			//MeterValues
			//StatusNotification
			// var req Payload = v16.TriggerMessageReq{
			// 	RequestedMessage: "MeterValues",
			// }
			// res, err := cp.Call("TriggerMessage", req)
			// if err != nil {
			// 	log.Printf("aa %s", err)
			// }
			// log.Print(res)
		}

		if cp == nil {
			log.Print("cp is nil")
		}
		// log.Printf("cp2 : %s, b : %v", cp.Id, (b))
	}

}

// // 로그 다운 핸들러
// func get_logfile_Handler(w http.ResponseWriter, r *http.Request) {

// 	uploadFile, header, err := r.FormFile("upload_log")
// 	if err != nil { // 3
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprint(w, err)
// 		return
// 	}
// 	dirname := "./logfiles"                                    // 4
// 	os.MkdirAll(dirname, 0777)                                 // 5
// 	filepath := fmt.Sprintf("%s/%s", dirname, header.Filename) // 6
// 	file, err := os.Create(filepath)                           // 7

// 	if err != nil { // 9
// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprint(w, err)
// 		return
// 	}

// 	io.Copy(file, uploadFile) // 10

// 	defer file.Close() // 8

// 	w.WriteHeader(http.StatusOK) // 11
// 	fmt.Fprint(w, filepath)      // 12
// }

//////////////////////////////////////////////////////////////////////////////////////////////

// 대쉬보드 정보
func Dashboard_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("Dashboard_Page_Handler")

	// 한국 시간대 로드
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		fmt.Println("시간대 로딩 에러:", err)
		return
	}

	// 현재 시간을 한국 시간대로 설정
	now := time.Now().In(loc)

	Device_Id := "전체"
	StartTime := fmt.Sprintf("%s %s", now.AddDate(0, 0, -3).Format("2006-01-02"), "00:00:00")
	EndTime := fmt.Sprintf("%s %s", now.Format("2006-01-02"), "23:59:59")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DashboardPage{List: adminpagedbconn.Device_list_select("0", "3"), State: adminpagedbconn.Device_state_list_select(ConnectedDeivceList_Handler(), "0", "3"), Price: adminpagedbconn.Price_list_select("0", "3"), Daily: adminpagedbconn.Charge_daily_list_select("0", "4"), Maps: adminpagedbconn.Map_list_select(), Graph: adminpagedbconn.Charge_daily_graph_select(StartTime, EndTime, Device_Id)})
	w.Write(json)
}

// 기기 상태 정보
func DeviceStateList_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("DeviceStateList_Page_Handler %v,%v", r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceStateListPage{Count: adminpagedbconn.DeviceState_list_count_select(), Items: adminpagedbconn.Device_state_list_select(ConnectedDeivceList_Handler(), r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))})
	w.Write(json)
}

// 기기 상태 정보 검색
func DeviceStateList_Page_Search_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("DeviceStateList_Page_Search_Handler %v,%v", r.URL.Query().Get("column"), r.URL.Query().Get("data"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceStateListPage{Items: adminpagedbconn.Device_state_list_select_search(ConnectedDeivceList_Handler(), r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	w.Write(json)
}

// 원격 연결되어있는 기기 리스트 반환
func ConnectedDeivceList_Handler() []string {

	var returnValue []string

	for k := range csms.chargepoints {
		// log.Printf("csms.chargepoints ----- key = %v \n", k)

		returnValue = append(returnValue, k)
	}

	return returnValue
}

func main() {

	err := godotenv.Load("./info.env")
	if err != nil {
		log.Println(err)
		return
	}

	runtime.NumCPU() // 실행 환경의 물리적 CPU 최대 개수로 설정된다.

	//랜덤값 기준시드 정해주기
	rand.Seed(time.Now().UnixNano())

	logger, _ := zap.NewDevelopment()
	logd = logger.Sugar()
	defer logd.Sync()

	// set ocpp library's logger to zap logger
	SetLogger(logd)

	// start csms server with default configurations
	csms = NewServer()

	csms.AddSubProtocol("ocpp1.6")
	// csms.SetCheckOriginHandler(func(r *http.Request) bool { return true })
	csms.SetCheckOriginHandler(SetCheckOriginHanler_H)
	// csms.SetPreUpgradeHandler(customPreUpgradeHandler)	// 접속시 헤더의 Authorization 값으로 id, pw 를 보내 인증하는 방식
	csms.SetCallQueueSize(32)
	// csms.SetTimeoutConfig(ServerTimeoutConfig{OcppWait: 10000}) // 연결 제한시간
	// csms.SetTimeoutConfig(ServerTimeoutConfig{PingWait: 5000}) // 핑 제한시간

	// register charge-point-initiated action handlers
	csms.On("Authorize", AuthorizationHandler)           //	충전 시작, 끝에 사용자가 회원인지 - 18p 타임아웃은 충전포인트에서 해야함
	csms.On("BootNotification", BootNotificationHandler) //  부트
	// csms.After("BootNotification", SendChangeConfigration)
	csms.On("Heartbeat", HeartBeatHandler)                                         // 하트비트
	csms.On("StartTransaction", StartTransactionHandler)                           // 충전시작
	csms.On("StopTransaction", StopTransactionHandler)                             // 충전 중단
	csms.On("DataTransfer", DataTransferHandler)                                   // 커스텀 정보값
	csms.On("MeterValues", MeterValuesHandler)                                     // 계량기값에 대한 정보(41p) - ChangeConfiguration을 이용, 데이터 수집간격을 구성, 데이터를 지정할수도있음
	csms.On("StatusNotification", StatusNotificationHanler)                        //진단 상태 알림 (43p) - 상태가 변하거나 오류가나면, state = (82p)
	csms.On("DiagnosticsStatusNotification", DiagnosticsStatusNotificationHandler) // 업로드 상태를 알리는 알림 (39p)
	csms.On("FirmwareStatusNotification", FirmwareStatusNotificationHandler)       // 펌웨어 업데이트 진행사항 (39p)

	// 원격 선 명령	- core
	// clear cache - 충전 지역 인증 캐쉬(18p) 지우기
	// ChangeAvailability - 가용성 변경(50p)	- 기기 상태 원격 명령
	// ChangeConfiguration - 구성 변경 (51p)
	// GetConfiguration	- 구성 가져오기(53p)
	// RemoteStartTransaction - 원격 시작 트랜잭션(54p)
	// RemoteStopTransaction - 원격 중지 트랜잭션(54p) -
	// Reset	- 초기화 (57p)	- 진행중이면 stoptransaction 보내야함 -> 응용프로그램 or 하드웨어 다시 시작
	// UnlockConnector	- 커넥터 잠금 해제(61p)	- 커넥터 케이블 고정장치가 오작동 하는 경우 분리
	// DataTransfer	--- 커스텀 정보값

	// FIRMWARE MANAGEMENT
	// GetDiagnostics
	// UpdateFirmware

	go csms.Start(os.Getenv("CSMSSERVER"), os.Getenv("CSMSSERVER_PATH"), nil)

	//연결이 되어 있는지 확인
	// log.Printf("aa : %s", strconv.FormatBool(csms.IsConnected("device_01")))
	// log.Printf("aa : %s", strconv.FormatBool(csms.IsConnected("device_02")))

	//////////////////////////////////////////////////////////////////////////////////////////////

	mux := mux.NewRouter()

	// mux.Use(accessControlMiddleware)
	// mux := pat.New()
	// mux.Get("/ws", wsHandler)

	//negroni를 쓰면 따로 명시하지 않을시 "/"에 접근하면 index.html or default.html 파일에 액세스함
	// n := negroni.Classic()

	//////////////////////////////////////////////////////////////////////////////////////////////

	//관리자 페이지
	//로그인
	// mux.Path("/").HandlerFunc(adminpage.LoginPage_Handler)

	mux.Path("/login").HandlerFunc(adminpage.Login_Handler).Methods("POST")

	privateRouter := mux.PathPrefix("/").Subrouter()
	privateRouter.Use(jwtMiddleware)

	// mux.Path("/login").HandlerFunc(Command_Handler)
	//대쉬보드
	privateRouter.Path("/dashboardlist").HandlerFunc(Dashboard_Page_Handler).Methods("GET")
	//충전기 위치 정보
	privateRouter.Path("/maplist").HandlerFunc(adminpage.MapList_Page_Handler).Methods("GET")
	//충전기 설치 정보
	privateRouter.Path("/devicelist").HandlerFunc(adminpage.DeviceList_Page_Handler).Methods("GET")
	//충전기 상태 정보
	privateRouter.Path("/devicestatelist").HandlerFunc(DeviceStateList_Page_Handler).Methods("GET")
	//충전요금 단가 관리
	privateRouter.Path("/pricemng").HandlerFunc(adminpage.Pricemng_Page_Handler).Methods("GET")
	//충전기 업데이트
	privateRouter.Path("/deviceupdatelist").HandlerFunc(adminpage.Updatemng_Page_Handler).Methods("GET")
	//진단파일 요청 관리
	privateRouter.Path("/devicelogcmdpage").HandlerFunc(adminpage.Logmng_Page_Handler).Methods("GET")
	//충전량 및 이력조회
	privateRouter.Path("/chargelist").HandlerFunc(adminpage.Chargelist_Page_Handler).Methods("GET")
	//일자별 충전량 조회
	privateRouter.Path("/chargedailylist").HandlerFunc(adminpage.ChargeDailylist_Page_Handler).Methods("GET")
	//충전기 그래프 조회
	privateRouter.Path("/chargedailygraph").HandlerFunc(adminpage.ChargeDailyGraph_Page_Handler).Methods("POST")
	//충전기 장애 이력
	privateRouter.Path("/errlist").HandlerFunc(adminpage.Errlist_Page_Handler).Methods("GET")

	//소프트웨어 업데이트 관리
	// mux.Path("/swmng").HandlerFunc(adminpage.Swmng_Page_Handler)
	//광고물 관리
	//충전기별 광고 관리
	privateRouter.Path("/admng").HandlerFunc(adminpage.Admng_Page_Handler)

	//매출실적 및 정산 조회
	//사용자 관리/권한
	//공통코드

	//////////////////////////////////관리자페이지 검색///////////////////////////////////

	//충전기 설치 정보 검색
	privateRouter.Path("/devicelist/search").HandlerFunc(adminpage.DeviceList_Page_Search_Handler).Methods("GET")
	//충전기 상태 정보 검색
	privateRouter.Path("/devicestatelist/search").HandlerFunc(DeviceStateList_Page_Search_Handler).Methods("GET")
	//충전요금 단가 관리 검색
	privateRouter.Path("/pricemng/search").HandlerFunc(adminpage.Pricemng_Page_Search_Handler).Methods("GET")
	//충전기 업데이트 검색
	privateRouter.Path("/deviceupdatelist/search").HandlerFunc(adminpage.Updatemng_Page_Search_Handler).Methods("GET")
	//진단파일 요청 관리 검색
	privateRouter.Path("/devicelogcmdpage/search").HandlerFunc(adminpage.Logmng_Page_Search_Handler).Methods("GET")
	//충전량 및 이력조회 검색
	privateRouter.Path("/chargelist/search").HandlerFunc(adminpage.Chargelist_Page_Search_Handler).Methods("GET")
	//일자별 충전량 검색
	privateRouter.Path("/chargedailylist/search").HandlerFunc(adminpage.ChargeDailylist_Page_Search_Handler).Methods("GET")
	//충전기 장애 이력 검색
	privateRouter.Path("/errlist/search").HandlerFunc(adminpage.Errlist_Page_Search_Handler).Methods("GET")

	//////////////////////////////////관리자페이지 기능///////////////////////////////////

	//기기 등록
	privateRouter.Path("/device/register").HandlerFunc(adminpage.Register_Device_Handler).Methods("POST")
	//기기 수정
	privateRouter.Path("/device/updater").HandlerFunc(adminpage.Updater_Device_Handler).Methods("POST")
	//기기 삭제
	privateRouter.Path("/device/deleter").HandlerFunc(adminpage.Deleter_Device_Handler).Methods("POST")

	//단가 등록
	privateRouter.Path("/price/register").HandlerFunc(adminpage.Register_Price_Handler).Methods("POST")
	//단가 삭제
	privateRouter.Path("/price/deleter").HandlerFunc(adminpage.Deleter_Price_Handler).Methods("POST")

	//기기 조회
	// mux.Path("/device/selecter").HandlerFunc(adminpage.Selecter_Device_Handler)

	// //광고업로드
	// mux.Path("/upload/ad").HandlerFunc(adminpage.Errlist_Page_Handler)
	// //apk업로드
	// mux.Path("/upload/apk").HandlerFunc(adminpage.Errlist_Page_Handler)

	//////////////////////////////////////////////////////////////////////////////////////////////
	// 정보값 db

	//결제정보
	privateRouter.Path("/pay/history").HandlerFunc(payinfo.Set_Pay_Handler).Methods("POST")

	//충전시작
	privateRouter.Path("/charging/start").HandlerFunc(payinfo.Set_charging_start_Handler).Methods("POST")

	//충전종료
	privateRouter.Path("/charging/end").HandlerFunc(payinfo.Set_charging_end_Handler).Methods("POST")

	//////////////////////////////////////////////////////////////////////////////////////////////

	//앱 업데이트

	//버전 확인
	privateRouter.Path("/api/apk/version").HandlerFunc(updateapk.Get_apk_version_Handler).Methods("POST")

	//앱 apk 파일 다운
	privateRouter.Path("/file/apk/operation/1").HandlerFunc(updateapk.App_apk_download_Handler).Methods("POST")

	//앱 pole apk 파일 다운
	privateRouter.Path("/file/apk/operation/2").HandlerFunc(updateapk.App_pole_apk_download_Handler).Methods("POST")

	//앱 업데이트 완료
	privateRouter.Path("/file/apk/consequnce").HandlerFunc(updateapk.App_update_complete_Handler).Methods("POST")

	//////////////////////////////////////////////////////////////////////////////////////////////

	//광고 리스트
	// mux.Path("/api/ad").HandlerFunc(get_ad_Handler).Methods("POST")

	// //광고 접근 후 다운
	// //경로끝에 /를 붙여줘야 html상에서 보임
	// mux.Path("/file/ad").HandlerFunc(access_ad_Handler).Methods("POST")

	//////////////////////////////////////////////////////////////////////////////////////////////

	//로그 파일
	// mux.Path("/file/log/operation").HandlerFunc(get_logfile_Handler).Methods("POST")

	//////////////////////////////////////////////////////////////////////////////////////////////

	//가격 리턴
	privateRouter.Path("/api/pay").HandlerFunc(get_pay_Handler).Methods("POST")

	//////////////////////////////////////////////////////////////////////////////////////////////

	//명령
	// mux.Path("/command/cdl").HandlerFunc(ConnectedDeivceList_Handler).Methods("POST")

	//원격 시작
	privateRouter.Path("/command/rmstt").HandlerFunc(RemoteStartTransaction_Handler).Methods("POST")

	//원격 멈춤
	privateRouter.Path("/command/rmstp").HandlerFunc(RemoteStopTransaction_Handler).Methods("POST")

	//기기 사용 가능 설정
	// mux.Path("/command/abdc").HandlerFunc(AbledDevice_Handler).Methods("POST")

	//캐쉬 삭제
	privateRouter.Path("/command/cc").HandlerFunc(ClearCache_Handler).Methods("POST")

	//커넥터 사용 가능 설정
	privateRouter.Path("/command/cgab").HandlerFunc(ChangeAvailability_Handler).Methods("POST")

	//커넥터 해체

	//리셋
	privateRouter.Path("/command/rs").HandlerFunc(Reset_Handler).Methods("POST")

	//기기 설정값 리턴
	privateRouter.Path("/command/gcfr").HandlerFunc(GetConfiguration_Handler).Methods("POST")

	//기기 설정값 설정
	privateRouter.Path("/command/ccfr").HandlerFunc(ChangeConfiguration_Handler).Methods("POST")

	//nvcat 설정
	privateRouter.Path("/command/ncs").HandlerFunc(NVCAT_Setting_Handler).Methods("POST")

	//nvcat 재시작
	privateRouter.Path("/command/ncrs").HandlerFunc(NVCAT_Restart_Handler).Methods("POST")

	//nvcat 부분취소
	privateRouter.Path("/command/nccp").HandlerFunc(NVCAT_CanclePay_Handler).Methods("POST")

	//소프트웨어 업데이트
	privateRouter.Path("/command/udfw").HandlerFunc(UpdateFirmware_Handler).Methods("POST")

	//기기 메시지 트리거
	privateRouter.Path("/command/tgmg").HandlerFunc(TriggerMessage_Handler).Methods("POST")

	//로그 읽어오기
	privateRouter.Path("/command/gdns").HandlerFunc(GetDiagnostics_Handler).Methods("POST")

	//////////////////////////////////////////////////////////////////////////////////////////////

	// corsOpts := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"http://localhost:5005"}, //you service is available and allowed for this base url
	// 	AllowedMethods: []string{
	// 		http.MethodGet, //http methods for your app
	// 		http.MethodPost,
	// 		// http.MethodPut,
	// 		// http.MethodPatch,
	// 		// http.MethodDelete,
	// 		// http.MethodOptions,
	// 		// http.MethodHead,
	// 	},
	// 	AllowCredentials: true,
	// 	AllowedHeaders:   []string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "application/json"},

	// 	// AllowedHeaders: []string{
	// 	// 	"*", //or you can your header key values which you are using in your application
	// 	// },
	// })

	// muxhandler := corsOpts.Handler(mux)

	muxhandler := accessControlMiddleware(mux)

	// n.UseHandler(mux)
	http.ListenAndServe(os.Getenv("HTTPSERVER"), muxhandler)

	// http.ListenAndServe(":8099", mux)

	//////////////////////////////////////////////////////////////////////////////////////////////

	// logger, _ := zap.NewDevelopment()
	// logd = logger.Sugar()
	// defer logd.Sync()

	// // set ocpp library's logger to zap logger
	// SetLogger(logd)

	// // start csms server with default configurations
	// csms = NewServer()

	// csms.AddSubProtocol("ocpp1.6")
	// // csms.SetCheckOriginHandler(func(r *http.Request) bool { return true })
	// csms.SetCheckOriginHandler(SetCheckOriginHanler_H)
	// // csms.SetPreUpgradeHandler(customPreUpgradeHandler)	// 접속시 헤더의 Authorization 값으로 id, pw 를 보내 인증하는 방식
	// csms.SetCallQueueSize(32)
	// // csms.SetTimeoutConfig(ServerTimeoutConfig{OcppWait: 10000}) // 연결 제한시간
	// // csms.SetTimeoutConfig(ServerTimeoutConfig{PingWait: 5000}) // 핑 제한시간

	// // register charge-point-initiated action handlers
	// csms.On("Authorize", AuthorizationHandler)           //	충전 시작, 끝에 사용자가 회원인지 - 18p 타임아웃은 충전포인트에서 해야함
	// csms.On("BootNotification", BootNotificationHandler) //  부트
	// // csms.After("BootNotification", SendChangeConfigration)
	// csms.On("Heartbeat", HeartBeatHandler)                                         // 하트비트
	// csms.On("StartTransaction", StartTransactionHandler)                           // 충전시작
	// csms.On("StopTransaction", StopTransactionHandler)                             // 충전 중단
	// csms.On("DataTransfer", DataTransferHandler)                                   // 커스텀 정보값
	// csms.On("MeterValues", MeterValuesHandler)                                     // 계량기값에 대한 정보(41p) - ChangeConfiguration을 이용, 데이터 수집간격을 구성, 데이터를 지정할수도있음
	// csms.On("StatusNotification", StatusNotificationHanler)                        //진단 상태 알림 (43p) - 상태가 변하거나 오류가나면, state = (82p)
	// csms.On("DiagnosticsStatusNotification", DiagnosticsStatusNotificationHandler) // 업로드 상태를 알리는 알림 (39p)
	// csms.On("FirmwareStatusNotification", FirmwareStatusNotificationHandler)       // 펌웨어 업데이트 진행사항 (39p)

	// // 원격 선 명령	- core
	// // clear cache - 충전 지역 인증 캐쉬(18p) 지우기
	// // ChangeAvailability - 가용성 변경(50p)	- 기기 상태 원격 명령
	// // ChangeConfiguration - 구성 변경 (51p)
	// // GetConfiguration	- 구성 가져오기(53p)
	// // RemoteStartTransaction - 원격 시작 트랜잭션(54p)
	// // RemoteStopTransaction - 원격 중지 트랜잭션(54p) -
	// // Reset	- 초기화 (57p)	- 진행중이면 stoptransaction 보내야함 -> 응용프로그램 or 하드웨어 다시 시작
	// // UnlockConnector	- 커넥터 잠금 해제(61p)	- 커넥터 케이블 고정장치가 오작동 하는 경우 분리
	// // DataTransfer	--- 커스텀 정보값

	// // FIRMWARE MANAGEMENT
	// // GetDiagnostics
	// // UpdateFirmware

	// csms.Start("0.0.0.0:8999", "/ws/", nil)

	//연결이 되어 있는지 확인
	// log.Printf("aa : %s", strconv.FormatBool(csms.IsConnected("device_01")))
	// log.Printf("aa : %s", strconv.FormatBool(csms.IsConnected("device_02")))

	//////////////////////////////////////////////////////////////////////////////////////////////
}

// // HandlerFuncType은 파라미터를 받는 핸들러 함수의 타입입니다.	-- 대쉬보드, 상태정보등 현제 csms에 연결된 상태를 같이 반환해주기위해
// type HandlerFuncType func(http.ResponseWriter, *http.Request, []string)

// // withParam은 핸들러 함수에 파라미터를 추가해서 넘기는 래퍼 함수입니다.
// func withParam(handler HandlerFuncType, param []string) http.HandlerFunc {

// 	aa := param
// 	log.Printf("params : %v", aa)

// 	return func(w http.ResponseWriter, r *http.Request) {
// 		handler(w, r, param) // 파라미터를 추가하여 핸들러 함수를 호출합니다.
// 	}
// }

func jwtMiddleware(next http.Handler) http.Handler {

	// log.Print("jwtMiddleware")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			// log.Print("jwtMiddleware 1")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := adminjwt.ParseToken(tokenString)
		if err != nil {
			// log.Print("jwtMiddleware 2")
			fmt.Println("Error parsing token")
			//유효기간 시간 초과시
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err != nil {
			// log.Print("jwtMiddleware 3")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			// log.Print("jwtMiddleware 4")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			// log.Print("jwtMiddleware 5")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user"].(string)
		if !ok {
			// log.Print("jwtMiddleware 6")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// log.Print("jwtMiddleware 7")
		ctx := context.WithValue(r.Context(), structs.ContextStruct{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func accessControlMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5005")
		// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type ,application/json")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})

}

////////////////////////////////////////////////////////////////////////////////////////////////////////

func SetCheckOriginHanler_H(r *http.Request) bool {
	// log.Printf("SetCheckOriginHanler %v	\n", r)
	log.Println("SetCheckOriginHanler")

	// for k, v := range csms.chargepoints {
	// 	log.Printf("csms.chargepoints ----- key = %v , value = %v \n", k, v)
	// }

	// log.Printf("access : %v", r.Header)
	log.Printf("access : %v", r.Host)

	return true
}

// func SendChangeConfigration(cp *ChargePoint, payload Payload) {

// 	var req Payload = v16.ChangeConfigurationReq{
// 		Key:   "WebSocketPingInterval",
// 		Value: "20",
// 	}
// 	res, err := cp.Call("ChangeConfiguration", req)
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	log.Print(res)
// }

// func customPreUpgradeHandler(w http.ResponseWriter, r *http.Request) bool {
// 	u, p, ok := r.BasicAuth()

// 	if !ok {
// 		path := strings.Split(r.URL.Path, "/")
// 		ida := path[len(path)-1]
// 		log.Printf("error parsing basic auth %s, %s , %s", u, p, ida)

// 		w.WriteHeader(401)
// 		return false
// 	}
// 	path := strings.Split(r.URL.Path, "/")
// 	id := path[len(path)-1]
// 	log.Printf("%s is trying to connect with %s:%s", id, u, p)
// 	if u != id {
// 		log.Printf("username provided is correct: %s", u)
// 		w.WriteHeader(401)
// 		return false
// 	}
// 	return true
// }

func BootNotificationHandler(cp *ChargePoint, p Payload) Payload {
	// csms.Store(cp)

	// cp2,b := csms.Load("device_01")

	req := p.(*v16.BootNotificationReq)
	log.Printf("\nid: %s\nBootNotification: %v", cp.Id, req)

	// csms.Store(cp) // cp 저장
	// cp2, b := csms.Load("device_01")
	// log.Printf("cp2 : %s, b : %v", cp2.Id, (b))

	rb := ocpp_db_handler.BootNotification_db_handler(req.ChargePointModel, req.FirmwareVersion)

	var res Payload

	switch rb {
	case "1":
		res = &v16.BootNotificationConf{
			CurrentTime: time.Now().Format("2006-01-02T15:04:05.000Z"),
			Interval:    10,
			Status:      "Accepted",
		}

	case "2":
		res = &v16.BootNotificationConf{
			CurrentTime: time.Now().Format("2006-01-02T15:04:05.000Z"),
			Interval:    10,
			Status:      "Pending",
		}

	case "3":
		res = &v16.BootNotificationConf{
			CurrentTime: time.Now().Format("2006-01-02T15:04:05.000Z"),
			Interval:    10,
			Status:      "Rejected",
		}

	}

	return res
}

func AuthorizationHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*v16.AuthorizeReq)
	log.Printf("\nid: %s\nAuthorizeReq: %v", cp.Id, req)

	rb := ocpp_db_handler.Authorize_select_payuser(req.IdTag)

	var res Payload

	if rb {
		res = v16.AuthorizeConf{ ///Accepted/ConcurrentTx, Expired, Blocked 및 Invalid 에 따라 유효, 만료, (일시적으로) 차단 또는 블랙리스트 로 표시
			IdTagInfo: &v16.IdTagInfo{
				ExpiryDate: time.Now().Format("2006-01-02T15:04:05.000Z"),
				Status:     "Accepted",
			},
		}
	} else {
		res = v16.AuthorizeConf{ ///Accepted/ConcurrentTx, Expired, Blocked 및 Invalid 에 따라 유효, 만료, (일시적으로) 차단 또는 블랙리스트 로 표시
			IdTagInfo: &v16.IdTagInfo{
				ExpiryDate: time.Now().Format("2006-01-02T15:04:05.000Z"),
				Status:     "Invalid",
			},
		}
	}
	return res
}

func HeartBeatHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*v16.HeartbeatReq)
	log.Printf("\nid: %s\nHeartbeat: %v", cp.Id, req)

	// cp.ResetPingPong(30000) // ping 리셋

	var res Payload = &v16.HeartbeatConf{
		CurrentTime: time.Now().Format("2006-01-02T15:04:05.000Z"),
	}
	return res
}

func StartTransactionHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*v16.StartTransactionReq)
	log.Printf("\nid: %s\nStartTransaction: %v", cp.Id, req)
	log.Printf("time stemp : %s", req.Timestamp)
	// csms.Store(cp) // cp 저장

	rb := ocpp_db_handler.StartTransaction_set_start(req.IdTag, req.ConnectorId, float32(req.MeterStart)/1000.0, req.Timestamp)

	var res Payload
	Ti, err := strconv.Atoi(req.IdTag)

	if err != nil {
		log.Printf("StartTransactionHandler err %s", err)
		rb = false
	}

	if rb {
		res = v16.StartTransactionConf{
			IdTagInfo: v16.IdTagInfo{
				ExpiryDate:  time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05.000Z"), //만료일
				ParentIdTag: cp.Id,
				Status:      "Accepted",
			},

			TransactionId: Ti,
		}

	} else {
		res = v16.StartTransactionConf{
			IdTagInfo: v16.IdTagInfo{
				ExpiryDate:  time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05.000Z"), //만료일
				ParentIdTag: cp.Id,
				Status:      "Expired",
			},

			TransactionId: Ti,
		}
	}

	return res

}

func StopTransactionHandler(cp *ChargePoint, p Payload) Payload {
	log.Print("StopTransactionHandler")

	// req := p.(*v16.StopTransactionReq)
	req := p.(*StopTransactionRequestTT)
	log.Printf("\nid: %s\nStopTransaction: %v", cp.Id, req)
	// csms.Store(cp) // cp 저장

	// v16.MeterValue
	// v16.SampledValue

	// reason 가지고 기록을 해야하나 - 오류나 멈춤 기록에 넣어야 할듯

	rb := ocpp_db_handler.StopTransaction_update_end(strconv.Itoa(req.TransactionId), req.Timestamp, float32(req.MeterStop)/1000.0)

	var res Payload

	if rb {
		res = v16.StopTransactionConf{
			IdTagInfo: v16.IdTagInfo{
				ExpiryDate:  time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05.000Z"), //만료일
				ParentIdTag: cp.Id,
				Status:      "Accepted",
			},
		}
	} else {
		res = v16.StopTransactionConf{
			IdTagInfo: v16.IdTagInfo{
				ExpiryDate:  time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05.000Z"), //만료일
				ParentIdTag: cp.Id,
				Status:      "Expired",
			},
		}
	}

	return res

}

func DataTransferHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*v16.DataTransferReq)
	log.Printf("\nid: %s\nDataTransfer: %v", cp.Id, req)
	// csms.Store(cp) // cp 저장

	var res Payload

	// vendorid 가 없을경우 UnknownVendor 상태를 반환해야함
	if req.VendorId == "" {
		res = v16.DataTransferConf{
			Status: "UnKnownVendorId",
		}

		return res
	}

	switch req.MessageId {
	// case "adlist":
	// 	data, _ := json.Marshal(structs.ST_ADList{MessageType: 1, List: get_ADlist_DB(req.VendorId)})

	// 	res = v16.DataTransferConf{
	// 		Status: "Accepted",
	// 		Data:   string(data),
	// 	}

	case "i_pay": // 결제 정보 db 저장

		result := structs.PayItem{}
		err := json.Unmarshal([]byte(req.Data), &result) // json

		if err != nil { // 데이터가 안맞을때
			log.Print(" DataTransferHandler Failed to json.Unmarshal : ", err)
			errdata, _ := json.Marshal(structs.ST_ADList{MessageType: -1})
			res = v16.DataTransferConf{
				Status: "UnKnownData",
				Data:   string(errdata),
			}

			return res
		}

		e := ocpp_db_handler.DataTransfer_set_pay(result.Device_Id, result.Charge_Pay, result.Pay_By, result.Credit_Num, result.Pay_Num, result.Pay_Day)
		if !e { // db 오류가 있을때

			errdata, _ := json.Marshal(structs.ST_ADList{MessageType: -1})

			res = v16.DataTransferConf{
				Status: "UnKnownData",
				Data:   string(errdata),
			}

			return res
		}

		sccdata, _ := json.Marshal(structs.ST_ADList{MessageType: 2})
		res = v16.DataTransferConf{
			Status: "Accepted",
			Data:   string(sccdata),
		}

	case "c_pay": // 결제 정보 db 저장, 부분 취소
		result := structs.CanclePayItem{}
		err := json.Unmarshal([]byte(req.Data), &result) // json

		if err != nil { // 데이터가 안맞을때
			log.Print(" DataTransferHandler Failed to json.Unmarshal : ", err)
			errdata, _ := json.Marshal(structs.ST_ADList{MessageType: -1})
			res = v16.DataTransferConf{
				Status: "UnKnownData",
				Data:   string(errdata),
			}

			return res
		}

		e := ocpp_db_handler.DataTransfer_update_pay(result.Cancle_Pay, result.C_Pay_Day, result.C_C_Num, result.Pay_Num)
		if !e { // db 오류가 있을때

			errdata, _ := json.Marshal(structs.ST_ADList{MessageType: -1})

			res = v16.DataTransferConf{
				Status: "UnKnownData",
				Data:   string(errdata),
			}

			return res
		}

		sccdata, _ := json.Marshal(structs.ST_ADList{MessageType: 2})
		res = v16.DataTransferConf{
			Status: "Accepted",
			Data:   string(sccdata),
		}

	case "g_pay": // 전기 시세 가져오기

		pydata, _ := json.Marshal(structs.ST_G_Pay{MessageType: 3, Data: structs.ST_GetPay{Pay: get_pay_DB()}})

		res = v16.DataTransferConf{
			Status: "Accepted",
			Data:   string(pydata),
		}

	default:
		errdata, _ := json.Marshal(structs.ST_ADList{MessageType: -1})

		res = v16.DataTransferConf{
			Status: "UnKnownMessageId",
			Data:   string(errdata),
		}
	}

	return res

}

func MeterValuesHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*MeterValuesReqTT)

	log.Printf("\nid: %s\nMeterValues: %v ", cp.Id, req)
	log.Printf("size : %d, values : %s", len(req.MeterValue), req.MeterValue[0].SampledValue[0].Value)
	// csms.Store(cp) // cp 저장

	var TID string

	if len(strconv.Itoa(req.TransactionId)) < 8 {
		TID = fmt.Sprintf("%08d", req.TransactionId)
	}

	ocpp_db_handler.Metervalue_update(TID, req.ConnectorId, req.MeterValue[0].Timestamp, req.MeterValue[0].SampledValue[0].Value)

	var res Payload = v16.MeterValuesConf{} //72p
	return res

}

func StatusNotificationHanler(cp *ChargePoint, p Payload) Payload {
	req := p.(*StatusNotificationReqTT)
	log.Printf("\nid: %s\nStatusNotification: %v", cp.Id, req)
	// csms.Store(cp) // cp 저장

	var res Payload

	ocpp_db_handler.StatusNotification_update(cp.Id, req.ConnectorId, req.ErrorCode, req.Status, req.VendorId, req.VendorErrorCode, req.Info, req.Timestamp)
	//만약 에러코드가 존재하면 에러테이블에 저장하기
	if req.ErrorCode != "NoError" || strings.Trim(req.VendorId, " ") != "" || strings.Trim(req.Info, " ") != "" {
		ocpp_db_handler.Device_Error_insert(cp.Id, req.ErrorCode, req.VendorId, req.VendorErrorCode, req.Info)
	}

	res = v16.StatusNotificationConf{} //76p

	return res

}

func DiagnosticsStatusNotificationHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*v16.DiagnosticsStatusNotificationReq)
	log.Printf("\nid: %s\nDiagnosticsStatus: %v", cp.Id, req)

	ocpp_db_handler.DiagnosticsStatusNotification_update(cp.Id, req.Status)

	var res Payload = &v16.DiagnosticsStatusNotificationConf{}
	return res
}

func FirmwareStatusNotificationHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*v16.FirmwareStatusNotificationReq)
	log.Printf("\nid: %s\nFirmwareStatusNotification: %v", cp.Id, req)

	// cp.ResetPingPong(30000) // ping 리셋

	ocpp_db_handler.FirmwareStatusNotification_update(cp.Id, req.Status)

	var res Payload = &v16.FirmwareStatusNotificationConf{}
	return res
}

//////////////////////////////////////////////////////////////////////////////////////////////

// 어드민핸들러
type AbledDeviceReq struct {
	Enabled string `json:"enabled"`
}
type AbledDeviceConf struct {
	Status string `json:"status"`
}

// type ConnctedDevice struct {
// 	Id string `json:"id"`
// }

// 원격 시작 핸들러
func RemoteStartTransaction_Handler(w http.ResponseWriter, r *http.Request) {

	id := r.Header.Get("Id")
	msg := v16.RemoteStartTransactionReq{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		return
	}
	defer r.Body.Close()

	log.Printf("RemoteStartTransaction_Handler msg : %v", msg)
	var returnValue *v16.RemoteStartTransactionConf

	if csms != nil {
		cp, b := csms.Load(id)
		if b {

			////////////////RemoteStartTransaction			54,59p 설명		-- 안보내지고 설정해야할게 많음
			// randomIdTag := strconv.Itoa(rand.Intn(1000000000)) // 랜덤 IdTag 만들고 - db에 거래 추가하기

			ocpp_db_handler.DataTransfer_set_pay(id, "", "command", "", msg.IdTag, time.Now().Format("060102150405"))

			var req Payload = RemoteStartTransactionReqTT{
				IdTag:       msg.IdTag,
				ConnectorId: msg.ConnectorId,
			}
			res, err := cp.Call("RemoteStartTransaction", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)

			returnValue = res.(*v16.RemoteStartTransactionConf)

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &v16.RemoteStartTransactionConf{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 원격 멈춤 핸들러
func RemoteStopTransaction_Handler(w http.ResponseWriter, r *http.Request) {

	id := r.Header.Get("Id")
	msg := v16.RemoteStopTransactionReq{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("RemoteStopTransaction_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	log.Printf("RemoteStopTransaction_Handler msg : %v", msg)
	var returnValue *v16.RemoteStopTransactionConf

	if csms != nil {
		cp, b := csms.Load(id)
		if b {

			var req Payload = v16.RemoteStopTransactionReq{
				TransactionId: msg.TransactionId, // 검색 선택에 의해 유동적으로 넣어야함
			}
			res, err := cp.Call("RemoteStopTransaction", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)

			returnValue = res.(*v16.RemoteStopTransactionConf)

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &v16.RemoteStopTransactionConf{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// // 기기 사용 가능 설정	-- boot시 결정됨	"1" : 정상, "2": 보류, "3":불가
// func AbledDevice_Handler(w http.ResponseWriter, r *http.Request) {

// 	id := r.Header.Get("Id")
// 	msg := AbledDeviceReq{}
// 	decoder := json.NewDecoder(r.Body)
// 	err := decoder.Decode(&msg)
// 	if err != nil {
// 		log.Printf("AbledDevice_Handler err : %s", err)
// 		return
// 	}
// 	defer r.Body.Close()

// 	log.Printf("AbledDevice_Handler msg : %v", msg)
// 	var returnValue string

// 	rb := ocpp_db_handler.BootNotification_enable_update_handler(id, msg.Enabled)
// 	if rb {
// 		returnValue = "Accepted"
// 	} else {
// 		returnValue = "Rejected"
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(AbledDeviceConf{returnValue})
// }

// 캐쉬삭제
func ClearCache_Handler(w http.ResponseWriter, r *http.Request) {

	id := r.Header.Get("Id")
	msg := v16.ClearCacheReq{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("ClearCache_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	var returnValue *v16.ClearCacheConf

	if csms != nil {
		cp, b := csms.Load(id)
		if b {

			////////////////clear cache
			var req Payload = v16.ClearCacheReq{}
			res, err := cp.Call("ClearCache", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)

			returnValue = res.(*v16.ClearCacheConf)

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &v16.ClearCacheConf{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 커넥터 사용 가능 설정 핸들러
func ChangeAvailability_Handler(w http.ResponseWriter, r *http.Request) {

	id := r.Header.Get("Id")
	msg := v16.ChangeAvailabilityReq{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("ChangeAvailability_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	log.Printf("ChangeAvailability_Handler msg : %v", msg)
	var returnValue *v16.ChangeAvailabilityConf

	if csms != nil {
		cp, b := csms.Load(id)
		if b {

			////////////////ChangeAvailability - 81p Inoperative = 불가능. Operative = 가능		// 기기의 메인일때의 일회용 명령. 재부팅하면 초기화되버림. db에 device테이블에서 enable을 바꿔줘야 로딩에서 재제가 된다
			var req Payload = v16.ChangeAvailabilityReq{
				ConnectorId: msg.ConnectorId,
				Type:        msg.Type,
			}
			res, err := cp.Call("ChangeAvailability", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)

			returnValue = res.(*v16.ChangeAvailabilityConf)

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &v16.ChangeAvailabilityConf{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

//커넥터 해체

// 리셋
func Reset_Handler(w http.ResponseWriter, r *http.Request) {

	id := r.Header.Get("Id")
	msg := v16.ResetReq{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Reset_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	var returnValue *v16.ResetConf

	if csms != nil {
		cp, b := csms.Load(id)
		if b {

			////////////////Reset	- 97p
			//- Soft : 진행중인 충전 정상적으로 중지, stopTransaction을 보냄. 그다음 응용소프트웨어 다시시작
			//- Hard : 하드웨어 다시 시작
			var req Payload = v16.ResetReq{
				Type: msg.Type,
			}
			res, err := cp.Call("Reset", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)

			returnValue = res.(*v16.ResetConf)

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &v16.ResetConf{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

//기기 설정값 리턴
//- 없으면 모든 설정 반환
//- AllowOfflineTxForUnknownId	-- 오프라인 일때 알수없는 id가 충전을 할수있는지.
//- AuthorizationCacheEnabled  -- (옵션) - 인증캐쉬가 구현되어있는 경우, 인증캐쉬를 사용할건지 아닌지
//- AuthorizeRemoteTxRequests	-- 원격 충전을 위해 승인을 받아야하는지.
//- ClockAlignedDataInterval  -- 충전 스테이
//- ConnectionTimeOut	-- 시작부터 충전 전까지의 취소 시간
//- ConnectorPhaseRotation	-- 해당기기엔 해당없음 - NotApplicable로 답하면됨
//- GetConfigurationMaxKeys -- GetConfiguration에서 요청된 구성키의 최대수
//- HeartbeatInterval	-- heartbeat간격
//- LocalAuthorizeOffline  -오프라인 일 때 충전 포인트가 인증 캐시 및/또는 로컬 인증 목록을 사용하여 사용자에게 인증할지 여부를 제어합니다 .
//- LocalPreAuthorize		- 청구 지점이 중앙 시스템의 승인 응답을 기다리지 않고 트랜잭션을 시작하기 위해 승인 캐시 및/또는 로컬 승인 목록을 사용할지 여부를 제어합니다
//- MeterValuesAlignedData - ClockAlignedDataInterval 초 마다 MeterValues.req PDU에 포함될 측정량 세트를 규정하는 쉼표로 구분된 목록입니다. (35p)
//- MeterValuesSampledData	- MeterValueSampleInterval 초 마다 MeterValues.req PDU에 포함될 측정량 세트를 규정하는 쉼표로 구분된 목록입니다 . (35p)
//- MeterValueSampleInterval	-- "MeterValues" PDU에 의해 전송되는 계량(또는 기타) 데이터 샘플링 사이의 시간(초)입니다. 서로 다른 시간 간격으로 데이터를 보내는 대신, 일관된 간격으로 데이터를 전송하도록 설정할 때 사용됩니다. 이것은 데이터의 정확성과 타임스탬프 기반의 동기화를 보장합니다.
//- NumberOfConnectors	-- 충전포인터의 물리적 충전 커넥터 수
//- ResetRetries	--실패한 충전 포인트 재설정을 재시도하는 횟수입니다
//- StopTransactionOnEVSideDisconnect --true 로 설정되면 충전 포인트는 케이블이 EV에서 분리될 때 관리상 트랜잭션을 중지해야 합니다.	(49p)  false 로 설정된 경우 케이블이 EV에서 분리될 때 트랜잭션이 중지되지 않아야 합니다. EV가 다시 연결되면 에너지 전송이 다시 허용됩니다.
//- StopTransactionOnInvalidId --승인되지 않은 승인 상태(Invalid)를 수신한 경우 청구 지점이 진행 중인 거래를 중지할지 여부,이 트랜잭션에 대한 StartTransaction.conf	(20p)  -True/False
//- StopTxnAlignedData	-- StopTransaction.req MeterValues.req 의 TransactionData 요소에 포함될 시계 정렬 주기 측정량, 트랜잭션의 모든 ClockAlignedDataInterval 에 대한 PDU
//- StopTxnSampledData -- StopTransaction.req PDU 의 TransactionData 요소에 포함될 샘플링된 측정 량,충전 세션 시작부터 MeterValueSampleInterval 초
//- SupportedFeatureProfiles --지원되는 기능 프로필 목록입니다 ,Core, FirmwareManagement, LocalAuthListManagement, Reservation, SmartCharging and RemoteTrigger.
//- TransactionMessageAttempts -- 중앙 시스템이 거래 관련 메시지를 처리하지 못할 때 청구 지점이 거래 관련 메시지 제출을 시도해야 하는 빈도입니다.
//- TransactionMessageRetryInterval -- 중앙 시스템이 실패한 거래 관련 메시지를 다시 제출하기 전에 충전 지점이 기다려야 하는 시간
//- UnlockConnectorOnEVSideDisconnect -- true 로 설정된 경우 충전 포인트는 케이블이 EV에서 분리될 때 충전 포인트 측 케이블의 잠금을 해제해야 합니다. (49p)
//- WebSocketPingInterval -- 웹소켓 구현에만 해당됩니다. 0은 클라이언트 측 웹소켓 Ping/Pong을 비활성화합니다

func GetConfiguration_Handler(w http.ResponseWriter, r *http.Request) {

	id := r.Header.Get("Id")
	msg := v16.GetConfigurationReq{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("GetConfiguration_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	log.Printf("GetConfiguration_Handler msg : %v , length : %d", msg, len(msg.Key))
	var returnValue *GetConfigurationConfTT

	if csms != nil {
		cp, b := csms.Load(id)
		if b {
			////////////////GetConfiguration
			var req Payload = v16.GetConfigurationReq{
				Key: msg.Key,
			}
			res, err := cp.Call("GetConfiguration", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)

			returnValue = res.(*GetConfigurationConfTT)

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &GetConfigurationConfTT{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 기기 설정값 설정
func ChangeConfiguration_Handler(w http.ResponseWriter, r *http.Request) {

	id := r.Header.Get("Id")
	msg := ChangeConfigurationReqTT{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("ChangeConfiguration_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	log.Printf("ChangeConfiguration_Handler msg : %v", msg)
	var returnValue *v16.ChangeConfigurationConf

	if csms != nil {
		cp, b := csms.Load(id)
		if b {
			//////////////// ChangeConfiguration
			var req Payload = ChangeConfigurationReqTT{
				Key:   msg.Key,
				Value: msg.Value,
			}
			res, err := cp.Call("ChangeConfiguration", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)

			returnValue = res.(*v16.ChangeConfigurationConf)

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &v16.ChangeConfigurationConf{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// nvcat 설정
func NVCAT_Setting_Handler(w http.ResponseWriter, r *http.Request) {

	id := r.Header.Get("Id")
	msg := v16.DataTransferReq{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("NVCAT_Setting_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	log.Printf("NVCAT_Setting_Handler msg : %v", msg)
	var returnValue *v16.DataTransferConf

	if csms != nil {
		cp, b := csms.Load(id)
		if b {

			var req Payload = v16.DataTransferReq{ // nvcat 옵션 변경
				VendorId:  msg.VendorId,
				MessageId: msg.MessageId,
				Data:      msg.Data,
				// "{'Catid':'7283730001','Bizno':'8158100527','Serverip':'211.33.136.19','ServerPort':'9701'}"
				// Value: "{'Catid':'2393300001','Bizno':'2208115770','Serverip':'211.33.136.2','ServerPort':'9701'}",
			}
			res, err := cp.Call("DataTransfer", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)

			returnValue = res.(*v16.DataTransferConf)
		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &v16.DataTransferConf{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// nvcat 재시작
func NVCAT_Restart_Handler(w http.ResponseWriter, r *http.Request) {

	id := r.Header.Get("Id")
	msg := v16.DataTransferReq{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("NVCAT_Restart_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	log.Printf("NVCAT_Restart_Handler msg : %v", msg)
	var returnValue *v16.DataTransferConf

	if csms != nil {
		cp, b := csms.Load(id)
		if b {

			var req Payload = v16.DataTransferReq{ //nvcat 재시작 - 뭔가 오류났을때
				VendorId:  msg.VendorId,
				MessageId: msg.MessageId,
			}
			res, err := cp.Call("DataTransfer", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)

			returnValue = res.(*v16.DataTransferConf)
		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &v16.DataTransferConf{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// nvcat 부분취소
func NVCAT_CanclePay_Handler(w http.ResponseWriter, r *http.Request) {

	id := r.Header.Get("Id")
	msg := v16.DataTransferReq{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("NVCAT_CanclePay_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	log.Printf("NVCAT_CanclePay_Handler msg : %v", msg)
	var returnValue *v16.DataTransferConf

	if csms != nil {
		cp, b := csms.Load(id)
		if b {

			var req Payload = v16.DataTransferReq{ //nvcat 부분취소
				VendorId:  msg.VendorId,
				MessageId: msg.MessageId,
				Data:      msg.Data,
			}
			res, err := cp.Call("DataTransfer", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)

			returnValue = res.(*v16.DataTransferConf)
		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &v16.DataTransferConf{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

//소프트웨어 업데이트

func UpdateFirmware_Handler(w http.ResponseWriter, r *http.Request) {

	id := r.Header.Get("Id")
	msg := v16.UpdateFirmwareReq{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("UpdateFirmware_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	log.Printf("UpdateFirmware_Handler msg : %v", msg)
	var returnValue *v16.UpdateFirmwareConf

	if csms != nil {
		cp, b := csms.Load(id)
		if b {

			////////////////UpdateFirmware
			var req Payload = v16.UpdateFirmwareReq{
				// //	Location: "1.1" + "$$" + string("file/apk/operation/1"), // 두개짜리
				// Location:     "1.1" + "$$" + string("file/apk/operation/2"), //스마트 폴
				// RetrieveDate: time.Now().Format("2006-01-02T15:04:05.000Z"),
				Location:      msg.Location,
				Retries:       msg.Retries,
				RetryInterval: msg.RetryInterval,
				RetrieveDate:  msg.RetrieveDate,
			}
			res, err := cp.Call("UpdateFirmware", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)

			returnValue = res.(*v16.UpdateFirmwareConf)
		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &v16.UpdateFirmwareConf{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 기기 메시지 트리거
func TriggerMessage_Handler(w http.ResponseWriter, r *http.Request) {

	id := r.Header.Get("Id")
	msg := v16.TriggerMessageReq{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("TriggerMessage_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	log.Printf("TriggerMessage_Handler msg : %v", msg)
	var returnValue *v16.TriggerMessageConf

	if csms != nil {
		cp, b := csms.Load(id)
		if b {

			////////////////TriggerMessage
			//BootNotification
			//DiagnosticsStatusNotification
			//FirmwareStatusNotification
			//Heartbeat
			//MeterValues
			//StatusNotification
			var req Payload = v16.TriggerMessageReq{
				RequestedMessage: msg.RequestedMessage,
			}
			res, err := cp.Call("TriggerMessage", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)

			returnValue = res.(*v16.TriggerMessageConf)
		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &v16.TriggerMessageConf{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

//로그 읽어오기

func GetDiagnostics_Handler(w http.ResponseWriter, r *http.Request) {

	id := r.Header.Get("Id")
	msg := v16.GetDiagnosticsReq{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("GetDiagnostics_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	log.Printf("GetDiagnostics_Handler msg : %v", msg)
	var returnValue *v16.GetDiagnosticsConf

	if csms != nil {
		cp, b := csms.Load(id)
		if b {

			////////////////GetDiagnostics	// 진단 파일 가져오기 (로그 등등)
			// startTime이 없으면 최근/stoptime 기준 3일전까지의 데이터
			// stoptime이 없으면 최근까지의 데이터
			var req Payload = v16.GetDiagnosticsReq{
				Location: msg.Location, // 두개짜리
				// Location:      string("file/log/operation"),
				Retries:       msg.Retries,
				RetryInterval: msg.RetryInterval,
				StartTime:     msg.StartTime,
				StopTime:      msg.StopTime,
			}
			res, err := cp.Call("GetDiagnostics", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res) // 리턴값 파일명. 없으면 오류가 났거나 해당 조건에 맞는 진단이 없어 업로드 안함

			returnValue = res.(*v16.GetDiagnosticsConf)
		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &v16.GetDiagnosticsConf{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}
