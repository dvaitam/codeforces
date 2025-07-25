package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	n, m int
	A, B []int
}

func (tc testCase) input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	sb.WriteString(fmt.Sprintf("%d", len(tc.A)))
	for _, v := range tc.A {
		sb.WriteString(fmt.Sprintf(" %d", v))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%d", len(tc.B)))
	for _, v := range tc.B {
		sb.WriteString(fmt.Sprintf(" %d", v))
	}
	sb.WriteString("\n")
	return sb.String()
}

func solve(n, m int, A, B []int) string {
	k := len(A)
	l := len(B)
	sort.Ints(A)
	sort.Ints(B)
	maxA, maxB := 0, 0
	if k > 0 {
		maxA = A[k-1]
	}
	if l > 0 {
		maxB = B[l-1]
	}
	onlyA := make([]int, 0, n*m)
	onlyB := make([]int, 0, n*m)
	both := make([][2]int, 0, n*m)
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			dA := i + j
			dB := i + (m + 1 - j)
			okA := dA <= maxA
			okB := dB <= maxB
			if !okA && !okB {
				return "NO"
			}
			if okA && !okB {
				onlyA = append(onlyA, dA)
			} else if !okA && okB {
				onlyB = append(onlyB, dB)
			} else {
				both = append(both, [2]int{dA, dB})
			}
		}
	}
	if len(onlyA) > k || len(onlyB) > l {
		return "NO"
	}
	sort.Ints(onlyA)
	for i, d := range onlyA {
		if A[i] < d {
			return "NO"
		}
	}
	sort.Ints(onlyB)
	for i, d := range onlyB {
		if B[i] < d {
			return "NO"
		}
	}
	remA := k - len(onlyA)
	Arem := A[len(onlyA):]
	Brem := B[len(onlyB):]
	sort.Slice(both, func(i, j int) bool { return both[i][0] < both[j][0] })
	assignA := make([]int, 0, remA)
	assignB := make([]int, 0, len(both)-remA)
	for idx, db := range both {
		if idx < remA {
			assignA = append(assignA, db[0])
		} else {
			assignB = append(assignB, db[1])
		}
	}
	sort.Ints(assignA)
	for i, d := range assignA {
		if Arem[i] < d {
			return "NO"
		}
	}
	sort.Ints(assignB)
	for i, d := range assignB {
		if Brem[i] < d {
			return "NO"
		}
	}
	return "YES"
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	total := n * m
	k := rng.Intn(total + 1)
	l := total - k
	A := make([]int, k)
	B := make([]int, l)
	for i := 0; i < k; i++ {
		A[i] = rng.Intn(n+m) + 1
	}
	for i := 0; i < l; i++ {
		B[i] = rng.Intn(n+m) + 1
	}
	return testCase{n: n, m: m, A: A, B: B}
}

func deterministicCases() []testCase {
	cases := []testCase{
		{n: 1, m: 1, A: []int{2}, B: []int{}},
		{n: 2, m: 2, A: []int{2, 3}, B: []int{2, 3}},
		{n: 3, m: 3, A: []int{2, 3, 4}, B: []int{2, 3, 4, 5, 6}},
	}
	return cases
}

func runCase(bin string, tc testCase) error {
	input := tc.input()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	result := strings.TrimSpace(out.String())
	expect := solve(tc.n, tc.m, append([]int(nil), tc.A...), append([]int(nil), tc.B...))
	if result != expect {
		return fmt.Errorf("expected %s got %s", expect, result)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := deterministicCases()
	for len(tests) < 100 {
		tests = append(tests, randomCase(rng))
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
