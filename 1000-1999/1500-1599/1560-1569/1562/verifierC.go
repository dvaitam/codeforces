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

func solveCase(n int, s string) string {
	idx := -1
	for i, ch := range s {
		if ch == '0' {
			idx = i + 1
		}
	}
	mid := n / 2
	if idx > mid {
		return fmt.Sprintf("1 %d 1 %d\n", idx, idx-1)
	} else if idx < mid {
		return fmt.Sprintf("%d %d %d %d\n", mid+1, n, mid, n-1)
	} else {
		return fmt.Sprintf("%d %d %d %d\n", mid+1, n, mid, n)
	}
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(18) + 2
	sb := make([]byte, n)
	for i := range sb {
		if rng.Intn(2) == 0 {
			sb[i] = '0'
		} else {
			sb[i] = '1'
		}
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
