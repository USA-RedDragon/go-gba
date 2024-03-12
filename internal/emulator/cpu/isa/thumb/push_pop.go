package thumb

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type PUSH struct {
	instruction uint16
}

func (p PUSH) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("PUSH")

	// Bit 8 denotes storing LR
	storeLR := p.instruction&(1<<8)>>8 == 1

	// Bits 7-0 are the registers to push
	registers := p.instruction & 0xFF

	var pushRegisters []uint8
	// Push LR if needed
	if storeLR {
		pushRegisters = append(pushRegisters, 14)
	}

	// Collect the registers to push in backwards order so that they are pushed in the correct order
	for i := 7; i >= 0; i-- {
		if registers&(1<<i)>>i == 1 {
			pushRegisters = append(pushRegisters, uint8(i))
		}
	}

	// Push the registers
	for _, reg := range pushRegisters {
		cpu.WriteSP(cpu.ReadSP() - 4)
		err := cpu.GetMMIO().Write32(cpu.ReadSP(), cpu.ReadRegister(reg))
		if err != nil {
			panic(err)
		}
		fmt.Printf("Pushing register r%d @ %08X\n", reg, cpu.ReadSP())
	}
	return
}

type POP struct {
	instruction uint16
}

func (p POP) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	fmt.Println("POP")

	// Bit 8 denotes loading PC
	loadPC := p.instruction&(1<<8)>>8 == 1

	// Bits 7-0 are the registers to push
	registers := p.instruction & 0xFF

	var popRegisters []uint8

	// Collect the registers to pop
	for i := 0; i < 8; i++ {
		if registers&(1<<i)>>i == 1 {
			popRegisters = append(popRegisters, uint8(i))
		}
	}

	if loadPC {
		popRegisters = append(popRegisters, 15)
	}

	// Pop the registers
	for _, reg := range popRegisters {
		contents, err := cpu.GetMMIO().Read32(cpu.ReadSP())
		if err != nil {
			panic(err)
		}
		cpu.WriteRegister(reg, contents)
		fmt.Printf("Popping register r%d @ %08X\n", reg, cpu.ReadSP())
		cpu.WriteSP(cpu.ReadSP() + 4)
	}
	return
}
