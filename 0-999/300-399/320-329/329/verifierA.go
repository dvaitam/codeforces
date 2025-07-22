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

type testA struct {
	input  string
	output string
}

func solveA(n int, grid []string) string {
	rowAns := make([][2]int, n)
	okRow := true
	for i := 0; i < n; i++ {
		pos := -1
		for j := 0; j < n; j++ {
			if grid[i][j] == '.' {
				pos = j
				break
			}
		}
		if pos == -1 {
			okRow = false
			break
		}
		rowAns[i] = [2]int{i + 1, pos + 1}
	}
	if okRow {
		var sb strings.Builder
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%d %d\n", rowAns[i][0], rowAns[i][1])
		}
		return strings.TrimSpace(sb.String())
	}
	colAns := make([][2]int, n)
	okCol := true
	for j := 0; j < n; j++ {
		pos := -1
		for i := 0; i < n; i++ {
			if grid[i][j] == '.' {
				pos = i
				break
			}
		}
		if pos == -1 {
			okCol = false
			break
		}
		colAns[j] = [2]int{pos + 1, j + 1}
	}
	if okCol {
		var sb strings.Builder
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%d %d\n", colAns[i][0], colAns[i][1])
		}
		return strings.TrimSpace(sb.String())
	}
	return "-1"
}

func generateCaseA(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	grid := make([]string, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				row[j] = '.'
			} else {
				row[j] = 'E'
			}
		}
		grid[i] = string(row)
		sb.WriteString(grid[i])
		sb.WriteByte('\n')
	}
	out := solveA(n, grid)
	return sb.String(), out
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, got)
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
		in, exp := generateCaseA(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
