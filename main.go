package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Board struct {
	board   [19][19]int
	heights [19]int
	player  int
}

type State struct {
	id       int
	visits   int
	value    int
	ucb1     float64
	board    *Board
	parent   *State
	children []State
}

type mctsTree struct {
	root     *State
	elements int
}

func updateOrder(s []State) {
	for i := 0; i < 18; i++ {
		if s[i].ucb1 > s[i+1].ucb1 {
			break
		}
		s[i], s[i+1] = s[i+1], s[i]
	}

}

func displayBoard(b *Board) {
	var top [19]int
	for i := 0; i < 19; i++ {
		top[i] = i % 10
	}
	fmt.Println(top)
	fmt.Print("\n")
	for i := 18; i >= 0; i-- {
		fmt.Println(b.board[i], " ", i)
	}
	fmt.Println("\n")
}

func checkDiagonalDown(s *[19][19]int) int {
	var total int
	for i := 4; i < 19; i++ {
		for j := 0; j < i-3; i++ {
			total = 0
			for k := j; k < j+5; j++ {
				total = total + s[i-k][k]
			}
			if total == 5 {
				return 1
			} else if total == -5 {
				return -1
			}
		}
	}

	for j := 0; j < 15; j++ {
		for i := 0; i < 19-j; i++ {
			total = 0
			for k := i; k < i+5; k++ {
				total = total + s[18-k][j+k]
			}
			if total == 5 {
				return 1
			} else if total == -5 {
				return -1
			}
		}
	}
	return 0
}

func checkDiagonalUp(s *[19][19]int) int {
	var total int
	for i := 0; i < 15; i++ {
		for j := 0; j < 19-i; i++ {
			total = 0
			for k := j; k < j+5; j++ {
				total = total + s[k+i][k]
			}
			if total == 5 {
				return 1
			} else if total == -5 {
				return -1
			}
		}

	}

	for j := 1; j < 15; j++ {
		for i := 0; i < 15-j; i++ {
			total = 0
			for k := i; k < i+5; k++ {
				total = total + s[k][j+k]
			}
			if total == 5 {
				return 1
			} else if total == -5 {
				return -1
			}
		}
	}
	return 0

}

func checkHorizontal(s *[19][19]int) int {
	var total int
	for i := 0; i < 19; i++ {
		for j := 0; j < 15; j++ {
			total = 0
			for k := j; k < j+5; k++ {
				total = total + s[i][k]
			}
			if total == 5 {
				return 1
			} else if total == -5 {
				return -1
			}
		}
	}
	return 0
}

func checkVertical(s *[19][19]int) int {
	var total int
	for j := 0; j < 19; j++ {
		for i := 0; i < 15; i++ {
			total = 0
			for k := i; k < i+5; k++ {
				total = total + s[k][j]
			}
			if total == 5 {
				return 1
			} else if total == -5 {
				return -1
			}
		}
	}
	return 0
}

func checkFull(heights *[19]int) bool {
	var fillCount int
	full := false
	for i := 0; i < 19; i++ {
		fillCount = fillCount + heights[i]
	}
	if fillCount == 19*19 {
		full = true
	}
	return full
}

func checkState(b *Board) string {
	state := "in_play"
	score := 0

	// Check if board is full.
	if checkFull(&b.heights) {
		state = "tie"
	}
	//score = score +
	//	checkHorizontal(&b.board) +
	//	checkVertical(&b.board) +
	//	checkDiagonalDown(&b.board) +
	//	checkDiagonalUp(&b.board)
	if score <= -1 {
		state = "lose"
	} else if score >= 1 {
		state = "win"
	}
	return state
}

func ucb1(s *State) float64 {
	visits := s.visits
	value := s.value
	pVisits := s.parent.visits
	score := float64(value)/float64(visits) +
		2*math.Sqrt(math.Log(float64(pVisits))/float64(visits))
	return score
}

func mctsInit() *mctsTree {
	root_board := Board{player: 1}
	root := State{id: 0, visits: 0, ucb1: math.Inf(0), board: &root_board}
	root.children = createChildren(&root, 1)
	tree := mctsTree{&root, 20}

	return &tree
}

func selectPath(s *State) *State {
	//fmt.Println("new path starting at", s.id)
	current := s

	for {
		current.visits++
		if current.children == nil {
			//fmt.Println("end path at", current.id)
			break
		}
		//fmt.Println("next_node", current.children[0].id)

		current = &current.children[0]

	}
	return current
}

func nextMove(b *Board) int {
	var next int

	for {
		//rand.Seed(time.Now().Unix())
		next = rand.Int() % len(b.heights)
		if b.heights[next] != 19 {
			break
		}
	}
	return next
}

func expand(b Board) int {

	var next int
	var player int
	var state string

	score := make(map[string]int)
	score["win"] = 1
	score["tie"] = 0
	score["lose"] = -1

	for {
		state = checkState(&b)
		if state != "in_play" {
			break
		}
		player = b.player

		next = nextMove(&b)
		b.board[b.heights[next]][next] = (player+2)%2 + 1
		b.player = (player+2)%2 + 1
		b.heights[next]++
	}

	return score[state]
}

func backpropagate(s *State, score int) {
	current := s
	for {
		current.value = current.value + score
		if current.parent == nil {
			break
		}
		current.ucb1 = ucb1(current)
		current = current.parent
		updateOrder(current.children)
	}
}

func createChildren(current *State, id int) []State {
	//fmt.Println("Creating children for node", current.id)
	boards := make([]Board, 0, 19)
	children := make([]State, 0, 19)
	heights := current.board.heights
	var height int
	nextId := id

	nextPlayer := (current.board.player+2)%2 + 1
	for i := 0; i < 19; i++ {
		height = heights[i]
		if height < 19 {
			boards = append(boards, Board{current.board.board,
				current.board.heights,
				nextPlayer})

			boards[i].board[height][i] = nextPlayer
			boards[i].heights[i]++
			boards[i].player = nextPlayer
			//fmt.Println("creating child", nextId, "at position", i)
			children = append(children, State{id: nextId,
				ucb1:   math.Inf(0),
				parent: current,
				board:  &boards[i]})
			nextId++
		}

	}
	return children
}

func main() {

	iterations := 1000
	//Create root node and initialize MCTree
	tree := mctsInit()
	displayBoard(tree.root.board)
	var leaf *State
	var score int

	for i := 0; i < iterations; i++ {
		fmt.Println("\nNew iteration starting at", tree.root.id)
		leaf = selectPath(tree.root)
		if leaf.visits > 1 {
			fmt.Println("Creating children for leaf", leaf.id)
			leaf.children = createChildren(leaf, tree.elements)
			tree.elements = tree.elements + 19
			leaf = &leaf.children[0]
			leaf.visits++
		}
		fmt.Println("Expanding from leaf", leaf.id)
		score = expand(*leaf.board)
		fmt.Println("Backpropagating")
		backpropagate(leaf, score)

	}
	for i := 0; i < 19; i++ {
		fmt.Print(tree.root.children[i].id, " ")
	}
	fmt.Print("\n")
	for i := 0; i < 19; i++ {
		fmt.Print(tree.root.children[i].visits, " ")
	}
}
