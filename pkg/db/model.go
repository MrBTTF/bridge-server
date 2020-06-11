package db

import "github.com/MrBTTF/gophercises/deck"

type Session struct {
	id      string
	players []Player
	laid    []deck.Card
	deck    []deck.Card
}

type Player struct {
	id   string
	name string
	hand []deck.Card
}
