package thumb

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type ADDH struct {
	instruction uint16
}

func (a ADDH) Execute(cpu interfaces.CPU) {
	fmt.Println("ADDH")

	panic("Not implemented")
}

type CMPH struct {
	instruction uint16
}

func (c CMPH) Execute(cpu interfaces.CPU) {
	fmt.Println("CMPH")

	// This one needs to set condition flags

	panic("Not implemented")
}

type MOVH struct {
	instruction uint16
}

func (m MOVH) Execute(cpu interfaces.CPU) {
	fmt.Println("MOVH")

	panic("Not implemented")
}
