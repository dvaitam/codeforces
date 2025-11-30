package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n      int
	weights []int
	marked map[int]bool
}

// Embedded testcases from testcasesA.txt.
const testcaseData = `
100
6 21 88 38 90 28 22
1 6
5
1 2
2 3
3 4
4 5
5 6
5 84 89 78 67 57
4 3 5 4 2
4
1 2
2 3
3 4
4 5
6 67 1 23 97 33 22
5 4 5 1 3 2
5
1 2
2 3
3 4
4 5
5 6
6 47 13 38 61 4 5
1 1
5
1 2
2 3
3 4
4 5
5 6
4 97 3 28 49
2 2 3
3
1 2
2 3
3 4
6 95 35 43 85 32 73
2 4 6
5
1 2
2 3
3 4
4 5
5 6
3 43 41 39
0
2
1 2
2 3
4 11 31 69 37
3 1 4 3
3
1 2
2 3
3 4
6 85 17 36 81 48 4
2 6 3
5
1 2
2 3
3 4
4 5
5 6
2 96 27
0
1
1 2
4 38 58 28 37
1 3
3
1 2
2 3
3 4
3 20 56 14
3 3 2 1
2
1 2
2 3
2 53 20
1 2
1
1 2
6 48 15 32 91 50 44
4 1 3 4 6
5
1 2
2 3
3 4
4 5
5 6
3 23 29 42
0
2
1 2
2 3
2 22 22
1 2
1
1 2
3 11 80 81
3 1 2 3
2
1 2
2 3
3 63 93 32
0
2
1 2
2 3
3 2 1 93
2 3 1
2
1 2
2 3
3 7 80 50
3 3 2 1
2
1 2
2 3
4 46 56 17 70
0
3
1 2
2 3
3 4
3 43 62 27
2 2 1
2
1 2
2 3
3 6 10 18
3 1 2 3
2
1 2
2 3
6 17 84 62 66 99 17
6 1 3 2 5 6 4
5
1 2
2 3
3 4
4 5
5 6
5 60 87 75 35 85
4 3 2 4 5
4
1 2
2 3
3 4
4 5
5 52 99 67 89 14
0
4
1 2
2 3
3 4
4 5
6 25 13 42 14 48 34
5 1 3 4 6 2
5
1 2
2 3
3 4
4 5
5 6
5 68 36 91 15 4
1 1
4
1 2
2 3
3 4
4 5
6 69 80 96 9 4 4
6 5 2 3 1 6 4
5
1 2
2 3
3 4
4 5
5 6
3 32 36 19
2 3 1
2
1 2
2 3
3 10 88 79
0
2
1 2
2 3
3 76 59 50
1 2
2
1 2
2 3
4 13 95 15 43
0
3
1 2
2 3
3 4
3 95 15 47
0
2
1 2
2 3
4 7 55 40 24
1 2
3
1 2
2 3
3 4
5 24 18 11 38 91
2 1 2
4
1 2
2 3
3 4
4 5
4 51 28 94 98
4 4 2 1 3
3
1 2
2 3
3 4
3 17 41 98
0
2
1 2
2 3
3 95 58 61
2 1 3
2
1 2
2 3
5 80 89 38 17 19
1 3
4
1 2
2 3
3 4
4 5
5 85 12 54 64 42
1 1
4
1 2
2 3
3 4
4 5
2 55 78
0
1
1 2
2 21 94
2 2 1
1
1 2
6 66 34 87 52 12 33
2 5 1
5
1 2
2 3
3 4
4 5
5 6
6 36 86 58 1 5 31
3 5 1 3
5
1 2
2 3
3 4
4 5
5 6
6 11 80 93 73 27 63
1 4
5
1 2
2 3
3 4
4 5
5 6
3 3 84 79
0
2
1 2
2 3
4 16 5 61 78
1 1
3
1 2
2 3
3 4
2 54 78
0
1
1 2
4 56 86 61 86
0
3
1 2
2 3
3 4
3 49 7 91
1 2
2
1 2
2 3
5 27 21 17 84 87
2 5 2
4
1 2
2 3
3 4
4 5
2 84 86
2 2 1
1
1 2
2 65 92
0
1
1 2
5 74 2 77 45 91
5 4 2 5 1 3
4
1 2
2 3
3 4
4 5
6 87 23 66 96 11 66
6 6 2 1 3 5 4
5
1 2
2 3
3 4
4 5
5 6
3 27 99 92
0
2
1 2
2 3
3 65 14 33
2 3 2
2
1 2
2 3
6 40 86 17 17 13 75
6 5 4 6 1 2 3
5
1 2
2 3
3 4
4 5
5 6
4 82 18 16 91
2 1 3
3
1 2
2 3
3 4
4 50 80 14 53
4 2 1 3 4
3
1 2
2 3
3 4
4 86 42 99 60
0
3
1 2
2 3
3 4
4 60 58 63 1
1 1
3
1 2
2 3
3 4
2 97 6
2 2 1
1
1 2
2 65 96
0
1
1 2
6 34 38 50 79 25 96
2 6 5
5
1 2
2 3
3 4
4 5
5 6
2 53 94
1 2
1
1 2
3 8 2 97
1 1
2
1 2
2 3
4 97 87 86 56
0
3
1 2
2 3
3 4
5 46 16 75 92 11
0
4
1 2
2 3
3 4
4 5
6 66 81 10 49 88 67
0
5
1 2
2 3
3 4
4 5
5 6
2 76 97
2 1 2
1
1 2
6 18 86 54 39 4 61
1 6
5
1 2
2 3
3 4
4 5
5 6
6 14 59 95 52 60 52
5 2 5 6 1 3
5
1 2
2 3
3 4
4 5
5 6
3 28 33 8
0
2
1 2
2 3
6 54 23 36 59 43 74
3 4 2 3
5
1 2
2 3
3 4
4 5
5 6
6 49 88 19 1 16 74
0
5
1 2
2 3
3 4
4 5
5 6
6 94 94 12 5 51 23
3 4 2 1
5
1 2
2 3
3 4
4 5
5 6
4 32 31 26 98
1 1
3
1 2
2 3
3 4
2 62 39
0
1
1 2
6 53 23 74 94 88 59
6 4 2 6 3 1 5
5
1 2
2 3
3 4
4 5
5 6
6 83 47 47 14 13 7
4 5 4 6 3
5
1 2
2 3
3 4
4 5
5 6
3 89 41 6
2 3 1
2
1 2
2 3
3 13 42 53
1 1
2
1 2
2 3
2 84 18
2 2 1
1
1 2
2 27 36
0
1
1 2
3 50 7 14
0
2
1 2
2 3
6 96 16 59 64 72 45
1 1
5
1 2
2 3
3 4
4 5
5 6
4 29 31 29 19
0
3
1 2
2 3
3 4
6 35 93 39 23 89 84
4 5 3 4 6
5
1 2
2 3
3 4
4 5
5 6
6 54 95 47 97 36 82
6 4 5 1 6 3 2
5
1 2
2 3
3 4
4 5
5 6
3 8 43 64
2 1 3
2
1 2
2 3
5 63 56 1 17 42
0
4
1 2
2 3
3 4
4 5
4 91 23 81 100
4 3 1 2 4
3
1 2
2 3
3 4
3 42 40 74
0
2
1 2
2 3
3 54 28 72
3 2 3 1
2
1 2
2 3
4 10 55 73 52
1 3
3
1 2
2 3
3 4
6 73 82 33 99 44 44
4 4 3 5 2
5
1 2
2 3
3 4
4 5
5 6
5 80 77 45 79 58
1 3
4
1 2
2 3
3 4
4 5
3 38 27 60
1 3
2
1 2
2 3`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	pos := 0
	nextInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected end of data")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}
	t, err := nextInt()
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	res := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		weights := make([]int, n)
		for j := 0; j < n; j++ {
			v, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d bad weight %d: %v", i+1, j+1, err)
			}
			weights[j] = v
		}
		k, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d bad k: %v", i+1, err)
		}
		marked := make(map[int]bool, k)
		for j := 0; j < k; j++ {
			v, err := nextInt()
			if err != nil {
				return nil, fmt.Errorf("case %d bad marked %d: %v", i+1, j+1, err)
			}
			marked[v] = true
		}
		m, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d bad m: %v", i+1, err)
		}
		for j := 0; j < m; j++ {
			if _, err := nextInt(); err != nil {
				return nil, fmt.Errorf("case %d bad edge u%d: %v", i+1, j+1, err)
			}
			if _, err := nextInt(); err != nil {
				return nil, fmt.Errorf("case %d bad edge v%d: %v", i+1, j+1, err)
			}
		}
		res = append(res, testCase{n: n, weights: weights, marked: marked})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra tokens at end of test data")
	}
	return res, nil
}

// solve mirrors 1524A.go: list unmarked vertices then output dummy edges.
func solve(tc testCase) string {
	unmarked := make([]int, 0)
	for i := 1; i <= tc.n; i++ {
		if !tc.marked[i] {
			unmarked = append(unmarked, i)
		}
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(unmarked)))
	for _, v := range unmarked {
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i := 1; i <= tc.n; i++ {
		sb.WriteString(fmt.Sprintf("1 %d\n", i))
	}
	return strings.TrimRight(sb.String(), "\n")
}

func runCandidate(bin string, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	input := strings.TrimSpace(testcaseData) + "\n"
	var expectedLines []string
	for _, tc := range tests {
		expectedLines = append(expectedLines, solve(tc))
	}
	expect := strings.TrimRight(strings.Join(expectedLines, "\n"), "\n")

	got, err := runCandidate(bin, input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if got != expect {
		fmt.Println("output mismatch")
		fmt.Println("expected:")
		fmt.Println(expect)
		fmt.Println("got:")
		fmt.Println(got)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
