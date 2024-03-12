package arm

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

type LDR struct {
	instruction uint32
}

func (ldr LDR) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	immediate := ldr.instruction&(1<<25)>>25 == 0
	pre := ldr.instruction&(1<<24)>>24 == 1
	up := ldr.instruction&(1<<23)>>23 == 1
	word := ldr.instruction&(1<<22)>>22 == 0
	writeback := ldr.instruction&(1<<21)>>21 == 1

	rn := uint8((ldr.instruction >> 16) & 0xF)
	rd := uint8((ldr.instruction >> 12) & 0xF)

	memory := cpu.GetMMIO()

	if cpu.GetConfig().Debug {
		fmt.Printf("Immediate: %t, Pre: %t, Up: %t, Word: %t, Writeback: %t\n", immediate, pre, up, word, writeback)
	}

	var offset uint32
	if immediate {
		offset = ldr.instruction & 0xFFF
	} else {
		offset, _ = unshiftRegister(ldr.instruction&0xFFF, cpu)
	}
	if cpu.GetConfig().Debug {
		fmt.Printf("LDR r%d, [r%d, 0x%X]\n", rd, rn, offset)
	}

	address := cpu.ReadRegister(rn)
	if pre {
		if up {
			address += offset
		} else {
			address -= offset
		}
	}

	var read uint32
	if !word {
		read8, err := memory.Read8(address)
		if err != nil {
			panic(err)
		}
		read = uint32(read8)
	} else {
		var err error
		read, err = memory.Read32(address)
		if err != nil {
			panic(err)
		}
	}
	cpu.WriteRegister(rd, read)

	if pre {
		if writeback && rn != rd {
			cpu.WriteRegister(rn, address)
		}
	} else {
		if up {
			address += offset
		} else {
			address -= offset
		}
		if rn != rd {
			cpu.WriteRegister(rn, address)
		}
	}

	if cpu.GetConfig().Debug {
		fmt.Printf("Address: 0x%X\n", address)
	}
	return
}

type STR struct {
	instruction uint32
}

func (str STR) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	immediate := str.instruction&(1<<25)>>25 == 0
	pre := str.instruction&(1<<24)>>24 == 1
	up := str.instruction&(1<<23)>>23 == 1
	word := str.instruction&(1<<22)>>22 == 0
	writeback := str.instruction&(1<<21)>>21 == 1

	rn := uint8((str.instruction >> 16) & 0xF)
	rd := uint8((str.instruction >> 12) & 0xF)
	rdVal := cpu.ReadRegister(rd)

	memory := cpu.GetMMIO()

	if cpu.GetConfig().Debug {
		fmt.Printf("Immediate: %t, Pre: %t, Up: %t, Word: %t, Writeback: %t\n", immediate, pre, up, word, writeback)
	}
	var offset uint32
	if immediate {
		offset = str.instruction & 0xFFF
	} else {
		offset, _ = unshiftRegister(str.instruction&0xFFF, cpu)
	}
	if cpu.GetConfig().Debug {
		fmt.Printf("STR r%d, [r%d, 0x%X]\n", rd, rn, offset)
	}
	if rd == 15 {
		rdVal += 4
	}
	address := cpu.ReadRegister(rn)
	if pre {
		if up {
			address += offset
		} else {
			address -= offset
		}
	}

	// If the datatype is a byte, check if it's on a word boundary
	if !word {
		// On a word boundary, bits 7-0 of the value in the register
		// are moved to the bottom bits of the memory at address
		err := memory.Write8(address, uint8(rdVal&0xFF))
		if err != nil {
			panic(err)
		}
	} else {
		// On a word boundary plus one, bits 15-8 of the value in the register
		// are moved to the bottom bits of the memory at address
		err := memory.Write32(address, rdVal)
		if err != nil {
			panic(err)
		}
	}
	if pre {
		if writeback {
			cpu.WriteRegister(rn, address)
		}
	} else {
		if up {
			address += offset
		} else {
			address -= offset
		}
		cpu.WriteRegister(rn, address)
	}
	if cpu.GetConfig().Debug {
		fmt.Printf("Address: 0x%X\n", address)
	}
	return
}

