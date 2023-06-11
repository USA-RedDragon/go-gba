package arm

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/interfaces"
)

type AND struct {
	instruction uint32
}

func (a AND) Execute(cpu interfaces.CPU) {
	fmt.Println("AND")
	panic("Not implemented")
}

type EOR struct {
	instruction uint32
}

func (e EOR) Execute(cpu interfaces.CPU) {
	fmt.Println("EOR")
	panic("Not implemented")
}

type SUB struct {
	instruction uint32
}

func (s SUB) Execute(cpu interfaces.CPU) {
	fmt.Println("SUB")
	// If bit 25 is set, then the instruction is an immediate operation.
	immediate := (s.instruction&(1<<25))>>25 == 1

	// Rn is bits 19-16
	rn := uint8((s.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	op2 := uint32(0)
	if immediate {
		op2, _ = unshiftImmediate(s.instruction & 0x00000FFF)
	} else {
		op2, _ = unshiftRegister(s.instruction&0x00000FFF, cpu)
	}

	// Subtract op2 from Rn and update the condition flags, but do not store the result.
	diff := rnVal - op2

	if s.instruction&(1<<20)>>20 == 1 {
		// Set carry flag if the subtraction would make a positive number.
		carry := rnVal >= op2
		// Set overflow flag if the subtraction would overflow.
		overflow := (rnVal^op2)>>31 == 1 && (rnVal^diff)>>31 == 1
		cpu.SetConditionCodes(diff, carry, overflow)
	}
}

type RSB struct {
	instruction uint32
}

func (rsb RSB) Execute(cpu interfaces.CPU) {
	fmt.Println("RSB")
	panic("Not implemented")
}

type ADD struct {
	instruction uint32
}

func (a ADD) Execute(cpu interfaces.CPU) {
	fmt.Println("ADD")

	immediate := (a.instruction&(1<<25))>>25 == 1

	// Rn is bits 19-16
	rn := uint8((a.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	// Rd is bits 15-12
	rd := uint8((a.instruction & 0x0000F000) >> 12)

	op2 := uint32(0)
	if immediate {
		op2, _ = unshiftImmediate(a.instruction & 0x00000FFF)
	} else {
		op2, _ = unshiftRegister(a.instruction&0x00000FFF, cpu)
	}

	res := rnVal + op2
	fmt.Printf("r%d = r%d [%08X] + %08X = %08X\n", rd, rn, rnVal, op2, res)

	cpu.WriteRegister(rd, res)

	if a.instruction&(1<<20)>>20 == 1 {
		carry := (rnVal>>31)+(op2>>31) > (res >> 31)
		overflow := (rnVal^op2)>>31 == 0 && (rnVal^res)>>31 == 1
		cpu.SetConditionCodes(res, carry, overflow)
	}
}

type ADC struct {
	instruction uint32
}

func (a ADC) Execute(cpu interfaces.CPU) {
	fmt.Println("ADC")
	panic("Not implemented")
}

type SBC struct {
	instruction uint32
}

func (s SBC) Execute(cpu interfaces.CPU) {
	fmt.Println("SBC")
	panic("Not implemented")
}

type RSC struct {
	instruction uint32
}

func (rsc RSC) Execute(cpu interfaces.CPU) {
	fmt.Println("RSC")
	panic("Not implemented")
}

type TST struct {
	instruction uint32
}

func (t TST) Execute(cpu interfaces.CPU) {
	fmt.Println("TST")
	panic("Not implemented")
}

type TEQ struct {
	instruction uint32
}

func (t TEQ) Execute(cpu interfaces.CPU) {
	// If bit 25 is set, then the instruction is an immediate operation.
	immediate := (t.instruction&(1<<25))>>25 == 1

	// Rn is bits 19-16
	rn := uint8((t.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	op2 := uint32(0)
	if immediate {
		op2, _ = unshiftImmediate(t.instruction & 0x00000FFF)
	} else {
		op2, _ = unshiftRegister(t.instruction&0x00000FFF, cpu)
	}

	res := rnVal ^ op2

	if t.instruction&(1<<20)>>20 == 1 {
		cpu.SetZ(res == 0)
		cpu.SetN((res&(1<<31))>>31 == 1)
	}
}

type CMP struct {
	instruction uint32
}

func (c CMP) Execute(cpu interfaces.CPU) {
	fmt.Println("CMP")

	// If bit 25 is set, then the instruction is an immediate operation.
	immediate := (c.instruction&(1<<25))>>25 == 1

	// Rn is bits 19-16
	rn := uint8((c.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	op2 := uint32(0)
	if immediate {
		op2, _ = unshiftImmediate(c.instruction & 0x00000FFF)
	} else {
		op2, _ = unshiftRegister(c.instruction&0x00000FFF, cpu)
	}

	// Subtract op2 from Rn and update the condition flags, but do not store the result.
	diff := rnVal - op2

	if c.instruction&(1<<20)>>20 == 1 {
		// Set carry flag if the subtraction would make a positive number.
		carry := rnVal >= op2
		// Set overflow flag if the subtraction would overflow.
		overflow := (rnVal^op2)>>31 == 1 && (rnVal^diff)>>31 == 1
		cpu.SetConditionCodes(diff, carry, overflow)
	}
}

type CMN struct {
	instruction uint32
}

func (c CMN) Execute(cpu interfaces.CPU) {
	fmt.Println("CMN")
	panic("Not implemented")
}

type ORR struct {
	instruction uint32
}

func (o ORR) Execute(cpu interfaces.CPU) {
	fmt.Println("ORR")
	panic("Not implemented")
}

type MOV struct {
	instruction uint32
}

func (m MOV) Execute(cpu interfaces.CPU) {
	// Destination register is bits 15-12
	destination := uint8((m.instruction & 0x0000F000) >> 12)

	// If bit 25 is set, then the instruction is an immediate operation.
	immediate := (m.instruction&(1<<25))>>25 == 1

	// 2nd operand is bits 11-0
	val := uint32(0)
	carry := false
	if immediate {
		val, carry = unshiftImmediate(m.instruction & 0x00000FFF)
		cpu.WriteRegister(destination, val)
	} else {
		val, carry = unshiftRegister(m.instruction&0x00000FFF, cpu)
		cpu.WriteRegister(destination, val)
	}

	// If bit 20 is set, then the instruction sets the condition codes.
	if m.instruction&(1<<20)>>20 == 1 {
		cpu.SetZ(val == 0)
		cpu.SetN((val&(1<<31))>>31 == 1)
		cpu.SetC(carry)
	}
}

type BIC struct {
	instruction uint32
}

func (b BIC) Execute(cpu interfaces.CPU) {
	fmt.Println("BIC")
	panic("Not implemented")
}

type MVN struct {
	instruction uint32
}

func (m MVN) Execute(cpu interfaces.CPU) {
	fmt.Println("MVN")
	panic("Not implemented")
}
