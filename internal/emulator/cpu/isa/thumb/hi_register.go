package thumb

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type ADDH struct {
	instruction uint16
}

func (a ADDH) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("ADDH")

	panic("Not implemented")
	return
}

type CMPH struct {
	instruction uint16
}

func (c CMPH) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("CMPH")

	// Bits 7-6 are the hi operand flags
	hof := c.instruction & (1<<7 | 1<<6) >> 6
	// Bits 5-3 are the source register
	rs := uint8(c.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)
	// Bits 2-0 are the destination register
	rd := uint8(c.instruction & (1<<2 | 1<<1 | 1<<0))

	switch hof {
	case 0b01:
		// cmp Rd, Hs
		fmt.Printf("cmp r%d, r%d\n", rd, rs+8)
		res := cpu.ReadRegister(rd) - cpu.ReadHighRegister(rs)
		cpu.SetZ(res == 0)
		cpu.SetN(res&0x80000000 != 0)
		cpu.SetC(cpu.ReadRegister(rd) >= cpu.ReadHighRegister(rs))
	case 0b10:
		// cmp Hd, Rs
		fmt.Printf("cmp r%d, r%d\n", rd+8, rs)
		res := cpu.ReadHighRegister(rd) - cpu.ReadRegister(rs)
		fmt.Printf("res: %08X\n", res)
		cpu.SetZ(res == 0)
		cpu.SetN(res&0x80000000 != 0)
		cpu.SetC(cpu.ReadHighRegister(rd) >= cpu.ReadRegister(rs))
	case 0b11:
		// cmp Hd, Hs
		fmt.Printf("cmp r%d, r%d\n", rd+8, rs+8)
		res := cpu.ReadHighRegister(rd) - cpu.ReadHighRegister(rs)
		fmt.Printf("res: %08X\n", res)
		cpu.SetZ(res == 0)
		cpu.SetN(res&0x80000000 != 0)
		cpu.SetC(cpu.ReadHighRegister(rd) >= cpu.ReadHighRegister(rs))
	default:
		panic("Invalid hi operand flag")
	}

	return
}

type MOVH struct {
	instruction uint16
}

func (m MOVH) Execute(cpu interfaces.CPU) (repipeline bool) {
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

	return
}
