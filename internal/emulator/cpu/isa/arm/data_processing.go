package arm

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type AND struct {
	instruction uint32
}

func (a AND) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	if cpu.GetConfig().Debug {
		fmt.Println("AND")
	}

	// Rn is bits 19-16
	rn := uint8((a.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	// Rd is bits 15-12
	rd := uint8((a.instruction & 0x0000F000) >> 12)

	op2 := ALUOp2(a.instruction, cpu)
	res := rnVal & op2

	if cpu.GetConfig().Debug {
		fmt.Printf("r%d = r%d [%08X] & %08X = %08X\n", rd, rn, rnVal, op2, res)
	}

	cpu.WriteRegister(rd, res)

	if a.instruction&(1<<20)>>20 == 1 {
		overflow := (rnVal^op2)>>31 == 0 && (rnVal^res)>>31 == 1
		cpu.SetZ(res == 0)
		cpu.SetN(res&(1<<31)>>31 != 0)
		cpu.SetV(overflow)
		if rd == 15 {
			cpu.WriteCPSR(cpu.ReadSPSR())
		}
	}

	return
}

type EOR struct {
	instruction uint32
}

func (e EOR) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Rn is bits 19-16
	rn := uint8((e.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	op2 := ALUOp2(e.instruction, cpu)

	// Rd is bits 15-12
	rd := uint8((e.instruction & 0x0000F000) >> 12)

	res := rnVal ^ op2

	if cpu.GetConfig().Debug {
		fmt.Printf("eor r%d, r%d, %d = %08X\n", rd, rn, op2, res)
	}

	cpu.WriteRegister(rd, res)

	if e.instruction&(1<<20)>>20 == 1 {
		cpu.SetZ(res == 0)
		cpu.SetN(res&(1<<31)>>31 != 0)
		if rd == 15 {
			cpu.WriteCPSR(cpu.ReadSPSR())
		}
	}

	return
}

type SUB struct {
	instruction uint32
}

func (s SUB) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Rn is bits 19-16
	rn := uint8((s.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	op2 := ALUOp2(s.instruction, cpu)

	// Subtract op2 from Rn and update the condition flags, but do not store the result.
	diff := rnVal - op2

	if cpu.GetConfig().Debug {
		fmt.Printf("sub r%d, %d = %08X\n", rn, op2, diff)
	}

	cpu.WriteRegister(rn, diff)

	if s.instruction&(1<<20)>>20 == 1 {
		// Set carry flag if the subtraction would make a positive number.
		carry := rnVal >= op2
		// Set overflow flag if the subtraction would overflow.
		overflow := (rnVal^op2)>>31 == 1 && (rnVal^diff)>>31 == 1
		cpu.SetN(diff&(1<<31)>>31 != 0)
		cpu.SetZ(diff == 0)
		cpu.SetV(overflow)
		cpu.SetC(carry)
	}
	return
}

type RSB struct {
	instruction uint32
}

func (r RSB) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	if cpu.GetConfig().Debug {
		fmt.Println("RSB")
	}
	// Rn is bits 19-16
	rn := uint8((r.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	// Rd is bits 15-12
	rd := uint8((r.instruction & 0x0000F000) >> 12)

	op2 := ALUOp2(r.instruction, cpu)

	// Reverse subtract op2 from Rn
	res := op2 - rnVal

	if cpu.GetConfig().Debug {
		fmt.Printf("r%d = %08X - r%d [%08X] = %08X\n", rd, op2, rn, rnVal, res)
	}

	cpu.WriteRegister(rd, res)

	if r.instruction&(1<<20)>>20 == 1 {
		// Set carry flag if the subtraction would make a positive number.
		carry := op2 >= rnVal
		// Set overflow flag if the subtraction would overflow.
		overflow := (op2^rnVal)>>31 == 1 && (op2^res)>>31 == 1
		cpu.SetN(res&(1<<31)>>31 != 0)
		cpu.SetZ(res == 0)
		cpu.SetV(overflow)
		cpu.SetC(carry)
		if rd == 15 {
			cpu.WriteCPSR(cpu.ReadSPSR())
		}
	}
	return
}

type ADD struct {
	instruction uint32
}

func (a ADD) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	if cpu.GetConfig().Debug {
		fmt.Println("ADD")
	}

	// Rn is bits 19-16
	rn := uint8((a.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	// Rd is bits 15-12
	rd := uint8((a.instruction & 0x0000F000) >> 12)

	op2 := ALUOp2(a.instruction, cpu)

	res := rnVal + op2
	if cpu.GetConfig().Debug {
		fmt.Printf("r%d = r%d [%08X] + %08X = %08X\n", rd, rn, rnVal, op2, res)
	}

	cpu.WriteRegister(rd, res)

	if a.instruction&(1<<20)>>20 == 1 {
		carry := (rnVal>>31)+(op2>>31) > (res >> 31)
		overflow := (rnVal^op2)>>31 == 0 && (rnVal^res)>>31 == 1
		cpu.SetN(res&(1<<31)>>31 != 0)
		cpu.SetZ(res == 0)
		cpu.SetV(overflow)
		cpu.SetC(carry)
		if rd == 15 {
			cpu.WriteCPSR(cpu.ReadSPSR())
		}
	}
	return
}

type ADC struct {
	instruction uint32
}

func (a ADC) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	if cpu.GetConfig().Debug {
		fmt.Println("ADC")
	}

	// Rn is bits 19-16
	rn := uint8((a.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	// Rd is bits 15-12
	rd := uint8((a.instruction & 0x0000F000) >> 12)

	op2 := ALUOp2(a.instruction, cpu)

	res := rnVal + op2
	if cpu.GetConfig().Debug {
		fmt.Printf("r%d = r%d [%08X] + %08X = %08X\n", rd, rn, rnVal, op2, res)
	}

	if cpu.GetC() {
		res++
	}

	cpu.WriteRegister(rd, res)

	if a.instruction&(1<<20)>>20 == 1 {
		carry := (rnVal>>31)+(op2>>31) > (res >> 31)
		overflow := (rnVal^op2)>>31 == 0 && (rnVal^res)>>31 == 1
		cpu.SetN(res&(1<<31)>>31 != 0)
		cpu.SetZ(res == 0)
		cpu.SetV(overflow)
		cpu.SetC(carry)
		if rd == 15 {
			cpu.WriteCPSR(cpu.ReadSPSR())
		}
	}
	return
}

type SBC struct {
	instruction uint32
}

func (s SBC) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	if cpu.GetConfig().Debug {
		fmt.Println("SBC")
	}

	// Rn is bits 19-16
	rn := uint8((s.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	op2 := ALUOp2(s.instruction, cpu)

	// Subtract op2 from Rn and update the condition flags, but do not store the result.
	diff := rnVal - op2

	if cpu.GetConfig().Debug {
		fmt.Printf("sbc r%d, %d = %08X\n", rn, op2, diff)
	}

	if !cpu.GetC() {
		diff--
	}

	cpu.WriteRegister(rn, diff)

	if s.instruction&(1<<20)>>20 == 1 {
		// Set carry flag if the subtraction would make a positive number.
		carry := rnVal >= op2
		// Set overflow flag if the subtraction would overflow.
		overflow := (rnVal^op2)>>31 == 1 && (rnVal^diff)>>31 == 1
		cpu.SetN(diff&(1<<31)>>31 != 0)
		cpu.SetZ(diff == 0)
		cpu.SetV(overflow)
		cpu.SetC(carry)
	}

	return
}

type RSC struct {
	instruction uint32
}

func (rsc RSC) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("RSC")
	// Rn is bits 19-16
	rn := uint8((rsc.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	// Rd is bits 15-12
	rd := uint8((rsc.instruction & 0x0000F000) >> 12)

	op2 := ALUOp2(rsc.instruction, cpu)

	// Rd:= Op2 - Op1 + C - 1
	c := uint32(0)
	if cpu.GetC() {
		c = 1
	}
	res := op2 - rnVal + c - 1

	cpu.WriteRegister(rd, res)

	if rsc.instruction&(1<<20)>>20 == 1 {
		// Set carry flag if the subtraction would make a positive number.
		carry := (rnVal + c - 1) >= op2
		// Set overflow flag if the subtraction would overflow.
		overflow := (rnVal^op2)>>31 == 1 && (rnVal^res)>>31 == 1

		cpu.SetN(res&(1<<31)>>31 != 0)
		cpu.SetZ(res == 0)
		cpu.SetV(overflow)
		cpu.SetC(carry)
		if rd == 15 {
			cpu.WriteCPSR(cpu.ReadSPSR())
		}
	}

	return
}

type TST struct {
	instruction uint32
}

func (t TST) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Rn is bits 19-16
	rn := uint8((t.instruction & 0x000F0000) >> 16)

	rnVal := cpu.ReadRegister(rn)

	op2 := ALUOp2(t.instruction, cpu)

	if cpu.GetConfig().Debug {
		fmt.Printf("tst r%d [0x%08X] & 0x%08X = 0x%08X\n", rn, rnVal, op2, rnVal&op2)
	}

	res := rnVal & op2

	cpu.SetZ(res == 0)
	cpu.SetN(res&(1<<31)>>31 != 0)

	return
}

type TEQ struct {
	instruction uint32
}

func (t TEQ) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Rn is bits 19-16
	rn := uint8((t.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	op2 := ALUOp2(t.instruction, cpu)

	res := rnVal ^ op2

	if t.instruction&(1<<20)>>20 == 1 {
		cpu.SetZ(res == 0)
		cpu.SetN(res&(1<<31)>>31 != 0)
	}
	return
}

type CMP struct {
	instruction uint32
}

func (c CMP) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	if cpu.GetConfig().Debug {
		fmt.Println("CMP")
	}

	// Rn is bits 19-16
	rn := uint8((c.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	op2 := ALUOp2(c.instruction, cpu)

	// Subtract op2 from Rn and update the condition flags, but do not store the result.
	diff := rnVal - op2

	fmt.Printf("cmp r%d, %d = %08X\n", rn, op2, diff)

	if c.instruction&(1<<20)>>20 == 1 {
		// Set carry flag if the subtraction would make a positive number.
		carry := rnVal >= op2

		// Set overflow flag if the subtraction would overflow.
		overflow := (rnVal^op2)>>31 == 1 && (rnVal^diff)>>31 == 1
		cpu.SetN(diff&(1<<31)>>31 != 0)
		cpu.SetZ(diff == 0)
		cpu.SetV(overflow)
		cpu.SetC(carry)
	}
	return
}

type CMN struct {
	instruction uint32
}

func (c CMN) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	if cpu.GetConfig().Debug {
		fmt.Println("CMN")
	}

	// Rn is bits 19-16
	rn := uint8((c.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	op2 := ALUOp2(c.instruction, cpu)

	res := rnVal + op2

	if c.instruction&(1<<20)>>20 == 1 {
		carry := (rnVal>>31)+(op2>>31) > (res >> 31)
		overflow := (rnVal^op2)>>31 == 0 && (rnVal^res)>>31 == 1
		cpu.SetN(res&(1<<31)>>31 != 0)
		cpu.SetZ(res == 0)
		cpu.SetV(overflow)
		cpu.SetC(carry)
	}
	return
}

type ORR struct {
	instruction uint32
}

func (o ORR) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Rn is bits 19-16
	rn := uint8((o.instruction & 0x000F0000) >> 16)
	rnVal := cpu.ReadRegister(rn)

	// Rd is bits 15-12
	rd := uint8((o.instruction & 0x0000F000) >> 12)

	op2 := ALUOp2(o.instruction, cpu)

	if cpu.GetConfig().Debug {
		fmt.Printf("orr r%d [0x%08X] | 0x%08X = 0x%08X\n", rd, rnVal, op2, rnVal|op2)
	}

	res := rnVal | op2

	cpu.WriteRegister(rd, res)

	cpu.SetN(res&(1<<31)>>31 != 0)
	cpu.SetZ(res == 0)

	return
}

type MOV struct {
	instruction uint32
}

func (m MOV) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Destination register is bits 15-12
	destination := uint8((m.instruction & 0x0000F000) >> 12)

	// 2nd operand is bits 11-0
	op2 := ALUOp2(m.instruction, cpu)

	fmt.Printf("mov r%d 0x%08X\n", destination, op2)

	cpu.WriteRegister(destination, op2)

	// If bit 20 is set, then the instruction sets the condition codes.
	if m.instruction&(1<<20)>>20 == 1 {
		cpu.SetZ(op2 == 0)
		cpu.SetN(op2&(1<<31)>>31 != 0)
		if destination == 15 {
			cpu.WriteCPSR(cpu.ReadSPSR())
		}
	}
	return
}

type BIC struct {
	instruction uint32
}

func (b BIC) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	if cpu.GetConfig().Debug {
		fmt.Println("BIC")
	}

	// Destination register is bits 15-12
	rd := uint8((b.instruction & 0x0000F000) >> 12)

	// 2nd operand is bits 11-0
	op2 := ALUOp2(b.instruction, cpu)

	// Rd = Rn AND NOT Op2
	rdVal := cpu.ReadRegister(rd)
	res := rdVal &^ op2
	cpu.WriteRegister(rd, res)

	// If bit 20 is set, then the instruction sets the condition codes.
	if b.instruction&(1<<20)>>20 == 1 {
		cpu.SetZ(res == 0)
		cpu.SetN(op2&(1<<31)>>31 != 0)
		if rd == 15 {
			cpu.WriteCPSR(cpu.ReadSPSR())
		}
	}

	return
}

type MVN struct {
	instruction uint32
}

func (m MVN) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("MVN")
	// Destination register is bits 15-12
	destination := uint8((m.instruction & 0x0000F000) >> 12)

	// 2nd operand is bits 11-0
	op2 := ALUOp2(m.instruction, cpu)
	val := uint32(0)
	cpu.WriteRegister(destination, ^op2)

	// If bit 20 is set, then the instruction sets the condition codes.
	if m.instruction&(1<<20)>>20 == 1 {
		cpu.SetZ(val == 0)
		cpu.SetN(op2&(1<<31)>>31 != 0)
		if destination == 15 {
			cpu.WriteCPSR(cpu.ReadSPSR())
		}
	}

	return
}
