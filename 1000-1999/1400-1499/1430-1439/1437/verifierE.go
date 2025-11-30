package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const infInt64 = int64(2e18)

// solveCase mirrors 1437E.go for a single test.
func solveCase(n int, k int, a []int64, fixed []int) (int64, error) {
	if len(a) != n {
		return 0, fmt.Errorf("len(a) %d != n %d", len(a), n)
	}
	if len(fixed) != k {
		return 0, fmt.Errorf("len(fixed) %d != k %d", len(fixed), k)
	}
	arr := make([]int64, n+2)
	arr[0] = -infInt64
	arr[n+1] = infInt64
	for i := 1; i <= n; i++ {
		arr[i] = a[i-1]
	}
	bpos := make([]int, 0, k+2)
	bpos = append(bpos, 0)
	bpos = append(bpos, fixed...)
	bpos = append(bpos, n+1)

	var res int64
	for i := 1; i < len(bpos); i++ {
		l, r := bpos[i-1], bpos[i]
		if arr[l] >= arr[r] {
			return -1, nil
		}
		d := make([]int64, 0, r-l-1)
		for j := l + 1; j < r; j++ {
			v := arr[j]
			if v <= arr[l] || v >= arr[r] {
				continue
			}
			idx := sort.Search(len(d), func(x int) bool { return d[x] >= v })
			if idx == len(d) {
				d = append(d, v)
			} else {
				d[idx] = v
			}
		}
		segLen := int64(r - l - 1)
		res += segLen - int64(len(d))
	}
	return res, nil
}

// referenceSolve reproduces I/O behaviour of 1437E.go.
func referenceSolve(input string) (string, error) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return "", fmt.Errorf("empty input")
	}
	pos := 0
	nextInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected EOF")
		}
		val, err := strconv.Atoi(fields[pos])
		pos++
		return val, err
	}
	n, err := nextInt()
	if err != nil {
		return "", fmt.Errorf("bad n: %v", err)
	}
	k, err := nextInt()
	if err != nil {
		return "", fmt.Errorf("bad k: %v", err)
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		v, err := nextInt()
		if err != nil {
			return "", fmt.Errorf("bad a[%d]: %v", i, err)
		}
		a[i] = int64(v)
	}
	fixed := make([]int, k)
	for i := 0; i < k; i++ {
		v, err := nextInt()
		if err != nil {
			return "", fmt.Errorf("bad fixed[%d]: %v", i, err)
		}
		fixed[i] = v
	}
	if pos != len(fields) {
		return "", fmt.Errorf("extra tokens at end")
	}
	ans, err := solveCase(n, k, a, fixed)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(ans, 10), nil
}

type testCase struct {
	input string
}

// Embedded testcases from testcasesE.txt.
const testcaseData = `
7 2 6 17 48 47 47 25 21 3 6
3 3 22 42 13 1 2 3
2 1 5 9 1
5 1 31 20 38 3 29 4
5 1 42 33 25 6 39 5
2 1 41 27 2
5 1 3 2 18 40 26 2
2 0 42 40
4 4 33 49 12 10 1 2 3 4
5 0 16 38 8 5 42
1 0 16
2 0 22 27
5 3 6 24 38 6 21 1 3 4
3 0 5 4 19
3 0 11 7 24
3 2 11 30 43 1 2
6 5 16 6 8 27 11 32 1 3 4 5 6
7 0 46 7 2 31 10 14 22
5 0 6 33 34 17 9
3 2 28 8 48 1 2
6 2 36 45 20 27 39 21 3 4
1 0 32
6 1 8 43 3 14 17 15 1
6 4 15 46 13 33 19 28 3 4 5 6
6 4 27 14 14 44 1 46 2 3 4 6
7 3 3 36 28 27 35 22 14 1 2 3
2 0 22 4
2 2 5 34 1 2
6 0 8 10 25 40 20 2
4 0 13 39 47 39
5 4 1 15 11 1 37 1 2 3 4
7 6 21 7 11 8 36 10 31 1 2 4 5 6 7
1 1 17 1
4 1 18 44 39 36 4
4 2 29 15 39 35 1 4
4 3 13 8 48 24 1 2 3
5 5 5 2 46 27 23 1 2 3 4 5
1 1 28 1
3 3 39 18 25 1 2 3
1 0 12
1 1 39 1
5 3 28 13 26 4 23 3 4 5
6 3 37 39 26 37 17 39 1 2 5
3 1 41 43 36 3
6 4 26 17 9 42 37 6 2 3 4 6
5 0 32 14 34 19 13
1 1 7 1
6 6 24 47 25 24 47 39 1 2 3 4 5 6
1 1 8 1
3 1 22 22 46 2
3 1 21 22 6 2
3 2 41 19 10 2 3
2 2 7 15 1 2
6 1 19 9 18 4 1 16 3
3 0 45 29 41
1 1 10 1
2 0 22 6
1 1 32 1
3 3 48 24 38 1 2 3
7 1 9 26 18 6 35 38 46 7
3 1 46 31 27 1
3 2 18 49 22 1 2
2 0 19 50
4 3 22 30 47 7 1 3 4
2 1 13 38 2
7 4 3 9 1 45 24 29 19 3 5 6 7
2 2 1 38 1 2
1 1 29 1
4 3 50 6 3 30 1 2 4
2 2 31 29 1 2
6 3 35 20 32 31 29 21 1 3 5
1 0 5
2 2 28 29 1 2
7 0 35 12 10 50 43 27 36
6 2 11 45 22 11 16 31 2 4
1 1 11 1
1 0 22
7 0 45 8 20 20 18 12 38
6 1 26 31 46 36 2 18 6
7 6 37 8 20 24 35 36 25 1 3 4 5 6 7
3 3 2 16 2 1 2 3
4 1 20 3 7 7 3
4 2 48 9 28 17 2 3
4 3 38 34 20 37 1 2 4
6 5 12 43 30 19 9 7 2 3 4 5 6
3 0 21 8 1
2 0 37 45
3 0 18 39 42
4 3 8 11 17 44 1 2 4
2 0 4 1
2 0 47 44
1 1 6 1
2 2 23 46 1 2
7 3 46 33 45 36 38 11 30 1 2 4
5 4 40 28 37 46 32 1 2 3 4
2 2 11 21 1 2
7 5 42 40 28 43 27 19 14 3 4 5 6 7
6 1 47 30 6 1 44 44 5
4 4 3 21 22 50 1 2 3 4
1 0 45
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
			return nil, fmt.Errorf("line %d: not enough tokens", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		k, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad k: %v", i+1, err)
		}
		if len(fields) != 2+n+k {
			return nil, fmt.Errorf("line %d: expected %d tokens, got %d", i+1, 2+n+k, len(fields))
		}
		// validate integers
		for idx := 2; idx < len(fields); idx++ {
			if _, err := strconv.Atoi(fields[idx]); err != nil {
				return nil, fmt.Errorf("line %d: bad int at %d: %v", i+1, idx+1, err)
			}
		}
		res = append(res, testCase{input: line + "\n"})
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expected, err := referenceSolve(tc.input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed\ninput: %sexpected: %s\ngot: %s\n", idx+1, tc.input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
