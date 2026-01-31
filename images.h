#pragma once

#include <Arduino.h>

/**
 * Placeholder images for Dashboard 3 (Image Rotation)
 *
 * These are simple 120x100 pixel images encoded as RGB565 arrays.
 * Replace these with your own images converted to C arrays.
 *
 * To convert images:
 * 1. Resize to 120x100 pixels
 * 2. Use online tool: https://lvgl.io/tools/imageconverter
 * 3. Select "C array" output and "RGB565" color format
 */

// Image dimensions
#define IMG_WIDTH 120
#define IMG_HEIGHT 100

// ========================================
// Image 1: Gradient (Blue to Cyan)
// ========================================
const uint16_t image1[IMG_WIDTH * IMG_HEIGHT] PROGMEM = {
    // Generated programmatically-a blue to cyan gradient
    // This is a placeholder - in production, use actual image data
};

// Helper function to generate Image 1 at runtime
void generateImage1(uint16_t *buffer) {
  for (int y = 0; y < IMG_HEIGHT; y++) {
    for (int x = 0; x < IMG_WIDTH; x++) {
      // Blue to Cyan gradient (vertical)
      uint8_t b = 255;
      uint8_t g = (y * 255) / IMG_HEIGHT;
      uint8_t r = 0;

      // Convert RGB888 to RGB565
      uint16_t color = ((r & 0xF8) << 8) | ((g & 0xFC) << 3) | (b >> 3);
      buffer[y * IMG_WIDTH + x] = color;
    }
  }
}

// ========================================
// Image 2: Gradient (Red to Yellow)
// ========================================
void generateImage2(uint16_t *buffer) {
  for (int y = 0; y < IMG_HEIGHT; y++) {
    for (int x = 0; x < IMG_WIDTH; x++) {
      // Red to Yellow gradient (vertical)
      uint8_t r = 255;
      uint8_t g = (y * 255) / IMG_HEIGHT;
      uint8_t b = 0;

      // Convert RGB888 to RGB565
      uint16_t color = ((r & 0xF8) << 8) | ((g & 0xFC) << 3) | (b >> 3);
      buffer[y * IMG_WIDTH + x] = color;
    }
  }
}

// ========================================
// Image 3: Gradient (Green to Cyan)
// ========================================
void generateImage3(uint16_t *buffer) {
  for (int y = 0; y < IMG_HEIGHT; y++) {
    for (int x = 0; x < IMG_WIDTH; x++) {
      // Green to Cyan gradient (vertical)
      uint8_t r = 0;
      uint8_t g = 255;
      uint8_t b = (y * 255) / IMG_HEIGHT;

      // Convert RGB888 to RGB565
      uint16_t color = ((r & 0xF8) << 8) | ((g & 0xFC) << 3) | (b >> 3);
      buffer[y * IMG_WIDTH + x] = color;
    }
  }
}

// ========================================
// Image 4: Checkerboard Pattern
// ========================================
void generateImage4(uint16_t *buffer) {
  for (int y = 0; y < IMG_HEIGHT; y++) {
    for (int x = 0; x < IMG_WIDTH; x++) {
      // Checkerboard pattern (20x20 squares)
      int squareSize = 20;
      bool isWhite = ((x / squareSize) + (y / squareSize)) % 2 == 0;

      uint16_t color = isWhite ? 0xFFFF : 0x8410; // White or Gray
      buffer[y * IMG_WIDTH + x] = color;
    }
  }
}

// ========================================
// Image Buffer Array (Generated at runtime)
// ========================================
// Note: These are generated at runtime to save flash memory
// For production, replace with actual pre-encoded images

static uint16_t imageBuffer1[IMG_WIDTH * IMG_HEIGHT];
static uint16_t imageBuffer2[IMG_WIDTH * IMG_HEIGHT];
static uint16_t imageBuffer3[IMG_WIDTH * IMG_HEIGHT];
static uint16_t imageBuffer4[IMG_WIDTH * IMG_HEIGHT];

static uint16_t *placeholderImages[NUM_ROTATION_IMAGES] = {
    imageBuffer1, imageBuffer2, imageBuffer3, imageBuffer4};

// ========================================
// Initialize placeholder images
// ========================================
void initPlaceholderImages() {
  generateImage1(imageBuffer1);
  generateImage2(imageBuffer2);
  generateImage3(imageBuffer3);
  generateImage4(imageBuffer4);
}
