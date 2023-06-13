package thumb

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type MOV struct {
	instruction uint16
}

func (m MOV) Execute(cpu interfaces.CPU) {
	// Bits 10-8 are the destination register
	rd := uint8(m.instruction & (1<<10 | 1<<9 | 1<<8) >> 8)
	// Bits 7-0 are the immediate value
	imm := uint32(m.instruction & 0xFF)

	fmt.Printf("mov r%d, #%d\n", rd, imm)

	cpu.WriteRegister(rd, imm)

	cpu.SetN(imm&(1<<31)>>31 == 1)
	cpu.SetZ(imm == 0)
}

type CMP struct {
	instruction uint16
}

func (c CMP) Execute(cpu interfaces.CPU) {
	fmt.Println("CMP")
	// Bits 10-8 are the destination register
	rd := c.instruction & (1<<10 | 1<<9 | 1<<8) >> 8
	// Bits 7-0 are the immediate value
	imm := c.instruction & 0xFF

	fmt.Printf("Destination register: %d\n", rd)
	fmt.Printf("Immediate value: %d\n", imm)

	panic("Not implemented")
}

type ADD struct {
	instruction uint16
}

func (a ADD) Execute(cpu interfaces.CPU) {
	fmt.Println("ADD")
	// Bits 10-8 are the destination register
	rd := a.instruction & (1<<10 | 1<<9 | 1<<8) >> 8
	// Bits 7-0 are the immediate value
	imm := a.instruction & 0xFF

	fmt.Printf("Destination register: %d\n", rd)
	fmt.Printf("Immediate value: %d\n", imm)

	panic("Not implemented")
}

type SUB struct {
	instruction uint16
}

func (s SUB) Execute(cpu interfaces.CPU) {
	fmt.Println("SUB")
	// Bits 10-8 are the destination register
	rd := s.instruction & (1<<10 | 1<<9 | 1<<8) >> 8
	// Bits 7-0 are the immediate value
	imm := s.instruction & 0xFF

	fmt.Printf("Destination register: %d\n", rd)
	fmt.Printf("Immediate value: %d\n", imm)

	panic("Not implemented")
}

type ADD2 struct {
	instruction uint16
}

func (a ADD2) Execute(cpu interfaces.CPU) {
	fmt.Println("ADD2")

	// Bits 8-6 are the immediate value
	imm := uint32((a.instruction & (1<<8 | 1<<7 | 1<<6)) >> 6)
	// Bits 5-3 are the source register
	rs := uint8((a.instruction & (1<<5 | 1<<4 | 1<<3)) >> 3)
	// Bits 2-0 are the destination register
	rd := uint8((a.instruction & (1<<2 | 1<<1 | 1<<0)))

	// bit 10 == 1 means the operand is an immediate value
	if a.instruction&(1<<10)>>10 == 1 {
		fmt.Printf("ADD2 r%d, r%d, #%d\n", rd, rs, imm)
		cpu.WriteRegister(rd, cpu.ReadRegister(rs)+imm)
	} else {
		fmt.Printf("ADD2 r%d, r%d, r%d\n", rd, rs, imm)
		cpu.WriteRegister(rd, cpu.ReadRegister(rs)+cpu.ReadRegister(rd))
	}

	// Save condition flags
	cpu.SetN(cpu.ReadRegister(rd)&(1<<31)>>31 == 1)
	cpu.SetZ(cpu.ReadRegister(rd) == 0)
	fmt.Println("Not setting C or V")
}

type SUB2 struct {
	instruction uint16
}

func (s SUB2) Execute(cpu interfaces.CPU) {
	fmt.Println("SUB2")

	panic("Not implemented")
}
