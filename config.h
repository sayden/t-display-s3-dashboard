#pragma once

// ========================================
// WiFi Configuration
// ========================================
// Include sensitive data from non-committed secrets.h
// If you don't have this file, copy secrets.h.example to secrets.h
#include "secrets.h"

// WiFi connection timeout (milliseconds)
#define WIFI_CONNECT_TIMEOUT_MS 10000

// ========================================
// Server Configuration
// ========================================
// Replace with your server's IP address and port
#define SERVER_BASE_URL "http://192.168.1.201:8081"

// API endpoints for each dashboard
#define MESSAGE_ENDPOINT                                                       \
  "/api/message" // Expected JSON: {"message": "Hello World"}
#define WEATHER_ENDPOINT                                                       \
  "/api/weather" // Expected JSON: {"temp": 22.5, "condition": "Sunny",
                 // "humidity": 60}
// Dashboard 4 (Tamagotchi) endpoint - currently unused but available for future
// integration
#define TAMAGOTCHI_ENDPOINT "/api/tamagotchi"
#define INTERVALS_ENDPOINT "/api/intervals" // Intervals.icu training data

// HTTP timeout (milliseconds)
#define HTTP_TIMEOUT_MS 5000

// ========================================
// Dashboard Configuration
// ========================================
// Auto-rotation interval in milliseconds (1 minute by default)
#define AUTO_ROTATION_MS 60000

// Number of dashboards
#define NUM_DASHBOARDS 4

// Dashboard 3 (Image Animation) settings
#define IMAGE_ROTATION_MS 83  // Time each image is displayed
#define NUM_ROTATION_IMAGES 4 // Number of images to rotate through

// ========================================
// Time Configuration
// ========================================
#define NTP_SERVER "pool.ntp.org"
#define GMT_OFFSET_SEC 3600      // GMT+1 (UTC+1)
#define DAYLIGHT_OFFSET_SEC 3600 // Daylight saving (+1h)

// ========================================
// Touch Configuration
// ========================================
// Touch detection area for "next dashboard" gesture
// X coordinate threshold (right half of screen = x > 60)
#define TOUCH_NEXT_THRESHOLD_X 160

// Minimum touch duration to register (milliseconds)
#define TOUCH_MIN_DURATION_MS 50

// ========================================
// Display Configuration
// ========================================
// Screen dimensions (T-Display-S3)
#define SCREEN_WIDTH 320
#define SCREEN_HEIGHT 170

// Color scheme
#define COLOR_BACKGROUND 0x0000 // Black
#define COLOR_PRIMARY TFT_CYAN
#define COLOR_SECONDARY TFT_YELLOW
#define COLOR_TEXT TFT_WHITE
#define COLOR_ERROR TFT_RED
#define COLOR_INDICATOR_ON TFT_ORANGE
#define COLOR_INDICATOR_OFF 0x4208 // Dark gray
