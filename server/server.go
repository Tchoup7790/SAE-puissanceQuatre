package main

import "net"

// Structure de données pour représenter le server
type server struct {
	listener     net.Listener
	conn         [2]net.Conn
	gameState    int
	playersReady int
	chanRead     chan map[string]string
	chanWrite    chan string
	gameOn       bool
}

// Constantes pour représenter la séquence de jeu actuelle mais via le serveur
const (
	titleState int = iota
	colorSelectState
	playState
	resultState
)

// fonction créant un serveur en initialisant ses valeurs de bases
func createServer() server {
	s := server{}
	s.listener = createListener()

	s.chanRead = make(chan map[string]string, 1)
	s.chanWrite = make(chan string, 1)

	s.conn = [2]net.Conn{
		connectionAccept(s.listener, 1),
		connectionAccept(s.listener, 2),
	}
	return s
}