type LDM struct {
	instruction uint32
}

func (ldm LDM) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bit 24 == 1 means pre-indexed addressing
	pre := ldm.instruction&(1<<24)>>24 == 1
	// Bit 23 == 1 means the offset is added to the base register (up)
	up := ldm.instruction&(1<<23)>>23 == 1
	// Bit 22 == 1 means to load the PSR or force user mode
	psr := ldm.instruction&(1<<22)>>22 == 1
	// Bit 21 == 1 means the base register is written back to
	writeback := ldm.instruction&(1<<21)>>21 == 1

	if cpu.GetConfig().Debug {
		fmt.Printf("Pre: %t, Up: %t, PSR: %t, Writeback: %t\n", pre, up, psr, writeback)
	}

	// Bits 19-16 are the base register
	rn := uint8((ldm.instruction >> 16) & 0xF)

	// Bits 15-0 are the register list
	registerList := ldm.instruction & 0xFFFF

	address := cpu.ReadRegister(rn)

	// If the PSR bit is set, we need to load the PSR
	if psr {
		panic("Not implemented")
	}

	var registers []uint8
	if up {
		for i := 0; i < 16; i++ {
			if registerList&(1<<i)>>i == 1 {
				registers = append(registers, uint8(i))
			}
		}
	} else {
		for i := 15; i >= 0; i-- {
			if registerList&(1<<i)>>i == 1 {
				registers = append(registers, uint8(i))
			}
		}
	}

	for _, register := range registers {
		if pre {
			if up {
				address += 4
			} else {
				address -= 4
			}
		}
		value, err := cpu.GetMMIO().Read32(address)
		if err != nil {
			panic(err)
		}
		cpu.WriteRegister(register, value)
		if cpu.GetConfig().Debug {
			fmt.Printf("Pulling register r%d @ %08X\n", register, address)
		}
		if !pre {
			if up {
				address += 4
			} else {
				address -= 4
			}
		}
	}

	if (!pre || writeback) && registerList&(1<<rn)>>rn == 0 {
		cpu.WriteRegister(rn, address)
	}

	writebackStr := ""
	if writeback {
		writebackStr = "!"
	}

	if cpu.GetConfig().Debug {
		fmt.Printf("ldm r%d%s, {%v}\t # %08x\n", rn, writebackStr, registers, address)
	}
	return
}

type STM struct {
	instruction uint32
}

func (stm STM) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bit 24 == 1 means pre-indexed addressing
	p := stm.instruction&(1<<24)>>24 == 1
	// Bit 23 == 1 means the offset is added to the base register (up)
	u := stm.instruction&(1<<23)>>23 == 1
	rn := uint8((stm.instruction >> 16) & 0b1111)
	rnval := cpu.ReadRegister(rn)

	n := 0
	switch {
	case p && u: // IB
		for rs := 0; rs < 16; rs++ {
			if stm.instruction&(1<<rs)>>rs == 1 {
				cpu.WriteRegister(rn, cpu.ReadRegister(rn)+4)
				err := cpu.GetMMIO().Write32(cpu.ReadRegister(rn), cpu.ReadRegister(uint8(rs)))
				if err != nil {
					panic(err)
				}
				n++
			}
		}
	case !p && u: // IA
		for rs := 0; rs < 16; rs++ {
			if stm.instruction&(1<<rs)>>rs == 1 {
				err := cpu.GetMMIO().Write32(cpu.ReadRegister(rn), cpu.ReadRegister(uint8(rs)))
				if err != nil {
					panic(err)
				}
				cpu.WriteRegister(rn, cpu.ReadRegister(rn)+4)
				n++
			}
		}
	case p && !u: // DB, push
		for rs := 15; rs >= 0; rs-- {
			if stm.instruction&(1<<rs)>>rs == 1 {
				cpu.WriteRegister(rn, cpu.ReadRegister(rn)-4)
				err := cpu.GetMMIO().Write32(cpu.ReadRegister(rn), cpu.ReadRegister(uint8(rs)))
				if err != nil {
					panic(err)
				}
				n++
			}
		}
	case !p && !u: // DA
		for rs := 15; rs >= 0; rs-- {
			if stm.instruction&(1<<rs)>>rs == 1 {
				err := cpu.GetMMIO().Write32(cpu.ReadRegister(rn), cpu.ReadRegister(uint8(rs)))
				if err != nil {
					panic(err)
				}
				cpu.WriteRegister(rn, cpu.ReadRegister(rn)-4)
				n++
			}
		}
	}

	// Pre-indexing, write-back is optional
	writeBack := stm.instruction&(1<<21)>>21 == 1
	if p && !writeBack {
		cpu.WriteRegister(rn, rnval)
	}

	writebackStr := ""
	if writeBack {
		writebackStr = "!"
	}

	if cpu.GetConfig().Debug {
		fmt.Printf("stm r%d%s\n", rn, writebackStr)
	}

	return
}

