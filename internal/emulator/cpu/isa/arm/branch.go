package arm

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type B struct {
	instruction uint32
}

func (b B) Execute(cpu interfaces.CPU) (repipeline bool) {
	offset := b.instruction & 0x00FFFFFF
	// Sign extend the offset
	if offset&0x00800000 != 0 {
		offset |= 0xFF000000
	}
	offset <<= 2
	// if bit 0 of the offset is set, we're in THUMB mode
	cpu.SetThumbMode(offset&0b11 != 0)
	cpu.WritePC(cpu.ReadPC() + offset&0xFFFFFFFC)
	if cpu.GetConfig().Debug {
		fmt.Printf("New PC 0x%X\n", cpu.ReadPC())
	}
	return
}

type BL struct {
	instruction uint32
}

func (bl BL) Execute(cpu interfaces.CPU) (repipeline bool) {
	offset := bl.instruction & 0x00FFFFFF
	// Sign extend the offset
	if offset&0x00800000 != 0 {
		offset |= 0xFF000000
	}
	offset <<= 2
	cpu.WriteLR(cpu.ReadPC() - 4)
	// if bit 0 of the offset is set, we're in THUMB mode
	if offset&0b11 != 0 {
		if cpu.GetConfig().Debug {
			fmt.Println("Setting THUMB mode")
		}
		cpu.WritePC(cpu.ReadPC() + offset - 1)
		cpu.SetThumbMode(true)
	} else {
		cpu.WritePC(cpu.ReadPC() + offset)
		cpu.SetThumbMode(false)
	}
	if cpu.GetConfig().Debug {
		fmt.Printf("Branching by 0x%X\n", offset)
		fmt.Printf("New PC 0x%X\n", cpu.ReadPC())
	}
	return
}

type BX struct {
	instruction uint32
}

func (bx BX) Execute(cpu interfaces.CPU) (repipeline bool) {
	// Bits 3-0 are the register to branch to
	rm := uint8(bx.instruction & 0x0000000F)

	// if bit 0 of the register is set, we're in THUMB mode
	if cpu.ReadRegister(rm)&0b11 != 0 {
		if cpu.GetConfig().Debug {
			fmt.Println("Setting THUMB mode")
		}
		cpu.SetThumbMode(true)
		cpu.WritePC(cpu.ReadRegister(rm) - 1)
	} else {
		cpu.SetThumbMode(false)
		cpu.WritePC(cpu.ReadRegister(rm))
	}

	if cpu.GetConfig().Debug {
		fmt.Printf("New PC 0x%X\n", cpu.ReadPC())
	}
	return
}
