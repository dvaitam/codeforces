package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `
2 2 1
6 4 4 2 1 1 1
4 3 1 2 3
3 1 1 2
2 1 2
3 1 1 2
3 3 3 2
1 1
6 4 5 2 2 2 4
3 1 3 2
1 1
5 3 5 2 4 4
5 3 4 4 2 2
3 2 1 1
1 1
6 3 5 5 6 4 6
3 1 3 1
1 1
2 2 2
2 2 2
6 5 3 6 5 2 3
1 1
6 2 3 5 5 2 1
3 1 2 2
1 1
3 3 1 2
6 6 3 1 3 3 3
2 2 1
3 3 1 2
3 1 2 2
5 2 3 5 1 3
1 1
2 2 2
3 3 1 2
2 2 1
1 1
1 1
6 2 5 6 2 5 1
5 4 5 2 3 1
1 1
4 2 4 2 2
4 4 4 1 2
4 4 2 4 2
4 2 1 1 3
3 1 3 1
2 2 2
2 2 1
3 3 1 3
4 1 4 4 1
4 2 2 3 3
6 4 6 3 4 5 2
6 6 3 3 4 4 1
3 3 3 1
1 1
5 2 3 1 2 4
5 4 4 4 2 1
2 1 1
5 3 1 4 4 2
5 1 2 2 5 3
5 4 5 4 1 1
1 1
4 3 2 1 3
1 1
3 2 1 1
5 4 4 2 3 4
2 2 2
1 1
1 1
5 4 4 4 1 2
6 3 4 4 1 5 2
3 1 1 3
2 2 2
3 2 3 1
6 2 1 6 3 1 5
1 1
6 5 4 5 1 5 2
4 3 3 2 2
6 1 6 1 5 3 5
4 3 4 4 4
2 2 2
6 3 2 4 1 5 2
6 5 2 3 3 2 5
3 3 3 1
1 1
6 2 3 6 3 3 2
3 1 3 1
5 1 3 1 2 2
5 4 5 1 1 2
6 4 6 6 2 5 6
3 3 2 3
2 2 1
3 2 1 3
3 3 3 3
4 4 4 3 3
4 4 1 3 2
5 4 5 5 3 4
4 1 2 1 4
6 3 5 3 1 3 4
2 2 1
2 1 1
`

type testCase struct {
	n   int
	arr []int
}

func parseTests(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	tests := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil || len(fields) != n+1 {
			return nil, fmt.Errorf("bad testcase line: %q", line)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("bad value in line: %q", line)
			}
			arr[i] = v
		}
		tests = append(tests, testCase{n: n, arr: arr})
	}
	return tests, nil
}

// solve replicates the placeholder logic from 1969E.go: it always outputs 0 for a test case.
func solve(tc testCase) int {
	_ = tc
	return 0
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\\n")
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		input := buildInput(tc)
		want := strconv.Itoa(solve(tc))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed: expected %s got %s\ninput:\n%s", idx+1, want, strings.TrimSpace(got), input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
