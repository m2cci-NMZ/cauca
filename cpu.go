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
	clock int
}

var ticks [256]int = [256]int{
	2, 6, 4, 4, 2, 2, 4, 4, 10, 4, 4, 4, 2, 2, 4, 4,
	2, 6, 4, 4, 2, 2, 4, 4, 4, 4, 4, 4, 2, 2, 4, 4,
	0, 6, 4, 4, 2, 2, 4, 2, 0, 4, 4, 4, 2, 2, 4, 2,
	4, 6, 4, 4, 6, 6, 6, 2, 0, 4, 4, 4, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	4, 4, 4, 4, 4, 4, 2, 4, 2, 2, 2, 2, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	0, 6, 0, 6, 0, 8, 4, 8, 0, 2, 0, 0, 0, 6, 4, 8,
	0, 6, 0, 0, 0, 8, 4, 8, 0, 8, 0, 0, 0, 0, 4, 8,
	6, 6, 4, 0, 0, 8, 4, 8, 8, 2, 8, 0, 0, 0, 4, 8,
	6, 6, 4, 2, 0, 8, 4, 8, 6, 4, 8, 2, 0, 0, 4, 8}

/* *************************************** */
/* Helper functions                        */
/* *************************************** */

// Takes a and b and concatenates them to return a 16 bit word
func concatenateBytes(a byte, b byte) uint16 {
	result := (uint16(b) << 8) + uint16(a)
	return result
}

// Separates value and returns the two bytes.
func separateWord(value uint16) (byte, byte) {
	b := byte(value >> 8)
	a := byte(value)
	return a, b
}

// Returns the HL register, which is a 16 bit combination from registers H and L
func (reg *Register) getHLregister() uint16 {
	return concatenateBytes(reg.l, reg.h)
}

// Sets registers H and L by separating value
func (reg *Register) setHLregisters(value uint16) {
	a, b := separateWord(value)
	reg.h = b
	reg.l = a
}

/* *************************************** */
/* Flags setting function                  */
/* *************************************** */

// Sets bit at position of F register to value
func (register *Register) setRegisterFlag(value bool, position byte) {
	if value {
		register.flags |= (1 << position)
	} else {
		register.flags &^= (1 << position)
	}
}

// Returns the value of bit pos in n
func hasBit(n uint16, pos uint16) bool {
	val := n & (1 << pos)
	return (val > 0)
}

/* *************************************** */
/* 8 bit loads                             */
/* *************************************** */

// Loads value in register destination
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

// Loads register source in register destination. For 16 bit registers, read memory instead.
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
		*(reg_map[destination]) = mem.readByte(reg.getHLregister())
	} else if destination == "HL" {
		if len(source) == 1 {
			value := mem.readWord(uint16(*(reg_map[source])))
			reg.setHLregisters(value)
		} else {
			value := mem.readWord(reg.pc)
			reg.setHLregisters(value)
		}
	} else {
		*(reg_map[destination]) = *(reg_map[source])
	}
}

// Loads value at memory address source and put it in register A
func (reg *Register) ldAn(source string, mem *Memory) {
	reg_map := map[string]uint16{
		"A":  uint16(reg.a),
		"B":  uint16(reg.b),
		"C":  uint16(reg.c),
		"D":  uint16(reg.d),
		"E":  uint16(reg.e),
		"H":  uint16(reg.h),
		"L":  uint16(reg.l),
		"HL": reg.getHLregister(),
		"BC": concatenateBytes(reg.c, reg.b),
		"DE": concatenateBytes(reg.e, reg.d),
		"PC": reg.pc,
	}
	reg.a = mem.readByte(reg_map[source])
}

// Loads value od register source in register A
func (reg *Register) ldnA(source string, mem *Memory) {
	reg_map := map[string]byte{
		"A": reg.a,
		"B": reg.b,
		"C": reg.c,
		"D": reg.d,
		"E": reg.e,
		"H": reg.h,
		"L": reg.l,
	}
	if len(source) == 1 {
		reg.a = reg_map[source]
	} else {
		switch source {
		case "BC":
			mem.writeByte(concatenateBytes(reg.c, reg.b), reg.a)
		case "DE":
			mem.writeByte(concatenateBytes(reg.e, reg.d), reg.a)
		case "HL":
			mem.writeByte(reg.getHLregister(), reg.a)
		case "PC":
			mem.writeByte(reg.pc, reg.a)
		}
	}
}

// Load value in the io memory bank at address in register C on register A
func (reg *Register) ldAC(mem *Memory) {
	reg.a = mem.io[reg.c]
}

// Load value in register A in io memory bank at address in register C
func (reg *Register) ldCA(mem *Memory) {
	mem.io[reg.c] = reg.a
}

// Load value at address HL in register A and decrement HL
func (reg *Register) lddAHL(mem *Memory) {
	address := reg.getHLregister()
	reg.a = mem.readByte(address)
	address--
	reg.setHLregisters(address)
}

// Loads value in register A in address HL and decrement HL
func (reg *Register) lddHLA(mem *Memory) {
	address := reg.getHLregister()
	mem.writeByte(address, reg.a)
	address--
	reg.setHLregisters(address)
}

// Load value at address HL in register A and increment HL
func (reg *Register) ldiAHL(mem *Memory) {
	address := reg.getHLregister()
	reg.a = mem.readByte(address)
	address++
	reg.setHLregisters(address)
}

