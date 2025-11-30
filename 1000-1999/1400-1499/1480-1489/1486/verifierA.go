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
const testcasesRaw = `7 6 0 4 8 7 6 4
8 5 9 3 8 2 4 2 1
10 4 8 9 2 4 1 1 10 5 7
9 1 5 6 5 9 10 3 8 7
8 8 4 0 8 0 1 6 10
1 9
8 5 3 5 1 3 9 3 3
3 8 7 1
2 5 8
8 1 4 8 4 1 8 5 8
4 9 8 9 4
8 1 9 6 5 9 3 4 2
4 2 0 9 10
5 7 1 1 10 2
3 0 1 8
7 8 4 8 3 3 10 9
7 9 4 7 7 10 10 5
2 5 9
2 7 9
6 3 3 0 4 1 3
6 2 5 6 0 1 2
4 0 9 10 8
10 10 1 0 1 10 3 9 9 1 6
2 5 1
1 9
1 3
3 1 7 3
1 10
1 8
7 9 1 4 1 3 1 10
5 5 6 2 0 8
8 0 9 1 6 3 4 5 7
10 2 10 3 0 10 2 2 5 8 4
2 9 7
3 0 7 10
7 9 8 4 10 5 6 10
5 2 8 0 7 1
6 0 8 4 2 3 7
6 9 4 10 5 9 10
10 2 4 6 6 10 1 0 9 3 5
3 3 3 10
8 6 10 9 6 0 6 9 6
1 2
8 1 4 2 7 8 7 8 9
1 0
8 5 4 7 0 6 3 8 10
2 2 0
7 10 6 5 0 3 0 0
9 9 1 3 1 9 10 3 4 4
3 1 7 6
2 0 4
8 1 4 2 10 8 10 10 5
2 2 4
1 0
1 3
5 8 5 5 9 0
10 10 7 10 7 10 6 5 8 2 3
7 9 4 0 2 2 4 5
6 5 1 5 9 0 0
5 2 2 9 4 5
7 8 2 4 1 7 3 0
5 2 8 1 4 6
6 4 6 1 1 8 7
8 5 5 1 7 1 7 6 0
5 5 10 2 2 10
10 6 10 1 1 1 3 3 0 6 0
2 6 8
9 4 7 7 9 10 3 6 1 5
4 4 9 2 6
4 5 1 1 0
9 7 10 3 1 7 6 4 3 10
1 3
10 2 1 3 7 6 5 8 2 1 9
8 2 9 6 10 10 6 8 7
6 7 7 10 10 3 8
10 3 0 5 5 5 0 8 2 4 9
3 6 9 4
8 1 1 8 0 1 3 2 0
5 0 7 5 2 2
8 5 8 6 8 8 0 9 1
9 9 1 6 3 4 8 9 6 7
7 9 9 3 0 10 0 2
5 8 9 4 5 1
8 4 4 6 6 6 0 2 10
3 3 4 5
1 0
8 6 2 7 9 1 10 2 5
7 0 9 7 6 7 0 1
8 2 0 0 9 9 2 10 5
2 8 10
6 3 6 7 1 0 9
8 9 10 5 10 1 10 9 4
3 6 4 10
2 8 3
1 6
8 5 3 7 5 10 1 0 0
8 4 0 8 10 9 9 3 3
2 10 8
9 6 8 4 1 2 6 9 6 1
2 6 1
`

type testCase struct {
	arr []int64
}

func solveCase(h []int64) string {
	var carry int64
	for i, v := range h {
		carry += v
		need := int64(i)
		if carry < need {
			return "NO"
		}
		carry -= need
	}
	return "YES"
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var res []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid n: %v", idx+1, err)
		}
		if len(parts) != n+1 {
			return nil, fmt.Errorf("line %d: expected %d numbers, got %d", idx+1, n+1, len(parts))
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(parts[i+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse value %d: %v", idx+1, i+1, err)
			}
			arr[i] = v
		}
		res = append(res, testCase{arr: arr})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	return res, nil
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

	for i, tc := range cases {
		var input strings.Builder
		fmt.Fprintf(&input, "1\n%d\n", len(tc.arr))
		for idx, v := range tc.arr {
			if idx > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.FormatInt(v, 10))
		}
		input.WriteByte('\n')

		expected := solveCase(tc.arr)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(strings.ToUpper(got)) != expected {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
