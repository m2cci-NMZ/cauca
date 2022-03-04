package main

type Gpu struct {
	mode         int
	mode_clock   int
	line         byte
	frame_buffer [160][144]int
	rendering    bool
}

const numtiles int = 256

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
	var data_region uint16
	var map_region uint16
	bgmapselector := hasBit(uint16(mem.readByte(0xff40)), 4)
	bgmemselector := hasBit(uint16(mem.readByte(0xff40)), 3)
	if !bgmapselector {
		data_region = 0x8800
		unsigned = false
	} else {
		data_region = 0x8000
	}
	if bgmemselector {
		map_region = 0x9C00
	} else {
		map_region = 0x9800
	}

	y := gpu.line + scrollY
	map_line := y / 8 * 32
	tile_line := y % 8

	for pixel := 1; pixel < 160; pixel += 8 {
		x := byte(pixel) + scrollX
		map_col := x / 8
		tileaddress := map_region + uint16(map_line) + uint16(map_col)
		tile_id := mem.readByte(tileaddress)
		if !unsigned && tile_id < 128 {
			tile_id += 128
		}
		tile_data_location := data_region + uint16(tile_id)*16
		tile_data_line := tile_data_location + uint16(tile_line)*2

		data_byte_1 := mem.readByte(tile_data_line)
		data_byte_2 := mem.readByte(tile_data_line + 1)

		for i := 0; i < 8; i++ {
			pixel_1 := hasBit(uint16(data_byte_1), uint16(8-i))
			pixel_2 := hasBit(uint16(data_byte_2), uint16(8-i))
			if pixel_1 || pixel_2 {
				gpu.frame_buffer[map_col*8+byte(i)][y] = 1
			} else {
				gpu.frame_buffer[map_col*8+byte(i)][y] = 0
			}
		}
	}
}

func (gpu *Gpu) getVram(mem Memory) [numtiles * 8][8]byte {
	var vram [numtiles * 8][8]byte
	for byte_index := 0; byte_index < numtiles*8; byte_index++ {
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
