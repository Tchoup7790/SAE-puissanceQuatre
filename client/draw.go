package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"strconv"
)

// Affichage des graphismes à l'écran selon l'état actuel du jeu.
func (g *game) Draw(screen *ebiten.Image) {

	screen.Fill(globalBackgroundColor)

	switch g.gameState {
	case waitingConnection:
		g.connectionWaitingDraw(screen)
	case titleState:
		g.titleDraw(screen)
	case colorSelectState:
		g.colorSelectDraw(screen)
	case waitingColor:
		g.colorWaitingDraw(screen)
	case playState:
		g.playDraw(screen)
	case resultState:
		g.resultDraw(screen)
	}

}

// Affichage des graphismes de l'écran d'attente du serveur.
func (g game) connectionWaitingDraw(screen *ebiten.Image) {
	text.Draw(screen, "Puissance 4 en réseau", largeFont, 90, 150, globalTextColor)
	text.Draw(screen, "Projet de programmation système", normalFont, 105, 190, globalTextColor)
	text.Draw(screen, "Année 2023-2024", normalFont, 210, 230, globalTextColor)

	if globalDebug {
		text.Draw(screen, "debug", smallFont, 20, 20, globalTextColor)
	}

	if g.stateFrame >= globalBlinkDuration/3 {
		text.Draw(screen, "Attente de la connection du serveur", normalFont, 110, 500,
			globalTextColor)
		text.Draw(screen, strconv.Itoa(g.nbrPlayerReady)+" joueurs connectés", smallFont, 250, 600, globalTextColor)
	}
}

// Affichage des graphismes de l'écran d'attente de la deuxième couleur.
func (g game) colorWaitingDraw(screen *ebiten.Image) {
	text.Draw(screen, "Couleur choisie", largeFont, 170, 150, globalTextColor)

	if g.stateFrame >= globalBlinkDuration/3 {
		text.Draw(screen, "Attente du choix du second joueur", normalFont, 120, 190, globalTextColor)
		text.Draw(screen, "Presser esc pour changer de couleur", normalFont, 110, 600, globalTextColor)
	}

	text.Draw(screen, "Joueur 1", normalFont, 200, 350, globalTextColor)
	vector.DrawFilledCircle(screen, 250, 450, globalTileSize/2, globalSelectColor, true)
	vector.DrawFilledCircle(screen, 250, 450, globalTileSize/2-globalCircleMargin,
		globalTokenColors[g.p1Color], true)

	text.Draw(screen, "Joueur 2", normalFont, 390, 350, globalTextColor)
	vector.DrawFilledCircle(screen, 450, 450, globalTileSize/2, globalSelectColorp2, true)
	vector.DrawFilledCircle(screen, 450, 450, globalTileSize/2-globalCircleMargin,
		globalTokenColors[g.p2Color], true)
}

// Affichage des graphismes de l'écran titre.
func (g game) titleDraw(screen *ebiten.Image) {
	text.Draw(screen, "Puissance 4 en réseau", largeFont, 90, 150, globalTextColor)
	text.Draw(screen, "Projet de programmation système", normalFont, 105, 190, globalTextColor)
	text.Draw(screen, "Année 2023-2024", normalFont, 210, 230, globalTextColor)
	text.Draw(screen, strconv.Itoa(g.nbrPlayerReady)+" joueurs connectés", smallFont, 250, 600, globalTextColor)

	if globalDebug {
		text.Draw(screen, "debug", smallFont, 640, 20, globalTextColor)
	}

	if g.stateFrame >= globalBlinkDuration/3 {
		if g.playerReady {
			text.Draw(screen, "Presser entrer pour commencer", normalFont, 130, 500, globalTextColor)
		} else {
			text.Draw(screen, "Attente d'un second joueur", normalFont, 160, 500, globalTextColor)
		}
	}
}

