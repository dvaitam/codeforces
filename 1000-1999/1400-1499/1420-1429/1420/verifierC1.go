package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC1.txt so the verifier is self-contained.
const testcasesRaw = `1 1
5 2 5 3 4 1
20 9 11 4 5 1 16 3 19 12 17 20 14 2 10 6 7 15 8 18 13
3 3 2 1
7 1 4 6 5 3 2 7
9 3 4 5 1 2 9 7 6 8
1 1
15 4 3 5 2 6 11 14 9 10 13 7 1 15 12 8
14 5 3 13 14 6 11 1 12 9 7 2 10 8 4
9 2 8 7 3 5 1 4 9 6
10 2 4 1 7 3 6 5 10 9 8
4 3 2 4 1
8 2 7 5 8 1 6 3 4
13 13 8 1 3 12 7 9 4 6 10 2 5 11
3 1 2 3
11 11 7 3 1 10 2 9 5 8 6 4
12 6 9 4 12 5 8 2 7 3 1 11 10
16 10 14 2 9 11 5 1 8 13 6 16 15 3 12 4 7
10 2 10 5 8 7 9 6 3 4 1
9 4 1 8 5 9 2 6 7 3
7 4 7 5 6 2 3 1
9 2 4 5 9 7 1 8 3 6
12 5 11 6 3 2 7 9 4 10 12 8 1
12 12 6 1 5 10 2 9 8 4 3 7 11
14 4 3 2 6 11 13 12 8 10 9 5 14 1 7
6 4 6 3 2 5 1
3 3 1 2
6 4 5 3 6 1 2
3 3 1 2
20 6 19 10 4 17 20 3 18 8 16 1 7 14 12 13 11 9 2 5 15
6 1 5 2 3 4 6
13 2 7 4 3 6 10 13 8 11 12 9 5 1
15 5 13 4 14 10 11 3 9 1 8 6 15 2 12 7
7 1 4 2 3 6 5 7
7 7 4 6 2 3 5 1
6 6 4 5 1 3 2
19 17 4 13 9 14 2 1 3 15 16 7 8 6 19 11 10 18 12 5
11 2 6 7 3 10 9 11 4 5 8 1
3 2 1 3
19 9 11 12 8 19 2 3 6 10 18 7 1 16 15 5 14 13 17 4
19 14 5 18 19 8 16 17 3 6 1 4 10 11 7 9 2 15 12 13
19 5 7 3 17 1 11 16 12 15 18 13 4 19 10 2 14 6 9 8
9 4 2 1 6 9 3 8 7 5
1 1
16 8 6 9 2 3 15 14 16 10 11 12 7 4 5 13 1
11 4 6 2 9 5 10 11 3 8 7 1
3 2 3 1
18 13 10 4 8 9 17 12 18 5 3 7 6 1 14 15 11 2 16
17 14 13 16 17 4 5 7 10 9 1 11 12 8 15 3 6 2
3 2 1 3
9 1 2 9 8 7 6 3 5 4
5 5 3 1 2 4
17 7 1 9 6 3 15 5 12 13 10 8 16 4 17 14 11 2
18 16 15 9 10 5 4 18 8 6 7 12 2 13 14 1 11 17 3
5 2 5 4 1 3
1 1
11 2 5 8 11 3 9 6 7 1 4 10
12 10 12 4 9 3 11 5 6 8 1 7 2
14 9 10 12 11 6 13 14 2 7 4 5 8 1 3
6 6 3 2 1 5 4
18 4 8 6 11 17 10 12 2 13 15 5 7 9 1 14 3 18 16
2 2 1
17 4 5 15 8 13 16 14 11 17 12 9 3 6 7 1 2 10
13 13 9 11 10 7 8 2 3 6 5 12 4 1
12 12 8 6 2 4 9 5 11 1 7 3 10
1 1
11 3 4 6 2 5 8 7 11 10 1 9
1 1
14 8 9 1 13 14 2 5 12 3 4 6 11 10 7
14 3 5 1 12 13 9 8 4 14 7 2 6 11 10
13 3 12 8 4 7 1 13 6 9 11 5 2 10
2 1 2
13 13 12 8 5 6 4 10 7 3 1 9 2 11
20 20 1 12 4 14 8 13 17 18 2 9 5 10 3 7 16 11 19 6 15
7 1 4 6 7 5 2 3
2 1 2
18 13 11 16 1 8 9 15 17 10 3 6 4 12 7 5 18 2 14
15 4 8 1 5 6 9 10 3 7 15 14 13 12 11 2
17 9 1 4 11 8 3 14 13 7 2 15 17 12 10 6 16 5
13 1 10 2 8 6 11 3 7 4 5 9 13 12
4 4 3 1 2
19 15 6 9 16 11 13 8 1 2 19 3 10 14 12 5 18 7 4 17
15 9 14 15 2 1 7 13 3 4 8 11 5 12 10 6
17 9 6 15 10 13 4 2 17 3 1 5 16 11 12 8 7 14
19 16 19 6 18 12 3 4 7 14 9 5 2 10 17 1 15 8 11 13
6 6 1 4 2 5 3
3 1 3 2
19 19 15 7 12 16 1 10 11 13 6 17 5 2 4 8 9 3 18 14
13 13 4 7 1 6 3 5 11 2 8 10 12 9
8 5 7 1 2 6 3 8 4
12 11 1 7 3 12 4 8 5 10 2 6 9
8 1 8 7 3 4 6 2 5
20 10 8 2 1 9 15 3 20 7 18 14 6 17 11 19 16 12 13 4 5
4 3 2 1 4
8 6 5 7 1 2 4 8 3
13 6 1 11 2 12 3 5 7 4 10 13 9 8
3 1 3 2
13 5 13 8 6 2 3 9 12 11 10 1 7 4
6 2 5 4 1 3 6
1 1`

type testCase struct {
	n   int
	arr []int
}

func solveCase(tc testCase) string {
	prev := 0
	var ans int64
	for _, v := range tc.arr {
		diff := v - prev
		if diff > 0 {
			ans += int64(diff)
		}
		prev = v
	}
	return strconv.FormatInt(ans, 10)
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid n: %v", idx+1, err)
		}
		if len(fields) != 1+n {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, 1+n, len(fields))
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse a[%d]: %v", idx+1, i, err)
			}
			arr[i] = v
		}
		cases = append(cases, testCase{n: n, arr: arr})
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/binary")
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

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&input, "%d 0\n", tc.n)
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
	}

	got, err := runCandidate(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var expected strings.Builder
	for i, tc := range cases {
		if i > 0 {
			expected.WriteByte('\n')
		}
		expected.WriteString(solveCase(tc))
	}

	if strings.TrimSpace(got) != strings.TrimSpace(expected.String()) {
		fmt.Printf("output mismatch\nexpected:\n%s\n\ngot:\n%s\n", expected.String(), got)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
