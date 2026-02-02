/**
 * @file      Dashboards.ino
 * @brief     Multi-screen scrolling dashboard for Lilygo T-Display-S3
 * @date      2026-01-30
 *
 * Features:
 * - Dashboard 1: Message display (HTTP fetched)
 * - Dashboard 2: Weather display (HTTP fetched)
 * - Dashboard 3: Rotating image animation
 * - Dashboard 4: Tamagotchi game placeholder
 *
 * Navigation: Auto-rotates every X seconds or tap right side of screen
 */

#include "config.h"
#include "images.h"
#include "pin_config.h"
#include <Arduino.h>
#include <ArduinoJson.h>
#include <HTTPClient.h>
#include <TFT_eSPI.h>
#include <TouchDrvCSTXXX.hpp>
#include <WiFi.h>

// ========================================
// Global Objects
// ========================================
TFT_eSPI tft = TFT_eSPI();
TFT_eSprite sprite = TFT_eSprite(&tft);
TouchDrvCSTXXX touch;

// ========================================
// State Variables
// ========================================
int currentDashboard = 0;
unsigned long lastRotationTime = 0;
unsigned long lastImageRotationTime = 0;
int currentImageIndex = 0;
bool wifiConnected = false;

// Cached data
String cachedMessage = "";
String cachedWeatherTemp = "";
String cachedWeatherCondition = "";
String cachedWeatherHumidity = "";
unsigned long lastMessageFetch = 0;
unsigned long lastWeatherFetch = 0;

// Tamagotchi cached data
String cachedDogName = "";
int cachedHunger = 0;
int cachedHappiness = 0;
int cachedHygiene = 0;
int cachedHealth = 0;
String cachedState = "";
bool cachedNeedsAttention = false;
uint16_t dogSpriteBuffer[80 * 80]; // 80x80 sprite buffer
bool spriteLoaded = false;
unsigned long lastTamagotchiFetch = 0;
int tamagotchiTouchAction = 0; // 0=none, 1=feed, 2=play, 3=clean

// ========================================
// HTTP Helper Functions
// ========================================

/**
 * Performs HTTP GET request
 * @param url Full URL to request
 * @return Response body as String, empty string on failure
 */
String httpGet(String url) {
  if (!wifiConnected) {
    Serial.println("WiFi not connected");
    return "";
  }

  HTTPClient http;
  http.setTimeout(HTTP_TIMEOUT_MS);

  Serial.print("HTTP GET: ");
  Serial.println(url);

  if (!http.begin(url)) {
    Serial.println("HTTP begin failed");
    return "";
  }

  int httpCode = http.GET();
  String payload = "";

  if (httpCode > 0) {
    Serial.printf("HTTP response code: %d\n", httpCode);
    if (httpCode == HTTP_CODE_OK) {
      payload = http.getString();
      Serial.println("Response: " + payload);
    }
  } else {
    Serial.printf("HTTP GET failed: %s\n",
                  http.errorToString(httpCode).c_str());
  }

  http.end();
  return payload;
}

/**
 * Fetch message from server for Dashboard 1
 */
void fetchMessage() {
  String url = String(SERVER_BASE_URL) + String(MESSAGE_ENDPOINT);
  String response = httpGet(url);

  if (response.length() > 0) {
    // Parse JSON: {"message": "Hello World"}
    StaticJsonDocument<256> doc;
    DeserializationError error = deserializeJson(doc, response);

    if (!error) {
      cachedMessage = doc["message"].as<String>();
      lastMessageFetch = millis();
      Serial.println("Message updated: " + cachedMessage);
    } else {
      Serial.println("JSON parse error: " + String(error.c_str()));
      cachedMessage = "Error: Invalid JSON";
    }
  } else {
    cachedMessage = "Error: Connection failed";
  }
}

String cachedWindSpeed = "";
String cachedForecast3h = "";
String cachedForecastTom = "";

/**
 * Fetch weather from server for Dashboard 2
 */
