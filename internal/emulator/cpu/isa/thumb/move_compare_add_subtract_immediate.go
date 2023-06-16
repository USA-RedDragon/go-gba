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
	rd := uint8(a.instruction & (1<<10 | 1<<9 | 1<<8) >> 8)
	// Bits 7-0 are the immediate value
	imm := uint32(a.instruction & 0xFF)

	fmt.Printf("Destination register: %d\n", rd)
	fmt.Printf("Immediate value: %d\n", imm)

	// Add the immediate value to the destination register
	cpu.WriteRegister(rd, cpu.ReadRegister(rd)+imm)
}

type SUB struct {
	instruction uint16
}

func (s SUB) Execute(cpu interfaces.CPU) {
	fmt.Println("SUB")
	// Bits 10-8 are the destination register
	rd := uint8(s.instruction & (1<<10 | 1<<9 | 1<<8) >> 8)
	// Bits 7-0 are the immediate value
	imm := uint32(s.instruction & 0xFF)

	fmt.Printf("Destination register: %d\n", rd)
	fmt.Printf("Immediate value: %d\n", imm)

	// Subtract the immediate value from the destination register
	cpu.WriteRegister(rd, cpu.ReadRegister(rd)-imm)

	cpu.SetN(cpu.ReadRegister(rd)&(1<<31)>>31 == 1)
	cpu.SetZ(cpu.ReadRegister(rd) == 0)
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

	// Bits 8-6 are the immediate value
	imm := uint32((s.instruction & (1<<8 | 1<<7 | 1<<6)) >> 6)
	// Bits 5-3 are the source register
	rs := uint8((s.instruction & (1<<5 | 1<<4 | 1<<3)) >> 3)
	// Bits 2-0 are the destination register
	rd := uint8((s.instruction & (1<<2 | 1<<1 | 1<<0)))

	// bit 10 == 1 means the operand is an immediate value
	if s.instruction&(1<<10)>>10 == 1 {
		fmt.Printf("SUB2 r%d, r%d, #%d\n", rd, rs, imm)
		cpu.WriteRegister(rd, cpu.ReadRegister(rs)-imm)
	} else {
		fmt.Printf("SUB2 r%d, r%d, r%d\n", rd, rs, imm)
		cpu.WriteRegister(rd, cpu.ReadRegister(rs)-cpu.ReadRegister(rd))
	}

	// Save condition flags
	cpu.SetN(cpu.ReadRegister(rd)&(1<<31)>>31 == 1)
	cpu.SetZ(cpu.ReadRegister(rd) == 0)
}

type SUBSP struct {
	instruction uint16
}

func (a SUBSP) Execute(cpu interfaces.CPU) {
	fmt.Println("SUBSP")

	// Bit 7 == 1 if the offset is negative
	negative := a.instruction&(1<<7)>>7 == 1

	// Bits 6-0 are the immediate offset
	imm := uint32(a.instruction & 0b111111)

	offset := imm << 2

	if negative {
		fmt.Printf("SUBSP #-%d\n", offset)
		cpu.WriteSP(cpu.ReadSP() - offset)
	} else {
		fmt.Printf("SUBSP #%d\n", offset)
		cpu.WriteSP(cpu.ReadSP() + offset)
	}
}
