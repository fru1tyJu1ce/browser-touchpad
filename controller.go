package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
)

type ClientMessage struct {
	Type string `json:"Type"`
	DX   int    `json:"dX"`
	DY   int    `json:"dY"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var fileServer = http.FileServer(http.Dir("./static"))

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
			actualPosX, actualPosY := robotgo.GetMousePos()
			robotgo.MoveMouse(actualPosX+msg.DX, actualPosY+msg.DY)
		case "click":
			robotgo.MouseClick("left", true)
		case "doubleClick":
			robotgo.MouseClick("left", true)
			robotgo.MouseClick("left", true)
		default:
			fmt.Println("unknown message type")
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
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
	fmt.Println("server is running")
	wg.Wait()
}
