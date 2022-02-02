package main

import (
	"os"
)

// see https://gbdev.gg8.se/wiki/articles/Memory_Map
type Memory struct {
	rom  [0x8000]byte
	vram [0x2000]byte
	eram [0x2000]byte
	wram [0x2000]byte
	oam  [0x100]byte
	io   [0x100]byte
	hram [0x80]byte
}

func (mem Memory) readByte(address uint16) byte {
	if address < 0x8000 {
		return mem.rom[address]
	} else if address >= 0x8000 && address < 0xA000 {
		return mem.vram[address-0x8000]
	} else if address >= 0xA000 && address < 0xC000 {
		return mem.eram[address-0xA000]
	} else if address >= 0xC000 && address < 0xE000 {
		return mem.wram[address-0xC000]
	} else if address >= 0xFE00 && address < 0xFF00 {
		return mem.oam[address-0xFE00]
	} else if address >= 0xFF00 && address < 0xFF80 {
		return mem.io[address-0xFF00]
	} else if address >= 0xFF80 && address < 0xFFFF {
		return mem.hram[address-0xFF80]
	} else {
		return 0
	}
}

func (mem Memory) readWord(address uint16) uint16 {
	data := concatenateBytes(mem.readByte(address), mem.readByte(address+1))
	return data
}

func (mem *Memory) writeByte(address uint16, value byte) {
	if address < 0x8000 {
		mem.rom[address] = value
	} else if address >= 0x8000 && address < 0xA000 {
		mem.vram[address-0x8000] = value
	} else if address >= 0xA000 && address < 0xC000 {
		mem.eram[address-0xA000] = value
	} else if address > 0xC000 && address < 0xFE00 {
		mem.wram[address-0xC000] = value
	} else if address > 0xFE00 && address <= 0xFF00 {
		mem.oam[address-0xFE00] = value
	} else if address > 0xFF00 && address <= 0xFF80 {
		mem.io[address-0xFF00] = value
	} else if address > 0xFF80 && address <= 0xFFFF {
		mem.hram[address-0xFF80] = value
	}
}

func (mem *Memory) writeWord(address uint16, value uint16) {
	r1, r2 := separateWord(value)
	mem.writeByte(address, r1)
	mem.writeByte(address+1, r2)
}

func (mem *Memory) loadRom(f string) {
	data, error := os.ReadFile(f)
	if error != nil {
		panic(error)
	}
	for i := 0x0000; i < 0x8000; i++ {
		mem.rom[i] = data[i]
	}
}
