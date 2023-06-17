package isa

import "github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"

type Instruction interface {
	Execute(cpu interfaces.CPU) (repipeline bool)
}
