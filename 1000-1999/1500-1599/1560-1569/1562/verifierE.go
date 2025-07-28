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

type sub struct{ i, j int }

func solveCase(n int, s string) string {
	arr := []byte(s)
	lcp := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		lcp[i] = make([]int, n+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			if arr[i] == arr[j] {
				lcp[i][j] = lcp[i+1][j+1] + 1
			}
		}
	}
	less := func(x, y sub) bool {
		len1 := x.j - x.i + 1
		len2 := y.j - y.i + 1
		l := lcp[x.i][y.i]
		if l >= len1 || l >= len2 {
			if len1 == len2 {
				return false
			}
			return len1 < len2
		}
		return arr[x.i+l] < arr[y.i+l]
	}
	tails := make([]sub, 0)
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			cur := sub{i, j}
			lo, hi := 0, len(tails)
			for lo < hi {
				mid := (lo + hi) / 2
				if less(tails[mid], cur) {
					lo = mid + 1
				} else {
					hi = mid
				}
			}
			if lo == len(tails) {
				tails = append(tails, cur)
			} else {
				tails[lo] = cur
			}
		}
	}
	return fmt.Sprintf("%d\n", len(tails))
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2 // 2..6
	sb := make([]byte, n)
	for i := range sb {
		sb[i] = byte(rng.Intn(26) + 'a')
	}
	s := string(sb)
	input := fmt.Sprintf("1\n%d\n%s\n", n, s)
	expected := solveCase(n, s)
	return input, expected
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
