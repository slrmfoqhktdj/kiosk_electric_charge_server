package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"time"

	v16 "github.com/aliml92/ocpp/v16"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/pat"
	"github.com/joho/godotenv"
	"github.com/urfave/negroni"

	"kiosk_electoric_charge/gosrc/adminpage"
	"kiosk_electoric_charge/gosrc/ocpp_db_handler"
	"kiosk_electoric_charge/gosrc/payinfo"
	"kiosk_electoric_charge/gosrc/structs"
	"kiosk_electoric_charge/gosrc/updateapk"

	// "github.com/gorilla/websocket"
	// "time"
	// ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	// "github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	// "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"go.uber.org/zap"
)

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// type Message struct {
// 	Type string      `json:"type"`
// 	Data interface{} `json:"data"`
// }

// func wsHandler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	for {
// 		m := &Message{}
// 		err := conn.ReadJSON(m)
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		err = conn.WriteJSON(m)
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 	}
// }

// 가격 반환 db
func get_pay_DB() float64 {
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	conn, err := sql.Open(dbname, dbdata)
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
func get_ADlist_DB(device_id string) []structs.ST_AD {
	dbname := os.Getenv("DBNAME")
	dbdata := os.Getenv("DBDATA")

	var r []structs.ST_AD

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

		s := structs.ST_AD{ADId: ad_id}
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

	log.Printf("list = %s", get_ADlist_DB(a_id))

	fmt.Fprint(w, string(data))

}

// 가격 반환 핸들러
func get_pay_Handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(structs.ST_GetPay{Pay: get_pay_DB()})
	fmt.Fprint(w, string(data))
}

// const (
// 	defaultListenPort        = 8887
// 	defaultHeartbeatInterval = 600
// )

// func setupCentralSystem() ocpp16.CentralSystem {
// 	return ocpp16.NewCentralSystem(nil, nil)
// }

// type CentralSystemHandler struct {
// 	// ... your own state variables
// }

// func (handler *CentralSystemHandler) OnAuthorize(chargePointId string, request *core.AuthorizeRequest) (confirmation *core.AuthorizeConfirmation, err error) {
// 	// ... your own custom logic
// 	log.Printf("OnAuthorize %s", chargePointId)
// 	return core.NewAuthorizationConfirmation(types.NewIdTagInfo(types.AuthorizationStatusAccepted)), nil
// }

// func (handler *CentralSystemHandler) OnBootNotification(chargePointId string, request *core.BootNotificationRequest) (confirmation *core.BootNotificationConfirmation, err error) {
// 	// ... your own custom logic
// 	log.Printf("OnBootNotification %s", chargePointId)
// 	return core.NewBootNotificationConfirmation(types.NewDateTime(time.Now()), defaultHeartbeatInterval, core.RegistrationStatusAccepted), nil
// }

//////////////////////////////////////////////////////////////////////////////////////////////

// log replaces standard log
var logd *zap.SugaredLogger

var csms *Server

// 원격 선 명령
// clear cache - 충전 지역 인증 캐쉬(18p) 지우기
// ChangeAvailability - 가용성 변경(50p)	- 기기 상태 원격 명령
// ChangeConfiguration - 구성 변경 (51p)
// GetConfiguration	- 구성 가져오기(53p)
// RemoteStartTransaction - 원격 시작 트랜잭션(54p)
// RemoteStopTransaction - 원격 중지 트랜잭션(54p) -
// Reset	- 초기화 (57p)	- 진행중이면 stoptransaction 보내야함 -> 응용프로그램 or 하드웨어 다시 시작
// UnlockConnector	- 커넥터 잠금 해제(61p)	- 커넥터 케이블 고정장치가 오작동 하는 경우 분리
// DataTransfer	--- 커스텀 정보값

