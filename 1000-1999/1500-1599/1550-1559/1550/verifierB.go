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

type testCase struct {
	n        int
	a        int
	b        int
	s        string
	expected int
}

func computeExpected(n, a, b int, s string) int {
	if b >= 0 {
		return (a + b) * n
	}
	zeroGroups := 0
	oneGroups := 0
	prev := byte(0)
	for i := 0; i < n; i++ {
		if i == 0 || s[i] != prev {
			if s[i] == '0' {
				zeroGroups++
			} else {
				oneGroups++
			}
			prev = s[i]
		}
	}
	ops := 1 + minInt(zeroGroups, oneGroups)
	return a*n + b*ops
}

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	a := rng.Intn(201) - 100
	b := rng.Intn(201) - 100
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	s := sb.String()
	exp := computeExpected(n, a, b, s)
	return testCase{n: n, a: a, b: b, s: s, expected: exp}
}

func (tc testCase) input() string {
	return fmt.Sprintf("1\n%d %d %d\n%s\n", tc.n, tc.a, tc.b, tc.s)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != tc.expected {
		return fmt.Errorf("expected %d got %d", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