void fetchWeather() {
  String url = String(SERVER_BASE_URL) + String(WEATHER_ENDPOINT);
  String response = httpGet(url);

  if (response.length() > 0) {
    // Parse JSON with extended forecast data
    StaticJsonDocument<1024> doc; // Increased buffer for larger JSON
    DeserializationError error = deserializeJson(doc, response);

    if (!error) {
      cachedWeatherTemp = String(doc["temp"].as<float>(), 1);
      cachedWeatherCondition = doc["condition"].as<String>();
      cachedWeatherHumidity = String(doc["humidity"].as<int>());
      cachedWindSpeed = String(doc["wind_speed"].as<float>(), 1);

      // Format 3h Forecast
      float t3 = doc["forecast_3h"]["temp"];
      String c3 = doc["forecast_3h"]["condition"];
      float p3 = doc["forecast_3h"]["precip"];
      cachedForecast3h =
          "+3h: " + String(t3, 0) + "C " + c3 + " Rain:" + String(p3, 0) + "%";

      // Format Tomorrow Forecast
      float tT = doc["forecast_tom"]["temp"];
      String cT = doc["forecast_tom"]["condition"];
      float pT = doc["forecast_tom"]["precip"];
      cachedForecastTom =
          "Tom: " + String(tT, 0) + "C " + cT + " Rain:" + String(pT, 0) + "%";

      lastWeatherFetch = millis();
      Serial.println("Weather updated for Aix-les-Bains");
    } else {
      Serial.println("JSON parse error: " + String(error.c_str()));
      cachedWeatherCondition = "Error: Invalid JSON";
    }
  } else {
    cachedWeatherCondition = "Error: Connection failed";
  }
}

// ========================================
// Base64 Decoding for Sprite Images
// ========================================

// Base64 character lookup table
static const char base64_chars[] =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

int base64_decode_char(char c) {
  if (c >= 'A' && c <= 'Z')
    return c - 'A';
  if (c >= 'a' && c <= 'z')
    return c - 'a' + 26;
  if (c >= '0' && c <= '9')
    return c - '0' + 52;
  if (c == '+')
    return 62;
  if (c == '/')
    return 63;
  return -1;
}

/**
 * Decode base64 string directly into sprite buffer
 * @param base64 Base64 encoded string
 * @param output Pointer to uint16_t buffer
 * @param maxBytes Maximum bytes to decode
 * @return Number of bytes decoded
 */
int decodeBase64ToBuffer(const String &base64, uint16_t *output, int maxBytes) {
  int len = base64.length();
  int outIdx = 0;
  uint8_t *outBytes = (uint8_t *)output;

  for (int i = 0; i < len && outIdx < maxBytes; i += 4) {
    int a = base64_decode_char(base64[i]);
    int b = (i + 1 < len) ? base64_decode_char(base64[i + 1]) : 0;
    int c = (i + 2 < len) ? base64_decode_char(base64[i + 2]) : 0;
    int d = (i + 3 < len) ? base64_decode_char(base64[i + 3]) : 0;

    if (a < 0)
      a = 0;
    if (b < 0)
      b = 0;
    if (c < 0)
      c = 0;
    if (d < 0)
      d = 0;

    uint32_t triple = (a << 18) | (b << 12) | (c << 6) | d;

    if (outIdx < maxBytes)
      outBytes[outIdx++] = (triple >> 16) & 0xFF;
    if (outIdx < maxBytes && base64[i + 2] != '=')
      outBytes[outIdx++] = (triple >> 8) & 0xFF;
    if (outIdx < maxBytes && base64[i + 3] != '=')
      outBytes[outIdx++] = triple & 0xFF;
  }

  return outIdx;
}

// ========================================
// Tamagotchi API Functions
// ========================================

/**
 * Fetch tamagotchi state from server
 */
