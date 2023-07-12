package emulator

import (
	"github.com/USA-RedDragon/go-gba/internal/config"
	"github.com/USA-RedDragon/go-gba/internal/emulator/cpu"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Emulator struct {
	config *config.Config
	cpu    *cpu.ARM7TDMI
}

func New(config *config.Config) *Emulator {
	emu := &Emulator{
		config: config,
		cpu:    cpu.NewARM7TDMI(config),
	}
	return emu
}

func (e *Emulator) Update() error {
	for {
		e.cpu.Step()
		if e.cpu.PPU.FrameReady() {
			e.cpu.PPU.ClearFrameReady()
			break
		}
	}
	return nil
}

func (e *Emulator) Draw(screen *ebiten.Image) {
	fb := e.cpu.PPU.FrameBuffer()
	if fb != nil {
		screen.WritePixels(fb)
	}
	ebitenutil.DebugPrint(screen, e.cpu.DebugRegisters())
}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(e.config.Scale * 240), int(e.config.Scale * 160)
}
