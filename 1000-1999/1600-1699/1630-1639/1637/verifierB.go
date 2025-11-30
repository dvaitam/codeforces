package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `
8
6 12 0 4 14 10 4 11
4
10 19 5 12
3
7 15 4
5
17 6 18 13 6
4
19 3 15 4
2
15 10
2
3 11
8
9 13 1 16 17 8 0 14
9
13 8 11 16 7 9 12 17 4
5
13 11 1 15 15
8
9 16 5 17 12 14 7 14
6
14 6 17 7 5 1
7
14 5 15 16 6 16 15
3
14 7 13
4
14 14 19 4
8
19 2 8 7 2 1 0 12
8
4 5 11 8 16 12 12 4
2
13 7
1
9
4
6 13 0 13
7
12 14 12 2 17 10 15
2
8 14
7
10 19 0 6 18 7 7
8
12 6 15 1 0 1 11 19
9
19 17 9 12 4 4 3 19 8
6
16 13 5 17 0 13
3
11 3 13
5
9 1 4 0 3
9
0 10 17 13 1 7 9 3 19
1
6
5
8 19 6 3 4
6
15 7 17 11 12 13
8
8 5 2 9 12 14 1 9
8
5 8 4 15 11 13 5 9
7
13 12 13 14 5 4 18
9
0 2 1 2 5 8 6 6 14
3
2 5 18
7
13 13 2 18 2 17 15
9
9 10 8 4 12 14 12 9 5
5
13 1 2 14 17
4
12 4 8 12
9
17 4 16 9 7 4 5 8 10
7
0 6 19 8 4 5 16
7
14 7 14 16 4 12 13
1
7
8
19 7 0 17 15 5 16 7
5
16 18 13 4 4
7
13 3 1 8 6 7 4
5
18 14 2 14 0
9
19 19 3 18 14 7 12 8 0
5
8 7 8 17 18
5
10 8 15 17 3
9
0 9 0 15 19 19 15 2 11
7
1 5 17 12 14 14 8
9
11 8 19 8 5 7 2 4 16
1
12
1
15
4
18 4 16 6
1
5
7
17 17 11 4 9 12 19
4
10 9 18 10
2
16 12
8
13 3 4 2 3 8 3 15
1
5
9
15 19 5 3 0 1 3 14 1
2
1 4
5
7 9 18 10 15
4
19 16 15 2
2
1 17
2
18 5
2
15 6
8
14 13 3 18 2 0 17 10
1
8
4
13 6 18 0
1
10
7
14 0 15 0 14 9 3
3
5 13 4
2
0 16
9
16 15 7 13 3 13 12 2 1
1
17
5
9 10 10 3 12
8
9 7 15 10 14 3 17 6
5
0 3 16 0 7
2
17 4
9
19 17 18 19 12 2 5 2 17
7
15 16 17 3 0 13 14
6
2 10 13 15 0 10
4
4 16 11 11
9
3 11 2 12 16 9 3 2 17
4
9 7 7 11
8
3 1 9 16 13 3 0 17
4
3 3 9 0
5
16 16 2 19 1
6
2 13 7 4 8 15
8
10 5 15 1 1 17 2 14
2
17 2
2
16 7
2
0 17
3
6 12 6
4
14 19 2 2
`

type testCase struct {
	n   int
	arr []int
}

func solveCase(tc testCase) int64 {
	n := tc.n
	totalLen := int64(n*(n+1)*(n+2)) / 6
	addZero := int64(0)
	for i, v := range tc.arr {
		if v == 0 {
			left := i + 1
			right := n - i
			addZero += int64(left * right)
		}
	}
	return totalLen + addZero
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines)/2)
	for i := 0; i < len(lines); {
		nLine := strings.TrimSpace(lines[i])
		if nLine == "" {
			i++
			continue
		}
		n, err := strconv.Atoi(nLine)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", i+1, err)
		}
		i++
		if i >= len(lines) {
			return nil, fmt.Errorf("line %d: missing array line", i)
		}
		arrFields := strings.Fields(lines[i])
		if len(arrFields) != n {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", i+1, n, len(arrFields))
		}
		arr := make([]int, n)
		for j, f := range arrFields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse arr[%d]: %v", i+1, j, err)
			}
			arr[j] = v
		}
		cases = append(cases, testCase{n: n, arr: arr})
		i++
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
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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

	for i, tc := range cases {
		expected := solveCase(tc)

		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for idx, v := range tc.arr {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strconv.FormatInt(expected, 10) {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