void fetchTamagotchi() {
  String url = String(SERVER_BASE_URL) + String(TAMAGOTCHI_ENDPOINT);
  String response = httpGet(url);

  if (response.length() > 0) {
    // Use larger buffer for image data - parse in chunks
    DynamicJsonDocument doc(32768); // 32KB for JSON with base64 image
    DeserializationError error = deserializeJson(doc, response);

    if (!error) {
      cachedDogName = doc["dog"]["name"].as<String>();
      cachedHunger = doc["dog"]["hunger"].as<int>();
      cachedHappiness = doc["dog"]["happiness"].as<int>();
      cachedHygiene = doc["dog"]["hygiene"].as<int>();
      cachedHealth = doc["dog"]["health"].as<int>();
      cachedState = doc["state"].as<String>();
      cachedNeedsAttention = doc["needs_attention"].as<bool>();

      // Decode sprite image
      String imageBase64 = doc["image"].as<String>();
      if (imageBase64.length() > 0) {
        int decoded =
            decodeBase64ToBuffer(imageBase64, dogSpriteBuffer, 80 * 80 * 2);
        spriteLoaded = (decoded > 0);
        Serial.printf("Sprite decoded: %d bytes\n", decoded);
      }

      lastTamagotchiFetch = millis();
      Serial.printf("Tamagotchi updated: %s (State: %s)\n",
                    cachedDogName.c_str(), cachedState.c_str());
    } else {
      Serial.println("JSON parse error: " + String(error.c_str()));
      cachedState = "error";
    }
  } else {
    cachedState = "offline";
  }
}

/**
 * Send action to tamagotchi (feed, play, clean)
 */
void sendTamagotchiAction(const char *action) {
  String url = String(SERVER_BASE_URL) + "/api/tamagotchi/" + String(action);

  if (!wifiConnected) {
    Serial.println("WiFi not connected");
    return;
  }

  HTTPClient http;
  http.setTimeout(HTTP_TIMEOUT_MS);

  Serial.print("POST: ");
  Serial.println(url);

  if (!http.begin(url)) {
    Serial.println("HTTP begin failed");
    return;
  }

  int httpCode = http.POST("");

  if (httpCode > 0) {
    Serial.printf("Action response code: %d\n", httpCode);
    // Refresh state after action
    fetchTamagotchi();
  } else {
    Serial.printf("Action failed: %s\n", http.errorToString(httpCode).c_str());
  }

  http.end();
}

// ========================================
// Drawing Helper Functions
// ========================================

/**
 * Draw dashboard indicator dots
 */
void drawIndicator() {
  int dotSize = 6;
  int spacing = 12;
  int totalWidth = NUM_DASHBOARDS * spacing - (spacing - dotSize);
  int startX = (SCREEN_WIDTH - totalWidth) / 2;
  int y = SCREEN_HEIGHT - 12;

  for (int i = 0; i < NUM_DASHBOARDS; i++) {
    int x = startX + (i * spacing);
    uint16_t color =
        (i == currentDashboard) ? COLOR_INDICATOR_ON : COLOR_INDICATOR_OFF;
    sprite.fillCircle(x, y, dotSize / 2, color);
  }
}

/**
 * Draw WiFi status icon
 */
void drawWiFiIcon() {
  int x = 5;
  int y = 5;
  uint16_t color = wifiConnected ? COLOR_PRIMARY : COLOR_ERROR;

  // Simple WiFi icon (three arcs)
  sprite.drawArc(x + 5, y + 8, 10, 8, 180, 360, color, COLOR_BACKGROUND);
  sprite.drawArc(x + 5, y + 8, 7, 5, 180, 360, color, COLOR_BACKGROUND);
  sprite.fillCircle(x + 5, y + 8, 2, color);
}

/**
 * Draw header with title and indicators
 */
void drawHeader(String title) {
  sprite.setTextColor(COLOR_PRIMARY, COLOR_BACKGROUND);
  sprite.setTextDatum(TC_DATUM);
  sprite.drawString(title, SCREEN_WIDTH / 2, 5, 2);

  drawWiFiIcon();
  drawIndicator();
}

/**
 * Draw centered text with word wrapping
 */
