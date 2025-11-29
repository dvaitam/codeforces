package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `1 9
5 7 7 3 7 8
1 4
2 0 3
3 8 4 0
4 5 3 3 3
3 0 9 1
1 1
2 5 9
6 1 8 6 4 10 3
2 7 7
4 3 9 1 8
5 2 1 5 9 2
1 4
6 6 4 3 2 1 10
2 10 3
3 4 0 1
2 8 7
3 2 1 1
2 9 8
2 3 2
4 3 7 9 7
1 5
4 8 2 0 6
6 9 9 2 9 2 6
3 8 8 3
6 3 5 3 9 7 10
1 6
5 10 7 0 8 5
2 5 2
1 5
4 1 8 9 6
3 9 8 2
2 10 0
4 1 3 7 4
5 4 4 0 9 5
3 10 1 0
5 0 7 0 2 1
1 4
3 1 1 1
2 9 9
1 10
3 5 9 9
4 0 8 10 9
6 7 6 4 9 0 3
3 5 5 2
4 3 6 0 0
3 9 8 3
3 2 9 0
6 6 4 3 0 0 9
2 7 9
2 3 9
6 5 4 10 6 8 1
5 7 3 4 4 0
3 9 4 9
3 3 4 5
5 4 4 6 7 1
6 1 6 3 10 2 10
6 3 4 5 10 8 7
6 7 9 6 3 8 7
3 1 3 9
4 4 9 1 0
1 1
6 3 8 6 10 5 8
1 1
4 4 7 7 0
3 8 10 8
4 5 9 7 2
6 8 3 10 0 7 9
2 3 9
2 10 0
3 1 6 5
4 5 2 9 6
4 2 6 2 3
1 5
4 10 7 1 10
3 4 9 2
6 4 0 3 8 4 4
2 1 7
1 8
4 5 2 4 8
1 0
2 3 5
5 2 5 5 3 4
2 10 8
1 10
4 4 4 6 4
1 10
1 0
5 2 3 8 0 8
6 2 9 9 4 10 2
2 9 7
1 10
4 9 4 8 3
6 10 5 1 3 3 2
2 4 8
4 4 7 4 3
3 3 2 7
4 2 8 3 6
4 4 3 5 10`

type testCase struct {
	n   int
	arr []int
}

func median3(a, b, c int) int {
	if a > b {
		if b > c {
			return b
		} else if a > c {
			return c
		}
		return a
	}
	if a > c {
		return a
	} else if b > c {
		return c
	}
	return b
}

// Embedded reference logic from 1322E.go.
func solve(tc testCase) (int, []int) {
	n := tc.n
	a0 := append([]int(nil), tc.arr...)
	if n <= 2 {
		return 0, a0
	}

	b1 := make([]int, n)
	b1[0] = a0[0]
	b1[n-1] = a0[n-1]
	changed := false
	for i := 1; i < n-1; i++ {
		v := median3(a0[i-1], a0[i], a0[i+1])
		b1[i] = v
		if v != a0[i] {
			changed = true
		}
	}
	if !changed {
		return 0, a0
	}

	safe := make([]bool, n)
	for i := 0; i < n; i++ {
		if i == 0 || i == n-1 {
			safe[i] = true
			continue
		}
		bi := b1[i]
		if (bi > b1[i-1] && bi > b1[i+1]) || (bi < b1[i-1] && bi < b1[i+1]) {
			safe[i] = false
		} else {
			safe[i] = true
		}
	}

	time2 := make([]int, n)
	finalVal := make([]int, n)
	for i := 0; i < n; i++ {
		if safe[i] {
			time2[i] = 0
			finalVal[i] = b1[i]
		} else {
			time2[i] = -1
		}
	}

	q := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if safe[i] {
			q = append(q, i)
		}
	}
	maxT := 0
	for head := 0; head < len(q); head++ {
		u := q[head]
		t := time2[u]
		for _, d := range []int{-1, 1} {
			v := u + d
			if v < 0 || v >= n {
				continue
			}
			if time2[v] == -1 {
				time2[v] = t + 1
				finalVal[v] = finalVal[u]
				if time2[v] > maxT {
					maxT = time2[v]
				}
				q = append(q, v)
			} else if time2[v] == t+1 && v > 0 && v < n-1 {
				bi := b1[v]
				if bi > b1[v-1] && bi > b1[v+1] {
					if finalVal[u] > finalVal[v] {
						finalVal[v] = finalVal[u]
					}
				} else if bi < b1[v-1] && bi < b1[v+1] {
					if finalVal[u] < finalVal[v] {
						finalVal[v] = finalVal[u]
					}
				}
			}
		}
	}

	return 1 + maxT, finalVal
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	var res []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid n", idx+1)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d: expected %d numbers, got %d", idx+1, n, len(fields)-1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[i+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid value", idx+1)
			}
			arr[i] = val
		}
		res = append(res, testCase{n: n, arr: arr})
	}
	return res, nil
}

func run(bin string, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid test data:", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expSteps, expArr := solve(tc)
		input := fmt.Sprintf("%d\n", tc.n)
		for i, v := range tc.arr {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"

		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(got)
		if len(fields) < tc.n+1 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers, got %d\n", idx+1, tc.n+1, len(fields))
			os.Exit(1)
		}
		gotSteps, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid step count\n", idx+1)
			os.Exit(1)
		}
		if gotSteps != expSteps {
			fmt.Fprintf(os.Stderr, "case %d failed: expected steps %d got %d\n", idx+1, expSteps, gotSteps)
			os.Exit(1)
		}
		for i := 0; i < tc.n; i++ {
			val, err := strconv.Atoi(fields[i+1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: invalid number in output\n", idx+1)
				os.Exit(1)
			}
			if val != expArr[i] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %v\n", idx+1, expArr, fields[1:])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
