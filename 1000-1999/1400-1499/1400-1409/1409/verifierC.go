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

func solve(n int, x, y int64) []int64 {
	diff := y - x
	var step int64
	for i := n - 1; i >= 1; i-- {
		if diff%int64(i) == 0 {
			step = diff / int64(i)
			break
		}
	}
	bestMax := int64(1 << 62)
	var start int64
	for i := 0; i < n; i++ {
		cand := y - int64(i)*step
		if cand <= 0 {
			continue
		}
		maxTerm := cand + int64(n-1)*step
		if maxTerm < bestMax {
			bestMax = maxTerm
			start = cand
		}
	}
	res := make([]int64, n)
	for i := 0; i < n; i++ {
		res[i] = start + int64(i)*step
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(49) + 2
	x := int64(rng.Intn(49) + 1)
	y := int64(rng.Intn(50-int(x)) + int(x) + 1)
	ans := solve(n, x, y)
	input := fmt.Sprintf("1\n%d %d %d\n", n, x, y)
	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return input, sb.String()
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
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
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
