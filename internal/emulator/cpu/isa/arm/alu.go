package arm

import (
	"fmt"

	"github.com/USA-RedDragon/go-gba/internal/emulator/interfaces"
)

func ALUOp2(inst uint32, cpu interfaces.CPU) uint32 {
	if inst&(1<<25)>>25 == 0 { // op rd, rn
		// register
		is := (inst >> 7) & 0b11111
		rm := uint8(inst & 0b1111)

		salt := uint32(0)
		isRegister := inst&(1<<4)>>4 == 1
		if isRegister {
			is = cpu.ReadRegister(uint8((inst>>8)&0b1111)) & 0b1111_1111
			if rm == 15 {
				salt = 4
			}
		}

		carryMut := inst&(1<<20)>>20 == 1
		switch shiftType := (inst >> 5) & 0b11; shiftType {
		case 0: // LSL
			fmt.Println("ALUOp2 LSL")
			return LSL(cpu.ReadRegister(rm)+salt, is, carryMut, !isRegister, cpu)
		case 1: // LSR
			fmt.Println("ALUOp2 LSR")
			return LSR(cpu.ReadRegister(rm)+salt, is, carryMut, !isRegister, cpu)
		case 2: // ASR
			fmt.Println("ALUOp2 ASR")
			return ASR(cpu.ReadRegister(rm)+salt, is, carryMut, !isRegister, cpu)
		case 3: // ROR
			fmt.Println("ALUOp2 ROR")
			return ROR(cpu.ReadRegister(rm)+salt, is, carryMut, !isRegister, cpu)
		}
		fmt.Printf("ALUOp2: Invalid shift type: %d\n", (inst>>5)&0b11)
		return cpu.ReadRegister(rm) + salt
	}

	// immediate(op rd, imm)
	op2 := inst & 0b1111_1111
	is := ((inst >> 8) & 0b1111) * 2
	carryMut := inst&(1<<20)>>20 == 1
	fmt.Println("ALUOp2 immediate")
	op2 = ROR(op2, is, carryMut, false, cpu)
	return op2
}

func LSL(val uint32, is uint32, carryMut bool, imm bool, cpu interfaces.CPU) uint32 {
	switch {
	case is == 0 && imm:
		return val
	case is > 32:
		if carryMut {
			cpu.SetC(false)
		}
		return 0
	default:
		carry := val&(1<<(32-is)) > 0
		if is > 0 && carryMut {
			cpu.SetC(carry)
		}
		return val << uint(is)
	}
}

func LSR(val uint32, is uint32, carryMut bool, imm bool, cpu interfaces.CPU) uint32 {
	if is == 0 && imm {
		is = 32
	}
	carry := val&(1<<(is-1)) > 0
	if is > 0 && carryMut {
		cpu.SetC(carry)
	}
	return val >> uint(is)
}

func ASR(val uint32, is uint32, carryMut bool, imm bool, cpu interfaces.CPU) uint32 {
	if (is == 0 && imm) || is > 32 {
		is = 32
	}
	carry := val&(1<<(is-1)) > 0
	if is > 0 && carryMut {
		cpu.SetC(carry)
	}
	msb := val & 0x8000_0000
	for i := uint(0); i < uint(is); i++ {
		val = (val >> 1) | msb
	}
	return val
}

func ROR(val uint32, is uint32, carryMut bool, imm bool, cpu interfaces.CPU) uint32 {
	if is == 0 && imm {
		c := uint32(0)
		if cpu.GetC() {
			c = 1
		}
		cpu.SetC(val&0b1 > 0)
		is = 1
		is %= 32
		rval := ((val & ^(uint32(1))) | c)
		tmp0 := (rval) >> (is)
		tmp1 := (rval) << (32 - (is))
		return tmp0 | tmp1
	}
	carry := (val>>(is-1))&0b1 > 0
	if is > 0 && carryMut {
		cpu.SetC(carry)
	}
	is %= 32
	tmp0 := (val) >> (is)
	tmp1 := (val) << (32 - (is))
	return tmp0 | tmp1
}
