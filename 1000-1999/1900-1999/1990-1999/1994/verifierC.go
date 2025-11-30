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
1 2 2
6 5 5 5 10 4 10 1
10 5 7 7 9 6 9 8 9 5 1 1
6 14 6 7 7 9 3 9
3 7 4 1 3
6 5 3 9 9 6 9 9
3 14 7 9 6
10 11 6 8 3 7 8 9 4 8 5 8
9 16 6 8 8 6 10 9 8 8 4
6 5 10 5 8 5 5 9
9 16 9 10 10 7 5 4 8 9 6
10 2 6 1 4 2 1 10 1 5 10 4
2 16 3 5
4 6 1 7 1 1
6 11 3 4 1 2 2 2
1 1 1
6 8 3 3 3 9 1 7
10 1 4 3 1 1 6 10 2 5 6 8
1 9 8
9 19 1 5 7 10 3 8 4 2 6
2 0 8 3
9 18 7 8 9 6 3 6 5 5 10
7 20 1 9 3 1 5 1 3
3 5 2 8 4
9 1 4 4 8 2 5 2 10 4 10
10 11 5 7 5 9 1 3 1 7 7 3
2 16 2 4
2 3 1 3
4 3 4 1 9 8
8 9 9 7 4 4 7 7 9 1
10 18 1 7 9 10 3 2 8 6 1 9
2 19 6 5
6 9 1 7 2 2 5 4
1 14 1
7 20 8 8 4 10 10 2 1
5 0 6 5 2 4 8
4 3 10 6 7 8
3 11 7 2 5
2 3 2 10
6 20 7 4 2 1 10 8
1 15 5
6 14 3 6 5 8 9 8
7 15 5 7 4 3 8 10 5
9 13 2 10 10 2 2 6 3 9 3
7 2 2 1 3 5 7 4 6
8 5 9 5 2 3 9 7 2 6
9 7 9 5 3 3 8 4 7 6 10
3 14 8 1 10
7 5 7 9 1 8 5 7 5
7 20 8 6 9 6 2 4 9
10 6 7 7 1 6 8 9 8 3 2 1
7 6 10 10 7 4 2 7 9
4 8 10 10 4 8
10 4 1 10 7 8 5 9 10 3 8 4
2 11 1 8
9 2 10 8 6 8 5 9 8 1 2
10 11 3 7 5 3 1 3 8 7 8 5
3 0 5 9 8
1 11 1
9 12 10 8 4 5 8 3 8 9 5
2 8 6 5
6 20 5 7 9 2 9 4
7 19 9 3 9 2 5 1 4
8 17 4 9 5 1 2 2 7 6
4 10 6 2 6 8
6 5 8 8 5 8 3 8
4 8 6 3 2 4
8 6 6 3 6 3 3 4 5 9
7 12 6 5 10 9 10 6 7
5 17 10 2 6 5 7
8 5 5 6 8 8 2 3 6 7
3 0 2 6 3
6 2 7 1 9 6 4 10
7 17 5 8 3 6 6 4 8
2 4 4 6
5 4 7 6 5 2 6
4 7 4 10 1 6
6 20 10 1 3 3 2 7
8 8 3 6 9 10 2 6 10 7
4 1 7 8 8 10
6 17 10 10 2 10 9 9
8 12 8 3 7 7 9 8 1 2
8 18 3 2 9 3 2 7 5 8
1 8 2
6 7 3 1 3 7 2 6
8 1 8 4 2 8 3 9 1 3
9 17 1 1 4 9 1 9 6 9 4
3 11 8 1 3
9 3 4 2 8 4 1 10 4 7 6
10 20 7 9 9 3 9 2 3 4 3 7
4 9 6 7 3 7
3 12 6 5 2
9 3 8 5 5 9 8 5 4 7 3
9 3 1 10 9 4 4 4 7 10 1
3 20 1 5 8
9 1 4 3 10 6 1 4 2 3 9
3 2 8 5 4
3 10 5 9 10
2 13 7 1
8 9 2 5 1 4 7 6 5 9
`

type testCase struct {
	n   int
	k   int64
	arr []int64
}

func parseTests(raw string) ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	tests := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("bad line: %q", line)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("bad n in line: %q", line)
		}
		kVal, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("bad k in line: %q", line)
		}
		if len(fields) != n+2 {
			return nil, fmt.Errorf("expected %d numbers got %d in line: %q", n+2, len(fields), line)
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[2+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("bad value in line: %q", line)
			}
			arr[i] = v
		}
		tests = append(tests, testCase{n: n, k: kVal, arr: arr})
	}
	return tests, nil
}

func solve(tc testCase) string {
	n := tc.n
	k := tc.k
	a := make([]int64, n+2)
	b := make([]int64, n+2)
	f := make([]int64, n+3)
	for i := 1; i <= n; i++ {
		a[i] = tc.arr[i-1]
		b[i] = b[i-1] + a[i]
	}
	var ans int64
	j := n
	for i := n; i >= 1; i-- {
		for j > i && b[j-1]-b[i-1] > k {
			j--
		}
		if b[j]-b[i-1] > k {
			f[i] = int64(j-i) + f[j+1]
		} else {
			f[i] = int64(n - i + 1)
		}
		ans += f[i]
	}
	return strconv.FormatInt(ans, 10)
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
		want := solve(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
