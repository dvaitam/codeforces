package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Embedded copy of testcasesD.txt to avoid external dependency.
const testcasesDData = `
5 85 35 44 0 48 61 78 85
3 87 1 62 0 34 87
9 39 7 38 0 10 15 26 31 34 35 36 39
6 11 1 4 0 1 5 8 10 11
8 64 26 63 0 9 24 29 57 60 62 64
3 14 2 10 0 4 14
6 96 28 67 0 50 54 65 74 96
7 78 38 65 0 4 30 36 44 75 78
4 99 21 91 0 73 74 99
3 93 14 88 0 35 93
6 25 2 18 0 3 12 16 21 25
3 62 29 39 0 2 62
6 64 27 35 0 3 39 40 49 64
2 58 23 45 0 58
9 45 17 25 0 1 3 5 7 20 35 39 45
5 62 10 50 0 10 17 45 62
2 53 11 35 0 53
4 58 13 43 0 34 56 58
8 92 39 83 0 14 35 56 65 72 80 92
5 48 14 31 0 20 34 36 48
6 11 4 9 0 1 6 7 10 11
3 17 6 14 0 12 17
6 72 2 10 0 3 33 48 59 72
6 85 39 60 0 23 24 41 47 85
7 86 17 56 0 4 14 17 49 73 86
6 74 15 57 0 24 31 35 42 74
8 93 45 52 0 14 29 42 43 77 87 93
9 31 2 13 0 7 9 15 19 21 24 29 31
5 25 1 18 0 7 11 19 25
4 45 11 17 0 23 40 45
4 63 10 44 0 51 55 63
5 69 23 64 0 38 53 54 69
2 62 5 18 0 62
2 71 33 61 0 71
9 38 2 32 0 5 8 15 19 22 34 35 38
5 15 1 14 0 4 9 12 15
8 83 4 6 0 16 22 31 39 62 65 83
2 77 35 62 0 77
2 88 8 52 0 88
4 42 18 34 0 4 23 42
5 35 4 22 0 8 11 16 35
6 26 1 17 0 2 13 19 21 26
6 41 9 37 0 1 4 21 31 41
2 26 1 5 0 26
2 18 8 9 0 18
3 75 33 65 0 41 75
4 50 3 26 0 25 42 50
8 85 20 67 0 16 17 25 34 43 55 85
8 10 4 5 0 1 3 6 7 8 9 10
2 57 21 53 0 57
7 63 31 58 0 2 14 16 30 35 63
5 98 38 43 0 17 29 55 98
2 51 12 48 0 51
6 25 8 12 0 13 17 22 24 25
3 50 19 37 0 7 50
2 70 10 26 0 70
8 15 5 7 0 2 3 7 10 11 14 15
2 53 4 6 0 53
3 96 31 68 0 75 96
5 21 1 20 0 4 8 17 21
9 22 9 10 0 3 6 8 11 18 19 21 22
5 68 26 43 0 45 48 51 68
10 63 31 37 0 16 25 27 33 48 54 60 62 63
4 63 23 60 0 38 49 63
10 97 31 51 0 13 20 21 52 62 64 83 96 97
9 66 12 21 0 19 26 30 35 38 41 53 66
5 49 1 19 0 13 25 31 49
4 82 24 40 0 42 62 82
4 63 23 54 0 39 45 63
5 69 2 64 0 6 10 52 69
8 39 8 29 0 5 9 13 14 16 17 39
4 89 44 47 0 22 33 89
2 50 6 34 0 50
3 20 2 5 0 9 20
6 14 3 11 0 6 10 11 12 14
2 13 3 9 0 13
8 58 16 21 0 14 26 32 38 42 48 58
4 79 21 29 0 10 36 79
7 24 8 17 0 4 12 17 22 23 24
9 47 22 44 0 7 17 22 37 42 43 44 47
10 77 8 72 0 8 20 23 24 38 46 66 73 77
7 93 30 38 0 14 19 43 72 83 93
8 81 20 62 0 9 23 24 40 59 62 81
3 33 13 25 0 7 33
6 44 13 15 0 3 9 31 33 44
6 41 17 29 0 22 26 29 35 41
3 55 16 24 0 10 55
6 85 7 22 0 15 24 25 73 85
8 95 26 43 0 19 25 51 70 76 78 95
10 31 10 16 0 1 7 9 10 12 26 27 28 31
9 62 31 58 0 20 21 25 36 38 41 58 62
9 77 20 76 0 1 4 14 25 30 62 63 77
4 77 30 43 0 25 68 77
5 14 7 12 0 2 8 11 14
6 94 10 28 0 7 12 60 80 94
2 56 20 35 0 56
8 19 8 17 0 1 2 3 5 11 12 19
3 53 26 33 0 5 53
5 65 15 47 0 6 14 41 65
7 19 4 16 0 3 6 13 16 18 19
`

type testCase struct {
	n int
	l int
	x int
	y int
	a []int
}

func parseTestcases() ([]testCase, error) {
	var cases []testCase
	reader := bufio.NewReader(strings.NewReader(testcasesDData))
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err != nil && len(line) == 0 {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			return nil, fmt.Errorf("line %d: expected at least 5 fields", lineNum)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", lineNum, err)
		}
		l, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad l: %v", lineNum, err)
		}
		x, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad x: %v", lineNum, err)
		}
		y, err := strconv.Atoi(fields[3])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad y: %v", lineNum, err)
		}
		if len(fields) != 4+n {
			return nil, fmt.Errorf("line %d: expected %d marks got %d", lineNum, n, len(fields)-4)
		}
		marks := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[4+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad mark %d: %v", lineNum, i+1, err)
			}
			marks[i] = val
		}
		cases = append(cases, testCase{n: n, l: l, x: x, y: y, a: marks})
	}
	return cases, nil
}

// solve mirrors 479D.go logic.
func solve(tc testCase) string {
	a := append([]int(nil), tc.a...)
	n := tc.n
	l := tc.l
	x := tc.x
	y := tc.y

	judge := func(d int) bool {
		r := 0
		for i := 0; i < n; i++ {
			for r < n && a[r]-a[i] < d {
				r++
			}
			if r < n && a[r]-a[i] == d {
				return true
			}
		}
		return false
	}

	f1 := judge(x)
	f2 := judge(y)
	if f1 && f2 {
		return "0"
	}
	if f1 || f2 {
		if f1 {
			return fmt.Sprintf("1\n%d", y)
		}
		return fmt.Sprintf("1\n%d", x)
	}

	ck := func(pos int) bool {
		if pos < 0 || pos > l {
			return false
		}
		idx := sort.Search(n, func(i int) bool { return a[i] >= pos+y })
		if idx < n && a[idx] == pos+y {
			return true
		}
		idx = sort.Search(n, func(i int) bool { return a[i] >= pos-y })
		if idx < n && a[idx] == pos-y {
			return true
		}
		return false
	}
	for i := 0; i < n; i++ {
		p := a[i] + x
		if ck(p) {
			return fmt.Sprintf("1\n%d", p)
		}
		p = a[i] - x
		if ck(p) {
			return fmt.Sprintf("1\n%d", p)
		}
	}
	return fmt.Sprintf("2\n%d %d", x, y)
}

func runCandidate(bin string, tc testCase) (string, error) {
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d %d %d\n", tc.n, tc.l, tc.x, tc.y)
	for i, v := range tc.a {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.Itoa(v))
	}
	input.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, tc := range tests {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
