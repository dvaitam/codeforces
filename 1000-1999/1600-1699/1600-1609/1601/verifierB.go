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
1 5 3
4 0 2 2 5 5 1 1 5
4 5 0 0 0 1 5 1 0
4 0 0 3 4 4 2 3 3
2 3 0 2 1
3 4 1 3 1 2 0
1 5 0
4 5 1 0 3 3 2 1 5
1 1 4
2 0 1 3 3
3 4 1 0 4 3 1
4 5 5 3 4 3 5 2 3
4 5 5 1 4 4 1 0 2
3 2 0 4 1 2 4
2 3 4 2 5
4 0 0 4 0 0 1 1 0
3 0 3 2 1 1 5
4 2 4 3 4 4 0 4 0
1 5 3
2 2 4 4 3
4 3 4 4 1 0 5 0 5
2 2 4 4 2
3 0 3 2 2 3 3
4 0 1 5 1 1 2 5 2
1 0 3
4 1 3 4 5 0 5 5 1
3 3 0 4 3 3 3
1 0 3
2 0 0 4 4
2 5 2 0 5
3 1 3 3 0 0 4
4 4 5 2 5 0 5 5 4
3 1 3 2 5 5 0
2 0 3 3 2
2 3 2 5 0
1 0 3
3 0 4 5 4 4 1
2 0 5 4 5
4 4 2 0 1 3 4 3 0
1 3 0
1 3 1
1 3 3
4 0 3 2 5 2 0 2 0
1 2 5
1 2 2
2 0 1 2 0
2 1 0 1 5
1 5 0
3 2 5 0 4 1 1
2 3 0 3 2
3 1 0 1 2 2 3
3 2 4 5 2 1 4
1 0 4
3 1 3 1 1 1 2
2 1 1 2 2
4 5 0 1 4 0 3 0 5
1 1 3
3 4 3 5 1 4 3
3 5 2 0 1 3 5
3 5 4 0 3 3 0
4 5 2 3 1 2 2 3 0
2 0 2 0 4
2 5 3 3 1
4 3 1 1 3 2 4 1 2
4 5 5 0 3 1 2 0 5
4 4 3 0 1 2 0 5 2
2 3 5 3 0
4 1 4 3 2 5 0 0 2
1 0 2
4 4 4 5 3 3 0 5 2
3 2 5 1 4 0 0
1 2 2
3 0 4 1 1 0 3
3 2 4 1 4 4 5
2 4 0 3 5
4 5 2 2 3 2 4 5 1
2 0 5 0 3
4 4 3 1 4 5 2 2 5
4 5 3 1 3 3 5 4 2
4 5 0 3 2 1 5 3 0
2 0 2 3 3
1 4 0
1 5 5
4 0 2 0 0 4 0 2 5
3 5 1 1 4 2 1
1 3 3
3 3 1 2 3 5 5
4 1 3 5 1 4 2 1 1
2 3 2 3 3
4 3 5 1 1 3 1 4 5
1 3 0
2 5 0 1 2
1 5 5
2 1 4 2 4
1 5 4
3 2 3 3 0 5 5
4 4 3 3 2 5 3 1 2
3 0 0 0 1 2 0
3 5 0 1 0 3 5
2 4 3 4 1
`

type testCase struct {
	n int
	a []int
	b []int
}

func solveCase(tc testCase) (int, []int, bool) {
	n := tc.n
	a := make([]int, n+2)
	b := make([]int, n+2)
	for i := 1; i <= n; i++ {
		a[i] = tc.a[i-1]
		b[i] = tc.b[i-1]
	}

	f := make([]int, n+2)
	pre := make([]int, n+2)
	for i := 0; i <= n; i++ {
		f[i] = -1
		pre[i] = 0
	}
	f[n] = 0
	pre[n] = n + 1
	l, r := n, n
	for r > 0 {
		for l > 0 && f[l] == f[r] {
			l--
		}
		for i := r; i > l; i-- {
			base := i + b[i]
			if base < 0 || base > n {
				continue
			}
			low := base - a[base]
			if low < 0 {
				low = 0
			}
			for j := low; j <= base; j++ {
				if pre[j] == 0 {
					f[j] = f[r] + 1
					pre[j] = i
				}
			}
		}
		if l < 0 || pre[l] == 0 {
			return -1, nil, false
		}
		r = l
	}
	steps := f[0]
	var path []int
	var build func(int)
	build = func(x int) {
		if pre[x] <= n {
			build(pre[x])
			path = append(path, x)
		}
	}
	build(0)
	return steps, path, true
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		expected := 1 + 2*n
		if len(fields) != expected {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, expected, len(fields))
		}
		tc := testCase{n: n, a: make([]int, n), b: make([]int, n)}
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse a: %v", idx+1, err)
			}
			tc.a[i] = v
		}
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[1+n+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse b: %v", idx+1, err)
			}
			tc.b[i] = v
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
		expectedSteps, expectedPath, _ := solveCase(tc)

		var sb strings.Builder
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for idx, v := range tc.a {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for idx, v := range tc.b {
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
		var expectedBuilder strings.Builder
		expectedBuilder.WriteString(strconv.Itoa(expectedSteps))
		expectedBuilder.WriteByte('\n')
		if expectedSteps != -1 {
			for _, v := range expectedPath {
				expectedBuilder.WriteString(strconv.Itoa(v))
				expectedBuilder.WriteByte(' ')
			}
		}
		expectedStr := strings.TrimSpace(expectedBuilder.String())
		if strings.TrimSpace(got) != expectedStr {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, expectedStr, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
