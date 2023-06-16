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

	// Bits 7-6 are the hi operand flags
	hof := m.instruction & (1<<7 | 1<<6) >> 6
	// Bits 5-3 are the source register
	rs := uint8(m.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)
	// Bits 2-0 are the destination register
	rd := uint8(m.instruction & (1<<2 | 1<<1 | 1<<0))

	switch hof {
	case 0b01:
		// move hi register source to low register destination
		fmt.Printf("mov r%d, r%d\n", rd, rs+8)
		cpu.WriteRegister(rd, cpu.ReadHighRegister(rs))
	case 0b10:
		// move low register source to hi register destination
		fmt.Printf("mov r%d, r%d\n", rd+8, rs)
		cpu.WriteHighRegister(rd, cpu.ReadRegister(rs))
	case 0b11:
		// move hi register source to hi register destination
		fmt.Printf("mov r%d, r%d\n", rd+8, rs+8)
		cpu.WriteHighRegister(rd, cpu.ReadHighRegister(rs))
	default:
		panic("Invalid hi operand flag")
	}
}
