package arm

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type UMULL struct {
	instruction uint32
}

func (u UMULL) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bit 20 is the S flag to update the CPSR
	s := u.instruction&(1<<20) != 0

	// Bits 19-16 are the RdHi register
	rdHi := uint8((u.instruction >> 16) & 0xf)

	// Bits 15-12 are the RdLo register
	rdLo := uint8((u.instruction >> 12) & 0xf)

	// Bits 11-8 are the Rs register
	rs := uint8((u.instruction >> 8) & 0xf)
	rsVal := uint64(cpu.ReadRegister(rs))

	// Buts 3-0 are the Rm register
	rm := uint8(u.instruction & 0xf)
	rmVal := uint64(cpu.ReadRegister(rm))

	fmt.Printf("umull r%d, r%d, r%d, r%d\n", rdLo, rdHi, rm, rs)

	res := rsVal * rmVal

	cpu.WriteRegister(rdHi, uint32((res>>32)&0xffffffff))
	cpu.WriteRegister(rdLo, uint32(res&0xffffffff))

	if s {
		cpu.SetN(res&(1<<63) != 0)
		cpu.SetZ(res == 0)
		// Both the C and V flags are set to meaningless values.
	}
	return
}

type UMLAL struct {
	instruction uint32
}

func (u UMLAL) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bit 20 is the S flag to update the CPSR
	s := u.instruction&(1<<20) != 0

	// Bits 19-16 are the RdHi register
	rdHi := uint8((u.instruction >> 16) & 0xf)
	rdHiVal := uint64(cpu.ReadRegister(rdHi))

	// Bits 15-12 are the RdLo register
	rdLo := uint8((u.instruction >> 12) & 0xf)
	rdLoVal := uint64(cpu.ReadRegister(rdLo))

	accVal := (rdHiVal << 32) | rdLoVal

	// Bits 11-8 are the Rs register
	rs := uint8((u.instruction >> 8) & 0xf)
	rsVal := uint64(cpu.ReadRegister(rs))

	// Buts 3-0 are the Rm register
	rm := uint8(u.instruction & 0xf)
	rmVal := uint64(cpu.ReadRegister(rm))

	fmt.Printf("umlal r%d, r%d, r%d, r%d\n", rdLo, rdHi, rm, rs)

	res := rsVal*rmVal + accVal

	cpu.WriteRegister(rdHi, uint32((res>>32)&0xffffffff))
	cpu.WriteRegister(rdLo, uint32(res&0xffffffff))

	if s {
		cpu.SetN(res&(1<<63) != 0)
		cpu.SetZ(res == 0)
		// Both the C and V flags are set to meaningless values.
	}
	return
}

type SMULL struct {
	instruction uint32
}

func (s SMULL) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bit 20 is the S flag to update the CPSR
	save := s.instruction&(1<<20) != 0

	// Bits 19-16 are the RdHi register
	rdHi := uint8((s.instruction >> 16) & 0xf)

	// Bits 15-12 are the RdLo register
	rdLo := uint8((s.instruction >> 12) & 0xf)

	// Bits 11-8 are the Rs register
	rs := uint8((s.instruction >> 8) & 0xf)
	rsVal := int32(cpu.ReadRegister(rs))

	// Buts 3-0 are the Rm register
	rm := uint8(s.instruction & 0xf)
	rmVal := int32(cpu.ReadRegister(rm))

	fmt.Printf("smull r%d, r%d, r%d, r%d\n", rdLo, rdHi, rm, rs)

	var res int64 = int64(rsVal) * int64(rmVal)

	cpu.WriteRegister(rdHi, uint32((res>>32)&0xffffffff))
	cpu.WriteRegister(rdLo, uint32(res&0xffffffff))

	if save {
		res := uint64(res)
		cpu.SetN(res&(1<<63) != 0)
		cpu.SetZ(res == 0)
		// Both the C and V flags are set to meaningless values.
	}
	return
}

type SMLAL struct {
	instruction uint32
}

func (s SMLAL) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bit 20 is the S flag to update the CPSR
	save := s.instruction&(1<<20) != 0

	// Bits 19-16 are the RdHi register
	rdHi := uint8((s.instruction >> 16) & 0xf)
	rdHiVal := int32(cpu.ReadRegister(rdHi))

	// Bits 15-12 are the RdLo register
	rdLo := uint8((s.instruction >> 12) & 0xf)
	rdLoVal := int32(cpu.ReadRegister(rdLo))

	accVal := int64((uint64(rdHiVal) << 32) | uint64(rdLoVal))

	// Bits 11-8 are the Rs register
	rs := uint8((s.instruction >> 8) & 0xf)
	rsVal := int32(cpu.ReadRegister(rs))

	// Buts 3-0 are the Rm register
	rm := uint8(s.instruction & 0xf)
	rmVal := int32(cpu.ReadRegister(rm))

	fmt.Printf("smlal r%d, r%d, r%d, r%d\n", rdLo, rdHi, rm, rs)

	var res int64 = int64(rsVal)*int64(rmVal) + accVal

	cpu.WriteRegister(rdHi, uint32((res>>32)&0xffffffff))
	cpu.WriteRegister(rdLo, uint32(res&0xffffffff))

	if save {
		res := uint64(res)
		cpu.SetN(res&(1<<63) != 0)
		cpu.SetZ(res == 0)
		// Both the C and V flags are set to meaningless values.
	}
	return
}
