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
	n   int
	arr []int
}

// Embedded testcases from testcasesA.txt.
const testcaseData = `
5 98 73 73 73 73
7 16 16 16 16 16 16 64
10 61 61 61 61 61 61 84 61 61 61
6 63 13 13 13 13 13
9 78 56 56 56 56 56 56 56 56
10 35 35 35 93 35 35 35 35 35 35
4 4 41 41 41
3 70 84 84
9 88 88 88 88 88 88 28 88 88
3 68 29 68
10 71 71 71 71 71 30 71 71 71 71
6 87 87 87 29 87 87
7 3 3 3 3 3 3 54
4 24 24 81 24
4 96 96 96 43
6 39 39 39 39 37 39
10 65 65 65 65 65 65 65 65 65 51
3 62 62 32
9 54 54 86 54 54 54 54 54 54
8 71 71 71 71 71 90 71 71
4 85 57 57 57
5 67 67 51 67 67
10 94 94 94 94 94 94 94 4 94 94
3 40 40 91
9 83 83 22 83 83 83 83 83 83
6 2 99 2 2 2 2
6 52 52 66 52 52 52
8 35 59 59 59 59 59 59 59
9 95 95 66 95 95 95 95 95 95
6 55 55 55 8 55 55
8 73 73 73 71 73 73 73 73
9 63 63 63 63 63 63 46 63 63
8 1 1 1 1 1 69 1 1
10 77 77 77 4 77 77 77 77 77 77
5 71 75 71 71 71
4 33 71 71 71
4 11 11 11 3
3 36 97 97
7 15 80 15 15 15 15 15
8 38 38 9 38 38 38 38 38
5 33 68 33 33 33
7 83 83 92 83 83 83 83
10 90 90 90 90 90 90 90 42 90 90
10 15 15 15 15 4 15 15 15 15 15
9 44 44 44 54 44 44 44 44 44
7 14 14 14 14 14 33 14
6 56 78 78 78 78 78
6 3 51 3 3 3 3
3 93 21 93
9 70 70 70 70 70 70 70 70 29
10 68 29 29 29 29 29 29 29 29 29
9 87 87 87 87 87 74 87 87 87
9 8 8 8 8 95 8 8 8 8
5 28 28 7 28 28
4 10 10 40 10
5 54 54 73 54 54
5 72 2 2 2 2
6 73 59 73 73 73 73
3 49 26 49
4 27 27 27 74
6 64 64 64 64 64 14
9 38 38 38 38 38 38 38 65 38
3 42 79 42
7 3 21 3 3 3 3 3
8 73 73 73 73 73 18 73 73
9 28 35 28 28 28 28 28 28 28
9 71 71 71 71 71 71 71 71 45
10 99 99 99 69 99 99 99 99 99 99
4 6 93 93 93
5 22 69 22 22 22
7 98 98 98 98 43 98 98
7 48 48 44 48 48 48 48
4 38 38 38 31
5 71 75 75 75 75
8 6 53 6 6 6 6 6 6
9 19 19 19 19 19 17 19 19 19
4 79 79 79 76
4 74 71 74 74
4 35 35 47 35
4 36 59 59 59
3 38 38 2
3 53 12 12
3 25 25 31
9 21 21 21 21 21 21 21 15 21
5 88 31 88 88 88
4 56 56 49 56
7 92 92 62 92 92 92 92
4 27 27 84 27
3 4 2 4
8 58 58 58 58 58 51 58 58
9 9 9 9 9 9 9 9 41 9
4 33 33 33 28
8 34 34 34 24 34 34 34 34
7 26 26 32 26 26 26 26
4 36 36 36 12
4 84 84 74 84
6 40 50 50 50 50 50
8 24 24 24 24 41 24 24 24
6 43 43 43 43 13 43
4 29 32 32 32
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
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		if len(fields)-1 != n {
			return nil, fmt.Errorf("case %d expected %d values, got %d", i+1, n, len(fields)-1)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(fields[j+1])
			if err != nil {
				return nil, fmt.Errorf("case %d bad value %d: %v", i+1, j+1, err)
			}
			arr[j] = v
		}
		res = append(res, testCase{n: n, arr: arr})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	return res, nil
}

// solve mirrors 1512A.go for one test case.
func solve(tc testCase) int {
	freq := make(map[int]int)
	for _, v := range tc.arr {
		freq[v]++
	}
	uniqueVal := 0
	for v, c := range freq {
		if c == 1 {
			uniqueVal = v
			break
		}
	}
	for i, v := range tc.arr {
		if v == uniqueVal {
			return i + 1
		}
	}
	return -1
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != strconv.Itoa(expect) {
			fmt.Printf("case %d failed: expected %d got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
