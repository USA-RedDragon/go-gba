package emulator

import (
	"os"

	"github.com/USA-RedDragon/go-gba/internal/config"
	"github.com/USA-RedDragon/go-gba/internal/emulator/cpu"
	"github.com/hajimehoshi/ebiten/v2"
)

type Emulator struct {
	config  *config.Config
	cpu     *cpu.ARM7TDMI
	stopped bool
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
		if e.stopped {
			break
		}
		e.cpu.Step()
		if e.cpu.PPU.FrameReady() {
			e.cpu.PPU.ClearFrameReady()
			break
		}
	}
	return nil
}

func (e *Emulator) Draw(screen *ebiten.Image) {
	if e.stopped {
		return
	}
	fb := e.cpu.PPU.FrameBuffer()
	if fb != nil {
		screen.WritePixels(fb)
	}
	// ebitenutil.DebugPrint(screen, e.cpu.DebugRegisters())
}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(e.config.Scale * 240), int(e.config.Scale * 160)
}

func (e *Emulator) Stop() {
	e.stopped = true
	e.cpu.Halt()
	os.Exit(0)
}
