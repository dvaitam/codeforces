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

func expected(n int, row1, row2, b []int) int {
	prefix1 := make([]int, n)
	for i := 1; i < n; i++ {
		prefix1[i] = prefix1[i-1] + row1[i-1]
	}
	suffix2 := make([]int, n+1)
	for i := n - 1; i >= 1; i-- {
		suffix2[i] = suffix2[i+1] + row2[i-1]
	}
	const inf = int(1e9)
	minTotal := inf
	for j1 := 1; j1 <= n; j1++ {
		timeThere := suffix2[j1] + b[j1-1] + prefix1[j1-1]
		for j2 := 1; j2 <= n; j2++ {
			if j1 == j2 {
				continue
			}
			timeBack := prefix1[j2-1] + b[j2-1] + suffix2[j2]
			total := timeThere + timeBack
			if total < minTotal {
				minTotal = total
			}
		}
	}
	return minTotal
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(49) + 2
	row1 := make([]int, n-1)
	row2 := make([]int, n-1)
	b := make([]int, n)
	for i := range row1 {
		row1[i] = rng.Intn(100) + 1
	}
	for i := range row2 {
		row2[i] = rng.Intn(100) + 1
	}
	for i := range b {
		b[i] = rng.Intn(100) + 1
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, v := range row1 {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range row2 {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	exp := fmt.Sprintf("%d", expected(n, row1, row2, b))
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input, exp := genCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %s got %s\ninput:\n%s", t+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
