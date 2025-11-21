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
	N      int
	A, B   string
	expect string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := buildTests()
	input := renderInput(tests)

	out, err := run(bin, input)
	if err != nil {
		fmt.Printf("runtime error: %v\ninput:\n%s\n", err, input)
		os.Exit(1)
	}
	gotLines := strings.Split(strings.TrimSpace(out), "\n")
	if len(gotLines) != len(tests) {
		fmt.Printf("expected %d lines but got %d\ninput:\n%s\nactual:\n%s\n", len(tests), len(gotLines), input, out)
		os.Exit(1)
	}
	for i, tc := range tests {
		got := strings.TrimSpace(gotLines[i])
		if got != tc.expect {
			fmt.Printf("case %d failed: expected %s got %s\nN=%d\nA=%s\nB=%s\nfull input:\n%s\n", i+1, tc.expect, got, tc.N, tc.A, tc.B, input)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func renderInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n%s\n%s\n", tc.N, tc.A, tc.B)
	}
	return sb.String()
}

func buildTests() []testCase {
	rand.Seed(time.Now().UnixNano())
	var tests []testCase

	// Exhaustive small cases for N=3 and N=4
	for n := 3; n <= 4; n++ {
		states := allStates(n)
		for _, a := range states {
			for _, b := range states {
				tests = append(tests, makeCase(n, a, b))
			}
		}
	}

	// Hand-crafted edge cases
	tests = append(tests,
		makeCase(5, "00000", "00000"),
		makeCase(5, "00000", "10000"),
		makeCase(5, "01000", "00001"),
		makeCase(6, "100001", "001100"),
		makeCase(6, "010000", "111111"),
	)

	// Random cases with small N so brute force stays fast
	for i := 0; i < 200; i++ {
		n := 3 + rand.Intn(6) // 3..8
		a := randomState(n)
		b := randomState(n)
		tests = append(tests, makeCase(n, a, b))
	}

	return tests
}

func makeCase(n int, a, b string) testCase {
	exp := "NO"
	if reachable(a, b) {
		exp = "YES"
	}
	return testCase{
		N:      n,
		A:      a,
		B:      b,
		expect: exp,
	}
}

func reachable(A, B string) bool {
	n := len(A)
	if n != len(B) {
		return false
	}
	start := 0
	target := 0
	for i := 0; i < n; i++ {
		if A[i] == '1' {
			start |= 1 << i
		}
		if B[i] == '1' {
			target |= 1 << i
		}
	}
	if start == target {
		return true
	}
	seen := make(map[int]struct{})
	queue := []int{start}
	seen[start] = struct{}{}
	for len(queue) > 0 {
		mask := queue[0]
		queue = queue[1:]
		for i := 0; i < n; i++ {
			if mask&(1<<i) == 0 {
				continue
			}
			next := mask
			if i > 0 {
				next ^= 1 << (i - 1)
			}
			if i+1 < n {
				next ^= 1 << (i + 1)
			}
			if _, ok := seen[next]; ok {
				continue
			}
			if next == target {
				return true
			}
			seen[next] = struct{}{}
			queue = append(queue, next)
		}
	}
	return false
}

func allStates(n int) []string {
	total := 1 << n
	states := make([]string, 0, total)
	for mask := 0; mask < total; mask++ {
		var sb strings.Builder
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
		}
		states = append(states, sb.String())
	}
	return states
}

func randomState(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rand.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return sb.String()
}
