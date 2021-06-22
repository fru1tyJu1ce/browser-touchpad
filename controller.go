package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var fileServer = http.FileServer(http.Dir("./static"))

func delChar(s []rune, index int) []rune {
	return append(s[0:index], s[index+1:]...)
}

func reader(conn *websocket.Conn) {

	//	lastX, lastY := robotgo.GetMousePos()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))

		//if string(p)[0] != '!' {
		//	continue
		//}

		c := strings.Split(string(p), " ")

		x, err := strconv.Atoi(c[0])
		if err != nil {
			// handle error
			fmt.Println(err)
		}

		y, err := strconv.Atoi(c[1])
		if err != nil {
			// handle error
			fmt.Println(err)
		}

		actX, actY := robotgo.GetMousePos()

		robotgo.MoveMouse(actX+x, actY+y)

		if err := conn.WriteMessage(messageType, p); err != nil {
			return
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

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
