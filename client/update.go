package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
	"strconv"
)

// Mise à jour de l'état du jeu en fonction des entrées au clavier.
func (g *game) Update() error {

	g.stateFrame++

	switch g.gameState {
	case waitingConnection:
		if g.gameConnect() {
			g.gameState++
		}
	case titleState:
		if g.titleUpdate() {
			g.gameState++
			g.nbrPlayerReady = 0
			g.playerReady = false
		}
	case colorSelectState:
		if g.colorSelectUpdate() {
			g.gameState++
		}
	case waitingColor:
		if g.waitingColor() {
			g.gameState++
		}
	case playState:
		g.tokenPosUpdate()
		var lastXPositionPlayed int
		var lastYPositionPlayed int
		if g.turn == p1Turn {
			lastXPositionPlayed, lastYPositionPlayed = g.p1Update()
		} else {
			lastXPositionPlayed, lastYPositionPlayed = g.p2Update()
		}
		if lastXPositionPlayed >= 0 {
			finished, result := g.checkGameEnd(lastXPositionPlayed, lastYPositionPlayed)
			if finished {
				g.result = result
				g.chanWrite <- "finish--" + strconv.Itoa(result)
				g.gameState++
			}
		}
	case resultState:
		if g.resultUpdate() {
			g.reset()
			g.playerReady = false
			g.gameState = playState
		}
	}

	return nil
}

// Mise à jour de l'état du jeu à l'écran titre et gestion du premier ou second joueur.
func (g *game) titleUpdate() bool {
	g.stateFrame = g.stateFrame % globalBlinkDuration

	if !g.playerReady {
		select {
		case readMap := <-g.chanRead:
			if readMap["params"] == "1" {
				g.turn = p1Turn
				g.p1Color = 0
				g.p2Color = 1
			} else {
				g.turn = p2Turn
				g.p1Color = 1
				g.p2Color = 0
			}
			g.playerReady = true
			g.nbrPlayerReady++
		default:
		}
	}

	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) && g.playerReady
}

