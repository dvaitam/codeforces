package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesRaw = `8 25 49 27 3 17 33 32 26
8 20 31 23 38 14 33 9 19
3 49 7 40
8 17 35 46 39 10 20 7 47
2 44 22
5 36 7 23 28 21
6 41 14 36 31 29 34
4 4 36 1 6
7 26 46 43 41 1 40 32
8 22 16 47 21 46 5 13 37
3 16 10 35
5 6 6 21 33 32
2 20 36
4 46 8 36 22
8 35 14 39 36 38 19 29 6
6 25 21 37 16 19 12
3 12 3 40
7 17 31 5 6 44 49 9
3 3 6 45
8 35 44 26 46 34 18 34 16
8 14 44 38 27 38 18 29 32
7 42 45 23 6 21 40 8
5 38 41 22 13 16
2 47 18
2 46 15
4 11 22 28 4
2 10 45
3 3 37 41
6 39 44 5 2 8 41
3 39 37 8
5 6 24 8 3 39
2 13 12
7 8 31 14 47 4 44 2
6 28 40 7 17 5 15
2 42 20
4 28 12 4 33
5 3 39 7 45 26
3 17 23 47
5 37 11 45 44 14
8 4 44 11 11 22 34 17 8
6 29 43 12 1 31 44
5 37 33 20 42 23
5 43 17 10 36 45
2 30 48
2 22 48
2 35 18
3 16 49 31
4 40 19 44 23
6 41 40 9 46 20 25
7 27 42 6 1 39 13 45
4 11 16 15 41
5 25 46 44 37 27
2 26 45
6 27 50 43 46 3 11
5 5 17 45 11 29
6 32 36 39 49 1 3
5 21 20 30 4 27
3 36 41 6
8 47 9 1 26 44 27 21 1
3 1 46 49
2 44 34
6 7 13 8 39 42 13
8 20 18 45 12 7 31 26 41
2 2 18
5 8 17 9 42 34
8 42 42 23 8 10 18 2 3
2 14 44
4 36 21 24 37
8 3 48 45 39 42 32 46 42
5 41 28 24 35 12
3 25 38 19
2 9 10
4 22 22 24 46
2 22 50
6 3 3 18 11 10 38
4 24 26 36 9
4 8 31 47 16
2 20 12
8 34 47 5 20 26 22 20 27
2 7 36
5 31 22 22 8 31
2 45 32
5 3 20 22 48 44
3 11 41 37
5 41 6 5 6 13
7 15 4 25 1 7 26 36
6 19 29 32 38 46 44
3 28 6 24
3 17 38 50
3 28 13 23
2 5 45
2 34 29
8 44 13 8 32 26 17 14 42
2 14 40
3 7 13 30
5 24 35 10 7 39
5 10 37 26 41 44
5 34 32 44 21 32
5 41 43 13 35 40
3 1 22 46`

type testCase struct {
	n   int
	arr []int
}

func solveCase(tc testCase) string {
	for i := 1; i < tc.n; i++ {
		if tc.arr[i] >= tc.arr[i-1] {
			return "YES"
		}
	}
	return "NO"
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
		fmt.Fprintf(&input, "%d", tc.n)
		for _, v := range tc.arr {
			fmt.Fprintf(&input, " %d", v)
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
