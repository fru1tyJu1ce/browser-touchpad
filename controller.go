package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"strconv"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
)

var Run bool = true // App iss runing while run == true
var Addr string     // Connection data

var fileServerServer = http.FileServer(http.Dir("./static/server"))
var fileServerClient = http.FileServer(http.Dir("./static/client"))

type ClientMessage struct {
	Type string `json:"Type"`
	DX   int    `json:"dx"`
	DY   int    `json:"dy"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  128,
	WriteBufferSize: 128,
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func mouseMove(dx int, dy int) {
	actualPosX, actualPosY := robotgo.GetMousePos()
	robotgo.MoveMouse(actualPosX+dx, actualPosY+dy)
}

func scroll(dy int) {
	if dy < 0 {
		robotgo.ScrollMouse(dy, "down")
	} else if dy > 0 {
		robotgo.ScrollMouse(dy, "up")
	}
}

func reader(conn *websocket.Conn) {

	//	lastX, lastY := robotgo.GetMousePos()
	var msg ClientMessage
	for {
		messageType, rawMsg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(rawMsg)
		errjSON := json.Unmarshal([]byte(rawMsg), &msg)
		if errjSON != nil {
			log.Println(err)
			continue
		}

		log.Println(msg)

		switch msgType := msg.Type; msgType {
		case "movement":
			mouseMove(msg.DX, msg.DY)
		case "click":
			robotgo.Click("left")
		case "rightClick":
			robotgo.Click("right")
		case "scroll":
			scroll(msg.DY)
		case "toggle":
			robotgo.MouseToggle("down", "left")
		case "toggleUp":
			robotgo.MouseToggle("up")
		default:
			fmt.Println("unknown command recived")
		}

		if err := conn.WriteMessage(messageType, rawMsg); err != nil {
			return
		}

	}

}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } //TODO cyber sec muy importante

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("client connected")
	reader(ws)
}

func setupRoutes() {
	http.Handle("/client/", http.StripPrefix("/client", fileServerClient))
	http.Handle("/server/", http.StripPrefix("/server", fileServerServer))
	http.HandleFunc("/ws", wsEndpoint)
}

func setup() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	Addr = GetOutboundIP().String() + ":" + strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	parseAdressToJs(Addr)
	createQr(Addr + "/client")
	//open("http://" + Addr + "/client")
	go open("http://" + Addr + "/server")
	fmt.Println("server started", Addr)
	panic(http.Serve(listener, nil))

}

func runCheck() {
	for {
		if !Run {
			return
		}
	}
}

func main() {
	InitTrayHandler()
	setupRoutes()
	go setup()
	runCheck()
}
