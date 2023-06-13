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
	return &Emulator{
		config: config,
		cpu:    cpu.NewARM7TDMI(config),
	}
}

func (e *Emulator) Update() error {
	return nil
}

func (e *Emulator) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
	// Get the framebuffer from the CPU
	// Then screen.WritePixels(fb)
}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(e.config.Scale * 240), int(e.config.Scale * 160)
}
