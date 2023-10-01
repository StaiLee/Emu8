package emul8

import (
	"fmt"
	"math/rand"
)

const (
	screenWidth  = 64
	screenHeight = 32
)

type GoChip8 struct {
	memory     [4096]byte
	v          [16]byte
	i          uint16
	pc         uint16
	stack      []uint16
	sp         int
	display    [screenWidth * screenHeight]byte
	key        [16]bool
	delay      uint8
	sound      uint8
	shouldDraw bool
}

var fontSet = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, //0
	0x20, 0x60, 0x20, 0x20, 0x70, //1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
	0x90, 0x90, 0xF0, 0x10, 0x10, //4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
	0xF0, 0x10, 0x20, 0x40, 0x40, //7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
	0xF0, 0x90, 0xF0, 0x90, 0x90, //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
	0xF0, 0x80, 0x80, 0x80, 0xF0, //C
	0xE0, 0x90, 0x90, 0x90, 0xE0, //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
	0xF0, 0x80, 0xF0, 0x80, 0x80, //F
}

func (c *GoChip8) emulateOpcode(opcode uint16) {
	// Extraire les parties de l'opcode
	x := int((opcode & 0x0F00) >> 8)
	y := int((opcode & 0x00F0) >> 4)
	kk := byte(opcode & 0x00FF)
	nnn := opcode & 0x0FFF

	switch opcode & 0xF000 {
	case 0x00E0:
		// 0x00E0: Effacer l'écran (mettre tous les pixels à zéro)
		for i := 0; i < len(c.display); i++ {
			c.display[i] = 0
		}
		c.shouldDraw = true
		c.pc += 2

	case 0x00EE:
		// 0x00EE: Revenir d'une sous-routine
		if c.sp > 0 {
			c.sp--
			c.pc = c.stack[c.sp]
		}

	default:
		fmt.Printf("Opcode non fonctionel: %X\n", opcode)
	}

	switch opcode & 0xF000 {
	case 0x1000:
		// 0x1NNN: Sauter à l'adresse NNN
		c.pc = nnn

	case 0x2000:
		// 0x2NNN: Appeler une sous-routine à l'adresse NNN
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = nnn

	case 0x3000:
		// 0x3XKK: Passer à l'instruction suivante si VX est égal à KK
		if c.v[x] == kk {
			c.pc += 2
		}

	case 0x4000:
		// 0x4XKK: Passer à l'instruction suivante si VX n'est pas égal à KK
		if c.v[x] != kk {
			c.pc += 2
		}

	case 0x5000:
		// 0x5XY0: Passer à l'instruction suivante si VX est égal à VY
		if c.v[x] == c.v[y] {
			c.pc += 2
		}

	case 0x6000:
		// 0x6XKK: Définir VX sur KK
		c.v[x] = kk

	case 0x7000:
		// 0x7XKK: Ajouter KK à VX
		c.v[x] += kk

	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000:
			// 0x8XY0: Définir VX sur la valeur de VY
			c.v[x] = c.v[y]

		case 0x0001:
			// 0x8XY1: Définir VX sur VX OU VY
			c.v[x] |= c.v[y]

		case 0x0002:
			// 0x8XY2: Définir VX sur VX ET VY
			c.v[x] &= c.v[y]

		case 0x0003:
			// 0x8XY3: Définir VX sur VX XOR VY
			c.v[x] ^= c.v[y]

		case 0x0004:
			// 0x8XY4: Ajouter VY à VX, définir VF à 1 s'il y a une retenue
			if c.v[x] > 0xFF-c.v[y] {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[x] += c.v[y]

		case 0x0005:
			// 0x8XY5: Soustraire VY de VX, définir VF à 0 s'il y a un emprunt
			if c.v[x] > c.v[y] {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[x] -= c.v[y]

		case 0x0006:
			// 0x8XY6: Décaler VY vers la droite de 1 et stocker le résultat dans VX, définir VF sur le bit le moins significatif de VY avant le décalage
			c.v[0xF] = c.v[y] & 0x01
			c.v[x] = c.v[y] >> 1

		case 0x0007:
			// 0x8XY7: Définir VX sur VY moins VX, définir VF à 0 s'il y a un emprunt
			if c.v[y] > c.v[x] {
				c.v[0xF] = 0
			} else {
				c.v[0xF] = 1
			}
			c.v[x] = c.v[y] - c.v[x]

		case 0x000E:
			// 0x8XYE: Décaler VY vers la gauche de 1 et stocker le résultat dans VX, définir VF sur le bit le plus significatif de VY avant le décalage
			c.v[0xF] = (c.v[y] >> 7) & 0x01
			c.v[x] = c.v[y] << 1

		default:
			fmt.Printf("op code non fonctionel: %X\n", opcode)
		}

	case 0x9000:
		// 0x9XY0: Passer à l'instruction suivante si VX n'est pas égal à VY
		if c.v[x] != c.v[y] {
			c.pc += 2
		}

	case 0xA000:
		// 0xANNN: Définir I sur l'adresse NNN
		c.i = nnn

	case 0xB000:
		// 0xBNNN: Sauter à l'adresse NNN plus V0
		c.pc = nnn + uint16(c.v[0])

	case 0xC000:
		// 0xCXKK: Définir VX sur le résultat d'une opération ET binaire sur un nombre aléatoire et KK
		randomValue := byte(rand.Intn(256))
		c.v[x] = randomValue & kk

	case 0xD000:
		// 0xDXYN: Dessiner un sprite aux coordonnées (VX, VY)
		// La logique pour dessiner un sprite doit être implémentée ici

	default:
		fmt.Printf("op code non fonctionel: %X\n", opcode)
	}
}
