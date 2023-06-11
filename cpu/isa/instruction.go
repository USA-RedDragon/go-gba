package isa

import "github.com/USA-RedDragon/go-gba/interfaces"

type Instruction interface {
	Execute(cpu interfaces.CPU)
}
