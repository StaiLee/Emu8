package main

// Importation des packages
import (
	"emu8/emul8"
	"log"
	"os"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	if len(os.Args) < 2 {    
		log.Fatal("roms8/file.ch8") // 
	}
	programName := os.Args[1] 
	chip := emul8.InitiateChip8()
	chip.LoadGUI(programName)
	ebiten.SetWindowSize(64*10, 32*10)
	ebiten.SetWindowTitle(programName)
	if err := ebiten.RunGame(&chip); err != nil {
		log.Fatal(err)
	}
}
