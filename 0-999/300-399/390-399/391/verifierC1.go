package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC1.txt so the verifier is self-contained.
const testcasesC1 = `14 7
3 1
0 0
2 4
4 2
3 4
3 5
2 4
3 1
2 2
4 0
1 5
4 0
1 5
2 2
1
---
8 9
0 3
2 3
3 4
2 2
4 0
5 4
1 0
4 0
0
---
12 7
5 1
5 4
5 1
0 1
4 1
3 4
2 3
1 0
5 0
2 2
5 4
0 3
2
---
2 2
4 5
2 0
5
---
9 6
4 3
1 2
4 4
4 1
4 2
2 1
3 1
0 1
4 0
2
---
10 11
2 4
3 5
0 3
0 5
5 4
1 2
1 4
0 1
0 1
5 5
0
---
10 7
4 4
2 0
3 3
3 4
5 5
5 2
5 1
2 1
0 0
2 5
2
---
5 1
5 3
1 0
2 0
1 1
2 5
9
---
4 1
4 5
5 0
4 0
4 0
-1
---
11 4
4 1
4 5
0 0
3 3
0 1
2 5
0 0
0 5
4 0
0 4
1 3
1
---
10 2
2 4
0 3
1 0
0 4
5 0
2 5
2 3
3 1
1 2
0 2
3
---
15 16
4 5
1 1
5 0
5 3
1 5
0 3
5 4
1 4
1 2
2 5
4 2
2 3
0 5
4 2
3 1
0
---
9 1
3 1
5 3
0 2
2 4
5 2
0 5
4 2
2 4
1 5
10
---
14 10
1 1
5 1
2 5
3 3
5 3
3 5
5 5
0 4
0 3
4 0
1 3
5 5
2 4
1 3
2
---
13 11
5 0
0 0
1 3
3 2
0 2
2 3
5 0
1 3
3 1
4 4
3 5
4 0
4 5
0
---
3 1
3 2
5 0
3 1
-1
---
1 1
5 4
-1
---
10 2
1 3
0 3
4 5
5 0
1 0
2 2
2 3
5 0
1 2
0 1
3
---
11 9
5 2
5 4
2 2
0 2
1 4
2 0
0 5
0 5
0 4
1 5
5 3
2
---
12 11
3 1
5 2
3 2
2 2
4 2
1 5
1 0
3 2
4 4
2 0
0 0
1 2
0
---
3 2
4 3
2 4
2 1
5
---
5 1
3 1
5 4
1 5
0 0
2 2
12
---
7 6
2 2
3 2
0 0
0 3
4 0
3 5
3 3
0
---
7 1
2 4
2 3
5 5
5 0
1 0
1 0
5 1
6
---
12 4
0 5
3 1
0 3
0 0
3 2
4 1
4 2
2 4
3 1
3 3
4 1
5 2
4
---
2 1
5 4
0 3
-1
---
13 11
1 3
0 3
3 2
3 4
2 1
1 0
5 4
0 3
1 1
4 4
1 3
0 5
1 5
1
---
7 8
5 4
2 4
3 1
3 0
5 2
5 5
1 5
0
---
15 11
2 4
0 0
4 0
1 1
2 1
4 0
1 2
3 0
4 3
2 2
5 1
5 1
3 5
0 3
0 2
0
---
9 7
4 5
4 3
0 1
4 2
0 4
5 4
4 3
4 3
0 3
3
---
10 10
1 2
0 0
5 3
0 2
5 2
1 3
2 3
4 3
4 0
2 1
0
---
11 3
1 5
2 0
5 5
2 5
0 1
0 2
3 3
3 0
1 4
3 3
4 3
6
---
8 1
0 5
3 2
1 0
0 5
0 4
4 5
4 2
1 1
9
---
7 8
0 2
0 5
4 0
5 5
3 5
4 4
5 2
0
---
13 14
1 3
3 2
2 5
5 0
5 0
0 0
4 3
1 2
0 0
3 4
3 5
2 4
1 4
0
---
15 7
1 0
0 3
5 0
4 0
5 3
4 1
3 5
4 0
2 3
0 3
1 5
3 3
4 0
3 3
0 2
0
---
12 5
0 2
2 0
0 4
0 1
2 1
5 0
0 1
2 5
2 5
1 5
0 0
1 5
0
---
1 2
2 5
0
---
1 1
1 1
1
---
8 2
3 2
2 3
5 2
2 2
1 4
0 5
1 2
2 1
7
---
10 2
0 4
4 1
4 1
2 1
1 2
3 2
1 3
1 5
1 0
2 1
3
---
10 1
3 1
0 4
5 3
0 2
1 5
3 2
2 0
4 1
3 3
5 5
10
---
6 6
4 5
0 2
3 3
3 1
0 2
3 2
2
---
8 2
1 5
0 3
2 3
0 1
4 3
4 3
5 1
1 1
8
---
8 6
4 1
1 2
2 0
3 5
5 3
5 4
0 3
3 0
1
---
4 3
0 4
5 1
2 3
4 5
4
---
13 8
0 0
5 2
3 3
1 4
4 4
3 5
2 3
5 3
0 0
0 5
2 2
5 2
0 2
2
---
14 13
5 0
1 3
4 2
0 2
0 4
0 1
2 4
2 4
4 5
2 1
0 4
4 0
1 3
1 5
0
---
9 7
5 0
2 5
2 0
3 3
2 3
4 4
5 3
1 1
1 4
1
---
11 5
2 3
5 5
3 0
5 3
3 2
1 1
3 5
3 3
5 0
4 4
2 1
4
---
1 2
3 3
0
---
15 1
4 5
0 5
5 2
0 5
5 1
5 1
5 4
3 2
0 1
2 0
0 3
0 3
4 5
0 2
2 3
9
---
3 3
3 3
5 1
5 3
7
---
15 5
4 1
2 4
1 5
1 0
1 3
3 0
2 1
3 5
3 0
3 1
3 2
5 0
1 5
1 5
3 5
1
---
3 2
4 0
2 5
4 4
-1
---
13 5
2 3
3 3
3 2
0 5
5 3
5 1
4 2
5 2
5 0
4 0
5 0
3 1
4 2
2
---
1 2
5 0
0
---
3 1
3 4
5 3
1 4
-1
---
15 8
3 0
1 1
2 2
4 0
0 1
4 4
0 2
2 1
2 1
4 2
3 5
2 5
2 5
0 2
4 2
2
---
9 7
2 0
3 3
2 3
5 5
1 3
0 0
2 3
4 3
0 3
0
---
8 2
0 4
0 0
1 3
0 4
1 3
3 0
1 1
3 1
1
---
8 4
1 4
2 4
2 1
2 5
3 2
0 3
4 4
2 4
7
---
8 9
5 0
2 1
2 5
1 3
0 1
0 1
4 2
5 5
0
---
11 1
5 4
4 2
4 1
3 5
0 5
0 5
3 4
5 5
2 2
0 1
5 5
17
---
14 15
1 0
0 4
4 4
2 2
5 2
2 4
1 5
2 0
5 1
3 3
2 1
4 1
1 4
1 0
0
---
3 2
3 1
4 1
5 5
-1
---
3 2
3 4
2 1
4 5
10
---
8 2
4 0
0 1
4 4
4 5
2 1
3 5
3 3
2 5
9
---
8 9
2 1
5 2
0 2
2 3
1 0
4 1
3 0
5 2
0
---
12 12
1 2
2 0
3 2
5 4
5 4
5 0
5 5
1 5
2 1
2 2
5 0
1 3
0
---
11 1
0 3
1 1
0 1
0 3
1 3
5 1
2 2
0 2
0 1
2 5
5 2
6
---
12 13
5 0
3 4
3 0
0 3
3 3
2 1
4 0
4 4
5 1
5 4
5 5
3 4
0
---
14 12
5 0
1 3
0 4
2 0
1 4
1 5
1 3
0 5
1 1
5 4
5 4
4 0
1 2
5 2
0
---
13 12
5 4
3 2
0 1
5 5
0 3
3 2
0 2
2 3
2 3
5 0
1 2
3 4
1 2
1
---
9 8
0 5
4 0
4 3
4 2
2 3
5 0
0 0
3 2
3 5
0
---
1 2
5 5
0
---
10 5
5 2
1 1
4 4
4 5
2 5
3 4
2 5
1 2
0 1
4 1
5
---
6 7
5 2
0 2
4 2
0 5
4 2
1 2
0
---
8 7
4 1
3 2
1 0
0 5
1 3
4 5
0 5
3 2
1
---
12 10
1 4
0 0
4 4
3 5
2 4
1 4
4 3
1 3
0 2
5 2
5 2
1 0
0
---
7 5
2 5
4 3
4 1
4 0
2 4
2 1
1 5
5
---
12 13
4 4
3 1
2 5
3 1
0 5
4 2
3 5
5 5
4 3
1 0
0 5
2 5
0
---
4 5
1 0
2 4
1 3
3 1
0
---
12 12
3 2
4 3
5 2
0 4
1 2
1 0
3 2
0 0
3 3
5 0
3 2
2 1
0
---
13 7
3 4
3 5
5 2
5 2
5 0
3 3
5 3
0 4
4 1
3 0
2 0
4 4
1 2
3
---
6 1
0 3
1 0
5 2
0 0
5 2
4 3
7
---
15 5
5 5
2 5
3 5
5 5
1 0
4 5
1 4
3 0
2 4
4 5
0 5
1 5
1 5
1 3
3 1
6
---
6 7
4 1
0 2
0 2
4 5
1 5
2 2
0
---
13 5
4 2
3 0
1 1
3 5
4 3
1 4
0 1
5 3
4 4
1 1
5 5
4 2
0 5
6
---
3 4
4 0
0 4
4 2
0
---
1 1
0 3
3
---
10 6
4 0
3 1
2 5
3 1
4 0
2 2
5 4
0 1
1 0
2 3
1
---
11 1
5 3
2 2
3 4
0 4
4 2
5 5
1 2
2 4
0 3
3 3
0 1
13
---
1 2
2 1
0
---
7 2
2 3
0 4
0 5
0 5
2 0
0 2
1 3
3
---
13 9
5 0
5 2
3 2
3 0
0 3
0 0
4 1
3 4
2 5
0 2
5 0
4 0
0 3
0
---
6 2
4 5
5 5
1 1
1 4
1 2
1 0
8
---
8 7
0 2
1 1
1 1
5 3
2 3
2 2
1 4
1 1
2
---
7 1
1 1
4 0
0 0
5 2
3 4
5 0
1 0
2
---
1 1
4 4
-1`

