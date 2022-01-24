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

func (reg *Register) ldr1r2(destination string, source string) {
	reg_map := map[string]*byte{
		"A": &(reg.a),
		"B": &(reg.b),
		"C": &(reg.c),
		"D": &(reg.d),
		"E": &(reg.e),
		"H": &(reg.h),
		"L": &(reg.l),
	}
	*(reg_map[destination]) = *(reg_map[source])
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
		value = ((value & 0x0F) << 4) | ((value & 0xF0) >> 4)
		reg.a = value
	case "B":
		value := reg.b
		value = ((value & 0x0F) << 4) | ((value & 0xF0) >> 4)
		reg.b = value
	case "C":
		value := reg.c
		value = ((value & 0x0F) << 4) | ((value & 0xF0) >> 4)
		reg.c = value
	case "D":
		value := reg.d
		value = ((value & 0x0F) << 4) | ((value & 0xF0) >> 4)
		reg.d = value
	case "E":
		value := reg.e
		value = ((value & 0x0F) << 4) | ((value & 0xF0) >> 4)
		reg.e = value
	case "H":
		value := reg.h
		value = ((value & 0x0F) << 4) | ((value & 0xF0) >> 4)
		reg.h = value
	case "L":
		value := reg.l
		value = ((value & 0x0F) << 4) | ((value & 0xF0) >> 4)
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
			value = (value - 0x06) & 0xFF
		}
		if hasBit(uint16(reg.flags), 4) {
			value -= 0x60
		}
	} else {
		if hasBit(uint16(reg.flags), 5) || ((value & 0xF) > 9) {
			value += 0x06
		}
		if hasBit(uint16(reg.flags), 4) || value > 0x9F {
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
	// haf carry flag
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
func (reg *Register) rlcA() {
	value := uint16(reg.a) << 1
	if hasBit(uint16(reg.flags), 4) {
		var pos uint16 = 7
		value |= pos
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

func (reg *Register) rlA() {
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
	reg.a = byte(value)
	//carry flag
	if (reg.a & 0x80) != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
}

func (reg *Register) rrcA() {
	value := uint16(reg.a) >> 1
	if hasBit(uint16(reg.flags), 4) {
		var pos uint16 = 1
		value |= pos
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

func (reg *Register) rrA() {
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
	reg.a = byte(value)
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
	value := uint16(register) >> 1
	if hasBit(uint16(reg.flags), 4) {
		var pos uint16 = 1
		value |= (pos << 7)
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
	if hasBit(uint16(reg.flags), 4) {
		var pos uint16 = 1
		value |= pos
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
	reg.sp++
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

func (reg *Register) execute(opcode byte, mem *Memory) {
	switch opcode {
	case 0x00:
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
	}
}
