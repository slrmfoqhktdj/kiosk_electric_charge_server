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

///////////////////////////////

// 페이지 구조체
type Page struct {
	Title string
}

type DeviceListPage struct { // 충전기 설치정보 페이지
	Title string
	Page  int
	Items []DeviceListItem
}

type PricemngPage struct { // 충전요금 단가 관리 페이지
	Title string
	Page  int
	Items []PriceListItem
}

type SwmngPage struct { // 소프트웨어 업데이트 관리 페이지
	Title string
	Page  int
	Items []SwListItem
}

type AdmngPage struct { // 충전기별 광고 관리 페이지
	Title string
	Page  int
	Items []AdListItem
}

type ChargeListPage struct { // 충전량 및 이력조회 관리 페이지
	Title string
	Page  int
	Items []ChargeListItem
}

type ErrListPage struct { // 충전기 장애 이력 관리 페이지
	Title string
	Page  int
	Items []ErrListItem
}

type DeviceregPage struct {
	Title string
	Page  int
}

///////////////////////////////		// 페이지 정보 구조체

type DeviceListItem struct { // 충전기 설치정보 구조체
	Id        string
	Latitude  string
	Longitude string
	Addr      string
	S_Version string
	H_C_M     string
	H_P_M     string
	H_B_M     string
	State     string
}

type PriceListItem struct { // 충전요금 단가 관리 구조체
	Id    string
	Price string
	Time  string
}

type SwListItem struct { // 소프트웨어 업데이트 관리 구조체
	Id      string
	Version string
	Time    string
}

type AdListItem struct { // 충전기별 광고 관리 구조체
	D_Id  string
	AD_Id string
	Time  string
}

type ChargeListItem struct { // 충전량 및 이력조회 관리 구조체
	Pay_Num  string
	D_Id     string
	C_Time   string
	C_Amount string
	C_Pay    string
	Pay_By   string
	C_Num    string
	Pay_Day  string
	C_Time_S string
	C_Time_E string
}

type ErrListItem struct { // 충전기 장애 이력 관리 구조체
	D_Id   string
	O_Time string
	R_Time string
	Cause  string
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
