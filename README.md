
# BrowserTouchpad

Use your touch device as a touchpad for your computer, 
without an extra app on your mobile device.

**Important:** *Both devices have to be in the same network.*


## Usage 

1. Download or clone this repository and change in to the repository folder.

2. Windows: Run the BrowserTouchpad.exe and **allow network access** if you get a windows defender firewall warning.
  
    Linux: Compile and Run the program with Go. *(Go 1.18 or higher required)*

3. A browser tab will open, if not click on the tray icon -> *Show QR Code*.

4. Scan the QR Code with your your mobile device.


## Screenshots

<img src="https://github.com/onnbt/BrowserTouchpad/blob/master/screenshots/tray.png" width=13% height=13%>
<img src="https://github.com/onnbt/BrowserTouchpad/blob/master/screenshots/browser.png" width=13% height=13%>
<img src="https://github.com/onnbt/BrowserTouchpad/blob/master/screenshots/mobile.jpg" width=13% height=13%>

## Tech

**Client:** Js, HTML, CSS

**Server:** [Go](https://github.com/golang), [Gorilla WebSocket](https://github.com/gorilla/websocket), [Robotgo](https://github.com/go-vgo/robotgo)

