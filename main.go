package main

import (
	"fmt"
	"os"
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
	for i < 100000 {
		if (cpu.pc >= 0x2817) && (cpu.pc <= 0x282a) {
			fmt.Print(cpu.pc)
		}
		cpu.execute(memory.readByte(cpu.pc), &memory)
		if cpu.pc == 0x282a {
			var tmp []byte = memory.vram[:]
			os.WriteFile("tile.bin", tmp, 0)
		}
		//if cpu.pc > 10280 && cpu.pc < 10284 {
		//fmt.Print(cpu.pc, "\n")
		//}
		i++
	}
	//	fmt.Print(cpu.pc)
}
