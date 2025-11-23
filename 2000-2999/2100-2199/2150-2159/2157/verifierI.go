package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	binaryPath := os.Args[1]

	rand.Seed(time.Now().UnixNano())

	// We will generate a batch of tests
	numTests := 50
	type TestCase struct {
		n, m int
	}
	tests := make([]TestCase, numTests)

	for i := 0; i < numTests; i++ {
		// Constrain N and M to be small enough for DP
		// m up to 50
		// n up to 1000
		m := rand.Intn(49) + 2 // 2 to 50
		n := rand.Intn(1000) + 1 // 1 to 1000
		tests[i] = TestCase{n, m}
	}

	// Prepare input for binary
	var inputBuf bytes.Buffer
	fmt.Fprintf(&inputBuf, "%d\n", numTests)
	for _, t := range tests {
		fmt.Fprintf(&inputBuf, "%d %d\n", t.n, t.m)
	}

	// Run binary
	cmd := exec.Command(binaryPath)
	cmd.Stdin = &inputBuf
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Binary failed to run: %v\n", err)
		os.Exit(1)
	}

	// Parse output
	scanner := bufio.NewScanner(&outBuf)
	for i, tCase := range tests {
		if !scanner.Scan() {
			fmt.Printf("Test %d: Missing output\n", i)
			os.Exit(1)
		}
		outLine := strings.TrimSpace(scanner.Text())
		outLine = strings.ToUpper(outLine)

		expected := solve(tCase.n, tCase.m)
		expectedStr := "NO"
		if expected {
			expectedStr = "YES"
		}

		if outLine != expectedStr {
			fmt.Printf("Test %d FAILED\n", i)
			fmt.Printf("Input: n=%d, m=%d\n", tCase.n, tCase.m)
			fmt.Printf("Expected: %s\n", expectedStr)
			fmt.Printf("Actual:   %s\n", outLine)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed!\n", numTests)
}

// cache for DP: key = (n << 10) | blocked
// max n = 1000, max blocked = 50. 
// key space roughly 1000 * 64 < 2^16.
// We reset cache for each test case or include m in key?
// Easier to just reallocate or clear map for each test case since they are small.

var memo map[int]bool

func solve(n, m int) bool {
	memo = make(map[int]bool)
	return canWin(n, m, 0) // 0 means no blocked move
}

func canWin(n, m, blocked int) bool {
	// If we can make health <= 0 immediately, we win.
	// We can pick any x in [1, m], x != blocked.
	// To win immediately, we need n - x <= 0 => x >= n.
	// So if there is a valid x >= n, we win.
	// Valid x are [1, m] excluding blocked.
	
	// Best move for immediate win is to pick largest possible x.
	// But we just need *any* x >= n.
	// Does there exist x in [1, m], x != blocked such that x >= n?
	
	// Range of valid moves: [1, m] \ {blocked}.
	// If n <= 0, game is already over (previous player won), so this state is Lose.
	// But per problem, we check if current player makes move to <= 0.
	
	if n <= 0 {
		// Should not happen if called correctly, but implies we already lost.
		return false
	}

	// Check immediate win
	// If we can play x >= n
	if m >= n {
		// We can definitely reach <= 0 unless the ONLY winning move is blocked.
		// Winning moves: x in [n, m].
		// If blocked is not in [n, m], we can use any x in [n, m].
		// If blocked is in [n, m], we can use any other x in [n, m].
		// So we win if [n, m] has size > 1 OR (size == 1 and blocked != the single element).
		
		countWinning := m - n + 1
		if countWinning > 1 {
			return true
		}
		if countWinning == 1 {
			// The only winning move is 'n' (if n=m) or generally x=n (if n <= m).
			// Wait, range is [n, m].
			// If count == 1, then n == m.
			// The move is x = n.
			if blocked != n {
				return true
			}
			// If blocked == n, we cannot take the immediate win.
			// But we might still win by playing something else < n.
			// So we don't return false yet, we just can't win *immediately*.
		}
	}

	// Memoization key
	key := (n << 10) | blocked
	if res, ok := memo[key]; ok {
		return res
	}

	// Try all possible moves
	canForceWin := false
	for x := 1; x <= m; x++ {
		if x == blocked {
			continue
		}
		
		// If this move wins immediately, we already handled it?
		// Not fully. If n <= m and blocked == n (and n==m), we fell through.
		if n - x <= 0 {
			canForceWin = true
			break
		}

		// Otherwise, recursive check
		// If opponent cannot win from (n-x) with blocked x, then we win.
		if !canWin(n-x, m, x) {
			canForceWin = true
			break
		}
	}

	memo[key] = canForceWin
	return canForceWin
}