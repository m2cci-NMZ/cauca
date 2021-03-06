package main

type Gpu struct {
	mode           int
	mode_clock     int
	line           int
	frame_buffer   [160][144]int
	rendering      bool
	lcd            bool
	scrollX        int
	scrollY        int
	sprite         bool
	background     bool
	sprite_size    bool
	background_map bool
	background_set bool
	window         bool
	window_map     bool
	display        bool
}

const numtiles int = 512

func (gpu *Gpu) step(cpu Register, mem *Memory) {
	gpu.mode_clock += cpu.clock
	gpu.rendering = false
	switch gpu.mode {
	// Hblank
	case 0:
		if gpu.mode_clock >= 204 {
			gpu.mode_clock = 0
			gpu.line++
			mem.writeByte(0xff44, byte(gpu.line))
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
			mem.writeByte(0xff44, byte(gpu.line))
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
	scrollY := int(mem.readByte(0xff42))
	scrollX := int(mem.readByte(0xff43))
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
		x := pixel + scrollX
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
			color := gpu.getColor(pixel_1, pixel_2, mem)
			gpu.frame_buffer[map_col*8+i][y] = color

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

func (gpu *Gpu) getColor(pixel1 bool, pixel2 bool, mem Memory) int {
	var palette byte = mem.readByte(0xff47)
	var hi, lo byte
	switch true {
	case (!pixel1 && !pixel2):
		hi = 1
		lo = 0
	case (!pixel1 && pixel2):
		hi = 3
		lo = 2
	case (pixel1 && !pixel2):
		hi = 5
		lo = 4
	default:
		hi = 7
		lo = 6
	}
	color1 := hasBit(uint16(palette), uint16(hi))
	color2 := hasBit(uint16(palette), uint16(lo))
	if color1 && color2 {
		return 3
	}
	if color1 && !color2 {
		return 2
	}
	if !color1 && color2 {
		return 1
	} else {
		return 0
	}
}

func (gpu *Gpu) setGpuControl(mem Memory) {
	gpuRegister := mem.readByte(0xff40)
	gpu.lcd = hasBit(uint16(gpuRegister), 0)
	gpu.sprite = hasBit(uint16(gpuRegister), 1)
	gpu.sprite_size = hasBit(uint16(gpuRegister), 2)
	gpu.background_map = hasBit(uint16(gpuRegister), 3)
	gpu.background_set = hasBit(uint16(gpuRegister), 4)
	gpu.window = hasBit(uint16(gpuRegister), 5)
	gpu.window_map = hasBit(uint16(gpuRegister), 6)
	gpu.display = hasBit(uint16(gpuRegister), 7)
}
