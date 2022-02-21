package main

type Gpu struct {
	mode         int
	mode_clock   int
	line         byte
	frame_buffer [160][144]int
	rendering    bool
}

func (gpu *Gpu) step(cpu Register, mem *Memory) {
	gpu.mode_clock += cpu.clock
	gpu.rendering = false
	switch gpu.mode {
	// Hblank
	case 0:
		if gpu.mode_clock >= 204 {
			gpu.mode_clock = 0
			gpu.line++
			mem.writeByte(0xff44, gpu.line)
			if gpu.line == 143 {
				// last vblank, render the framebuffer
				gpu.mode = 1
				gpu.rendering = true
			} else {
				gpu.mode = 2
			}
		}
	// Vblank
	case 1:
		if gpu.mode_clock >= 456 {
			gpu.mode_clock = 0
			gpu.line++
			mem.writeByte(0xff44, gpu.line)
			if gpu.line > 153 {
				// Restart scanning
				gpu.mode = 2
				gpu.line = 0
			}
		}
	// Scanline (OAM access)
	case 2:
		if gpu.mode_clock >= 80 {
			gpu.mode_clock = 0
			gpu.mode = 3
		}
	// Scanline (VRAM access)
	case 3:
		if gpu.mode_clock >= 172 {
			gpu.mode_clock = 0
			gpu.mode = 0
			//write scanline to frame buffer
			gpu.writeScanline(*mem)
		}
	}
}

func (gpu *Gpu) writeScanline(mem Memory) {
	scrollY := mem.readByte(0xff42)
	scrollX := mem.readByte(0xff43)
	var unsigned bool = true
	var mapoffset uint16
	var bgmemory uint16
	// background map has a base offeset in vram that is
	// set according to the value of bit 4 of the LCD register
	// (adress 0xff40)
	bgmap := hasBit(uint16(mem.readByte(0xff40)), 4)
	bgmem := hasBit(uint16(mem.readByte(0xff40)), 3)
	if bgmap {
		mapoffset = 0x8800
		unsigned = false
	} else {
		mapoffset = 0x8000
	}
	if bgmem {
		bgmemory = 0x9C00
	} else {
		bgmemory = 0x9800
	}

	y := gpu.line + scrollY
	row := y / 8 * 32

	for pixel := 1; pixel < 160; pixel++ {
		x := byte(pixel) + scrollX
		col := x / 8
		tileaddress := bgmemory + uint16(row) + uint16(col)
		tilenum := mem.readByte(tileaddress)
		var tileloc byte
		if !unsigned && tilenum < 128 {
			tileloc = (tilenum+128)/16 + byte(mapoffset)
		} else {
			tileloc = tilenum/16 + byte(mapoffset)
		}
		line := y % 8
		line *= 2
		byte1 := mem.readByte(uint16(tileloc + line))
		byte2 := mem.readByte(uint16(tileloc+line) + 1)
		bitpos := ((int(x) % 8) - 7) * -1

		data1 := hasBit(uint16(byte1), uint16(bitpos))
		data2 := hasBit(uint16(byte2), uint16(bitpos))
		var color int
		if data1 && data2 {
			color = 3
		} else if data1 && !data2 {
			color = 2
		} else if !data1 && data2 {
			color = 1
		} else {
			color = 0
		}
		gpu.frame_buffer[x][y] = color
	}
}

func (gpu *Gpu) getVram(mem Memory) [26][8]byte {
	var vram [26][8]byte
	for byte_index := 0; byte_index < 26; byte_index++ {
		for i := 1; i < 8; i++ {
			var mask byte = 1 << i
			pixel1 := mem.vram[byte_index*2]
			pixel2 := mem.vram[byte_index*2+1]
			a := (pixel1 & mask) >> i
			b := (pixel2 & mask) >> i
			vram[byte_index][i] = a + b
		}
	}
	return vram
}
