package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
	qrcode "github.com/yeqown/go-qrcode"
)

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func createQr(addr string) {
	qrc, err := qrcode.New("http://" + addr)
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
	}

	if err := qrc.Save("./static/server/qrcode.jpeg"); err != nil {
		fmt.Printf("could not save image: %v", err)
	}
}

func parsesAdressToJs(addr string) {
	socketJS := "const socket = new WebSocket(" + string('"') + "ws://" + addr + "/ws" + string('"') + ");\n"

	file, err := os.OpenFile("./static/client/script.js", os.O_RDWR, 0644)
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
	ReadBufferSize:  128,
	WriteBufferSize: 128,
}

var fileServerServer = http.FileServer(http.Dir("./static/server"))
var fileServerClient = http.FileServer(http.Dir("./static/client"))

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
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } //TODO

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

func main() {
	var wg sync.WaitGroup
	setupRoutes()
	wg.Add(1)
	go func() {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			panic(err)
		}
		addr := GetOutboundIP().String() + ":" + strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
		parsesAdressToJs(addr)
		createQr(addr + "/client")
		open("http://" + addr + "/client")
		open("http://" + addr + "/server")
		fmt.Println("server started", addr)
		panic(http.Serve(listener, nil))
	}()
	wg.Wait()
}
