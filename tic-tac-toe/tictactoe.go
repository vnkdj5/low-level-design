package tictactoe

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const BOARD_SIZE int = 3

type Player struct {
	Name   string
	Symbol rune
}

func NewPlayer(name string, symbol rune) *Player {
	return &Player{
		Name:   name,
		Symbol: symbol,
	}
}

type Board struct {
	Grid       [3][3]rune
	MovesCount int
}

func NewBoard() *Board {
	b := &Board{}
	b.InitializeBoard()
	return b
}

func (b *Board) InitializeBoard() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			b.Grid[i][j] = '-'
		}
	}
	b.MovesCount = 0
}

func (b *Board) MakeMove(row, col int, symbol rune) error {
	if row < 0 || row >= BOARD_SIZE || col < 0 || col >= BOARD_SIZE {
		return fmt.Errorf("Invalid row or column")
	}
	b.Grid[row][col] = symbol
	b.MovesCount++
	return nil
}

func (b *Board) PrintBoard() {
	for i := 0; i < BOARD_SIZE; i++ {
		for j := 0; j < BOARD_SIZE; j++ {
			fmt.Printf("%c ", b.Grid[i][j])
		}
		fmt.Println()
	}
}

func (b *Board) IsFull() bool {
	return b.MovesCount == BOARD_SIZE*BOARD_SIZE
}

func (b *Board) HasWinner() bool {
	//fmt.Println("Checking winner")

	rowsWin := b.checkRows()
	//fmt.Printf("Rows win: %v\n", rowsWin)

	columnsWin := b.checkColumns()
	//fmt.Printf("Columns win: %v\n", columnsWin)

	diagonalsWin := b.checkDiagonals()
	//fmt.Printf("Diagonals win: %v\n", diagonalsWin)

	//fmt.Printf("Combined Win: %v\n", rowsWin || columnsWin || diagonalsWin)

	return rowsWin || columnsWin || diagonalsWin
}

func (b *Board) checkRows() bool {
	for row := 0; row < BOARD_SIZE; row++ {
		char := b.Grid[row][0]
		if char == '-' {
			continue
		}
		for col := 0; col < BOARD_SIZE; col++ {
			if b.Grid[row][col] != char {
				break
			}
			if col == BOARD_SIZE-1 {
				return true
			}
		}
	}
	return false
}

func (b *Board) checkColumns() bool {
	for col := 0; col < BOARD_SIZE; col++ {
		char := b.Grid[0][col]
		if char == '-' {
			continue
		}
		for row := 0; row < BOARD_SIZE; row++ {
			if b.Grid[row][col] != char {
				break
			}
			if row == BOARD_SIZE-1 {
				return true
			}
		}
	}
	return false
}

func (b *Board) checkDiagonals() bool {
	// Check 1st diagonal
	//fmt.Println("Checking 1st diagonal")
	char := b.Grid[0][0]
	if char != '-' {
		for i := 0; i < BOARD_SIZE; i++ {
			//fmt.Printf("Checking 1st diagonal: b.Grid[%d][%d] = %c\n", i, i, b.Grid[i][i])

			if b.Grid[i][i] != char {
				break
			}
			fmt.Printf("i = %d  BOARD_SIZE-1= %d\n", i, BOARD_SIZE-1)
			if i == BOARD_SIZE-1 {
				return true
			}
		}
	}

	// Check 2nd diagonal
	char = b.Grid[0][BOARD_SIZE-1]
	if char != '-' {
		for i := 0; i < BOARD_SIZE; i++ {
			//fmt.Printf("Checking 2nd diagonal: b.Grid[%d][%d] = %c\n", i, BOARD_SIZE-i-1, b.Grid[i][BOARD_SIZE-i-1])

			if b.Grid[i][BOARD_SIZE-i-1] != char {
				break
			}
			if i == BOARD_SIZE-1 {
				return true
			}
		}
	}

	return false
}

type Game struct {
	Player1       *Player
	Player2       *Player
	Board         *Board
	CurrentPlayer *Player
}

func NewGame(player1, player2 *Player) *Game {
	return &Game{
		Player1:       player1,
		Player2:       player2,
		Board:         NewBoard(),
		CurrentPlayer: player1,
	}
}

func (g *Game) switchPlayer() {
	if g.CurrentPlayer == g.Player1 {
		g.CurrentPlayer = g.Player2
	} else {
		g.CurrentPlayer = g.Player1
	}
}

func (g *Game) getValidInput(prompt string) int {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(prompt)
		scanner.Scan()
		input, err := strconv.Atoi(scanner.Text())
		if err == nil && input >= 0 && input <= BOARD_SIZE-1 {
			return input
		}
		fmt.Println(fmt.Sprintf("Invalid input! Please enter a number between 0 and %d.", BOARD_SIZE-1))
	}
}

func (g *Game) Play() {
	n := BOARD_SIZE - 1
	g.Board.PrintBoard()
	for !g.Board.IsFull() {
		if g.Board.HasWinner() {
			break
		}
		fmt.Printf("%s's turn.\n", g.CurrentPlayer.Name)
		row := g.getValidInput(fmt.Sprintf("Enter row (0-%d): ", n))
		col := g.getValidInput(fmt.Sprintf("Enter column (0-%d): ", n))

		err := g.Board.MakeMove(row, col, g.CurrentPlayer.Symbol)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		g.Board.PrintBoard()
		g.switchPlayer()
	}

	if g.Board.HasWinner() {
		g.switchPlayer() // Switch back to the winner
		fmt.Printf("%s wins!\n", g.CurrentPlayer.Name)
	} else {
		fmt.Println("It's a draw!")
	}

}
