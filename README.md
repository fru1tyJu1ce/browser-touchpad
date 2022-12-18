
# BrowserTouchpad

Use your touch device as a touchpad for your computer, 
without an extra app on your mobile device.

**Important:** *Both devices have to be in the same network.*


## Usage 

1. Download as zip or clone this repository, unzip and/or change to the repository folder.

2. Windows: Run the BrowserTouchpad.exe and **allow network access** if you get a windows defender firewall warning.
  
    Linux: Run the program with Go. *(Go 1.18 or higher required)*

```bash
  go run .
```

3. A browser tab will open automatically, if not click on the tray icon -> *Show QR Code*.

4. Scan the QR Code with your your mobile device and open the url.


## Screenshots

<img src="https://github.com/onnbt/BrowserTouchpad/blob/master/screenshots/tray.png" width=15% height=15%>
<img src="https://github.com/onnbt/BrowserTouchpad/blob/master/screenshots/browser.png" width=15% height=15%>
<img src="https://github.com/onnbt/BrowserTouchpad/blob/master/screenshots/mobile.jpg" width=15% height=15%>

## Tech Stack

**Client:** Js, HTML, CSS

**Server:** [Go](https://github.com/golang), [Gorilla WebSocket](https://github.com/gorilla/websocket), [Robotgo](https://github.com/go-vgo/robotgo)


## TODO

- Installer for Windows
- Install script for Linux
- Add a license
