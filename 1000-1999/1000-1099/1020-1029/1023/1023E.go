package main

import (
	"fmt"
)

// q sends a query to the interactor: "Is there a path from (a, b) to (c, d)?"
// Returns true if the answer is "YES".
func q(a, b, c, d int) bool {
	fmt.Printf("? %d %d %d %d\n", a, b, c, d)
	var s string
	fmt.Scan(&s)
	return len(s) > 0 && s[0] == 'Y'
}

func main() {
	var N int
	if _, err := fmt.Scan(&N); err != nil {
		return
	}

	var a string
	x, y := 1, 1

	// Phase 1: Forward from (1, 1) to the anti-diagonal (x + y == N + 1)
	for x+y <= N {
		// Check if we can reach the goal (N, N) if we move Down to (x+1, y)
		if q(x+1, y, N, N) {
			a += "D"
			x++
		} else {
			// If not, we must move Right
			a += "R"
			y++
		}
	}

	var bBytes []byte
	x, y = N, N

	// Phase 2: Backward from (N, N) to the anti-diagonal
	for x+y > N+1 {
		// Check if we can reach (x, y-1) from the start (1, 1)
		// If YES, the move leading to (x, y) was 'R' (from Left)
		if q(1, 1, x, y-1) {
			bBytes = append(bBytes, 'R')
			y--
		} else {
			// If NO, the move leading to (x, y) was 'D' (from Up)
			bBytes = append(bBytes, 'D')
			x--
		}
	}

	// Reverse the backward path string since we constructed it from end to start
	for i, j := 0, len(bBytes)-1; i < j; i, j = i+1, j-1 {
		bBytes[i], bBytes[j] = bBytes[j], bBytes[i]
	}

	// Output the combined path
	fmt.Printf("! %s%s\n", a, string(bBytes))
}
