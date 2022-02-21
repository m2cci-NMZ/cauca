package main

import (
	"fmt"
)

func main() {
	var memory Memory
	var cpu Register
	var gpu Gpu
	var display Display
	memory.loadRom("roms/tetris")
	defer display.close()
	defer display.vramClose()
	//var i int = 0
	cpu.a = 0x01
	cpu.flags = 0xb0
	cpu.c = 0x13
	cpu.e = 0xd8
	cpu.h = 0x01
	cpu.l = 0x4d
	cpu.sp = 0xfffe
	cpu.pc = 0x100
	display.init()
	display.initVramViewer()
	//memory.writeByte(0xff44, 0x94)
	for display.running {
		//for cpu.pc != 0x2a0 {
		//debug
		if cpu.pc == 0x02a6 {
			fmt.Println(cpu.pc)
		}
		cpu.execute(memory.readByte(cpu.pc), &memory)
		gpu.step(cpu, &memory)
		display.display(gpu)
		display.displayVram(gpu, memory)
		//i++
	}
	//os.WriteFile("tile.bin", memory.vram[:], 0777)
}
