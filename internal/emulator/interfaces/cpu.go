package interfaces

import "github.com/USA-RedDragon/go-gba/internal/emulator/memory"

type CPU interface {
	ReadRegister(reg uint8) uint32
	WriteRegister(reg uint8, value uint32)
	ReadSP() uint32
	WriteSP(value uint32)
	ReadLR() uint32
	WriteLR(value uint32)
	ReadPC() uint32
	WritePC(value uint32)
	ReadCPSR() uint32
	WriteCPSR(value uint32)
	ReadSPSR() uint32
	WriteSPSR(value uint32)

	GetMMIO() *memory.MMIO

	SetZ(value bool)
	SetN(value bool)
	SetC(value bool)
	SetV(value bool)
	GetZ() bool
	GetN() bool
	GetC() bool
	GetV() bool
	SetConditionCodes(res uint32, carry bool, overflow bool)

	SetThumbMode(value bool)
	GetThumbMode() bool
}
