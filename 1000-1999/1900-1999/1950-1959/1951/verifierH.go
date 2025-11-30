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
3 19 50 83 6 9 68 12 46
5 7 64 27 4 11 55 53 8 30 11 70 54 7 72 15 28 80 80 74 7 73 74 50 6 28 5 71 17 37 53 18 69
1 73 39
5 87 23 13 74 73 81 24 47 12 70 91 8 72 7 79 26 63 87 68 54 99 40 59 74 58 46 38 31 23 89 99 31
1 73 38
5 63 43 93 57 36 77 9 15 65 53 21 96 43 19 62 53 5 85 9 97 71 73 40 43 88 44 76 63 74 58 8 11
3 60 89 85 8 7 93 89 39
5 87 57 36 91 49 85 44 2 59 45 21 78 14 63 7 27 98 36 16 94 31 50 50 63 10 21 57 51 70 35 17 55
5 35 90 53 45 87 48 29 19 10 22 19 29 84 29 1 62 75 23 33 36 0 18 53 68 47 78 72 40 16 88 65 79
1 58 99
5 50 50 51 50 13 61 81 51 7 24 8 26 56 20 14 43 76 6 13 0 72 19 68 12 46 78 3 9 26 78 48 19
3 44 77 46 60 15 14 62 59
4 61 39 10 18 13 95 43 94 33 61 88 20 66 2 26 67
3 18 88 69 3 97 67 38 82
1 89 33
5 46 21 45 98 28 68 69 99 64 42 81 28 78 100 97 24 30 51 94 29 25 66 63 45 93 3 3 35 60 33 24 88
5 44 57 92 44 46 10 28 13 29 60 25 43 26 61 79 78 0 61 83 44 82 10 84 15 49 100 91 96 25 61 22 55
3 11 92 50 59 51 95 10 92
2 21 16 3 19
5 59 83 18 78 76 60 84 44 19 70 70 16 2 1 92 83 13 67 95 17 55 24 27 3 32 27 37 64 30 97 75 41
3 69 53 16 7 94 45 58 84
5 66 53 64 16 68 19 67 65 2 56 99 23 77 0 99 19 22 18 60 79 92 15 71 7 41 87 66 67 71 61 100 99
1 71 7
2 24 35 5 98
1 64 57
5 3 97 8 56 41 78 64 77 65 25 88 35 57 65 68 61 64 31 89 66 33 71 25 57 17 53 15 50 56 40 9 85
2 54 9 27 85
3 100 15 99 19 91 82 84 46
2 32 17 59 28
1 50 62
2 85 28 20 90
4 65 51 43 53 25 45 40 11 92 46 2 43 70 58 56 90
1 49 42
5 79 37 65 8 14 100 29 13 10 33 34 5 99 23 34 96 16 54 86 33 51 19 68 65 73 63 89 41 11 35 7 88
2 54 9 34 2
1 33 10
5 28 8 33 15 58 1 43 70 53 34 79 16 5 67 90 30 14 20 33 6 23 25 39 80 39 67 97 26 37 57 64 86
2 34 44 2 32
1 1 2
5 70 24 65 60 31 57 13 84 83 55 84 63 69 50 64 39 88 27 29 43 25 90 93 81 17 51 44 6 16 1 9 80
3 55 20 7 10 85 48 64 85
3 76 31 88 37 5 58 23 20
3 57 0 33 46 42 70 41 31
1 39 27
3 23 0 42 48 10 60 35 64
2 31 64 99 0
1 33 11
2 51 75 5 50
1 38 38
2 10 74 67 96
2 84 91 100 76
4 97 41 92 63 19 36 92 79 82 18 5 91 65 80 54 93
5 17 67 96 64 72 2 87 74 91 87 88 82 29 10 3 5 17 81 46 13 48 57 71 6 80 2 80 68 87 31 62 33
1 58 8
5 68 11 84 67 8 95 94 60 32 9 33 30 93 96 26 29 94 83 58 63 48 9 61 87 36 98 5 78 80 82 25 9
5 18 42 32 83 95 88 38 79 72 17 1 61 7 62 34 86 12 88 27 86 62 37 90 66 36 59 59 59 98 15 70 25
3 10 60 2 37 58 9 64 57
3 49 26 26 9 74 11 18 95
5 33 46 16 77 80 65 35 14 90 46 29 63 62 50 3 20 0 62 87 57 51 38 93 18 53 44 48 40 15 42 0 41
3 50 15 25 91 1 94 37 32
3 8 50 49 75 9 46 54 96
3 6 35 13 6 84 36 81 19
2 34 55 65 40
2 98 47 100 54
1 97 80
4 70 70 26 92 10 6 93 52 57 78 96 17 82 36 62 6
5 16 21 60 53 43 36 38 32 94 94 83 33 51 83 30 38 61 71 85 50 15 21 82 20 9 26 64 63 70 28 57 42
4 54 17 70 24 31 11 22 43 71 11 40 30 47 33 72 25
1 95 52
4 52 95 67 26 48 34 43 96 7 63 35 73 46 16 87 64
5 80 27 11 34 31 49 51 82 57 55 39 2 16 4 54 90 97 60 75 62 0 9 50 67 59 57 31 100 13 28 19 19
5 87 13 92 89 82 97 58 10 70 99 5 0 100 16 29 72 4 82 91 38 16 80 32 67 81 55 89 97 14 12 9 38
5 74 24 49 33 28 76 0 1 68 38 58 35 40 82 31 60 67 30 70 31 3 52 90 83 39 7 2 24 63 86 82 53
1 32 29
4 47 29 63 4 89 43 91 53 46 87 50 25 0 37 94 64
1 26 63
2 39 98 24 29
4 28 33 97 37 13 79 63 78 23 28 62 53 85 7 76 18
4 6 27 3 76 18 53 6 90 7 23 50 57 91 40 93 14
1 21 42
2 23 83 67 95
4 4 39 85 92 48 47 42 56 21 13 0 10 35 10 44 53
1 71 97
2 48 45 98 39
4 11 6 90 60 25 47 69 57 24 41 46 94 60 3 80 52
2 80 98 51 5
4 4 59 8 7 32 24 95 8 77 43 46 34 42 78 5 33
3 35 38 0 92 96 76 81 8
1 29 13
4 91 59 99 49 32 55 63 16 63 23 1 94 38 88 98 19
5 30 41 40 58 46 100 100 76 10 65 25 50 96 20 31 52 8 83 4 61 70 69 41 20 54 13 9 33 79 10 26 12
4 63 90 57 22 29 17 53 58 79 86 30 95 68 99 85 97
1 99 37
3 35 72 34 47 32 94 33 25
4 31 23 31 30 19 36 74 24 41 8 50 32 31 64 67 29
1 83 59
1 13 0
4 29 57 47 5 37 29 15 6 24 76 74 24 9 47 65 22
4 77 33 99 99 85 0 13 81 76 90 79 44 27 4 47 43
2 5 26 32 4
`

type testCase struct {
	k   int
	arr []int
}

func solve(k int) []int {
	n := 1 << k
	res := make([]int, k)
	for t := 1; t <= k; t++ {
		res[t-1] = n - (1 << t) + 1
	}
	return res
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	cases := make([]testCase, 0)
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		k, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad k", idx+1)
		}
		n := 1 << k
		if len(parts) != 1+n {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, 1+n, len(parts))
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(parts[1+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad value", idx+1)
			}
			arr[i] = v
		}
		cases = append(cases, testCase{k: k, arr: arr})
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(tc.k))
	sb.WriteByte('\n')
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expected := solve(tc.k)
		var exp strings.Builder
		for i, v := range expected {
			if i > 0 {
				exp.WriteByte(' ')
			}
			exp.WriteString(strconv.Itoa(v))
		}

		input := buildInput(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp.String() {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, exp.String(), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
