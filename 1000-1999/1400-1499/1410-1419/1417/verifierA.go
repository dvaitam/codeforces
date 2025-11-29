package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesRaw = `7 49 27 3 17 33 32 26 20
8 23 19 7 17 5 10 5 4 20
5 35 10 20 7 5 22
8 36 7 23 28 21 14 36 31 29
9 17 2 1 3 13 1 16 11 8 11
2 13 10 4
4 10 9 8 2 2
6 33 32 7 20 19 8 22
9 14 13 10 9 10 5 8 2 10 13
7 21 19 8 10 6 7 6 2
10 43 17 31 5 6 9 10 3 6 35 26
9 18 17 8 7 14 9 15 16 12 3
6 40 8 32 38 22 13 16
1 47 18
2 46 15 24
3 22 14 2 4
3 45 15 3 37
9 39 5 2 8 13 39 37 8 26 6
6 8 1 1 4 3 2 8
4 47 4 44 2 35
7 40 7 17 5 15 5 20 23
7 12 1 9 8 1 10 2 12
7 13 5 6 12 8 10 3 12
4 50 4 44 11 11
6 34 17 8 29 12 1 31
7 37 33 20 23 25 17 10 36
1 30 24
2 22 2 18
5 9 4 8 6 5 6
10 41 40 9 20 25 27 6 1 39 13 22
3 16 8 15 13
10 27 2 13 23 19 14 25 22 23 2 6
8 5 3 2 4 5 4 5 5 1
1 32 21
5 30 2 26 27 29 26
7 13 9 11 2 12 3 1 7
7 21 1 7 1 1 17 20 4
4 8 4 5 5 3
2 31 28 30
7 41 6 2 18 29 8 17 9
9 42 42 23 8 10 18 2 3 3 14
5 36 21 24 3 32 30
7 24 18 6 7 13 19 10 1
3 10 5 6 6
6 46 6 22 40 3 3 18
3 10 10 5 6
7 36 9 19 8 31 16 4 20
3 34 5 20 26
6 20 14 4 4 18 16 16
6 22 4 16 4 16 14 2
5 22 22 5 6 21 19
7 41 6 5 6 13 15 4 25
1 7 4
9 34 19 29 32 14 28 6 24 15 17
10 50 11 28 13 23 8 5 45 2 34 29
4 8 8 7 5 4
1 14 10
3 7 2 4 4
6 35 10 7 32 10 26 28
9 32 21 32 32 13 15 1 22 21 21
1 34 10
5 39 10 25 38 19 31
2 6 5 1
2 15 3 1
5 1 1 1 1 1 1
6 33 25 33 3 6 5 28
4 19 18 14 16 13
10 38 15 2 1 12 20 33 37 17 22 5
8 17 10 14 13 13 2 6 5 8
5 47 22 4 3 31 27
3 32 6 10 23
7 3 3 2 2 2 1 1 2
3 2 1 1 2
2 45 36 42
6 13 7 13 13 13 8 2
1 40 30
10 41 22 8 40 19 9 25 19 8 34 13
1 26 15
6 49 13 30 23 41 5 3
1 32 17
1 34 14
4 6 6 5 6 5
7 33 20 8 10 28 28 6 7
7 5 1 4 2 1 4 4 4
1 32 21
5 6 3 1 1 3 6
1 23 12
3 1 1 1 1
10 10 4 1 4 2 1 5 6 1 10 4
3 12 8 2 8
6 46 17 9 2 14 24 22
8 19 10 18 11 6 19 3 4 18
10 20 6 13 5 5 8 11 17 8 8 6
5 24 14 22 2 5 20
1 26 3
2 9 7 5
9 27 24 5 19 14 10 21 12 3 8
8 41 24 41 34 4 25 27 1 27
6 29 7 12 10 16 3 6
2 18 4 18
10 45 10 45 29 26 12 27 28 12 16 30
6 34 10 23 30 6 31 14
5 1 1 1 1 1 1
2 50 41 20
9 39 10 28 31 6 32 15 35 26 18
1 8 5
1 1 1
7 34 26 29 7 17 23 19 13
10 6 1 1 3 3 5 3 1 5 2 2
2 27 10 10
9 9 9 4 9 2 7 9 7 5 5
8 24 19 21 5 6 4 23 4 13
7 38 30 9 36 20 23 31 27
4 31 16 23 17 11
8 42 4 29 20 10 32 4 40 14
1 23 16
7 1 1 1 1 1 1 1 1
10 1 1 1 1 1 1 1 1 1 1 1
7 11 6 7 11 11 7 3 8
3 34 21 9 14`

type testCase struct {
	n int
	k int
	a []int
}

func solveCase(tc testCase) int {
	sorted := make([]int, len(tc.a))
	copy(sorted, tc.a)
	sort.Ints(sorted)
	mn := sorted[0]
	ans := 0
	for i := 1; i < len(sorted); i++ {
		ans += (tc.k - sorted[i]) / mn
	}
	return ans
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
		if len(fields) < 2 {
			return nil, fmt.Errorf("invalid test line %d", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("parse n on line %d: %v", idx+1, err)
		}
		k, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("parse k on line %d: %v", idx+1, err)
		}
		if len(fields) != 2+n {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, 2+n, len(fields))
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[2+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse a[%d]: %v", idx+1, i, err)
			}
			arr[i] = v
		}
		cases = append(cases, testCase{n: n, k: k, a: arr})
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
		for i, v := range tc.a {
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
		expected.WriteString(strconv.Itoa(solveCase(tc)))
	}

	if strings.TrimSpace(got) != strings.TrimSpace(expected.String()) {
		fmt.Printf("output mismatch\nexpected:\n%s\n\ngot:\n%s\n", expected.String(), got)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
