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

func solveC(x, y int64) int64 {
	if x == y {
		return 0
	}
	prev := y
	curr := 2*y - 1
	t := int64(1)
	for curr < x {
		next := curr + prev - 1
		prev = curr
		curr = next
		t++
	}
	return t + 2
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(3))
	// Start with a handful of deterministic edge cases (including the one the
	// user reported) so we are not relying on RNG for coverage.
	seeds := []testCase{
		{in: "10230 14\n", out: fmt.Sprintf("%d\n", solveC(10230, 14))},
		{in: "6 3\n", out: fmt.Sprintf("%d\n", solveC(6, 3))},
		{in: "8 5\n", out: fmt.Sprintf("%d\n", solveC(8, 5))},
		{in: "100000 99999\n", out: fmt.Sprintf("%d\n", solveC(100000, 99999))},
	}

	tests := make([]testCase, 0, 100+len(seeds))
	tests = append(tests, seeds...)

	for i := 0; i < 100; i++ {
		y := int64(rng.Intn(99997) + 3) // 3 .. 99999
		x := int64(rng.Intn(int(100000-y)) + int(y) + 1)
		expect := solveC(x, y)
		in := fmt.Sprintf("%d %d\n", x, y)
		out := fmt.Sprintf("%d\n", expect)
		tests = append(tests, testCase{in: in, out: out})
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
