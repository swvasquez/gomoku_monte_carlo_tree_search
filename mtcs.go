package monte_carlo_tree_search

import (
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
	children []State
}

type mctsTree struct {
	root     *State
	elements int
}

func updateOrder(s []State) {
	for i := 0; i < 18; i++ {
		if s[i].ucb1 < s[i+1].ucb1 {
			break
		}
		s[i], s[i+1] = s[i+1], s[i]
	}

}
func max(x int, y int) int {
	if x >= y {
		return x
	}
	return y
}

func min(x int, y int) int {
	if x <= y {
		return x
	}
	return y
}

func checkDiagonalDown(s *[19][19]int) int {
	var total int
	for i := 5 - 1; i < 19; i++ {
		for j := 0; j < i-5; i++ {
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

	for j := 1; j < 19-5; j++ {
		for i := 0; i < 19-j-5+1; i++ {
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
	for i := 5 - 1; i < 19; i++ {
		for j := 0; j < i-5; i++ {
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

	for j := 1; j < 19-5; j++ {
		for i := 0; i < 19-j-5+1; i++ {
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

func checkHorizontal(s *[19][19]int) int {
	var total int
	for i := 0; i < 19; i++ {
		for j := 0; j < 15; j++ {
			total = 0
			for k := j; k <= j+5; k++ {
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
	for i := 0; i < 19; i++ {
		for j := 0; j < 15; j++ {
			total = 0
			for k := j; k <= j+5; k++ {
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
func checkFull(heights *[19]int) bool {
	var fillCount int
	full := true
	for i := 0; i <= 19; i++ {
		fillCount = fillCount + heights[i]
	}
	if fillCount == 19*19 {
		full = false
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
	score = score +
		checkHorizontal(&b.board) +
		checkVertical(&b.board) +
		checkDiagonalDown(&b.board) +
		checkDiagonalUp(&b.board)

	if score >= -1 {
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
	root := State{id: 0, visits: 0, ucb1: math.Inf(0)}
	tree := mctsTree{&root, 0}
	return &tree
}

func selectPath(s *State) *State {
	current := s
	for {
		if current.children == nil {
			break
		}
		current = &current.children[0]
	}
	return current
}

func nextMove(b *Board) int {
	var next int
	for {
		rand.Seed(time.Now().Unix())
		next = rand.Int() % len(b.heights)
		if b.heights[next] != 18 {
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

		b.heights[next]++
		b.board[next][b.heights[next]] = player * -1
		b.player = player * -1
	}
	return score[state]
}

func backpropogate(s *State, score int) {
	current := s
	for {
		current.visits++
		current.value = current.value + score
		updateOrder(current.children)
		if current.parent == nil {
			break
		}
		current.ucb1 = ucb1(current)
		current = current.parent
	}
}

func createChildren(current *State, id int) []State {
	boards := make([]Board, 0, 19)
	children := make([]State, 0, 19)

	heights := current.board.heights
	var height int
	nextId := id + 1

	for i := 0; i < 19; i++ {
		height = heights[i]
		if height < 19 {
			boards = append(boards, Board{current.board.board,
				current.board.heights,
				current.board.player * -1})

			boards[i].board[height+1][i] = -1 * current.board.player
			boards[i].heights[i]++
			boards[i].player = -1 * current.board.player
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
	iterations := 100
	//Create root node and initialize MCTree
	tree := mctsInit()

	var leaf *State
	var score int

	for i := 0; i < iterations; i++ {

		leaf = selectPath(tree.root)
		if leaf.visits > 0 || leaf.id == 0 {
			createChildren(leaf, tree.elements)
			tree.elements = tree.elements + 19
			leaf = &leaf.children[0]
		}
		score = expand(*leaf.board)
		backpropogate(leaf, score)

	}
}
