var form_ = document.createElement('form');

//충전기 설치 정보
function menu_devicelist_func(){
    form_.setAttribute("action","/devicelist");
    document.body.appendChild(form_);
    form_.submit();
}

//충전요금 단가 관리
function menu_pricemng_func(){
    form_.setAttribute("action","/pricemng");
    document.body.appendChild(form_);
    form_.submit();
}

//소프트웨어 업데이트 관리
function menu_swmng_func(){
    form_.setAttribute("action","/swmng");
    document.body.appendChild(form_);
    form_.submit();
}

//광고물 관리

//충전기별 광고 관리
function menu_admng_func(){
    form_.setAttribute("action","/admng");
    document.body.appendChild(form_);
    form_.submit();
}

//충전량 및 이력조회
function menu_chargelist_func(){
    // form_.attr("method","post");
    form_.setAttribute("action","/chargelist");
    document.body.appendChild(form_);
    form_.submit();
}

//일자별 충전량 조회
//충전기 그래프 조회

//충전기 장애 이력
function menu_errlist_func(){
    form_.setAttribute("action","/errlist");
    document.body.appendChild(form_);
    form_.submit();
}

//매출실적 및 정산 조회

//사용자 관리/권한
//공통코드




