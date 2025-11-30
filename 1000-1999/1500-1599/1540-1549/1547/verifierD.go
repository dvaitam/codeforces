package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesD.txt so the verifier is self-contained.
const testcasesRaw = `
100
2
4 11
5
15 2 0 15 8
5
7 6 15 15 12
2
7 4
5
12 0 2 5 1
3
0 8 15
5
12 13 12 14 4
3
3 1 4
4
6 8 13 9
4
12 11 13 7
3
0 8 5
3
3 6 8
3
3 2 15
4
2 11 2 13
2
0 9
4
13 3 1 1
4
10 8 7 1
3
0 2 3
5
1 6 13 9 8
2
1 10
3
11 4 12
4
14 12 3 8
4
7 9 13 8
5
9 10 0 13 10
1
12
5
4 1 10 14 11
3
8 15 0
5
1 0 11 8 14
3
10 5 11
2
10 11
5
8 9 12 3 0
5
4 9 7 8 7
3
5 13 3
1
10
3
7 14 5
1
10
2
14 8
2
3 1
5
6 10 5 8 10
1
11
5
4 13 9 8 14
3
13 9 13
5
13 1 13 4 6
1
15
5
13 7 1 14 9
5
10 7 2 9 3
2
1 1
5
6 13 1 0 15
1
5
5
9 7 0 13 1
5
3 10 4 8 15
1
11
2
6 3
5
3 5 7 8 4
1
15
5
12 1 8 7 8
5
13 1 15 10 0
1
4
1
3
1
2
4
1 2 15 10
2
10 2
3
12 12 9
3
8 6 10
4
3 4 0 12
1
5
1
11
4
12 1 13 1
3
15 10 13
4
14 0 7 6
5
8 2 13 7 13
2
0 10
3
8 3 14
1
12
1
10
5
3 0 15 4 7
4
1 2 3 12
2
0 10
1
0
1
15
3
9 2 1
5
7 3 3 1 10
5
5 2 7 5 7
4
12 8 11 12
3
13 2 12
5
7 13 5 13 15
2
12 4
2
3 15
4
14 5 4 8
2
4 10
2
9 13
5
8 6 9 0 8
4
12 6 5 11
2
10 15
2
13 15
5
6 14 0 15 2
4
1 14 7 7
1
6
3
7 6 8
2
5 1

`

const mask = (1 << 30) - 1

type testCase struct {
	arr []int
}

func solveCase(tc testCase) []int {
	n := len(tc.arr)
	y := make([]int, n)
	prev := tc.arr[0]
	for i := 1; i < n; i++ {
		y[i] = (^tc.arr[i] & mask) & prev
		prev = tc.arr[i] ^ y[i]
	}
	return y
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
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n", i+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", i+1, err)
		}
		pos++
		if pos+n > len(fields) {
			return nil, fmt.Errorf("case %d: not enough elements", i+1)
		}
		tc := testCase{arr: make([]int, n)}
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(fields[pos+j])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse arr[%d]: %v", i+1, j, err)
			}
			tc.arr[j] = v
		}
		pos += n
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		inputBuilder.WriteString(strconv.Itoa(len(tc.arr)))
		inputBuilder.WriteByte('\n')
		for i, v := range tc.arr {
			if i > 0 {
				inputBuilder.WriteByte(' ')
			}
			inputBuilder.WriteString(strconv.Itoa(v))
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
		expectedArr := solveCase(tc)
		strs := make([]string, len(expectedArr))
		for j, v := range expectedArr {
			strs[j] = strconv.Itoa(v)
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
