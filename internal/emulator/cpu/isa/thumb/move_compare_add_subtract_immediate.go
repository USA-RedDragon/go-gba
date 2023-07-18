package thumb

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type MOV struct {
	instruction uint16
}

func (m MOV) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bits 10-8 are the destination register
	rd := uint8(m.instruction & (1<<10 | 1<<9 | 1<<8) >> 8)
	// Bits 7-0 are the immediate value
	imm := uint32(m.instruction & 0xFF)

	fmt.Printf("mov r%d, #%d\n", rd, imm)

	cpu.WriteRegister(rd, imm)

	cpu.SetN(imm&(1<<31)>>31 != 0)
	cpu.SetZ(imm == 0)
	return
}

type CMP struct {
	instruction uint16
}

func (c CMP) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("CMP")
	// Bits 10-8 are the destination register
	rd := uint8(c.instruction & (1<<10 | 1<<9 | 1<<8) >> 8)
	rdVal := cpu.ReadRegister(rd)
	// Bits 7-0 are the immediate value
	imm := uint32(c.instruction & 0xFF)

	fmt.Printf("Destination register: %d\n", rd)
	fmt.Printf("Immediate value: %d\n", imm)

	// Subtract the immediate value from the destination register
	res := rdVal - imm

	// Set carry flag if the subtraction would make a positive number.
	carry := rdVal >= imm

	// Set overflow flag if the subtraction would overflow.
	overflow := (rdVal^imm)>>31 == 1 && (rdVal^res)>>31 == 1
	cpu.SetN(res&(1<<31)>>31 != 0)
	cpu.SetZ(res == 0)
	cpu.SetV(overflow)
	cpu.SetC(carry)

	return
}

type ADD struct {
	instruction uint16
}

func (a ADD) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("ADD")
	// Bits 10-8 are the destination register
	rd := uint8(a.instruction & (1<<10 | 1<<9 | 1<<8) >> 8)
	rdVal := cpu.ReadRegister(rd)
	// Bits 7-0 are the immediate value
	imm := uint32(a.instruction & 0xFF)

	fmt.Printf("Destination register: %d\n", rd)
	fmt.Printf("Immediate value: %d\n", imm)

	// Add the immediate value to the destination register
	res := rdVal + imm
	cpu.WriteRegister(rd, res)

	cpu.SetN(res&(1<<31)>>31 != 0)
	cpu.SetZ(res == 0)
	cpu.SetC(rdVal > imm)
	cpu.SetV((rdVal^imm)>>31 == 0 && (rdVal^res)>>31 == 1)
	return
}

type SUB struct {
	instruction uint16
}

func (s SUB) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("SUB")
	// Bits 10-8 are the destination register
	rd := uint8(s.instruction & (1<<10 | 1<<9 | 1<<8) >> 8)
	rdVal := cpu.ReadRegister(rd)
	// Bits 7-0 are the immediate value
	imm := uint32(s.instruction & 0xFF)

	fmt.Printf("Destination register: %d\n", rd)
	fmt.Printf("Immediate value: %d\n", imm)

	// Subtract the immediate value from the destination register
	res := rdVal - imm
	cpu.WriteRegister(rd, res)

	cpu.SetN(res&(1<<31)>>31 != 0)
	cpu.SetZ(res == 0)
	cpu.SetC(rdVal >= imm)
	cpu.SetV((rdVal^imm)>>31 == 1 && (rdVal^res)>>31 == 1)
	return
}

type ADD2 struct {
	instruction uint16
}

func (a ADD2) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("ADD2")

	// Bits 8-6 are the immediate value
	imm := uint8((a.instruction & (1<<8 | 1<<7 | 1<<6)) >> 6)
	// Bits 5-3 are the source register
	rs := uint8((a.instruction & (1<<5 | 1<<4 | 1<<3)) >> 3)
	// Bits 2-0 are the destination register
	rd := uint8((a.instruction & (1<<2 | 1<<1 | 1<<0)))

	rsVal := cpu.ReadRegister(rs)

	uint64Val := uint64(0)

	// bit 10 == 1 means the operand is an immediate value
	if a.instruction&(1<<10)>>10 == 1 {
		fmt.Printf("ADD2 r%d, r%d, #%d\n", rd, rs, imm)
		cpu.WriteRegister(rd, rsVal+uint32(imm))

		// Set the C flag if the addition overflowed
		uint64Val = uint64(rsVal) + uint64(imm)
	} else {
		fmt.Printf("ADD2 r%d, r%d, r%d\n", rd, rs, imm)
		immVal := cpu.ReadRegister(imm)
		cpu.WriteRegister(rd, rsVal+immVal)

		// Set the C flag if the addition overflowed
		uint64Val = uint64(rsVal) + uint64(immVal)
	}

	if uint64Val > 0xFFFFFFFF {
		cpu.SetC(true)
	} else {
		cpu.SetC(false)
	}

	rdVal := cpu.ReadRegister(rd)

	// Save condition flags
	cpu.SetN(rdVal&(1<<31)>>31 == 1)
	cpu.SetZ(rdVal == 0)

	fmt.Println("Not setting V flag")
	return
}

type SUB2 struct {
	instruction uint16
}

func (s SUB2) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("SUB2")

	// Bits 8-6 are the immediate value
	imm := uint32((s.instruction & (1<<8 | 1<<7 | 1<<6)) >> 6)
	// Bits 5-3 are the source register
	rs := uint8((s.instruction & (1<<5 | 1<<4 | 1<<3)) >> 3)
	rsVal := cpu.ReadRegister(rs)
	// Bits 2-0 are the destination register
	rd := uint8((s.instruction & (1<<2 | 1<<1 | 1<<0)))

	// bit 10 == 1 means the operand is an immediate value
	if s.instruction&(1<<10)>>10 == 1 {
		fmt.Printf("SUB2 r%d, r%d, #%d\n", rd, rs, imm)
	} else {
		fmt.Printf("SUB2 r%d, r%d, r%d\n", rd, rs, imm)
		imm = cpu.ReadRegister(uint8(imm))
	}

	res := rsVal - imm
	cpu.WriteRegister(rd, res)

	cpu.SetN(res&(1<<31)>>31 != 0)
	cpu.SetZ(res == 0)
	cpu.SetC(rsVal >= imm)
	cpu.SetV((rsVal^imm)>>31 == 1 && (rsVal^res)>>31 == 1)

	return
}

type SUBSP struct {
	instruction uint16
}

func (a SUBSP) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("SUBSP")

	// Bit 7 == 1 if the offset is negative
	negative := a.instruction&(1<<7)>>7 == 1

	// Bits 6-0 are the immediate offset
	imm := uint32(a.instruction & 0b111111)

	offset := imm << 2

	sp := cpu.ReadSP()

	if negative {
		fmt.Printf("SUBSP #-%d\n", offset)
		cpu.WriteSP(sp - offset)
		// Set the C flag if sp - offset underflows
		if sp-offset > sp {
			cpu.SetC(false)
		} else {
			cpu.SetC(true)
		}
	} else {
		fmt.Printf("SUBSP #%d\n", offset)
		cpu.WriteSP(sp + offset)
		// Set the C flag if sp + offset overflows
		if sp+offset < sp {
			cpu.SetC(false)
		} else {
			cpu.SetC(true)
		}
	}
	return
}