void drawCenteredText(String text, int y, int fontsize) {
  sprite.setTextColor(COLOR_TEXT, COLOR_BACKGROUND);
  sprite.setTextDatum(MC_DATUM);

  // Simple word wrapping for long text
  int maxWidth = SCREEN_WIDTH - 20;
  if (sprite.textWidth(text, fontsize) > maxWidth) {
    // Split into words and wrap
    int lineHeight = 20;
    int currentY = y - lineHeight;
    String line = "";
    String word = "";

    for (int i = 0; i < text.length(); i++) {
      char c = text.charAt(i);
      if (c == ' ' || i == text.length() - 1) {
        if (i == text.length() - 1)
          word += c;

        String testLine = line + (line.length() > 0 ? " " : "") + word;
        if (sprite.textWidth(testLine, fontsize) > maxWidth &&
            line.length() > 0) {
          sprite.drawString(line, SCREEN_WIDTH / 2, currentY, fontsize);
          currentY += lineHeight;
          line = word;
        } else {
          line = testLine;
        }
        word = "";
      } else {
        word += c;
      }
    }
    if (line.length() > 0) {
      sprite.drawString(line, SCREEN_WIDTH / 2, currentY, fontsize);
    }
  } else {
    sprite.drawString(text, SCREEN_WIDTH / 2, y, fontsize);
  }
}

// ========================================
// Dashboard Rendering Functions
// ========================================

/**
 * Dashboard 1: Message Display
 */
void drawDashboard1() {
  sprite.fillSprite(COLOR_BACKGROUND);
  drawHeader("X FEED");

  // Fetch message if not cached or old
  if (cachedMessage.length() == 0 || millis() - lastMessageFetch > 60000) {
    sprite.setTextColor(COLOR_SECONDARY, COLOR_BACKGROUND);
    sprite.setTextDatum(MC_DATUM);
    sprite.drawString("Loading...", SCREEN_WIDTH / 2, SCREEN_HEIGHT / 2, 2);
    sprite.pushSprite(0, 0);

    fetchMessage();
  }

  // Draw message
  sprite.fillSprite(COLOR_BACKGROUND);
  drawHeader("MESSAGE");

  if (cachedMessage.startsWith("Error:")) {
    sprite.setTextColor(COLOR_ERROR, COLOR_BACKGROUND);
  } else {
    sprite.setTextColor(COLOR_TEXT, COLOR_BACKGROUND);
  }

  drawCenteredText(cachedMessage, SCREEN_HEIGHT / 2, 4);

  sprite.pushSprite(0, 0);
}

/**
 * Dashboard 2: Weather Display
 */
void drawDashboard2() {
  sprite.fillSprite(COLOR_BACKGROUND);
  drawHeader("AIX-LES-BAINS"); // Updated title for specific location

  // Fetch weather if not cached or old
  if (cachedWeatherCondition.length() == 0 ||
      millis() - lastWeatherFetch > 600000) { // 10 min cache
    sprite.setTextColor(COLOR_SECONDARY, COLOR_BACKGROUND);
    sprite.setTextDatum(MC_DATUM);
    sprite.drawString("Loading Weather...", SCREEN_WIDTH / 2, SCREEN_HEIGHT / 2,
                      2);
    sprite.pushSprite(0, 0);

    fetchWeather();
  }

  // Draw weather
  sprite.fillSprite(COLOR_BACKGROUND);
  drawHeader("AIX-LES-BAINS");

  if (cachedWeatherCondition.startsWith("Error:")) {
    sprite.setTextColor(COLOR_ERROR, COLOR_BACKGROUND);
    sprite.setTextDatum(MC_DATUM);
    sprite.drawString(cachedWeatherCondition, SCREEN_WIDTH / 2,
                      SCREEN_HEIGHT / 2, 2);
  } else {
    // Current Weather (Top Section)
    sprite.setTextDatum(TL_DATUM);

    // Temperature (Large)
    sprite.setTextColor(COLOR_SECONDARY, COLOR_BACKGROUND);
    sprite.drawString(cachedWeatherTemp + " C", 20, 35, 6);

    // Condition & Wind (Right side)
    sprite.setTextColor(COLOR_PRIMARY, COLOR_BACKGROUND);
    sprite.drawString(cachedWeatherCondition, 140, 40, 4);

    sprite.setTextColor(COLOR_TEXT, COLOR_BACKGROUND);
    sprite.drawString("Wind: " + cachedWindSpeed + " km/h", 140, 65, 2);
    sprite.drawString("Hum: " + cachedWeatherHumidity + "%", 140, 80, 2);

    // Forecast Separator
    sprite.drawLine(10, 105, SCREEN_WIDTH - 10, 105, COLOR_SECONDARY);

    // Forecasts (Bottom Section)
    sprite.setTextColor(COLOR_TEXT, COLOR_BACKGROUND);
    sprite.setTextDatum(TC_DATUM);

    // +3h Forecast
    sprite.drawString(cachedForecast3h, SCREEN_WIDTH / 2, 115, 2);

    // Tomorrow Forecast
    sprite.drawString(cachedForecastTom, SCREEN_WIDTH / 2, 135, 2);
  }

  sprite.pushSprite(0, 0);
}

