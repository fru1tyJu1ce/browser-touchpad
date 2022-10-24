package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

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

// Parsing address(connection data) to ./static/client/script.js needed for the frontend websocket connection
func parseAdressToJs(addr string) {
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
