package main

import (
	"log"
)

// Mise à jour de l'état du jeu après avoir envoyer l'id du joueur
func (server *server) connection() bool {
	server.chanWrite <- "1--connection__1"
	server.chanWrite <- "2--connection__2"
	return true
}

// Mise à jour de l'état du jeu après que les couleurs soient choisi
func (server *server) colorChosen() bool {
	if server.playersReady != 2 {
		select {
		case readMap := <-server.chanRead:
			switch readMap["function"] {
			case "updateColor":
				if readMap["id"] == "1" {
					server.chanWrite <- "2--" + readMap["function"] + "__" + readMap["params"]
				} else {
					server.chanWrite <- "1--" + readMap["function"] + "__" + readMap["params"]
				}
			case "changeColor":
				server.playersReady--
			case "colorChoose":
				server.playersReady++

				if server.playersReady == 2 {
					server.chanWrite <- "1--" + "colorChoose"
					server.chanWrite <- "2--" + "colorChoose"
				}
			case "":
				log.Fatal("client " + readMap["id"] + " a quitté subitement, fin de la partie")
			default:
				log.Fatal("unknow function " + readMap["function"])
			}
		default:
		}
		return false
	} else {
		server.playersReady = 0
		return true
	}
}

// Mise à jour de l'état du jeu lors de la partie
func (server *server) game() bool {
	if server.playersReady != 2 {
		select {
		case readMap := <-server.chanRead:
			switch readMap["function"] {
			case "finish":
				server.playersReady++
			case "":
				log.Fatal("client " + readMap["id"] + " a quitté subitement, fin de la partie")
			default:
				//envoi des coordonnées
				if readMap["id"] == "1" {
					server.chanWrite <- "2--" + readMap["function"] + "__" + readMap["params"]
				} else {
					server.chanWrite <- "1--" + readMap["function"] + "__" + readMap["params"]
				}
			}
		default:
		}
		return false
	} else {
		server.playersReady = 0
		return true
	}
}

// Mise à jour de l'état du jeu lors de la fin de la partie
func (server *server) endGame() bool {
	if server.playersReady != 2 {
		select {
		case readMap := <-server.chanRead:
			switch readMap["function"] {
			case "wantRestart":
				if readMap["id"] == "1" {
					server.chanWrite <- "2--wantRestart"
				} else {
					server.chanWrite <- "1--wantRestart"
				}
				server.playersReady++
			case "colorChange":
				server.chanWrite <- "1--colorChange"
				server.chanWrite <- "2--colorChange"
				server.playersReady = 0
				server.gameState = colorSelectState
			case "quit":
				server.chanWrite <- "1--quit"
				server.chanWrite <- "2--quit"
			default:
				log.Println("unknow function " + readMap["function"])
				server.chanWrite <- "1--quickQuit"
				server.chanWrite <- "2--quickQuit"
			}
		default:
		}
		return false
	} else {
		server.playersReady = 0
		return true
	}
	return true
}