/**
 * Dashboard 3: Image Rotation
 */
void drawDashboard3() {
  sprite.fillSprite(COLOR_BACKGROUND);
  drawHeader("GALLERY");

  // Rotate through images
  if (millis() - lastImageRotationTime > IMAGE_ROTATION_MS) {
    currentImageIndex = (currentImageIndex + 1) % NUM_ROTATION_IMAGES;
    lastImageRotationTime = millis();
  }

  // Draw current image
  int imgWidth = 120;
  int imgHeight = 100;
  int x = (SCREEN_WIDTH - imgWidth) / 2;
  int y = (SCREEN_HEIGHT - imgHeight) / 2 + 5;

  sprite.pushImage(x, y, imgWidth, imgHeight,
                   placeholderImages[currentImageIndex]);

  // Image counter
  sprite.setTextColor(COLOR_TEXT, COLOR_BACKGROUND);
  sprite.setTextDatum(BC_DATUM);
  String counter =
      String(currentImageIndex + 1) + " / " + String(NUM_ROTATION_IMAGES);
  sprite.drawString(counter, SCREEN_WIDTH / 2, SCREEN_HEIGHT - 20, 2);

  sprite.pushSprite(0, 0);
}

/**
 * Dashboard 4: Tamagotchi Placeholder
 */
