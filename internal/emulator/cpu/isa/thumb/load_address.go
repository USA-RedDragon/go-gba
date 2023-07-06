package thumb

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type ADDSP struct {
	instruction uint16
}

func (a ADDSP) Execute(cpu interfaces.CPU) (repipeline bool) {
	// Bits 10-8 are the destination register
	rd := uint8(a.instruction & (1<<10 | 1<<9 | 1<<8) >> 8)

	// Bits 7-0 are the immediate value
	imm := (a.instruction & 0xFF) << 2

	fmt.Printf("add r%d, sp, #0x%X\n", rd, imm)

	cpu.WriteRegister(rd, cpu.ReadSP()+uint32(imm))

	return
}

type ADDPC struct {
	instruction uint16
}

func (a ADDPC) Execute(cpu interfaces.CPU) (repipeline bool) {
	// Bits 10-8 are the destination register
	rd := uint8(a.instruction & (1<<10 | 1<<9 | 1<<8) >> 8)

	// Bits 7-0 are the immediate value
	imm := (a.instruction & 0xFF) << 2

	fmt.Printf("add r%d, pc, #0x%X\n", rd, imm)

	cpu.WriteRegister(rd, cpu.ReadPC()+uint32(imm))

	return
}
