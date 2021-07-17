const socket = new WebSocket("ws://192.168.188.23:46551/ws");  


console.log('attempting websocket connection');

socket.onopen = () => {
  console.log('succesfully connected');
  //socket.send('client connected');
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



// mouse + touch detection

const diffClick = 400;
const diffDownUp = 190;
const diffMouseMovement = 100;

var lastMouseMove = new Date(); 
var lastClick = new Date();
var toggle = false;

var lastX = 0;
var lastY = 0;

var lastDX = 0;
var lastDY = 0;

var lastMouseDown = 0;
var lastMouseUp = 0;




//Mousemoement
(function () {
  document.onmousemove = handleMouseMove;
  function handleMouseMove(event) {
    var eventDoc, doc, body;

    //console.log('mouse movement detected');

    event = event || window.event; // IE-ism

    // If pageX/Y aren't available and clientX/Y are,
    // calculate pageX/Y - logic taken from jQuery.
    // (This is to support old IE)
    if (event.pageX == null && event.clientX != null) {
      eventDoc = (event.target && event.target.ownerDocument) || document;
      doc = eventDoc.documentElement;
      body = eventDoc.body;

      event.pageX = event.clientX +
        (doc && doc.scrollLeft || body && body.scrollLeft || 0) -
        (doc && doc.clientLeft || body && body.clientLeft || 0);
      event.pageY = event.clientY +
        (doc && doc.scrollTop || body && body.scrollTop || 0) -
        (doc && doc.clientTop || body && body.clientTop || 0);


    }
    document.getElementById("indec").innerHTML = "mouse input detected";
    document.getElementById("indecx").innerHTML = event.pageX;
    document.getElementById("indecy").innerHTML = event.pageY;

    
    if ((new Date() - lastMouseMove) >= diffMouseMovement) {
      lastX = event.pageX;
      lastY = event.pageY;
      lastMouseMove = new Date();
    }
    

    let x = event.pageX;
    let y = event.pageY;
    let dx = (x - lastX)*2;
    let dy = (y - lastY)*2;



    socket.send(JSON.stringify({
      Type: "mouseMove",
      dx: dx,
      dy: dy
    }));


    lastX = x;
    lastY = y;
  }
})();

/*
function isTouchScreendevice() {
  return 'ontouchstart' in window || navigator.maxTouchPoints;
};

if (isTouchScreendevice()) {
  alert("I am a touch screen device")
}
*/

function touchHandler(event) {
  var touches = event.changedTouches,
    first = touches[0],
    type = "";
  switch (event.type) {
    case "touchstart": type = "mousedown"; break;
    case "touchmove": type = "mousemove"; break;
    case "touchend": type = "mouseup"; break;
    default: return;
  }

  // initMouseEvent(type, canBubble, cancelable, view, clickCount, 
  //                screenX, screenY, clientX, clientY, ctrlKey, 
  //                altKey, shiftKey, metaKey, button, relatedTarget);

  var simulatedEvent = document.createEvent("MouseEvent");
  simulatedEvent.initMouseEvent(type, true, true, window, 1,
    first.screenX, first.screenY,
    first.clientX, first.clientY, false,
    false, false, false, 0/*left*/, null);

  first.target.dispatchEvent(simulatedEvent);
  //event.preventDefault();
}

function click() {
  document.getElementById("mouseC").innerHTML = "click detected";

  lastClick = new Date();
  socket.send(JSON.stringify({
    type: "click"
  }));
}


window.addEventListener("mouseup", function (event) {
  document.getElementById("mouseB").innerHTML = "mouse up";

  if(toggle){
    toggle = false;
    socket.send(JSON.stringify({
      type: "toggleUp"
    }));
  }

  lastMouseUp = new Date();
  let diff = (lastMouseUp - lastMouseDown);
  if (diff <= diffDownUp) {
    click();
  }

}, false);


window.addEventListener("mousedown", function (event) {
  document.getElementById("mouseB").innerHTML = "mouse down";
  lastMouseDown = new Date();

  let diff = (new Date - lastClick);
  if (diff <= diffClick) {
    toggle = true;
    socket.send(JSON.stringify({
      type: "toggle"
    }));
  }
  console.log('mouseDown ' + lastMouseDown)
}, false);




function init() {
  document.addEventListener("touchstart", touchHandler, false);
  document.addEventListener("touchmove", touchHandler, false);
  document.addEventListener("touchend", touchHandler, false);
  document.addEventListener("touchcancel", touchHandler, false);
}

init()

/*

document.getElementById("swtch").addEventListener('change', e => {

  if(e.target.checked){
    document.getElementById("demo").innerHTML = "Hello World";
  }

  

});

*/
