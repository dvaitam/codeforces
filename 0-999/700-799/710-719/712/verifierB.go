package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	in  string
	out string
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solveB(s string) int {
	n := len(s)
	if n%2 == 1 {
		return -1
	}
	l0, r0, u0, d0 := 0, 0, 0, 0
	for i := 0; i < n; i++ {
		switch s[i] {
		case 'L':
			l0++
		case 'R':
			r0++
		case 'U':
			u0++
		case 'D':
			d0++
		}
	}
	half := n / 2
	best := 2 * n
	for k := 0; k <= half; k++ {
		m := half - k
		dev := abs(l0-k) + abs(r0-k) + abs(u0-m) + abs(d0-m)
		if dev < best {
			best = dev
		}
	}
	return best / 2
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(2))
	tests := make([]testCase, 100)
	chars := []byte{'L', 'R', 'U', 'D'}
	for t := 0; t < 100; t++ {
		n := rng.Intn(20) + 1
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[i] = chars[rng.Intn(4)]
		}
		s := string(b)
		expect := solveB(s)
		tests[t] = testCase{in: s + "\n", out: fmt.Sprintf("%d\n", expect)}
	}
	return tests
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := strings.TrimSpace(tc.out)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
