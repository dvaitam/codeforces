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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solve(a []int) string {
	n := len(a)
	pos := make([]int, n+1)
	for i, v := range a {
		pos[v] = i
	}
	L := make([]int, n)
	R := make([]int, n)
	stack := []int{}
	for i := 0; i < n; i++ {
		for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			L[i] = -1
		} else {
			L[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	stack = stack[:0]
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			R[i] = n
		} else {
			R[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	span := make([]int, n+1)
	for x := 1; x <= n; x++ {
		idx := pos[x]
		span[x] = R[idx] - L[idx] - 1
	}
	minSpan := make([]int, n+2)
	minSpan[0] = n + 1
	for i := 1; i <= n; i++ {
		if span[i] < minSpan[i-1] {
			minSpan[i] = span[i]
		} else {
			minSpan[i] = minSpan[i-1]
		}
	}
	res := make([]byte, n)
	for k := 1; k <= n; k++ {
		m := n - k + 1
		if minSpan[m] >= k {
			res[k-1] = '1'
		} else {
			res[k-1] = '0'
		}
	}
	return string(res)
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
	expected := solve(a)
	return testCase{input: sb.String(), expected: expected}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 1
	perm := rng.Perm(n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = perm[i] + 1
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
