package structs

type ST_AD struct {
	ADId string
}

type ST_GetPay struct {
	Pay float64
}

///////////////////////////////

type ST_GetVersion struct {
	Version string
}

type ContextStruct struct{}

///////////////////////////////

// Post data 구조체
type LoginData struct { //로그인 데이터
	Email string `json:"email"`
	PWD   string `json:"pwd"`
}

type DeviceGraphTimeData struct { //충전기 그래프 조회 데이터
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
	DeviceID  string `json:"device_id"`
}

///////////////////////////////

// 페이지 구조체
type Page struct {
	Title string
}

type LoginPage struct { // 사용자 로그인 페이지
	Items LoginDataForm `json:"item"`
}

type DashboardPage struct { // 대쉬보드 페이지
	List  []DeviceListItem       `json:"list"`
	State []DeviceStateListItem  `json:"state"`
	Price []PriceListItem        `json:"price"`
	Daily []ChargeDailyListItem  `json:"daily"`
	Graph []ChargeDailyGraphItem `json:"graph"`
	Maps  []MapListItem          `json:"maps"`
}

type MapListPage struct { // 충전기 위치정보 페이지
	Items []MapListItem `json:"item"`
}

type DeviceListPage struct { // 충전기 설치정보 페이지
	Count int              `json:"count"`
	Items []DeviceListItem `json:"item"`
}

type DeviceStateListPage struct { // 충전기 상태정보 페이지
	Count int                   `json:"count"`
	Items []DeviceStateListItem `json:"item"`
}

type PricemngPage struct { // 충전요금 단가 관리 페이지
	Count int             `json:"count"`
	Items []PriceListItem `json:"item"`
}

type UpdatemngPage struct { // 업데이트 관리 페이지
	Count int                    `json:"count"`
	Items []DeviceUpdateListItem `json:"item"`
}

type LogmngPage struct { // 로그 관리 페이지
	Count int                 `json:"count"`
	Items []DeviceLogListItem `json:"item"`
}

type ChargeListPage struct { // 충전량 및 이력조회 관리 페이지
	Count int              `json:"count"`
	Items []ChargeListItem `json:"item"`
}

type ChargeDailyListPage struct { // 일자별 충전량 조회 페이지
	Count int                   `json:"count"`
	Items []ChargeDailyListItem `json:"item"`
}

type ChargeDailyGraphPage struct { // 충전기 그래프 조회 페이지
	ListItem []ChargeDailyGraphListItem `json:"listitem"`
	Items    []ChargeDailyGraphItem     `json:"item"`
}

type ErrListPage struct { // 충전기 장애 이력 관리 페이지
	Count int           `json:"count"`
	Items []ErrListItem `json:"item"`
}

// type SwmngPage struct { // 소프트웨어 업데이트 관리 페이지
// 	Title string
// 	Page  int
// 	Items []SwListItem
// }

type AdmngPage struct { // 충전기별 광고 관리 페이지
	Title string
	Page  int
	Items []AdListItem
}

type DeviceregPage struct {
	Title string
	Page  int
}

// /////////////////////////////		// 페이지 정보 구조체

type LoginDataForm struct { // 로그인 시도 대답 폼
	Email string `json:"email"`
	Name  string `json:"name"`
}

type DeviceIUDDataForm struct { // 기기 추가, 수정, 삭제 결과 대답 구조체
	Items bool `json:"item"`
}

