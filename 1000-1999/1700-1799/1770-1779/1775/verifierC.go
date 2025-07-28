package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func expected(n, x int64) int64 {
	if x&^n != 0 {
		return -1
	}
	const LIMIT int64 = 1 << 60
	lower := n
	upper := LIMIT
	for b := 0; b < 61; b++ {
		bit := int64(1) << uint(b)
		if n&bit != 0 {
			nextZero := (n | (bit - 1)) + 1
			if x&bit == 0 {
				if nextZero > lower {
					lower = nextZero
				}
			} else {
				if nextZero-1 < upper {
					upper = nextZero - 1
				}
			}
		} else {
			if x&bit != 0 {
				return -1
			}
		}
	}
	if lower <= upper {
		return lower
	}
	return -1
}

func generateTests() [][2]int64 {
	rand.Seed(2)
	t := 100
	tests := make([][2]int64, t)
	for i := 0; i < t; i++ {
		n := rand.Int63n(100000) + 1
		x := rand.Int63n(n + 1)
		tests[i] = [2]int64{n, x}
	}
	return tests
}

func verifyCase(bin string, n, x int64) error {
	exp := expected(n, x)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(fmt.Sprintf("1\n%d %d\n", n, x))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("execution error: %v", err)
	}
	outStr := strings.TrimSpace(string(out))
	got, err := strconv.ParseInt(outStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid integer output: %q", outStr)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		if err := verifyCase(bin, tc[0], tc[1]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (n=%d x=%d)\n", i+1, err, tc[0], tc[1])
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
