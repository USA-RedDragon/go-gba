package arm

import (
	"fmt"
	"math"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type SWI struct {
	instruction uint32
}

func (s SWI) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bits 23-0 are the comment field
	comment := s.instruction & 0x00FFFFFF
	// DIV BIOS call
	if comment == 0x60000 {
		numerator := cpu.ReadRegister(0)
		denominator := cpu.ReadRegister(1)
		cpu.WriteRegister(0, numerator/denominator)
		cpu.WriteRegister(1, numerator%denominator)
		cpu.WriteRegister(3, uint32(math.Abs(float64(numerator)/float64(denominator))))
	} else {
		panic(fmt.Sprintf("SWI: %d", comment))
	}
	return
}
