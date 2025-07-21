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

type pair struct{ x, y int }

func solve(n, m int, pres byte, grid []string) int {
	seen := make(map[byte]bool)
	dirs := []pair{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] != pres {
				continue
			}
			for _, d := range dirs {
				ni, nj := i+d.x, j+d.y
				if ni >= 0 && ni < n && nj >= 0 && nj < m {
					c := grid[ni][nj]
					if c != '.' && c != pres {
						seen[c] = true
					}
				}
			}
		}
	}
	return len(seen)
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	pres := byte('A' + rng.Intn(26))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %c\n", n, m, pres))
	grid := make([]string, n)
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(5) == 0 {
				row[j] = '.'
			} else {
				row[j] = letters[rng.Intn(len(letters))]
			}
		}
		grid[i] = string(row)
		sb.WriteString(grid[i])
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	sb.WriteByte('\n')
	return sb.String(), solve(n, m, pres, grid)
}

func runCase(exe, input string, expected int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
