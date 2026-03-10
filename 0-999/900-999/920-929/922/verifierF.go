package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCaseF struct {
	N int
	K int
}

func maxPairsF(N int) int {
	// Count total valid pairs (a,b) with 1<=a<b<=N and a|b
	A := make([]int, N+2)
	total := 0
	for i := 1; i <= N; i++ {
		for j := i * 2; j <= N; j += i {
			A[j]++
		}
		total += A[i]
	}
	return total
}

func countPairs(elems []int) int {
	if len(elems) == 0 {
		return 0
	}
	set := make(map[int]bool, len(elems))
	maxElem := 0
	for _, v := range elems {
		set[v] = true
		if v > maxElem {
			maxElem = v
		}
	}
	cnt := 0
	for _, a := range elems {
		for b := a * 2; b <= maxElem; b += a {
			if set[b] {
				cnt++
			}
		}
	}
	return cnt
}

func validate(tc testCaseF, got string) error {
	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	answer := strings.TrimSpace(lines[0])

	mp := maxPairsF(tc.N)
	if tc.K > mp {
		if answer != "No" {
			return fmt.Errorf("expected No (impossible), got %q", answer)
		}
		return nil
	}

	if answer != "Yes" {
		return fmt.Errorf("expected Yes, got %q", answer)
	}
	if len(lines) < 3 {
		return fmt.Errorf("Yes but missing size/elements lines")
	}
	m, err := strconv.Atoi(strings.TrimSpace(lines[1]))
	if err != nil {
		return fmt.Errorf("bad size line: %v", err)
	}
	fields := strings.Fields(lines[2])
	if len(fields) != m {
		return fmt.Errorf("declared size %d but got %d elements", m, len(fields))
	}
	elems := make([]int, m)
	seen := make(map[int]bool, m)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("bad element %q: %v", f, err)
		}
		if v < 1 || v > tc.N {
			return fmt.Errorf("element %d out of range [1,%d]", v, tc.N)
		}
		if seen[v] {
			return fmt.Errorf("duplicate element %d", v)
		}
		seen[v] = true
		elems[i] = v
	}
	sort.Ints(elems)
	got_k := countPairs(elems)
	if got_k != tc.K {
		return fmt.Errorf("f(S)=%d but expected %d", got_k, tc.K)
	}
	return nil
}

func genTestsF() []testCaseF {
	rand.Seed(6)
	tests := make([]testCaseF, 0, 100)
	for len(tests) < 100 {
		N := rand.Intn(100) + 2
		maxPairs := N * (N - 1) / 2
		K := rand.Intn(maxPairs+1) + 1
		if rand.Intn(5) == 0 { // sometimes impossible
			K = maxPairs + rand.Intn(100) + 1
		}
		tests = append(tests, testCaseF{N: N, K: K})
	}
	return tests
}

func runCase(bin string, tc testCaseF) error {
	input := fmt.Sprintf("%d %d\n", tc.N, tc.K)
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var outBuf strings.Builder
	var errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return validate(tc, outBuf.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsF()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: N=%d K=%d %v\n", i+1, tc.N, tc.K, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
