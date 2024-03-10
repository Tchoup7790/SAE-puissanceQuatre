package main

import "net"

// Structure de données pour représenter l'état courant du jeu.
type game struct {
	gameState      int
	stateFrame     int
	grid           [globalNumTilesX][globalNumTilesY]int
	p1Color        int
	p2Color        int
	turn           int
	tokenPosition  int
	result         int
	conn           net.Conn
	playerReady    bool
	nbrPlayerReady int
	chanRead       chan map[string]string
	chanWrite      chan string
	//TODO Compteur de point ?
}

// Constantes pour représenter la séquence de jeu actuelle (écran titre,
// écran de sélection des couleurs, jeu, écran de résultats).
const (
	waitingConnection int = iota
	titleState
	colorSelectState
	waitingColor
	playState
	resultState
)

// Constantes pour représenter les pions dans la grille de puissance 4
// (absence de pion, pion du joueur 1, pion du joueur 2).
const (
	noToken int = iota
	p1Token
	p2Token
)

// Constantes pour représenter le tour de jeu (joueur 1 ou joueur 2).
const (
	p1Turn int = iota
	p2Turn
)

// Constantes pour représenter le résultat d'une partie (égalité si
// la grille a été remplie sans qu'un joueur n'ait gagné, joueur 1
// gagnant ou joueur 2 gagnant).
const (
	equality int = iota
	p1wins
	p2wins
)

// Remise à 0 du jeu pour recommencer une partie. Le joueur qui a
// perdu la dernière partie commence.
func (g *game) reset() {
	g.nbrPlayerReady = 0
	for x := 0; x < globalNumTilesX; x++ {
		for y := 0; y < globalNumTilesY; y++ {
			g.grid[x][y] = noToken
		}
	}
}
