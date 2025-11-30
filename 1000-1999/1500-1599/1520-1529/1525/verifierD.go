package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solve(arr []int) int {
	ones := make([]int, 0, len(arr))
	zeros := make([]int, 0, len(arr))
	for i, v := range arr {
		if v == 1 {
			ones = append(ones, i)
		} else {
			zeros = append(zeros, i)
		}
	}
	m := len(ones)
	if m == 0 {
		return 0
	}
	const inf = int(1e9)
	dp := make([]int, m+1)
	for i := 1; i <= m; i++ {
		dp[i] = inf
	}
	for _, pos0 := range zeros {
		for i := m; i >= 1; i-- {
			d := abs(ones[i-1] - pos0)
			if dp[i-1]+d < dp[i] {
				dp[i] = dp[i-1] + d
			}
		}
	}
	return dp[m]
}

func runProg(prog, input string) (string, error) {
	cmd := exec.Command(prog)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func lineToInput(line string) (string, error) {
	parts := strings.Fields(line)
	if len(parts) < 1 {
		return "", fmt.Errorf("empty line")
	}
	n, _ := strconv.Atoi(parts[0])
	if len(parts) != 1+n {
		return "", fmt.Errorf("expected %d numbers got %d", n+1, len(parts))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(parts[1+i])
	}
	sb.WriteByte('\n')
	return sb.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load embedded testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range tests {
		input, err := lineToInput(tc.raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		want := solve(tc.arr)
		gotStr, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(gotStr))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid integer output %q\n", idx+1, gotStr)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%d\ngot:\n%d\n", idx+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

type testCase struct {
	raw string
	arr []int
}

// Embedded copy of testcasesD.txt (each line is \"n a1 a2 ... an\").
const testcaseData = `
9 1 0 0 1 1 0 0 1 0
17 0 1 0 0 0 1 0 1 1 0 0 0 0 0 0 0 0
4 0 0 0 0
10 1 1 1 0 0 0 0 0 0 0
5 0 0 0 0 0
15 1 0 0 0 1 0 1 1 0 0 0 0 0 0 0
12 0 0 0 0 0 1 1 0 1 1 0 0
4 0 0 0 1
2 1 0
15 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1
5 0 1 0 1 0
11 0 0 0 1 1 0 0 0 1 1 0
18 0 1 1 1 0 1 0 0 0 1 0 0 0 0 0 0 1 0
20 0 1 0 1 0 0 0 0 0 0 0 0 0 0 0 0 1 1 0 1
11 1 1 0 0 0 0 1 0 1 0 0
5 0 0 0 0 0
10 0 0 0 0 1 0 0 0 0 0
12 0 1 1 0 0 0 0 0 1 1 0 1
5 0 0 0 0 0
10 1 0 0 0 0 0 0 0 0 1
15 1 0 0 0 0 0 0 0 1 1 0 1 0 0 0
3 1 0 0
19 1 0 0 1 0 0 0 0 0 0 0 0 1 0 0 1 0 1 0
9 0 0 0 0 0 0 0 0 0
19 0 0 0 0 0 1 0 1 0 1 1 1 0 0 0 1 0 0 1
9 1 0 0 0 0 0 1 0 0
3 0 0 0
17 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
6 0 0 0 0 0 0
3 1 0 0
3 0 0 1
17 0 1 0 0 0 0 0 0 1 1 0 0 1 1 0 0 0
19 0 1 0 0 0 0 0 1 0 0 0 0 0 1 0 0 0 0 1
20 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0
19 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
4 0 1 0 0
15 0 0 1 0 1 1 1 0 0 1 0 0 1 0 1
8 1 0 0 0 0 0 0 0
20 0 1 0 0 1 0 0 0 1 0 1 1 0 0 1 0 1 0 1 1
8 1 0 1 0 0 0 0 1
16 0 1 0 0 0 1 0 0 0 0 1 0 0 0 0 0
10 0 0 0 0 0 0 0 0 0 1
11 0 0 0 0 0 0 0 0 0 0 0
4 0 0 0 0
19 0 0 1 0 0 0 0 0 0 0 0 1 0 0 1 0 1 0 1
20 1 1 0 0 0 1 1 0 0 1 0 0 1 0 1 0 0 0 1 0
12 1 1 0 0 1 1 0 1 0 0 0 0
4 0 0 0 0
13 0 0 0 0 0 0 0 0 0 0 0 0 0
12 0 0 0 1 0 1 0 0 0 1 0 0
5 0 0 0 1 1
20 0 0 1 0 1 0 0 1 0 1 0 0 1 0 0 0 0 0 1 0
12 1 0 1 1 1 0 0 0 0 0 0 0
2 0 0
17 0 0 0 1 0 0 0 0 0 1 0 0 0 0 0 0 0
4 0 1 0 1
9 0 1 0 0 1 1 0 0 1
6 0 0 0 0 0 0
8 0 0 0 0 0 0 0 0
3 0 0 1
7 1 1 0 0 1 0 0
11 0 0 0 0 0 0 0 0 0 0 0
16 0 1 0 1 0 0 0 1 1 1 0 1 1 1 0 0
2 0 0
18 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1
9 1 0 0 0 1 0 0 0 1
17 0 0 0 0 1 0 0 0 0 1 0 0 0 1 0 0 0
2 1 0
12 0 0 0 0 0 0 0 0 0 0 0 0
17 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
15 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0
16 0 0 0 0 1 1 0 0 0 1 0 0 0 0 0 0
16 0 0 0 0 0 0 0 0 0 0 1 0 0 0 0 0
8 0 0 0 1 0 0 1 1
8 0 0 0 0 0 0 0 1
11 0 1 0 0 0 0 0 0 0 0 0
4 0 0 0 0
5 0 0 0 0 0
14 0 1 0 1 0 0 0 0 0 0 1 0 0 0
14 0 0 0 0 0 0 1 0 0 0 0 1 0 1
15 0 0 0 0 0 0 0 0 0 1 0 0 0 0 0
7 1 0 0 0 0 0 0
13 0 0 0 0 0 0 0 1 0 0 0 0 0
11 0 1 0 1 0 0 1 0 0 1 0
3 0 1 0
11 0 0 0 0 0 1 0 0 0 0 0
14 0 0 1 1 0 0 0 1 0 0 0 1 0 1
13 0 0 1 0 0 1 1 0 1 0 1 0 1
7 0 0 0 0 0 0 0
17 0 0 0 0 0 0 0 0 1 0 0 0 0 0 0 1 0
12 0 0 0 0 0 0 0 0 0 0 0 0
10 0 0 0 0 0 0 0 0 0 0
13 0 0 0 1 1 0 0 0 1 0 0 1 0
16 0 0 0 0 1 0 1 1 0 0 0 1 0 0 0 0
2 0 1
17 0 0 0 0 0 0 0 0 0 1 0 0 0 0 1 0 0
7 0 0 1 1 1 0 0
3 1 0 0
19 1 0 0 0 1 0 0 0 0 0 0 1 0 0 0 0 0 0 1
15 0 0 0 0 0 1 1 1 1 0 0 0 0 0 0
`

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tc, err := parseTestcaseLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", i+1, err)
		}
		tests = append(tests, tc)
	}
	return tests, nil
}

func parseTestcaseLine(line string) (testCase, error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return testCase{}, fmt.Errorf("empty line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return testCase{}, fmt.Errorf("bad n: %v", err)
	}
	if len(fields) != n+1 {
		return testCase{}, fmt.Errorf("expected %d values got %d", n+1, len(fields))
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		val, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return testCase{}, fmt.Errorf("bad value at position %d: %v", i+1, err)
		}
		arr[i] = val
	}
	return testCase{raw: line, arr: arr}, nil
}
