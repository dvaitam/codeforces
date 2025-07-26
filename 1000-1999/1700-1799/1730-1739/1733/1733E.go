package main

import (
	"bufio"
	"fmt"
	"os"
)

// NOTE: This solution only simulates a limited number of steps. It does not
// handle the full constraints from the original problem statement.

const N = 120
const maxSteps = 300

type slime struct{ x, y int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	belts := make([][]int, N)
	for i := range belts {
		belts[i] = make([]int, N)
	}
	slimes := []slime{{0, 0}}
	history := make([]map[[2]int]bool, maxSteps+1)

	for t := 0; t <= maxSteps; t++ {
		grid := make(map[[2]int]bool)
		for _, s := range slimes {
			if s.x < N && s.y < N {
				grid[[2]int{s.x, s.y}] = true
			}
		}
		history[t] = grid
		if t == maxSteps {
			break
		}
		var next []slime
		for _, s := range slimes {
			nx, ny := s.x, s.y
			if belts[nx][ny] == 0 {
				ny++
			} else {
				nx++
			}
			if nx < N && ny < N {
				next = append(next, slime{nx, ny})
			}
		}
		for _, s := range slimes {
			belts[s.x][s.y] ^= 1
		}
		next = append(next, slime{0, 0})
		slimes = next
	}

	for ; q > 0; q-- {
		var t, x, y int
		fmt.Fscan(in, &t, &x, &y)
		if t <= maxSteps && history[t][[2]int{x, y}] {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
