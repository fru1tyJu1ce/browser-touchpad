const socket = new WebSocket("ws://192.168.188.20:62723/ws");


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

function sendControlInstruction(dx, dy, type) { // <- sending controll instructions to the server 
  msg = JSON.stringify({
    Type: type,
    dx: ~~dx,
    dy: ~~dy
  })
  socket.send(msg);
}

function init() { // <- initialisation
  document.getElementById('version').innerText = 'v.010';
  canvas = document.getElementById("mainCanvas");
  canvas.width = document.body.clientWidth;
  canvas.height = document.body.clientHeight;
  canvasW = canvas.width;
  canvasH = canvas.height;
}
init();

/* Settings */

/* global variables for settings & configuaration */
var accelerator = 2;
var nightMode = false;

/* methods for settings & configuaration */
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
activateNightMode() // <- autostart nightmode by default 

function changeSens(newSens) {
  accelerator = newSens;
  document.getElementById("valSens").innerHTML = newSens;
}

/* input detection */

/* global variables for the input detection */
const CLICKDETECTIONTRASHHOLD = 99;
const TOGGLETRASHHOLD = 400;
const RELOADTRASHHOLD = 10000;

var lastX = 0;
var lastY = 0;

var lastTouchStart = 0;
var lastTouchEnd = 0;

var lastClick = 0;

var lastTouchstartTouchpointsCount = 0;

var scroller = false;
var toggle = false;
var rightClick = false

var scrollArea = document.getElementById('sclr');
var rightClickArea = document.getElementById('rghtClck');

/* functions for the input detection */
function clickDetection() {
  if (lastTouchEnd - lastTouchStart
    <= CLICKDETECTIONTRASHHOLD) {
    var type;
    if (lastTouchstartTouchpointsCount === 2 || rightClick) {
      type = "rightClick";
      rightClick = false;
    }
    else {
      type = "click";
      lastClick = new Date().getTime();
    }
    sendControlInstruction(0, 0, type);
  }
}

function toggleDetection() {
  let now = new Date().getTime();
  if (!toggle && now - lastClick <= TOGGLETRASHHOLD) {
    toggle = true;
    sendControlInstruction(0, 0, "toggle");
  } else if (toggle) {
    toggle = false;
    sendControlInstruction(0, 0, "toggleUp");
  }
}

/* touchevents */
document.ontouchstart = (event) => {
  let now = new Date().getTime();
  lastX = event.touches[0].pageX;
  lastY = event.touches[0].pageY;
  lastTouchstartTouchpointsCount = event.touches.length;
  toggleDetection();
  lastTouchStart = now;
};

document.ontouchend = (event) => {
  lastTouchEnd = new Date().getTime();
  lastX = event.pageX;
  lastY = event.pageY;
  if (toggle) toggleDetection();
  else clickDetection();
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
    sendControlInstruction(dx, dy, type);
  }
})();

scrollArea.ontouchstart = function (event) {
  scroller = true;
}

scrollArea.ontouchend = function (event) {
  scroller = false;
}

rightClickArea.ontouchstart = function (event) {
  rightClick = true;
}

