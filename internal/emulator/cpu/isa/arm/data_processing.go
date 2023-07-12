package arm

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type AND struct {
	instruction uint32
}

func (a AND) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("AND")

	// If bit 25 is set, then the instruction is an immediate operation.
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

	res := rnVal & op2

	fmt.Printf("r%d = r%d [%08X] & %08X = %08X\n", rd, rn, rnVal, op2, res)

	cpu.WriteRegister(rd, res)

	if a.instruction&(1<<20)>>20 == 1 {
		carry := (rnVal>>31)+(op2>>31) > (res >> 31)
		overflow := (rnVal^op2)>>31 == 0 && (rnVal^res)>>31 == 1
		cpu.SetConditionCodes(res, carry, overflow)
	}

	return
}

type EOR struct {
	instruction uint32
}

func (e EOR) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("EOR")
	panic("Not implemented")
	return
}

type SUB struct {
	instruction uint32
}

func (s SUB) Execute(cpu interfaces.CPU) (repipeline bool) {
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

	fmt.Printf("sub r%d, %d = %08X\n", rn, op2, diff)

	cpu.WriteRegister(rn, diff)

	if s.instruction&(1<<20)>>20 == 1 {
		// Set carry flag if the subtraction would make a positive number.
		carry := rnVal >= op2
		// Set overflow flag if the subtraction would overflow.
		overflow := (rnVal^op2)>>31 == 1 && (rnVal^diff)>>31 == 1
		cpu.SetConditionCodes(diff, carry, overflow)
	}
	return
}

type RSB struct {
	instruction uint32
}

func (r RSB) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("RSB")
	// If bit 25 is set, then the instruction is an immediate operation.
	immediate := (r.instruction&(1<<25))>>25 == 1

	// Rn is bits 19-16
	rn := uint8((r.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	// Rd is bits 15-12
	rd := uint8((r.instruction & 0x0000F000) >> 12)

	op2 := uint32(0)
	if immediate {
		op2, _ = unshiftImmediate(r.instruction & 0x00000FFF)
	} else {
		op2, _ = unshiftRegister(r.instruction&0x00000FFF, cpu)
	}

	// Reverse subtract op2 from Rn
	res := op2 - rnVal

	fmt.Printf("r%d = %08X - r%d [%08X] = %08X\n", rd, op2, rn, rnVal, res)

	cpu.WriteRegister(rd, res)

	if r.instruction&(1<<20)>>20 == 1 {
		// Set carry flag if the subtraction would make a positive number.
		carry := op2 >= rnVal
		// Set overflow flag if the subtraction would overflow.
		overflow := (op2^rnVal)>>31 == 1 && (op2^res)>>31 == 1
		cpu.SetConditionCodes(res, carry, overflow)
	}
	return
}

type ADD struct {
	instruction uint32
}

func (a ADD) Execute(cpu interfaces.CPU) (repipeline bool) {
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
	return
}

type ADC struct {
	instruction uint32
}

func (a ADC) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("ADC")
	panic("Not implemented")
	return
}

type SBC struct {
	instruction uint32
}

func (s SBC) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("SBC")
	panic("Not implemented")
	return
}

type RSC struct {
	instruction uint32
}

func (rsc RSC) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("RSC")
	panic("Not implemented")
	return
}

type TST struct {
	instruction uint32
}

func (t TST) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("TST")

	// If bit 25 is set, then the instruction is an immediate operation.
	immediate := (t.instruction&(1<<25))>>25 == 1

	// Rn is bits 19-16
	rn := uint8((t.instruction & 0x000F0000) >> 16)

	rnVal := cpu.ReadRegister(rn)

	op2 := uint32(0)
	carry := false
	if immediate {
		op2, carry = unshiftImmediate(t.instruction & 0x00000FFF)
	} else {
		op2, carry = unshiftRegister(t.instruction&0x00000FFF, cpu)
	}

	res := rnVal & op2

	cpu.SetZ(res == 0)
	cpu.SetN(res>>31 == 1)
	cpu.SetC(carry)

	return
}

type TEQ struct {
	instruction uint32
}

func (t TEQ) Execute(cpu interfaces.CPU) (repipeline bool) {
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
		cpu.SetN(res&(1<<31) != 0)
	}
	return
}

type CMP struct {
	instruction uint32
}

func (c CMP) Execute(cpu interfaces.CPU) (repipeline bool) {
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
	return
}

type CMN struct {
	instruction uint32
}

func (c CMN) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("CMN")
	panic("Not implemented")
	return
}

type ORR struct {
	instruction uint32
}

func (o ORR) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("ORR")
	// If bit 25 is set, then the instruction is an immediate operation.
	immediate := (o.instruction&(1<<25))>>25 == 1

	// Rn is bits 19-16
	rn := uint8((o.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	// Rd is bits 15-12
	rd := uint8((o.instruction & 0x0000F000) >> 12)

	op2 := uint32(0)
	if immediate {
		op2, _ = unshiftImmediate(o.instruction & 0x00000FFF)
	} else {
		op2, _ = unshiftRegister(o.instruction&0x00000FFF, cpu)
	}

	res := rnVal | op2

	cpu.WriteRegister(rd, res)

	cpu.SetN(res>>31 == 1)
	cpu.SetZ(res == 0)

	return
}

type MOV struct {
	instruction uint32
}

func (m MOV) Execute(cpu interfaces.CPU) (repipeline bool) {
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
		cpu.SetN(val&(1<<31) != 0)
		cpu.SetC(carry)
	}
	return
}

type BIC struct {
	instruction uint32
}

func (b BIC) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("BIC")

	// Destination register is bits 15-12
	rd := uint8((b.instruction & 0x0000F000) >> 12)

	// If bit 25 is set, then the instruction is an immediate operation.
	immediate := (b.instruction&(1<<25))>>25 == 1

	// 2nd operand is bits 11-0
	op2 := uint32(0)
	carry := false
	if immediate {
		op2, carry = unshiftImmediate(b.instruction & 0x00000FFF)
	} else {
		op2, carry = unshiftRegister(b.instruction&0x00000FFF, cpu)
	}

	// Rd = Rn AND NOT Op2
	rdVal := cpu.ReadRegister(rd)
	res := rdVal &^ op2
	cpu.WriteRegister(rd, res)

	// If bit 20 is set, then the instruction sets the condition codes.
	if b.instruction&(1<<20)>>20 == 1 {
		cpu.SetZ(res == 0)
		cpu.SetN(res&(1<<31) != 0)
		cpu.SetC(carry)
	}

	return
}

type MVN struct {
	instruction uint32
}

func (m MVN) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("MVN")
	panic("Not implemented")
	return
}
