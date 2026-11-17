[개요] 2023.10~2024.06 / 과충전 방지 기능 전기차 충전기.

[사용기술]
 - Go 1.19, gorilla/websocket, gorilla/mux, aliml92/ocpp(OCPP 1.6 메시지 모델),
   go-playground validator, zap 로깅, jwt-go, godotenv, go-sql-driver/mysql
 - MySQL, JWT(HS256, 만료 1시간), CORS 미들웨어

[주요기능]
 - OCPP 1.6 CSMS(중앙 시스템): 충전기와 WebSocket("ocpp1.6" 서브프로토콜) 연결
   · 충전기 → 서버: BootNotification, Authorize, Heartbeat, StartTransaction,
     StopTransaction, MeterValues, StatusNotification, DataTransfer,
     DiagnosticsStatusNotification, FirmwareStatusNotification
   · 서버 → 충전기(원격 명령 API): RemoteStart/StopTransaction, ClearCache,
     ChangeAvailability, UnlockConnector, Reset, Get/ChangeConfiguration,
     UpdateFirmware, GetDiagnostics, TriggerMessage
 - DataTransfer 확장: 결제 정보 저장(i_pay), 부분 취소(c_pay), 전기 단가 조회(g_pay),
   결제 단말(NVCAT) 설정·재시작·결제취소 명령
 - 충전 이력/일별 통계/그래프, 단가 관리, 오류 이력, 장치 등록·수정·삭제
 - 앱 APK 버전 확인·다운로드·업데이트 완료 보고
 - 관리자 REST API(JWT 인증): 대시보드, 지도, 장치 목록/상태, 로그, 충전 이력 등

[핵심 시스템 구조와 구현 방식]
 - OCPP 라이브러리를 직접 포함해(client.go/server.go/charge_point.go/messages.go)
   Call(2)/CallResult(3)/CallError(4) 프레임을 처리하고, Call 대기 큐(크기 32),
   핑퐁/타임아웃 설정, 액션별 핸들러 등록(csms.On(액션, 핸들러))을 구현.
   MeterValues, StatusNotification, StopTransaction은 현장 단말 포맷에 맞춘
   커스텀 요청 구조체(...ReqTT)로 파싱.
 - 패키지 분리(gosrc): db, ocpp_db_handler(OCPP→DB 반영), adminpage/adminpagedbconn
   (관리자 조회·CRUD), payinfo(결제/충전 시작·종료 보고), updateapk(APK 배포),
   adminjwt(토큰), ocpplogger.
 - 인증: /login만 공개, 나머지는 jwtMiddleware로 Bearer 토큰 검증 → 컨텍스트에 사용자 주입.
 - 검색 API는 허용 컬럼 화이트리스트(isSafeColumnName)로 SQL 컬럼명 주입을 방지.
 - 연결된 충전기는 csms.Store/Load로 관리, 관리자 명령은 해당 연결에 Call 후 결과를 REST로 반환.

[데이터 흐름]
 충전기(Android 키오스크/PLC) --WebSocket OCPP 1.6--> Go CSMS
   BootNotification: DB에서 장치 enabled 확인 → Accepted/Pending/Rejected 응답
   Authorize(결제 번호) → DB 조회 → Accepted/Invalid
   Start/StopTransaction, MeterValues → 충전 테이블 INSERT/UPDATE(계량값 kWh 변환)
   StatusNotification 오류 → 오류 테이블 저장
 관리자 웹 --REST(JWT)--> Go --Call--> 충전기 (원격 시작·중지·재부팅·펌웨어)
 키오스크 앱 --REST--> 결제 이력/충전 시작·종료/APK 업데이트


<img width="2992" height="2992" alt="KakaoTalk_20250417_163948276_04" src="https://github.com/user-attachments/assets/dc649e78-5abb-441f-87e9-3b29e17c5c87" />


<img width="2992" height="2992" alt="KakaoTalk_20250417_163948276_02" src="https://github.com/user-attachments/assets/1af46489-74e9-4762-8904-6264f6773840" />



 <img width="991" height="1244" alt="전기차 충전기 DB 다이어그램" src="https://github.com/user-attachments/assets/12fdb484-d7b3-4784-ad68-d8d2f6ea041f" />

