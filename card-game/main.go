package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Card represents a single card
type Card struct {
	Value int
}

// Player represents a player in the game
type Player struct {
	ID    int
	Cards []Card
}

// PlayCard plays a card from the player's hand
func (p *Player) PlayCard() Card {
	if len(p.Cards) == 0 {
		return Card{-1} // No card to play
	}
	card := p.Cards[0]
	p.Cards = p.Cards[1:]
	return card
}

// AddCards adds cards to the player's deck
func (p *Player) AddCards(cards []Card) {
	p.Cards = append(p.Cards, cards...)
}

// HasWon checks if the player has all cards
func (p *Player) HasWon(totalCards int) bool {
	return len(p.Cards) == totalCards
}

// Game represents the game logic
type Game struct {
	Players     []Player
	TotalCards  int
	CurrentPool []Card
}

// InitializePlayers initializes the players with a shuffled deck
func InitializePlayers(numPlayers int, deck []Card) []Player {
	players := make([]Player, numPlayers)
	numCardsPerPlayer := len(deck) / numPlayers

	for i := 0; i < numPlayers; i++ {
		players[i] = Player{
			ID:    i + 1,
			Cards: deck[i*numCardsPerPlayer : (i+1)*numCardsPerPlayer],
		}
	}
	return players
}

// ShuffleDeck creates and shuffles a deck of cards
func ShuffleDeck() []Card {
	deck := make([]Card, 52)
	for j := 0; j < 4; j++ {
		for i := 0; i < 13; i++ {
			deck[j*13+i] = Card{Value: i + 1}
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	return deck
}

// PlayRound executes one round of the game
func (g *Game) PlayRound() (winner *Player) {
	g.CurrentPool = []Card{}
	roundCards := map[int]Card{}

	// Each player plays a card
	for i := range g.Players {
		card := g.Players[i].PlayCard()
		if card.Value != -1 {
			roundCards[g.Players[i].ID] = card
			g.CurrentPool = append(g.CurrentPool, card)
			fmt.Printf("Player %d plays card with value %d\n", g.Players[i].ID, card.Value)
		}
	}

	// Determine the winner
	minValue := 999
	winnerID := -1
	for playerID, card := range roundCards {
		if card.Value < minValue {
			minValue = card.Value
			winnerID = playerID
		}
	}

	// Add the pool to the winner's hand
	if winnerID != -1 {
		for i := range g.Players {
			if g.Players[i].ID == winnerID {
				g.Players[i].AddCards(g.CurrentPool)
				fmt.Printf("Player %d wins this round and collects %d cards\n", winnerID, len(g.CurrentPool))
				winner = &g.Players[i]
				fmt.Printf("Player %d has %d cards\n", g.Players[i].ID, len(g.Players[i].Cards))

				//break
			} else {
				fmt.Printf("Player %d has %d cards\n", g.Players[i].ID, len(g.Players[i].Cards))

			}
		}
	}
	return
}

// HasGameEnded checks if the game has a winner
func (g *Game) HasGameEnded() bool {
	for _, player := range g.Players {
		if player.HasWon(g.TotalCards) {
			return true
		}
	}
	return false
}

// PlayGame starts the game loop
func (g *Game) PlayGame() {
	for !g.HasGameEnded() {
		g.PlayRound()
		fmt.Println()
	}
	// Declare the winner
	for _, player := range g.Players {
		fmt.Printf("Player %d has %d cards\n", player.ID, len(player.Cards))
		if player.HasWon(g.TotalCards) {
			fmt.Printf("Player %d wins the game!\n", player.ID)
			return
		}
	}
}

func main() {
	// Initialize game
	deck := ShuffleDeck()
	numPlayers := 4
	players := InitializePlayers(numPlayers, deck)
	totalCards := len(deck)

	// Start the game
	game := Game{
		Players:    players,
		TotalCards: totalCards,
	}
	game.PlayGame()
}