func Command_Handler(w http.ResponseWriter, r *http.Request) {
	if csms != nil {
		log.Print("nil")

		cp, b := csms.Load("device_01")
		if b {
			////ChangeConfiguration
			// var req Payload = v16.ChangeConfigurationReq{
			// 	Key:   "WebSocketPingInterval",
			// 	Value: "30",
			// }
			// res, err := cp.Call("ChangeConfiguration", req)
			// if err != nil {
			// 	log.Printf("aa %s", err)
			// }
			// log.Print(res)

			////RemoteStartTransaction
			// var req Payload = v16.RemoteStartTransactionReq{
			// 	IdTag:       "device_01",
			// 	ConnectorId: 1,
			// }
			// res, err := cp.Call("RemoteStartTransaction", req)
			// if err != nil {
			// 	log.Printf("aa %s", err)
			// }
			// log.Print(res)

			////RemoteStopTransaction
			var req Payload = v16.RemoteStopTransactionReq{
				TransactionId: 456789,
			}
			res, err := cp.Call("RemoteStartTransaction", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)

			////clear cache

			////ChangeAvailability

			////Reset	- 97p
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

			////UnlockConnector

			////GetConfiguration
		}

		if cp == nil {
			log.Print("cp is nil")
		}
		// log.Printf("cp2 : %s, b : %v", cp.Id, (b))
	}

}

//////////////////////////////////////////////////////////////////////////////////////////////