type LDRSH struct {
	instruction uint32
}

func (ldrsh LDRSH) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bit 24 == 1 means pre-indexed addressing
	pre := ldrsh.instruction&(1<<24)>>24 == 1
	// // Bit 23 == 1 means the offset is added to the base register (up)
	up := ldrsh.instruction&(1<<23)>>23 == 1
	// // Bit 21 == 1 means the base register is written back to
	writeback := ldrsh.instruction&(1<<21)>>21 == 1

	// Bits 19-16 are the base register
	rn := uint8((ldrsh.instruction >> 16) & 0xF)

	// Bits 15-12 are the destination register
	rd := uint8((ldrsh.instruction >> 12) & 0xF)

	// Bits 11-8 are the offset's high nibble
	offsetHigh := uint8((ldrsh.instruction >> 8) & 0xF)

	// Bits 3-0 are the offset's low nibble
	offsetLow := uint8(ldrsh.instruction & 0xF)

	address := cpu.ReadRegister(rn)
	if pre {
		if up {
			address += uint32(offsetHigh<<4 | offsetLow)
		} else {
			address -= uint32(offsetHigh<<4 | offsetLow)
		}
	}

	fmt.Printf("ldrsh r%d, [r%d, #%d]\n", rd, rn, offsetHigh<<4|offsetLow)

	halfword, err := cpu.GetMMIO().Read16(address & 0xFFFFFFFE)
	if err != nil {
		panic(err)
	}
	signedHalfword := int32(int16(halfword))

	if address&1 == 1 {
		val := signedHalfword
		// Right rotate the halfword by 8 bits
		is := 8
		is %= 32
		tmp0 := (val) >> (is)
		tmp1 := (val) << (32 - (is))
		cpu.WriteRegister(rd, uint32(tmp0|tmp1))
	} else {
		cpu.WriteRegister(rd, uint32(signedHalfword))
	}

	if pre {
		if writeback && rn != rd {
			cpu.WriteRegister(rn, address)
		}
	} else {
		if up {
			address += uint32(offsetHigh<<4 | offsetLow)
		} else {
			address -= uint32(offsetHigh<<4 | offsetLow)
		}
		if rn != rd {
			cpu.WriteRegister(rn, address)
		}
	}

	return
}

type LDRSB struct {
	instruction uint32
}