// Affichage des graphismes de l'écran de sélection des couleurs des joueurs.
func (g game) colorSelectDraw(screen *ebiten.Image) {
	text.Draw(screen, "Quelle couleur pour vos pions ?", normalFont, 110, 80, globalTextColor)

	line := 0
	col := 0
	for numColor := 0; numColor < globalNumColor; numColor++ {

		xPos := (globalNumTilesX-globalNumColorCol)/2 + col
		yPos := (globalNumTilesY-globalNumColorLine)/2 + line

		if numColor == g.p1Color {
			vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize),
				float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2, globalSelectColor, true)
		} else if numColor == g.p2Color {
			vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize),
				float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2, globalSelectColorp2,
				true)
		}

		vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2-globalCircleMargin, globalTokenColors[numColor], true)

		col++
		if col >= globalNumColorCol {
			col = 0
			line++
		}
	}
}

// Affichage des graphismes durant le jeu.
func (g game) playDraw(screen *ebiten.Image) {
	if g.turn == p1Turn {
		g.drawGrid(screen)

		vector.DrawFilledCircle(screen, float32(globalTileSize/2+g.tokenPosition*globalTileSize), float32(globalTileSize/2), globalTileSize/2-globalCircleMargin, globalTokenColors[g.p1Color], true)
	} else {
		g.drawGrid(offScreenImage)

		options := &ebiten.DrawImageOptions{}
		options.ColorScale.ScaleAlpha(0.2)

		screen.DrawImage(offScreenImage, options)

		vector.DrawFilledCircle(screen, float32(globalTileSize/2+g.tokenPosition*globalTileSize), float32(globalTileSize/2), globalTileSize/2-globalCircleMargin, globalTokenColors[g.p2Color], true)
		if g.stateFrame >= globalBlinkDuration/3 {
			text.Draw(screen, "Attente de l'adversaire", normalFont, 165, 600, globalTextColor)
		}
	}
}

// Affichage des graphismes à l'écran des résultats.
func (g game) resultDraw(screen *ebiten.Image) {
	g.drawGrid(offScreenImage)

	options := &ebiten.DrawImageOptions{}
	options.ColorScale.ScaleAlpha(0.2)
	screen.DrawImage(offScreenImage, options)

	message := "Égalité"
	if g.result == p1wins {
		message = "Gagné !"
	} else if g.result == p2wins {
		message = "Perdu…"
	}
	text.Draw(screen, message, normalFont, 300, 400, globalTextColor)
	text.Draw(screen, strconv.Itoa(g.nbrPlayerReady)+" joueurs prêt", normalFont, 250, 645, globalTextColor)
	if g.stateFrame >= globalBlinkDuration/3 {
		text.Draw(screen, "Presser c pour changer les couleurs", normalFont, 120, 525, globalTextColor)
		text.Draw(screen, "Presser esc pour quitter le jeu", normalFont, 120, 575, globalTextColor)
		if !g.playerReady {
			text.Draw(screen, "Presser entrer pour recommencer", normalFont, 120, 475, globalTextColor)
		} else {
			text.Draw(screen, "Attente du second joueur", normalFont, 170, 475, globalTextColor)
		}
	}
}

// Affichage de la grille de puissance 4, incluant les pions déjà joués.
func (g game) drawGrid(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, globalTileSize, globalTileSize*globalNumTilesX, globalTileSize*globalNumTilesY, globalGridColor, true)

	for x := 0; x < globalNumTilesX; x++ {
		for y := 0; y < globalNumTilesY; y++ {

			var tileColor color.Color
			switch g.grid[x][y] {
			case p1Token:
				tileColor = globalTokenColors[g.p1Color]
			case p2Token:
				tileColor = globalTokenColors[g.p2Color]
			default:
				tileColor = globalBackgroundColor
			}

			vector.DrawFilledCircle(screen, float32(globalTileSize/2+x*globalTileSize), float32(globalTileSize+globalTileSize/2+y*globalTileSize), globalTileSize/2-globalCircleMargin, tileColor, true)
		}
	}
}
