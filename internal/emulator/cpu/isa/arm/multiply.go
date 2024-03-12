package arm

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type MLA struct {
	instruction uint32
}

func (m MLA) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bit 20 determines whether or not to update the condition codes
	updateConditionCodes := (m.instruction&(1<<20))>>20 == 1

	// Bits 19-16 are the Rd register
	rd := uint8((m.instruction & 0x000F0000) >> 16)

	// Bits 15-12 are the Rm register
	rn := uint8((m.instruction & 0x0000F000) >> 12)

	// Bits 11-8 are the Rs register
	rs := uint8((m.instruction & 0x00000F00) >> 8)

	// Bits 3-0 are the Rn register
	rm := uint8(m.instruction & 0x0000000F)

	if cpu.GetConfig().Debug {
		//nolint:golint,dupword
		fmt.Printf("mla r%d, r%d, r%d, r%d\n", rd, rm, rs, rn)
	}

	// Rd := Rm * Rs + Rn
	rmVal := cpu.ReadRegister(rm)
	rsVal := cpu.ReadRegister(rs)
	rnVal := cpu.ReadRegister(rn)

	res := rmVal*rsVal + rnVal

	if updateConditionCodes {
		cpu.SetZ(res == 0)
		cpu.SetN(res&(1<<31)>>31 != 0)
	}

	cpu.WriteRegister(rd, res)

	return
}

type MUL struct {
	instruction uint32
}

func (m MUL) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bit 20 determines whether or not to update the condition codes
	updateConditionCodes := (m.instruction&(1<<20))>>20 == 1

	// Bits 19-16 are the Rd register
	rd := uint8((m.instruction & 0x000F0000) >> 16)

	// Bits 11-8 are the Rs register
	rs := uint8((m.instruction & 0x00000F00) >> 8)

	// Bits 3-0 are the Rn register
	rm := uint8(m.instruction & 0x0000000F)

	if cpu.GetConfig().Debug {
		//nolint:golint,dupword
		fmt.Printf("mul r%d, r%d, r%d\n", rd, rm, rs)
	}

	// Rd := Rm * Rs + Rn
	rmVal := cpu.ReadRegister(rm)
	rsVal := cpu.ReadRegister(rs)

	res := rmVal * rsVal

	if updateConditionCodes {
		cpu.SetZ(res == 0)
		cpu.SetN(res&(1<<31)>>31 != 0)
	}

	cpu.WriteRegister(rd, res)

	return
}
