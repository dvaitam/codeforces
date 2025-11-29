package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesC.txt to remove external dependency.
const testcasesRaw = `244 10
558 3
379 10
486 10
68 10
14 8
266 9
240 4
735 8
554 9
488 7
655 3
238 3
889 9
400 1
688 2
164 10
44 5
799 1
844 5
485 10
737 7
732 7
405 10
456 3
900 6
100 1
140 8
223 5
989 7
798 5
432 9
854 7
588 6
547 10
418 10
238 6
699 1
877 5
621 3
716 6
988 9
927 10
583 2
731 4
649 10
274 5
128 2
494 8
91 6
820 2
421 3
21 5
438 7
894 2
46 10
630 1
387 10
339 9
903 5
518 4
37 5
8 2
111 10
549 1
972 4
995 7
299 10
270 3
707 1
889 6
322 6
982 3
919 7
386 8
891 9
396 10
698 9
106 10
997 9
278 7
650 4
959 5
448 5
534 5
562 6
12 7
594 6
21 7
631 10
648 3
62 6
478 6
696 6
624 5
756 8
23 10
63 1
987 6
258 8`

const mod int64 = 1000000007

type testCase struct {
	n int
	m int
}

// referenceSolution embeds the logic from 1288C.go so no external oracle is needed.
func referenceSolution(n, m int) int64 {
	maxN := n + 2*m
	fact := make([]int64, maxN+1)
	invFact := make([]int64, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[maxN] = modPow(fact[maxN], mod-2)
	for i := maxN; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
	comb := func(n, k int) int64 {
		if k < 0 || k > n {
			return 0
		}
		return fact[n] * invFact[k] % mod * invFact[n-k] % mod
	}
	return comb(n+2*m-1, 2*m)
}

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func parseTestcases() []testCase {
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scanner.Split(bufio.ScanWords)
	tests := make([]testCase, 0)
	for {
		if !scanner.Scan() {
			break
		}
		n, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic("invalid n")
		}
		if !scanner.Scan() {
			panic("unexpected end of testcases")
		}
		m, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic("invalid m")
		}
		tests = append(tests, testCase{n: n, m: m})
	}
	return tests
}

func runBin(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	cand := os.Args[1]

	tests := parseTestcases()
	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.m)
		exp := referenceSolution(tc.n, tc.m)
		gotStr, err := runBin(cand, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(gotStr, 10, 64)
		if err != nil {
			fmt.Printf("Test %d: cannot parse output %q: %v\n", i+1, gotStr, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("Test %d failed: input %d %d expected %d got %d\n", i+1, tc.n, tc.m, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
