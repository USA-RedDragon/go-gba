package thumb

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type B struct {
	instruction uint16
}

func (a B) Execute(cpu interfaces.CPU) {
	fmt.Println("B")
	fmt.Println("ConditionalBranch")
	// Bits 11-8 are the condition
	cond := a.instruction & (1<<11 | 1<<10 | 1<<9 | 1<<8) >> 8

	// Bits 7-0 are the 8-bit signed offset
	offset := int8(a.instruction & 0xFF)

	conditionPassed := false
	switch cond {
	case 0b0:
		conditionPassed = cpu.GetZ()
	case 0b1:
		conditionPassed = !cpu.GetZ()
	case 0b10:
		conditionPassed = cpu.GetC()
	case 0b11:
		conditionPassed = !cpu.GetC()
	case 0b100:
		conditionPassed = cpu.GetN()
	case 0b101:
		conditionPassed = !cpu.GetN()
	case 0b110:
		conditionPassed = cpu.GetV()
	case 0b111:
		conditionPassed = !cpu.GetV()
	case 0b1000:
		conditionPassed = cpu.GetC() && !cpu.GetZ()
	case 0b1001:
		conditionPassed = !cpu.GetC() || cpu.GetZ()
	case 0b1010:
		conditionPassed = cpu.GetN() == cpu.GetV()
	case 0b1011:
		conditionPassed = cpu.GetN() != cpu.GetV()
	case 0b1100:
		conditionPassed = !cpu.GetZ() && cpu.GetN() == cpu.GetV()
	case 0b1101:
		conditionPassed = cpu.GetZ() || cpu.GetN() != cpu.GetV()
	}

	if conditionPassed {
		cpu.WritePC(cpu.ReadPC() + uint32(offset))
	} else {
		fmt.Println("Branch condition not met")
	}
}

type BX struct {
	instruction uint16
}

func (a BX) Execute(cpu interfaces.CPU) {
	fmt.Println("BX")

	// Bits 7-6 are the hi operand flags
	hof := a.instruction & (1<<7 | 1<<6) >> 6
	rshs := uint8(a.instruction & (1<<5 | 1<<4 | 1<<3) >> 3)

	if hof == 0b00 {
		fmt.Printf("Branching to r%d\n", rshs)
		addr := cpu.ReadRegister(rshs)
		fmt.Printf("Address: 0x%08X\n", addr)
		if addr&1 == 1 {
			fmt.Println("Setting THUMB mode")
			cpu.WritePC(addr - 1)
			cpu.SetThumbMode(true)
		} else {
			cpu.WritePC(addr)
			cpu.SetThumbMode(false)
		}
	} else if hof == 0b01 {
		fmt.Printf("Branching to r%d\n", rshs+8)
		addr := cpu.ReadRegister(rshs + 8)
		fmt.Printf("Address: 0x%08X\n", addr)
		if addr&1 == 1 {
			fmt.Println("Setting THUMB mode")
			cpu.WritePC(addr - 1)
			cpu.SetThumbMode(true)
		} else {
			cpu.WritePC(addr)
			cpu.SetThumbMode(false)
		}
	}
}