func (ldrsb LDRSB) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bit 24 == 1 means pre-indexed addressing
	pre := ldrsb.instruction&(1<<24)>>24 == 1
	// Bit 23 == 1 means the offset is added to the base register (up)
	up := ldrsb.instruction&(1<<23)>>23 == 1
	// Bit 21 == 1 means the base register is written back to
	writeback := ldrsb.instruction&(1<<21)>>21 == 1

	// Bits 19-16 are the base register
	rn := uint8((ldrsb.instruction >> 16) & 0xF)

	// Bits 15-12 are the destination register
	rd := uint8((ldrsb.instruction >> 12) & 0xF)

	// Bits 11-8 are the offset's high nibble
	offsetHigh := uint8((ldrsb.instruction >> 8) & 0xF)

	// Bits 3-0 are the offset's low nibble
	offsetLow := uint8(ldrsb.instruction & 0xF)

	address := cpu.ReadRegister(rn)
	if pre {
		if up {
			address += uint32(offsetHigh<<4 | offsetLow)
		} else {
			address -= uint32(offsetHigh<<4 | offsetLow)
		}
	}

	fmt.Printf("ldrsb r%d, [r%d, #%d]\n", rd, rn, offsetHigh<<4|offsetLow)

	b, err := cpu.GetMMIO().Read8(address)
	if err != nil {
		panic(err)
	}
	signedByte := int32(int8(b))

	cpu.WriteRegister(rd, uint32(signedByte))

	if pre {
		if writeback && rn != rd {
			cpu.WriteRegister(rn, address)
		}
	} else {
		if up {
			address += uint32(offsetHigh<<4 | offsetLow)
		} else {
			address -= uint32(offsetHigh<<4 | offsetLow)
		}
		if rn != rd {
			cpu.WriteRegister(rn, address)
		}
	}

	if pre {
		if writeback && rn != rd {
			cpu.WriteRegister(rn, address)
		}
	} else {
		if up {
			address += uint32(offsetHigh<<4 | offsetLow)
		} else {
			address -= uint32(offsetHigh<<4 | offsetLow)
		}
		if rn != rd {
			cpu.WriteRegister(rn, address)
		}
	}

	return
}

type LDRH struct {
	instruction uint32
}

func (ldrh LDRH) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bit 24 == 1 means pre-indexed addressing
	pre := ldrh.instruction&(1<<24)>>24 == 1
	// Bit 23 == 1 means the offset is added to the base register (up)
	up := ldrh.instruction&(1<<23)>>23 == 1
	// Bit 21 == 1 means the base register is written back to
	writeback := ldrh.instruction&(1<<21)>>21 == 1

	// Bits 19-16 are the base register
	rn := uint8((ldrh.instruction >> 16) & 0xF)

	// Bits 15-12 are the destination register
	rd := uint8((ldrh.instruction >> 12) & 0xF)

	// Bits 11-8 are the offset's high nibble
	offsetHigh := uint8((ldrh.instruction >> 8) & 0xF)

	// Bits 3-0 are the offset's low nibble
	offsetLow := uint8(ldrh.instruction & 0xF)

	address := cpu.ReadRegister(rn)
	if pre {
		if up {
			address += uint32(offsetHigh<<4 | offsetLow)
		} else {
			address -= uint32(offsetHigh<<4 | offsetLow)
		}
	}

	if cpu.GetConfig().Debug {
		fmt.Printf("ldrh r%d, [r%d, #%d]\n", rd, rn, offsetHigh<<4|offsetLow)
	}

	// Load halfword from memory
	halfword, err := cpu.GetMMIO().Read16(address & 0xFFFFFFFE)
	if err != nil {
		panic(err)
	}

	if address&1 == 1 {
		val := uint32(halfword)
		// Right rotate the halfword by 8 bits
		is := 8
		is %= 32
		tmp0 := (val) >> (is)
		tmp1 := (val) << (32 - (is))
		cpu.WriteRegister(rd, tmp0|tmp1)
	} else {
		cpu.WriteRegister(rd, uint32(halfword))
	}

	if pre {
		if writeback && rn != rd {
			cpu.WriteRegister(rn, address)
		}
	} else {
		if up {
			address += uint32(offsetHigh<<4 | offsetLow)
		} else {
			address -= uint32(offsetHigh<<4 | offsetLow)
		}
		if rn != rd {
			cpu.WriteRegister(rn, address)
		}
	}

	return
}

type STRSH struct {
	instruction uint32
}