void drawDashboard4() {
  sprite.fillSprite(COLOR_BACKGROUND);

  // Header with attention indicator
  if (cachedNeedsAttention) {
    sprite.setTextColor(COLOR_ERROR, COLOR_BACKGROUND);
    drawHeader("! BUDDY NEEDS YOU !");
  } else {
    drawHeader("TAMAGOTCHI");
  }

  // Fetch data if needed (every 30 seconds)
  if (cachedDogName.length() == 0 || millis() - lastTamagotchiFetch > 30000) {
    sprite.setTextColor(COLOR_SECONDARY, COLOR_BACKGROUND);
    sprite.setTextDatum(MC_DATUM);
    sprite.drawString("Loading...", SCREEN_WIDTH / 2, SCREEN_HEIGHT / 2, 2);
    sprite.pushSprite(0, 0);

    fetchTamagotchi();
    sprite.fillSprite(COLOR_BACKGROUND);
    if (cachedNeedsAttention) {
      sprite.setTextColor(COLOR_ERROR, COLOR_BACKGROUND);
      drawHeader("! BUDDY NEEDS YOU !");
    } else {
      drawHeader("TAMAGOTCHI");
    }
  }

  // Draw dog sprite (left side)
  if (spriteLoaded) {
    int spriteX = 10;
    int spriteY = 25;
    sprite.pushImage(spriteX, spriteY, 80, 80, dogSpriteBuffer);
  } else {
    // Fallback simple dog face
    int cx = 50;
    int cy = 65;
    sprite.fillCircle(cx, cy, 25, COLOR_PRIMARY);
    sprite.fillCircle(cx - 8, cy - 5, 4, COLOR_BACKGROUND);
    sprite.fillCircle(cx + 8, cy - 5, 4, COLOR_BACKGROUND);
    sprite.drawArc(cx, cy, 15, 12, 200, 340, COLOR_BACKGROUND, COLOR_PRIMARY);
  }

  // Draw dog name and state
  sprite.setTextColor(COLOR_SECONDARY, COLOR_BACKGROUND);
  sprite.setTextDatum(TC_DATUM);
  sprite.drawString(cachedDogName, 50, 110, 2);

  sprite.setTextColor(COLOR_TEXT, COLOR_BACKGROUND);
  sprite.drawString(cachedState, 50, 125, 1);

  // Draw stat bars (right side)
  int barX = 105;
  int barY = 28;
  int barWidth = 80;
  int barHeight = 10;
  int barSpacing = 22;

  // Helper lambda-like function for drawing bars
  auto drawStatBar = [&](int y, const char *label, int value, uint16_t color) {
    // Label
    sprite.setTextColor(COLOR_TEXT, COLOR_BACKGROUND);
    sprite.setTextDatum(TL_DATUM);
    sprite.drawString(label, barX, y - 12, 1);

    // Bar background
    sprite.fillRect(barX, y, barWidth, barHeight, 0x2104);

    // Bar fill
    int fillWidth = (value * barWidth) / 100;
    sprite.fillRect(barX, y, fillWidth, barHeight, color);

    // Value text
    sprite.setTextDatum(TR_DATUM);
    sprite.drawString(String(value), barX + barWidth + 20, y - 1, 1);
  };

  // Draw each stat bar
  drawStatBar(barY, "Hunger", cachedHunger, TFT_GREEN);
  drawStatBar(barY + barSpacing, "Happy", cachedHappiness, TFT_YELLOW);
  drawStatBar(barY + barSpacing * 2, "Clean", cachedHygiene, TFT_CYAN);
  drawStatBar(barY + barSpacing * 3, "Health", cachedHealth,
              cachedHealth > 50 ? TFT_GREEN : TFT_RED);

  // Draw action buttons at bottom
  int btnY = 138;
  int btnWidth = 100; // Wider buttons
  int btnHeight = 22;
  int btnSpacing = 6;
  int btnStartX = 5;

  auto drawButton = [&](int x, const char *label, uint16_t color) {
    sprite.fillRoundRect(x, btnY, btnWidth, btnHeight, 4, color);
    sprite.setTextColor(COLOR_BACKGROUND, color);
    sprite.setTextDatum(MC_DATUM);
    sprite.drawString(label, x + btnWidth / 2, btnY + btnHeight / 2, 2);
  };

  // Three action buttons distributed across 320px
  // 5 + 100 + 6 + 100 + 6 + 100 = 317px
  drawButton(btnStartX, "FEED", TFT_GREEN);
  drawButton(btnStartX + btnWidth + btnSpacing, "PLAY", TFT_ORANGE);
  drawButton(btnStartX + 2 * (btnWidth + btnSpacing), "CLEAN", TFT_CYAN);

  sprite.pushSprite(0, 0);
}

/**
 * Handle touch specifically for Tamagotchi dashboard
 * Returns true if touch was handled
 */
