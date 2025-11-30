package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesF.txt so the verifier is self-contained.
const testcasesF = `100
1 6
2 7
3 8
4 9
5 10
6 11
7 12
8 13
9 14
10 15
11 16
12 17
13 18
14 19
15 20
16 21
17 22
18 23
19 24
20 25
21 26
22 27
23 28
24 29
25 30
26 31
27 32
28 33
29 34
30 35
31 36
32 37
33 38
34 39
35 40
36 41
37 42
38 43
39 44
40 45
41 46
42 47
43 48
44 49
45 50
46 51
47 52
48 53
49 54
50 55
1 6
2 7
3 8
4 9
5 10
6 11
7 12
8 13
9 14
10 15
11 16
12 17
13 18
14 19
15 20
16 21
17 22
18 23
19 24
20 25
21 26
22 27
23 28
24 29
25 30
26 31
27 32
28 33
29 34
30 35
31 36
32 37
33 38
34 39
35 40
36 41
37 42
38 43
39 44
40 45
41 46
42 47
43 48
44 49
45 50
46 51
47 52
48 53
49 54
50 55`

// Embedded sieve/solver from 661F.go.
func sieve(n int) []bool {
	prime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		prime[i] = true
	}
	for p := 2; p*p <= n; p++ {
		if prime[p] {
			for j := p * p; j <= n; j += p {
				prime[j] = false
			}
		}
	}
	return prime
}

func countPrimes(primes []bool, l, r int) int {
	cnt := 0
	for x := l; x <= r; x++ {
		if primes[x] {
			cnt++
		}
	}
	return cnt
}

type testCase struct {
	l int
	r int
}

func parseCases() ([]testCase, error) {
	fields := strings.Fields(testcasesF)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	pos := 0
	readInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected end of data")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}
	t, err := readInt()
	if err != nil {
		return nil, fmt.Errorf("bad test count: %v", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		l, err1 := readInt()
		r, err2 := readInt()
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("case %d: bad l or r", i+1)
		}
		cases = append(cases, testCase{l: l, r: r})
	}
	return cases, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	primes := sieve(500)

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&sb, "%d %d\n", tc.l, tc.r)
	}

	outStr, err := runCandidate(bin, sb.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	outFields := strings.Fields(outStr)
	if len(outFields) != len(cases) {
		fmt.Printf("expected %d outputs, got %d\n", len(cases), len(outFields))
		os.Exit(1)
	}
	for i, tc := range cases {
		exp := countPrimes(primes, tc.l, tc.r)
		got, err := strconv.Atoi(outFields[i])
		if err != nil {
			fmt.Printf("case %d: bad output\n", i+1)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed: expected %d got %d\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
