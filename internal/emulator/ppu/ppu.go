package ppu

import (
	"fmt"
	"image"

	"golang.org/x/image/draw"

	"github.com/USA-RedDragon/go-gba/internal/config"
	"github.com/USA-RedDragon/go-gba/internal/emulator/memory"
)

const (
	// VRAMSize is 96KB
	VRAMSize = 96 * 1024
	// OAMSize is 1KB
	OAMSize = 1 * 1024
	// PaletteRAMSize is 1KB
	PaletteRAMSize = 1 * 1024
	NUM_PIXELS     = 240 * 160
)

type PPU struct {
	virtualMemory *memory.MMIO
	vRAM          [VRAMSize]byte
	oam           [OAMSize]byte
	paletteRAM    [PaletteRAMSize]byte
	cycle         int
	frameReady    bool
	config        *config.Config
}

func NewPPU(config *config.Config, mmio *memory.MMIO) *PPU {
	ppu := &PPU{
		virtualMemory: mmio,
		cycle:         0,
		frameReady:    false,
		config:        config,
	}

	mmio.AddMMIO(ppu.paletteRAM[:], 0x05000000, PaletteRAMSize)
	mmio.AddMMIO(ppu.vRAM[:], 0x06000000, VRAMSize)
	mmio.AddMMIO(ppu.oam[:], 0x07000000, OAMSize)

	return ppu
}

func (p *PPU) FrameReady() bool {
	return p.frameReady
}

func (p *PPU) ClearFrameReady() {
	p.frameReady = false
}

func (p *PPU) FrameBuffer() []byte {
	dispCNT, err := p.virtualMemory.Read16(0x04000000)
	if err != nil {
		panic(err)
	}

	// Grab bits 0-2 of dispCNT to get the display mode
	displayMode := dispCNT & 0x7

	switch displayMode {
	case 0:
		fmt.Println("Mode 0: Tiled 240x160 8-bpp with 4 backgrounds")
	case 1:
		fmt.Println("Mode 1: Tiled 240x160 8-bpp with 3 backgrounds")
	case 2:
		fmt.Println("Mode 2: Tiled 240x160 8-bpp with 2 backgrounds")
	case 3:
		fmt.Println("Mode 3: Bitmap 240x160 16-bpp with 1 background")
		// Convert vram contents from 16-bit pixels to 8-bit pixels (RGBA)
		// RGBA bitmap
		bitmap := make([]byte, NUM_PIXELS*4)
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
		copy(originalImage.Pix, bitmap)

		// Calculate the target dimensions
		targetWidth := int(float64(240) * p.config.Scale)
		targetHeight := int(float64(160) * p.config.Scale)

		// Create a new blank image with the target dimensions
		upscaledImage := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

		// Perform bicubic interpolation
		draw.CatmullRom.Scale(upscaledImage, upscaledImage.Bounds(), originalImage, originalImage.Bounds(), draw.Src, nil)

		// Convert the upscaled image to a byte array
		upscaled := make([]byte, targetWidth*targetHeight*4)
		copy(upscaled, upscaledImage.Pix)

		return upscaled
	case 4:
		fmt.Println("Mode 4: Bitmap 240x160 8-bpp with 2 backgrounds")
	case 5:
		fmt.Println("Mode 5: Bitmap 160x128 16-bpp with 2 backgrounds")
	default:
		panic(fmt.Sprintf("Invalid display mode: %d", displayMode))
	}

	return nil
}

func (p *PPU) Step() {
	dispCNT, err := p.virtualMemory.Read16(0x04000000)
	if err != nil {
		panic(err)
	}

	// Grab bits 0-2 of dispCNT to get the display mode
	displayMode := dispCNT & 0x7

	switch displayMode {
	case 0:
		fmt.Println("Mode 0: Tiled 240x160 8-bpp with 4 backgrounds")
	case 1:
		fmt.Println("Mode 1: Tiled 240x160 8-bpp with 3 backgrounds")
	case 2:
		fmt.Println("Mode 2: Tiled 240x160 8-bpp with 2 backgrounds")
	case 3:
		fmt.Println("Mode 3: Bitmap 240x160 16-bpp with 1 background")
	case 4:
		fmt.Println("Mode 4: Bitmap 240x160 8-bpp with 2 backgrounds")
	case 5:
		fmt.Println("Mode 5: Bitmap 160x128 16-bpp with 2 backgrounds")
	default:
		panic(fmt.Sprintf("Invalid display mode: %d", displayMode))
	}

	fmt.Printf("Cycle: %d\n", p.cycle)

	// Every 4 cycles is a pixel
	if p.cycle%4 == 0 {
		// Grab the current pixel
		fmt.Println("Pixel")
	} else

	// Every 308 cycles is a scanline
	if p.cycle > 0 && p.cycle%308 == 0 {
		// Grab the current scanline
		fmt.Println("Scanline")
	}

	// Every 280896 cycles is a frame
	if p.cycle > 0 && p.cycle%280896 == 0 {
		// Frame is done
		fmt.Println("Frame")
		p.frameReady = true
		p.cycle = 0
	} else {
		p.cycle++
	}
}
