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

func maxPerimeter(grid [][]int) int {
	n := len(grid)
	m := len(grid[0])
	maxP := 0
	for r1 := 0; r1 < n; r1++ {
		for r2 := r1; r2 < n; r2++ {
			curr := 0
			for c := 0; c < m; c++ {
				zero := true
				for r := r1; r <= r2; r++ {
					if grid[r][c] != 0 {
						zero = false
						break
					}
				}
				if zero {
					curr++
					p := 2 * (curr + r2 - r1 + 1)
					if p > maxP {
						maxP = p
					}
				} else {
					curr = 0
				}
			}
		}
	}
	return maxP
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		grid := make([][]int, n)
		for i := 0; i < n; i++ {
			grid[i] = make([]int, m)
			for j := 0; j < m; j++ {
				grid[i][j] = rng.Intn(2)
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] == 0 {
					sb.WriteByte('0')
				} else {
					sb.WriteByte('1')
				}
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		expected := fmt.Sprintf("%d", maxPerimeter(grid))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", t+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
