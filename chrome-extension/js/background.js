chrome.webRequest.onBeforeRequest.addListener(
  function(details) {
  },
  {urls: ["<all_urls>"]},
  ["blocking"])

chrome.webRequest.onBeforeSendHeaders.addListener(
  function(details){
    if (localStorage.status !== "on") {
      return 
    }

    var bg = chrome.extension.getBackgroundPage();
    var data = localStorage.data 
    data = data?data: '{}'      
    details.requestHeaders.push({
      name: REQUESTHEADER_NAME,
      value: data
    })
    
    return {requestHeaders: details.requestHeaders}
  },
  {
    urls: [
      "<all_urls>"
    ]
  },
  [
    "blocking",
    "requestHeaders"
  ]
)

function set_status_on() {
  if (!check_data()) {
    return
  }

  localStorage.status = "on"
  ws.open(false, true)
}

function set_status_off() {
  localStorage.status = "off"
  ws.close()
}

function get_status() {
  var bg = chrome.extension.getBackgroundPage()
  bg.console.log(localStorage)
  if (localStorage.status !== "on") {
    localStorage.status = "off"
  }
  return localStorage.status
}

function service_is_running() {
  if (localStorage.running != "ok") {
    return false
  }
  return true
}

function check_data() {
  if ( localStorage.data.length == 0) {
    var options = {
      type:"basic",
      title:"提示",
      message:"没有设置数据选项",
      iconUrl:"images/warn.jpg",
    };

    var notification = chrome.notifications.create(options, null)
    return false
  } else {
    return true
  }
}

var cpuInfo;
var memoryInfo;
var storageInfo;

chrome.system.cpu.getInfo(function(info){
  cpuInfo = info;
});

chrome.system.memory.getInfo(function(info){
  memoryInfo = info;
});

chrome.system.storage.getInfo(function(info){
  storageInfo = info;
});

// init
var bg = chrome.extension.getBackgroundPage()
if (localStorage.data ==undefined) {
  localStorage.data = ''
}
localStorage.running = 'error'

var ws = new ReconnectingWebSocket(WS_URI, null, {debug: true, reconnectInterval: 1000, automaticOpen: false})
if (localStorage.status === "on") {
  if(!check_data()) {
    retrun
  }
  ws.open(false, true)
}

ws.onopen = function(e) {
  var dev = { 'cpu': cpuInfo, 'memory': memoryInfo, 'storage': storageInfo }
  var obj = JSON.parse(localStorage.data)
  obj['dev'] =  dev
  
  var bg = chrome.extension.getBackgroundPage();
  bg.console.log(JSON.stringify(obj))
  ws.send(JSON.stringify(obj))
  localStorage.running = 'ok'
}

ws.onmessage = function(e) {
  ws.send(localStorage.data)
  localStorage.running = 'ok'
}

ws.onerror = function(e) {
  localStorage.running = 'error'
}

ws.onclose = function(e) {
  localStorage.running = 'error'
}
