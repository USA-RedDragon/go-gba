package thumb

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type AND struct {
	instruction uint16
}

func (a AND) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("AND")

	// Bits 5-3 are the source register
	rs := uint8(a.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(a.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("and r%d, r%d\n", rd, rs)

	panic("Not implemented")
	return
}

type EOR struct {
	instruction uint16
}

func (e EOR) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("EOR")

	// Bits 5-3 are the source register
	rs := uint8(e.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(e.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("eor r%d, r%d\n", rd, rs)

	panic("Not implemented")
	return
}

type LSL struct {
	instruction uint16
}

func (l LSL) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("LSL")

	// Bits 5-3 are the source register
	rs := uint8(l.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(l.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("lsl r%d, r%d\n", rd, rs)

	panic("Not implemented")
	return
}

type LSR struct {
	instruction uint16
}

func (l LSR) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("LSR")

	// Bits 5-3 are the source register
	rs := uint8(l.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(l.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("lsr r%d, r%d\n", rd, rs)

	panic("Not implemented")
	return
}

type ASR struct {
	instruction uint16
}

func (a ASR) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("ASR")

	// Bits 5-3 are the source register
	rs := uint8(a.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(a.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("asr r%d, r%d\n", rd, rs)

	panic("Not implemented")
	return
}

type ADC struct {
	instruction uint16
}

func (a ADC) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("ADC")

	// Bits 5-3 are the source register
	rs := uint8(a.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(a.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("adc r%d, r%d\n", rd, rs)

	panic("Not implemented")
	return
}

type SBC struct {
	instruction uint16
}

func (s SBC) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("SBC")

	// Bits 5-3 are the source register
	rs := uint8(s.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(s.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("sbc r%d, r%d\n", rd, rs)

	panic("Not implemented")
	return
}

type ROR struct {
	instruction uint16
}

func (r ROR) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("ROR")

	// Bits 5-3 are the source register
	rs := uint8(r.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(r.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("ror r%d, r%d\n", rd, rs)

	panic("Not implemented")
	return
}

type TST struct {
	instruction uint16
}

func (t TST) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("TST")

	// Bits 5-3 are the source register
	rs := uint8(t.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(t.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("tst r%d, r%d\n", rd, rs)

	// TST performs a bitwise AND on the two registers, but does not store the result
	// It only updates the status register
	res := cpu.ReadRegister(rd) & cpu.ReadRegister(rs)

	// Update the status registers
	cpu.SetZ(res == 0)
	// Set N if the result is negative
	cpu.SetN(res&(1<<31) != 0)
	return
}

type NEG struct {
	instruction uint16
}

func (n NEG) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("NEG")

	// Bits 5-3 are the source register
	rs := uint8(n.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(n.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("neg r%d, r%d\n", rd, rs)

	panic("Not implemented")
	return
}

type CMPALU struct {
	instruction uint16
}

func (c CMPALU) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("CMP")

	// Bits 5-3 are the source register
	rs := uint8(c.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(c.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("cmp r%d, r%d\n", rd, rs)

	panic("Not implemented")
	return
}

type CMN struct {
	instruction uint16
}

func (c CMN) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("CMN")

	// Bits 5-3 are the source register
	rs := uint8(c.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(c.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("cmn r%d, r%d\n", rd, rs)

	panic("Not implemented")
	return
}

type ORR struct {
	instruction uint16
}

func (o ORR) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("ORR")

	// Bits 5-3 are the source register
	rs := uint8(o.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(o.instruction & (1<<2 | 1<<1 | 1<<0))

	rdVal := cpu.ReadRegister(rd)
	rsVal := cpu.ReadRegister(rs)
	res := rdVal | rsVal
	cpu.WriteRegister(rd, res)

	fmt.Printf("orr r%d, r%d\n", rd, rs)

	// update the status registers
	cpu.SetZ(res == 0)
	cpu.SetN(res&(1<<31) != 0)
	return
}

type MUL struct {
	instruction uint16
}

func (m MUL) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("MUL")

	// Bits 5-3 are the source register
	rs := uint8(m.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(m.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("mul r%d, r%d\n", rd, rs)

	panic("Not implemented")
	return
}

type BIC struct {
	instruction uint16
}

func (b BIC) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("BIC")

	// Bits 5-3 are the source register
	rs := uint8(b.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(b.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("bic r%d, r%d\n", rd, rs)

	panic("Not implemented")
	return
}

type MVN struct {
	instruction uint16
}

func (m MVN) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("MVN")

	// Bits 5-3 are the source register
	rs := uint8(m.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(m.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("mvns r%d, r%d\n", rd, rs)

	// Store the bitwise inverse of the source register in the destination register
	cpu.WriteRegister(rd, ^cpu.ReadRegister(rs))

	// N flag is set if the result is negative
	cpu.SetN(cpu.ReadRegister(rd)&(1<<31) != 0)
	return
}

type LSLMoveShifted struct {
	instruction uint16
}

func (l LSLMoveShifted) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("LSL")

	// Bits 10-6 are the offset
	offset := uint8(l.instruction & (1<<10 | 1<<9 | 1<<8 | 1<<7 | 1<<6) >> 6)

	// Bits 5-3 are the source register
	rs := uint8(l.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(l.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("lsls r%d, r%d, #0x%X\n", rd, rs, offset)

	// Shift the source register left by the offset and store the result in the destination register
	cpu.WriteRegister(rd, cpu.ReadRegister(rs)<<offset)

	// Update the CPSR
	cpu.SetZ(cpu.ReadRegister(rd) == 0)
	cpu.SetN(cpu.ReadRegister(rd)&(1<<31) != 0)
	return
}

type LSRMoveShifted struct {
	instruction uint16
}

func (l LSRMoveShifted) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("LSR")

	// Bits 10-6 are the offset
	offset := uint8(l.instruction & (1<<10 | 1<<9 | 1<<8 | 1<<7 | 1<<6) >> 6)

	// Bits 5-3 are the source register
	rs := uint8(l.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(l.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("lsrs r%d, r%d, #0x%X\n", rd, rs, offset)

	// Shift the source register right by the offset and store the result in the destination register
	cpu.WriteRegister(rd, cpu.ReadRegister(rs)>>offset)

	// Update the CPSR
	cpu.SetZ(cpu.ReadRegister(rd) == 0)
	cpu.SetN(cpu.ReadRegister(rd)&(1<<31) != 0)
	return
}

type ASRMoveShifted struct {
	instruction uint16
}

func (a ASRMoveShifted) Execute(cpu interfaces.CPU) (repipeline bool) {
	fmt.Println("ASR")

	// Bits 10-6 are the offset
	offset := uint8(a.instruction & (1<<10 | 1<<9 | 1<<8 | 1<<7 | 1<<6) >> 6)

	// Bits 5-3 are the source register
	rs := uint8(a.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	// Bits 2-0 are the destination register
	rd := uint8(a.instruction & (1<<2 | 1<<1 | 1<<0))

	fmt.Printf("asrs r%d, r%d, #0x%X\n", rd, rs, offset)

	panic("Not implemented")
	return
}
