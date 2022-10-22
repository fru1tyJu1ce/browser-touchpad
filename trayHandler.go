package main

import (
	"io/ioutil"
	"log"

	"github.com/getlantern/systray"
)

func getIcon() []byte { //TODO Changeable icon color

	bytes, err := ioutil.ReadFile("./static/res/iconwhite.ico")
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func onReady() {
	systray.SetIcon(getIcon())
	systray.SetTitle("BrowserTouchpad")
	systray.SetTooltip("active")
	controll := systray.AddMenuItem("Show QR Code", "Open the webinterface")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		<-controll.ClickedCh
		open("http://" + Addr + "/server")
	}()

	go func() {
		<-mQuit.ClickedCh
		Run = false
		systray.Quit()
	}()
}

func onExit() { // TODO Chek 4 real cleanup
	log.Println("exit trayHandler")
}

func InitTrayHandler() {
	go systray.Run(onReady, onExit)
}
