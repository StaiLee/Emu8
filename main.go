package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// Les constantes pour la taille de l'écran du Chip-8
const (
	screenWidth  = 64
	screenHeight = 32
)

// GoChip8 représente l'émulateur Chip-8
type GoChip8 struct {
	memory  [4096]byte
	v       [16]byte
	i       uint16
	pc      uint16
	stack   []uint16
	sp      int
	display [screenWidth * screenHeight]byte
	key     [16]bool
	delay   uint8
	sound   uint8
}

// NewGoChip8 initialise un nouvel émulateur Chip-8
func NewGoChip8() *GoChip8 {
	// Créez une instance de GoChip8 et initialisez ses membres ici
	// Par exemple, initialisez la mémoire, les registres, etc.
	return &GoChip8{}
}

// LoadROM charge un programme Chip-8 depuis un fichier ROM
func (c8 *GoChip8) LoadROM(filename string) error {
	// Lisez le contenu du fichier ROM et chargez-le dans la mémoire de l'émulateur
	// Assurez-vous que l'émulateur est configuré pour exécuter le programme depuis l'adresse mémoire appropriée

	file, err := os.Open("roms8/c8_test.c8")
	if err != nil {
		return err
	}
	defer file.Close()

	// Lisez le contenu du fichier et chargez-le dans la mémoire
	// Utilisez un tampon pour stocker les données lues du fichier
	// Initialisez également le compteur de programme (PC) avec l'adresse de début du programme

	// Exemple (vous devrez adapter cela à votre structure de données) :
	// addr := 0x200 // Adresse de début standard pour les programmes Chip-8
	// buffer := make([]byte, 0x1000-addr)
	// _, err = file.Read(buffer)
	// if err != nil {
	// 	return err
	// }
	// copy(c8.memory[addr:], buffer)

	return nil
}

func main() {
	// Initialisez l'émulateur Chip-8
	goChip8 := NewGoChip8()

	// Chargez un programme Chip-8 depuis un fichier ROM
	if len(os.Args) != 2 {
		fmt.Println("Usage: chip8-emu <ROM file>")
		return
	}
	if err := goChip8.LoadROM(os.Args[1]); err != nil {
		log.Fatal(err)
	}

	// Initialisez Ebiten
	ebiten.SetWindowSize(screenWidth*10, screenHeight*10)
	ebiten.SetWindowTitle("Chip-8 Emulator")
	if err := ebiten.RunGame(goChip8); err != nil {
		log.Fatal(err)
	}
}
