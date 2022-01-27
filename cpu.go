package main

import "strconv"

type Register struct {
	a     byte
	b     byte
	c     byte
	d     byte
	e     byte
	h     byte
	l     byte
	flags byte
	sp    uint16
	pc    uint16
}

func concatenateBytes(a byte, b byte) uint16 {
	result := (uint16(a) << 8) + uint16(b)
	return result
}

/* *************************************** */
/* Flags setting function                  */
/* *************************************** */

func (register *Register) setRegisterFlag(value bool, position byte) {
	if value {
		register.flags |= (1 << position)
	} else {
		register.flags &^= (1 << position)
	}
}

func hasBit(n uint16, pos uint16) bool {
	val := n & (1 << pos)
	return (val > 0)
}

/* *************************************** */
/* 8 bit loads                             */
/* *************************************** */

func (reg *Register) ldnnn(value byte, destination string) {
	switch destination {
	case "B":
		reg.b = value
	case "C":
		reg.c = value
	case "D":
		reg.d = value
	case "E":
		reg.e = value
	case "H":
		reg.h = value
	case "L":
		reg.l = value
	}
}

func (reg *Register) ldr1r2(destination string, source string, mem *Memory) {
	reg_map := map[string]*byte{
		"A": &(reg.a),
		"B": &(reg.b),
		"C": &(reg.c),
		"D": &(reg.d),
		"E": &(reg.e),
		"H": &(reg.h),
		"L": &(reg.l),
	}
	if source == "HL" {
		*(reg_map[destination]) = mem.readByte(concatenateBytes(reg.h, reg.l))
	} else if destination == "HL" {
		value := mem.readWord(uint16(*(reg_map[source])))
		reg.h = byte(value & 0xff00)
		reg.l = byte(value & 0x00ff)
	} else {
		*(reg_map[destination]) = mem.readByte(uint16(*(reg_map[source])))
	}
}

func (reg *Register) ldAn(source string, mem *Memory) {
	switch source {
	case "A":
		reg.a = mem.readByte(uint16(reg.a))
	case "B":
		reg.a = mem.readByte(uint16(reg.b))
	case "C":
		reg.a = mem.readByte(uint16(reg.c))
	case "D":
		reg.a = mem.readByte(uint16(reg.d))
	case "E":
		reg.a = mem.readByte(uint16(reg.e))
	case "H":
		reg.a = mem.readByte(uint16(reg.h))
	case "L":
		reg.a = mem.readByte(uint16(reg.l))
	case "BC":
		reg.a = mem.readByte(concatenateBytes(reg.b, reg.c))
	case "DE":
		reg.a = mem.readByte(concatenateBytes(reg.d, reg.e))
	case "HL":
		reg.a = mem.readByte(concatenateBytes(reg.h, reg.l))
	default:
		address, _ := strconv.ParseInt(source, 16, 16)
		reg.a = mem.readByte(uint16(address))
	}
}

func (reg *Register) ldnA(source string, mem *Memory) {
	switch source {
	case "A":
		mem.writeByte(uint16(reg.a), reg.a)
	case "B":
		mem.writeByte(uint16(reg.b), reg.a)
	case "C":
		mem.writeByte(uint16(reg.c), reg.a)
	case "D":
		mem.writeByte(uint16(reg.d), reg.a)
	case "E":
		mem.writeByte(uint16(reg.e), reg.a)
	case "H":
		mem.writeByte(uint16(reg.h), reg.a)
	case "L":
		mem.writeByte(uint16(reg.l), reg.a)
	case "BC":
		mem.writeByte(concatenateBytes(reg.b, reg.c), reg.a)
	case "DE":
		mem.writeByte(concatenateBytes(reg.d, reg.e), reg.a)
	case "HL":
		mem.writeByte(concatenateBytes(reg.h, reg.l), reg.a)
	default:
		address, _ := strconv.ParseInt(source, 16, 16)
		mem.writeByte(uint16(address), reg.a)
	}
}

func (reg *Register) ldAC(mem *Memory) {
	reg.a = mem.io[reg.c]
}

func (reg *Register) ldCA(mem *Memory) {
	mem.io[reg.c] = reg.a
}

func (reg *Register) lddAHL(mem *Memory) {
	address := (uint16(reg.h) << 8) + uint16(reg.l)
	reg.a = mem.readByte(address)
	address--
	reg.h = byte(address >> 8)
	reg.l = byte(address)
}

func (reg *Register) lddHLA(mem *Memory) {
	address := (uint16(reg.h) << 8) + uint16(reg.l)
	mem.io[address] = reg.a
	address--
	reg.h = byte(address >> 8)
	reg.l = byte(address)
}

func (reg *Register) ldiAHL(mem *Memory) {
	address := (uint16(reg.h) << 8) + uint16(reg.l)
	reg.a = mem.readByte(address)
	address++
	reg.h = byte(address >> 8)
	reg.l = byte(address)
}

func (reg *Register) ldiHLA(mem *Memory) {
	address := (uint16(reg.h) << 8) + uint16(reg.l)
	mem.writeByte(address, reg.a)
	address++
	reg.h = byte(address >> 8)
	reg.l = byte(address)
}

func (reg *Register) ldhnA(value byte, mem *Memory) {
	mem.io[value] = reg.a
}

func (reg *Register) ldhAn(value byte, mem *Memory) {
	reg.a = mem.io[value]
}

/* *************************************** */
/* 16 bit loads                            */
/* *************************************** */

func (reg *Register) ldnnn16(value uint16, destination string) {
	r1 := byte(value >> 8)
	r2 := byte(value)
	switch destination {
	case "BC":
		reg.b = r1
		reg.c = r2
	case "DE":
		reg.d = r1
		reg.e = r2
	case "HL":
		reg.h = r1
		reg.l = r2
	case "SP":
		reg.sp = value
	}
}

func (reg *Register) ldSPHL() {
	value := (uint16(reg.h) << 8) + uint16(reg.l)
	reg.sp = value
}

func (reg *Register) ldHLSPn(value byte) {
	result := uint16(value) + reg.sp
	r1 := byte(result >> 8)
	r2 := byte(result)
	reg.h = r1
	reg.l = r2
	// reset Z flag
	reg.setRegisterFlag(false, 7)
	// reset N flag
	reg.setRegisterFlag(false, 6)
	// set H flag
	if (reg.sp&0x000F + uint16(value&0x0F)) > 0x0F {
		reg.setRegisterFlag(true, 5)
	} else {
		reg.setRegisterFlag(false, 5)
	}
	//set C flag
	if (result & 0xFF00) != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
}

