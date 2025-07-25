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

func solveA(n, m, k int, moves [][2]int) int {
	grid := make([][]bool, n+2)
	for i := range grid {
		grid[i] = make([]bool, m+2)
	}
	for t := 1; t <= k; t++ {
		r := moves[t-1][0]
		c := moves[t-1][1]
		grid[r][c] = true
		for dr := -1; dr <= 0; dr++ {
			for dc := -1; dc <= 0; dc++ {
				x := r + dr
				y := c + dc
				if x >= 1 && x+1 <= n && y >= 1 && y+1 <= m {
					if grid[x][y] && grid[x+1][y] && grid[x][y+1] && grid[x+1][y+1] {
						return t
					}
				}
			}
		}
	}
	return 0
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(20) + 1
	m := rng.Intn(20) + 1
	k := rng.Intn(100) + 1
	moves := make([][2]int, k)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < k; i++ {
		r := rng.Intn(n) + 1
		c := rng.Intn(m) + 1
		moves[i] = [2]int{r, c}
		sb.WriteString(fmt.Sprintf("%d %d\n", r, c))
	}
	expect := solveA(n, m, k, moves)
	return sb.String(), expect
}

func runCase(bin, input string, expected int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(outStr)
	if err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		if err := runCase(bin, input, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
