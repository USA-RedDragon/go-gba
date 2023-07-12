package emulator

import (
	"github.com/USA-RedDragon/go-gba/internal/config"
	"github.com/USA-RedDragon/go-gba/internal/emulator/cpu"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Emulator struct {
	config        *config.Config
	cpu           *cpu.ARM7TDMI
	prevFB        []byte
	prevFBPresent bool
}

func New(config *config.Config) *Emulator {
	emu := &Emulator{
		config: config,
		cpu:    cpu.NewARM7TDMI(config),
	}
	go emu.cpu.Run()
	return emu
}

func (e *Emulator) Update() error {
	if e.cpu.PPU.FrameReady() {
		fb := e.cpu.PPU.FrameBuffer()
		e.cpu.PPU.ClearFrameReady()
		// Print vram to stderr
		// e.cpu.PPU.DumpVRAM()
		// fmt.Fprint(os.Stderr, e.cpu.PPU.DumpVRAM())
		if fb != nil {
			if !e.prevFBPresent {
				e.prevFB = fb
				e.prevFBPresent = true
			}
		}
	}
	return nil
}

func (e *Emulator) Draw(screen *ebiten.Image) {
	if e.prevFBPresent {
		screen.WritePixels(e.prevFB)
	}
	ebitenutil.DebugPrint(screen, e.cpu.DebugRegisters())
}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(e.config.Scale * 240), int(e.config.Scale * 160)
}
