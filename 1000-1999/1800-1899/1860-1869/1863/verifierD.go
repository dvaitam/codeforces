package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseD struct {
	n    int
	m    int
	grid []string
}

func generateTestsD(num int) []testCaseD {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testCaseD, num)
	for i := 0; i < num; i++ {
		n := rand.Intn(3)*2 + 2
		m := rand.Intn(3)*2 + 2
		board := make([][]byte, n)
		for r := 0; r < n; r++ {
			board[r] = make([]byte, m)
		}
		for r := 0; r < n; r += 2 {
			for c := 0; c < m; c += 2 {
				if rand.Intn(2) == 0 {
					board[r][c] = 'U'
					board[r+1][c] = 'D'
					board[r][c+1] = 'U'
					board[r+1][c+1] = 'D'
				} else {
					board[r][c] = 'L'
					board[r][c+1] = 'R'
					board[r+1][c] = 'L'
					board[r+1][c+1] = 'R'
				}
			}
		}
		grid := make([]string, n)
		for r := 0; r < n; r++ {
			grid[r] = string(board[r])
		}
		tests[i] = testCaseD{n: n, m: m, grid: grid}
	}
	return tests
}

func solveD(tc testCaseD) []string {
	n := tc.n
	m := tc.m
	grid := tc.grid
	hor := make([][]int, m+2)
	ver := make([][]int, n+2)
	for i := 1; i <= n; i++ {
		row := grid[i-1]
		for j := 1; j <= m; j++ {
			c := row[j-1]
			if c == 'U' {
				ver[i] = append(ver[i], j)
			} else if c == 'L' {
				hor[j] = append(hor[j], i)
			}
		}
	}
	board := make([][]byte, n+2)
	for i := range board {
		board[i] = make([]byte, m+2)
	}
	fail := false
	for j := 1; j <= m; j++ {
		if len(hor[j])%2 != 0 {
			fail = true
			break
		}
		for idx, irow := range hor[j] {
			if idx%2 == 0 {
				board[irow][j] = 'B'
				board[irow][j+1] = 'W'
			} else {
				board[irow][j] = 'W'
				board[irow][j+1] = 'B'
			}
		}
	}
	if !fail {
		for i := 1; i <= n; i++ {
			if len(ver[i])%2 != 0 {
				fail = true
				break
			}
			for idx, jcol := range ver[i] {
				if idx%2 == 0 {
					board[i][jcol] = 'B'
					board[i+1][jcol] = 'W'
				} else {
					board[i][jcol] = 'W'
					board[i+1][jcol] = 'B'
				}
			}
		}
	}
	if fail {
		return []string{"-1"}
	}
	res := make([]string, n)
	for i := 1; i <= n; i++ {
		b := make([]byte, m)
		for j := 1; j <= m; j++ {
			b[j-1] = board[i][j]
		}
		res[i-1] = string(b)
	}
	return res
}

func validPainting(tc testCaseD, ans []string) bool {
	n := tc.n
	m := tc.m
	if len(ans) != n {
		return false
	}
	grid := tc.grid
	for i := 0; i < n; i++ {
		if len(ans[i]) != m {
			return false
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			cell := grid[i][j]
			switch cell {
			case 'U':
				if i+1 >= n || grid[i+1][j] != 'D' {
					return false
				}
				a, b := ans[i][j], ans[i+1][j]
				if !((a == 'B' && b == 'W') || (a == 'W' && b == 'B')) {
					return false
				}
			case 'L':
				if j+1 >= m || grid[i][j+1] != 'R' {
					return false
				}
				a, b := ans[i][j], ans[i][j+1]
				if !((a == 'B' && b == 'W') || (a == 'W' && b == 'B')) {
					return false
				}
			case 'D':
				if i == 0 || grid[i-1][j] != 'U' {
					return false
				}
			case 'R':
				if j == 0 || grid[i][j-1] != 'L' {
					return false
				}
			case '.':
				if ans[i][j] != '.' {
					return false
				}
			}
		}
	}
	for i := 0; i < n; i++ {
		bcnt, wcnt := 0, 0
		for j := 0; j < m; j++ {
			if ans[i][j] == 'B' {
				bcnt++
			}
			if ans[i][j] == 'W' {
				wcnt++
			}
		}
		if bcnt != wcnt {
			return false
		}
	}
	for j := 0; j < m; j++ {
		bcnt, wcnt := 0, 0
		for i := 0; i < n; i++ {
			if ans[i][j] == 'B' {
				bcnt++
			}
			if ans[i][j] == 'W' {
				wcnt++
			}
		}
		if bcnt != wcnt {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsD(100)
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
		for _, row := range tc.grid {
			fmt.Fprintln(&input, row)
		}
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("binary execution failed:", err)
		os.Exit(1)
	}
	outputs := strings.Split(strings.TrimSpace(string(out)), "\n")
	idx := 0
	for tnum, tc := range tests {
		expected := solveD(tc)
		if expected[0] == "-1" {
			if idx >= len(outputs) || strings.TrimSpace(outputs[idx]) != "-1" {
				fmt.Printf("test %d failed: expected -1\n", tnum+1)
				os.Exit(1)
			}
			idx++
		} else {
			if idx+tc.n > len(outputs) {
				fmt.Printf("test %d failed: insufficient output lines\n", tnum+1)
				os.Exit(1)
			}
			cand := make([]string, tc.n)
			for i := 0; i < tc.n; i++ {
				cand[i] = strings.TrimSpace(outputs[idx+i])
			}
			if !validPainting(tc, cand) {
				fmt.Printf("test %d produced invalid painting\n", tnum+1)
				os.Exit(1)
			}
			idx += tc.n
		}
	}
	if idx != len(outputs) {
		fmt.Printf("expected %d output lines, got %d\n", idx, len(outputs))
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
