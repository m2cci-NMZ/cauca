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

mem Memory

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
/* 16 bit loads                            */
/* *************************************** */

func (reg *Register) ldnnn(value uint16, destination string) {
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
	value := (uint16(H) << 8) + uint16(L)
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
	if (reg.sp&0x0F + value&0x0F) > 0x0F {
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

func (reg *Register) ldnnSP(value uint16) {
	mem.writeWord(value,reg.sp)
}

func (reg *Register) pushnn(registers string) {
	r1,r2 byte
	switch registers {
	case "AF":
		r1=unit16(reg.A)<<8
		r2=uint16(reg.b)
	case "BC":
		r1=unit16(reg.B)<<8
		r2=uint16(reg.C)
	case "DE":
		r1=unit16(reg.D)<<8
		r2=uint16(reg.E)
	case "HL":
		r1=unit16(reg.H)<<8
		r2=uint16(reg.L)
	}
	value:=r1+r2
	mem.writeWord(reg.sp,value)
	reg.sp = reg.sp - 2
}

func (reg *Register) popnn(registers string){
	r1:=mem.readByte(reg.sp)
	r2:=mem.readByte(reg.sp+1)
	switch registers{
	case "AF":
		reg.a =r1
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

func (reg *Register) subn(value byte){
	result := uint16(reg.a) - uint16(value)
	// negative flag
	if (result < 0){
		reg.setRegisterFlag(true, 6)
	}
	else{
		reg.setRegisterFlag(false, 6)
	}
	// zero flag
	if (result == 0){
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

func (reg *Register) sbcAn(value byte){
	if hasBit(uint16(reg.flags), 4) {
		value++
	}
	reg.subn(value)	
}

func (reg *Register) andn(value byte){
	result:= reg.A & value
	// zero flag
	if (result == 0){
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

func (reg *Register) orn(value byte){
	result:= reg.A | value
	// zero flag
	if (result == 0){
		reg.setRegisterFlag(true, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	// carry flag
	reg.setRegisterFlag(false, 4)
	reg.A = result
}

func (reg *Register) xorn(value byte){
	result:= reg.A ^ value
	// zero flag
	if (result == 0){
		reg.setRegisterFlag(true, 7)
	}
	// negative flag
	reg.setRegisterFlag(false, 6)
	// half carry flag
	reg.setRegisterFlag(false, 5)
	// carry flag
	reg.setRegisterFlag(false, 4)
	reg.A = result
}

func (reg *Register) cpn(value byte){
	tmp:=reg.A
	subn(value)
	reg.A = tmp
}

func (reg *Register) incn(register string){
	result int16
	result++
	switch register {
		case "A":
			reg.A = result
		case "B":
			reg.A = result
		case "C":
			reg.A = result
		case "D":
			reg.A = result
		case "E":
			reg.A = result
		case "F":
			reg.A = result
		case "H":
			reg.A = result
		case "L":
			reg.A = result
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
	if (reg.a&0x0F + value&0x0F) > 0x0F {
		reg.setRegisterFlag(true, 5)
	} else {
		reg.setRegisterFlag(false, 5)
	}
}

func (reg *Register) decn(register string){
	result int16
	result--
	switch register {
		case "A":
			reg.A = result
		case "B":
			reg.A = result
		case "C":
			reg.A = result
		case "D":
			reg.A = result
		case "E":
			reg.A = result
		case "F":
			reg.A = result
		case "H":
			reg.A = result
		case "L":
			reg.A = result
	}
	// negative flag
	if (result < 0){
		reg.setRegisterFlag(true, 6)
	}
	else{
		reg.setRegisterFlag(false, 6)
	}
	// zero flag
	if (result == 0){
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
}
/* *************************************** */
/* 16 bit ALU                               */
/* *************************************** */

func (reg *Register) addHLn (value int16){
	result int32
	HL:= int16(H)<<8 + int16(L)
	result = int32(HL) + int32(value)
	// negative flag
	reg.setRegisterFlag(false, 6)
	// carry flag
	if (result & 0xFFFF0000) != 0{
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	// half carry flag
	if ((int16(result) & 0x0F) + (value & 0x0F)) > 0x0F{
		reg.setRegisterFlag(true, 5)
	} else {
		reg.setRegisterFlag(false, 5)
	}
	reg.H = byte(result<<8)
	reg.L = byte(result)
}

func (reg *Register) addSPn(value int16){
	result:= int32(reg.SP) + int32(value)
	// zero flag
	reg.setRegisterFlag(false, 7)
	// negative flag
	reg.setRegisterFlag(false, 6)
	// carry flag
	if (result & 0xFFFF0000) != 0{
		reg.setRegisterFlag(true, 4)
	} else {
		reg.setRegisterFlag(false, 4)
	}
	// half carry flag
	if ((int16(result) & 0x0F) + (value & 0x0F)) > 0x0F{
		reg.setRegisterFlag(true, 5)
	} else {
		reg.setRegisterFlag(false, 5)
	}
	reg.sp = int16(result)
}

func (reg *Register) incnn(register string){
	switch register{
	case "BC":
		result := int16(reg.B)<<8 + int16(reg.C) + 1
		reg.B = byte(result>>8)
		reg.C = byte(result)
	case "DE":
		result := int16(reg.D)<<8 + int16(reg.E) + 1
		reg.D = byte(result>>8)
		reg.E = byte(result)
	case "HL":
		result := int16(reg.H)<<8 + int16(reg.L) + 1
		reg.H = byte(result>>8)
		reg.L = byte(result)
	case "SP":
		reg.SP++
	}
}

func (reg *Register) decnn(register string){
	switch register{
	case "BC":
		result := int16(reg.B)<<8 + int16(reg.C) - 1
		reg.B = byte(result>>8)
		reg.C = byte(result)
	case "DE":
		result := int16(reg.D)<<8 + int16(reg.E) - 1
		reg.D = byte(result>>8)
		reg.E = byte(result)
	case "HL":
		result := int16(reg.H)<<8 + int16(reg.L) - 1
		reg.H = byte(result>>8)
		reg.L = byte(result)
	case "SP":
		reg.SP--
	}
}