package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `
8 97 53 5 33 65 62 51 100
6 61 45 74 27 64 17
6 17 96 12 79 32 68
4 39 12 93 9
7 60 71 12 45 55 40 78
5 70 61 56 66 33
2 70 1
3 92 51 90
2 78 63
7 31 93 41 90 8 24 72
5 30 18 69 57 11
3 40 65 62
3 38 70 37
3 70 42 69
5 77 70 75 36 56
3 76 49 40
5 37 23 24 23 4
6 60 8 11 86 96 16
4 4 10 89 69
8 90 67 35 66 30 27 86 75
8 74 35 57 63 84 82 89 45
3 41 78 14
9 75 80 42 24 31 2 93 34 14
5 47 21 42 54 7
3 100 18 89
5 5 73 81 68 77
3 3 15 81
5 77 73 15 50 11
7 14 4 77 2 24 23 91
3 61 26 93
2 86 2
10 54 79 12 33 8 28 9 82 38 44
8 23 7 64 59 5 76 12 89
8 25 33 45 93 60 72 21 89
5 98 7 100 86 20
4 43 67 32 15
9 85 22 1 60 87 52 72 65 39
7 49 84 32 19 71 88 1
9 94 10 42 94 5 69 35 17 30
9 45 78 36 86 45 75 81 79 16
5 36 79 100 55 23
5 46 34 10 67 78
5 64 14 60 88 72
5 32 1 40 7 74
5 46 39 92 41 75
8 50 60 12 62 77 86 44 99
3 34 16 19
10 71 51 1 63 81 75 68 81 12 57
5 73 63 69 83 78
4 77 37 49 66
8 84 82 81 40 43 30 79 6
5 36 63 22 16 18
7 37 5 11 94 27 21 75
4 42 83 86 67
10 99 23 12 49 21 73 62 47 99 42
8 71 47 39 39 29 47 82 82
6 50 52 61 64 27 69
2 75 64
7 55 23 76 5 31 29 21
3 57 35 76
9 46 71 40 46 14 30 22 34 57
2 23 45
4 89 86 91 40
3 57 87 13
10 10 63 22 13 81 10 59 50 47 21
4 86 4 12 22
10 34 75 21 76 40 61 38 22 29 19
8 24 40 95 70 96 48 19 76
5 57 21 33 93 34
7 65 24 3 28 26 51 83
8 63 9 8 78 90 93 89 92
4 9 27 20 32
10 40 62 73 60 14 57 89 26 41 21
8 56 79 63 28 17 51 57 31
2 23 23
6 73 92 8 57 40 51
4 5 44 51 65
7 48 81 1 77 38 30 65
10 65 16 4 92 63 96 17 80 99 46
3 45 76 56
2 34 21
5 10 21 100 59 57
4 3 19 66 60
7 41 42 39 58 100 97 89
5 18 64 16 34 90
9 89 71 62 67 41 70 49 35 31
3 52 83 32
5 35 33 41 12 72
10 46 99 12 33 37 74 57 74 49 100
4 100 79 80 55
3 70 89 79
7 62 21 27 26 6 15 13
6 55 19 17 16 100 43
2 2 66
10 21 5 21 64 47 96 13 20 57 98
2 9 53
10 47 6 76 54 98 49 73 26 24 47
10 12 41 30 77 43 35 49 88 32 55
6 33 2 37 65 59 85
3 58 77 33
5 80 73 10 10 98
6 84 41 57 80 90 4
2 71 69
`

type testCase struct {
	arr   []int
	input string
}

// solve implements the logic from 1869A.go: it returns a sequence of ranges.
func solve(n int) [][2]int {
	if n%2 == 0 {
		return [][2]int{{1, n}, {1, n}}
	}
	return [][2]int{{1, n - 1}, {1, n - 1}, {2, n}, {2, n}}
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func applyOps(arr []int, ops [][2]int) error {
	n := len(arr)
	for _, op := range ops {
		l := op[0] - 1
		r := op[1] - 1
		if l < 0 || r < l || r >= n {
			return fmt.Errorf("invalid range %d %d", op[0], op[1])
		}
		x := 0
		for i := l; i <= r; i++ {
			x ^= arr[i]
		}
		for i := l; i <= r; i++ {
			arr[i] = x
		}
	}
	for _, v := range arr {
		if v != 0 {
			return fmt.Errorf("array not zero after ops")
		}
	}
	return nil
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	cases := []testCase{}
	pos := 0
	for pos < len(fields) {
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("bad n at token %d: %w", pos, err)
		}
		pos++
		if pos+n > len(fields) {
			return nil, fmt.Errorf("not enough elements for n=%d", n)
		}
		arr := make([]int, n)
		strs := make([]string, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[pos+i])
			if err != nil {
				return nil, fmt.Errorf("bad value at token %d: %w", pos+i, err)
			}
			arr[i] = v
			strs[i] = fields[pos+i]
		}
		pos += n
		input := fmt.Sprintf("1\n%d\n%s\n", n, strings.Join(strs, " "))
		cases = append(cases, testCase{arr: arr, input: input})
	}
	return cases, nil
}

func parseOutput(out string) ([][2]int, error) {
	parts := strings.Fields(out)
	if len(parts) == 0 {
		return nil, fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid k: %w", err)
	}
	if k < 0 || k > 8 {
		return nil, fmt.Errorf("invalid k: %d", k)
	}
	if len(parts) != 1+2*k {
		return nil, fmt.Errorf("expected %d numbers, got %d", 1+2*k, len(parts))
	}
	ops := make([][2]int, k)
	for i := 0; i < k; i++ {
		l, err := strconv.Atoi(parts[1+2*i])
		if err != nil {
			return nil, fmt.Errorf("invalid l: %w", err)
		}
		r, err := strconv.Atoi(parts[2+2*i])
		if err != nil {
			return nil, fmt.Errorf("invalid r: %w", err)
		}
		ops[i] = [2]int{l, r}
	}
	return ops, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	_ = solve // keep embedded reference implementation available

	for idx, tc := range cases {
		out, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		ops, err := parseOutput(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		arrCopy := make([]int, len(tc.arr))
		copy(arrCopy, tc.arr)
		if err := applyOps(arrCopy, ops); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
