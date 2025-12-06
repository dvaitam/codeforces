package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var (
	n        int
	solution [][2]int
	found    bool
)

func main() {
	// Configure fast I/O
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	// Read n
	if !scanner.Scan() {
		return
	}
	n, _ = strconv.Atoi(scanner.Text())

	// Read permutation and pad with 0 at start and n+1 at end
	// p indices will be 0 to n+1
	p := make([]int, n+2)
	p[0] = 0
	for i := 1; i <= n; i++ {
		scanner.Scan()
		val, _ := strconv.Atoi(scanner.Text())
		p[i] = val
	}
	p[n+1] = n + 1

	// Start DFS
	solve(0, p)
}

func solve(depth int, p []int) {
	if found {
		return
	}

	// Calculate number of "jumps" (non-adjacent values neighbors)
	// A jump is where absolute difference != 1
	jumps := 0
	for i := 0; i <= n; i++ {
		if abs(p[i+1]-p[i]) != 1 {
			jumps++
		}
	}

	// If no jumps, check if it is sorted (strictly increasing)
	if jumps == 0 {
		isSorted := true
		for i := 0; i <= n; i++ {
			if p[i+1]-p[i] != 1 {
				isSorted = false
				break
			}
		}
		if isSorted {
			found = true
			fmt.Println(depth)
			// Print commands in reverse order of discovery
			// because we discovered operations from Target -> Sorted
			// so the operations from Sorted -> Target are the reverse
			for i := len(solution) - 1; i >= 0; i-- {
				fmt.Printf("%d %d\n", solution[i][0], solution[i][1])
			}
			return
		}
	}

	// Limit depth to 3
	if depth == 3 {
		return
	}

	// Pruning: each operation can remove at most 2 jumps.
	// We need to reduce jumps to 0 in (3 - depth) moves.
	if jumps > 2*(3-depth) {
		return
	}

	// Identify candidate cut points
	// A candidate is an index i such that the adjacency property changes
	// or it is a jump boundary.
	candidates := make([]int, 0, 32)
	for i := 0; i <= n; i++ {
		isCand := false
		d := p[i+1] - p[i]
		
		// If it's a jump, it's a candidate
		if abs(d) != 1 {
			isCand = true
		} else {
			// If it's 1 or -1, check if the "run" changes
			if i > 0 {
				dPrev := p[i] - p[i-1]
				if d != dPrev {
					isCand = true
				}
			} else {
				// i=0. If starts with -1, it's a boundary of a -1 run.
				if d == -1 {
					isCand = true
				}
			}
		}
		if isCand {
			candidates = append(candidates, i)
		}
	}

	// Try all pairs of candidates as reversal boundaries
	for i := 0; i < len(candidates); i++ {
		u := candidates[i]
		for j := i + 1; j < len(candidates); j++ {
			v := candidates[j]
			
			// Reversing a single element (u+1 == v) does nothing, skip
			if v == u+1 {
				continue
			}

			// Perform reversal on range [u+1, v]
			reverse(p, u+1, v)
			solution = append(solution, [2]int{u + 1, v})

			solve(depth+1, p)
			if found {
				return
			}

			// Backtrack
			solution = solution[:len(solution)-1]
			reverse(p, u+1, v)
		}
	}
}

func reverse(p []int, l, r int) {
	for l < r {
		p[l], p[r] = p[r], p[l]
		l++
		r--
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