bool handleTamagotchiTouch(int16_t touchX, int16_t touchY) {
  if (currentDashboard != 3)
    return false;

  // Touch driver with rotation 1 has X and Y SWAPPED and X INVERTED:
  // - Raw touchY corresponds to horizontal position (screen X)
  // - Raw touchX corresponds to vertical position (screen Y) inverted?
  // User reports: "reaction happening on top of screen instead of bottom"
  // This means high touchX (Top) passed the check check, while low touchX
  // (Bottom) didn't. So Top = High X, Bottom = Low X. To map to Screen Y
  // (0=Top, 170=Bottom):
  int16_t screenX = touchY;       // Horizontal (0-320)
  int16_t screenY = 170 - touchX; // Vertical (inverted)

  Serial.printf("Tamagotchi touch: raw(%d,%d) -> screen(%d,%d)\n", touchX,
                touchY, screenX, screenY);

  // Button layout: 3 buttons spread across 320px width
  // ~106px per button
  int btn1Boundary = 106;
  int btn2Boundary = 212;

  // Check if touch is in button vertical zone (Bottom of screen)
  // Screen Y is 0-170. Buttons are at bottom > 110
  if (screenY >= 110) {
    Serial.printf("In button zone. screenX=%d\n", screenX);

    // FEED button: Left (0-106)
    if (screenX < btn1Boundary) {
      Serial.println("Touch: FEED");
      sendTamagotchiAction("feed?type=meal");
      return true;
    }
    // PLAY button: Middle (106-212)
    else if (screenX >= btn1Boundary && screenX < btn2Boundary) {
      Serial.println("Touch: PLAY");
      sendTamagotchiAction("play");
      return true;
    }
    // CLEAN button: Right (212-320)
    else {
      Serial.println("Touch: CLEAN");
      sendTamagotchiAction("clean?type=bath");
      return true;
    }

    return true;
  }

  return false;
}

/**
 * Render current dashboard
 */
void renderDashboard() {
  switch (currentDashboard) {
  case 0:
    drawDashboard1();
    break;
  case 1:
    drawDashboard2();
    break;
  case 2:
    drawDashboard3();
    break;
  case 3:
    drawDashboard4();
    break;
  }
}

// ========================================
// Navigation Functions
// ========================================

/**
 * Advance to next dashboard
 */
void nextDashboard() {
  currentDashboard = (currentDashboard + 1) % NUM_DASHBOARDS;
  lastRotationTime = millis();

  // Reset image rotation when entering Dashboard 3
  if (currentDashboard == 2) {
    lastImageRotationTime = millis();
  }

  Serial.println("Dashboard: " + String(currentDashboard));
  renderDashboard();
}

/**
 * Handle touch input
 */
void handleTouch() {
  int16_t x[1], y[1];
  uint8_t touched = touch.getPoint(x, y, 1);

  if (touched) {
    Serial.printf("Touch detected: x=%d, y=%d\n", x[0], y[0]);

    // Check for Tamagotchi button presses first
    if (handleTamagotchiTouch(x[0], y[0])) {
      delay(TOUCH_MIN_DURATION_MS * 4); // Longer debounce for actions
      renderDashboard();                // Refresh display
      return;
    }

    // Otherwise, any tap = next dashboard
    nextDashboard();
    delay(TOUCH_MIN_DURATION_MS * 2); // Debounce
  }
}

// ========================================
// WiFi Functions
// ========================================

/**
 * Initialize WiFi connection
 */
void initWiFi() {
  Serial.println("Connecting to WiFi...");
  sprite.fillSprite(COLOR_BACKGROUND);
  sprite.setTextColor(COLOR_TEXT, COLOR_BACKGROUND);
  sprite.setTextDatum(MC_DATUM);
  sprite.drawString("Connecting to WiFi...", SCREEN_WIDTH / 2,
                    SCREEN_HEIGHT / 2, 2);
  sprite.pushSprite(0, 0);

  WiFi.mode(WIFI_STA);
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

  unsigned long startTime = millis();
  while (WiFi.status() != WL_CONNECTED &&
         millis() - startTime < WIFI_CONNECT_TIMEOUT_MS) {
    delay(500);
    Serial.print(".");
  }

  if (WiFi.status() == WL_CONNECTED) {
    wifiConnected = true;
    Serial.println("\nWiFi connected!");
    Serial.print("IP address: ");
    Serial.println(WiFi.localIP());

    sprite.fillSprite(COLOR_BACKGROUND);
    sprite.setTextColor(COLOR_PRIMARY, COLOR_BACKGROUND);
    sprite.drawString("WiFi Connected!", SCREEN_WIDTH / 2,
                      SCREEN_HEIGHT / 2 - 10, 2);
    sprite.setTextColor(COLOR_TEXT, COLOR_BACKGROUND);
    sprite.drawString(WiFi.localIP().toString(), SCREEN_WIDTH / 2,
                      SCREEN_HEIGHT / 2 + 10, 2);
    sprite.pushSprite(0, 0);
    delay(2000);
  } else {
    wifiConnected = false;
    Serial.println("\nWiFi connection failed!");

    sprite.fillSprite(COLOR_BACKGROUND);
    sprite.setTextColor(COLOR_ERROR, COLOR_BACKGROUND);
    sprite.drawString("WiFi Failed", SCREEN_WIDTH / 2, SCREEN_HEIGHT / 2, 2);
    sprite.setTextColor(COLOR_TEXT, COLOR_BACKGROUND);
    sprite.drawString("Check config.h", SCREEN_WIDTH / 2,
                      SCREEN_HEIGHT / 2 + 20, 2);
    sprite.pushSprite(0, 0);
    delay(3000);
  }
}

