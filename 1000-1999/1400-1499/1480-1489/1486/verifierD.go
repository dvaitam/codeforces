package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func check(a []int, k, mid int) bool {
	n := len(a)
	prefix := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if a[i-1] >= mid {
			prefix[i] = prefix[i-1] + 1
		} else {
			prefix[i] = prefix[i-1] - 1
		}
	}
	minPrefix := 0
	for i := k; i <= n; i++ {
		if prefix[i]-minPrefix > 0 {
			return true
		}
		if prefix[i-k+1] < minPrefix {
			minPrefix = prefix[i-k+1]
		}
	}
	return false
}

func expected(a []int, k int) int {
	low, high := 1, len(a)
	res := 1
	for low <= high {
		mid := (low + high) / 2
		if check(a, k, mid) {
			res = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return res
}

type testCase struct {
	n    int
	k    int
	arr  []int
}

// Embedded testcases from testcasesD.txt.
const testcaseData = `
1 1 1
6 2 6 6 3 3 5 2
10 1 10 3 7 7 9 6 9 8 9 5
1 1 1
8 6 7 7 3 3 4 4 1 3
6 2 2 5 5 3 5 6
9 3 8 7 9 6 6 6 8 3 7
8 4 8 5 8 6 8 8 6 8
8 4 6 3 5 8 5 5 7 5
4 4 3 1 3 1
4 1 1 1 3 2
2 1 2 1
4 1 4 1 1 3
6 2 2 6 1 1 1 1
1 1 1
6 3 2 2 6 2 5 6
1 1 1
4 2 1 1 3 1
5 3 4 1 3 4 5
10 1 5 7 10 3 8 4 2 6 2 1
8 3 7 8 6 3 6 5 5 7
1 1 1
5 1 2 2 2 1 4
4 1 2 2 4 1
5 1 5 2 5 5 3
5 4 3 5 1 2 1
7 4 2 1 5 6 1 2 1
2 1 1 1
2 1 1 2
8 5 7 4 4 7 7 1 1 7
9 3 2 8 6 1 9 2 6 5 6
5 1 4 1 1 3 2
1 1 1
7 6 4 4 2 5 5 1 1
5 1 3 3 1 2 4
4 1 3 4 4 2
6 4 1 3 1 1 1 5
6 6 4 2 6 1 1 5
8 1 8 5 6 8 3 6 5 8
9 8 7 8 5 7 4 3 8 5 9
7 6 6 6 1 5 6 7 5
2 1 2 1
9 3 7 2 2 1 3 5 7 4 6
8 3 5 2 3 7 2 6 4 5
3 1 2 3 1
7 7 3 7 7 5 6 2 4
8 1 7 3 7 1 8 5 7 5
7 6 6 4 3 5 3 6 6
2 1 1 2
7 6 1 3 4 5 6 4 6
3 1 1 2 1
10 10 7 4 2 7 9 4 5 10 10 4
8 3 1 7 8 5 3 8 4 2
6 1 4 5 6 6 1 5
8 6 8 5 8 1 2 6 3 7
5 2 1 2 4 4 4
5 2 1 3 5 4 1
6 1 5 4 5 4 2 6
5 4 2 4 5 3 1
5 3 3 3 3 4 5
2 1 2 1
9 2 5 1 4 8 9 4 9 5 1
2 1 2 2
4 3 3 1 3 4
6 2 4 4 3 4 2 6
8 4 5 6 3 2 4 8 4 6
3 2 1 1 1
5 5 4 4 3 3 5
9 6 7 5 9 2 6 5 7 8 3
5 3 4 4 1 2 3
7 2 1 1 3 2 3 1 6
7 1 5 3 2 7 7 5 4
9 5 8 3 6 6 4 8 2 3 4
6 3 2 4 3 3 1 3
4 2 2 1 3 3
10 1 3 3 2 7 8 5 3 6 9 10
2 2 2 1
1 1 1
8 6 2 8 7 8 3 7 7 8
1 1 1
10 3 2 9 3 2 7 5 8 1 5 2
6 2 2 1 2 4 6 1
6 6 4 1 4 2 1 4
3 3 1 1 3
9 9 1 1 4 9 1 9 6 9 4
3 2 2 1 1
9 2 4 2 8 4 1 4 7 6 7
9 9 3 9 2 3 4 3 7 4 5
6 4 2 4 2 4 3 3
2 1 2 2
5 5 4 3 2 4 2
9 2 1 9 4 4 4 7 1 3 1
5 4 5 1 2 2 5
6 1 6 2 1 2 6 6
9 3 2 8 5 4 3 6 5 9 2
7 4 6 6 1 4 3 6 7
2 2 1 1
7 3 3 5 6 4 5 5 6
4 4 2 2 4 4
6 4 4 5 3 5 2 5
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: not enough fields", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		k, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad k: %v", i+1, err)
		}
		if len(fields) != 2+n {
			return nil, fmt.Errorf("line %d: expected %d numbers, got %d", i+1, 2+n, len(fields))
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			val, err := strconv.Atoi(fields[2+j])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad value %d: %v", i+1, j+1, err)
			}
			arr[j] = val
		}
		res = append(res, testCase{n: n, k: k, arr: arr})
	}
	return res, nil
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		want := fmt.Sprintf("%d", expected(tc.arr, tc.k))
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
