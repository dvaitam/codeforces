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

type gridCase struct {
	n, m   int
	lines  []string
	expect []string
}

func generateCase(rng *rand.Rand) gridCase {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	// create grid with '.'
	grid := make([][]byte, n)
	for i := range grid {
		grid[i] = make([]byte, m)
		for j := 0; j < m; j++ {
			grid[i][j] = '.'
		}
	}
	// ensure at least one star
	starCount := rng.Intn(n*m) + 1
	for k := 0; k < starCount; k++ {
		i := rng.Intn(n)
		j := rng.Intn(m)
		grid[i][j] = '*'
	}
	// convert grid to lines
	lines := make([]string, n)
	for i := range grid {
		lines[i] = string(grid[i])
	}
	// compute bounding rectangle
	minR, maxR := n-1, 0
	minC, maxC := m-1, 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '*' {
				if i < minR {
					minR = i
				}
				if i > maxR {
					maxR = i
				}
				if j < minC {
					minC = j
				}
				if j > maxC {
					maxC = j
				}
			}
		}
	}
	expect := make([]string, maxR-minR+1)
	for i := minR; i <= maxR; i++ {
		expect[i-minR] = string(grid[i][minC : maxC+1])
	}
	return gridCase{n, m, lines, expect}
}

func runCase(bin string, c gridCase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", c.n, c.m))
	for _, line := range c.lines {
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(outLines) != len(c.expect) {
		return fmt.Errorf("expected %d lines got %d", len(c.expect), len(outLines))
	}
	for i, exp := range c.expect {
		if strings.TrimSpace(outLines[i]) != exp {
			return fmt.Errorf("line %d: expected %s got %s", i+1, exp, outLines[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		c := generateCase(rng)
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%v\n", i+1, err, c)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
