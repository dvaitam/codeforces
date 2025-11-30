package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Embedded copy of testcasesD.txt.
const testcasesDData = `
100
43 64
50 31
11 73
36 11
50 67
11 5
42 23
34 54
39 86
14 58
46 89
26 35
2 77
9 51
11 58
37 8
25 13
42 77
26 44
15 66
30 7
31 81
7 36
33 64
36 84
26 62
17 25
15 71
24 22
20 79
10 60
5 10
31 52
37 54
36 13
17 63
15 16
19 20
24 14
9 86
4 87
9 77
36 27
1 6
26 73
45 98
39 64
46 90
7 62
36 46
50 45
7 89
1 32
15 65
20 37
15 3
32 47
33 45
6 11
20 75
28 30
48 49
25 99
10 31
19 27
48 63
43 47
19 51
40 18
50 17
26 47
33 62
15 85
46 49
41 47
28 37
23 53
50 100
46 38
7 63
19 17
29 21
23 33
48 25
22 65
15 16
44 51
25 61
33 61
37 81
15 88
26 66
20 64
15 42
34 89
1 13
31 42
26 31
28 8
37 100
`

type testCase struct {
	n int
	x int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesDData), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no data")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("parse T: %w", err)
	}
	if len(lines)-1 < t {
		return nil, fmt.Errorf("expected %d cases, got %d", t, len(lines)-1)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		fields := strings.Fields(lines[i+1])
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 numbers", i+2)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", i+2, err)
		}
		x, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse x: %w", i+2, err)
		}
		cases = append(cases, testCase{n: n, x: x})
	}
	return cases, nil
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n%2 == 0 {
		return n == 2
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func powMod(a, e, mod int) int {
	res := 1
	a %= mod
	for e > 0 {
		if e&1 != 0 {
			res = int((int64(res) * int64(a)) % int64(mod))
		}
		a = int((int64(a) * int64(a)) % int64(mod))
		e >>= 1
	}
	return res
}

// solve mirrors 303D.go for a single test case.
func solve(tc testCase) int {
	n, x := tc.n, tc.x
	if n == 1 {
		if x > 2 {
			return x - 1
		}
		return -1
	}
	p := n + 1
	if !isPrime(p) {
		return -1
	}
	m := n
	var primes []int
	for i := 2; i*i <= m; i++ {
		if m%i == 0 {
			primes = append(primes, i)
			for m%i == 0 {
				m /= i
			}
		}
	}
	if m > 1 {
		primes = append(primes, m)
	}

	M := x - 1
	kMax := M / p
	rem := M - kMax*p

	isRoot := func(r int) bool {
		if r%p == 0 {
			return false
		}
		for _, q := range primes {
			if powMod(r, n/q, p) == 1 {
				return false
			}
		}
		return true
	}

	// search in kMax block
	if rem > 0 {
		for r := rem; r >= 1; r-- {
			if isRoot(r) {
				return r + kMax*p
			}
		}
	}
	k := kMax - 1
	if k >= 0 {
		for r := p - 1; r >= 1; r-- {
			if isRoot(r) {
				ans := r + k*p
				if ans >= 2 && ans < x {
					return ans
				}
				break
			}
		}
	}
	return -1
}

func runCandidate(bin string, tc testCase) (int, error) {
	input := fmt.Sprintf("%d %d\n", tc.n, tc.x)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	valStr := strings.TrimSpace(out.String())
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("parse output %q: %v", valStr, err)
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
