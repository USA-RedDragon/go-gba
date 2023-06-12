package arm

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type B struct {
	instruction uint32
}

func (b B) Execute(cpu interfaces.CPU) {
	offset := b.instruction & 0x00FFFFFF
	// Sign extend the offset
	if offset&0x00800000 != 0 {
		offset |= 0xFF000000
	}
	offset <<= 2
	// if bit 0 of the offset is set, we're in THUMB mode
	if offset&0x00000001 == 1 {
		fmt.Println("Setting THUMB mode")
		cpu.WritePC(cpu.ReadPC() + offset - 1)
		cpu.SetThumbMode(true)
	} else {
		cpu.WritePC(cpu.ReadPC() + offset)
		cpu.SetThumbMode(false)
	}
	fmt.Printf("New PC 0x%X\n", cpu.ReadPC())
}

type BL struct {
	instruction uint32
}

func (bl BL) Execute(cpu interfaces.CPU) {
	fmt.Println("Branch With Link")
	offset := bl.instruction & 0x00FFFFFF
	// Sign extend the offset
	if offset&0x00800000 != 0 {
		offset |= 0xFF000000
	}
	offset <<= 2
	cpu.WriteRegister(14, cpu.ReadPC())
	// if bit 0 of the offset is set, we're in THUMB mode
	if offset&0x00000001 == 1 {
		fmt.Println("Setting THUMB mode")
		cpu.WritePC(cpu.ReadPC() + offset - 1)
		cpu.SetThumbMode(true)
	} else {
		cpu.WritePC(cpu.ReadPC() + offset)
		cpu.SetThumbMode(false)
	}
	fmt.Printf("Branching by 0x%X\n", offset)
	fmt.Printf("New PC 0x%X\n", cpu.ReadPC())
}

type BX struct {
	instruction uint32
}

func (bx BX) Execute(cpu interfaces.CPU) {
	fmt.Println("Branch Exchange")

	// Bits 3-0 are the register to branch to
	rm := uint8(bx.instruction & 0x0000000F)

	// if bit 0 of the register is set, we're in THUMB mode
	if cpu.ReadRegister(rm)&0x00000001 == 1 {
		fmt.Println("Setting THUMB mode")
		cpu.WritePC(cpu.ReadRegister(rm) - 1)
		cpu.SetThumbMode(true)
	} else {
		cpu.WritePC(cpu.ReadRegister(rm))
		cpu.SetThumbMode(false)
	}

	fmt.Printf("New PC 0x%X\n", cpu.ReadPC())
}
