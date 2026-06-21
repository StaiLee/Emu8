package emul8

import "testing"

// step sets an opcode and executes exactly one instruction.
func step(c *Chip8, opcode uint16) {
	c.Opcode = opcode
	c.emulateOpcode()
}

func TestInitLoadsFullFontSet(t *testing.T) {
	// Regression: the font is 80 bytes (16 glyphs x 5). Only loading the first 10
	// left glyphs 2..F blank, so Fx29 (I = Vx*5) pointed at zeroed memory.
	c := InitiateChip8()
	for i := 0; i < len(fontSet); i++ {
		if c.Memory[i] != fontSet[i] {
			t.Fatalf("font byte %d = 0x%02X, want 0x%02X", i, c.Memory[i], fontSet[i])
		}
	}
	if c.PC != 0x200 {
		t.Fatalf("PC = 0x%X, want 0x200", c.PC)
	}
}

func TestOpcode6xLoadAnd7xAdd(t *testing.T) {
	c := InitiateChip8()
	step(&c, 0x6005) // LD V0, 0x05
	if c.V[0] != 0x05 {
		t.Fatalf("V0 = 0x%02X, want 0x05", c.V[0])
	}
	step(&c, 0x7003) // ADD V0, 0x03
	if c.V[0] != 0x08 {
		t.Fatalf("V0 = 0x%02X, want 0x08", c.V[0])
	}
}

func TestOpcode5xy0SkipsWhenEqual(t *testing.T) {
	// Regression: 5xy0 must skip when Vx == Vy (it was inverted to !=).
	c := InitiateChip8()
	c.V[0], c.V[1] = 7, 7
	c.PC = 0x200
	step(&c, 0x5010)
	if c.PC != 0x204 {
		t.Fatalf("equal Vx/Vy: PC = 0x%X, want 0x204 (skipped)", c.PC)
	}

	c.V[1] = 8
	c.PC = 0x200
	step(&c, 0x5010)
	if c.PC != 0x202 {
		t.Fatalf("unequal Vx/Vy: PC = 0x%X, want 0x202 (no skip)", c.PC)
	}
}

func TestOpcode9xy0SkipsWhenNotEqual(t *testing.T) {
	c := InitiateChip8()
	c.V[0], c.V[1] = 1, 2
	c.PC = 0x200
	step(&c, 0x9010)
	if c.PC != 0x204 {
		t.Fatalf("PC = 0x%X, want 0x204", c.PC)
	}
}

func TestOpcode8xy4AddCarry(t *testing.T) {
	c := InitiateChip8()
	c.V[0], c.V[1] = 0xF0, 0x20
	step(&c, 0x8014)
	if c.V[0] != 0x10 {
		t.Fatalf("V0 = 0x%02X, want 0x10", c.V[0])
	}
	if c.V[0xF] != 1 {
		t.Fatalf("VF = %d, want 1 (carry)", c.V[0xF])
	}
}

func TestOpcode8xy5SubBorrow(t *testing.T) {
	c := InitiateChip8()
	c.V[0], c.V[1] = 0x50, 0x20
	step(&c, 0x8015)
	if c.V[0] != 0x30 {
		t.Fatalf("V0 = 0x%02X, want 0x30", c.V[0])
	}
	if c.V[0xF] != 1 {
		t.Fatalf("VF = %d, want 1 (no borrow)", c.V[0xF])
	}
}

func TestOpcodeAnnnAndBnnn(t *testing.T) {
	c := InitiateChip8()
	step(&c, 0xA321) // LD I, 0x321
	if c.I != 0x321 {
		t.Fatalf("I = 0x%X, want 0x321", c.I)
	}
	c.V[0] = 0x02
	step(&c, 0xB300) // JP V0, 0x300
	if c.PC != 0x302 {
		t.Fatalf("PC = 0x%X, want 0x302", c.PC)
	}
}

func TestOpcodeFx33BCD(t *testing.T) {
	c := InitiateChip8()
	c.I = 0x300
	c.V[0] = 234
	step(&c, 0xF033)
	if c.Memory[0x300] != 2 || c.Memory[0x301] != 3 || c.Memory[0x302] != 4 {
		t.Fatalf("BCD = %d%d%d, want 234",
			c.Memory[0x300], c.Memory[0x301], c.Memory[0x302])
	}
}

func TestOpcodeFx55StoreAndFx65Load(t *testing.T) {
	c := InitiateChip8()
	c.I = 0x400
	c.V[0], c.V[1], c.V[2] = 0xAA, 0xBB, 0xCC
	step(&c, 0xF255) // store V0..V2
	if c.Memory[0x400] != 0xAA || c.Memory[0x401] != 0xBB || c.Memory[0x402] != 0xCC {
		t.Fatalf("store failed: % X", c.Memory[0x400:0x403])
	}

	c2 := InitiateChip8()
	c2.I = 0x400
	c2.Memory[0x400], c2.Memory[0x401], c2.Memory[0x402] = 0x11, 0x22, 0x33
	step(&c2, 0xF265) // load V0..V2
	if c2.V[0] != 0x11 || c2.V[1] != 0x22 || c2.V[2] != 0x33 {
		t.Fatalf("load failed: % X", c2.V[0:3])
	}
}

func TestCallAndReturn(t *testing.T) {
	c := InitiateChip8()
	c.PC = 0x200
	step(&c, 0x2300) // CALL 0x300
	if c.PC != 0x300 || c.SP != 1 || c.Stack[0] != 0x200 {
		t.Fatalf("after CALL: PC=0x%X SP=%d Stack[0]=0x%X", c.PC, c.SP, c.Stack[0])
	}
	step(&c, 0x00EE) // RET
	if c.PC != 0x202 || c.SP != 0 {
		t.Fatalf("after RET: PC=0x%X SP=%d, want PC=0x202 SP=0", c.PC, c.SP)
	}
}

func TestOpcode3xnnAnd4xnnSkip(t *testing.T) {
	c := InitiateChip8()
	c.V[0] = 0x42
	c.PC = 0x200
	step(&c, 0x3042) // SE V0, 0x42 -> skip
	if c.PC != 0x204 {
		t.Fatalf("3xnn equal: PC=0x%X, want 0x204", c.PC)
	}
	c.PC = 0x200
	step(&c, 0x4042) // SNE V0, 0x42 -> no skip
	if c.PC != 0x202 {
		t.Fatalf("4xnn equal: PC=0x%X, want 0x202", c.PC)
	}
}
