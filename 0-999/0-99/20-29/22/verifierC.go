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

func solveCase(n, m, v int) string {
	maxExtra := (n-1)*(n-2)/2 + 1
	if m < n-1 || m > maxExtra {
		return "-1"
	}
	edges := 0
	var a, b int
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i == v {
			continue
		}
		a, b = i, v
		edges++
		if edges == n-1 {
			a, b = b, a
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
	}
	for i := 1; i <= n && edges < m; i++ {
		if i == v || i == b {
			continue
		}
		for j := i + 1; j <= n && edges < m; j++ {
			if j == v || j == b {
				continue
			}
			sb.WriteString(fmt.Sprintf("%d %d\n", i, j))
			edges++
		}
	}
	return strings.TrimSpace(sb.String())
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(6) + 3
		maxM := n * (n - 1) / 2
		m := rng.Intn(maxM + 1)
		v := rng.Intn(n) + 1
		input := fmt.Sprintf("%d %d %d\n", n, m, v)
		expected := solveCase(n, m, v)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\ngot:\n%s\ninput:\n%s", t+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
