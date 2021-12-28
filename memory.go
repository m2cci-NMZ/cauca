package main

import (
	//"fmt"
	"bytes"
	"os"
)

/*
From http://marc.rawer.de/Gameboy/Docs/GBCPUman.pdf


 Interrupt Enable Register
 --------------------------- FFFF
 Internal RAM
 --------------------------- FF80
 Empty but unusable for I/O
 --------------------------- FF4C
 I/O ports
 --------------------------- FF00
 Empty but unusable for I/O
 --------------------------- FEA0
 Sprite Attrib Memory (OAM)
 --------------------------- FE00
 Echo of 8kB Internal RAM
 --------------------------- E000
 8kB Internal RAM
 --------------------------- C000
 8kB switchable RAM bank
 --------------------------- A000
 8kB Video RAM
 --------------------------- 8000 --
 16kB switchable ROM bank |
 --------------------------- 4000 |= 32kB Cartrigbe
 16kB ROM bank #0 |
 --------------------------- 0000 --
*/

var rom [0x8000]byte
var vram [0x2000]byte
var sram [0x2000]byte
var iram1 [0x2000]byte
var eiram [0x2000]byte
var oam [0x100]byte
var io [0x100]byte
var iram2 [0x80]byte

func readByte(address uint16) byte {
	switch address & 0xF000 {
	// rom bank 0
	case 0x0000:
	case 0x1000:
	case 0x2000:
	case 0x3000:
	// switchable rom bank
	case 0x4000:
	case 0x5000:
	case 0x6000:
	case 0x7000:
		return rom[address]
	// video ram
	case 0x8000:
	case 0x9000:
		return vram[address&0x1FFF]
	// switchable ram
	case 0xA000:
	case 0xB000:
		return sram[address&0x1FFF]
	// internal ram
	case 0xC000:
	case 0xD000:
		return iram1[address&0x1FFF]
	// echo of internal ram
	case 0xE000:
		return eiram[address&0x1FFF]
	case 0xF000:
		switch address & 0x0F00 {
		case 0x000:
		case 0x100:
		case 0x200:
		case 0x300:
		case 0x400:
		case 0x500:
		case 0x600:
		case 0x700:
		case 0x800:
		case 0x900:
		case 0xA00:
		case 0xB00:
		case 0xC00:
		case 0xD00:
			return eiram[address&0x1FFF]
		// OAM
		case 0xE00:
			return oam[address&0xFF]
		// zero page
		case 0xF00:
			if address >= 0xFF80 {
				return iram2[address&0x7F]
			} else {
				return io[address&0x7F]
			}
		}
	}
	return 0
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
