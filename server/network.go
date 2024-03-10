package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
)

// fonction permettant de créer un net.Listener grâce à l'addresse en variable global
func createListener() net.Listener {
	listener, err := net.Listen("tcp", connAddr)
	if err != nil {
		log.Fatal("Listener error", err)
	}
	log.Println("Listen :", listener.Addr().String())
	return listener
}

// fonction permettant d'accepter la connection via un listener en renvoyant une connection net.Conn
func connectionAccept(listener net.Listener, id int) net.Conn {
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal("Connection error", err)
	}
	log.Println("Connection Player", id, ":", conn.RemoteAddr())
	return conn
}

// méthode sur un server permettant de créer un bufio.Writer et d'envoyer un message stocker dans le chanWrite avec un int id particulier
func (server *server) handleWrite() {
	var writerConns = []*bufio.Writer{
		bufio.NewWriter(server.conn[0]),
		bufio.NewWriter(server.conn[1]),
	}
	for {
		select {
		case input := <-server.chanWrite:
			// formattage du texte reçus
			// Attention, le texte côtés serveur à pour séparateur -- alors que le client à pour séparateur __
			inputs := strings.Split(input, "--")
			id, err := strconv.Atoi(inputs[0])
			if err != nil || id < 1 || id > 2 {
				log.Println("Writer id error")
				return
			}
			message := inputs[1]

			// envoie du texte
			_, err = writerConns[id-1].WriteString(message + "\n")
			if err != nil {
				log.Println(message, "not send to", id)
				return
			}
			err = writerConns[id-1].Flush()
			if err != nil {
				log.Println(message, "not send to", id)
				return
			}
			log.Println(message, "send to", id)
		default:
		}
	}
}

// méthode sur un server permettant de créer un bufio.Reader et de récupérer un message avec un int id particulier stocker dans le chanRead
func (server *server) handleRead(id int) {
	var readConn = bufio.NewReader(server.conn[id-1])
	for {
		message, _ := readConn.ReadString('\n')
		message = strings.Replace(message, "\n", "", 1)

		log.Println(message, "read from", id)

		inputs := strings.Split(message, "--")

		m := make(map[string]string)
		m["id"] = strconv.Itoa(id)
		m["function"] = inputs[0]

		if len(inputs) == 2 {
			m["params"] = inputs[1]
		}

		server.chanRead <- m
	}
}
