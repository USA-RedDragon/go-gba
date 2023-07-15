package arm

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type MSR struct {
	instruction uint32
}

func (m MSR) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	if cpu.GetConfig().Debug {
		fmt.Println("PSR Transfer MSR")
	}
	// Bit 22 is destination
	spsr := m.instruction&(1<<22)>>22 == 1
	// If bit 16 is 1, we're doing a register transfer
	test := m.instruction & (1 << 16) >> 16
	if test == 1 {
		// Register contents to PSR
		rm := uint8(m.instruction & 0x0000000F)
		// Do the thing
		if spsr {
			cpu.WriteSPSR(cpu.ReadRegister(rm))
		} else {
			cpu.WriteRegister(16, cpu.ReadRegister(rm))
		}
	} else {
		// register contents or immediate value to PSR flag bits
		// Bit 25 is immediate flag
		immediate := m.instruction&(1<<25)>>25 == 1
		if immediate {
			val, _ := unshiftImmediate(m.instruction & 0x00000FFF)
			if spsr {
				cpu.WriteSPSR(val)
			} else {
				// only set the top 4 flag bits
				cpu.WriteCPSR(cpu.ReadCPSR()&0x0FFFFFFF | val&0xF0000000)
			}
		} else {
			val, _ := unshiftRegister(m.instruction&0x00000FFF, cpu)
			if spsr {
				if cpu.GetConfig().Debug {
					fmt.Printf("Register: r%d [%08X] to SPSR\n", m.instruction&0x00000F, val)
				}
			} else {
				if cpu.GetConfig().Debug {
					fmt.Printf("Register: r%d [%08X] to CPSR\n", m.instruction&0x00000F, val)
				}
			}
		}
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
