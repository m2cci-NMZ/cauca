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
	var mem Memory
	mem.writeByte(10, 10)
	reg.c = 10
	reg.ldAn("C", &mem)
	if reg.a != mem.readByte(uint16(reg.c)) {
		t.Errorf("%d for register A, expected %d", reg.a, mem.readByte(uint16(reg.c)))
	}
}

func TestLdnA(t *testing.T) {
	var reg Register
	var mem Memory
	reg.a = 10
	reg.c = 20
	reg.ldnA("C", &mem)
	if mem.readByte(20) != reg.a {
		t.Errorf("%d for memory at adress 20, expected %d", mem.readByte(20), reg.a)
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

func TestAddHLn(t *testing.T) {
	var reg Register
	var value uint16 = 0x0102
	reg.addHLn(value)
	if reg.h != 1 {
		t.Errorf("%d in register H, expected 1", reg.h)
	}
	if reg.l != 2 {
		t.Errorf("%d in register L, expected 2", reg.l)
	}
}

func TestAddSPn(t *testing.T) {
	var reg Register
	var value uint16 = 0x0102
	reg.addSPn(value)
	if reg.sp != value {
		t.Errorf("%d in stack pointer, expected %d", reg.sp, value)
	}
}

func TestIncnn(t *testing.T) {
	var reg Register
	reg.b = 1
	reg.c = 1
	reg.incnn("BC")
	if reg.b != 1 {
		t.Errorf("%d in register b, expected 1", reg.b)
	}
	if reg.c != 2 {
		t.Errorf("%d in register c, expected 2", reg.c)
	}
}

func TestDecnn(t *testing.T) {
	var reg Register
	reg.b = 1
	reg.c = 1
	reg.decnn("BC")
	if reg.b != 1 {
		t.Errorf("%d in register b, expected 1", reg.b)
	}
	if reg.c != 0 {
		t.Errorf("%d in register c, expected 0", reg.c)
	}
}

func TestSwapn(t *testing.T) {
	var reg Register
	reg.a = 1
	reg.swapn("A")
	if reg.a != 0x10 {
		t.Errorf("%d in register A, expected 128", reg.a)
	}
}

func TestDaa(t *testing.T) {
	var reg Register
	reg.a = 1
	reg.setRegisterFlag(true, 5)
	reg.dAA()
	if reg.a != 7 {
		t.Errorf("%d in register A, expected 97", reg.a)
	}
}

func TestCpl(t *testing.T) {
	var reg Register
	reg.a = 1
	reg.cpl()
	if reg.a != 254 {
		t.Errorf("%d in register A, expexted 254", reg.a)
	}
}

func TestCcf(t *testing.T) {
	var reg Register
	reg.setRegisterFlag(true, 4)
	reg.ccf()
	if reg.flags != 0 {
		t.Errorf("failed to reset carry flag")
	}
}

func TestScf(t *testing.T) {
	var reg Register
	reg.scf()
	if reg.flags != 16 {
		t.Errorf("failed to set carry flag")
	}
}

func TestRlcA(t *testing.T) {
	var reg Register
	reg.a = 129
	reg.rlcA()
	if reg.a != 2 {
		t.Errorf("%d in register A, expected 2", reg.a)
	}
	if !hasBit(uint16(reg.flags), 4) {
		t.Errorf("bit 7 of register A was not carried")
	}
}

func TestRlA(t *testing.T) {
	var reg Register
	reg.a = 129
	reg.setRegisterFlag(true, 4)
	reg.rlA()
	if reg.a != 3 {
		t.Errorf("%d in register A, expected 3", reg.a)
	}
	if !hasBit(uint16(reg.flags), 4) {
		t.Errorf("bit 7 of register A was not carried")
	}
}

func TestRrcA(t *testing.T) {
	var reg Register
	reg.a = 129
	reg.rrcA()
	if reg.a != 64 {
		t.Errorf("%d in register A, expected 64", reg.a)
	}
	if !hasBit(uint16(reg.flags), 4) {
		t.Errorf("carry flag not set")
	}
}

func TestRrA(t *testing.T) {
	var reg Register
	reg.a = 129
	reg.setRegisterFlag(true, 4)
	reg.rrA()
	if reg.a != 192 {
		t.Errorf("%d in register A, expected 64", reg.a)
	}
	if !hasBit(uint16(reg.flags), 4) {
		t.Errorf("carry flag not set")
	}
}

func TestRlcn(t *testing.T) {
	var reg Register
	reg.b = 129
	reg.rlcn("B")
	if reg.b != 2 {
		t.Errorf("%d in register A, expected 2", reg.a)
	}
	if !hasBit(uint16(reg.flags), 4) {
		t.Errorf("bit 7 of register A was not carried")
	}
}

func TestRln(t *testing.T) {
	var reg Register
	reg.b = 129
	reg.setRegisterFlag(true, 4)
	reg.rln("B")
	if reg.b != 3 {
		t.Errorf("%d in register A, expected 3", reg.a)
	}
	if !hasBit(uint16(reg.flags), 4) {
		t.Errorf("bit 7 of register A was not carried")
	}
}

func TestRrcn(t *testing.T) {
	var reg Register
	reg.b = 129
	reg.rrcn("B")
	if reg.b != 64 {
		t.Errorf("%d in register A, expected 64", reg.a)
	}
	if !hasBit(uint16(reg.flags), 4) {
		t.Errorf("carry flag not set")
	}
}

func TestRrn(t *testing.T) {
	var reg Register
	reg.b = 129
	reg.setRegisterFlag(true, 4)
	reg.rrn("B")
	if reg.b != 192 {
		t.Errorf("%d in register A, expected 64", reg.a)
	}
	if !hasBit(uint16(reg.flags), 4) {
		t.Errorf("carry flag not set")
	}
}

func TestSlAn(t *testing.T) {
	var reg Register
	reg.a = 1
	reg.slan("A")
	if reg.a != 2 {
		t.Errorf("%d in register A, expected 2", reg.a)
	}
}

func TestSrAn(t *testing.T) {
	var reg Register
	reg.a = 130
	reg.sran("A")
	if reg.a != 193 {
		t.Errorf("%d in register A, expected 1", reg.a)
	}
}

func TestSrln(t *testing.T) {
	var reg Register
	reg.a = 2
	reg.srln("A")
	if reg.a != 1 {
		t.Errorf("%d in register A, expected 1", reg.a)
	}
}

func TestBitBr(t *testing.T) {
	var reg Register
	reg.a = 1
	reg.bitBr("A", 0)
	if !hasBit(uint16(reg.flags), 7) {
		t.Errorf("zero flag not set")
	}
}

func TestSetBr(t *testing.T) {
	var reg Register
	reg.setBr("A", 0)
	if reg.a != 1 {
		t.Errorf("%d in register A, expected 1", reg.a)
	}
}

func TestResBr(t *testing.T) {
	var reg Register
	reg.a = 1
	reg.resBr("A", 0)
	if reg.a != 0 {
		t.Errorf("%d in register A, expected 0", reg.a)
	}
}

func TestJpnn(t *testing.T) {
	var reg Register
	var value uint16 = 10
	reg.jpnn(value)
	if reg.pc != value {
		t.Errorf("%d in program counter, expected %d", reg.pc, value)
	}
}

func TestJpccnn(t *testing.T) {
	var reg Register
	var value uint16 = 10
	reg.setRegisterFlag(false, 7)
	reg.jpccnn(value, "Z")
	if reg.pc != 0 {
		t.Errorf("%d in program counter, expected 0", reg.pc)
	}
}

func TestJpHL(t *testing.T) {
	var reg Register
	reg.h = 1
	reg.l = 1
	reg.jpHL()
	if reg.pc != concatenateBytes(reg.h, reg.l) {
		t.Errorf("%d in program counter, expected %d", reg.pc, concatenateBytes(reg.h, reg.l))
	}
}

func TestJrn(t *testing.T) {
	var reg Register
	reg.pc = 1
	reg.jrn(10)
	if reg.pc != 11 {
		t.Errorf("%d in program counter, expected 11", reg.pc)
	}
}

func TestJrccnn(t *testing.T) {
	var reg Register
	reg.pc = 1
	reg.setRegisterFlag(false, 7)
	reg.jrccn(10, "Z")
	if reg.pc != 1 {
		t.Errorf("%d in program counter, expected 1", reg.pc)
	}
}
