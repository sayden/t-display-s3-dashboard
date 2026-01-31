# Dashboard Server

A simple Go server that provides REST API endpoints for the Lilygo T-Display-S3 Dashboard project.

## Quick Start

```bash
# Run the server
go run main.go

# Or build and run
go build -o dashboard-server
./dashboard-server
```

The server will start on `http://0.0.0.0:8080`.

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/message` | GET | Returns a text message |
| `/api/weather` | GET | Returns weather data |
| `/api/tamagotchi` | GET | Returns tamagotchi game state |
| `/health` | GET | Returns server health status |

## Response Examples

### GET /api/message
```json
{"message": "Hello from Go Server! 🚀"}
```

### GET /api/weather
```json
{"temp": 22.5, "condition": "Sunny", "humidity": 60}
```

### GET /api/tamagotchi
```json
{"name": "Pixel", "hunger": 75, "happy": 80, "energy": 90}
```

## Configuration

Update the T-Display-S3 `config.h` with your server's IP address:

```cpp
#define SERVER_BASE_URL "http://YOUR_IP:8080"
```

To find your IP address:
```bash
# Linux/Mac
ip addr show | grep inet
# or
hostname -I
```

## Testing

```bash
curl http://localhost:8080/api/message
curl http://localhost:8080/api/weather
curl http://localhost:8080/api/tamagotchi
curl http://localhost:8080/health
```
