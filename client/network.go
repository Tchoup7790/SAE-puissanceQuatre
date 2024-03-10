package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

// méthode permettant de connecter le client au server en renvoyant false si la connection échoue
func (game *game) gameConnect() bool {
	conn, err := net.Dial("tcp", globalConnAddr)

	if err != nil {
		return false
	}

	game.conn = conn

	game.chanRead = make(chan map[string]string, 1)
	game.chanWrite = make(chan string, 1)

	go game.handleRead()
	go game.handleWrite()

	if globalDebug {
		log.Println(globalConnAddr)
	}

	game.nbrPlayerReady++

	return true
}

// méthode sur un server permettant de créer un bufio.Writer et d'envoyer un message stocker dans le chanWrite au serveur
func (game *game) handleWrite() {
	var writerConn = bufio.NewWriter(game.conn)
	for {
		select {
		case message := <-game.chanWrite:
			_, err := writerConn.WriteString(message + "\n")
			err = writerConn.Flush()
			if globalDebug {
				if err != nil {
					log.Println(message, "not send to server")
					return
				} else {
					log.Println(message, "send to server")
				}
			}
		default:
		}
	}
}

// méthode sur un server permettant de créer un bufio.Reader et de récupérer un message du serveur en le stockant dans le chanRead
func (game *game) handleRead() {
	var readConn = bufio.NewReader(game.conn)
	for {
		message, _ := readConn.ReadString('\n')
		message = strings.Replace(message, "\n", "", 1)

		if globalDebug {
			log.Println(message, "read from server")
		}

		inputs := strings.Split(message, "__")

		m := make(map[string]string)
		m["function"] = inputs[0]

		if len(inputs) == 2 {
			m["params"] = inputs[1]
		}

		game.chanRead <- m
	}
}
