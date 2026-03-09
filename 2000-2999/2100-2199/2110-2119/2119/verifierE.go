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

type testCase struct {
	n int
	a []int64
	b []int64
}

type dpState struct {
	x  int64
	dp int64
}

// solve returns the minimum number of increments needed so that b[i]&b[i+1]==a[i]
// for all i, or -1 if impossible. Implements the same DP as the reference C++ solution.
func solve(n int, a, b []int64) int64 {
	const inf = int64(1e18)

	// getA(i) returns the C++ a[i] (1-indexed): 0 for i<=0 or i>=n, else a[i-1].
	getA := func(i int) int64 {
		if i <= 0 || i >= n {
			return 0
		}
		return a[i-1]
	}

	f := []dpState{{x: 0, dp: 0}}

	for i := 1; i <= n; i++ {
		ap := getA(i - 1) // a[i-1] in C++: constraint between b[i-1] and b[i]
		an := getA(i)     // a[i]   in C++: constraint between b[i] and b[i+1]
		bi := b[i-1]      // b[i]   in C++ (0-indexed in Go)

		var x int64
		var g []dpState

		for j := 30; j >= -1; j-- {
			y := x | ap | an
			if j != -1 {
				y |= int64(1) << j
			}
			if y >= bi {
				mn := inf
				for _, l := range f {
					if (l.x & y) != ap {
						continue
					}
					cost := l.dp + y - bi
					if cost < mn {
						mn = cost
					}
				}
				if mn < inf {
					g = append(g, dpState{x: y, dp: mn})
				}
			}
			if j != -1 {
				x |= (int64(1) << j) & bi
			}
		}
		f = g
	}

	mn := inf
	for _, l := range f {
		if l.dp < mn {
			mn = l.dp
		}
	}
	if mn < inf {
		return mn
	}
	return -1
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase

	// Fixed cases from problem examples and edge cases.
	tests = append(tests, testCase{n: 4, a: []int64{1, 4, 4}, b: []int64{1, 2, 3, 4}})        // expected 4
	tests = append(tests, testCase{n: 4, a: []int64{4, 0, 4}, b: []int64{1, 1, 1, 1}})        // expected -1
	tests = append(tests, testCase{n: 3, a: []int64{0, 1}, b: []int64{1, 1, 1}})              // expected 1 (NOT impossible)
	tests = append(tests, testCase{n: 2, a: []int64{0}, b: []int64{0, 0}})                    // expected 0
	tests = append(tests, testCase{n: 2, a: []int64{7}, b: []int64{5, 2}})                    // expected 7
	tests = append(tests, testCase{n: 5, a: []int64{3, 0, 5, 1}, b: []int64{1, 2, 3, 4, 5}}) // mixed
	limit := int64(1<<29) - 1
	tests = append(tests, testCase{n: 4, a: []int64{limit, limit, limit}, b: []int64{0, 0, limit, 0}})
	tests = append(tests, testCase{n: 5, a: []int64{limit, limit, limit, limit}, b: []int64{0, 0, 0, 0, 0}})

	for len(tests) < 120 {
		n := rng.Intn(59) + 2
		a := make([]int64, n-1)
		b := make([]int64, n)
		for i := range a {
			a[i] = rng.Int63n(1 << 29)
		}
		for i := range b {
			b[i] = rng.Int63n(1 << 29)
		}
		tests = append(tests, testCase{n: n, a: a, b: b})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	if strings.HasSuffix(candidate, ".go") {
		tmp, err := os.CreateTemp("", "verifierE-bin-*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create temp: %v\n", err)
			os.Exit(1)
		}
		tmp.Close()
		defer os.Remove(tmp.Name())
		out, err := exec.Command("go", "build", "-o", tmp.Name(), candidate).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "compile error: %v\n%s", err, out)
			os.Exit(1)
		}
		candidate = tmp.Name()
	} else if strings.HasSuffix(candidate, ".cpp") || strings.HasSuffix(candidate, ".cc") {
		tmp, err := os.CreateTemp("", "verifierE-bin-*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create temp: %v\n", err)
			os.Exit(1)
		}
		tmp.Close()
		defer os.Remove(tmp.Name())
		out, err := exec.Command("g++", "-O2", "-o", tmp.Name(), candidate).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "compile error: %v\n%s", err, out)
			os.Exit(1)
		}
		candidate = tmp.Name()
	}

	tests := generateTests()
	input := buildInput(tests)

	candOut, err := runBinary(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}

	fields := strings.Fields(candOut)
	if len(fields) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d integers in output, got %d\n", len(tests), len(fields))
		os.Exit(1)
	}

	for i, tc := range tests {
		expected := solve(tc.n, tc.a, tc.b)
		got, err := strconv.ParseInt(fields[i], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid output %q\n", i+1, fields[i])
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d, got %d\nn=%d a=%v b=%v\n",
				i+1, expected, got, tc.n, tc.a, tc.b)
			os.Exit(1)
		}
	}
	fmt.Printf("Accepted (%d tests).\n", len(tests))
}
