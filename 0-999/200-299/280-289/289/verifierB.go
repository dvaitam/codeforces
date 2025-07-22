package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func expectedAnswer(n, m, d int, grid []int) int64 {
	mod := grid[0] % d
	for _, v := range grid {
		if v%d != mod {
			return -1
		}
	}
	vals := make([]int, len(grid))
	for i, v := range grid {
		vals[i] = v / d
	}
	sort.Ints(vals)
	median := vals[len(vals)/2]
	var moves int64
	for _, v := range vals {
		if v > median {
			moves += int64(v - median)
		} else {
			moves += int64(median - v)
		}
	}
	return moves
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	d := rng.Intn(9) + 1
	total := n * m
	grid := make([]int, total)
	base := rng.Intn(50)
	base = base - base%d
	mod := base % d
	valid := rng.Intn(2) == 0
	for i := 0; i < total; i++ {
		val := rng.Intn(50)
		val = val - val%d + mod
		grid[i] = val
	}
	if !valid {
		idx := rng.Intn(total)
		grid[idx]++ // break mod
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, d)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			fmt.Fprintf(&sb, "%d", grid[i*m+j])
			if j+1 < m {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	expected := expectedAnswer(n, m, d, grid)
	return sb.String(), expected
}

func runCase(exe string, input string, expected int64) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	var got int64
	if _, err := fmt.Sscan(outStr, &got); err != nil {
		return fmt.Errorf("cannot parse output: %v\n%s", err, outStr)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []struct {
		input    string
		expected int64
	}{
		{"1 1 1\n5\n", 0},
		{"2 2 3\n1 4\n7 10\n", 4},
	}
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		cases = append(cases, struct {
			input    string
			expected int64
		}{in, exp})
	}
	for i, tc := range cases {
		if err := runCase(exe, tc.input, tc.expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
