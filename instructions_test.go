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

func TestLdnnn16(t *testing.T) {
	var reg Register
	var value uint16 = 0x0102
	reg.ldnnn16(value, "BC")
	if reg.b != 1 {
		t.Errorf("%d value in register B, expected 1", reg.b)
	}
	if reg.c != 2 {
		t.Errorf("%d value in register C, expected 2", reg.b)
	}
}

func TestLdSPHL(t *testing.T) {
	var reg Register
	reg.h = 1
	reg.l = 2
	var value uint16 = 0x0102
	reg.ldSPHL()
	if reg.sp != value {
		t.Errorf("%d value in register sp, expected %d", reg.b, value)
	}
}

func TestLdHLSPn(t *testing.T) {
	var reg Register
	reg.sp = 0x0102
	reg.ldHLSPn(1)
	if reg.h != 1 {
		t.Errorf("%d in register h, expected 1", reg.h)
	}
	if reg.l != 3 {
		t.Errorf("%d in register l, expected 3", reg.l)
	}
}

func TestLdnnSP(t *testing.T) {
	var reg Register
	var mem Memory
	reg.sp = 10
	reg.ldnnSP(10, &mem)
	if mem.readWord(10) != reg.sp {
		t.Errorf("%d at memory address 10, expected %d (sp)", mem.readWord(10), reg.sp)
	}
}

func TestPushnn(t *testing.T) {
	var reg Register
	var mem Memory
	reg.b = 1
	reg.c = 1
	reg.sp = 10
	reg.pushnn("BC", &mem)
	if mem.readWord(10) != 1 {
		t.Errorf("%d at memory address 10, expected %d", mem.readWord(10), reg.b)
	}
	if mem.readWord(11) != 1 {
		t.Errorf("%d at memory address 11, expected %d", mem.readWord(11), reg.c)
	}
	if reg.sp != 8 {
		t.Errorf("%d in sp, expected 12", reg.sp)
	}
}

func TestPopnn(t *testing.T) {
	var reg Register
	var mem Memory
	reg.sp = 10
	mem.rom[10] = 1
	mem.rom[11] = 1
	reg.popnn("BC", &mem)
	if reg.b != 1 {
		t.Errorf("%d at register B, expected %d", reg.b, mem.readWord(10))
	}
	if mem.readWord(11) != 1 {
		t.Errorf("%d at register C, expected %d", reg.c, mem.readWord(11))
	}
	if reg.sp != 12 {
		t.Errorf("%d in sp, expected 12", reg.sp)
	}
}

func TestAddAn(t *testing.T) {
	var value byte = 1
	var reg Register
	reg.addAn(value)
	if reg.a != value {
		t.Errorf("%d in register a, expected %d", reg.a, value)
	}
}

func TestAddcAn(t *testing.T) {
	var value byte = 1
	var reg Register
	reg.flags = 1 << 4
	reg.addcAn(value)
	if reg.a != (value + 1) {
		t.Errorf("%d in register a, expected %d", reg.a, value+1)
	}
}

func TestSubn(t *testing.T) {
	var value byte = 1
	var reg Register
	reg.a = 1
	reg.subn(value)
	if reg.a != 0 {
		t.Errorf("%d in register a, expected 0", reg.a)
	}
}

func TestSbcAn(t *testing.T) {
	var value byte = 1
	var reg Register
	reg.flags = 1 << 4
	reg.a = 1
	reg.sbcAn(value)
	if reg.a != 1 {
		t.Errorf("%d in register A, expected 1", reg.a)
	}
}

func TestAndn(t *testing.T) {
	var value byte = 0
	var reg Register
	reg.a = 1
	reg.andn(value)
	if reg.a != 0 {
		t.Errorf("%d in register A, expected 0", reg.a)
	}
}

func TestOrn(t *testing.T) {
	var value byte = 0
	var reg Register
	reg.a = 1
	reg.orn(value)
	if reg.a != 1 {
		t.Errorf("%d in register A, expected 1", reg.a)
	}
}

func TestXorn(t *testing.T) {
	var value byte = 1
	var reg Register
	reg.a = 1
	reg.xorn(value)
	if reg.a != 0 {
		t.Errorf("%d in register A, expected 1", reg.a)
	}
}

func TestCpn(t *testing.T) {
	var value byte = 1
	var reg Register
	reg.a = 1
	reg.cpn(value)
	if !hasBit(uint16(reg.flags), 7) {
		t.Errorf("expected zero flag to be 1, is 0")
	}
	value++
	reg.cpn(value)
	if !hasBit(uint16(reg.flags), 4) {
		t.Errorf("expected carry flag to be 1, is 0")
	}
}

func TestIncn(t *testing.T) {
	var reg Register
	reg.incn("A")
	if reg.a != 1 {
		t.Errorf("%d in register A, expected 1", reg.a)
	}
}

func TestDecn(t *testing.T) {
	var reg Register
	reg.a = 1
	reg.decn("A")
	if reg.a != 0 {
		t.Errorf("%d in register A, expected 0", reg.a)
	}
}
