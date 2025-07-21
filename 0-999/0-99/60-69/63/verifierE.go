package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

type coord struct{ q, r int }

func grundy(mask int, lines [][]int) int {
	dp := make([]int, mask+1)
	seen := make([]bool, 256)
	gList := make([]int, 0, 32)
	for m := 1; m <= mask; m++ {
		gList = gList[:0]
		for _, line := range lines {
			for i := 0; i < len(line); {
				if (m>>line[i])&1 == 0 {
					i++
					continue
				}
				j := i
				for j < len(line) && (m>>line[j])&1 == 1 {
					j++
				}
				for s := i; s < j; s++ {
					seg := 0
					for e := s; e < j; e++ {
						seg |= 1 << line[e]
						child := m &^ seg
						g := dp[child]
						if !seen[g] {
							seen[g] = true
							gList = append(gList, g)
						}
					}
				}
				i = j
			}
		}
		mex := 0
		for seen[mex] {
			mex++
		}
		dp[m] = mex
		for _, g := range gList {
			seen[g] = false
		}
	}
	return dp[mask]
}

func solve(board []string) string {
	coords := []coord{}
	for r := -2; r <= 2; r++ {
		minq := -2
		if -r-2 > minq {
			minq = -r - 2
		}
		maxq := 2
		if -r+2 < maxq {
			maxq = -r + 2
		}
		for q := minq; q <= maxq; q++ {
			coords = append(coords, coord{q, r})
		}
	}
	idx := make(map[coord]int)
	for i, c := range coords {
		idx[c] = i
	}
	var lines [][]int
	for qv := -2; qv <= 2; qv++ {
		var line []int
		for i, c := range coords {
			if c.q == qv {
				line = append(line, i)
			}
		}
		sort.Slice(line, func(i, j int) bool { return coords[line[i]].r < coords[line[j]].r })
		lines = append(lines, line)
	}
	for rv := -2; rv <= 2; rv++ {
		var line []int
		for i, c := range coords {
			if c.r == rv {
				line = append(line, i)
			}
		}
		sort.Slice(line, func(i, j int) bool { return coords[line[i]].q < coords[line[j]].q })
		lines = append(lines, line)
	}
	for sv := -2; sv <= 2; sv++ {
		var line []int
		for i, c := range coords {
			if -c.q-c.r == sv {
				line = append(line, i)
			}
		}
		sort.Slice(line, func(i, j int) bool { return coords[line[i]].q < coords[line[j]].q })
		lines = append(lines, line)
	}
	mask := 0
	idxOrder := []coord{}
	for r := -2; r <= 2; r++ {
		minq := -2
		if -r-2 > minq {
			minq = -r - 2
		}
		maxq := 2
		if -r+2 < maxq {
			maxq = -r + 2
		}
		for q := minq; q <= maxq; q++ {
			idxOrder = append(idxOrder, coord{q, r})
		}
	}
	for i, c := range idxOrder {
		if board[i] == "O" {
			mask |= 1 << idx[c]
		}
	}
	if grundy(mask, lines) != 0 {
		return "Karlsson"
	}
	return "Lillebror"
}

func generateCases() []testCase {
	rand.Seed(5)
	cases := make([]testCase, 100)
	for t := 0; t < 100; t++ {
		board := make([]string, 19)
		has := false
		for i := 0; i < 19; i++ {
			if rand.Intn(2) == 0 {
				board[i] = "O"
				has = true
			} else {
				board[i] = "."
			}
		}
		if !has {
			board[0] = "O"
		}
		var buf bytes.Buffer
		idx := 0
		layout := []int{3, 4, 5, 4, 3}
		for _, cnt := range layout {
			for j := 0; j < cnt; j++ {
				fmt.Fprintf(&buf, "%s ", board[idx])
				idx++
			}
			if cnt > 0 {
				buf.Truncate(buf.Len() - 1)
			}
			buf.WriteByte('\n')
		}
		expected := solve(board)
		cases[t] = testCase{input: buf.String(), expected: expected}
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE.go <binary>")
		os.Exit(1)
	}
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:%s\nactual:%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