// ========================================
// Setup & Loop
// ========================================

void setup() {
  Serial.begin(115200);
  Serial.println("Dashboards starting...");

  // Power on display
  pinMode(PIN_POWER_ON, OUTPUT);
  digitalWrite(PIN_POWER_ON, HIGH);

  // Initialize display
  tft.init();
  tft.setRotation(1);
  tft.fillScreen(TFT_BLACK);

  sprite.createSprite(SCREEN_WIDTH, SCREEN_HEIGHT);
  sprite.setSwapBytes(true);
  sprite.setTextDatum(MC_DATUM);

  // Set backlight
  ledcSetup(0, 10000, 8);
  ledcAttachPin(PIN_LCD_BL, 0);
  ledcWrite(0, 160);

  // Initialize touch
  touch.setPins(PIN_TOUCH_RES, PIN_TOUCH_INT);
  if (!touch.begin(Wire, CST328_SLAVE_ADDRESS, PIN_IIC_SDA, PIN_IIC_SCL)) {
    Serial.println("Touch CST328 failed, trying CST816...");
    if (!touch.begin(Wire, CST816_SLAVE_ADDRESS, PIN_IIC_SDA, PIN_IIC_SCL)) {
      Serial.println("Touch initialization failed!");
    }
  }
  touch.disableAutoSleep();

  // Initialize placeholder images
  initPlaceholderImages();

  // Initialize WiFi
  initWiFi();

  // Initialize timers
  lastRotationTime = millis();
  lastImageRotationTime = millis();

  // Render first dashboard
  renderDashboard();

  Serial.println("Setup complete!");
}

void loop() {
  // Auto-rotation
  if (millis() - lastRotationTime > AUTO_ROTATION_MS) {
    nextDashboard();
  }

  // Handle touch input
  handleTouch();

  // Re-render Dashboard 3 for animation
  if (currentDashboard == 2) {
    renderDashboard();
  }

  delay(50);
}

// ========================================
// TFT Pin Verification
// ========================================
#if PIN_LCD_WR != TFT_WR || PIN_LCD_RD != TFT_RD || PIN_LCD_CS != TFT_CS ||    \
    PIN_LCD_DC != TFT_DC || PIN_LCD_RES != TFT_RST || PIN_LCD_D0 != TFT_D0 ||  \
    PIN_LCD_D1 != TFT_D1 || PIN_LCD_D2 != TFT_D2 || PIN_LCD_D3 != TFT_D3 ||    \
    PIN_LCD_D4 != TFT_D4 || PIN_LCD_D5 != TFT_D5 || PIN_LCD_D6 != TFT_D6 ||    \
    PIN_LCD_D7 != TFT_D7 || PIN_LCD_BL != TFT_BL ||                            \
    TFT_BACKLIGHT_ON != HIGH || 170 != TFT_WIDTH || 320 != TFT_HEIGHT
#error                                                                         \
    "Error! Please make sure <User_Setups/Setup206_LilyGo_T_Display_S3.h> is selected in <TFT_eSPI/User_Setup_Select.h>"
#endif

#if ESP_IDF_VERSION >= ESP_IDF_VERSION_VAL(5, 0, 0)
#error                                                                         \
    "The current version is not supported for the time being, please use a version below Arduino ESP32 3.0"
#endif