func (reg *Register) ldnnSP(value uint16, mem *Memory) {
	mem.writeWord(value, reg.sp)
}

func (reg *Register) pushnn(registers string, mem *Memory) {
	var r1, r2 uint16
	switch registers {
	case "AF":
		r1 = uint16(reg.a) << 8
		r2 = uint16(reg.b)
	case "BC":
		r1 = uint16(reg.b) << 8
		r2 = uint16(reg.c)
	case "DE":
		r1 = uint16(reg.d) << 8
		r2 = uint16(reg.e)
	case "HL":
		r1 = uint16(reg.h) << 8
		r2 = uint16(reg.l)
	}
	value := r1 + r2
	mem.writeWord(reg.sp, value)
	reg.sp = reg.sp - 2
}

func (reg *Register) popnn(registers string, mem *Memory) {
	r1 := mem.readByte(reg.sp)
	r2 := mem.readByte(reg.sp + 1)
	switch registers {
	case "AF":
		reg.a = r1
		reg.b = r2
	case "BC":
		reg.b = r1
		reg.c = r2
	case "DE":
		reg.d = r1
		reg.e = r2
	case "HL":
		reg.h = r1
		reg.l = r2
	}
	reg.sp = reg.sp + 2
}

/* *************************************** */
/* 8 bit ALU                               */
/* *************************************** */

func (reg *Register) addAn(value byte) {
	result := uint16(reg.a) + uint16(value)
	// carry flag
	if (result & 0xFF00) != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	// zero flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	if (reg.a&0x0F + value&0x0F) > 0x0F {
		reg.setRegisterFlag(true, 5)
	} else {
		reg.setRegisterFlag(false, 5)
	}
	reg.a = byte(result & 0xFF)
}

func (reg *Register) addcAn(value byte) {
	if hasBit(uint16(reg.flags), 4) {
		value++
	}
	reg.addAn(value)
}

func (reg *Register) subn(value byte) {
	result := uint16(reg.a) - uint16(value)
	// negative flag
	if result < 0 {
		reg.setRegisterFlag(true, 6)
	} else {
		reg.setRegisterFlag(false, 6)
	}
	// zero flag
	if result == 0 {
		reg.setRegisterFlag(true, 7)
	}
	// half carry flag
	if (reg.a&0x0F + value&0x0F) > 0x0F {
		reg.setRegisterFlag(true, 5)
	} else {
		reg.setRegisterFlag(false, 5)
	}
	// carry flag
	if (result & 0xFF00) != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	reg.a = byte(result & 0xFF)
}

func (reg *Register) sbcAn(value byte) {
	if hasBit(uint16(reg.flags), 4) {
		value--
	}
	reg.subn(value)
}

