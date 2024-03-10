package main

import "net"

// Constantes définissant les paramètres généraux du serveur
var (
	connAddr    string       = ":8080"
	netListener net.Listener = nil
)
