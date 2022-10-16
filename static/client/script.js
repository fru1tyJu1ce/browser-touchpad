const socket = new WebSocket("ws://192.168.188.23:37733/ws");

console.log('attempting websocket connection');

socket.onopen = () => {
  console.log('succesfully connected');
  socket.send('client connected');
}

socket.onclose = (event) => {
  console.log('socket closed connection: ', event);
}

socket.onmessage = (msg) => {
  console.log(msg);
}

socket.onerror = (error) => {
  console.log('socket error: ', error);
}

function sendConrolInstruction(dx, dy, type) {
  msg = JSON.stringify({
    Type: type,
    dx: ~~dx,
    dy: ~~dy
  })
  socket.send(msg);
}

function init() { // everything 4 initialisation
  canvas = document.getElementById("mainCanvas");
  canvas.width = document.body.clientWidth;
  canvas.height = document.body.clientHeight;
  canvasW = canvas.width;
  canvasH = canvas.height;
}
init();

//Settings

const STANDARDACCELERATOR = 2;
var accelerator = STANDARDACCELERATOR;
var nightMode = false;

function activateNightMode() {
  if (!nightMode) {
    document.body.style.backgroundColor = "black";
    document.body.style.color = "white";
    nightMode = true;
  }
  else {
    document.body.style.backgroundColor = "white";
    document.body.style.color = "black";
    nightMode = false;
  }
}
activateNightMode()// aoutostart nightmode by default 

function changeSens(newSens) {
  accelerator = newSens;
  document.getElementById("valSens").innerHTML = newSens;
}

// Global vars needed for touch detection

const CLICKDETECTIONTRASHHOLD = 99;
const TOGGLETRASHHOLD = 200;
const RELOADTRASHHOLD = 500;


var lastX = 0;
var lastY = 0;

var lastTouchStart = 0;
var lastTouchEnd = 0;

var lastClick = 0;

var lastTouchstartTouchpointsCount = 0;

var scroller = false;

// touch detection
document.ontouchstart = (event) => {
  let now = new Date().getTime();
  // if (lastTouchStart != 0 && now - lastTouchStart >= RELOADTRASHHOLD) document.reload();
  lastTouchStart = now;
  lastX = event.touches[0].pageX;
  lastY = event.touches[0].pageY;
  lastTouchstartTouchpointsCount = event.touches.length;
  document.getElementById('action').innerText = lastClick + " / "+ now;
  if(now - lastClick < TOGGLETRASHHOLD + 400){
    sendConrolInstruction(0, 0, "toggle");
  }
};

var touchArea = document.getElementById('sclr');
touchArea.ontouchstart = function (event) {
  scroller = true;
}
touchArea.ontouchend = function (event) {
  endConrolInstruction(0, 0, "toggleUp");
  scroller = false;
}

function clickDetection(event) {
  document.getElementById('action').innerText = lastTouchstartTouchpointsCount;
  if (lastTouchEnd - lastTouchStart
    <= CLICKDETECTIONTRASHHOLD) {
    var type;
    if (lastTouchstartTouchpointsCount === 2) type = "rightClick";
    else {
      type = "click";
      lastClick = new Date().getTime();
    }
    sendConrolInstruction(0, 0, type);
  }
}

document.ontouchend = (event) => {
  //document.getElementById('action').innerText = "touchendStart";
  document.getElementById('action').innerText = lastTouchEnd + " / " + event.pageX;

  lastTouchEnd = new Date().getTime();
  lastX = event.pageX;
  lastY = event.pageY;

  clickDetection(event);
};

(function () {
  document.ontouchmove = handleMouseMove;
  function handleMouseMove(event) {
    var eventDoc, doc, body;

    event = event || window.event; // IE-ism

    if (event.touches[0].pageX == null && event.clientX != null) {
      eventDoc = (event.target && event.target.ownerDocument) || document;
      doc = eventDoc.documentElement;
      body = eventDoc.body;

      event.touches[0].pageX = event.clientX +
        (doc && doc.scrollLeft || body && body.scrollLeft || 0) -
        (doc && doc.clientLeft || body && body.clientLeft || 0);
      event.touches[0].pageY = event.clientY +
        (doc && doc.scrollTop || body && body.scrollTop || 0) -
        (doc && doc.clientTop || body && body.clientTop || 0);
    }

    dx = (event.touches[0].pageX - lastX) * accelerator;
    dy = (event.touches[0].pageY - lastY) * accelerator;

    lastX = event.touches[0].pageX;
    lastY = event.touches[0].pageY;

    var type;
    if (event.targetTouches.length === 2 || scroller) {
      type = "scroll"
      if (dy < 0) dy = -1;
      if (dy > 0) dy = 1;
    }
    else type = "movement";

    sendConrolInstruction(dx, dy, type);
    document.getElementById('action').innerText = event.targetTouches.length;
  }
})();

document.getElementById('action').innerText = 'v.09';
