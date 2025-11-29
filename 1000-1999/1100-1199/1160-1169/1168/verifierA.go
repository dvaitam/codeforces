package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
7 8 0 4 7 6 4 7 5
4 10 2 4 2 1
5 10 9 2 4 1 1
6 9 8 1 5 6 5 3
8 9 8 4 0 8 0 1 6 0
8 7 1 5 2 5 6 0 1 4
4 5 1 4 3 0
2 7 4 3
2 6 4 2
2 10 5 8
4 10 9 4 7 1
7 7 4 1 2 1 1 6 1
1 6 3
2 3 2 0
3 2 0 1 1
4 5 4 3 4 2
8 9 5 1 5 1 7 5 3 3
1 6 0
4 7 6 1 2 3
1 3 0
4 2 0 0 0 0
2 8 1 5
2 2 0 0
3 3 1 0 2
1 2 1
2 6 0 1
2 6 2 3
3 2 1 0 0
7 5 2 2 3 4 1 1 0
3 4 2 2 0
8 4 0 3 3 2 2 3 2 1
1 9 1
6 2 1 0 0 1 1 1
6 4 2 3 3 0 0 1
6 4 1 1 3 3 3 0
7 8 0 2 7 1 4 2 7
8 10 9 0 0 7 5 4 7 0
7 5 4 0 1 0 3 3 2
1 5 0
1 10 9
2 5 0 4
4 6 2 5 1 0
8 8 1 0 4 7 1 4 2 5
2 4 2 0
1 2 0
5 10 5 5 9 0 9
8 9 6 5 8 2 3 6 4 0
3 4 2 2 2
6 3 1 2 0 0 1 0
3 6 2 3 4
3 6 0 3 5
4 2 1 0 0 1
7 7 2 3 0 0 4 3 3
6 7 0 3 0 5 3 3
1 6 2
3 4 3 0 0
2 5 1 0
7 2 0 1 1 1 1 0 1
2 7 1 2
3 8 3 5 1
2 2 1 0
2 9 6 4
4 2 0 0 0 0
8 8 5 2 1 7 2 6 6 7
6 9 7 3 8 3 0 5
6 7 0 4 1 6 2 4
3 8 4 7 1
2 10 0 1
4 4 0 2 0 3
6 4 1 3 2 3 0 0
2 8 3 4
7 9 6 3 0 0 2 4 8
5 7 0 3 6 2 6
5 8 6 6 0 2 2
4 6 5 2 0 0
8 8 2 7 1 2 5 6 0 7
7 9 0 1 7 2 0 0 2
6 3 2 2 2 1 0 1
8 3 0 2 2 1 2 2 1 2
2 6 1 3
5 3 2 0 0 1 1
6 5 3 2 0 0 0 3
5 2 0 0 0 1 1
2 4 3 3
2 3 1 0
2 8 2 0
8 8 6 0 7 5 4 1 5 1
2 7 5 0
6 7 1 0 6 1 6 2
2 4 1 0
4 3 2 0 1 1
1 5 1
3 9 1 7 5
5 4 0 1 2 2 3
5 6 4 5 2 1 4
2 3 2 2
5 4 3 1 1 1 2
4 5 1 2 2 3
1 4 0
7 3 2 0 0 1 1 2 1`

type testCase struct {
	n   int
	m   int
	arr []int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ok checks if array can be made non-decreasing with x operations (from 1168A.go).
func ok(a []int, n, m, x int) bool {
	cur := 0
	for i := 0; i < n; i++ {
		ai := a[i]
		if ai+x < m {
			if ai < cur {
				return false
			}
			cur = ai
		} else {
			wrap := (ai + x) % m
			if wrap < cur {
				if ai < cur {
					return false
				}
				cur = ai
			}
		}
	}
	return true
}

// solve computes minimal x via binary search.
func solve(a []int, n, m int) int {
	lo, hi := -1, m
	for lo+1 < hi {
		mid := (lo + hi) >> 1
		if ok(a, n, m, mid) {
			hi = mid
		} else {
			lo = mid
		}
	}
	return hi
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	want := strconv.Itoa(solve(tc.arr, tc.n, tc.m))
	if got != want {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("invalid test data")
	}
	idx := 0
	total, err := strconv.Atoi(strings.TrimSpace(lines[idx]))
	if err != nil {
		return nil, err
	}
	idx++
	tests := make([]testCase, 0, total)
	for idx < len(lines) && len(tests) < total {
		line := strings.TrimSpace(lines[idx])
		idx++
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("invalid testcase line: %q", line)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		if len(fields) != n+2 {
			return nil, fmt.Errorf("expected %d numbers got %d", n+2, len(fields))
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[2+i])
			if err != nil {
				return nil, err
			}
			arr[i] = v
		}
		tests = append(tests, testCase{n: n, m: m, arr: arr})
	}
	if len(tests) != total {
		return nil, fmt.Errorf("expected %d testcases got %d", total, len(tests))
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
