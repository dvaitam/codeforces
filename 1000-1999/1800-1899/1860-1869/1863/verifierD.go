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
	chars := []byte{'U', 'D', 'L', 'R'}
	for i := 0; i < num; i++ {
		n := rand.Intn(6) + 1
		m := rand.Intn(6) + 1
		grid := make([]string, n)
		for r := 0; r < n; r++ {
			b := make([]byte, m)
			for c := 0; c < m; c++ {
				b[c] = chars[rand.Intn(len(chars))]
			}
			grid[r] = string(b)
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
			for i := 0; i < tc.n; i++ {
				if idx >= len(outputs) || strings.TrimSpace(outputs[idx]) != expected[i] {
					fmt.Printf("test %d line %d mismatch\n", tnum+1, i+1)
					os.Exit(1)
				}
				idx++
			}
		}
	}
	if idx != len(outputs) {
		fmt.Printf("expected %d output lines, got %d\n", idx, len(outputs))
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
