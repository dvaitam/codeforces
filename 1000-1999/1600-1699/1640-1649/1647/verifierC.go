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

func solveCase(n, m int, grid []string) string {
	if grid[0][0] == '1' {
		return "-1\n"
	}
	s := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '1' {
				s++
			}
		}
	}
	var out strings.Builder
	out.WriteString(fmt.Sprintf("%d\n", s))
	for i := n - 1; i >= 1; i-- {
		for j := 0; j < m; j++ {
			if grid[i][j] == '1' {
				fmt.Fprintf(&out, "%d %d %d %d\n", i, j+1, i+1, j+1)
			}
		}
	}
	for j := m - 1; j >= 1; j-- {
		if grid[0][j] == '1' {
			fmt.Fprintf(&out, "1 %d 1 %d\n", j, j+1)
		}
	}
	return out.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	grid := make([]string, n)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				row[j] = '0'
			} else {
				row[j] = '1'
			}
		}
		grid[i] = string(row)
		sb.WriteString(grid[i])
		sb.WriteByte('\n')
	}
	input := sb.String()
	expected := solveCase(n, m, grid)
	return input, expected
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %q got %q", exp, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
