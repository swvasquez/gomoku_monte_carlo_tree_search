// A simple script that runs a Monte Carlo Tree Search
// on a 19 by 19 Gomoku board (five-in-a-row).
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
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
	children []*State
}

type mctsTree struct {
	root     *State
	elements int
	depth    int
}

func updateOrder(s []*State) {
	for i := len(s) - 1; i > 0; i-- {
		if s[i].ucb1 > s[i-1].ucb1 {
			break
		}
		s[i], s[i-1] = s[i-1], s[i]
	}

}

// Not used in code, but append use when
// you want to view the board state.
func displayBoard(b *[19][19]int) {
	var top [19]int
	for i := 0; i < 19; i++ {
		top[i] = i % 10
	}
	fmt.Printf("\n%v\n\n", top)
	for i := 18; i >= 0; i-- {
		fmt.Printf("%v  %d\n", b[i], i)
	}
	fmt.Printf("\n\n")
}

func checkDiagonalDown(s *[19][19]int) int {
	var total int
	for i := 4; i < 19; i++ {
		for j := 0; j < i-3; j++ {
			total = 0
			for k := j; k < j+5; k++ {
				total = total + s[i-k][k]
			}
			if total == 5 {
				return 1
			} else if total == 35 {
				return -1
			}
		}
	}

	for j := 0; j < 15; j++ {
		for i := 0; i < 15-j; i++ {
			total = 0
			for k := i; k < i+5; k++ {
				total = total + s[18-k][j+k]
			}
			if total == 5 {
				return 1
			} else if total == 35 {
				return -1
			}
		}
	}
	return 0
}

func checkDiagonalUp(s *[19][19]int) int {
	var total int
	for i := 0; i < 15; i++ {
		for j := 0; j < 15-i; j++ {
			total = 0
			for k := j; k < j+5; k++ {
				total = total + s[k+i][k]
			}
			if total == 5 {
				return 1
			} else if total == 35 {
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
			} else if total == 35 {
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
			} else if total == 35 {
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
			} else if total == 35 {
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

	for {
		score = checkHorizontal(&b.board)
		if score == 1 || score == -1 {
			break
		}
		score = checkVertical(&b.board)
		if score == 1 || score == -1 {
			break
		}
		score = checkDiagonalDown(&b.board)
		if score == 1 || score == -1 {
			break
		}
		score = checkDiagonalUp(&b.board)
		if score == 1 || score == -1 {
			break
		}
		break
	}

	if score == -1 {
		state = "lose"
	} else if score == 1 {
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

func ucb1Update(s *State) {
	for _, val := range s.children {
		val.ucb1 = ucb1(val)
	}
}

func ucb1Max(s *State) *State {
	maxState := s.children[0]
	max := maxState.ucb1

	for _, val := range s.children {
		if val.ucb1 > max {
			maxState = val
		}
	}
	return maxState
}

func mctsInit(root *[19][19]int, player int) *mctsTree {
	rootBoard := Board{}
	rootBoard.board = *root
	rootBoard.player = player
	rootState := State{id: 0, visits: 0, ucb1: math.Inf(0), board: &rootBoard, children: make([]*State, 0, 19)}
	tree := mctsTree{root: &rootState}

	return &tree
}

func selectPath(tree *mctsTree) (*State, int) {
	current := tree.root
	depth := 0
	for {
		current.visits++
		children := len(current.children)
		if current.children == nil {
			break
		}
		if children < 19 {
			current = createChild(current, tree)
		} else {
			ucb1Update(current)
			current = ucb1Max(current)
		}
		depth++
	}
	return current, depth
}

func nextMove(b *Board) int {
	var next int

	for {
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

		b.board[b.heights[next]][next] = 8 - player
		b.player = 8 - player
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
		current = current.parent
	}
}

func createChild(current *State, tree *mctsTree) *State {
	if current.children == nil {
		current.children = make([]*State, 0, 19)
	}
	count := len(current.children)
	nextPlayer := 8 - current.board.player
	var height int
	for i := 0; i < 18; i++ {
		height = current.board.heights[count]
		if height < 18 {
			break
		}
		count++
	}

	nextId := current.id*19 + count + 1
	nextBoard := Board{current.board.board,
		current.board.heights,
		nextPlayer}
	nextBoard.board[height][count] = nextPlayer
	nextBoard.heights[count]++
	nextBoard.player = nextPlayer
	nextState := State{id: nextId,
		ucb1:   math.Inf(0),
		parent: current,
		board:  &nextBoard}
	current.children = append(current.children, &nextState)
	tree.elements++
	return &nextState
}
func expansion(s *State) {

}

func main() {

	// Define how long you want to run the algorithm.
	iterations := 100000
	rand.Seed(time.Now().Unix())

	// Create root node and initialize MCTree.
	// Set another initial value other than the empty
	// board here.
	var init_board [19][19]int
	init_player := 1
	tree := mctsInit(&init_board, init_player)
	var leaf *State
	var score int
	var depth int

	// Selects path from root to leaf.
	for i := 0; i < iterations; i++ {
		leaf, depth = selectPath(tree)
		if leaf.visits > 1 {
			createChild(leaf, tree)
			leaf = leaf.children[0]
			depth++
			leaf.visits++
		}

		// Expands the leaf node by playing randomly.
		score = expand(*leaf.board)

		// Backpropagate results all the way up to root.
		backpropagate(leaf, score)

		if depth > tree.depth {
			tree.depth = depth
		}
	}

	// Display results: node ids and visit counts.
	// The node with the max value corresponds to the column
	// the next move should take place at.
	var output []int
	for _, val := range tree.root.children {
		output = append(output, val.visits)
	}
	fmt.Printf("\n%d Nodes, %d Max Depth", tree.elements, tree.depth)
	fmt.Printf("\n%v\n", "RESULTS")
	fmt.Println("N  V")
	for idx, val := range output {
		fmt.Printf("%02d %d\n", idx, val)
	}
}
