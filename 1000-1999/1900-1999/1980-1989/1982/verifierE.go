package main

import (
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesE = `100
10 3
5 3
10 6
11 2
20 3
13 4
13 5
17 2
19 6
8 1
14 1
9 4
6 4
5 6
20 4
18 6
17 6
16 2
9 3
10 3
17 0
7 2
15 5
3 0
7 4
4 6
9 6
1 6
7 1
20 1
12 4
9 0
1 0
20 1
10 6
8 3
16 5
12 6
12 0
19 3
4 5
1 6
6 0
20 4
9 2
7 6
13 1
11 1
6 3
5 4
19 4
14 4
8 0
2 3
20 0
5 4
3 3
10 2
20 1
2 2
5 2
9 3
16 4
11 2
3 3
6 0
17 2
18 5
13 4
15 3
6 4
14 6
5 5
4 4
4 6
6 5
16 2
4 5
14 2
9 6
4 5
12 4
18 6
3 0
2 4
1 6
5 0
14 3
16 0
11 5
5 0
13 5
19 5
8 0
7 3
1 0
20 6
13 1
18 5
10 3`

const mod = 1000000007

func solveCase(n int64, k int) int64 {
	allowed := make([]bool, n)
	for i := int64(0); i < n; i++ {
		if bits.OnesCount64(uint64(i)) <= k {
			allowed[i] = true
		}
	}
	var count int64
	start := int64(0)
	for i := int64(0); i < n; i++ {
		if !allowed[i] {
			start = i + 1
			continue
		}
		count += i - start + 1
	}
	return count % mod
}

type testCase struct {
	n int64
	k int
}

func parseTests() ([]testCase, error) {
	fields := strings.Fields(testcasesE)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no data")
	}
	pos := 0
	readInt := func() (int64, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.ParseInt(fields[pos], 10, 64)
		pos++
		return v, err
	}
	t64, err := readInt()
	if err != nil {
		return nil, err
	}
	t := int(t64)
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		n64, err := readInt()
		if err != nil {
			return nil, err
		}
		k64, err := readInt()
		if err != nil {
			return nil, err
		}
		tests[i] = testCase{n: n64, k: int(k64)}
	}
	return tests, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.FormatInt(tc.n, 10))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(tc.k))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	tests, err := parseTests()
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}

	input := buildInput(tests)
	output, err := runCandidate(os.Args[1], input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	outFields := strings.Fields(output)
	if len(outFields) != len(tests) {
		fmt.Printf("expected %d outputs, got %d\n", len(tests), len(outFields))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := strconv.FormatInt(solveCase(tc.n, tc.k), 10)
		if outFields[i] != want {
			fmt.Printf("case %d failed\nn=%d k=%d\nexpected: %s\ngot: %s\n", i+1, tc.n, tc.k, want, outFields[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
