package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesRaw = `
100
2 2
1 2
16 5
1 1
1
13
5 3
1 2 3
18 12 9
2 1
2
7
1 1
1
7
2 2
2 1
3 20
3 3
2 1 3
8 16 9
1 1
1
10
5 3
5 2 4
14 20 10
4 4
2 1 4 3
2 3 2 15
3 3
3 2 1
5 7 3
4 2
4 2
6 12
4 3
2 4 1
2 8 9
5 5
2 1 5 4 3
15 1 2 12 3
3 3
3 2 1
11 10 11
2 2
1 2
20 7
4 3
2 4 3
20 6 11
5 1
3
2
4 2
3 2
10 19
1 1
1
14
2 1
1
2
1 1
1
20
1 1
1
11
1 1
1
14
2 2
1 2
15 14
4 1
2
14
4 2
4 1
16 7
1 1
1
9
2 1
1
14
3 1
2
2
3 3
1 2 3
16 13 3
4 2
2 4
10 16
3 2
3 1
9 11
4 4
1 2 4 3
13 20 5 9
1 1
1
19
4 4
4 1 3 2
6 1 20 9
1 1
1
8
5 1
2
6
5 3
5 4 3
15 1 3
1 1
1
18
3 3
1 3 2
3 17 1
3 2
1 3
18 15
4 2
3 2
8 16
4 1
1
4
5 3
5 4 2
15 3 7
3 2
2 1
18 6
3 1
1
5
3 2
2 3
18 1
2 1
2
4
5 1
4
20
4 1
2
14
3 2
1 3
1 2
5 3
5 4 3
10 17 15
5 5
4 5 1 2 3
11 5 14 3 20
2 1
2
12
2 2
1 2
13 6
3 3
2 3 1
10 1 20
1 1
1
4
2 1
2
19
1 1
1
16
2 2
2 1
16 8
3 2
1 2
18 17
4 4
4 2 3 1
14 19 1 9
2 2
2 1
13 20
1 1
1
15
3 3
2 1 3
16 8 16
5 1
2
8
1 1
1
2
2 2
2 1
17 8
2 1
1
2
3 1
1
9
1 1
1
5
4 2
1 3
12 1
5 2
5 4
5 7
3 2
3 1
13 6
2 2
2 1
13 11
2 2
1 2
1 9
1 1
1
3
2 1
1
8
2 1
2
18
2 2
2 1
6 17
5 2
1 4
20 12
3 3
2 1 3
9 19 19
3 2
3 2
3 8
3 1
3
2
5 3
5 4 2
3 6 2
4 2
1 4
3 11
5 2
3 1
17 17
1 1
1
20
3 3
1 3 2
16 5 7
3 1
3
12
3 2
3 2
2 6
4 4
1 3 2 4
7 14 11 6
2 1
2
14
3 2
2 3
11 20
2 2
2 1
1 13
3 3
2 1 3
3 6 20
1 1
1
7
4 1
3
1
4 2
2 1
8 17

`

const INF int64 = 1e18

type testCase struct {
	n   int
	pos []int
	val []int64
}

func solveCase(tc testCase) []int64 {
	ans := make([]int64, tc.n)
	for i := 0; i < tc.n; i++ {
		ans[i] = INF
	}
	for i := 0; i < len(tc.pos); i++ {
		idx := tc.pos[i] - 1
		if tc.val[i] < ans[idx] {
			ans[idx] = tc.val[i]
		}
	}
	for i := 1; i < tc.n; i++ {
		if ans[i-1]+1 < ans[i] {
			ans[i] = ans[i-1] + 1
		}
	}
	for i := tc.n - 2; i >= 0; i-- {
		if ans[i+1]+1 < ans[i] {
			ans[i] = ans[i+1] + 1
		}
	}
	return ans
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	pos := 1
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos+1 >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n/k", i+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", i+1, err)
		}
		k, err := strconv.Atoi(fields[pos+1])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse k: %v", i+1, err)
		}
		pos += 2
		if pos+k*2 > len(fields) {
			return nil, fmt.Errorf("case %d: insufficient data", i+1)
		}
		tc := testCase{n: n, pos: make([]int, k), val: make([]int64, k)}
		for j := 0; j < k; j++ {
			v, err := strconv.Atoi(fields[pos+j])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse pos[%d]: %v", i+1, j, err)
			}
			tc.pos[j] = v
		}
		pos += k
		for j := 0; j < k; j++ {
			v, err := strconv.Atoi(fields[pos+j])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse val[%d]: %v", i+1, j, err)
			}
			tc.val[j] = int64(v)
		}
		pos += k
		cases = append(cases, tc)
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("trailing data after parsing")
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

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimRight(out.String(), "\n"), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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

	var inputBuilder strings.Builder
	inputBuilder.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		inputBuilder.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.pos)))
		for i, v := range tc.pos {
			if i > 0 {
				inputBuilder.WriteByte(' ')
			}
			inputBuilder.WriteString(strconv.Itoa(v))
		}
		inputBuilder.WriteByte('\n')
		for i, v := range tc.val {
			if i > 0 {
				inputBuilder.WriteByte(' ')
			}
			inputBuilder.WriteString(strconv.FormatInt(v, 10))
		}
		inputBuilder.WriteByte('\n')
	}

	gotOut, err := runCandidate(bin, inputBuilder.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	lines := strings.Split(strings.TrimRight(gotOut, "\n"), "\n")
	if len(lines) != len(cases) {
		fmt.Printf("expected %d lines of output, got %d\n", len(cases), len(lines))
		os.Exit(1)
	}

	for i, tc := range cases {
		expectedAns := solveCase(tc)
		strs := make([]string, len(expectedAns))
		for j, v := range expectedAns {
			strs[j] = strconv.FormatInt(v, 10)
		}
		expected := strings.Join(strs, " ")
		got := strings.TrimSpace(lines[i])
		if got != expected {
			fmt.Printf("case %d failed\nexpected: %q\ngot: %q\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
