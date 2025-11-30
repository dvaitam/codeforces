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
const testcasesRaw = `100
3 0 2 0
8 7 7 6 3 1 7 0 6
7 0 7 4 3 1 5 0
1 0
9 0 6 3 6 0 8 3 7 7
9 3 5 3 3 7 4 0 6 8
2 0 2
5 0 5 2 5 5
9 6 8 3 4 4 9 7 8 6
10 0 7 3 6 6 10 2 5 8 10
6 0 3 5 4 0 6
3 3 2 3
1 1
1 1
10 9 9 6 10 2 2 8 3 0 3
9 8 3 6 8 5 9 5 7 4
9 9 0 6 8 2 8 8 3 6
1 1
6 4 4 1 4 3 3
6 3 2 0 4 4 4
10 5 7 9 0 3 10 2 8 9 2
2 2 1
1 0
2 0 1
1 1
4 2 0 4 1
6 2 0 1 1 2 4
3 2 2 3
6 3 3 0 0 2 3
6 3 6 1 2 0 2
9 3 9 6 0 3 0 6 2 0
3 3 3 1
9 7 3 8 0 6 9 5 6 0
5 1 1 0 2 0
2 1 1
3 3 2 1
1 0
10 3 9 7 2 9 8 0 6 3 5
2 0 2
7 3 7 1 6 4 7 0
6 4 6 3 2 0 1
4 2 4 1 2
7 3 4 1 6 5 7 3
2 2 0
2 0 0
3 1 2 2
10 8 4 5 5 5 1 4 3 9 7
3 0 2 0
7 1 6 2 2 5 1 6
2 2 2
4 4 0 2 2
5 4 4 0 3 2
2 0 1
1 0
2 1 0
1 0
4 4 3 1 0
8 2 3 2 1 6 6 8 4
9 4 7 5 1 3 5 0 0 0
5 5 4 2 3 3
6 3 0 0 2 4 3
2 1 0
10 8 7 10 5 4 2 8 3 4 3
4 2 0 2 0
8 1 5 3 6 4 0 5 2
6 6 6 4 2 1 2
2 2 2
10 9 1 3 3 0 3 6 1 4 8
2 2 0
1 0
5 2 3 3 1 0
9 5 1 8 2 2 2 2 5 4
2 2 2
10 4 2 3 2 8 0 5 9 10 8
4 1 2 3 4
3 0 1 2
2 2 1
7 4 7 7 0 6 5 2
5 3 0 5 3 4
1 0
6 4 1 4 1 1 2
5 3 4 3 1 4
2 0 1
1 0
9 5 8 7 3 3 5 7 7 3
7 5 4 3 0 1 5 2
9 3 4 4 4 8 5 2 7 9
2 0 2
9 9 6 2 2 4 6 3 9 0
8 6 5 6 8 2 8 0 8
2 1 2
2 1 2
2 0 2
2 1 0
7 6 6 2 5 7 2 7
4 0 3 4 4
7 1 4 4 3 6 0 3
9 7 9 0 0 9 3 4 3 2
5 1 4 1 2 2
10 4 10 7 2 8 5 7 6 1 3`

type testCase struct {
	n int
	a []int
}

func solveCase(tc testCase) []int {
	ans := make([]int, tc.n)
	cream := 0
	for i := tc.n - 1; i >= 0; i-- {
		if tc.a[i] > cream {
			cream = tc.a[i]
		}
		if cream > 0 {
			ans[i] = 1
			cream--
		} else {
			ans[i] = 0
		}
	}
	return ans
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(strings.Fields(lines[0])[0])
	if err != nil {
		return nil, fmt.Errorf("invalid t: %v", err)
	}
	res := make([]testCase, 0, t)
	for i := 1; i < len(lines); i++ {
		fields := strings.Fields(lines[i])
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", len(res)+1, err)
		}
		if len(fields)-1 != n {
			return nil, fmt.Errorf("case %d: expected %d values got %d", len(res)+1, n, len(fields)-1)
		}
		a := make([]int, n)
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(fields[j+1])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse a%d: %v", len(res)+1, j+1, err)
			}
			a[j] = v
		}
		res = append(res, testCase{n: n, a: a})
	}
	if len(res) != t {
		return nil, fmt.Errorf("expected %d cases, got %d", t, len(res))
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
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.a[j]))
		}
		sb.WriteByte('\n')

		expected := make([]string, tc.n)
		for idx, v := range solveCase(tc) {
			expected[idx] = strconv.Itoa(v)
		}
		expectedStr := strings.Join(expected, " ")

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expectedStr {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, expectedStr, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
