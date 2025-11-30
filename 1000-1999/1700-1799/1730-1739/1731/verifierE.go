package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `6 6
25 298
24 202
18 40
2 1
23 120
11 5
20 87
28 64
30 381
13 43
21 134
9 27
13 7
29 265
23 35
23 65
12 35
13 75
28 364
17 21
27 44
4 6
17 45
24 124
12 23
26 130
25 219
22 117
14 57
3 1
11 29
17 43
15 59
26 74
4 5
16 14
29 337
24 153
7 11
19 71
8 12
20 32
29 171
5 7
9 19
29 249
8 20
20 69
11 24
7 8
13 19`

type testCase struct {
	n int
	m int64
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func edgesCounts(n int) []int64 {
	e := make([]int64, n+1)
	for i := n; i >= 1; i-- {
		c := n / i
		e[i] = int64(c) * int64(c-1) / 2
		for j := i * 2; j <= n; j += i {
			e[i] -= e[j]
		}
	}
	return e
}

func minCost(n int, m int64) int64 {
	edges := edgesCounts(n)
	var cost int64
	for d := n; d >= 2 && m > 0; d-- {
		avail := edges[d] / int64(d-1)
		if avail == 0 {
			continue
		}
		need := m / int64(d-1)
		if need > avail {
			need = avail
		}
		cost += need * int64(d)
		m -= need * int64(d-1)
	}
	if m > 0 {
		return -1
	}
	return cost
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var cases []testCase
	for idx, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d malformed", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d parse n: %v", idx+1, err)
		}
		m, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d parse m: %v", idx+1, err)
		}
		cases = append(cases, testCase{n: n, m: m})
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		want := minCost(tc.n, tc.m)
		var input strings.Builder
		fmt.Fprintf(&input, "1\n%d %d\n", tc.n, tc.m)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotVal, err := strconv.ParseInt(strings.TrimSpace(got), 10, 64)
		if err != nil || gotVal != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
