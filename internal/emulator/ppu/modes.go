package ppu

import (
	"image"
)

func (p *PPU) renderMode3() *image.RGBA {
	// Convert vram contents from 16-bit pixels to 8-bit pixels (RGBA)
	// RGBA bitmap
	var bitmap [NUM_PIXELS * 4]byte
	for i := 0; i < NUM_PIXELS*2; i += 2 {
		pixel := uint16(p.vRAM[i+1])<<8 | uint16(p.vRAM[i])
		// Convert XBGR1555 to 32-bit RGBA
		destIndex := i * 2
		bitmap[destIndex] = byte((pixel & 0x1F) << 3)
		bitmap[destIndex+1] = byte(((pixel >> 5) & 0x1F) << 3)
		bitmap[destIndex+2] = byte(((pixel >> 10) & 0x1F) << 3)
		bitmap[destIndex+3] = 0xFF
	}

	originalImage := image.NewRGBA(image.Rect(0, 0, 240, 160))
	copy(originalImage.Pix, bitmap[:])

	return originalImage
}

func (p *PPU) renderMode4() *image.RGBA {
	dispCNT, err := p.virtualMemory.Read16(0x04000000)
	if err != nil {
		panic(err)
	}

	// bit 4 of dispCNT determines which page of vram to use
	page := dispCNT & 0x10

	startAddr := uint32(0x06000000)
	if page == 1 {
		startAddr = 0x0600A000
	}

	bitmap := make([]byte, NUM_PIXELS*4)
	// Each pixel is 8 bits
	for i := 0; i < NUM_PIXELS; i++ {
		paletteRamOffset, err := p.virtualMemory.Read8(startAddr + uint32(i))
		if err != nil {
			panic(err)
		}
		paletteRamOffset *= 2
		destIndex := i * 4
		// The value at the paletteRamAddr is a 16-bit color
		// Convert XBGR1555 to 32-bit RGBA
		pixel := uint16(p.paletteRAM[int(paletteRamOffset)+1])<<8 | uint16(p.paletteRAM[int(paletteRamOffset)])

		bitmap[destIndex] = byte((pixel & 0x1F) << 3)
		bitmap[destIndex+1] = byte(((pixel >> 5) & 0x1F) << 3)
		bitmap[destIndex+2] = byte(((pixel >> 10) & 0x1F) << 3)
		bitmap[destIndex+3] = 0xFF
	}

	originalImage := image.NewRGBA(image.Rect(0, 0, 240, 160))
	copy(originalImage.Pix, bitmap)

	return originalImage
}
