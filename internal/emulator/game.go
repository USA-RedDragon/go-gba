package emulator

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/USA-RedDragon/go-gba/internal/config"
	"github.com/USA-RedDragon/go-gba/internal/emulator/cpu"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Emulator struct {
	config    *config.Config
	cpu       *cpu.ARM7TDMI
	stopped   bool
	frametime int
}

func New(config *config.Config) *Emulator {
	emu := &Emulator{
		config: config,
		cpu:    cpu.NewARM7TDMI(config),
	}
	return emu
}

func (e *Emulator) Update() error {
	start := time.Now()
	for {
		if e.stopped {
			break
		}
		e.cpu.Step()
		if e.cpu.PPU.FrameReady() {
			e.cpu.PPU.ClearFrameReady()
			e.frametime = int(time.Since(start).Milliseconds())
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
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f\nFrame Time: %dms\nTPS: %0.2f", 1000.0/float64(e.frametime), e.frametime, ebiten.ActualTPS()))
	// ebitenutil.DebugPrint(screen, e.cpu.DebugRegisters())
}

func (e *Emulator) Layout(_, _ int) (int, int) {
	return int(e.config.Scale * 240), int(e.config.Scale * 160)
}

func (e *Emulator) Stop() {
	e.stopped = true
	pprof.StopCPUProfile()
	e.cpu.Halt()
	os.Exit(0)
}