func (strsh STRSH) Execute(_ interfaces.CPU) (repipeline bool, cycles uint16) {
	// // Bit 24 == 1 means pre-indexed addressing
	// pre := strsh.instruction&(1<<24)>>24 == 1
	// // Bit 23 == 1 means the offset is added to the base register (up)
	// up := strsh.instruction&(1<<23)>>23 == 1
	// // Bit 21 == 1 means the base register is written back to
	// writeback := strsh.instruction&(1<<21)>>21 == 1

	// // Bits 19-16 are the base register
	// rn := uint8((strsh.instruction >> 16) & 0xF)

	// // Bits 15-12 are the destination register
	// rd := uint8((strsh.instruction >> 12) & 0xF)

	// // Bits 11-8 are the offset's high nibble
	// offsetHigh := uint8((strsh.instruction >> 8) & 0xF)

	// // Bits 3-0 are the offset's low nibble
	// offsetLow := uint8(strsh.instruction & 0xF)
	panic("STRSH Not implemented")
}

type STRSB struct {
	instruction uint32
}

func (strsb STRSB) Execute(_ interfaces.CPU) (repipeline bool, cycles uint16) {
	// // Bit 24 == 1 means pre-indexed addressing
	// pre := strsb.instruction&(1<<24)>>24 == 1
	// // Bit 23 == 1 means the offset is added to the base register (up)
	// up := strsb.instruction&(1<<23)>>23 == 1
	// // Bit 21 == 1 means the base register is written back to
	// writeback := strsb.instruction&(1<<21)>>21 == 1

	// // Bits 19-16 are the base register
	// rn := uint8((strsb.instruction >> 16) & 0xF)

	// // Bits 15-12 are the destination register
	// rd := uint8((strsb.instruction >> 12) & 0xF)

	// // Bits 11-8 are the offset's high nibble
	// offsetHigh := uint8((strsb.instruction >> 8) & 0xF)

	// // Bits 3-0 are the offset's low nibble
	// offsetLow := uint8(strsb.instruction & 0xF)
	panic("STRSB Not implemented")
}

type STRH struct {
	instruction uint32
}

func (strh STRH) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bit 24 == 1 means pre-indexed addressing
	pre := strh.instruction&(1<<24)>>24 == 1
	// Bit 23 == 1 means the offset is added to the base register (up)
	up := strh.instruction&(1<<23)>>23 == 1
	// Bit 21 == 1 means the base register is written back to
	writeback := strh.instruction&(1<<21)>>21 == 1

	// Bits 19-16 are the base register
	rn := uint8((strh.instruction >> 16) & 0xF)

	// Bits 15-12 are the destination register
	rd := uint8((strh.instruction >> 12) & 0xF)

	// Bits 11-8 are the offset's high nibble
	offsetHigh := uint8((strh.instruction >> 8) & 0xF)

	// Bits 3-0 are the offset's low nibble
	offsetLow := uint8(strh.instruction & 0xF)

	offset := uint32(offsetHigh)<<4 | uint32(offsetLow)

	address := cpu.ReadRegister(rn)
	if pre {
		if up {
			address += offset
		} else {
			address -= offset
		}
		if writeback {
			cpu.WriteRegister(rn, address)
		}
	}

	if cpu.GetConfig().Debug {
		fmt.Printf("strh r%d, [r%d, #0x%X]  # 0x%08x\n", rd, rn, offset, address)
	}

	// Store unsigned halfword
	err := cpu.GetMMIO().Write16(address, uint16(cpu.ReadRegister(rd)))
	if err != nil {
		panic(err)
	}

	if !pre {
		if up {
			address += offset
		} else {
			address -= offset
		}
		cpu.WriteRegister(rn, address)
	}

	return
}

type LDRSHRegisterOffset struct {
	instruction uint32
}

func (ldrsh LDRSHRegisterOffset) Execute(_ interfaces.CPU) (repipeline bool, cycles uint16) {
	panic("LDRSHRegisterOffset Not implemented")
}

