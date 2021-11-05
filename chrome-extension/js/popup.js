function set_service_img() {
    var img = document.getElementById('img')
    if (bg.service_is_running()) {
        img.src = "images/online.png"
    } else {
        img.src = "images/offline.png"
    }  
}
document.addEventListener('DOMContentLoaded', function() {
    var link = document.getElementById('radio1')
    var bg = chrome.extension.getBackgroundPage()
    link.addEventListener('click', function() {
        bg.set_status_on()
        setTimeout(set_service_img, 5000)
    })
})

document.addEventListener('DOMContentLoaded', function() {
    var link = document.getElementById('radio0')
    var bg = chrome.extension.getBackgroundPage()
    link.addEventListener('click', function() {
        bg.set_status_off()
        setTimeout(set_service_img, 5000)
    })
})

document.addEventListener('DOMContentLoaded', function() {
    var link = document.getElementById('img')
    var bg = chrome.extension.getBackgroundPage()
    set_service_img()
})

document.getElementById('gw_div').innerHTML = "零信任网关客户端"

var radio0 = document.getElementById('radio0')
var radio1 = document.getElementById('radio1')
var bg = chrome.extension.getBackgroundPage()
if (bg.get_status() === "on") {
    radio0.checked = false
    radio1.checked = true
} else {
    radio0.checked = true
    radio1.checked = false
}
