package main

import (
	"fmt"
)

func main() {
	var memory Memory
	var cpu Register
	memory.loadRom("roms/tetris")
	var i int = 0
	cpu.a = 0x01
	cpu.flags = 0xb0
	cpu.c = 0x13
	cpu.e = 0xd8
	cpu.h = 0x01
	cpu.l = 0x4d
	cpu.sp = 0xfffe
	cpu.pc = 0x100
	for i < 1000000 {
		//debug
		if cpu.pc == 0x2795 {
			fmt.Println(cpu.pc)
		}
		cpu.execute(memory.readByte(cpu.pc), &memory)
		i++
	}
}
