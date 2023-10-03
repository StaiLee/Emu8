package emul8

import (
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// initialisation des constantes 
const Width = 64 * 10
const height = 32 * 10
const Font = 80

// initialisation de la struct chip8 
type Chip8 struct {
	Memory [4096]uint8
	V      [16]uint8
	GFX    [64 * 32]uint8
	Stack  [16]uint16
	Key    [16]bool
	Opcode uint16
	PC     uint16
	I      uint16
	Delay  uint8
	SP     uint16
	Sound  uint8
}

// initialisation tableau de la puce

var fontSet [Font]uint8 = [Font]uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

// fonction pour gerer les input avec chaques key
func (c *Chip8) input() {
	c.Key = [16]bool{
		ebiten.IsKeyPressed(ebiten.KeyX),
		ebiten.IsKeyPressed(ebiten.Key1),
		ebiten.IsKeyPressed(ebiten.Key2),
		ebiten.IsKeyPressed(ebiten.Key3),
		ebiten.IsKeyPressed(ebiten.KeyQ),
		ebiten.IsKeyPressed(ebiten.KeyW),
		ebiten.IsKeyPressed(ebiten.KeyE),
		ebiten.IsKeyPressed(ebiten.KeyA),
		ebiten.IsKeyPressed(ebiten.KeyS),
		ebiten.IsKeyPressed(ebiten.KeyD),
		ebiten.IsKeyPressed(ebiten.KeyZ),
		ebiten.IsKeyPressed(ebiten.KeyC),
		ebiten.IsKeyPressed(ebiten.Key4),
		ebiten.IsKeyPressed(ebiten.KeyR),
		ebiten.IsKeyPressed(ebiten.KeyF),
		ebiten.IsKeyPressed(ebiten.KeyV),
	}
}
// initialisation de la puce
func InitiateChip8() Chip8 {
	var mem [4096]uint8
	for i := 0; i < 10; i++ {
		mem[i] = fontSet[i]
	}
	return Chip8{
		PC:     0x200,
		Memory: mem,
	}
}

func (c *Chip8) opCoding() {
	c.Opcode = uint16(c.Memory[c.PC])<<8 | uint16(c.Memory[c.PC+1])
}
//emulation des opcode
func (c *Chip8) emulateOpcode() bool {
	switch c.Opcode & 0xF000 {
	case 0x0000:
		switch c.Opcode & 0x000F {
		case 0x0000:

			c.GFX = [2048]uint8{}
			c.PC += 2
		case 0x000E:

			c.SP--
			c.PC = c.Stack[c.SP]
			c.PC += 2
		default:
			panicUnknownOpcode(c.Opcode)
		}
	case 0x1000:

		c.PC = c.Opcode & 0x0FFF
	case 0x2000:

		c.Stack[c.SP] = c.PC
		c.SP++
		c.PC = c.Opcode & 0x0FFF
	case 0x3000:

		if c.V[(c.Opcode&0x0F00)>>8] == (uint8(c.Opcode) & 0x00FF) {
			c.PC += 4
		} else {
			c.PC += 2
		}
	case 0x4000:

		if c.V[(c.Opcode&0x0F00)>>8] != (uint8(c.Opcode) & 0x00FF) {
			c.PC += 4
		} else {
			c.PC += 2
		}
	case 0x5000:

		if c.V[(c.Opcode&0x0F00)>>8] != c.V[(uint8(c.Opcode)&0x00F0)>>4] {
			c.PC += 4
		} else {
			c.PC += 2
		}
	case 0x6000:

		c.V[(c.Opcode&0x0F00)>>8] = uint8(c.Opcode) & 0x00FF
		c.PC += 2
	case 0x7000:

		c.V[(c.Opcode&0x0F00)>>8] += uint8(c.Opcode) & 0x00FF
		c.PC += 2
	case 0x8000:
		switch c.Opcode & 0x000F {
		case 0x0000:

			c.V[(c.Opcode&0x0F00)>>8] = c.V[(c.Opcode&0x00F0)>>4]
			c.PC += 2
		case 0x0001:

			c.V[(c.Opcode&0x0F00)>>8] |= c.V[(c.Opcode&0x00F0)>>4]
			c.PC += 2
		case 0x0002:

			c.V[(c.Opcode&0x0F00)>>8] &= c.V[(c.Opcode&0x00F0)>>4]
			c.PC += 2
		case 0x0003:

			c.V[(c.Opcode&0x0F00)>>8] ^= c.V[(c.Opcode&0x00F0)>>4]
			c.PC += 2
		case 0x0004:

			if c.V[(c.Opcode&0x00F0)>>4] > (0xFF - c.V[(c.Opcode&0x0F00)>>8]) {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.V[(c.Opcode&0x0F00)>>8] += c.V[(c.Opcode&0x00F0)>>4]
			c.PC += 2
		case 0x0005:

			if c.V[(c.Opcode&0x00F0)>>4] > c.V[(c.Opcode&0x0F00)>>8] {

				c.V[0xF] = 0
			} else {

				c.V[0xF] = 1
			}
			c.V[(c.Opcode&0x0F00)>>8] -= c.V[(c.Opcode&0x00F0)>>4]
			c.PC += 2

		case 0x0006:

			c.V[0xF] = c.V[(c.Opcode&0x0F00)>>8] & 0x1
			c.V[(c.Opcode&0x0F00)>>8] >>= 1
			c.PC += 2
		case 0x0007:

			if c.V[(c.Opcode&0x0F00)>>8] > c.V[(c.Opcode&0x00F0)>>4] {
				c.V[0xF] = 0
			} else {
				c.V[0xF] = 1
			}
			c.V[(c.Opcode&0x0F00)>>8] = c.V[(c.Opcode&0x00F0)>>4] - c.V[(c.Opcode&0x0F00)>>8]
			c.PC += 2
		case 0x000E:

			c.V[0xF] = c.V[(c.Opcode&0x0F00)>>8] >> 7
			c.V[(c.Opcode&0x0F00)>>8] <<= 1
			c.PC += 2
		default:
			panicUnknownOpcode(c.Opcode)
		}
	case 0x9000:

		if c.V[(c.Opcode&0x0F00)>>8] != c.V[(c.Opcode&0x00F0)>>4] {
			c.PC += 4
		} else {
			c.PC += 2
		}
	case 0xA000:

		c.I = c.Opcode & 0x0FFF
		c.PC += 2
	case 0xB000:

		c.PC = (c.Opcode & 0x0FFF) + uint16(c.V[0])
	case 0xC000:

		c.V[(c.Opcode&0x0F00)>>8] = randomByte() & uint8(c.Opcode&0x00FF)
		c.PC += 2
	case 0xD000:

		x := uint16(c.V[(c.Opcode&0x0F00)>>8])
		y := uint16(c.V[(c.Opcode&0x00F0)>>4])
		height := uint16(c.Opcode & 0x000F)
		var pixel uint16

		c.V[0xF] = 0
		for yline := uint16(0); yline < height; yline++ {
			pixel = uint16(c.Memory[c.I+yline])
			for xline := uint16(0); xline < 8; xline++ {
				if (pixel & (0x80 >> xline)) != 0 {
					if c.GFX[x+xline+((y+yline)*64)] == 1 {
						c.V[0xF] = 1
					}
					c.GFX[x+xline+((y+yline)*64)] ^= 1
				}
			}
		}

		c.PC += 2

	case 0xE000:
		switch c.Opcode & 0x00FF {
		case 0x009E:

			if c.Key[c.V[(c.Opcode&0x0F00)>>8]] {
				c.PC += 4
			} else {
				c.PC += 2
			}
		case 0x00A1:

			if !c.Key[c.V[(c.Opcode&0x0F00)>>8]] {
				c.PC += 4
			} else {
				c.PC += 2
			}
		default:
			panicUnknownOpcode(c.Opcode)
		}
	case 0xF000:
		switch c.Opcode & 0x00FF {
		case 0x0007:

			c.V[(c.Opcode&0x0F00)>>8] = c.Delay
			c.PC += 2
		case 0x000A:
			keyPress := false
			for i := uint8(0); i < 16; i++ {
				if c.Key[i] {
					c.V[(c.Opcode&0x0F00)>>8] = i
					keyPress = true
				}
			}
			if !keyPress {
				return true
			}
			c.PC += 2
		case 0x0015:

			c.Delay = c.V[(c.Opcode&0x0F00)>>8]
			c.PC += 2
		case 0x0018:

			c.Sound = c.V[(c.Opcode&0x0F00)>>8]
			c.PC += 2
		case 0x001E:

			if c.I+uint16(c.V[(c.Opcode&0x0F00)>>8]) > 0xFFF {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.I += uint16(c.V[(c.Opcode&0x0F00)>>8])
			c.PC += 2
		case 0x0029:

			c.I = uint16(c.V[(c.Opcode&0x0F00)>>8]) * 0x5
			c.PC += 2
		case 0x0033:

			c.Memory[c.I] = c.V[(c.Opcode&0x0F00)>>8] / 100
			c.Memory[c.I+1] = (c.V[(c.Opcode&0x0F00)>>8] / 10) % 10
			c.Memory[c.I+2] = (c.V[(c.Opcode&0x0F00)>>8] % 100) % 10
			c.PC += 2
		case 0x0055:

			for i := uint16(0); i <= ((c.Opcode & 0x0F00) >> 8); i++ {
				c.Memory[c.I+i] = c.V[i]
			}
			c.I += ((c.Opcode & 0x0F00) >> 8) + 1
			c.PC += 2
		case 0x0065:

			for i := uint16(0); i <= ((c.Opcode & 0x0F00) >> 8); i++ {
				c.V[i] = c.Memory[c.I+i]
			}
			c.I += ((c.Opcode & 0x0F00) >> 8) + 1
			c.PC += 2

		default:
			panicUnknownOpcode(c.Opcode)
		}
	default:
		panicUnknownOpcode(c.Opcode)
	}
	return false
}
// opcode inconnus erreur 
func panicUnknownOpcode(opcode uint16) {
	log.Panicf("opcode non reconnu %v", opcode)
}

func (c *Chip8) timerRefresh() {
	if c.Delay > 0 {
		c.Delay--
	}
	if c.Sound > 0 {
		c.Sound--
	}
}
// methode dessiner a l'ecran
func (g *Chip8) Draw(screen *ebiten.Image) {
	for row := 0; row < height; row++ {
		for col := 0; col < Width; col++ {
			isOn := g.GFX[((row/10)*64)+(col/10)] == 1
			var colorToUse color.Color
			if isOn {
				colorToUse = color.White
			} else {
				colorToUse = color.Black
			}
			screen.Set(col, row, colorToUse)
		}
	}
}

func (g *Chip8) Update() error {
	g.input()
	g.Emulation()
	return nil
}

func (g *Chip8) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return Width, height
}

func (c *Chip8) LoadGUI(filename string) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("unable to read file %v", filename)
	}
	for i, b := range bytes {
		c.Memory[512+i] = b
	}
}
// créattion octet aléatoire 
func randomByte() uint8 {
	rand.Seed(time.Now().UTC().UnixNano())
	randint := rand.Intn(math.MaxUint8)
	return uint8(randint)
}
// execution de l'emulation
func (c *Chip8) Emulation() {
	c.opCoding()
	skip := c.emulateOpcode()
	if skip {
		return
	}
	c.timerRefresh()
}