func (reg *Register) andn(value byte) {
	result := reg.a & value
	// zero flag
	if result == 0 {
		reg.setRegisterFlag(true, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(true, 5)
	// carry flag
	reg.setRegisterFlag(false, 4)
	reg.a = result
}

func (reg *Register) orn(value byte) {
	result := reg.a | value
	// zero flag
	if result == 0 {
		reg.setRegisterFlag(true, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	// carry flag
	reg.setRegisterFlag(false, 4)
	reg.a = result
}

func (reg *Register) xorn(value byte) {
	result := reg.a ^ value
	// zero flag
	if result == 0 {
		reg.setRegisterFlag(true, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	// carry flag
	reg.setRegisterFlag(false, 4)
	reg.a = result
}

func (reg *Register) cpn(value byte) {
	tmp := reg.a
	reg.subn(value)
	reg.a = tmp
}

func (reg *Register) incn(register string) {
	var result uint16
	result++
	switch register {
	case "A":
		reg.a = byte(result)
	case "B":
		reg.b = byte(result)
	case "C":
		reg.c = byte(result)
	case "D":
		reg.d = byte(result)
	case "E":
		reg.e = byte(result)
	case "H":
		reg.h = byte(result)
	case "L":
		reg.l = byte(result)
	}
	// carry flag
	if (result & 0xFF00) != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	// zero flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	if (uint16(reg.a&0x0F) + result&0x000F) > 0x0F {
		reg.setRegisterFlag(true, 5)
	} else {
		reg.setRegisterFlag(false, 5)
	}
}

func (reg *Register) decn(register string) {
	var result int16
	switch register {
	case "A":
		result = int16(reg.a) - 1
		reg.a = byte(result)
	case "B":
		result = int16(reg.a) - 1
		reg.b = byte(result)
	case "C":
		result = int16(reg.a) - 1
		reg.c = byte(result)
	case "D":
		result = int16(reg.a) - 1
		reg.d = byte(result)
	case "E":
		result = int16(reg.a) - 1
		reg.e = byte(result)
	case "H":
		result = int16(reg.a) - 1
		reg.h = byte(result)
	case "L":
		result = int16(reg.a) - 1
		reg.l = byte(result)
	}
	// negative flag
	if result < 0 {
		reg.setRegisterFlag(true, 6)
	} else {
		reg.setRegisterFlag(false, 6)
	}
	// zero flag
	if result == 0 {
		reg.setRegisterFlag(true, 7)
	}
	// half carry flag
	if (int16(reg.a&0x0f) + result&0x000f) > 0x0f {
		reg.setRegisterFlag(true, 5)
	} else {
		reg.setRegisterFlag(false, 5)
	}
	// carry flag
	if (int32(result) & 0xff00) != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
}

/* *************************************** */
/* 16 bit ALU                               */
/* *************************************** */

func (reg *Register) addHLn(value uint16) {
	var result uint32
	HL := uint16(reg.h)<<8 + uint16(reg.l)
	result = uint32(HL) + uint32(value)
	// negative flag
	reg.setRegisterFlag(false, 6)
	// carry flag
	if (result & 0xFF0000) != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	// half carry flag
	if ((uint16(result) & 0x0F) + (value & 0x0F)) > 0x0F {
		reg.setRegisterFlag(true, 5)
	} else {
		reg.setRegisterFlag(false, 5)
	}
	reg.h = byte(result >> 8)
	reg.l = byte(result)
}

func (reg *Register) addSPn(value uint16) {
	result := uint32(reg.sp) + uint32(value)
	// zero flag
	reg.setRegisterFlag(false, 7)
	// negative flag
	reg.setRegisterFlag(false, 6)
	// carry flag
	if (result & 0xFF0000) != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	// half carry flag
	if ((uint16(result) & 0x0F) + (value & 0x0F)) > 0x0F {
		reg.setRegisterFlag(true, 5)
	} else {
		reg.setRegisterFlag(false, 5)
	}
	reg.sp = uint16(result)
}

func (reg *Register) incnn(register string) {
	switch register {
	case "BC":
		result := uint16(reg.b)<<8 + uint16(reg.c) + 1
		reg.b = byte(result >> 8)
		reg.c = byte(result)
	case "DE":
		result := uint16(reg.d)<<8 + uint16(reg.e) + 1
		reg.d = byte(result >> 8)
		reg.d = byte(result)
	case "HL":
		result := uint16(reg.h)<<8 + uint16(reg.l) + 1
		reg.h = byte(result >> 8)
		reg.l = byte(result)
	case "SP":
		reg.sp++
	}
}

func (reg *Register) decnn(register string) {
	switch register {
	case "BC":
		result := uint16(reg.b)<<8 + uint16(reg.c) - 1
		reg.b = byte(result >> 8)
		reg.c = byte(result)
	case "DE":
		result := uint16(reg.d)<<8 + uint16(reg.e) - 1
		reg.d = byte(result >> 8)
		reg.e = byte(result)
	case "HL":
		result := uint16(reg.h)<<8 + uint16(reg.l) - 1
		reg.h = byte(result >> 8)
		reg.l = byte(result)
	case "SP":
		reg.sp--
	}
}

/* *************************************** */
/* misc                                    */
/* *************************************** */

func (reg *Register) swapn(register string) {
	var value byte
	switch register {
	case "A":
		value := reg.a
		value = ((value & 0x0f) << 4) | ((value & 0xf0) >> 4)
		reg.a = value
	case "B":
		value := reg.b
		value = ((value & 0x0f) << 4) | ((value & 0xf0) >> 4)
		reg.b = value
	case "C":
		value := reg.c
		value = ((value & 0x0f) << 4) | ((value & 0xf0) >> 4)
		reg.c = value
	case "D":
		value := reg.d
		value = ((value & 0x0f) << 4) | ((value & 0xf0) >> 4)
		reg.d = value
	case "E":
		value := reg.e
		value = ((value & 0x0f) << 4) | ((value & 0xf0) >> 4)
		reg.e = value
	case "H":
		value := reg.h
		value = ((value & 0x0f) << 4) | ((value & 0xf0) >> 4)
		reg.h = value
	case "L":
		value := reg.l
		value = ((value & 0x0f) << 4) | ((value & 0xf0) >> 4)
		reg.l = value
	}
	// zero flag
	if value == 0 {
		reg.setRegisterFlag(true, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	// carry flag
	reg.setRegisterFlag(false, 4)
}

func (reg *Register) dAA() {
	var value uint16
	value = uint16(reg.a)
	if hasBit(uint16(reg.flags), 6) {
		if hasBit(uint16(reg.flags), 5) {
			value = (value - 0x06) & 0xff
		}
		if hasBit(uint16(reg.flags), 4) {
			value -= 0x60
		}
	} else {
		if hasBit(uint16(reg.flags), 5) || ((value & 0xf) > 9) {
			value += 0x06
		}
		if hasBit(uint16(reg.flags), 4) || value > 0x9f {
			value += 0x60
		}
	}
	// half carry flag
	reg.setRegisterFlag(false, 5)
	// zero flag
	if reg.a != 0 {
		reg.setRegisterFlag(false, 7)
	} else {
		reg.setRegisterFlag(true, 7)
	}
	// carry flag
	if value >= 0x100 {
		reg.setRegisterFlag(true, 4)
	}
	reg.a = byte(value)
}

func (reg *Register) cpl() {
	reg.a = ^reg.a
	// negative flag
	reg.setRegisterFlag(true, 6)
	// haf carry flag
	reg.setRegisterFlag(true, 5)
}

func (reg *Register) ccf() {
	if hasBit(uint16(reg.flags), 4) {
		reg.setRegisterFlag(false, 4)
	} else {
		reg.setRegisterFlag(true, 4)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
}

func (reg *Register) scf() {
	// carry flag
	reg.setRegisterFlag(true, 4)
	// negative flag
	reg.setRegisterFlag(false, 6)
	// haf carry flag
	reg.setRegisterFlag(false, 5)
}

/* *************************************** */
/* rotates and shifts                      */
/* *************************************** */

// this is VERY tricky, maybe worth writing unit tests
func (reg *Register) rlA() {
	value := uint16(reg.a) << 1
	if hasBit(uint16(reg.flags), 4) {
		value++
	}
	//zero flag
	if value != 0 {
		reg.setRegisterFlag(false, 7)
	} else {
		reg.setRegisterFlag(true, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	//carry flag
	if (reg.a & 0x80) != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	reg.a = byte(value)
}

func (reg *Register) rlcA() {
	value := uint16(reg.a) << 1
	//zero flag
	if value != 0 {
		reg.setRegisterFlag(false, 7)
	} else {
		reg.setRegisterFlag(true, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	//carry flag
	if (reg.a & 0x80) != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	reg.a = byte(value)
}

func (reg *Register) rrA() {
	value := uint16(reg.a) >> 1
	if hasBit(uint16(reg.flags), 4) {
		var pos uint16 = 7
		value |= 1 << pos
	}
	//zero flag
	if value != 0 {
		reg.setRegisterFlag(false, 7)
	} else {
		reg.setRegisterFlag(true, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	//carry flag
	if reg.a&0x01 != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	reg.a = byte(value)
}

func (reg *Register) rrcA() {
	value := uint16(reg.a) >> 1
	//zero flag
	if value != 0 {
		reg.setRegisterFlag(false, 7)
	} else {
		reg.setRegisterFlag(true, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	//carry flag
	if reg.a&0x01 != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	reg.a = byte(value)
}

func (reg *Register) rln(destination string) {
	var register byte
	switch destination {
	case "A":
		register = reg.a
	case "B":
		register = reg.b
	case "C":
		register = reg.c
	case "D":
		register = reg.d
	case "E":
		register = reg.e
	case "H":
		register = reg.h
	case "L":
		register = reg.l
	}
	value := uint16(register) << 1
	if hasBit(uint16(reg.flags), 4) {
		value++
	}
	//zero flag
	if value != 0 {
		reg.setRegisterFlag(false, 7)
	} else {
		reg.setRegisterFlag(true, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	//carry flag
	if (register & 0x80) != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	switch destination {
	case "A":
		reg.a = byte(value)
	case "B":
		reg.b = byte(value)
	case "C":
		reg.c = byte(value)
	case "D":
		reg.d = byte(value)
	case "E":
		reg.e = byte(value)
	case "H":
		reg.h = byte(value)
	case "L":
		reg.l = byte(value)
	}
}

func (reg *Register) rlcn(destination string) {
	var register byte
	switch destination {
	case "A":
		register = reg.a
	case "B":
		register = reg.b
	case "C":
		register = reg.c
	case "D":
		register = reg.d
	case "E":
		register = reg.e
	case "H":
		register = reg.h
	case "L":
		register = reg.l
	}
	value := uint16(register) << 1
	//zero flag
	if value != 0 {
		reg.setRegisterFlag(false, 7)
	} else {
		reg.setRegisterFlag(true, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	//carry flag
	if (register & 0x80) != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	switch destination {
	case "A":
		reg.a = byte(value)
	case "B":
		reg.b = byte(value)
	case "C":
		reg.c = byte(value)
	case "D":
		reg.d = byte(value)
	case "E":
		reg.e = byte(value)
	case "H":
		reg.h = byte(value)
	case "L":
		reg.l = byte(value)
	}
}

func (reg *Register) rrcn(destination string) {
	var register byte
	switch destination {
	case "A":
		register = reg.a
	case "B":
		register = reg.b
	case "C":
		register = reg.c
	case "D":
		register = reg.d
	case "E":
		register = reg.e
	case "H":
		register = reg.h
	case "L":
		register = reg.l
	}
	value := uint16(register) >> 1
	//zero flag
	if value != 0 {
		reg.setRegisterFlag(false, 7)
	} else {
		reg.setRegisterFlag(true, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	//carry flag
	if register&0x01 != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	switch destination {
	case "A":
		reg.a = byte(value)
	case "B":
		reg.b = byte(value)
	case "C":
		reg.c = byte(value)
	case "D":
		reg.d = byte(value)
	case "E":
		reg.e = byte(value)
	case "H":
		reg.h = byte(value)
	case "L":
		reg.l = byte(value)
	}
}
func (reg *Register) rrn(destination string) {
	var register byte
	switch destination {
	case "A":
		register = reg.a
	case "B":
		register = reg.b
	case "C":
		register = reg.c
	case "D":
		register = reg.d
	case "E":
		register = reg.e
	case "H":
		register = reg.h
	case "L":
		register = reg.l
	}
	value := uint16(register) >> 1
	if hasBit(uint16(reg.flags), 4) {
		var pos uint16 = 7
		value |= 1 << pos
	}
	//zero flag
	if value != 0 {
		reg.setRegisterFlag(false, 7)
	} else {
		reg.setRegisterFlag(true, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	//carry flag
	if register&0x01 != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	switch destination {
	case "A":
		reg.a = byte(value)
	case "B":
		reg.b = byte(value)
	case "C":
		reg.c = byte(value)
	case "D":
		reg.d = byte(value)
	case "E":
		reg.e = byte(value)
	case "H":
		reg.h = byte(value)
	case "L":
		reg.l = byte(value)
	}
}

func (reg *Register) slan(destination string) {
	var register byte
	switch destination {
	case "A":
		register = reg.a
	case "B":
		register = reg.b
	case "C":
		register = reg.c
	case "D":
		register = reg.d
	case "E":
		register = reg.e
	case "H":
		register = reg.h
	case "L":
		register = reg.l
	}
	value := register << 1
	var mask byte = 1
	value &^= mask
	//zero flag
	if value != 0 {
		reg.setRegisterFlag(false, 7)
	} else {
		reg.setRegisterFlag(true, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	//carry flag
	if register != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	switch destination {
	case "A":
		reg.a = value
	case "B":
		reg.b = value
	case "C":
		reg.c = value
	case "D":
		reg.d = value
	case "E":
		reg.e = value
	case "H":
		reg.h = value
	case "L":
		reg.l = value
	}
}

func (reg *Register) sran(destination string) {
	var register byte
	switch destination {
	case "A":
		register = reg.a
	case "B":
		register = reg.b
	case "C":
		register = reg.c
	case "D":
		register = reg.d
	case "E":
		register = reg.e
	case "H":
		register = reg.h
	case "L":
		register = reg.l
	}
	value := register >> 1
	if register&0x80 != 0 {
		value |= (1 << 7)
	}
	//zero flag
	if value != 0 {
		reg.setRegisterFlag(false, 7)
	} else {
		reg.setRegisterFlag(true, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	//carry flag
	if register != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	switch destination {
	case "A":
		reg.a = value
	case "B":
		reg.b = value
	case "C":
		reg.c = value
	case "D":
		reg.d = value
	case "E":
		reg.e = value
	case "H":
		reg.h = value
	case "L":
		reg.l = value
	}
}

func (reg *Register) srln(destination string) {
	var register byte
	switch destination {
	case "A":
		register = reg.a
	case "B":
		register = reg.b
	case "C":
		register = reg.c
	case "D":
		register = reg.d
	case "E":
		register = reg.e
	case "H":
		register = reg.h
	case "L":
		register = reg.l
	}
	value := register >> 1
	var mask byte = 1
	value &^= (mask << 7)
	//zero flag
	if value != 0 {
		reg.setRegisterFlag(false, 7)
	} else {
		reg.setRegisterFlag(true, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	//carry flag
	if register != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	switch destination {
	case "A":
		reg.a = value
	case "B":
		reg.b = value
	case "C":
		reg.c = value
	case "D":
		reg.d = value
	case "E":
		reg.e = value
	case "H":
		reg.h = value
	case "L":
		reg.l = value
	}
}

/* *************************************** */
/* Bit opcodes                             */
/* *************************************** */

func (reg *Register) bitBr(destination string, pos uint16) {
	var test bool
	switch destination {
	case "A":
		test = hasBit(uint16(reg.a), pos)
	case "B":
		test = hasBit(uint16(reg.b), pos)
	case "C":
		test = hasBit(uint16(reg.c), pos)
	case "D":
		test = hasBit(uint16(reg.d), pos)
	case "E":
		test = hasBit(uint16(reg.e), pos)
	case "H":
		test = hasBit(uint16(reg.h), pos)
	case "L":
		test = hasBit(uint16(reg.l), pos)
	}
	// zero flag
	reg.setRegisterFlag(test, 7)
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(true, 5)
}

func (reg *Register) setBr(destination string, pos uint16) {
	var mask byte = 1
	switch destination {
	case "A":
		reg.a |= (mask << pos)
	case "B":
		reg.b |= (mask << pos)
	case "C":
		reg.c |= (mask << pos)
	case "D":
		reg.d |= (mask << pos)
	case "E":
		reg.e |= (mask << pos)
	case "H":
		reg.h |= (mask << pos)
	case "L":
		reg.l |= (mask << pos)
	}
}

func (reg *Register) resBr(destination string, pos uint16) {
	var mask byte = 1
	switch destination {
	case "A":
		reg.a &^= (mask << pos)
	case "B":
		reg.b &^= (mask << pos)
	case "C":
		reg.c &^= (mask << pos)
	case "D":
		reg.d &^= (mask << pos)
	case "E":
		reg.e &^= (mask << pos)
	case "H":
		reg.h &^= (mask << pos)
	case "L":
		reg.l &^= (mask << pos)
	}
}

/* *************************************** */
/* Jumps                                   */
/* *************************************** */

func (reg *Register) jpnn(destination uint16) {
	reg.pc = destination
}

func (reg *Register) jpccnn(destination uint16, condition string) {
	switch condition {
	case "NZ":
		if !hasBit(uint16(reg.flags), 7) {
			reg.pc = destination
		}
	case "Z":
		if hasBit(uint16(reg.flags), 7) {
			reg.pc = destination
		}
	case "NC":
		if !hasBit(uint16(reg.flags), 4) {
			reg.pc = destination
		}
	case "C":
		if hasBit(uint16(reg.flags), 4) {
			reg.pc = destination
		}
	}
}

func (reg *Register) jpHL() {
	HL := (uint16(reg.h) << 8) + uint16(reg.l)
	reg.pc = HL
}

func (reg *Register) jrn(n uint16) {
	reg.pc += n
}

func (reg *Register) jrccn(n uint16, condition string) {
	switch condition {
	case "NZ":
		if !hasBit(uint16(reg.flags), 7) {
			reg.pc += n
		}
	case "Z":
		if hasBit(uint16(reg.flags), 7) {
			reg.pc += n
		}
	case "NC":
		if !hasBit(uint16(reg.flags), 4) {
			reg.pc += n
		}
	case "C":
		if hasBit(uint16(reg.flags), 4) {
			reg.pc += n
		}
	}
}

/* *************************************** */
/* Calls                                   */
/* *************************************** */

func (reg *Register) callnn(destination uint16, mem *Memory) {
	mem.writeWord(reg.sp, reg.pc+1)
	reg.sp += 2
	reg.pc = destination
}

func (reg *Register) callccnn(n uint16, condition string, mem *Memory) {
	switch condition {
	case "NZ":
		if !hasBit(uint16(reg.flags), 7) {
			mem.writeWord(reg.sp, reg.pc+1)
			reg.sp++
			reg.pc = n
		}
	case "Z":
		if hasBit(uint16(reg.flags), 7) {
			mem.writeWord(reg.sp, reg.pc+1)
			reg.sp++
			reg.pc += n
		}
	case "NC":
		if !hasBit(uint16(reg.flags), 4) {
			mem.writeWord(reg.sp, reg.pc+1)
			reg.sp++
			reg.pc += n
		}
	case "C":
		if hasBit(uint16(reg.flags), 4) {
			mem.writeWord(reg.sp, reg.pc+1)
			reg.sp++
			reg.pc += n
		}
	}
}

/* *************************************** */
/* Returns                                 */
/* *************************************** */

func (reg *Register) ret(mem *Memory) {
	address := mem.readWord(reg.sp)
	reg.pc = address
	reg.sp += 2
}

func (reg *Register) retcc(mem *Memory, condition string) {
	if condition == "NZ" && !hasBit(uint16(reg.flags), 7) {
		reg.ret(mem)
	} else if condition == "Z" && hasBit(uint16(reg.flags), 7) {
		reg.ret(mem)
	} else if condition == "NC" && !hasBit(uint16(reg.flags), 4) {
		reg.ret(mem)
	} else if condition == "C" && hasBit(uint16(reg.flags), 4) {
		reg.ret(mem)
	}
}

func (reg *Register) execute(opcode byte, mem *Memory) {
	switch opcode {
	case 0x00:
		//nop
	case 0x01:
		reg.ldnnn16(mem.readWord(reg.pc), "BC")
	case 0x02:
		reg.ldnA("BC", mem)
	case 0x03:
		reg.incnn("BC")
	case 0x04:
		reg.incn("B")
	case 0x05:
		reg.decn("B")
	case 0x06:
		value := mem.readByte(reg.pc)
		reg.ldnnn(value, "B")
	case 0x07:
		reg.rlcA()
	case 0x08:
		value := mem.readWord(reg.pc)
		reg.ldnnSP(value, mem)
	case 0x09:
		value := concatenateBytes(reg.b, reg.c)
		reg.addHLn(value)
	case 0x0a:
		reg.ldAn("BC", mem)
	case 0x0b:
		reg.decnn("BC")
	case 0x0c:
		reg.incn("C")
	case 0x0d:
		reg.decn("C")
	case 0x0e:
		value := mem.readByte(reg.pc)
		reg.ldnnn(value, "C")
	case 0x0f:
		reg.rrcA()
	case 0x10:
		//stop
	case 0x11:
		value := mem.readWord(reg.pc)
		reg.ldnnn16(value, "DE")
	case 0x12:
		reg.ldnA("DE", mem)
	case 0x13:
		reg.incnn("DE")
	case 0x14:
		reg.incn("D")
	case 0x15:
		reg.decn("D")
	case 0x16:
		value := mem.readByte(reg.pc)
		reg.ldnnn(value, "D")
	case 0x17:
		reg.rlA()
	case 0x18:
		value := mem.readWord(reg.pc)
		reg.jrn(value)
	case 0x19:
		value := mem.readWord(concatenateBytes(reg.d, reg.e))
		reg.addHLn(value)
	case 0x1a:
		reg.ldAn("DE", mem)
	case 0x1b:
		reg.decnn("DE")
	case 0x1c:
		reg.incn("E")
	case 0x1d:
		reg.decn("E")
	case 0x1e:
		value := mem.readByte(reg.pc)
		reg.ldnnn(value, "E")
	case 0x1f:
		reg.rrA()
	case 0x20:
		value := mem.readWord(reg.pc)
		reg.jrccn(value, "NZ")
	case 0x21:
		value := mem.readWord(reg.pc)
		reg.ldnnn16(value, "HL")
	case 0x22:
		reg.ldiHLA(mem)
	case 0x23:
		reg.incnn("HL")
	case 0x24:
		reg.incn("H")
	case 0x25:
		reg.decn("H")
	case 0x26:
		value := mem.readByte(reg.pc)
		reg.ldnnn(value, "H")
	case 0x27:
		reg.dAA()
	case 0x28:
		value := mem.readWord(reg.pc)
		reg.jrccn(value, "Z")
	case 0x29:
		value := mem.readWord(concatenateBytes(reg.h, reg.l))
		reg.addHLn(value)
	case 0x2a:
		reg.ldiAHL(mem)
	case 0x2b:
		reg.decnn("HL")
	case 0x2c:
		reg.incn("L")
	case 0x2d:
		reg.decn("L")
	case 0x2e:
		value := mem.readByte(reg.pc)
		reg.ldnnn(value, "L")
	case 0x2f:
		reg.cpl()
	case 0x30:
		value := mem.readWord(reg.pc)
		reg.jrccn(value, "NC")
	case 0x31:
		value := mem.readWord(reg.sp)
		reg.ldnnn16(value, "SP")
	case 0x32:
		reg.ldiHLA(mem)
	case 0x33:
		reg.incnn("SP")
	case 0x34:
		reg.incnn("HL")
	case 0x35:
		reg.decnn("HL")
	case 0x36:
		value := mem.readWord(reg.pc)
		reg.ldnnn16(value, "HL")
	case 0x37:
		reg.scf()
	case 0x38:
		value := mem.readWord(reg.pc)
		reg.jrccn(value, "C")
	case 0x39:
		value := mem.readWord(reg.sp)
		reg.addHLn(value)
	case 0x3a:
		reg.ldiAHL(mem)
	case 0x3b:
		reg.decnn("SP")
	case 0x3c:
		reg.incn("A")
	case 0x3d:
		reg.decn("A")
	case 0x3e:
		reg.ldAn("0x3e", mem)
	case 0x3f:
		reg.ccf()
	case 0x40:
		reg.ldr1r2("B", "B", mem)
	case 0x41:
		reg.ldr1r2("B", "C", mem)
	case 0x42:
		reg.ldr1r2("B", "D", mem)
	case 0x43:
		reg.ldr1r2("B", "E", mem)
	case 0x44:
		reg.ldr1r2("B", "H", mem)
	case 0x45:
		reg.ldr1r2("B", "L", mem)
	case 0x46:
		reg.ldr1r2("B", "HL", mem)
	case 0x47:
		reg.ldr1r2("B", "A", mem)
	case 0x48:
		reg.ldr1r2("C", "B", mem)
	case 0x49:
		reg.ldr1r2("C", "C", mem)
	case 0x4a:
		reg.ldr1r2("C", "D", mem)
	case 0x4b:
		reg.ldr1r2("C", "E", mem)
	case 0x4c:
		reg.ldr1r2("C", "H", mem)
	case 0x4d:
		reg.ldr1r2("C", "L", mem)
	case 0x4e:
		reg.ldr1r2("C", "HL", mem)
	case 0x4f:
		reg.ldr1r2("C", "A", mem)
	case 0x50:
		reg.ldr1r2("D", "B", mem)
	case 0x51:
		reg.ldr1r2("D", "C", mem)
	case 0x52:
		reg.ldr1r2("D", "D", mem)
	case 0x53:
		reg.ldr1r2("D", "E", mem)
	case 0x54:
		reg.ldr1r2("D", "H", mem)
	case 0x55:
		reg.ldr1r2("D", "L", mem)
	case 0x56:
		reg.ldr1r2("D", "HL", mem)
	case 0x57:
		reg.ldr1r2("D", "A", mem)
	case 0x58:
		reg.ldr1r2("E", "B", mem)
	case 0x59:
		reg.ldr1r2("E", "C", mem)
	case 0x5a:
		reg.ldr1r2("E", "D", mem)
	case 0x5b:
		reg.ldr1r2("E", "E", mem)
	case 0x5c:
		reg.ldr1r2("E", "H", mem)
	case 0x5d:
		reg.ldr1r2("E", "L", mem)
	case 0x5e:
		reg.ldr1r2("E", "HL", mem)
	case 0x5f:
		reg.ldr1r2("E", "A", mem)
	case 0x60:
		reg.ldr1r2("H", "B", mem)
	case 0x61:
		reg.ldr1r2("H", "C", mem)
	case 0x62:
		reg.ldr1r2("H", "D", mem)
	case 0x63:
		reg.ldr1r2("H", "E", mem)
	case 0x64:
		reg.ldr1r2("H", "H", mem)
	case 0x65:
		reg.ldr1r2("H", "L", mem)
	case 0x66:
		reg.ldr1r2("H", "HL", mem)
	case 0x67:
		reg.ldr1r2("H", "A", mem)
	case 0x68:
		reg.ldr1r2("L", "B", mem)
	case 0x69:
		reg.ldr1r2("L", "C", mem)
	case 0x6a:
		reg.ldr1r2("L", "D", mem)
	case 0x6b:
		reg.ldr1r2("L", "E", mem)
	case 0x6c:
		reg.ldr1r2("L", "H", mem)
	case 0x6d:
		reg.ldr1r2("L", "L", mem)
	case 0x6e:
		reg.ldr1r2("L", "HL", mem)
	case 0x6f:
		reg.ldr1r2("L", "A", mem)
	case 0x70:
		reg.ldr1r2("HL", "B", mem)
	case 0x71:
		reg.ldr1r2("HL", "C", mem)
	case 0x72:
		reg.ldr1r2("HL", "D", mem)
	case 0x73:
		reg.ldr1r2("HL", "E", mem)
	case 0x74:
		reg.ldr1r2("HL", "H", mem)
	case 0x75:
		reg.ldr1r2("HL", "L", mem)
	case 0x76:
		// halt
	case 0x77:
		reg.ldr1r2("HL", "A", mem)
	case 0x78:
		reg.ldr1r2("A", "B", mem)
	case 0x79:
		reg.ldr1r2("A", "C", mem)
	case 0x7a:
		reg.ldr1r2("A", "D", mem)
	case 0x7b:
		reg.ldr1r2("A", "E", mem)
	case 0x7c:
		reg.ldr1r2("A", "H", mem)
	case 0x7d:
		reg.ldr1r2("A", "L", mem)
	case 0x7e:
		reg.ldr1r2("A", "HL", mem)
	case 0x7f:
		reg.ldr1r2("A", "A", mem)
	case 0x80:
		value := mem.readByte(uint16(reg.b))
		reg.addAn(value)
	case 0x81:
		value := mem.readByte(uint16(reg.c))
		reg.addAn(value)
	case 0x82:
		value := mem.readByte(uint16(reg.d))
		reg.addAn(value)
	case 0x83:
		value := mem.readByte(uint16(reg.e))
		reg.addAn(value)
	case 0x84:
		value := mem.readByte(uint16(reg.h))
		reg.addAn(value)
	case 0x85:
		value := mem.readByte(uint16(reg.l))
		reg.addAn(value)
	case 0x86:
		value := mem.readByte(concatenateBytes(reg.h, reg.l))
		reg.addAn(value)
	case 0x87:
		value := mem.readByte(uint16(reg.a))
		reg.addAn(value)
	case 0x88:
		value := mem.readByte(uint16(reg.b))
		reg.addcAn(value)
	case 0x89:
		value := mem.readByte(uint16(reg.c))
		reg.addcAn(value)
	case 0x8a:
		value := mem.readByte(uint16(reg.d))
		reg.addcAn(value)
	case 0x8b:
		value := mem.readByte(uint16(reg.e))
		reg.addcAn(value)
	case 0x8c:
		value := mem.readByte(uint16(reg.h))
		reg.addcAn(value)
	case 0x8d:
		value := mem.readByte(uint16(reg.l))
		reg.addcAn(value)
	case 0x8e:
		value := mem.readByte(concatenateBytes(reg.h, reg.l))
		reg.addcAn(value)
	case 0x8f:
		value := mem.readByte(uint16(reg.a))
		reg.addcAn(value)
	case 0x90:
		value := mem.readByte(uint16(reg.b))
		reg.subn(value)
	case 0x91:
		value := mem.readByte(uint16(reg.c))
		reg.subn(value)
	case 0x92:
		value := mem.readByte(uint16(reg.d))
		reg.subn(value)
	case 0x93:
		value := mem.readByte(uint16(reg.e))
		reg.subn(value)
	case 0x94:
		value := mem.readByte(uint16(reg.h))
		reg.subn(value)
	case 0x95:
		value := mem.readByte(uint16(reg.l))
		reg.subn(value)
	case 0x96:
		value := mem.readByte(concatenateBytes(reg.h, reg.l))
		reg.subn(value)
	case 0x97:
		value := mem.readByte(uint16(reg.a))
		reg.subn(value)
	case 0x98:
		value := mem.readByte(uint16(reg.b))
		reg.sbcAn(value)
	case 0x99:
		value := mem.readByte(uint16(reg.c))
		reg.sbcAn(value)
	case 0x9a:
		value := mem.readByte(uint16(reg.d))
		reg.sbcAn(value)
	case 0x9b:
		value := mem.readByte(uint16(reg.e))
		reg.sbcAn(value)
	case 0x9c:
		value := mem.readByte(uint16(reg.h))
		reg.sbcAn(value)
	case 0x9d:
		value := mem.readByte(uint16(reg.l))
		reg.sbcAn(value)
	case 0x9e:
		value := mem.readByte(concatenateBytes(reg.h, reg.l))
		reg.sbcAn(value)
	case 0x9f:
		value := mem.readByte(uint16(reg.a))
		reg.sbcAn(value)
	case 0xa0:
		value := mem.readByte(uint16(reg.b))
		reg.andn(value)
	case 0xa1:
		value := mem.readByte(uint16(reg.c))
		reg.andn(value)
	case 0xa2:
		value := mem.readByte(uint16(reg.d))
		reg.andn(value)
	case 0xa3:
		value := mem.readByte(uint16(reg.e))
		reg.andn(value)
	case 0xa4:
		value := mem.readByte(uint16(reg.h))
		reg.andn(value)
	case 0xa5:
		value := mem.readByte(uint16(reg.l))
		reg.andn(value)
	case 0xa6:
		value := mem.readByte(concatenateBytes(reg.h, reg.l))
		reg.andn(value)
	case 0xa7:
		value := mem.readByte(uint16(reg.a))
		reg.andn(value)
	case 0xa8:
		value := mem.readByte(uint16(reg.b))
		reg.xorn(value)
	case 0xa9:
		value := mem.readByte(uint16(reg.c))
		reg.xorn(value)
	case 0xaa:
		value := mem.readByte(uint16(reg.d))
		reg.xorn(value)
	case 0xab:
		value := mem.readByte(uint16(reg.e))
		reg.xorn(value)
	case 0xac:
		value := mem.readByte(uint16(reg.h))
		reg.xorn(value)
	case 0xad:
		value := mem.readByte(uint16(reg.l))
		reg.xorn(value)
	case 0xae:
		value := mem.readByte(concatenateBytes(reg.h, reg.l))
		reg.xorn(value)
	case 0xaf:
		value := mem.readByte(uint16(reg.a))
		reg.xorn(value)
	case 0xb0:
		value := mem.readByte(uint16(reg.b))
		reg.orn(value)
	case 0xb1:
		value := mem.readByte(uint16(reg.c))
		reg.orn(value)
	case 0xb2:
		value := mem.readByte(uint16(reg.d))
		reg.orn(value)
	case 0xb3:
		value := mem.readByte(uint16(reg.e))
		reg.orn(value)
	case 0xb4:
		value := mem.readByte(uint16(reg.h))
		reg.orn(value)
	case 0xb5:
		value := mem.readByte(uint16(reg.l))
		reg.orn(value)
	case 0xb6:
		value := mem.readByte(concatenateBytes(reg.h, reg.l))
		reg.orn(value)
	case 0xb7:
		value := mem.readByte(uint16(reg.a))
		reg.orn(value)
	case 0xb8:
		value := mem.readByte(uint16(reg.b))
		reg.cpn(value)
	case 0xb9:
		value := mem.readByte(uint16(reg.c))
		reg.cpn(value)
	case 0xba:
		value := mem.readByte(uint16(reg.d))
		reg.cpn(value)
	case 0xbb:
		value := mem.readByte(uint16(reg.e))
		reg.cpn(value)
	case 0xbc:
		value := mem.readByte(uint16(reg.h))
		reg.cpn(value)
	case 0xbd:
		value := mem.readByte(uint16(reg.l))
		reg.cpn(value)
	case 0xbe:
		value := mem.readByte(concatenateBytes(reg.h, reg.l))
		reg.cpn(value)
	case 0xbf:
		value := mem.readByte(uint16(reg.a))
		reg.cpn(value)
	case 0xc0:
		reg.retcc(mem, "NZ")
	case 0xc1:
		reg.popnn("BC", mem)
	case 0xc2:
		value := mem.readWord(reg.pc)
		reg.jpccnn(value, "NZ")
	case 0xc3:
		value := mem.readWord(reg.pc)
		reg.jpnn(value)
	case 0xc4:
		value := mem.readWord(reg.pc)
		reg.callccnn(value, "NZ", mem)
	case 0xc5:
		reg.pushnn("BC", mem)
	case 0xc6:
		value := mem.readByte(reg.pc)
		reg.addAn(value)
	case 0xc7:
		//reset
	case 0xc8:
		reg.retcc(mem, "Z")
	case 0xc9:
		reg.ret(mem)
	case 0xca:
		value := mem.readWord(reg.pc)
		reg.jpccnn(value, "Z")
	case 0xcb:
		//cb prefix
	case 0xcc:
		value := mem.readWord(reg.pc)
		reg.callccnn(value, "Z", mem)
	case 0xcd:
		value := mem.readWord(reg.pc)
		reg.callnn(value, mem)
	case 0xce:
		value := mem.readByte(reg.pc)
		reg.addcAn(value)
	case 0xcf:
		//reset
	case 0xd0:
		reg.retcc(mem, "NC")
	case 0xd1:
		reg.popnn("DE", mem)
	case 0xd2:
		value := mem.readWord(reg.pc)
		reg.jpccnn(value, "NC")
	case 0xd3:
		//not used
	case 0xd4:
		value := mem.readWord(reg.pc)
		reg.callccnn(value, "NC", mem)
	case 0xd5:
		reg.pushnn("DE", mem)
	case 0xd6:
		value := mem.readByte(reg.pc)
		reg.subn(value)
	case 0xd7:
		//reset
	case 0xd8:
		reg.retcc(mem, "C")
	case 0xd9:
		//return and interrupt (not implemented)
	case 0xda:
		value := mem.readWord(reg.pc)
		reg.jpccnn(value, "C")
	case 0xdb:
		//not used
	case 0xdc:
		value := mem.readWord(reg.pc)
		reg.callccnn(value, "C", mem)
	case 0xdd:
		//not used
	case 0xde:
		value := mem.readByte(reg.pc)
		reg.sbcAn(value)
	case 0xdf:
		//reset
	case 0xe0:
		value := mem.readByte(reg.pc)
		reg.ldhnA(value, mem)
	case 0xe1:
		reg.popnn("HL", mem)
	case 0xe2:
		reg.ldCA(mem)
	case 0xe3:
		//not used
	case 0xe4:
		//not used
	case 0xe5:
		reg.pushnn("HL", mem)
	case 0xe6:
		value := mem.readByte(reg.pc)
		reg.andn(value)
	case 0xe7:
		//reset
	case 0xe8:
		value := mem.readWord(reg.pc)
		reg.addSPn(value)
	case 0xe9:
		reg.jpHL()
	case 0xea:
		reg.ldnA("HL", mem)
	case 0xeb:
		//not used
	case 0xec:
		//not used
	case 0xed:
		//not used
	case 0xee:
		value := mem.readByte(reg.pc)
		reg.xorn(value)
	case 0xef:
		//Reset
	case 0xf0:
		reg.ldAn("0xf0", mem)
	case 0xf1:
		reg.popnn("AF", mem)
	case 0xf2:
		reg.ldAn("C", mem)
	case 0xf3:
		//not implemented
	case 0xf4:
		//not used
	case 0xf5:
		reg.pushnn("AF", mem)
	case 0xf6:
		value := mem.readByte(reg.pc)
		reg.orn(value)
	case 0xf7:
		//reset
	case 0xf8:
		value := mem.readByte(reg.pc)
		reg.ldHLSPn(value)
	case 0xf9:
		reg.ldSPHL()
	case 0xfa:
		//implement this properly
		reg.ldAn("OxfA", mem)
	case 0xfb:
		//not implemented
	case 0xfc:
		//not used
	case 0xfd:
		//not used
	case 0xfe:
		value := mem.readByte(reg.pc)
		reg.cpn(value)
	case 0xff:
		//reset
	}
}

func (reg *Register) executeCb(opcode byte, mem *Memory) {
	switch opcode {
	case 0x00:
		reg.rlcn("B")
	case 0x01:
		reg.rlcn("C")
	case 0x02:
		reg.rlcn("D")
	case 0x03:
		reg.rlcn("E")
	case 0x04:
		reg.rlcn("H")
	case 0x05:
		reg.rlcn("L")
	case 0x06:
		reg.rlcn("HL")
	case 0x07:
		reg.rlcA()
	case 0x08:
		reg.rrcn("B")
	case 0x09:
		reg.rrcn("C")
	case 0x0a:
		reg.rrcn("D")
	case 0x0b:
		reg.rrcn("E")
	case 0x0c:
		reg.rrcn("H")
	case 0x0d:
		reg.rrcn("L")
	case 0x0e:
		reg.rrcn("HL")
	case 0x0f:
		reg.rrcA()
	}
}
