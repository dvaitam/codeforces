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
3 2 10 5
5 9 5 -9 -5 -2
4 -6 9 9 2
1 2
3 -3 -9 5
3 1 -10 0
3 -1 -2 5
6 -7 -3 -6 -1 4 0
3 3 9 10
1 -4
4 -4 3 5 6
3 -9 6 -5
1 -1
6 6 2 -6 6 8 -10
2 -4 -4
1 -3
1 4
1 1
6 -4 -2 1 4 6 2
6 -7 -10 -3 1 5 9
4 -5 5 8 7
3 1 -5 -2
6 -8 -1 -10 2 -9 -5
5 -4 -3 10 -3 -4
3 10 3 6
1 -10
4 -6 10 -5 9
1 -3
3 9 -1 10
3 3 2 1
4 -2 -4 4 -1
6 6 9 2 8 -7 -10
5 1 7 10 9 9
3 -1 -7 5
1 0
3 10 0 -2
3 10 10 -1
2 -6 6
2 -9 9
4 0 -6 -10 10
4 -1 -2 3 2
4 -9 8 8 -4
3 -3 7 10
4 1 6 -1 -5
6 7 -5 -1 -7 5 -6
6 -2 7 -5 0 10 -8
2 1 -3
6 -1 3 0 1 -2 8
3 4 -7 5
1 10
5 8 9 -8 5 -4
5 -7 2 6 -1 3
1 -6
2 -4 0
4 8 4 -6 0
6 -5 -8 5 0 10 -5
3 10 -9 -10
4 -2 -4 -5 8
2 5 -8
2 9 3
6 3 2 3 5 2 -10
1 7
2 1 -10
3 6 -4 -10
6 -10 10 -3 -3 1 -1
2 -7 2
5 8 -1 -5 -8 -9
3 -1 4 6
5 6 0 3 -6 0
4 1 -4 -5 2
1 -3
2 -6 -4
6 -10 8 6 -5 -7 1
6 10 -9 2 10 -2 7
5 -9 9 -9 -7 10
1 -9
6 -7 3 4 2 -7 7
3 5 -6 -4
6 10 -10 -1 3 -7 6
3 9 9 -6
4 -7 6 9 -7
3 -7 -7 5
2 9 -4
3 6 -4 1
6 3 -1 -5 -9 7 5
2 5 0
2 -10 -10
6 -8 -7 8 5 -6 -8
5 -8 -7 -2 10 -3
4 -1 -2 4 -9
1 -5
1 -1
3 0 3 -7
1 -9
1 -6
6 10 -5 0 1 4 9
3 -8 1 0
2 -7 2
4 4 -2 2 5
4 10 -5 -7 -6

`

type testCase struct {
	arr []int
}

func canForm(a []int) bool {
	n := len(a)
	var dfs func(int, int, bool) bool
	dfs = func(i, sum int, used bool) bool {
		if i == n {
			return used && sum == 0
		}
		if dfs(i+1, sum, used) {
			return true
		}
		if dfs(i+1, sum+a[i], true) {
			return true
		}
		if dfs(i+1, sum-a[i], true) {
			return true
		}
		return false
	}
	return dfs(0, 0, false)
}

func solveCase(tc testCase) string {
	if canForm(tc.arr) {
		return "YES"
	}
	return "NO"
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		if len(fields) != 1+n {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, 1+n, len(fields))
		}
		tc := testCase{arr: make([]int, n)}
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse value: %v", idx+1, err)
			}
			tc.arr[i] = v
		}
		cases = append(cases, tc)
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

	for i, tc := range cases {
		expected := solveCase(tc)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(len(tc.arr)))
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
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
