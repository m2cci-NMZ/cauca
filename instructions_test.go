package main

import "testing"

func TestLdnnn(t *testing.T) {
	var reg Register
	dest := "B"
	var value byte = 10
	reg.ldnnn(value, dest)
	if reg.b != value {
		t.Errorf("%d for register A, expected %d", reg.a, value)
	}
}

func TestLdr1r2(t *testing.T) {
	var reg Register
	reg.a = 10
	reg.ldr1r2("D", "A")
	if reg.d != reg.a {
		t.Errorf("%d for register D, expected %d", reg.a, reg.d)
	}
}

func TestLdAn(t *testing.T) {
	var reg Register
	reg.c = 10
	reg.ldAn("C")
	if reg.a != reg.c {
		t.Errorf("%d for register A, expected %d", reg.a, reg.c)
	}
}

func TestLdnA(t *testing.T) {
	var reg Register
	reg.a = 10
	reg.ldAn("C")
	if reg.c != reg.a {
		t.Errorf("%d for register C, expected %d", reg.c, reg.a)
	}
}
func TestLdAC(t *testing.T) {
	var reg Register
	var value byte = 10
	var mem Memory
	mem.io[30] = value
	reg.c = 30
	reg.ldAC(&mem)
	if reg.a != value {
		t.Errorf("%d for register A, expected %d", reg.a, value)
	}
}

func TestLdCA(t *testing.T) {
	var reg Register
	var mem Memory
	reg.a = 10
	reg.c = 30
	reg.ldCA(&mem)
	if mem.io[30] != reg.a {
		t.Errorf("%d for io memory at adress register C, expected %d", mem.io[30], reg.a)
	}
}

func TestLddAHL(t *testing.T) {
	var reg Register
	var mem Memory
	reg.h = 2
	reg.l = 1
	mem.writeByte(513, 10)
	reg.lddAHL(&mem)
	// HL is 000001000000001 in binary (513 in decimal)
	if mem.readByte(513) != reg.a {
		t.Errorf("%d in register A, expected %d", reg.a, mem.readByte(513))
	}
	if reg.l != 0 {
		t.Errorf("%d in register L, expected 0", reg.l)
	}
	if reg.h != 2 {
		t.Errorf("%d in register H, expected 2", reg.l)
	}
}

func TestLddHLA(t *testing.T) {
	var reg Register
	var mem Memory
	reg.a = 10
	reg.h = 2
	reg.l = 1
	reg.lddAHL(&mem)
	// HL is 000001000000001 in binary (513 in decimal)
	if mem.readByte(513) != reg.a {
		t.Errorf("%d at memory at adress HL, expected %d", mem.readByte(513), reg.a)
	}
	if reg.l != 0 {
		t.Errorf("%d in register L, expected 0", reg.l)
	}
	if reg.h != 2 {
		t.Errorf("%d in register H, expected 2", reg.l)
	}
}

func TestLdiAHL(t *testing.T) {
	var reg Register
	var mem Memory
	reg.h = 2
	reg.l = 1
	mem.writeByte(513, 10)
	reg.ldiAHL(&mem)
	// HL is 000001000000001 in binary (513 in decimal)
	if mem.readByte(513) != reg.a {
		t.Errorf("%d in register A, expected %d", reg.a, mem.readByte(513))
	}
	if reg.l != 2 {
		t.Errorf("%d in register L, expected 2", reg.l)
	}
	if reg.h != 2 {
		t.Errorf("%d in register H, expected 2", reg.l)
	}
}

func TestLdiHLA(t *testing.T) {
	var reg Register
	var mem Memory
	reg.a = 10
	reg.h = 2
	reg.l = 1
	reg.ldiHLA(&mem)
	// HL is 000001000000001 in binary (513 in decimal)
	if mem.readByte(513) != reg.a {
		t.Errorf("%d at memory at adress HL, expected %d", mem.readByte(513), reg.a)
	}
	if reg.l != 2 {
		t.Errorf("%d in register L, expected 0", reg.l)
	}
	if reg.h != 2 {
		t.Errorf("%d in register H, expected 2", reg.l)
	}
}

func TestLdhnA(t *testing.T) {
	var reg Register
	var mem Memory
	reg.a = 10
	reg.ldhnA(1, &mem)
	var address uint16 = 0xFF00 + 0x0001
	if mem.readByte(address) != reg.a {
		t.Errorf("%d value in memory address 0xFF01, expected register A", reg.l)
	}
}

func TestLdhAn(t *testing.T) {
	var reg Register
	var mem Memory
	var address uint16 = 0xFF00 + 0x0001
	mem.writeByte(address, 10)
	reg.ldhnA(1, &mem)
	if mem.readByte(address) != reg.a {
		t.Errorf("%d value in memory address 0xFF01, expected register A", reg.l)
	}
}