func main() {
	err := godotenv.Load("./info.env")
	if err != nil {
		log.Println(err)
		return
	}

	mux := pat.New()
	// mux.Get("/ws", wsHandler)

	//negroni를 쓰면 따로 명시하지 않을시 "/"에 접근하면 index.html or default.html 파일에 액세스함
	n := negroni.Classic()

	//////////////////////////////////////////////////////////////////////////////////////////////

	//관리자 페이지
	//로그인
	mux.Path("/").HandlerFunc(adminpage.LoginPage_Handler)

	// mux.Path("/login").HandlerFunc(adminpage.Login_Handler)
	mux.Path("/login").HandlerFunc(Command_Handler)

	//충전기 설치 정보
	mux.Path("/devicelist").HandlerFunc(adminpage.Devicelist_Page_Handler)
	//충전요금 단가 관리
	mux.Path("/pricemng").HandlerFunc(adminpage.Pricemng_Page_Handler)
	//소프트웨어 업데이트 관리
	mux.Path("/swmng").HandlerFunc(adminpage.Swmng_Page_Handler)
	//광고물 관리

	//충전기별 광고 관리
	mux.Path("/admng").HandlerFunc(adminpage.Admng_Page_Handler)
	//충전량 및 이력조회
	mux.Path("/chargelist").HandlerFunc(adminpage.Chargelist_Page_Handler)
	//일자별 충전량 조회
	//충전기 그래프 조회

	//충전기 장애 이력
	mux.Path("/errlist").HandlerFunc(adminpage.Errlist_Page_Handler)

	//매출실적 및 정산 조회
	//사용자 관리/권한
	//공통코드

	//////////////////////////////////관리자페이지 기능///////////////////////////////////

	//기기 등록
	mux.Path("/device/register").HandlerFunc(adminpage.Register_Device_Handler)
	//기기 수정
	mux.Path("/device/updater").HandlerFunc(adminpage.Updater_Device_Handler)
	//기기 조회
	mux.Path("/device/selecter").HandlerFunc(adminpage.Selecter_Device_Handler)
	//광고업로드
	mux.Path("/upload/ad").HandlerFunc(adminpage.Errlist_Page_Handler)
	//apk업로드
	mux.Path("/upload/apk").HandlerFunc(adminpage.Errlist_Page_Handler)

	//////////////////////////////////////////////////////////////////////////////////////////////
	// 정보값 db

	//결제정보
	mux.Path("/pay/history").HandlerFunc(payinfo.Set_Pay_Handler).Methods("POST")

	//충전시작
	mux.Path("/charging/start").HandlerFunc(payinfo.Set_charging_start_Handler).Methods("POST")

	//충전종료
	mux.Path("/charging/end").HandlerFunc(payinfo.Set_charging_end_Handler).Methods("POST")

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
	go http.ListenAndServe(":8080", n)

	//////////////////////////////////////////////////////////////////////////////////////////////

	// centralSystem := setupCentralSystem()

	// // Set callback handlers for connect/disconnect
	// centralSystem.SetNewChargePointHandler(func(chargePoint ocpp16.ChargePointConnection) {
	// 	log.Printf("new charge point %v connected", chargePoint)
	// })

	// centralSystem.SetChargePointDisconnectedHandler(func(chargePoint ocpp16.ChargePointConnection) {
	// 	log.Printf("charge point %v disconnected", chargePoint)
	// })

	// // Set handler for profile callbacks
	// handler := &CentralSystemHandler{}
	// centralSystem.SetCoreHandler(handler)
	// centralSystem.SetLocalAuthListHandler(handler)
	// centralSystem.SetFirmwareManagementHandler(handler)
	// centralSystem.SetReservationHandler(handler)
	// centralSystem.SetRemoteTriggerHandler(handler)
	// centralSystem.SetSmartChargingHandler(handler)

	// // Start central system
	// listenPort := defaultListenPort
	// log.Printf("starting central system")
	// centralSystem.Start(listenPort, "/{ws}") // This call starts server in daemon mode and is blocking
	// log.Println("stopped central system")

	//////////////////////////////////////////////////////////////////////////////////////////////

	logger, _ := zap.NewDevelopment()
	logd = logger.Sugar()
	defer logd.Sync()

	// set ocpp library's logger to zap logger
	SetLogger(logd)

	// start csms server with default configurations
	csms = NewServer()

	csms.AddSubProtocol("ocpp1.6")
	csms.SetCheckOriginHandler(func(r *http.Request) bool { return true })
	// csms.SetPreUpgradeHandler(customPreUpgradeHandler)
	csms.SetCallQueueSize(32)
	// csms.SetTimeoutConfig(ServerTimeoutConfig{OcppWait: 360000}) // 연결 제한시간

	// register charge-point-initiated action handlers
	csms.On("Authorize", AuthorizationHandler)           //	충전 시작, 끝에 사용자가 회원인지
	csms.On("BootNotification", BootNotificationHandler) //  부트
	// csms.After("BootNotification", SendChangeConfigration)
	csms.On("Heartbeat", HeartBeatHandler)                  // 하트비트
	csms.On("StartTransaction", StartTransactionHandler)    // 충전시작
	csms.On("StopTransaction", StopTransactionHandler)      // 충전 중단
	csms.On("DataTransfer", DataTransferHandler)            // 커스텀 정보값
	csms.On("MeterValues", MeterValuesHandler)              // 계량기값에 대한 정보(41p) - ChangeConfiguration을 이용, 데이터 수집간격을 구성, 데이터를 지정할수도있음
	csms.On("StatusNotification", StatusNotificationHanler) //진단 상태 알림 (43p) - 상태가 변하거나 오류가나면

	// 원격 선 명령
	// clear cache - 충전 지역 인증 캐쉬(18p) 지우기
	// ChangeAvailability - 가용성 변경(50p)	- 기기 상태 원격 명령
	// ChangeConfiguration - 구성 변경 (51p)
	// GetConfiguration	- 구성 가져오기(53p)
	// RemoteStartTransaction - 원격 시작 트랜잭션(54p)
	// RemoteStopTransaction - 원격 중지 트랜잭션(54p) -
	// Reset	- 초기화 (57p)	- 진행중이면 stoptransaction 보내야함 -> 응용프로그램 or 하드웨어 다시 시작
	// UnlockConnector	- 커넥터 잠금 해제(61p)	- 커넥터 케이블 고정장치가 오작동 하는 경우 분리
	// DataTransfer	--- 커스텀 정보값

	csms.Start("0.0.0.0:8999", "/ws/", nil)

	//연결이 되어 있는지 확인
	// log.Printf("aa : %s", strconv.FormatBool(csms.IsConnected("device_01")))
	// log.Printf("aa : %s", strconv.FormatBool(csms.IsConnected("device_02")))

	//////////////////////////////////////////////////////////////////////////////////////////////
}

// func SendChangeConfigration(cp *ChargePoint, payload Payload) {

