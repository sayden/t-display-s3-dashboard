package main

import (
	"encoding/base64"
)

// Sprite dimensions (80x80 pixels, RGB565 format = 2 bytes per pixel)
const (
	SpriteWidth  = 80
	SpriteHeight = 80
	SpriteSize   = SpriteWidth * SpriteHeight * 2 // 12800 bytes
)

// generateSprite creates an 80x80 RGB565 sprite for the given state
// These are simple procedural sprites - in production, replace with actual artwork
func generateSprite(state string) []byte {
	sprite := make([]byte, SpriteSize)

	// Define colors based on state (RGB565 format, big-endian)
	var bodyColor, eyeColor, mouthColor uint16

	switch state {
	case "happy":
		bodyColor = 0xFE60  // Orange
		eyeColor = 0x0000   // Black
		mouthColor = 0xF800 // Red (happy smile)
	case "sick":
		bodyColor = 0x8410  // Gray-green
		eyeColor = 0x0000   // Black
		mouthColor = 0x8410 // Gray (neutral/sick)
	case "dirty":
		bodyColor = 0x8A22  // Brown
		eyeColor = 0x0000   // Black
		mouthColor = 0x7BCF // Light gray
	case "hungry":
		bodyColor = 0xFE00  // Yellow-orange
		eyeColor = 0x0000   // Black
		mouthColor = 0x7800 // Dark red (open mouth)
	case "sad":
		bodyColor = 0xB5B6  // Light gray
		eyeColor = 0x001F   // Blue (tears)
		mouthColor = 0x0000 // Black (frown)
	default: // normal
		bodyColor = 0xFD20  // Golden
		eyeColor = 0x0000   // Black
		mouthColor = 0x0000 // Black
	}

	backgroundColor := uint16(0x0000) // Black background

	// Draw simple dog shape
	for y := 0; y < SpriteHeight; y++ {
		for x := 0; x < SpriteWidth; x++ {
			idx := (y*SpriteWidth + x) * 2

			color := backgroundColor

			// Body (ellipse from y=30 to y=75, x=15 to x=65)
			bodyY := float64(y-52) / 23.0
			bodyX := float64(x-40) / 25.0
			if bodyX*bodyX+bodyY*bodyY < 1.0 && y >= 30 {
				color = bodyColor
			}

			// Head (circle centered at 40,25 with radius 20)
			headX := float64(x - 40)
			headY := float64(y - 25)
			if headX*headX+headY*headY < 20*20 {
				color = bodyColor
			}

			// Left ear (triangle-ish at top-left of head)
			if x >= 18 && x <= 30 && y >= 5 && y <= 20 {
				earY := float64(y-5) / 15.0
				earLeft := 18.0 + earY*6
				earRight := 30.0 - earY*3
				if float64(x) >= earLeft && float64(x) <= earRight {
					color = bodyColor
				}
			}

			// Right ear (triangle-ish at top-right of head)
			if x >= 50 && x <= 62 && y >= 5 && y <= 20 {
				earY := float64(y-5) / 15.0
				earLeft := 50.0 + earY*3
				earRight := 62.0 - earY*6
				if float64(x) >= earLeft && float64(x) <= earRight {
					color = bodyColor
				}
			}

			// Left eye (circle at 32, 22)
			eyeX := float64(x - 32)
			eyeY := float64(y - 22)
			if eyeX*eyeX+eyeY*eyeY < 4*4 {
				color = 0xFFFF // White
			}
			if eyeX*eyeX+eyeY*eyeY < 2*2 {
				color = eyeColor // Pupil
			}

			// Right eye (circle at 48, 22)
			eyeX = float64(x - 48)
			eyeY = float64(y - 22)
			if eyeX*eyeX+eyeY*eyeY < 4*4 {
				color = 0xFFFF // White
			}
			if eyeX*eyeX+eyeY*eyeY < 2*2 {
				color = eyeColor // Pupil
			}

			// Nose (small triangle at 40, 30)
			if x >= 37 && x <= 43 && y >= 28 && y <= 32 {
				if y >= 28+(x-37)/2 && y >= 28+(43-x)/2 {
					color = 0x0000 // Black nose
				}
			}

			// Mouth (varies by state)
			if state == "happy" {
				// Smile arc
				smileX := float64(x - 40)
				smileY := float64(y - 35)
				if smileX*smileX/64+smileY*smileY/16 < 1.0 && y > 35 && y < 40 {
					color = mouthColor
				}
			} else if state == "sad" {
				// Frown arc
				frownX := float64(x - 40)
				frownY := float64(y - 40)
				if frownX*frownX/64+frownY*frownY/16 < 1.0 && y < 40 && y > 35 {
					color = mouthColor
				}
			} else if state == "hungry" {
				// Open mouth (circle)
				mouthX := float64(x - 40)
				mouthY := float64(y - 37)
				if mouthX*mouthX+mouthY*mouthY < 5*5 {
					color = mouthColor
				}
			}

			// Tail (right side, curved)
			if x >= 62 && x <= 78 && y >= 40 && y <= 55 {
				tailX := float64(x - 62)
				tailCenter := 47.0 + tailX*0.3
				if float64(y) >= tailCenter-3 && float64(y) <= tailCenter+3 {
					color = bodyColor
				}
			}

			// Legs (4 small rectangles)
			// Front left leg
			if x >= 22 && x <= 28 && y >= 65 && y <= 78 {
				color = bodyColor
			}
			// Front right leg
			if x >= 35 && x <= 41 && y >= 65 && y <= 78 {
				color = bodyColor
			}
			// Back left leg
			if x >= 42 && x <= 48 && y >= 65 && y <= 78 {
				color = bodyColor
			}
			// Back right leg
			if x >= 55 && x <= 61 && y >= 65 && y <= 78 {
				color = bodyColor
			}

			// Add special effects based on state
			if state == "sick" {
				// Sweat drops
				if (x == 55 && y >= 15 && y <= 18) || (x == 58 && y >= 18 && y <= 21) {
					color = 0x07FF // Cyan
				}
			}

			if state == "sad" {
				// Tear drops
				if (x == 30 && y >= 26 && y <= 30) || (x == 50 && y >= 26 && y <= 30) {
					color = 0x001F // Blue
				}
			}

			if state == "dirty" {
				// Dirt spots
				if (x == 25 && y == 35) || (x == 55 && y == 40) || (x == 45 && y == 55) {
					color = 0x4208 // Dark brown
				}
			}

			// Write RGB565 color (little-endian for ESP32)
			sprite[idx] = byte(color & 0xFF)
			sprite[idx+1] = byte(color >> 8)
		}
	}

	return sprite
}

// GetSprite returns the sprite data for the given state as base64
func GetSprite(state string) (string, int, int) {
	sprite := generateSprite(state)
	encoded := base64.StdEncoding.EncodeToString(sprite)
	return encoded, SpriteWidth, SpriteHeight
}

// GetSpriteRaw returns raw sprite bytes (for direct binary transfer)
func GetSpriteRaw(state string) ([]byte, int, int) {
	sprite := generateSprite(state)
	return sprite, SpriteWidth, SpriteHeight
}
