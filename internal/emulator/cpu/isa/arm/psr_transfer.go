package arm

import (
	"fmt"
	"math/bits"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type MSR struct {
	instruction uint32
}

func (m MSR) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	if cpu.GetConfig().Debug {
		fmt.Println("PSR Transfer MSR")
	}
	mask := uint32(0)
	if c := m.instruction&(1<<16)>>16 == 1; c {
		mask = 0x0000_00ff
	}
	if x := m.instruction&(1<<17)>>17 == 1; x {
		mask |= 0x0000_ff00
	}
	if s := m.instruction&(1<<18)>>18 == 1; s {
		mask |= 0x00ff_0000
	}
	if f := m.instruction&(1<<19)>>19 == 1; f {
		mask |= 0xff00_0000
	}

	secMask := uint32(0xf8ff03df)
	if cpu.ReadCPSR()&0x1F == 0b10000 {
		secMask = 0xf8ff0000
	}

	r := m.instruction&(1<<22)>>22 == 1
	if r {
		secMask |= 0x01000020
	}

	mask &= secMask
	psr := uint32(0)
	if m.instruction&(1<<25)>>25 == 1 {
		// register Psr[field] = Imm
		is, imm := ((m.instruction>>8)&0b1111)*2, m.instruction&0b1111_1111
		psr = bits.RotateLeft32(imm, -int(is))
	} else {
		// immediate Psr[field] = Rm
		rm := m.instruction & 0b1111
		psr = cpu.ReadRegister(uint8(rm))
	}
	psr &= mask

	if r {
		spsr := cpu.ReadSPSR()
		cpu.WriteSPSR((spsr & ^mask) | psr)
	} else {
		cpsr := cpu.ReadCPSR()
		cpsr &= ^mask
		cpsr |= psr
		cpu.WriteCPSR(cpsr)
	}

	return
}

type MRS struct {
	instruction uint32
}

func (m MRS) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	if cpu.GetConfig().Debug {
		fmt.Println("PSR Transfer MRS")
	}

	// Bits 15-12 are the destination register
	rd := uint8((m.instruction & 0x0000F000) >> 12)

	// Bit 22 is the source
	spsr := m.instruction&(1<<22)>>22 == 1

	if spsr {
		cpu.WriteRegister(rd, cpu.ReadSPSR())
	} else {
		cpu.WriteRegister(rd, cpu.ReadRegister(16))
	}
	return
}
