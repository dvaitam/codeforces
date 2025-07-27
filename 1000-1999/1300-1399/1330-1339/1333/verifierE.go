package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input  string
	n      int
	hasSol bool
}

func canRookMove(x1, y1, x2, y2 int) bool {
	return x1 == x2 || y1 == y2
}

func canQueenMove(x1, y1, x2, y2 int) bool {
	if x1 == x2 || y1 == y2 {
		return true
	}
	return abs(x1-x2) == abs(y1-y2)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func computeFee(board [][]int, canMove func(int, int, int, int) bool) int {
	n := len(board)
	pos := make(map[int][2]int, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			pos[board[i][j]] = [2]int{i, j}
		}
	}
	visited := make([]bool, n*n+1)
	cur := 1
	visited[1] = true
	fee := 0
	visitedCount := 1
	for visitedCount < n*n {
		x, y := pos[cur][0], pos[cur][1]
		next := 0
		best := n*n + 1
		for v := 1; v <= n*n; v++ {
			if visited[v] {
				continue
			}
			px, py := pos[v][0], pos[v][1]
			if canMove(x, y, px, py) {
				if v < best {
					best = v
					next = v
				}
			}
		}
		if next == 0 {
			for v := 1; v <= n*n; v++ {
				if !visited[v] {
					next = v
					break
				}
			}
			fee++
		}
		visited[next] = true
		visitedCount++
		cur = next
	}
	return fee
}

func parseBoard(n int, lines []string) ([][]int, error) {
	board := make([][]int, n)
	used := make([]bool, n*n+1)
	for i := 0; i < n; i++ {
		fields := strings.Fields(strings.TrimSpace(lines[i]))
		if len(fields) != n {
			return nil, fmt.Errorf("line %d: expected %d numbers", i+1, n)
		}
		row := make([]int, n)
		for j, f := range fields {
			val, err := strconv.Atoi(f)
			if err != nil || val < 1 || val > n*n || used[val] {
				return nil, fmt.Errorf("invalid number at line %d", i+1)
			}
			used[val] = true
			row[j] = val
		}
		board[i] = row
	}
	for v := 1; v <= n*n; v++ {
		if !used[v] {
			return nil, fmt.Errorf("missing number %d", v)
		}
	}
	return board, nil
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	if outStr == "-1" {
		if tc.hasSol {
			return fmt.Errorf("solution exists but -1 returned")
		}
		return nil
	}
	lines := strings.Split(outStr, "\n")
	if len(lines) != tc.n {
		return fmt.Errorf("expected %d lines got %d", tc.n, len(lines))
	}
	board, err := parseBoard(tc.n, lines)
	if err != nil {
		return err
	}
	rFee := computeFee(board, canRookMove)
	qFee := computeFee(board, canQueenMove)
	if rFee >= qFee {
		return fmt.Errorf("rook fee %d not less than queen %d", rFee, qFee)
	}
	return nil
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(3) + 1 // 1..3
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	has := n >= 3
	return testCase{input: sb.String(), n: n, hasSol: has}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		{input: "1\n", n: 1, hasSol: false},
		{input: "2\n", n: 2, hasSol: false},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
