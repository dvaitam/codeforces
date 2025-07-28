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
	for i := 0; i < n-1; i++ {
		for j := 0; j < m-1; j++ {
			cnt := 0
			if grid[i][j] == '1' {
				cnt++
			}
			if grid[i][j+1] == '1' {
				cnt++
			}
			if grid[i+1][j] == '1' {
				cnt++
			}
			if grid[i+1][j+1] == '1' {
				cnt++
			}
			if cnt == 3 {
				return "NO\n"
			}
		}
	}
	return "YES\n"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
