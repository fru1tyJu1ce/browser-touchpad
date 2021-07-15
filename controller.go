package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
	qrcode "github.com/yeqown/go-qrcode"
)

func createQr(addr string) {
	qrc, err := qrcode.New("http://" + addr)
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
	}

	// save file
	if err := qrc.Save("./static/qrcode.jpeg"); err != nil {
		fmt.Printf("could not save image: %v", err)
	}
}

func parsesAdressToJs(addr string) {
	socketJS := "const socket = new WebSocket(" + string('"') + "ws://" + addr + "/ws" + string('"') + ");  \n"

	file, err := os.OpenFile("./static/script.js", os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	_len, err := file.WriteAt([]byte(socketJS), 0)
	if err != nil && _len != len(socketJS) {
		log.Fatalf("failed writing to file: %s", err)
	}

}

type ClientMessage struct {
	Type string `json:"Type"`
	DX   int    `json:"dx"`
	DY   int    `json:"dy"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var fileServer = http.FileServer(http.Dir("./static"))

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

func reader(conn *websocket.Conn) {

	//	lastX, lastY := robotgo.GetMousePos()
	var msg ClientMessage
	for {
		messageType, rawMsg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(rawMsg))

		errjSON := json.Unmarshal([]byte(rawMsg), &msg)
		if errjSON != nil {
			log.Println(err)
			continue
		}

		switch msgType := msg.Type; msgType {
		case "mouseMove":
			mouseMove(msg.DX, msg.DY)
		case "click":
			robotgo.MouseClick()
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
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } //TODO

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("client connected")
	reader(ws)
}

func setupRoutes() {
	http.Handle("/", fileServer)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	var wg sync.WaitGroup
	fmt.Println("startup")
	setupRoutes()
	wg.Add(1)
	go func() {
		//log.Fatal(http.ListenAndServe(":8080", nil))

		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			panic(err)
		}

		addr := GetOutboundIP().String() + ":" + strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
		parsesAdressToJs(addr)
		createQr(addr)
		fmt.Println("server is running", addr)
		panic(http.Serve(listener, nil))

	}()
	wg.Wait()
}
