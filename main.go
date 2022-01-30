package main

import "fmt"

func main() {
	var memory Memory
	var cpu Register
	memory.loadRom("test")
	var i uint16 = 0
	for i < 1000 {
		cpu.execute(memory.readByte(cpu.pc), &memory)
		i++
	}
	fmt.Print(memory.readByte(cpu.pc))
}