// Mise à jour de l'état du jeu lors de la sélection des couleurs.
func (g *game) colorSelectUpdate() bool {

	col := g.p1Color % globalNumColorCol
	line := g.p1Color / globalNumColorLine

	updateColor := false

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		updateColor = true
		col = (col + 1) % globalNumColorCol
		g.p1Color = line*globalNumColorLine + col
		if g.p1Color == g.p2Color {
			col = (col + 1) % globalNumColorCol
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		updateColor = true
		col = (col - 1 + globalNumColorCol) % globalNumColorCol
		g.p1Color = line*globalNumColorLine + col
		if g.p1Color == g.p2Color {
			col = (col - 1 + globalNumColorCol) % globalNumColorCol
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		updateColor = true
		line = (line + 1) % globalNumColorLine
		g.p1Color = line*globalNumColorLine + col
		if g.p1Color == g.p2Color {
			line = (line + 1) % globalNumColorLine
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		updateColor = true
		line = (line - 1 + globalNumColorLine) % globalNumColorLine
		g.p1Color = line*globalNumColorLine + col
		if g.p1Color == g.p2Color {
			line = (line - 1 + globalNumColorLine) % globalNumColorLine
		}
	}

	if updateColor {
		g.p1Color = line*globalNumColorLine + col
		g.chanWrite <- "updateColor--" + strconv.Itoa(g.p1Color)
		updateColor = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && !g.playerReady {
		g.chanWrite <- "colorChoose--" + strconv.Itoa(g.p1Color)
		return true
	}

	select {
	case readMap := <-g.chanRead:
		if readMap["function"] == "updateColor" || readMap["function"] == "colorChoose" {
			g.p2Color, _ = strconv.Atoi(readMap["params"])
		} else {
			log.Fatal("error quick quit")
		}
	default:
	}
	return false
}

// Mise à jour de l'état du jeu à l'attente que le joueur 2 choisisse sa couleur.
func (g *game) waitingColor() bool {
	select {
	case readMap := <-g.chanRead:
		switch readMap["function"] {
		case "colorChoose":
			return true
		case "updateColor":
			g.p2Color, _ = strconv.Atoi(readMap["params"])
		default:
			log.Fatal("error quick quit")
		}
	default:
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.chanWrite <- "changeColor"
		g.gameState--
	}

	return false
}

// Gestion de la position du prochain pion à jouer par le joueur 1.
func (g *game) tokenPosUpdate() {
	if g.turn == p1Turn {
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			g.tokenPosition = (g.tokenPosition - 1 + globalNumTilesX) % globalNumTilesX
			g.chanWrite <- "updateMove--" + strconv.Itoa(g.tokenPosition)
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			g.tokenPosition = (g.tokenPosition + 1) % globalNumTilesX
			g.chanWrite <- "updateMove--" + strconv.Itoa(g.tokenPosition)
		}
	}
}

// Gestion du moment où le prochain pion est joué par le joueur 1.
func (g *game) p1Update() (int, int) {
	lastXPositionPlayed := -1
	lastYPositionPlayed := -1
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if updated, yPos := g.updateGrid(p1Token, g.tokenPosition); updated {
			g.turn = p2Turn
			lastXPositionPlayed = g.tokenPosition
			lastYPositionPlayed = yPos
			g.chanWrite <- "move--" + strconv.Itoa(lastXPositionPlayed)
		}
	}
	return lastXPositionPlayed, lastYPositionPlayed
}

// Gestion de la position du prochain pion joué par le joueur 2 et
// du moment où ce pion est joué.
func (g *game) p2Update() (int, int) {
	select {
	case readMap := <-g.chanRead:
		switch readMap["function"] {
		case "updateMove":
			g.tokenPosition, _ = strconv.Atoi(readMap["params"])
		case "quickQuit":
			log.Fatal("error quick quit")
		case "move":
			xposition, _ := strconv.Atoi(readMap["params"])
			if updated, yposition := g.updateGrid(p2Token, xposition); updated {
				g.turn = p1Turn
				return xposition, yposition
			}
		default:
			log.Fatal("error quick quit")
		}
	default:
	}
	return -1, -1
}

// Mise à jour de l'état du jeu à l'écran des résultats.
func (g *game) resultUpdate() bool {
	select {
	case readMap := <-g.chanRead:
		switch readMap["function"] {
		case "colorChange":
			g.gameState = colorSelectState
			g.playerReady = false
			g.reset()
			return false
		case "wantRestart":
			g.nbrPlayerReady++
		case "restart":
			return true
		case "quit":
			log.Fatal("fin de partie")
		default:
			log.Fatal("error quick quit")
		}
	default:
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.chanWrite <- "colorChange"
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.chanWrite <- "quit"
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && !g.playerReady {
		g.playerReady = true
		g.chanWrite <- "wantRestart"
		g.nbrPlayerReady++
	}
	return false
}

// Mise à jour de la grille de jeu lorsqu'un pion est inséré dans la
// colonne de coordonnée (x) position.
func (g *game) updateGrid(token, position int) (updated bool, yPos int) {
	for y := globalNumTilesY - 1; y >= 0; y-- {
		if g.grid[position][y] == noToken {
			updated = true
			yPos = y
			g.grid[position][y] = token
			return
		}
	}
	return
}

// Vérification de la fin du jeu : est-ce que le dernier joueur qui
// a placé un pion gagne ? est-ce que la grille est remplie sans gagnant
// (égalité) ? ou est-ce que le jeu doit continuer ?
func (g game) checkGameEnd(xPos, yPos int) (finished bool, result int) {

	tokenType := g.grid[xPos][yPos]

	// horizontal
	count := 0
	for x := xPos; x < globalNumTilesX && g.grid[x][yPos] == tokenType; x++ {
		count++
	}
	for x := xPos - 1; x >= 0 && g.grid[x][yPos] == tokenType; x-- {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// vertical
	count = 0
	for y := yPos; y < globalNumTilesY && g.grid[xPos][y] == tokenType; y++ {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut gauche/bas droit
	count = 0
	for x, y := xPos, yPos; x < globalNumTilesX && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x+1, y+1 {
		count++
	}

	for x, y := xPos-1, yPos-1; x >= 0 && y >= 0 && g.grid[x][y] == tokenType; x, y = x-1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut droit/bas gauche
	count = 0
	for x, y := xPos, yPos; x >= 0 && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x-1, y+1 {
		count++
	}

	for x, y := xPos+1, yPos-1; x < globalNumTilesX && y >= 0 && g.grid[x][y] == tokenType; x, y = x+1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// egalité ?
	if yPos == 0 {
		for x := 0; x < globalNumTilesX; x++ {
			if g.grid[x][0] == noToken {
				return
			}
		}
		return true, equality
	}

	return
}
