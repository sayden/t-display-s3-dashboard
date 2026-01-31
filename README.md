# Dashboards - Multi-Screen Scrolling Dashboard

A multi-screen dashboard application for the **Lilygo T-Display-S3** that displays HTTP-fetched data, rotating images, and a game placeholder.

## Features

### 4 Interactive Dashboards

1. **Message Display (D1)** - Fetches and displays a message from your server via HTTP
2. **Weather Display (D2)** - Shows temperature, weather condition, and humidity from your server
3. **Image Gallery (D3)** - Rotates through 4 placeholder images (TODO: replace with your own images)
4. **Tamagotchi Game (D4)** - Placeholder for future game implementation

### Navigation

- **Auto-rotation**: Dashboards automatically cycle every 10 seconds (configurable)
- **Touch navigation**: Tap the right side of the screen to advance to the next dashboard
- **Visual indicators**: Dots at the bottom show which dashboard is active

### Connectivity

- WiFi connection with status indicator
- HTTP client for fetching data from local server
- JSON parsing for message and weather data
- Error handling for connection failures

## Configuration

Before uploading, edit **`config.h`** to set:

### WiFi Credentials
```cpp
#define WIFI_SSID     "Your_WiFi_SSID"
#define WIFI_PASSWORD "Your_WiFi_Password"
```

### Server Settings
```cpp
#define SERVER_BASE_URL "http://192.168.1.100:8080"
```

### API Endpoints

Your server should provide these endpoints:

**Message Endpoint** (`/api/message`):
```json
{
  "message": "Hello World"
}
```

**Weather Endpoint** (`/api/weather`):
```json
{
  "temp": 22.5,
  "condition": "Sunny",
  "humidity": 60
}
```

### Timing Configuration

```cpp
#define AUTO_ROTATION_MS 10000     // Dashboard rotation interval (ms)
#define IMAGE_ROTATION_MS 2500     // Image rotation interval for D3 (ms)
```

## Dependencies

This project requires the following Arduino libraries:

- **TFT_eSPI** - Display driver (already configured for T-Display-S3)
- **TouchDrvCSTXXX** - Touch sensor driver (from LilyGo examples)
- **WiFi** - Built-in ESP32 library
- **HTTPClient** - Built-in ESP32 library
- **ArduinoJson** - Install via Library Manager

## Installation

1. Open Arduino IDE
2. Install required libraries (if not already installed)
3. Open `Dashboards.ino`
4. Edit `config.h` with your WiFi credentials and server URL
5. Select board: **ESP32S3 Dev Module** or **LilyGo T-Display-S3**
6. Upload to your device

## Server Implementation

The dashboard expects a simple REST API server on your local network. Here's an example structure:

### Go Example (TODO - to be implemented by user)
```go
package main

import (
    "encoding/json"
    "net/http"
)

type Message struct {
    Message string `json:"message"`
}

type Weather struct {
    Temp      float64 `json:"temp"`
    Condition string  `json:"condition"`
    Humidity  int     `json:"humidity"`
}

func main() {
    http.HandleFunc("/api/message", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(Message{Message: "Hello from Go!"})
    })
    
    http.HandleFunc("/api/weather", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(Weather{
            Temp:      22.5,
            Condition: "Sunny",
            Humidity:  60,
        })
    })
    
    http.ListenAndServe(":8080", nil)
}
```

### Rust Example (TODO - to be implemented by user)
```rust
// Similar REST API implementation in Rust
// Using frameworks like actix-web or axum
```

## Customization

### Replace Placeholder Images (Dashboard 3)

The current implementation uses runtime-generated gradients. To use real images:

1. Resize your images to **120x100 pixels**
2. Convert to C array using [LVGL Image Converter](https://lvgl.io/tools/imageconverter)
3. Select **RGB565** color format and **C array** output
4. Replace the `generateImageN()` functions in `images.h` with the generated arrays

### Touch Sensitivity

Edit `config.h` to adjust touch detection:
```cpp
#define TOUCH_NEXT_THRESHOLD_X 160    // X coordinate for "next" gesture
#define TOUCH_MIN_DURATION_MS 50      // Minimum touch duration
```

## Troubleshooting

### WiFi Connection Failed
- Check SSID and password in `config.h`
- Ensure your device is in range of the WiFi network
- Check Serial Monitor (115200 baud) for error messages

### HTTP Errors
- Verify server is running and accessible on the same network
- Test endpoints with `curl` or browser: `http://192.168.1.100:8080/api/message`
- Check Serial Monitor for HTTP response codes

### Touch Not Working
- Touch sensor may need recalibration
- Try different touch IC (CST328 vs CST816) - see CapacitiveTouch example
- Check Serial Monitor for touch initialization messages

### Display Issues
- Ensure TFT_eSPI is configured with Setup206_LilyGo_T_Display_S3.h
- Check pin definitions in `pin_config.h`

## Future Enhancements (TODO)

- [ ] Implement Tamagotchi game on Dashboard 4
- [ ] Add more dashboard types (clock, system stats, etc.)
- [ ] SD card support for image storage
- [ ] OTA update support
- [ ] Settings menu for runtime configuration
- [ ] Battery level indicator
- [ ] Sleep mode with wake on touch

## License

Based on LilyGo T-Display-S3 examples. See original repository for license details.

## Credits

- Display and touch examples from [Xinyuan-LilyGO/T-Display-S3](https://github.com/Xinyuan-LilyGO/T-Display-S3)
- Built with Arduino framework for ESP32
