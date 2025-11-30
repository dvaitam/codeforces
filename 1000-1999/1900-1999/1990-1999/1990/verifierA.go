package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesA = `13 97 53 5 33 65 62 51 100 38 61 45 74 27
17 17 36 17 96 12 79 32 68 90 77 18 39 12 93 9 87 42
16 71 12 45 55 40 78 81 26 70 61 56 66 33 7 70 1
3 92 51 90
1 78
16 42 31 93 41 90 8 24 72 28 30 18 69 57 11 10 40
17 62 13 38 70 37 90 15 70 42 69 26 77 70 75 36 56 11
20 49 40 73 30 37 23 24 23 4 78 84 33 60 8 11 86 96 16 19 4
3 89 69 87
13 90 67 35 66 30 27 86 75 53 74 35 57 63
12 10 41 78 14 62 75 80 42 24 31 2 93
9 14 90 28 47 21 42 54 7 12
5 89 28 5 73 81
18 77 87 9 3 15 81 24 77 73 15 50 11 47 14 4 77 2 24
6 91 15 61 26 93 7
1 69
14 79 12 33 8 28 9 82 38 44 55 23 7 64 59
2 76 12
13 25 33 45 93 60 72 21 89 86 26 98 7 100
6 20 43 67 32 15 76
18 56 58 68 48 49 42 98 71 55 8 7 17 45 11 23 82 9 77
7 28 98 22 69 56 17 54
16 80 73 73 3 63 79 36 19 80 42 58 91 14 25 34 44
17 54 46 88 9 79 70 45 2 19 98 84 90 54 37 24 35 65
12 5 79 69 64 74 90 60 40 12 28 59 42
15 66 90 37 61 61 95 75 75 42 11 24 2 44 75 78
19 13 94 40 29 61 16 7 93 79 20 81 31 43 20 95 22 65 23 48
10 55 99 97 76 16 40 7 42 12 58
13 52 35 6 14 14 5 72 100 48 24 56 56 33
3 73 13 62
16 41 87 64 74 56 18 64 56 37 90 35 31 7 88 85 46
5 40 58 92 10 100
4 23 78 25 66
18 18 5 72 51 76 66 72 63 74 79 5 87 3 14 88 26 54 73
7 86 79 56 50 67 60 71
7 91 17 3 5 66 41 66
2 49 57
15 58 68 48 30 29 91 24 37 33 58 51 39 78 100 30
12 20 98 91 21 78 95 65 14 98 31 84 41
2 94 55
15 100 3 41 36 48 71 97 58 66 45 69 68 37 80 50
9 21 36 34 39 24 44 6 79 22
3 100 23 86
4 46 63 59 3
9 60 71 42 15 26 29 44 56 65
8 22 21 34 32 54 37 90 50
6 48 35 19 8 41 62
9 11 77 81 25 80 19 73 23 29
14 58 83 66 87 15 78 26 99 26 61 40 4 42 100
1 21
5 1 1 1 1 1
12 40 91 47 55 15 40 37 46 99 16 19 90
15 14 31 96 66 2 26 55 21 82 93 76 16 14 90 3
2 96 10
17 15 8 79 32 78 28 100 3 58 23 79 100 92 48 51 37 32
7 51 31 45 83 96 39 33
13 47 61 96 29 11 29 33 21 91 15 73 92 29
9 41 74 16 19 32 20 4 84 13
3 7 88 7
10 10 47 99 3 57 40 66 71 22 27
12 40 5 62 44 56 16 3 42 5 35 53 10
5 49 1 60 40 81
2 63 64
1 97
7 35 81 76 39 76 81 89
5 33 89 35 53 93
16 45 2 10 8 63 64 98 17 23 20 57 63 56 18 99 36
14 95 45 16 67 39 35 66 83 24 70 96 17 71 67
12 48 69 47 5 15 20 10 99 11 56 36 21
5 37 40 42 6 19
19 100 86 4 6 7 92 56 99 70 38 8 20 13 88 98 41 45 91 56
1 94
5 31 30 79 54 39
15 46 100 34 83 82 45 33 49 42 98 32 69 83 42 11
17 17 17 44 44 55 45 76 71 43 44 73 34 50 84 76 43 83
15 31 60 86 81 13 59 70 76 52 89 51 6 58 53 53
12 21 79 62 92 30 14 23 19 20 25 53 67
6 9 76 89 96 54 91
19 90 88 7 2 85 21 10 17 67 21 88 75 78 19 45 63 48 61 74
13 36 75 73 96 13 17 89 40 44 33 83 55 53
5 22 30 22 34 36
12 3 55 97 50 7 22 2 26 55 37 31 5
14 32 90 59 73 54 48 4 93 7 16 39 6 83 71
16 29 18 88 95 76 17 75 44 89 93 51 2 86 74 43 81
5 29 7 20 9 36
10 79 55 55 24 12 52 21 87 17 69
11 6 60 68 60 21 26 3 79 93 66 60
4 27 53 10 55
3 69 11 6
6 35 12 11 76 64 90
6 72 40 77 80 94 68
8 13 58 67 71 82 54 62 58
20 39 77 33 39 1 6 54 9 23 91 89 93 19 25 58 19 40 15 45 64
9 56 39 10 56 45 18 32 97 18
10 6 22 10 6 61 19 6 70 29 58
11 39 76 99 32 42 78 40 79 78 56 27
16 78 24 42 38 12 5 38 74 11 74 44 31 80 67 45 19
5 56 1 40 53 22
11 33 74 82 10 25 53 78 99 75 27 100
12 83 100 14 41 42 37 97 57 92 9 95 43`

func expected(arr []int) string {
	cnt := make([]int, 101)
	for _, v := range arr {
		if v >= 0 && v < len(cnt) {
			cnt[v]++
		}
	}
	for i := 1; i <= 50; i++ {
		if cnt[i]%2 == 1 {
			return "YES"
		}
	}
	return "NO"
}

type testCase struct {
	arr []int
}

func parseTests() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesA), "\n")
	tests := make([]testCase, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 1 {
			return nil, fmt.Errorf("bad test line %d", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d length mismatch", i+1)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			val, err := strconv.Atoi(fields[1+j])
			if err != nil {
				return nil, err
			}
			arr[j] = val
		}
		tests[i] = testCase{arr: arr}
	}
	return tests, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.arr)))
		for _, v := range tc.arr {
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	tests, err := parseTests()
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}
	input := buildInput(tests)
	output, err := runCandidate(os.Args[1], input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	lines := strings.Fields(output)
	if len(lines) != len(tests) {
		fmt.Printf("expected %d outputs, got %d\n", len(tests), len(lines))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := expected(tc.arr)
		if lines[i] != want {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, want, lines[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
