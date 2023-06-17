package memory

import "fmt"

type mmioMapping struct {
	address uint32
	size    uint32

	data []byte
}

type MMIO struct {
	mmios []mmioMapping
}

func (h *MMIO) checkWritable(addr uint32) bool {
	// Addresses 0x00000000 - 0x00003FFF are not writable (BIOS)
	// Addresses 0x02000000 - 0x0203FFFF are writable (on-board WRAM)
	// Addresses 0x03000000 - 0x03007FFF are writable (on-chip WRAM)
	// Addresses 0x04000000 - 0x040003FE are writable (I/O registers)
	// Addresses 0x05000000 - 0x050003FF are writable (palette RAM)
	// Addresses 0x06000000 - 0x06017FFF are writable (VRAM)
	// Addresses 0x07000000 - 0x070003FF are writable (OAM)
	// Addresses 0x08000000 - 0x09FFFFFF are not writable (Game Pak ROM/FlashROM - waitstate 0)
	// Addresses 0x0A000000 - 0x0BFFFFFF are not writable (Game Pak ROM/FlashROM - waitstate 1)
	// Addresses 0x0C000000 - 0x0DFFFFFF are not writable (Game Pak ROM/FlashROM - waitstate 2)
	// Addresses 0x0E000000 - 0x0E00FFFF are writable (Game Pak SRAM)
	writable := true
	if addr >= 0x00000000 && addr < 0x00004000 {
		writable = false
	} else if addr >= 0x02000000 && addr < 0x02040000 {
		writable = true
	} else if addr >= 0x03000000 && addr < 0x03008000 {
		writable = true
	} else if addr >= 0x04000000 && addr < 0x04000400 {
		writable = true
	} else if addr >= 0x05000000 && addr < 0x05000400 {
		writable = true
	} else if addr >= 0x06000000 && addr < 0x06018000 {
		writable = true
	} else if addr >= 0x07000000 && addr < 0x07000400 {
		writable = true
	} else if addr >= 0x08000000 && addr < 0x0A000000 {
		writable = false
	} else if addr >= 0x0A000000 && addr < 0x0C000000 {
		writable = false
	} else if addr >= 0x0C000000 && addr < 0x0E000000 {
		writable = false
	} else if addr >= 0x0E000000 && addr < 0x0E010000 {
		writable = true
	}
	return writable
}

// findMMIO finds the MMIO device index that contains the given address
func (h *MMIO) findMMIOIndex(addr uint32) (int, error) {
	// account for memory mirroring
	// 0x03007FFF - 0x03FFFFFF should map repeatedly to 0x03000000 - 0x03007FFF
	if addr >= 0x03007FFF && addr < 0x04000000 {
		mod := addr % 0x8000
		addr = 0x03000000 + mod
	}
	for i, mmio := range h.mmios {
		if addr >= mmio.address && addr < mmio.address+mmio.size {
			return i, nil
		}
	}
	return 0, fmt.Errorf("MMIO address %08x not found", addr)
}

// Read8 reads a 8-bit value from the MMIO address space and returns it.
func (h *MMIO) Read8(addr uint32) (uint8, error) {
	index, err := h.findMMIOIndex(addr)
	if err != nil {
		return 0, err
	}
	nonMapped := addr - h.mmios[index].address
	if nonMapped >= h.mmios[index].size {
		return 0, fmt.Errorf("MMIO address %08x not found", addr)
	}
	return h.mmios[index].data[nonMapped], nil
}

// Write8 writes a 8-bit value to the MMIO address space.
func (h *MMIO) Write8(addr uint32, data uint8) error {
	index, err := h.findMMIOIndex(addr)
	if err != nil {
		return err
	}
	if !h.checkWritable(addr) {
		panic("MMIO address not writable")
		return fmt.Errorf("MMIO address %08x not writable", addr)
	}
	nonMapped := addr - h.mmios[index].address
	if nonMapped >= h.mmios[index].size {
		return fmt.Errorf("MMIO address %08x not found", addr)
	}
	h.mmios[index].data[nonMapped] = data
	return nil
}

// Read16 reads a 16-bit value from the MMIO address space and returns it.
func (h *MMIO) Read16(addr uint32) (uint16, error) {
	index, err := h.findMMIOIndex(addr)
	if err != nil {
		return 0, err
	}
	nonMapped := addr - h.mmios[index].address
	if nonMapped >= h.mmios[index].size {
		return 0, fmt.Errorf("MMIO address %08x not found", addr)
	}
	dataBytes := h.mmios[index].data[nonMapped : nonMapped+2]
	return uint16(dataBytes[0]) | uint16(dataBytes[1])<<8, nil
}

// Write16 writes a 16-bit value to the MMIO address space.
func (h *MMIO) Write16(addr uint32, data uint16) error {
	index, err := h.findMMIOIndex(addr)
	if err != nil {
		return err
	}
	if !h.checkWritable(addr) {
		panic("MMIO address not writable")
		return fmt.Errorf("MMIO address %08x not writable", addr)
	}
	nonMapped := addr - h.mmios[index].address
	if nonMapped >= h.mmios[index].size {
		return fmt.Errorf("MMIO address %08x not found", addr)
	}
	h.mmios[index].data[nonMapped] = byte(data)
	h.mmios[index].data[nonMapped+1] = byte(data >> 8)
	return nil
}

// Read32 reads a 32-bit value from the MMIO address space and returns it.
func (h *MMIO) Read32(addr uint32) (uint32, error) {
	index, err := h.findMMIOIndex(addr)
	if err != nil {
		return 0, err
	}
	if addr >= 0x03007FFF && addr < 0x04000000 {
		mod := addr % 0x8000
		addr = 0x03000000 + mod
	}
	nonMapped := addr - h.mmios[index].address
	if nonMapped >= h.mmios[index].size {
		return 0, fmt.Errorf("MMIO address %08x not found", addr)
	}
	dataBytes := h.mmios[index].data[nonMapped : nonMapped+4]
	return uint32(dataBytes[0]) | uint32(dataBytes[1])<<8 | uint32(dataBytes[2])<<16 | uint32(dataBytes[3])<<24, nil
}

// Write32 writes a 32-bit value to the MMIO address space.
func (h *MMIO) Write32(addr uint32, data uint32) error {
	index, err := h.findMMIOIndex(addr)
	if err != nil {
		return err
	}
	if addr >= 0x03007FFF && addr < 0x04000000 {
		mod := addr % 0x8000
		addr = 0x03000000 + mod
	}
	if !h.checkWritable(addr) {
		panic("MMIO address not writable")
		return fmt.Errorf("MMIO address %08x not writable", addr)
	}
	nonMapped := addr - h.mmios[index].address
	if nonMapped >= h.mmios[index].size {
		return fmt.Errorf("MMIO address %08x not found", addr)
	}
	h.mmios[index].data[nonMapped] = byte(data)
	h.mmios[index].data[nonMapped+1] = byte(data >> 8)
	h.mmios[index].data[nonMapped+2] = byte(data >> 16)
	h.mmios[index].data[nonMapped+3] = byte(data >> 24)
	return nil
}

// AddMMIO adds a new MMIO device to the MMIO handler.
func (h *MMIO) AddMMIO(data []byte, address uint32, size uint32) {
	h.mmios = append(h.mmios, mmioMapping{address, size, data})
}
