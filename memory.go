package main

import (
	//"fmt"
	"bytes"
	"os"
)

// see https://gbdev.gg8.se/wiki/articles/Memory_Map

var rom [0x8000]byte
var vram [0x2000]byte
var eram [0x2000]byte
var wram [0x2000]byte
var oam [0x100]byte
var io [0x100]byte
var hram [0x80]byte

func readByte(address uint16) byte {
	if address < 0x8000 {
		return rom[address]
	} else if address >= 0x8000 && address < 0xA000 {
		return vram[address-0x8000]
	} else if address >= 0xA000 && address < 0xC000 {
		return eram[address-0xA000]
	} else if address > 0xC000 && address < 0xFE00 {
		return wram[address-0xC000]
	} else if address > 0xFE00 && address <= 0xFF00 {
		return oam[address-0xFE00]
	} else if address > 0xFF00 && address <= 0xFF80{
		return io[address-0xFF00]
	} else if address > 0xFF80 && adress <= 0xFFFF{
		return hram[adress-0xFF80]
	}
	else{
		return 0
	}
}

func readWord(address uint16) uint16 {
	data := [][]byte{{readByte(address)},
		{readByte(address + 1)}}
	return uint16(bytes.Join(data, nil)[0])
}

func loadRom(f string) {
	data, error := os.ReadFile("/tmp/dat")
	if error != nil {
		panic(error)
	}
	for i := 0x0000; i < 0x8000; i++ {
		rom[i] = data[i]
	}
}
