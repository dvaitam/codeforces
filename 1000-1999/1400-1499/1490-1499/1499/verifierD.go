package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// solve mirrors the logic for 1499D: count pairs (a,b) with c*lcm(a,b) - d*gcd(a,b) = x.
func solve(c, d, x int) int {
	ans := 0
	for g := 1; g*g <= x; g++ {
		if x%g != 0 {
			continue
		}
		ans += countPairs(c, d, x, g)
		if g*g != x {
			ans += countPairs(c, d, x, x/g)
		}
	}
	return ans
}

func countPairs(c, d, x, g int) int {
	y := x / g
	if (y+d)%c != 0 {
		return 0
	}
	k := (y + d) / c
	if k <= 0 {
		return 0
	}
	cnt := countDistinctPrimeFactors(k)
	return 1 << cnt
}

func countDistinctPrimeFactors(n int) int {
	cnt := 0
	for p := 2; p*p <= n; p++ {
		if n%p == 0 {
			cnt++
			for n%p == 0 {
				n /= p
			}
		}
	}
	if n > 1 {
		cnt++
	}
	return cnt
}

type testCase struct {
	c, d, x int
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(42))
	var tests []testCase
	// Small manual cases
	tests = append(tests, testCase{2, 2, 24})
	tests = append(tests, testCase{3, 3, 3})
	tests = append(tests, testCase{1, 1, 1})
	tests = append(tests, testCase{1, 1, 10})
	// Random tests
	for i := 0; i < 200; i++ {
		c := rng.Intn(1000) + 1
		d := rng.Intn(1000) + 1
		x := rng.Intn(100000) + 1
		tests = append(tests, testCase{c, d, x})
	}
	return tests
}

func runCandidate(bin string, tests []testCase) ([]string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.c, tc.d, tc.x)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	lines := strings.Fields(strings.TrimSpace(out.String()))
	return lines, nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests := generateTests()
	results, err := runCandidate(bin, tests)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(results) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d results, got %d\n", len(tests), len(results))
		os.Exit(1)
	}
	for i, tc := range tests {
		expected := fmt.Sprintf("%d", solve(tc.c, tc.d, tc.x))
		got := strings.TrimSpace(results[i])
		if got != expected {
			fmt.Printf("test %d failed (c=%d d=%d x=%d)\nexpected: %s\ngot: %s\n", i+1, tc.c, tc.d, tc.x, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
