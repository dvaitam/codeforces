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
const testcasesRaw = `4 4 5 2 5 1 1 2 4 1
4 0 1 3 4 5 2 0 3 3
2 3 3 1 3 1
4 3 4 0 2 2 4 2 4 0
4 1 2 0 4 2 3 5 0 4
1 1 5 3
1 3 5 3
4 0 5 3 4 3 1 1 2 3
4 1 4 2 4 0 2 2 1 5
3 3 5 4 0 1 2 1
3 2 3 5 0 2 4 3
1 5 3 2
4 5 1 2 1 1 0 4 2 1
1 3 4 3
4 0 5 5 5 5 0 3 4 5
2 1 1 1 1 1
1 4 1 3
3 4 5 2 5 2 5 0
4 2 1 2 5 4 2 4 5 1
4 5 5 4 2 1 3 2 3 3
3 3 3 0 5 1 1 0
1 5 5 5
2 4 5 4 4 0
3 5 5 4 1 3 0 0
1 0 2 5
4 2 0 5 3 2 4 5 2 5
4 1 1 4 4 0 4 1 0 5
2 4 3 4 5 1
4 5 5 3 2 0 4 1 3 1
4 4 0 3 0 5 4 3 2 1
4 0 3 4 0 2 4 3 4 5
4 3 2 1 1 4 2 1 2 4
2 3 0 0 3 3
4 4 0 2 3 1 0 2 1 2
1 1 5 0
2 4 4 1 0 0
4 4 5 4 3 5 5 0 3 4
3 2 0 5 0 1 1 3
3 2 1 0 3 2 4 2
1 0 2 4
4 1 5 2 3 1 1 2 1 5
4 5 2 2 3 4 1 0 3 2
4 1 5 5 2 5 2 3 2 2
4 5 2 0 3 2 0 3 1 2
2 5 1 2 3 1
1 5 3 3
4 4 3 4 4 1 5 3 4 2
4 1 2 4 5 0 0 3 2 3
2 3 0 4 0 3
1 2 1 5
3 1 0 1 5 2 0 1
1 0 3 4
1 1 0 4
2 0 5 3 1 0
4 0 4 1 1 1 3 3 3 0
4 1 3 0 4 2 0 4 2 5
3 2 5 3 5 1 4 4
1 2 0 5
1 2 0 5
2 4 5 1 3 1
3 1 3 4 4 0 4 0
4 0 4 3 3 1 3 4 2 1
4 2 3 0 4 4 2 1 3 1
3 3 3 0 3 5 0 4
4 3 1 3 1 2 3 3 1 1
4 2 2 5 3 4 1 4 5 0
2 2 4 5 5 0
1 1 2 0
3 4 1 3 3 5 1 3
4 1 1 3 0 2 5 4 1 5
3 1 4 4 5 4 4 3
2 5 5 4 3 5
2 4 3 0 1 3
3 4 5 1 0 0 0 1
2 3 5 4 1 5
1 4 3 5
3 5 2 3 5 1 1 2
3 5 3 3 3 5 2 3
1 1 1 0
2 3 3 5 4 1
3 5 0 0 3 1 0 3
4 4 1 5 0 5 1 4 5 3
3 3 4 4 0 0 4 5
4 0 1 3 4 2 2 0 5 0
2 4 0 4 4 2
2 3 2 0 4 3
3 1 3 3 5 4 2 1
3 0 0 1 2 5 4 4
4 0 4 3 3 2 3 5 0 0
4 4 5 1 1 0 0 5 3 1
1 5 2 2
4 0 2 2 2 4 5 1 5 4
1 2 4 4
1 1 3 0
3 2 1 2 0 4 4 3
1 2 4 0
1 4 1 3
1 4 1 2
1 4 5 1
2 5 5 1 3 1`

type testCase struct {
	n int
	d int
	s []int
	a []int
}

func solveCase(tc testCase) int {
	used := make([]bool, tc.n)
	cur := tc.d
	ans := 0
	for {
		best := -1
		bestA := 0
		for i := 0; i < tc.n; i++ {
			if used[i] {
				continue
			}
			if cur <= tc.s[i] {
				if best == -1 || tc.a[i] < bestA {
					best = i
					bestA = tc.a[i]
				}
			}
		}
		if best == -1 {
			break
		}
		used[best] = true
		if tc.a[best] > cur {
			cur = tc.a[best]
		}
		ans++
	}
	return ans
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: not enough fields", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		d, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse d: %v", idx+1, err)
		}
		expected := 2 + 2*n
		if len(fields) != expected {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, expected, len(fields))
		}
		tc := testCase{n: n, d: d, s: make([]int, n), a: make([]int, n)}
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[2+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse s: %v", idx+1, err)
			}
			tc.s[i] = v
		}
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[2+n+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse a: %v", idx+1, err)
			}
			tc.a[i] = v
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
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.d))
		for idx := 0; idx < tc.n; idx++ {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.s[idx]))
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(tc.a[idx]))
		}
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strconv.Itoa(expected) {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
