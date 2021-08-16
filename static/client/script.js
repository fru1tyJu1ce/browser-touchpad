const socket = new WebSocket("ws://192.168.188.23:38741/ws");  

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

function sendMouse(type, dx, dy) {
  socket.send(JSON.stringify({
    Type: type,
    dx: ~~dx,
    dy: ~~dy
  }));

}


//touch detection

var mouse = false; //true if touched on touchpad obj


const diffClick = 400;
const diffDownUp = 400;
const diffMouseMovement = 100;

var lastMouseMove = new Date();
var lastMouseScroll = new Date();

var lastClick = new Date();
var toggle = false;

var lastX = 0;
var lastY = 0;

var lastDX = 0;
var lastDY = 0;

var lastMouseDown = 0;
var lastMouseUp = 0;

var lastMouseRightDown = 0;
var lastMouseRightUp = 0;

function mouseStatus(n) {
  mouse = n;
}

//Mousemovemmouseent
var src = document.getElementById("touchpad");

function moveMose(e) {
  if ((new Date() - lastMouseMove) >= diffMouseMovement) {
    lastX = e.changedTouches[0].pageX;
    lastY = e.changedTouches[0].pageY;
    lastMouseMove = new Date();
  }


  let x = e.changedTouches[0].pageX;
  let y = e.changedTouches[0].pageY;
  let dx = (x - lastX) * 2;
  let dy = (y - lastY) * 2;
  if (mouse) {
    sendMouse("mouseMove", dx, dy);
  }
}

function scroll(e) {
  if ((new Date() - lastMouseScroll) >= diffMouseMovement) {
    lastY = e.changedTouches[0].pageY;
    lastMouseScroll = new Date();
  }

  let y = e.changedTouches[0].pageY;
  let dy = (y - lastY) * 2;
  // if (mouse) {
  sendMouse("scroll", 0, dy);
  // }
}

window.addEventListener('touchmove', function (e) {

  document.getElementById("multiT").innerHTML = e.changedTouches.length;
  document.getElementById("indec").innerHTML = "mouse input detected";
  // Iterate through the touch points that have moved and log each
  // of the pageX/Y coordinates. The unit of each coordinate is CSS pixels.

  /*
  var i;
  for (i=0; i < e.changedTouches.length; i++) {
    console.log("touchpoint[" + i + "].pageX = " + e.changedTouches[i].pageX);
    console.log("touchpoint[" + i + "].pageY = " + e.changedTouches[i].pageY);
  }
  */

  if (e.changedTouches.length == 1) {
    moveMose(e);

  } else if (e.changedTouches.length == 2) {
    scroll(e);
  }



  document.getElementById("indecx").innerHTML = e.changedTouches[0].pageX;
  document.getElementById("indecy").innerHTML = e.changedTouches[0].pageY;

}, false);

window.addEventListener('touchstart', function (e) {
  if (mouse && e.touches.length == 1) {
    lastMouseDown = new Date();
  } else if (mouse && e.touches.length == 2) {
    lastMouseRightDown = new Date();
  }
}, false);

window.addEventListener('touchend', function (e) {
  if (e.touches.length == 1) {
    lastMouseUp = new Date();
    let diff = (lastMouseUp - lastMouseDown);
    if (diff <= diffDownUp) {
      document.getElementById("mouseC").innerHTML = "click detected";
      sendMouse("click", 0, 0);
    }
  } else if (e.touches == 2) {
    lastMouseRightUp = new Date();
    let diff = (lastMouseRightUp - lastMouseRightDown);
    if (diff <= diffDownUp) {
      sendMouse("rightClick", 0, 0);
    }

  }

}, false);


/*
function isTouchScreendevice() {
  return 'ontouchstart' in window || navigator.maxTouchPoints;
};

if (isTouchScreendevice()) {
  alert("I am a touch screen device")
}
*/



/*
document.getElementById("swtch").addEventListener('change', e => {

  if(e.target.checked){
    document.getElementById("demo").innerHTML = "Hello World";
  }

});
*/


