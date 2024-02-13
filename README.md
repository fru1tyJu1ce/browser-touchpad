# Browser Touchpad

Transform your touch device into a touchpad for your computer, all without needing an extra app on your mobile device.

- *Works seamlessly on Linux, Windows, and MacOS.*
- *Both devices must be connected to the same network.*

## Usage

1. **Download and Installation**

   Download or clone this repository and navigate to the repository folder.

2. **Build and Run**

   - **Windows:** Run the `BrowserTouchpad.exe` file and allow network access if prompted by Windows Defender Firewall.
   - **Alternative for Windows, Linux, and MacOS:** Build yourself and run the program with Go 1.18 or higher.

3. **Connect Your Devices**

   A browser tab will open automatically. If not, click on the tray icon and select "Show QR Code".

4. **Pair Your Devices**

   Scan the displayed QR code using your mobile device's camera.

## Screenshots

<img src="https://github.com/onnbt/BrowserTouchpad/blob/master/screenshots/tray.png" width=30% height=30% alt="Tray icon">
<img src="https://github.com/onnbt/BrowserTouchpad/blob/master/screenshots/browser.png" width=30% height=30% alt="Browser interface">
<img src="https://github.com/onnbt/BrowserTouchpad/blob/master/screenshots/mobile.jpg" width=30% height=30% alt="Mobile device">

## Tech Stack

**Client:** JavaScript, HTML, CSS

**Server:** [Go](https://github.com/golang), [Gorilla WebSocket](https://github.com/gorilla/websocket), [Robotgo](https://github.com/go-vgo/robotgo)
