package snakeandladdergame

import (
	"fmt"
	"sync"
)

type SnakeAndLadderGame struct {
	ID               int
	Board            *Board
	Players          []*Player
	Dice             *Dice
	CurrentPlayerIdx int
}

func NewSnakeAndLadderGame(ID int, playerNames []string) *SnakeAndLadderGame {
	game := &SnakeAndLadderGame{
		ID:               ID,
		Board:            NewBoard(),
		Dice:             NewDice(),
		Players:          []*Player{},
		CurrentPlayerIdx: 0,
	}

	for _, name := range playerNames {
		game.Players = append(game.Players, NewPlayer(name))
	}
	return game
}

func (g *SnakeAndLadderGame) Play(wg *sync.WaitGroup) {
	for !g.isGameOver() {
		player := g.Players[g.CurrentPlayerIdx]
		roll := g.Dice.Roll()
		newPosition := player.Position + roll

		if newPosition <= g.Board.Size {
			player.Position = g.Board.GetNewPosition(newPosition)
			fmt.Printf("Game: %d - %s rolled a %d and moved to position %d\n", g.ID, player.Name, roll, player.Position)
		}

		if player.Position == g.Board.Size {
			fmt.Printf("For Game %d - %s wins!\n", g.ID, player.Name)
			break
		}

		g.CurrentPlayerIdx = (g.CurrentPlayerIdx + 1) % len(g.Players)
	}
	wg.Done()
}

func (g *SnakeAndLadderGame) isGameOver() bool {
	for _, player := range g.Players {
		if player.Position == g.Board.Size {
			return true
		}
	}
	return false
}
