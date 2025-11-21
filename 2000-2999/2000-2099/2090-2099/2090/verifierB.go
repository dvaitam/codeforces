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

// Direct simulation of the pushing process for small grids to validate outputs.
func brute(grid [][]byte) bool {
	n := len(grid)
	m := len(grid[0])
	type state struct {
		mask int
	}
	total := 1 << (n * m)
	visited := make([]bool, total)

	encode := func(g [][]byte) int {
		mask := 0
		bit := 1
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if g[i][j] == '1' {
					mask |= bit
				}
				bit <<= 1
			}
		}
		return mask
	}

	decode := func(mask int) [][]byte {
		g := make([][]byte, n)
		bit := 1
		for i := 0; i < n; i++ {
			g[i] = make([]byte, m)
			for j := 0; j < m; j++ {
				if mask&bit != 0 {
					g[i][j] = '1'
				} else {
					g[i][j] = '0'
				}
				bit <<= 1
			}
		}
		return g
	}

	start := make([][]byte, n)
	for i := 0; i < n; i++ {
		start[i] = make([]byte, m)
		for j := 0; j < m; j++ {
			start[i][j] = '0'
		}
	}
	queue := []int{encode(start)}
	visited[queue[0]] = true
	target := encode(grid)

	for len(queue) > 0 {
		mask := queue[0]
		queue = queue[1:]
		if mask == target {
			return true
		}
		g := decode(mask)
		// push rows
		for i := 0; i < n; i++ {
			if g[i][m-1] == '1' {
				continue
			}
			newG := make([][]byte, n)
			for r := range newG {
				newG[r] = make([]byte, m)
				copy(newG[r], g[r])
			}
			carry := byte('1')
			for j := 0; j < m; j++ {
				if newG[i][j] == '0' {
					newG[i][j], carry = carry, '0'
					break
				} else {
					newG[i][j], carry = carry, newG[i][j]
				}
			}
			if carry == '1' {
				continue
			}
			newMask := encode(newG)
			if !visited[newMask] {
				visited[newMask] = true
				queue = append(queue, newMask)
			}
		}
		// push cols
		for j := 0; j < m; j++ {
			if g[n-1][j] == '1' {
				continue
			}
			newG := make([][]byte, n)
			for r := range newG {
				newG[r] = make([]byte, m)
				copy(newG[r], g[r])
			}
			carry := byte('1')
			for i := 0; i < n; i++ {
				if newG[i][j] == '0' {
					newG[i][j], carry = carry, '0'
					break
				} else {
					newG[i][j], carry = carry, newG[i][j]
				}
			}
			if carry == '1' {
				continue
			}
			newMask := encode(newG)
			if !visited[newMask] {
				visited[newMask] = true
				queue = append(queue, newMask)
			}
		}
	}
	return false
}

func solveReference(grid [][]byte) bool {
	n := len(grid)
	m := len(grid[0])
	if n*m <= 16 {
		return brute(grid)
	}
	rowFull := make([]bool, n)
	colFull := make([]bool, m)
	for i := 0; i < n; i++ {
		full := true
		for j := 0; j < m; j++ {
			if grid[i][j] == '0' {
				full = false
				break
			}
		}
		rowFull[i] = full
	}
	for j := 0; j < m; j++ {
		full := true
		for i := 0; i < n; i++ {
			if grid[i][j] == '0' {
				full = false
				break
			}
		}
		colFull[j] = full
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '1' && !rowFull[i] && !colFull[j] {
				return false
			}
		}
	}
	return true
}

type testCase struct {
	n, m int
	grid [][]byte
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 1
	m := rng.Intn(6) + 1
	if n*m > 16 && rng.Intn(3) == 0 {
		n = 4
		m = 4
	}
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				grid[i][j] = '0'
			} else {
				grid[i][j] = '1'
			}
		}
	}
	return testCase{n: n, m: m, grid: grid}
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for i := 0; i < tc.n; i++ {
			sb.Write(tc.grid[i])
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/2090B_binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	for i := 0; i < 200; i++ {
		cases = append(cases, genCase(rng))
	}
	input := buildInput(cases)
	output, err := runCandidate(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run candidate: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Fields(output)
	if len(lines) < len(cases) {
		fmt.Fprintf(os.Stderr, "not enough outputs: got %d expected %d\n", len(lines), len(cases))
		os.Exit(1)
	}
	for i, tc := range cases {
		ans := strings.ToLower(lines[i])
		expected := solveReference(tc.grid)
		want := "yes"
		if !expected {
			want = "no"
		}
		if (ans == "yes" || ans == "y" || ans == "yes\n") != expected {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %s got %s\n", i+1, want, lines[i])
			fmt.Fprintf(os.Stderr, "grid %dx%d:\n", tc.n, tc.m)
			for _, row := range tc.grid {
				fmt.Fprintf(os.Stderr, "%s\n", string(row))
			}
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
