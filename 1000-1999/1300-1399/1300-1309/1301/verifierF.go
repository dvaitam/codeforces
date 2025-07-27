package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleF")
	cmd := exec.Command("go", "build", "-o", oracle, "1301F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(4) + 1
	m := r.Intn(4) + 1
	maxColors := n * m
	k := r.Intn(min(5, maxColors)) + 1
	grid := make([][]int, n)
	for i := range grid {
		grid[i] = make([]int, m)
	}
	// ensure each color appears at least once
	for c := 1; c <= k; c++ {
		i := r.Intn(n)
		j := r.Intn(m)
		grid[i][j] = c
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 0 {
				grid[i][j] = r.Intn(k) + 1
			}
		}
	}
	q := r.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", grid[i][j]))
		}
		sb.WriteByte('\n')
	}
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		r1 := r.Intn(n) + 1
		c1 := r.Intn(m) + 1
		r2 := r.Intn(n) + 1
		c2 := r.Intn(m) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", r1, c1, r2, c2))
	}
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func run(cmdPath, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(cmdPath, ".go") {
		cmd = exec.Command("go", "run", cmdPath)
	} else if strings.HasSuffix(cmdPath, ".py") {
		cmd = exec.Command("python3", cmdPath)
	} else {
		cmd = exec.Command(cmdPath)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
