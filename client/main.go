package main

import (
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font/opentype"
)

// Mise en place des polices d'écritures utilisées pour l'affichage.
func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	smallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 20,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}

	normalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 30,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}

	largeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 50,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Création d'une image annexe pour l'affichage des résultats.
func init() {
	offScreenImage = ebiten.NewImage(globalWidth, globalHeight)
}

// Création, paramétrage et lancement du jeu.
func main() {
	globalDebug = false
	globalConnAddr = "localhost:8080"

	handleOption(os.Args)

	g := game{}

	ebiten.SetWindowTitle("Programmation système : projet puissance 4")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}

// Prise en compte de paramètre de lancement
func handleOption(params []string) {
	for _, option := range params {
		options := strings.Split(option, "=")
		switch options[0] {
		case "-debug":
			globalDebug = true
			log.Println("debug mod")
			return
		case "-addr":
			if len(options) != 2 {
				log.Println(options[0], "need a parameter, put on localhost:8080")
				globalConnAddr = "localhost:8080"
			} else {
				globalConnAddr = options[1]
				log.Println("new addr :" + globalConnAddr)
			}
			return
		default:
		}
	}
}
