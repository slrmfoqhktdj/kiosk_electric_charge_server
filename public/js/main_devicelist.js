const slcBtn = document.querySelector('#device-slc-btn-id');
const regBtn = document.querySelector('#device-reg-btn-id');
const updBtn = document.querySelector('#device-upd-btn-id');

//다이어로그 버튼
const SendBtn = document.querySelector('#device-reg-dialog-confirm-btn');
const CloseBtn = document.querySelector('#device-reg-dialog-close-btn');

const dialogOne = document.querySelector('#device-reg-dialog');

var form = document.querySelector('#device-reg-dialog-form');
var slc_form = document.querySelector('#device-slc-text-form-id');

var select_device_id = '';
var select_device_row;
var dialog_type = false;






//조회
slcBtn.addEventListener('click', () => {
    var device_name_select = $('input[name=device-name-select]').val();
    var device_state_select = $('input[name=device-state-select]').val();
    var device_addr_select = $('input[name=device-addr-select]').val();


    if(''.concat(device_name_select,device_state_select,device_addr_select).trim() != ''){  // 문자열을 다합친것이 공란일때
        menu_deviceslc_func()
    }else{
        window.alert('검색란 중 한가지 이상 채워주세요.')
    }

})

//등록
regBtn.addEventListener('click', () => {
    dialog_type = false;
    $('#device-reg-dialog-confirm-btn').text('등록')

    $('input[name=device-name]').val('');
    $('input[name=device-state').val('');
    $('input[name=device-latitude').val('');
    $('input[name=device-longitude').val('');
    $('input[name=device-addr]').val('');
    $('input[name=device-sw-v]').val('');
    $('input[name=device-charger-m]').val('');
    $('input[name=device-pay-m]').val('');
    $('input[name=device-hw-m]').val('');

    dialogOne.showModal();
})

//수정
updBtn.addEventListener('click', () => {
    dialog_type = true;
    $('#device-reg-dialog-confirm-btn').text('수정')

    if(select_device_id != ''){
        $('input[name=device-name]').val(select_device_row.children().eq(0).text());
        $('input[name=device-state]').val(select_device_row.children().eq(1).text());
        $('input[name=device-latitude').val(select_device_row.children().eq(2).text());
        $('input[name=device-longitude').val(select_device_row.children().eq(3).text());
        $('input[name=device-addr]').val(select_device_row.children().eq(4).text());
        $('input[name=device-sw-v]').val(select_device_row.children().eq(5).text());
        $('input[name=device-charger-m]').val(select_device_row.children().eq(6).text());
        $('input[name=device-pay-m]').val(select_device_row.children().eq(7).text());
        $('input[name=device-hw-m]').val(select_device_row.children().eq(8).text());

        dialogOne.showModal();
    }else{
        window.alert('수정할 기기를 선택해주세요.')
    }
})





$('table tbody tr').click(function(event){
    $('table tbody tr').css("background-color","rgba(0, 0, 0, 0)")      // 원래 색으로 되돌리기
    
    if(select_device_id != $(this).children().first().text()){
        $(this).css("background-color","rgba(0, 0, 0, 0.2)")                // 선택한 행 색 바꾸기
        select_device_row = $(this);
        select_device_id =  select_device_row.children().first().text();
    }else{
        select_device_row = null;
        select_device_id =  '';
    }
    
})









dialogOne.addEventListener('click', function (event) {
    const target = event.target;
    const rect = target.getBoundingClientRect();
    if (rect.left > event.clientX || rect.right < event.clientX ||
      rect.top > event.clientY || rect.bottom < event.clientY) {
        dialogOne.close();
    }
})



SendBtn.addEventListener('click',()=> {
    var name = $('input[name=device-name]').val();
    var state = $('input[name=device-state').val();
    var latitude = $('input[name=device-latitude').val();
    var longitude = $('input[name=device-longitude').val();
    var addr = $('input[name=device-addr]').val();
    var sw = $('input[name=device-sw-v]').val();
    var charger = $('input[name=device-charger-m]').val();
    var pay = $('input[name=device-pay-m]').val();
    var hw = $('input[name=device-hw-m]').val();

    if(name.trim() != '' &&
    latitude.trim() != '' &&
    longitude.trim() != '' &&
    addr.trim() != '' &&
    sw.trim() != '' &&
    charger.trim() != '' &&
    pay.trim() != '' &&
    hw.trim() != '' 
    ){
        if(dialog_type){
            menu_deviceupd_func()       // 수정
        }else{
            menu_devicereg_func();      // 등록
        }
        
    }else{
        window.alert('빈 곳을 채워주세요.')
    }

})

CloseBtn.addEventListener('click',()=> {
    dialogOne.close();
})

function menu_devicereg_func(){
    form.setAttribute("method","post");
    form.setAttribute("action","/device/register");


    form.submit();
}

function menu_deviceupd_func(){
    form.setAttribute("method","post");
    form.setAttribute("action","/device/updater");


    form.submit();
}

function menu_deviceslc_func(){
    slc_form.setAttribute("method","post");
    slc_form.setAttribute("action","/device/selecter");


    slc_form.submit();
}