// Embedded solver from 391C1.go.
func solve(n, k int, p, e []int) int64 {
	const INF = 1 << 60
	ans := INF
	totalMasks := 1 << n
	for mask := 0; mask < totalMasks; mask++ {
		wins := bits.OnesCount(uint(mask))
		effort := 0
		for i := 0; i < n; i++ {
			if (mask>>i)&1 == 1 {
				effort += e[i]
			}
		}
		if effort >= ans {
			continue
		}
		ahead := 0
		for i := 0; i < n; i++ {
			pi := p[i]
			if (mask>>i)&1 == 0 {
				pi++
			}
			if pi > wins || (pi == wins && (mask>>i)&1 == 0) {
				ahead++
			}
		}
		if ahead+1 <= k {
			ans = effort
		}
	}
	if ans == INF {
		return -1
	}
	return int64(ans)
}

type testCase struct {
	n, k int
	p    []int
	e    []int
}

func parseCases() ([]testCase, error) {
	parts := strings.Split(strings.TrimSpace(testcasesC1), "\n---\n")
	cases := make([]testCase, 0, len(parts))
	for idx, block := range parts {
		lines := strings.Split(strings.TrimSpace(block), "\n")
		if len(lines) == 0 {
			continue
		}
		var n, k int
		if _, err := fmt.Sscan(lines[0], &n, &k); err != nil {
			return nil, fmt.Errorf("case %d: bad n k", idx+1)
		}
		p := make([]int, n)
		e := make([]int, n)
		for i := 0; i < n; i++ {
			if _, err := fmt.Sscan(lines[i+1], &p[i], &e[i]); err != nil {
				return nil, fmt.Errorf("case %d: bad line %d", idx+1, i+2)
			}
		}
		cases = append(cases, testCase{n: n, k: k, p: p, e: e})
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
	for i := 0; i < tc.n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", tc.p[i], tc.e[i])
	}
	return sb.String()
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
		fmt.Fprintln(os.Stderr, "usage: verifierC1 /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Println("failed to load testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solve(tc.n, tc.k, tc.p, tc.e)
		input := buildInput(tc)
		gotStr, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(gotStr, 10, 64)
		if err != nil {
			fmt.Printf("case %d: bad output\n", idx+1)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d cases passed\n", len(cases))
}