// Loads value in register A in address HL and increment HL
func (reg *Register) ldiHLA(mem *Memory) {
	address := reg.getHLregister()
	mem.writeByte(address, reg.a)
	address++
	reg.setHLregisters(address)
}

// Load value in register A in io memory bank at address value
func (reg *Register) ldhnA(value byte, mem *Memory) {
	if value == 0 {
		if reg.a == 0x20 {
			mem.io[0] = 0xef
		} else if reg.a == 0x10 {
			mem.io[0] = 0xdf
		} else {
			mem.io[value] = reg.a
		}
	} else {
		mem.io[value] = reg.a
	}
}

// Load value in io memory bank at address value in register A
func (reg *Register) ldhAn(value byte, mem *Memory) {
	reg.a = mem.io[value]
}

/* *************************************** */
/* 16 bit loads                            */
/* *************************************** */

// Loads value in 16 bit register destination
func (reg *Register) ldnnn16(value uint16, destination string) {
	r2, r1 := separateWord(value)
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

// Load HL register in PC
func (reg *Register) ldSPHL() {
	value := reg.getHLregister()
	reg.sp = value
}

// Load SP+value into HL
func (reg *Register) ldHLSPn(value byte) {
	result := uint16(value) + reg.sp
	reg.setHLregisters(uint16(value))
	// reset Z flag
	reg.setRegisterFlag(false, 7)
	// reset N flag
	reg.setRegisterFlag(false, 6)
	// set H flag
	if (reg.sp&0x000f + uint16(value&0x0f)) > 0x0F {
		reg.setRegisterFlag(true, 5)
	} else {
		reg.setRegisterFlag(false, 5)
	}
	//set C flag
	if (result & 0xff00) != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
}

// Load SP at value address
func (reg *Register) ldnnSP(value uint16, mem *Memory) {
	mem.writeWord(value, reg.sp)
}

// Push pair of registers on top of stack and dercrement stack twice
func (reg *Register) pushnn(registers string, mem *Memory) {
	var value uint16
	switch registers {
	case "AF":
		value = concatenateBytes(reg.flags, reg.a)
	case "BC":
		value = concatenateBytes(reg.c, reg.b)

	case "DE":
		value = concatenateBytes(reg.e, reg.d)

	case "HL":
		value = reg.getHLregister()

	}
	mem.writeWord(reg.sp, value)
	reg.sp = reg.sp - 2
}

// Pop 16 bits on top of the stack and put in pair of registers and increment SP twice
func (reg *Register) popnn(registers string, mem *Memory) {
	reg.sp += 2
	r2 := mem.readByte(reg.sp)
	r1 := mem.readByte(reg.sp + 1)
	switch registers {
	case "AF":
		reg.a = r1
		reg.flags = r2
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
}

/* *************************************** */
/* 8 bit ALU                               */
/* *************************************** */

// Add value to register A
func (reg *Register) addAn(value byte) {
	result := uint16(reg.a) + uint16(value)
	// carry flag
	if (result & 0xFF00) != 0 {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	// zero flag
	if result == 0 {
		reg.setRegisterFlag(true, 7)
	} else {
		reg.setRegisterFlag(false, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	if (reg.a&0x0F + value&0x0F) > 0x0F {
		reg.setRegisterFlag(true, 5)
	} else {
		reg.setRegisterFlag(false, 5)
	}
	reg.a = byte(result & 0xFF)
}

// Add value to register A and adjust carry
func (reg *Register) addcAn(value byte) {
	if hasBit(uint16(reg.flags), 4) {
		value++
	}
	reg.addAn(value)
}

// Subtract value from register A
func (reg *Register) subn(value byte) {
	result := reg.a - value
	// negative flag
	reg.setRegisterFlag(true, 6)
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
	if reg.a < value {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	reg.a = result
}

// Subtract value from register A and carry correct
func (reg *Register) sbcAn(value byte) {
	if hasBit(uint16(reg.flags), 4) {
		value--
	}
	reg.subn(value)
}

// Perform bitwise AND of value with register A
func (reg *Register) andn(value byte) {
	result := reg.a & value
	// zero flag
	if result == 0 {
		reg.setRegisterFlag(true, 7)
	} else {
		reg.setRegisterFlag(false, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(true, 5)
	// carry flag
	reg.setRegisterFlag(false, 4)
	reg.a = result
}

// Perform bitwise OR of value with register A
func (reg *Register) orn(value byte) {
	result := reg.a | value
	// zero flag
	if result == 0 {
		reg.setRegisterFlag(true, 7)
	} else {
		reg.setRegisterFlag(false, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	// carry flag
	reg.setRegisterFlag(false, 4)
	reg.a = result
}

// Perform bitwise XOR of value with register A
func (reg *Register) xorn(value byte) {
	result := reg.a ^ value
	// zero flag
	if result == 0 {
		reg.setRegisterFlag(true, 7)
	} else {
		reg.setRegisterFlag(false, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	// carry flag
	reg.setRegisterFlag(false, 4)
	reg.a = result
}

// Compare register A with value and set flags
func (reg *Register) cpn(value byte) {
	tmp := reg.a
	reg.subn(value)
	reg.a = tmp
}

// Increment register
func (reg *Register) incn(register string) {
	var result uint16
	if len(register) == 1 {
		reg_map := map[string]*byte{
			"A": &reg.a,
			"B": &reg.b,
			"C": &reg.c,
			"D": &reg.d,
			"E": &reg.e,
			"H": &reg.h,
			"L": &reg.l,
		}
		*reg_map[register]++
		result = uint16(*reg_map[register])
	} else {
		result = reg.getHLregister()
		result++
		reg.setHLregisters(result)
	}
	// zero flag
	if result == 0 {
		reg.setRegisterFlag(true, 7)
	} else {
		reg.setRegisterFlag(false, 7)
	}
	// N flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	if (((byte(result) - 1) & 0x0f) + 1) > 0x0f {
		reg.setRegisterFlag(true, 5)
	} else {
		reg.setRegisterFlag(false, 5)
	}
}

// Decrement register
func (reg *Register) decn(register string) {
	var result uint16
	if len(register) == 1 {
		reg_map := map[string]*byte{
			"A": &reg.a,
			"B": &reg.b,
			"C": &reg.c,
			"D": &reg.d,
			"E": &reg.e,
			"H": &reg.h,
			"L": &reg.l,
		}
		*reg_map[register]--
		result = uint16(*reg_map[register])
	} else {
		result = reg.getHLregister()
		result--
		reg.setHLregisters(result)
	}
	// zero flag
	if result == 0 {
		reg.setRegisterFlag(true, 7)
	} else {
		reg.setRegisterFlag(false, 7)
	}
	// N flag
	reg.setRegisterFlag(true, 6)
	// half carry flag
	if (((byte(result) + 1) & 0x0f) - 1) > 0x0f {
		reg.setRegisterFlag(true, 5)
	} else {
		reg.setRegisterFlag(false, 5)
	}
}

/* *************************************** */
/* 16 bit ALU                               */
/* *************************************** */

// Add value to register HL
func (reg *Register) addHLn(value uint16) {
	var result uint32
	HL := reg.getHLregister()
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
	reg.setHLregisters(uint16(result))
}

// Add value to register SP
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

// Increment register
func (reg *Register) incnn(register string) {
	switch register {
	case "BC":
		result := concatenateBytes(reg.c, reg.b) + 1
		reg.c, reg.b = separateWord(result)
	case "DE":
		result := concatenateBytes(reg.e, reg.d) + 1
		reg.e, reg.d = separateWord(result)
	case "HL":
		result := concatenateBytes(reg.l, reg.h) + 1
		reg.l, reg.h = separateWord(result)
	case "SP":
		reg.sp++
	}
}

// Decrement register
func (reg *Register) decnn(register string) {
	switch register {
	case "BC":
		result := concatenateBytes(reg.c, reg.b) - 1
		reg.c, reg.b = separateWord(result)
	case "DE":
		result := concatenateBytes(reg.e, reg.d) - 1
		reg.e, reg.d = separateWord(result)
	case "HL":
		result := concatenateBytes(reg.l, reg.h) - 1
		reg.l, reg.h = separateWord(result)
	case "SP":
		reg.sp--
	}
}

/* *************************************** */
/* misc                                    */
/* *************************************** */

// Swap upper and lower nibbles of register
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
	case "HL":
		h := reg.h
		l := reg.l
		reg.h = l
		reg.l = h
		if concatenateBytes(reg.l, reg.h) != 0 {
			value = 1
		}
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

// BCD adjust register A
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

// Complement register A
func (reg *Register) cpl() {
	reg.a = ^reg.a
	// negative flag
	reg.setRegisterFlag(true, 6)
	// haf carry flag
	reg.setRegisterFlag(true, 5)
}

// Complement carry flag
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

// Set carry flag
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

// Rotate register A left
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

// Rotate register A left through carry flag
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

// Rotate register A right
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

// Rotate register A right through carry flag
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

// Rotate register destination left
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

// Rotate register destination left thorugh carry
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

// Rotate register destination right through carry
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

// Rotate register destination right.
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

// Shift register destination left
func (reg *Register) slan(destination string) {
	reg_map := map[string]*byte{
		"A": &(reg.a),
		"B": &(reg.b),
		"C": &(reg.c),
		"D": &(reg.d),
		"E": &(reg.e),
		"H": &(reg.h),
		"L": &(reg.l),
	}
	var register uint16
	if destination == "HL" {
		register = reg.getHLregister()
	} else {
		register = uint16(*reg_map[destination])
	}
	//carry flag
	if hasBit(register, 0) {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	value := register << 1
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
	if destination == "HL" {
		reg.setHLregisters(register)
	} else {
		*reg_map[destination] = byte(register)
	}
}

// Shift register destination right
func (reg *Register) sran(destination string) {
	reg_map := map[string]*byte{
		"A": &(reg.a),
		"B": &(reg.b),
		"C": &(reg.c),
		"D": &(reg.d),
		"E": &(reg.e),
		"H": &(reg.h),
		"L": &(reg.l),
	}
	var register uint16
	if destination == "HL" {
		register = reg.getHLregister()
	} else {
		register = uint16(*reg_map[destination])
	}
	msb := hasBit(register, 7)
	//carry flag
	if hasBit(register, 7) {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	value := register >> 1
	if msb {
		register |= (1 << 7)
	} else {
		register &^= (1 << 7)
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
	if destination == "HL" {
		reg.setHLregisters(register)
	} else {
		*reg_map[destination] = byte(register)
	}
}

// Shift register destination right
func (reg *Register) srln(destination string) {
	reg_map := map[string]*byte{
		"A": &(reg.a),
		"B": &(reg.b),
		"C": &(reg.c),
		"D": &(reg.d),
		"E": &(reg.e),
		"H": &(reg.h),
		"L": &(reg.l),
	}
	var register uint16
	if destination == "HL" {
		register = reg.getHLregister()
	} else {
		register = uint16(*reg_map[destination])
	}
	//carry flag
	if hasBit(register, 7) {
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
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
	if destination == "HL" {
		reg.setHLregisters(register)
	} else {
		*reg_map[destination] = byte(register)
	}
}

/* *************************************** */
/* Bit opcodes                             */
/* *************************************** */

// Set zero flag to register destination bit at position pos
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

// Set bit at position pos in register destination
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

// reset bit at position pos in register destination
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

//Reset: push address to stack and jump to destination
func (reg *Register) rst(destination uint16, mem *Memory) {
	//hotfix: rst has no arguments so next instruction is at pc, while it is at pc + 2 for call
	reg.pc -= 2
	reg.callnn(destination, mem)
}

// Jump at address destination
func (reg *Register) jpnn(destination uint16) {
	reg.pc = destination
}

// Jump at address destination if NZ flags
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

// Jump at address HL
func (reg *Register) jpHL() {
	HL := reg.getHLregister()
	reg.pc = HL
}

// Add n to PC
func (reg *Register) jrn(n byte) {
	var offset uint16 = uint16(int8(n))
	reg.pc += offset
}

// Add n to PC if NZ flags
func (reg *Register) jrccn(n byte, condition string) {
	var offset uint16 = uint16(int8(n))
	switch condition {
	case "NZ":
		if !hasBit(uint16(reg.flags), 7) {
			reg.pc += offset
		}
	case "Z":
		if hasBit(uint16(reg.flags), 7) {
			reg.pc += offset
		}
	case "NC":
		if !hasBit(uint16(reg.flags), 4) {
			reg.pc += offset
		}
	case "C":
		if hasBit(uint16(reg.flags), 4) {
			reg.pc += offset
		}
	}
}

/* *************************************** */
/* Calls                                   */
/* *************************************** */

// Push address of next instruction on top of stack and jump to destination
func (reg *Register) callnn(destination uint16, mem *Memory) {
	reg.pc += 2
	mem.writeWord(reg.sp, reg.pc)
	reg.sp -= 2
	reg.pc = destination
}

// Push address of next instruction on top of stack and jump to destination if condition
func (reg *Register) callccnn(n uint16, condition string, mem *Memory) {
	switch condition {
	case "NZ":
		if !hasBit(uint16(reg.flags), 7) {
			reg.callnn(n, mem)
		}
	case "Z":
		if hasBit(uint16(reg.flags), 7) {
			reg.callnn(n, mem)
		}
	case "NC":
		if !hasBit(uint16(reg.flags), 4) {
			reg.callnn(n, mem)
		}
	case "C":
		if hasBit(uint16(reg.flags), 4) {
			reg.callnn(n, mem)
		}
	}
}

/* *************************************** */
/* Returns                                 */
/* *************************************** */

// Pop value from stack and jump
func (reg *Register) ret(mem *Memory) {
	reg.sp += 2
	address := mem.readWord(reg.sp)
	reg.pc = address
}

// Pop value from stack and jump if condition
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

// Execute opcode
func (reg *Register) execute(opcode byte, mem *Memory) {
	reg.pc++
	reg.clock += ticks[opcode]
	switch opcode {
	case 0x00:
		//nop
	case 0x01:
		value := mem.readWord(reg.pc)
		reg.ldnnn16(value, "BC")
		reg.pc += 2
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
		reg.pc++
	case 0x07:
		reg.rlcA()
	case 0x08:
		value := mem.readWord(reg.pc)
		reg.ldnnSP(value, mem)
		reg.pc += 2
	case 0x09:
		value := concatenateBytes(reg.c, reg.b)
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
		reg.pc++
	case 0x0f:
		reg.rrcA()
	case 0x10:
		//stop
		reg.pc++
	case 0x11:
		reg.ldnnn16(reg.pc, "DE")
		reg.pc += 2
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
		reg.pc++
	case 0x17:
		reg.rlA()
	case 0x18:
		value := mem.readByte(reg.pc)
		reg.jrn(value)
		reg.pc++
	case 0x19:
		value := mem.readWord(concatenateBytes(reg.e, reg.d))
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
		reg.pc++
	case 0x1f:
		reg.rrA()
	case 0x20:
		value := mem.readByte(reg.pc)
		reg.jrccn(value, "NZ")
		reg.pc++
	case 0x21:
		value := mem.readWord(reg.pc)
		reg.ldnnn16(value, "HL")
		reg.pc += 2
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
		reg.pc++
	case 0x27:
		reg.dAA()
	case 0x28:
		value := mem.readByte(reg.pc)
		reg.jrccn(value, "Z")
		reg.pc++
	case 0x29:
		value := mem.readWord(concatenateBytes(reg.l, reg.h))
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
		reg.pc++
	case 0x2f:
		reg.cpl()
	case 0x30:
		value := mem.readByte(reg.pc)
		reg.jrccn(value, "NC")
		reg.pc++
	case 0x31:
		value := mem.readWord(reg.pc)
		reg.ldnnn16(value, "SP")
		reg.pc += 2
	case 0x32:
		reg.lddHLA(mem)
	case 0x33:
		reg.incnn("SP")
	case 0x34:
		reg.incnn("HL")
	case 0x35:
		reg.decnn("HL")
	case 0x36:
		reg.ldnnn16(reg.pc, "HL")
		reg.pc++
	case 0x37:
		reg.scf()
	case 0x38:
		value := mem.readByte(reg.pc)
		reg.jrccn(value, "C")
		reg.pc++
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
		reg.ldAn("PC", mem)
		reg.pc++
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
		value := reg.b
		reg.addAn(value)
	case 0x81:
		value := reg.c
		reg.addAn(value)
	case 0x82:
		value := reg.d
		reg.addAn(value)
	case 0x83:
		value := reg.e
		reg.addAn(value)
	case 0x84:
		value := reg.h
		reg.addAn(value)
	case 0x85:
		value := reg.l
		reg.addAn(value)
	case 0x86:
		value := mem.readByte(concatenateBytes(reg.l, reg.h))
		reg.addAn(value)
	case 0x87:
		value := reg.a
		reg.addAn(value)
	case 0x88:
		value := reg.b
		reg.addcAn(value)
	case 0x89:
		value := reg.c
		reg.addcAn(value)
	case 0x8a:
		value := reg.d
		reg.addcAn(value)
	case 0x8b:
		value := reg.e
		reg.addcAn(value)
	case 0x8c:
		value := reg.h
		reg.addcAn(value)
	case 0x8d:
		value := reg.l
		reg.addcAn(value)
	case 0x8e:
		value := mem.readByte(concatenateBytes(reg.l, reg.h))
		reg.addcAn(value)
	case 0x8f:
		value := reg.a
		reg.addcAn(value)
	case 0x90:
		value := reg.b
		reg.subn(value)
	case 0x91:
		value := reg.c
		reg.subn(value)
	case 0x92:
		value := reg.d
		reg.subn(value)
	case 0x93:
		value := reg.e
		reg.subn(value)
	case 0x94:
		value := reg.h
		reg.subn(value)
	case 0x95:
		value := reg.l
		reg.subn(value)
	case 0x96:
		value := mem.readByte(concatenateBytes(reg.l, reg.h))
		reg.subn(value)
	case 0x97:
		value := reg.a
		reg.subn(value)
	case 0x98:
		value := reg.b
		reg.sbcAn(value)
	case 0x99:
		value := reg.c
		reg.sbcAn(value)
	case 0x9a:
		value := reg.d
		reg.sbcAn(value)
	case 0x9b:
		value := reg.e
		reg.sbcAn(value)
	case 0x9c:
		value := reg.h
		reg.sbcAn(value)
	case 0x9d:
		value := reg.l
		reg.sbcAn(value)
	case 0x9e:
		value := mem.readByte(concatenateBytes(reg.l, reg.h))
		reg.sbcAn(value)
	case 0x9f:
		value := reg.a
		reg.sbcAn(value)
	case 0xa0:
		value := reg.b
		reg.andn(value)
	case 0xa1:
		value := reg.c
		reg.andn(value)
	case 0xa2:
		value := reg.d
		reg.andn(value)
	case 0xa3:
		value := reg.e
		reg.andn(value)
	case 0xa4:
		value := reg.h
		reg.andn(value)
	case 0xa5:
		value := reg.l
		reg.andn(value)
	case 0xa6:
		value := mem.readByte(concatenateBytes(reg.l, reg.h))
		reg.andn(value)
	case 0xa7:
		value := reg.a
		reg.andn(value)
	case 0xa8:
		value := reg.b
		reg.xorn(value)
	case 0xa9:
		value := reg.c
		reg.xorn(value)
	case 0xaa:
		value := reg.d
		reg.xorn(value)
	case 0xab:
		value := reg.e
		reg.xorn(value)
	case 0xac:
		value := reg.h
		reg.xorn(value)
	case 0xad:
		value := reg.l
		reg.xorn(value)
	case 0xae:
		value := mem.readByte(concatenateBytes(reg.l, reg.h))
		reg.xorn(value)
	case 0xaf:
		value := reg.a
		reg.xorn(value)
	case 0xb0:
		value := reg.b
		reg.orn(value)
	case 0xb1:
		value := reg.c
		reg.orn(value)
	case 0xb2:
		value := reg.d
		reg.orn(value)
	case 0xb3:
		value := reg.e
		reg.orn(value)
	case 0xb4:
		value := reg.h
		reg.orn(value)
	case 0xb5:
		value := reg.l
		reg.orn(value)
	case 0xb6:
		value := mem.readByte(concatenateBytes(reg.l, reg.h))
		reg.orn(value)
	case 0xb7:
		value := reg.a
		reg.orn(value)
	case 0xb8:
		value := reg.b
		reg.cpn(value)
	case 0xb9:
		value := reg.c
		reg.cpn(value)
	case 0xba:
		value := reg.d
		reg.cpn(value)
	case 0xbb:
		value := reg.e
		reg.cpn(value)
	case 0xbc:
		value := reg.h
		reg.cpn(value)
	case 0xbd:
		value := reg.l
		reg.cpn(value)
	case 0xbe:
		value := mem.readByte(concatenateBytes(reg.l, reg.h))
		reg.cpn(value)
	case 0xbf:
		value := reg.a
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
		reg.pc++
	case 0xc7:
		reg.rst(0, mem)
	case 0xc8:
		reg.retcc(mem, "Z")
	case 0xc9:
		reg.ret(mem)
	case 0xca:
		value := mem.readWord(reg.pc)
		reg.jpccnn(value, "Z")
	case 0xcb:
		opcode = mem.readByte(reg.pc)
		reg.executeCb(opcode)
	case 0xcc:
		value := mem.readWord(reg.pc)
		reg.callccnn(value, "Z", mem)
	case 0xcd:
		value := mem.readWord(reg.pc)
		reg.callnn(value, mem)
	case 0xce:
		value := mem.readByte(reg.pc)
		reg.addcAn(value)
		reg.pc++
	case 0xcf:
		reg.rst(0x08, mem)
	case 0xd0:
		reg.retcc(mem, "NC")
	case 0xd1:
		reg.popnn("DE", mem)
	case 0xd2:
		value := mem.readWord(reg.pc)
		reg.jpccnn(value, "NC")
		//reg.pc += 2
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
		reg.pc++
	case 0xd7:
		reg.rst(0x10, mem)
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
		reg.pc++
	case 0xdf:
		reg.rst(0x18, mem)
	case 0xe0:
		value := mem.readByte(reg.pc)
		reg.ldhnA(value, mem)
		reg.pc++
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
		reg.pc++
	case 0xe7:
		reg.rst(0x20, mem)
	case 0xe8:
		value := mem.readWord(reg.pc)
		reg.addSPn(value)
		reg.pc++
	case 0xe9:
		reg.jpHL()
	case 0xea:
		reg.ldnA("HL", mem)
		reg.pc += 2
	case 0xeb:
		//not used
	case 0xec:
		//not used
	case 0xed:
		//not used
	case 0xee:
		value := mem.readByte(reg.pc)
		reg.xorn(value)
		reg.pc++
	case 0xef:
		reg.rst(0x28, mem)
	case 0xf0:
		value := mem.readByte(reg.pc)
		reg.ldhAn(value, mem)
		reg.pc++
	case 0xf1:
		reg.popnn("AF", mem)
	case 0xf2:
		reg.ldAn("C", mem)
		reg.pc++
	case 0xf3:
		//not implemented
	case 0xf4:
		//not used
	case 0xf5:
		reg.pushnn("AF", mem)
	case 0xf6:
		value := mem.readByte(reg.pc)
		reg.orn(value)
		reg.pc++
	case 0xf7:
		reg.rst(0x30, mem)
	case 0xf8:
		value := mem.readByte(reg.pc)
		reg.ldHLSPn(value)
		reg.pc++
	case 0xf9:
		reg.ldSPHL()
	case 0xfa:
		var address uint16 = mem.readWord(reg.pc)
		reg.a = mem.readByte(address)
		reg.pc += 2
	case 0xfb:
		//not implemented
	case 0xfc:
		//not used
	case 0xfd:
		//not used
	case 0xfe:
		value := mem.readByte(reg.pc)
		reg.cpn(value)
		reg.pc++
	case 0xff:
		reg.rst(0x38, mem)
	}
}

func (reg *Register) executeCb(opcode byte) {
	reg.pc++
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
		reg.rrcn("A")
	case 0x10:
		reg.rln("B")
	case 0x11:
		reg.rln("C")
	case 0x12:
		reg.rln("D")
	case 0x13:
		reg.rln("E")
	case 0x14:
		reg.rln("H")
	case 0x15:
		reg.rln("L")
	case 0x16:
		reg.rln("HL")
	case 0x17:
		reg.rln("A")
	case 0x18:
		reg.rrn("B")
	case 0x19:
		reg.rrn("C")
	case 0x1a:
		reg.rrn("D")
	case 0x1b:
		reg.rrn("E")
	case 0x1c:
		reg.rrn("H")
	case 0x1d:
		reg.rrn("L")
	case 0x1e:
		reg.rrn("HL")
	case 0x1f:
		reg.rrn("A")
	case 0x20:
		reg.slan("B")
	case 0x21:
		reg.slan("C")
	case 0x22:
		reg.slan("D")
	case 0x23:
		reg.slan("E")
	case 0x24:
		reg.slan("H")
	case 0x25:
		reg.slan("L")
	case 0x26:
		reg.slan("HL")
	case 0x27:
		reg.slan("A")
	case 0x28:
		reg.sran("B")
	case 0x29:
		reg.sran("C")
	case 0x2a:
		reg.sran("D")
	case 0x2b:
		reg.sran("E")
	case 0x2c:
		reg.sran("H")
	case 0x2d:
		reg.sran("L")
	case 0x2e:
		reg.sran("HL")
	case 0x2f:
		reg.sran("A")
	case 0x30:
		reg.swapn("B")
	case 0x31:
		reg.swapn("C")
	case 0x32:
		reg.swapn("D")
	case 0x33:
		reg.swapn("E")
	case 0x34:
		reg.swapn("H")
	case 0x35:
		reg.swapn("L")
	case 0x36:
		reg.swapn("HL")
	case 0x37:
		reg.swapn("A")
	case 0x38:
		reg.srln("B")
	case 0x39:
		reg.srln("C")
	case 0x3a:
		reg.srln("D")
	case 0x3b:
		reg.srln("E")
	case 0x3c:
		reg.srln("H")
	case 0x3d:
		reg.srln("L")
	case 0x3e:
		reg.srln("HL")
	case 0x3f:
		reg.srln("A")
	case 0x40:
		reg.bitBr("B", 0)
	case 0x41:
		reg.bitBr("C", 0)
	case 0x42:
		reg.bitBr("D", 0)
	case 0x43:
		reg.bitBr("E", 0)
	case 0x44:
		reg.bitBr("H", 0)
	case 0x45:
		reg.bitBr("L", 0)
	case 0x46:
		reg.bitBr("HL", 0)
	case 0x47:
		reg.bitBr("A", 0)
	case 0x48:
		reg.bitBr("B", 1)
	case 0x49:
		reg.bitBr("C", 1)
	case 0x4a:
		reg.bitBr("D", 1)
	case 0x4b:
		reg.bitBr("E", 1)
	case 0x4c:
		reg.bitBr("H", 1)
	case 0x4d:
		reg.bitBr("L", 1)
	case 0x4e:
		reg.bitBr("HL", 1)
	case 0x4f:
		reg.bitBr("A", 1)
	case 0x50:
		reg.bitBr("B", 2)
	case 0x51:
		reg.bitBr("C", 2)
	case 0x52:
		reg.bitBr("D", 2)
	case 0x53:
		reg.bitBr("E", 2)
	case 0x54:
		reg.bitBr("H", 2)
	case 0x55:
		reg.bitBr("L", 2)
	case 0x56:
		reg.bitBr("HL", 2)
	case 0x57:
		reg.bitBr("A", 2)
	case 0x58:
		reg.bitBr("B", 3)
	case 0x59:
		reg.bitBr("C", 3)
	case 0x5a:
		reg.bitBr("D", 3)
	case 0x5b:
		reg.bitBr("E", 3)
	case 0x5c:
		reg.bitBr("H", 3)
	case 0x5d:
		reg.bitBr("L", 3)
	case 0x5e:
		reg.bitBr("HL", 3)
	case 0x5f:
		reg.bitBr("A", 3)
	case 0x60:
		reg.bitBr("B", 4)
	case 0x61:
		reg.bitBr("C", 4)
	case 0x62:
		reg.bitBr("D", 4)
	case 0x63:
		reg.bitBr("E", 4)
	case 0x64:
		reg.bitBr("H", 4)
	case 0x65:
		reg.bitBr("L", 4)
	case 0x66:
		reg.bitBr("HL", 4)
	case 0x67:
		reg.bitBr("A", 4)
	case 0x68:
		reg.bitBr("B", 5)
	case 0x69:
		reg.bitBr("C", 5)
	case 0x6a:
		reg.bitBr("D", 5)
	case 0x6b:
		reg.bitBr("E", 5)
	case 0x6c:
		reg.bitBr("H", 5)
	case 0x6d:
		reg.bitBr("L", 5)
	case 0x6e:
		reg.bitBr("HL", 5)
	case 0x6f:
		reg.bitBr("A", 5)
	case 0x70:
		reg.bitBr("B", 6)
	case 0x71:
		reg.bitBr("C", 6)
	case 0x72:
		reg.bitBr("D", 6)
	case 0x73:
		reg.bitBr("E", 6)
	case 0x74:
		reg.bitBr("H", 6)
	case 0x75:
		reg.bitBr("L", 6)
	case 0x76:
		reg.bitBr("HL", 6)
	case 0x77:
		reg.bitBr("A", 6)
	case 0x78:
		reg.bitBr("B", 7)
	case 0x79:
		reg.bitBr("C", 7)
	case 0x7a:
		reg.bitBr("D", 7)
	case 0x7b:
		reg.bitBr("E", 7)
	case 0x7c:
		reg.bitBr("H", 7)
	case 0x7d:
		reg.bitBr("L", 7)
	case 0x7e:
		reg.bitBr("HL", 7)
	case 0x7f:
		reg.bitBr("A", 7)
	case 0x80:
		reg.resBr("B", 0)
	case 0x81:
		reg.resBr("C", 0)
	case 0x82:
		reg.resBr("D", 0)
	case 0x83:
		reg.resBr("E", 0)
	case 0x84:
		reg.resBr("H", 0)
	case 0x85:
		reg.resBr("L", 0)
	case 0x86:
		reg.resBr("HL", 0)
	case 0x87:
		reg.resBr("A", 0)
	case 0x88:
		reg.resBr("B", 1)
	case 0x89:
		reg.resBr("C", 1)
	case 0x8a:
		reg.resBr("D", 1)
	case 0x8b:
		reg.resBr("E", 1)
	case 0x8c:
		reg.resBr("H", 1)
	case 0x8d:
		reg.resBr("L", 1)
	case 0x8e:
		reg.resBr("HL", 1)
	case 0x8f:
		reg.resBr("A", 1)
	case 0x90:
		reg.resBr("B", 2)
	case 0x91:
		reg.resBr("C", 2)
	case 0x92:
		reg.resBr("D", 2)
	case 0x93:
		reg.resBr("E", 2)
	case 0x94:
		reg.resBr("H", 2)
	case 0x95:
		reg.resBr("L", 2)
	case 0x96:
		reg.resBr("HL", 2)
	case 0x97:
		reg.resBr("A", 2)
	case 0x98:
		reg.resBr("B", 3)
	case 0x99:
		reg.resBr("C", 3)
	case 0x9a:
		reg.resBr("D", 3)
	case 0x9b:
		reg.resBr("E", 3)
	case 0x9c:
		reg.resBr("H", 3)
	case 0x9d:
		reg.resBr("L", 3)
	case 0x9e:
		reg.resBr("HL", 3)
	case 0x9f:
		reg.resBr("A", 3)
	case 0xa0:
		reg.resBr("B", 4)
	case 0xa1:
		reg.resBr("C", 4)
	case 0xa2:
		reg.resBr("D", 4)
	case 0xa3:
		reg.resBr("E", 4)
	case 0xa4:
		reg.resBr("H", 4)
	case 0xa5:
		reg.resBr("L", 4)
	case 0xa6:
		reg.resBr("HL", 4)
	case 0xa7:
		reg.resBr("A", 4)
	case 0xa8:
		reg.resBr("B", 5)
	case 0xa9:
		reg.resBr("C", 5)
	case 0xaa:
		reg.resBr("D", 5)
	case 0xab:
		reg.resBr("E", 5)
	case 0xac:
		reg.resBr("H", 5)
	case 0xad:
		reg.resBr("L", 5)
	case 0xae:
		reg.resBr("HL", 5)
	case 0xaf:
		reg.resBr("A", 5)
	case 0xb0:
		reg.resBr("B", 6)
	case 0xb1:
		reg.resBr("C", 6)
	case 0xb2:
		reg.resBr("D", 6)
	case 0xb3:
		reg.resBr("E", 6)
	case 0xb4:
		reg.resBr("H", 6)
	case 0xb5:
		reg.resBr("L", 6)
	case 0xb6:
		reg.resBr("HL", 6)
	case 0xb7:
		reg.resBr("A", 6)
	case 0xb8:
		reg.resBr("B", 7)
	case 0xb9:
		reg.resBr("C", 7)
	case 0xba:
		reg.resBr("D", 7)
	case 0xbb:
		reg.resBr("E", 7)
	case 0xbc:
		reg.resBr("H", 7)
	case 0xbd:
		reg.resBr("L", 7)
	case 0xbe:
		reg.resBr("HL", 7)
	case 0xbf:
		reg.setBr("A", 7)
	case 0xc0:
		reg.setBr("B", 0)
	case 0xc1:
		reg.setBr("C", 0)
	case 0xc2:
		reg.setBr("D", 0)
	case 0xc3:
		reg.setBr("E", 0)
	case 0xc4:
		reg.setBr("H", 0)
	case 0xc5:
		reg.setBr("L", 0)
	case 0xc6:
		reg.setBr("HL", 0)
	case 0xc7:
		reg.setBr("A", 0)
	case 0xc8:
		reg.setBr("B", 1)
	case 0xc9:
		reg.setBr("C", 1)
	case 0xca:
		reg.setBr("D", 1)
	case 0xcb:
		reg.setBr("E", 1)
	case 0xcc:
		reg.setBr("H", 1)
	case 0xcd:
		reg.setBr("L", 1)
	case 0xce:
		reg.setBr("HL", 1)
	case 0xcf:
		reg.setBr("A", 1)
	case 0xd0:
		reg.setBr("B", 2)
	case 0xd1:
		reg.setBr("C", 2)
	case 0xd2:
		reg.setBr("D", 2)
	case 0xd3:
		reg.setBr("E", 2)
	case 0xd4:
		reg.setBr("H", 2)
	case 0xd5:
		reg.setBr("L", 2)
	case 0xd6:
		reg.setBr("HL", 2)
	case 0xd7:
		reg.setBr("A", 2)
	case 0xd8:
		reg.setBr("B", 3)
	case 0xd9:
		reg.setBr("C", 3)
	case 0xda:
		reg.setBr("D", 3)
	case 0xdb:
		reg.setBr("E", 3)
	case 0xdc:
		reg.setBr("H", 3)
	case 0xdd:
		reg.setBr("L", 3)
	case 0xde:
		reg.setBr("HL", 3)
	case 0xdf:
		reg.setBr("A", 3)
	case 0xe0:
		reg.setBr("B", 4)
	case 0xe1:
		reg.setBr("C", 4)
	case 0xe2:
		reg.setBr("D", 4)
	case 0xe3:
		reg.setBr("E", 4)
	case 0xe4:
		reg.setBr("H", 4)
	case 0xe5:
		reg.setBr("L", 4)
	case 0xe6:
		reg.setBr("HL", 4)
	case 0xe7:
		reg.setBr("A", 4)
	case 0xe8:
		reg.setBr("B", 5)
	case 0xe9:
		reg.setBr("C", 5)
	case 0xea:
		reg.setBr("D", 5)
	case 0xeb:
		reg.setBr("E", 5)
	case 0xec:
		reg.setBr("H", 5)
	case 0xed:
		reg.setBr("L", 5)
	case 0xee:
		reg.setBr("HL", 5)
	case 0xef:
		reg.setBr("A", 5)
	case 0xf0:
		reg.setBr("B", 6)
	case 0xf1:
		reg.setBr("C", 6)
	case 0xf2:
		reg.setBr("D", 6)
	case 0xf3:
		reg.setBr("E", 6)
	case 0xf4:
		reg.setBr("H", 6)
	case 0xf5:
		reg.setBr("L", 6)
	case 0xf6:
		reg.setBr("HL", 6)
	case 0xf7:
		reg.setBr("A", 6)
	case 0xf8:
		reg.setBr("B", 7)
	case 0xf9:
		reg.setBr("C", 7)
	case 0xfa:
		reg.setBr("D", 7)
	case 0xfb:
		reg.setBr("E", 7)
	case 0xfc:
		reg.setBr("H", 7)
	case 0xfd:
		reg.setBr("L", 7)
	case 0xfe:
		reg.setBr("HL", 7)
	case 0xff:
		reg.setBr("A", 7)
	}
}
