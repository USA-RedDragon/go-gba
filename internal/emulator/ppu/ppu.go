package ppu

import (
	"fmt"
	"image"

	"github.com/USA-RedDragon/go-gba/internal/config"
	"github.com/USA-RedDragon/go-gba/internal/emulator/memory"
	"golang.org/x/image/draw"
)

const (
	// VRAMSize is 96KB
	VRAMSize = 96 * 1024
	// OAMSize is 1KB
	OAMSize = 1 * 1024
	// PaletteRAMSize is 1KB
	PaletteRAMSize = 1 * 1024
	NumPixels      = 240 * 160
)

type PPU struct {
	virtualMemory *memory.MMIO
	vRAM          [VRAMSize]byte
	oam           [OAMSize]byte
	paletteRAM    [PaletteRAMSize]byte
	ioRAM         []byte
	cycle         int
	pixelIndex    int
	scanlineIndex uint8
	frameReady    bool
	config        *config.Config
	HBlank        bool
	VBlank        bool
}

func NewPPU(config *config.Config, mmio *memory.MMIO, ioRAM []byte) *PPU {
	ppu := &PPU{
		virtualMemory: mmio,
		cycle:         0,
		frameReady:    false,
		config:        config,
		ioRAM:         ioRAM,
	}

	mmio.AddMMIO(ppu.paletteRAM[:], 0x05000000, PaletteRAMSize)
	mmio.AddMMIO(ppu.vRAM[:], 0x06000000, VRAMSize)
	mmio.AddMMIO(ppu.oam[:], 0x07000000, OAMSize)

	return ppu
}

// DumpVRAM returns a string representation of the VRAM
func (p *PPU) DumpVRAM() string {
	// The output should look like this:
	// 0x06000000: 0000 0000 0000 0000 0000 0000 0000 0000 ................

	var output string
	for i := 0; i < VRAMSize; i += 16 {
		output += fmt.Sprintf("0x%08X: ", 0x06000000+i)
		for j := 0; j < 16; j++ {
			output += fmt.Sprintf("%02X ", p.vRAM[i+j])
		}
		output += "\n"
	}

	return output
}

func (p *PPU) FrameReady() bool {
	return p.frameReady
}

func (p *PPU) ClearFrameReady() {
	p.frameReady = false
}

func (p *PPU) FrameBuffer() []byte {
	// Grab the first 16 bites of ioRAM
	dispCNT := uint16(p.ioRAM[0]) | uint16(p.ioRAM[1])<<8

	// Grab bits 0-2 of dispCNT to get the display mode
	displayMode := dispCNT & 0x7

	var originalRender *image.RGBA

	switch displayMode {
	case 0:
		fmt.Println("Mode 0: Tiled 240x160 8-bpp with 4 backgrounds")
	case 1:
		fmt.Println("Mode 1: Tiled 240x160 8-bpp with 3 backgrounds")
	case 2:
		fmt.Println("Mode 2: Tiled 240x160 8-bpp with 2 backgrounds")
	case 3:
		fmt.Println("Mode 3: Bitmap 240x160 16-bpp with 1 background")
		originalRender = p.renderMode3()
	case 4:
		fmt.Println("Mode 4: Bitmap 240x160 8-bpp with 2 backgrounds")
		originalRender = p.renderMode4()
	case 5:
		fmt.Println("Mode 5: Bitmap 160x128 16-bpp with 2 backgrounds")
	default:
		panic(fmt.Sprintf("Invalid display mode: %d", displayMode))
	}

	if originalRender == nil {
		return nil
	}

	upscaled := p.upscale(originalRender)

	return upscaled
}

func (p *PPU) upscale(render *image.RGBA) []byte {
	// Calculate the target dimensions
	// We use constants here because we want the scaled
	// image to always be a factor of 240x160
	targetWidth := int(240 * p.config.Scale)
	targetHeight := int(160 * p.config.Scale)

	// Create a new blank image with the target dimensions
	upscaledImage := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// Perform bicubic interpolation
	draw.CatmullRom.Scale(upscaledImage, upscaledImage.Bounds(), render, render.Bounds(), draw.Src, nil)

	// Convert the upscaled image to a byte array for the renderer
	upscaled := make([]byte, targetWidth*targetHeight*4)
	copy(upscaled, upscaledImage.Pix)

	return upscaled
}

func (p *PPU) Step() {
	dispCNT := uint16(p.ioRAM[0]) | uint16(p.ioRAM[1])<<8

	// Grab bits 0-2 of dispCNT to get the display mode
	displayMode := dispCNT & 0x7

	if displayMode > 5 {
		panic(fmt.Sprintf("Invalid display mode: %d", displayMode))
	}

	if p.config.Debug {
		fmt.Printf("PPU Cycle: %d\n", p.cycle)
	}
	// Every 4 cycles is a pixel
	if p.cycle%4 == 0 {
		// Grab the current pixel
		if p.config.Debug {
			fmt.Println("Pixel")
		}
		p.pixelIndex++
	}

	newlyHBlank := false
	newlyNotHBlank := false
	newlyVBlank := false
	newlyNotVBlank := false

	// Every 240+68 pixelIndexes is a scanline
	if p.pixelIndex > 240+68 {
		// Scanline is done
		if p.config.Debug {
			fmt.Println("Scanline")
		}
		p.scanlineIndex++
		p.ioRAM[0x06] = p.scanlineIndex
		p.pixelIndex = 0
		newlyNotHBlank = true
		p.HBlank = false
	} else if p.pixelIndex == 240 {
		// HBlank
		newlyHBlank = true
		p.HBlank = true
	}

	// Every 160+68 scanlines is a frame
	if p.scanlineIndex > 160+68 {
		// Frame is done
		if p.config.Debug {
			fmt.Println("Frame")
		}
		p.frameReady = true
		p.VBlank = false
		p.HBlank = false
		newlyNotVBlank = true
		newlyNotHBlank = true
		p.scanlineIndex = 0
		p.cycle = 0
	} else if p.scanlineIndex == 160 {
		// VBlank
		p.VBlank = true
		newlyVBlank = true
	}

	p.cycle++

	if newlyHBlank {
		p.ioRAM[0x04] |= 0x2
	}

	if newlyNotHBlank {
		p.ioRAM[0x04] &= 0xFD
	}

	if newlyVBlank {
		p.ioRAM[0x04] |= 0x1
	}

	if newlyNotVBlank {
		p.ioRAM[0x04] &= 0xFE
	}
}
