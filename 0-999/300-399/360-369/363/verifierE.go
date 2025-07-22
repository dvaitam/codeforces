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

func runCandidate(bin string, input string) (string, error) {
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

func oracle(input string) (string, error) {
	cmd := exec.Command("go", "run", "363E.go")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("oracle error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) (int, int, int, [][]int) {
	n := rng.Intn(4) + 3
	m := rng.Intn(4) + 3
	maxR := (n - 1) / 2
	if (m-1)/2 < maxR {
		maxR = (m - 1) / 2
	}
	r := rng.Intn(maxR + 1)
	grid := make([][]int, n)
	for i := range grid {
		grid[i] = make([]int, m)
		for j := range grid[i] {
			grid[i][j] = rng.Intn(10) + 1
		}
	}
	return n, m, r, grid
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, r, grid := generateCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, r))
		for ii := 0; ii < n; ii++ {
			for jj := 0; jj < m; jj++ {
				if jj > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprint(grid[ii][jj]))
			}
			sb.WriteByte('\n')
		}
		want, err := oracle(sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed: %v\n", err)
			os.Exit(1)
		}
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
