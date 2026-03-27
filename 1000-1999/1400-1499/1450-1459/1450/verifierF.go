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
	input    string
	expected string
}

// Embedded correct solver for 1450F
func solve(a []int) int {
	n := len(a)
	if n == 1 {
		return 0
	}

	cnt := make([]int, n+1)
	maxCnt := 0
	for i := 0; i < n; i++ {
		cnt[a[i]]++
		if cnt[a[i]] > maxCnt {
			maxCnt = cnt[a[i]]
		}
	}

	if 2*maxCnt > n+1 {
		return -1
	}

	S := 0
	M := make([]int, n+1)
	for i := 0; i < n-1; i++ {
		if a[i] != a[i+1] {
			S++
			M[a[i]]++
			M[a[i+1]]++
		}
	}

	D := 0
	for c := 1; c <= n; c++ {
		if cnt[c] > 0 {
			val := S - M[c] - (n + 1 - 2*cnt[c])
			if val > D {
				D = val
			}
		}
	}

	return n - 1 - S + D
}

func buildCase(a []int) testCase {
	n := len(a)
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	expected := fmt.Sprintf("%d", solve(a))
	return testCase{input: sb.String(), expected: expected}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(n) + 1
	}
	return buildCase(a)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != tc.expected {
		return fmt.Errorf("expected %s got %s", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
