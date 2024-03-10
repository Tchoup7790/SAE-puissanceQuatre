package main

// Création, paramétrage et lancement du serveur.
func main() {
	var s = createServer()

	go s.handleWrite()
	go s.handleRead(1)
	go s.handleRead(2)

	s.gameState = 0
	s.gameOn = true

	for s.gameOn {
		switch {
		case s.gameState == titleState:
			if s.connection() {
				s.gameState++
			}
		case s.gameState == colorSelectState:
			if s.colorChosen() {
				s.gameState++
			}
		case s.gameState == playState:
			if s.game() {
				s.gameState++
			}
		case s.gameState == resultState:
			if s.endGame() {
				s.gameState--
				s.chanWrite <- "1--restart"
				s.chanWrite <- "2--restart"
			}
		}
	}
}