// 	var req Payload = v16.ChangeConfigurationReq{
// 		Key:   "WebSocketPingInterval",
// 		Value: "30",
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
// 		log.Print("error parsing basic auth")
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
	csms.Store(cp)

	// cp2,b := csms.Load("device_01")

	req := p.(*v16.BootNotificationReq)
	log.Printf("\nid: %s\nBootNotification: %v", cp.Id, req)

	// csms.Store(cp) // cp 저장
	// cp2, b := csms.Load("device_01")
	// log.Printf("cp2 : %s, b : %v", cp2.Id, (b))

	rb := ocpp_db_handler.BootNotification_update(req.ChargePointModel, "1", req.FirmwareVersion)

	var res Payload

	if rb {
		res = &v16.BootNotificationConf{
			CurrentTime: time.Now().Format("2006-01-02T15:04:05.000Z"),
			Interval:    60,
			Status:      "Accepted",
		}
	} else {
		res = &v16.BootNotificationConf{
			CurrentTime: time.Now().Format("2006-01-02T15:04:05.000Z"),
			Interval:    60,
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
				Status: "Accepted",
			},
		}
	} else {
		res = v16.AuthorizeConf{ ///Accepted/ConcurrentTx, Expired, Blocked 및 Invalid 에 따라 유효, 만료, (일시적으로) 차단 또는 블랙리스트 로 표시
			IdTagInfo: &v16.IdTagInfo{
				Status: "Expired",
			},
		}
	}
	print("c")
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
	// csms.Store(cp) // cp 저장

	rb := ocpp_db_handler.StartTransaction_set_start(req.IdTag, req.ConnectorId, req.MeterStart, req.Timestamp)

	var res Payload
	Ti, err := strconv.Atoi(req.IdTag)

	if err != nil {
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

// type StopTransactionRequestTT struct {
// 	IdTag           string           `json:"idTag,omitempty" validate:"omitempty,max=20"`
// 	MeterStop       int              `json:"meterStop" validate:"required,gte=0"`
// 	Timestamp       string           `json:"timestamp" validate:"required,dateTime"`
// 	TransactionId   int              `json:"transactionId" validate:"required"`
// 	Reason          string           `json:"reason,omitempty" validate:"omitempty,reason"`
// 	TransactionData []v16.MeterValue `json:"transactionData,omitempty" validate:"omitempty,dive"`
// }

func StopTransactionHandler(cp *ChargePoint, p Payload) Payload {
	log.Print("StopTransactionHandler")

	// req := p.(*v16.StopTransactionReq)
	req := p.(*StopTransactionRequestTT)
	log.Printf("\nid: %s\nStopTransaction: %v", cp.Id, req)
	// csms.Store(cp) // cp 저장

	// v16.MeterValue
	// v16.SampledValue

	rb := ocpp_db_handler.StopTransaction_update_end(strconv.Itoa(req.TransactionId), req.Timestamp, req.MeterStop)

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

	ocpp_db_handler.Metervalue_update(strconv.Itoa(req.TransactionId), req.ConnectorId, req.MeterValue[0].Timestamp, req.MeterValue[0].SampledValue[0].Value)

	var res Payload = v16.MeterValuesConf{} //72p
	return res

}

func StatusNotificationHanler(cp *ChargePoint, p Payload) Payload {
	req := p.(*StatusNotificationReqTT)
	log.Printf("\nid: %s\nStatusNotification: %v", cp.Id, req)
	// csms.Store(cp) // cp 저장

	var res Payload

	if req.VendorErrorCode > "" {
		log.Printf("VendorErrorCode is not null")
		ocpp_db_handler.StatusNotification_device(cp.Id, req.Status, req.VendorErrorCode)
	} else {
		log.Printf("VendorErrorCode is null")
		ocpp_db_handler.StatusNotification_connector(cp.Id, req.ConnectorId, req.Status, req.ErrorCode)
	}

	res = v16.StatusNotificationConf{} //76p

	return res

}

//////////////////////////////////////////////////////////////////////////////////////////////
