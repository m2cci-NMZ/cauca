package main

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
func (reg *Register) ldAn(value byte) {
	reg.a = value
}

func (reg *Register) ldAC() {
	reg.a = io[reg.c]
}

func (reg *Register) ldCA() {
	io[reg.c] = reg.a
}

func (reg *Register) lddAHL() {
	address := (uint16(reg.h) << 8) + uint16(reg.l)
	reg.a = readByte(address)
	address--
	reg.h = byte(address >> 8)
	reg.l = byte(address)
}

func (reg *Register) lddHLA() {
	address := (uint16(reg.h) << 8) + uint16(reg.l)
	io[address] = reg.a
	address--
	reg.h = byte(address >> 8)
	reg.l = byte(address)
}

func (reg *Register) ldiAHL() {
	address := (uint16(reg.h) << 8) + uint16(reg.l)
	reg.a = readByte(address)
	address++
	reg.h = byte(address >> 8)
	reg.l = byte(address)
}

func (reg *Register) ldiHLA() {
	address := (uint16(reg.h) << 8) + uint16(reg.l)
	io[address] = reg.a
	address++
	reg.h = byte(address >> 8)
	reg.l = byte(address)
}

func (reg *Register) ldhnA(value byte) {
	io[value] = reg.a
}

func (reg *Register) ldhAn(value byte) {
	reg.a = io[value]
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

func (reg *Register) addCarry(value byte) {
	if hasBit(uint16(reg.flags), 4) {
		value++
	}
	reg.addAn(value)
}
