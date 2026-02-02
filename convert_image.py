
import sys
from PIL import Image
import os

def convert_image(image_path, output_path):
    img = Image.open(image_path)
    img = img.convert('RGB')
    
    # Force resize to 320x170
    target_width = 320
    target_height = 170
    img = img.resize((target_width, target_height), Image.Resampling.LANCZOS)
    
    width, height = img.size
    print(f"Processing image: {width}x{height}")
    
    with open(output_path, 'w') as f:
        f.write('#pragma once\n\n')
        f.write('#include <Arduino.h>\n\n')
        f.write(f'#define GALLERY_IMG_WIDTH {width}\n')
        f.write(f'#define GALLERY_IMG_HEIGHT {height}\n\n')
        f.write('const uint16_t gallery_image[{}] PROGMEM = {{\n'.format(width * height))
        
        data = []
        pixels = list(img.getdata())
        
        for r, g, b in pixels:
            # Convert to RGB565
            # R: 5 bits, G: 6 bits, B: 5 bits
            r5 = (r >> 3) & 0x1F
            g6 = (g >> 2) & 0x3F
            b5 = (b >> 3) & 0x1F
            
            rgb565 = (r5 << 11) | (g6 << 5) | b5
            data.append(f'0x{rgb565:04X}')
        
        # Write format: 16 pixels per line for neatness
        for i in range(0, len(data), 16):
            f.write('  ' + ', '.join(data[i:i+16]) + ',\n')
            
        f.write('};\n')

if __name__ == '__main__':
    if len(sys.argv) < 3:
        print("Usage: python convert_image.py <input_image> <output_header>")
        sys.exit(1)
        
    convert_image(sys.argv[1], sys.argv[2])