type LDRSBRegisterOffset struct {
	instruction uint32
}

func (ldrsb LDRSBRegisterOffset) Execute(_ interfaces.CPU) (repipeline bool, cycles uint16) {
	panic("LDRSBRegisterOffset Not implemented")
}

type LDRHRegisterOffset struct {
	instruction uint32
}

func (ldrh LDRHRegisterOffset) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bit 24 == 1 means pre-indexed addressing
	pre := ldrh.instruction&(1<<24)>>24 == 1
	// Bit 23 == 1 means the offset is added to the base register (up)
	up := ldrh.instruction&(1<<23)>>23 == 1
	// Bit 21 == 1 means the base register is written back to
	writeback := ldrh.instruction&(1<<21)>>21 == 1

	// Bits 19-16 are the base register
	rn := uint8((ldrh.instruction >> 16) & 0xF)

	// Bits 15-12 are the destination register
	rd := uint8((ldrh.instruction >> 12) & 0xF)

	// Bits 3-0 are the offset register
	rm := uint8(ldrh.instruction & 0xF)

	offset := cpu.ReadRegister(rm)

	address := cpu.ReadRegister(rn)
	if pre {
		if up {
			address += offset
		} else {
			address -= offset
		}
	}

	fmt.Printf("ldrh r%d, [r%d, r%d]  # 0x%08x\n", rd, rn, rm, address)

	// Load halfword from memory
	halfword, err := cpu.GetMMIO().Read16(address)
	if err != nil {
		panic(err)
	}

	cpu.WriteRegister(rd, uint32(halfword))

	if pre {
		if writeback && rn != rd {
			cpu.WriteRegister(rn, address)
		}
	} else {
		if up {
			address += offset
		} else {
			address -= offset
		}
		if rn != rd {
			cpu.WriteRegister(rn, address)
		}
	}

	return
}

type STRSHRegisterOffset struct {
	instruction uint32
}

func (strsh STRSHRegisterOffset) Execute(_ interfaces.CPU) (repipeline bool, cycles uint16) {
	panic("STRSHRegisterOffset Not implemented")
}

type STRSBRegisterOffset struct {
	instruction uint32
}

func (strsb STRSBRegisterOffset) Execute(_ interfaces.CPU) (repipeline bool, cycles uint16) {
	panic("STRSBRegisterOffset Not implemented")
}

type STRHRegisterOffset struct {
	instruction uint32
}

func (strh STRHRegisterOffset) Execute(cpu interfaces.CPU) (repipeline bool, cycles uint16) {
	// Bit 24 == 1 means pre-indexed addressing
	pre := strh.instruction&(1<<24)>>24 == 1
	// Bit 23 == 1 means the offset is added to the base register (up)
	up := strh.instruction&(1<<23)>>23 == 1
	// Bit 21 == 1 means the base register is written back to
	writeback := strh.instruction&(1<<21)>>21 == 1

	// Bits 19-16 are the base register
	rn := uint8((strh.instruction >> 16) & 0xF)

	// Bits 15-12 are the destination register
	rd := uint8((strh.instruction >> 12) & 0xF)

	// Bits 3-0 are the offset register
	rm := uint8(strh.instruction & 0xF)

	offset := cpu.ReadRegister(rm)

	address := cpu.ReadRegister(rn)
	if pre {
		if up {
			address += offset
		} else {
			address -= offset
		}
		if writeback {
			cpu.WriteRegister(rn, address)
		}
	}

	if cpu.GetConfig().Debug {
		fmt.Printf("strh r%d, [r%d, r%d]  # 0x%08x\n", rd, rn, rm, address)
	}

	// Store unsigned halfword
	err := cpu.GetMMIO().Write16(address, uint16(cpu.ReadRegister(rd)))
	if err != nil {
		panic(err)
	}

	if !pre {
		if up {
			address += offset
		} else {
			address -= offset
		}
		cpu.WriteRegister(rn, address)
	}

	return
}
