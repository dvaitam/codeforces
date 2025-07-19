package main

import (
	"bufio"
	"fmt"
	"os"
)

// Pair represents a coordinate (row, col)
type Pair struct{ x, y int }

// Step is a sequence of coordinates to rotate
type Step []Pair

var (
	arr   [][]int
	moves [][]int
)

func assert(cond bool) {
	if !cond {
		panic("assertion failed")
	}
}
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// flip transposes each step's coordinates
func flip(p []Step) []Step {
	q := make([]Step, len(p))
	for i, x := range p {
		y := make(Step, len(x))
		for j, a := range x {
			y[j] = Pair{a.y, a.x}
		}
		q[i] = y
	}
	return q
}

// getSwap returns sequence of rotations to swap adjacent cells
func getSwap(x1, y1, x2, y2 int) []Step {
	assert(abs(x1-x2)+abs(y1-y2) == 1)
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	if x1 == x2 {
		// same row: transpose and use swap on columns
		return flip(getSwap(y1, x1, y2, x2))
	}
	var ans []Step
	// x1 < x2
	yz := y1 - 1
	if yz < 0 {
		yz = y1 + 1
	}
	if x1 == 0 {
		ans = append(ans, Step{{x1, y1}, {x2, y1}, {x2 + 1, y1}, {x2 + 1, yz}, {x2, yz}, {x1, yz}})
		ans = append(ans, Step{{x2, yz}, {x2 + 1, yz}, {x2 + 1, y1}, {x2, y1}})
		ans = append(ans, Step{{x1, yz}, {x2, yz}, {x2, y1}, {x1, y1}})
	} else {
		ans = append(ans, Step{{x2, y1}, {x1, y1}, {x1 - 1, y1}, {x1 - 1, yz}, {x1, yz}, {x2, yz}})
		ans = append(ans, Step{{x1, yz}, {x1 - 1, yz}, {x1 - 1, y1}, {x1, y1}})
		ans = append(ans, Step{{x2, yz}, {x1, yz}, {x1, y1}, {x2, y1}})
	}
	return ans
}

// getMove returns rotation to move along same column or row
func getMove(x1, y1, x2, y2, dim int) []Step {
	assert(abs(x1-x2)+abs(y1-y2) == 1)
	if x1 == x2 {
		return flip(getMove(y1, x1, y2, x2, len(arr)))
	}
	var ans []Step
	yz := y1 + 1
	if yz >= dim {
		yz = y1 - 1
	}
	ans = append(ans, Step{{x1, y1}, {x2, y1}, {x2, yz}, {x1, yz}})
	return ans
}

// doMoves applies rotations and records moved values
func doMoves(g []Step) {
	for _, a := range g {
		b := make([]int, len(a))
		for p := range a {
			b[p] = arr[a[p].x][a[p].y]
		}
		moves = append(moves, b)
		// rotate values
		last := arr[a[len(a)-1].x][a[len(a)-1].y]
		for p := len(a) - 1; p > 0; p-- {
			arr[a[p].x][a[p].y] = arr[a[p-1].x][a[p-1].y]
		}
		arr[a[0].x][a[0].y] = last
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, m int
	fmt.Fscan(in, &n, &m)
	arr = make([][]int, n)
	for i := 0; i < n; i++ {
		arr[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &arr[i][j])
		}
	}
	// process each target position
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			num := 1 + i*m + j
			var x, y int
			for a := 0; a < n; a++ {
				for b := 0; b < m; b++ {
					if arr[a][b] == num {
						x, y = a, b
					}
				}
			}
			assert(x >= i)
			if x == i {
				assert(y >= j)
			}
			// move horizontally to correct column
			for y < j {
				if x == n-1 && x <= i+1 {
					doMoves(getSwap(x, y, x, y+1))
				} else {
					doMoves(getMove(x, y, x, y+1, m))
				}
				y++
			}
			for y > j {
				if x == n-1 && x <= i+1 {
					doMoves(getSwap(x, y, x, y-1))
				} else {
					doMoves(getMove(x, y, x, y-1, m))
				}
				y--
			}
			// move vertically
			for x > i+1 {
				doMoves(getMove(x, y, x-1, y, m))
				x--
			}
			for x > i {
				doMoves(getSwap(x, y, x-1, y))
				x--
			}
		}
	}
	// output
	fmt.Fprintln(out, len(moves))
	for _, mv := range moves {
		fmt.Fprintf(out, "%d", len(mv))
		for _, v := range mv {
			fmt.Fprintf(out, " %d", v)
		}
		fmt.Fprintln(out)
	}
}