type MapListPositionItem struct { // 충전기 위치정보 위경도 구조체
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type MapListItem struct { // 충전기 위치정보 구조체
	Id       string              `json:"id"`
	Position MapListPositionItem `json:"position"`
	Addr     string              `json:"addr"`
}

type DeviceListItem struct { // 충전기 설치정보 구조체
	Id        string `json:"id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Addr      string `json:"addr"`
	S_Version string `json:"s_version"`
	H_C_M     string `json:"h_c_m"`
	H_P_M     string `json:"h_p_m"`
	H_B_M     string `json:"h_b_m"`
	H_W_M     string `json:"h_w_m"`
	I_Time    string `json:"i_time"`
	U_Time    string `json:"u_time"`
	Enabled   string `json:"enabled"`
}

type DeviceStateListItem struct { // 충전기 상태정보 구조체
	Id      string  `json:"id"`
	Conn    bool    `json:"conn"`
	State   *string `json:"state"`
	State_e *string `json:"state_e"`
	V_Id    *string `json:"v_id"`
	V_Error *string `json:"v_error"`
	E_Info  *string `json:"e_info"`
	U_Time  *string `json:"u_time"`
	F_State *string `json:"f_state"`
	D_State *string `json:"d_state"`
}

type PriceListItem struct { // 충전요금 단가 관리 구조체
	Id     string `json:"id"`
	Price  string `json:"price"`
	Time   string `json:"time"`
	Writer string `json:"writer"`
}

type DeviceUpdateListItem struct {
	Id        string  `json:"id"`
	Addr      string  `json:"addr"`
	S_Version string  `json:"s_version"`
	F_State   *string `json:"f_state"`
	I_Time    string  `json:"i_time"`
	U_Time    string  `json:"u_time"`
}

type DeviceLogListItem struct {
	Id      string  `json:"id"`
	Addr    string  `json:"addr"`
	D_State *string `json:"d_state"`
}

type ChargeListItem struct { // 충전량 및 이력조회 관리 구조체
	Payid     string  `json:"payid"`
	D_id      *string `json:"d_id"`
	C_pay     *string `json:"c_pay"`
	Can_pay   *string `json:"can_pay"`
	Pay_by    *string `json:"pay_by"`
	C_num     *string `json:"c_num"`
	C_c_num   *string `json:"c_c_num"`
	Pay_day   *string `json:"pay_day"`
	C_pay_day *string `json:"c_pay_day"`
	C_conn    *string `json:"c_conn"`
	C_time    *string `json:"c_time"`
	C_amount  *string `json:"c_amount"`
	C_time_s  *string `json:"c_time_s"`
	C_time_e  *string `json:"c_time_e"`
	C_time_u  *string `json:"c_time_u"`
}

type ChargeDailyListItem struct { // 일자별 충전량 조회 구조체
	Day    string  `json:"day"`
	Amount *string `json:"amount"`
	Count  *string `json:"count"`
	Pay    *string `json:"pay"`
}

type ChargeDailyGraphListItem struct { // 충전기 그래프 기기 리스트 조회 구조체
	Id string `json:"id"`
}

type ChargeDailyGraphItem struct { // 충전기 그래프 조회 구조체
	Time   string  `json:"time"`
	Amount *string `json:"amount"`
	Count  *string `json:"count"`
	Pay    *string `json:"pay"`
}

type ErrListItem struct { // 충전기 장애 이력 관리 구조체
	Id      string `json:"id"`
	D_Id    string `json:"d_id"`
	D_Error string `json:"d_error"`
	V_Id    string `json:"v_id"`
	V_Error string `json:"v_error"`
	E_Info  string `json:"e_info"`
	O_Time  string `json:"o_time"`
}

// type SwListItem struct { // 소프트웨어 업데이트 관리 구조체
// 	Id      string
// 	Version string
// 	Time    string
// }

type AdListItem struct { // 충전기별 광고 관리 구조체
	D_Id  string
	AD_Id string
	Time  string
}

////////////////////////////////////////////////////////////////////////////////////////////
//ocpp

type ST_Error struct { // DataTransfer 에러
	MessageType int
}

type ST_G_Pay struct { //	DataTransfer 시세 반환
	MessageType int
	Data        ST_GetPay
}

type ST_ADList struct { // DataTransfer 광고 리스트 반환
	MessageType int
	Data        []ST_AD
}

type ST_Pay_Info struct { // DataTransfer 결제 정보
	MessageType int
	Data        []PayItem
}

type PayItem struct {
	Pay_Num    string `json:"Pay_Num"`
	Device_Id  string `json:"Device_Id"`
	Charge_Pay string `json:"Charge_Pay"`
	Pay_By     string `json:"Pay_By"`
	Credit_Num string `json:"Credit_Num"`
	Pay_Day    string `json:"Pay_Day"`
}

// can_pay string, c_pay_day string, c_c_num string, pay_num string
type CanclePayItem struct {
	Pay_Num    string `json:"Pay_Num"`
	Cancle_Pay string `json:"Cancle_Pay"`
	C_Pay_Day  string `json:"C_Pay_Day"`
	C_C_Num    string `json:"C_C_Num"`
}
