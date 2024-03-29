package ppu

import (
	"image"
)

func (p *PPU) renderMode3() *image.RGBA {
	// Convert vram contents from 16-bit pixels to 8-bit pixels (RGBA)
	// RGBA bitmap
	var bitmap [NumPixels * 4]byte
	for i := 0; i < NumPixels*2; i += 2 {
		pixel := uint16(p.vRAM[i+1])<<8 | uint16(p.vRAM[i])
		// Convert XBGR1555 to 32-bit RGBA
		destIndex := i * 2
		bitmap[destIndex] = byte((pixel & 0x1F) << 3)
		bitmap[destIndex+1] = byte(((pixel >> 5) & 0x1F) << 3)
		bitmap[destIndex+2] = byte(((pixel >> 10) & 0x1F) << 3)
		bitmap[destIndex+3] = 0xFF

		bitmap[destIndex] = byte((pixel&0x1F)<<3 | ((pixel & 0x1F) >> 2))                 // Red
		bitmap[destIndex+1] = byte(((pixel>>5)&0x1F)<<3 | (((pixel >> 5) & 0x1F) >> 2))   // Green
		bitmap[destIndex+2] = byte(((pixel>>10)&0x1F)<<3 | (((pixel >> 10) & 0x1F) >> 2)) // Blue
		bitmap[destIndex+3] = 0xFF                                                        // Alpha
	}

	originalImage := image.NewRGBA(image.Rect(0, 0, 240, 160))
	copy(originalImage.Pix, bitmap[:])

	return originalImage
}

func (p *PPU) renderMode4() *image.RGBA {
	dispCNT := uint16(p.ioRAM[0]) | uint16(p.ioRAM[1])<<8

	// bit 4 of dispCNT determines which page of vram to use
	page := dispCNT & 0x10

	startAddr := uint32(0x0000)
	if page == 1 {
		startAddr = 0xA000
	}

	bitmap := make([]byte, NumPixels*4)
	// Each pixel is 8 bits
	for i := 0; i < NumPixels; i++ {
		paletteRAMOffset := p.vRAM[startAddr+uint32(i)]
		paletteRAMOffset *= 2
		destIndex := i * 4
		// The value at the paletteRamAddr is a 16-bit color
		// Convert XBGR1555 to 32-bit RGBA
		pixel := uint16(p.paletteRAM[int(paletteRAMOffset)+1])<<8 | uint16(p.paletteRAM[int(paletteRAMOffset)])

		bitmap[destIndex] = byte((pixel&0x1F)<<3 | ((pixel & 0x1F) >> 2))                 // Red
		bitmap[destIndex+1] = byte(((pixel>>5)&0x1F)<<3 | (((pixel >> 5) & 0x1F) >> 2))   // Green
		bitmap[destIndex+2] = byte(((pixel>>10)&0x1F)<<3 | (((pixel >> 10) & 0x1F) >> 2)) // Blue
		bitmap[destIndex+3] = 0xFF                                                        // Alpha
	}

	originalImage := image.NewRGBA(image.Rect(0, 0, 240, 160))
	copy(originalImage.Pix, bitmap)

	return originalImage
}
