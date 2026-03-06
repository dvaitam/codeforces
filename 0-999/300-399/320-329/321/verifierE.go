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

func solveCase(n, k int, u [][]int) int {
	cost := make([][]int, n)
	for i := range cost {
		cost[i] = make([]int, n)
	}
	for l := 0; l < n; l++ {
		s := 0
		for r := l; r < n; r++ {
			for i := l; i < r; i++ {
				s += u[i][r]
			}
			cost[l][r] = s
		}
	}

	const INF = 1 << 60
	prev := make([]int, n)
	for i := 0; i < n; i++ {
		prev[i] = cost[0][i]
	}
	curr := make([]int, n)
	for j := 1; j < k; j++ {
		for i := 0; i < n; i++ {
			curr[i] = INF
		}
		for i := j; i < n; i++ {
			for p := j - 1; p < i; p++ {
				val := prev[p] + cost[p+1][i]
				if val < curr[i] {
					curr[i] = val
				}
			}
		}
		copy(prev, curr)
	}
	return prev[n-1]
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	k := rng.Intn(n) + 1
	if k > 4 {
		k = 4
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	u := make([][]int, n)
	for i := 0; i < n; i++ {
		u[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if j < i {
				u[i][j] = u[j][i]
			} else if j > i {
				u[i][j] = rng.Intn(10)
			}
			sb.WriteString(fmt.Sprintf("%d", u[i][j]))
			if j+1 < n {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	expected := strconv.Itoa(solveCase(n, k, u))
	return sb.String(), expected
}

func run(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expected := generateCase(rng)
		got, err := run(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expected, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
