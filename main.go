// A simple script that runs a Monte Carlo Tree Search on a 19 by 19 Gomoku board (five-in-a-row). Currently, moves are
// made similar to Connect 4, where pieces move to the lowest available position in a given column. This is just to
// reduce the complexity while developing the code.

package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"time"
)

type Board struct {
	board   [19][19]int
	heights [19]int
	player  int
	last    [2]int
}

type State struct {
	id       int
	visits   int
	value    int
	preval   int
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

type Packet struct {
	board  Board
	visits [19]int
}

// Displays the current board state. Players' moves are represented
// by 1s and 7s for visual simplicity.
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

// Checks if five consecutive pieces are found horizontally,
// vertically, or along a diagonal.
func checkBoard(s *[19][19]int) int {
	var memo [21][21][2][4]int
	row := -1
	val := 0
	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			val = s[i][j]
			row = (val % 5) - 1
			if row != -1 {
				memo[i+1][j+1][row][0] = 1 + memo[i+1][j][row][0]
				memo[i+1][j+1][row][1] = 1 + memo[i][j+1][row][1]
				memo[i+1][j+1][row][2] = 1 + memo[i][j][row][2]
				memo[i+1][j+1][row][3] = 1 + memo[i][j+2][row][3]
				for k := 0; k < 4; k++ {
					if memo[i+1][j+1][row][k] == 5 {
						return (-val + 4) / 3
					}
				}
			}
		}
	}
	return val
}

// Checks if the board is full.
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

// Declares whether a game is in play or finished.
func checkState(b *Board) string {
	state := "in_play"
	score := 0

	// Check if board is full.
	if checkFull(&b.heights) {
		state = "tie"
	}
	score = checkBoard(&b.board)

	if score == -1 {
		state = "lose"
	} else if score == 1 {
		state = "win"
	}
	return state
}

// The UCB1 score is used to determine whether exploration or
// exploitation should guide decision making.
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

func mctsInit(b *Board) *mctsTree {
	rootState := State{
		id: 0, visits: 0,
		ucb1:     math.Inf(0),
		board:    b,
		children: make([]*State, 0, 19),
	}
	tree := mctsTree{root: &rootState}
	return &tree
}

// This chooses a path, starting from the root and ending on a leaf.
// This path is determined by the UCB1 scores of a node's children.
func selectPath(tree *mctsTree) (*State, int) {
	current := tree.root
	depth := 0
	for {
		current.visits++
		if current.children == nil {
			break
		} else if len(current.children) < 19 {
			current = createChild(current, tree)
		} else {
			ucb1Update(current)
			current = ucb1Max(current)
		}
		depth++
	}
	return current, depth
}

// Chooses a random move consistent with the constraints of the board.
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

// When a path ends on an unvisited leaf, it preforms a simulation
// from that node, terminating when the game has ended. It returns
// a 1 if player 1 one, -1 if player 1 lost, or 0 if the
// game ended in a tie.
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

// Updates the values of all the nodes along the current path with
// the outcome from the above expansion.
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
	nextBoard := Board{board: current.board.board,
		heights: current.board.heights,
		player:  nextPlayer}
	nextBoard.board[height][count] = nextPlayer
	nextBoard.heights[count]++
	nextBoard.player = nextPlayer
	nextBoard.last = [2]int{height, count}
	nextState := State{id: nextId,
		ucb1:   math.Inf(0),
		parent: current,
		board:  &nextBoard,
		preval: 2}
	current.children = append(current.children, &nextState)
	tree.elements++
	return &nextState
}

// This goroutine runs one cycle of MTCS for a specified time limit.
func run(runtime int, c chan Packet) {
	b := <-c
	tree := mctsInit(&b.board)
	var leaf *State
	var score int
	var depth int

	// Selects path from root to leaf.
	start := time.Now()
	for {
		leaf, depth = selectPath(tree)
		if leaf.visits > 3 {
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

		end := time.Now()
		diff := end.Sub(start).Seconds()
		if diff > float64(runtime) {
			fmt.Println("tree depth", tree.depth)
			fmt.Println("nodes", tree.elements)
			break
		}
	}
	var max_visits int
	var best_move *Board
	for i, val := range tree.root.children {
		b.visits[i] = val.visits
		if b.visits[i] > max_visits {
			max_visits = b.visits[i]
			best_move = val.board
		}
	}
	b.board = *best_move
	c <- b
}

// This aggregates the results from many "run" goroutines and
// updates the board.
func play(b *Board, channels int, runtime int) {
	chans := make([]chan Packet, channels, 12)
	outputs := make([]Packet, channels, 12)

	for i := range chans {
		chans[i] = make(chan Packet)
	}
	// Create root node and initialize MCTree.
	// Set another initial value other than the empty
	// board here.
	root := Packet{board: *b}
	for _, channel := range chans {
		go run(runtime, channel)
		channel <- root
	}

	for i, channel := range chans {
		outputs[i] = <-channel
	}

	// Display results: node ids and visit counts.
	// The node with the max value corresponds to the column
	// the next move should take place at.
	var output [19]float32
	var max float32
	var next Board
	for i := 0; i < 19; i++ {
		for _, packet := range outputs {
			output[i] += float32(packet.visits[i]) / 3

			if output[i] > max {
				max = output[i]
				next = packet.board

			}
		}
	}
	displayBoard(&next.board)
	fmt.Println(next.last[0], next.last[1])
	b.board[next.last[0]][next.last[1]] = next.player
	b.heights[next.last[1]] += 1
	b.player = next.player

	for i := 0; i < 19; i++ {
		output[i] = output[i] / max
	}
	fmt.Printf("\n%v\n", "RESULTS")
	fmt.Println("N  V")
	for idx, val := range output {
		fmt.Printf("%02d %f\n", idx, val)
	}
}

func main() {
	// This pits two MTCS algorithms against one another. Each "player"
	// can use a different number of goroutines. This is done to see
	// how parallelization affects the decision process.

	// Define how long you want to run the algorithm.
	rand.Seed(time.Now().Unix())
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Create root node and initialize MCTree.
	// Set another initial value other than the empty
	// board here.
	var initState [19][19]int
	initPlayer := 1
	gameBoard := Board{board: initState, player: initPlayer}
	moveTime := 20
	p1Channels := 1
	p2Channels := 1

	for {
		play(&gameBoard, p1Channels, moveTime)
		if checkState(&gameBoard) != "in_play" {
			break
		}
		play(&gameBoard, p2Channels, moveTime)
		if checkState(&gameBoard) != "in_play" {
			break
		}

	}
}
