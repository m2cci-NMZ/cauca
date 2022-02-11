package main

type Gpu struct {
	mode         int
	mode_clock   int
	line         int
	frame_buffer [160][144]int
	rendering    bool
}

func (gpu *Gpu) step(cpu Register, mem Memory) {
	gpu.mode_clock += cpu.clock
	gpu.rendering = false
	switch gpu.mode {
	// Hblank
	case 0:
		if gpu.mode_clock >= 204 {
			gpu.mode_clock = 0
			gpu.line++
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
		}
	}
}
