package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// solve mirrors 1462D.go for one case.
func solve(arr []int64) int {
	n := len(arr)
	total := int64(0)
	for _, v := range arr {
		total += v
	}
	maxSeg := int64(1)
	var prefix int64
	for i := 0; i < n; i++ {
		prefix += arr[i]
		segSum := prefix
		if total%segSum != 0 {
			continue
		}
		need := total / segSum
		cur := int64(0)
		cnt := int64(0)
		ok := true
		for j := 0; j < n; j++ {
			cur += arr[j]
			if cur == segSum {
				cnt++
				cur = 0
			} else if cur > segSum {
				ok = false
				break
			}
		}
		if ok && cnt == need && need > maxSeg {
			maxSeg = need
		}
	}
	return n - int(maxSeg)
}

type testCase struct {
	input string
}

// Embedded testcases from testcasesD.txt.
const testcaseData = `
100
17 69 53 16 11 76 60 42 23 38 31 81 56 87 90 37 5 3
27 2 82 35 10 85 41 41 71 84 37 7 91 32 89 67 18 34 77 8 50 54 36 75 29 33 41 69
16 58 91 16 55 48 22 67 7 30 91 62 50 4 31 11 40
4 24 17 57 2
29 97 47 68 53 16 68 32 64 37 35 52 56 61 90 8 22 91 68 46 26 32 20 29 39 19 24 71 33 70
1 24
8 93 54 91 21 32 98 72 17
29 75 88 22 96 32 68 100 39 91 95 3 93 59 26 8 53 22 49 65 50 28 50 69 68 64 14 17 12 69
1 11
1 96
8 96 51 58 87 36 53 45 85
13 34 93 10 98 21 76 32 95 27 6 98 7 66
26 9 89 58 32 58 91 100 98 85 18 77 81 66 38 2 77 19 45 86 46 47 39 65 83 89 33
5 3 62 94 26 78
13 44 61 39 57 18 61 38 90 43 45 82 78 49
13 75 52 90 65 39 87 97 73 31 25 56 28 34
17 22 66 18 77 67 97 69 100 65 83 61 44 34 56 39 59 84
17 7 94 25 33 59 89 9 3 26 77 14 33 1 28 24 80 23
30 85 47 79 35 2 87 85 3 3 78 22 40 31 51 64 95 47 46 12 75 56 15 90 4 99 23 2 77 69 35
12 95 27 50 26 41 24 23 64 85 20 79 76
17 53 76 81 30 1 7 56 78 5 83 83 83 18 65 11 54 40
11 82 81 80 25 92 8 44 18 41 85 96
9 63 62 78 43 35 83 10 59 97
6 14 66 34 72 92 52
20 58 54 86 21 78 54 99 17 88 85 100 15 75 99 15 24 96 62 38 88
24 13 88 98 55 34 36 15 82 46 66 69 10 2 6 78 2 100 18 16 53 50 62 4 41
21 18 11 90 51 83 75 96 29 7 27 53 4 7 77 52 17 72 94 91 39 53
13 70 8 96 58 70 19 72 17 83 27 76 22 12
19 61 37 35 33 28 52 88 2 82 7 68 100 100 9 2 82 37 20 99
5 98 65 18 15 43
1 95
27 42 19 43 75 97 9 12 48 27 33 58 43 88 63 61 69 68 11 11 41 57 50 100 52 70 62 71
14 67 10 29 19 83 52 66 41 19 10 19 20 5 50
21 73 15 86 9 68 35 73 39 46 86 94 62 52 76 41 96 86 80 60 66 45
8 100 65 5 93 50 47 44 66
1 32
19 10 61 15 48 56 99 80 32 70 4 67 65 15 6 73 46 6 47 24
11 35 63 62 97 54 1 83 99 31 84 84
16 96 93 45 88 93 21 6 98 95 62 68 85 22 50 90 37
22 75 62 28 4 91 59 80 77 86 23 79 57 41 75 42 56 73 99 81 19 66 55
13 31 48 49 38 31 28 33 31 69 76 48 59 31
20 44 80 82 19 42 32 19 98 33 71 23 69 91 91 77 15 58 62 27 89
6 77 46 78 32 91 12
18 64 27 38 50 35 66 26 46 53 41 94 3 6 25 26 78 21 72
15 11 33 65 37 89 63 54 41 87 96 90 90 50 7 15
29 8 43 21 98 35 79 93 72 57 40 71 54 62 48 29 87 12 93 11 81 12 65 3 85 75 66 43 94 77
4 73 31 3 6
26 13 94 40 12 76 51 91 17 99 74 79 68 26 36 38 21 24 76 29 18 17 30 34 64 39 19
24 79 35 27 11 18 85 53 24 26 76 3 45 69 85 62 87 17 97 95 35 11 93 90 80
28 89 36 77 95 29 1 19 27 95 74 87 1 54 17 19 63 69 6 65 41 29 86 23 59 41 86 57 30
21 82 94 36 87 91 25 13 45 63 35 46 69 76 37 28 37 71 51 9 23 49
4 76 84 13 67
28 18 92 59 19 55 34 94 69 25 98 23 95 9 34 57 42 49 66 7 45 50 95 93 50 71 17 99 91
23 70 73 20 84 5 40 79 57 67 98 74 54 27 60 27 51 12 13 69 52 16 65 87
15 96 84 22 5 26 30 31 13 77 49 53 9 71 28 68
7 28 63 12 25 13 66 17
20 25 71 74 67 82 80 56 77 51 83 34 87 61 82 81 78 4 45 45 99
26 27 90 83 5 26 57 28 67 88 28 1 19 31 59 61 12 78 36 84 38 51 43 95 93 86 92
10 86 12 79 16 64 48 79 80 84 12
8 18 53 31 3 63 37 19 67
28 82 48 95 87 19 15 66 42 41 88 93 10 47 25 26 38 67 90 55 84 93 95 49 51 5 19 57 10
10 81 85 20 68 25 78 63 45 20 75
5 74 55 34 95 28
4 21 18 57 69
14 23 5 39 97 61 74 11 88 47 83 13 13 64 32
17 6 4 91 88 76 12 60 34 95 75 33 90 72 49 21 31 81
4 87 75 94 100
17 21 64 84 81 32 2 29 19 91 77 87 99 19 82 20 98 76
20 23 9 14 36 13 35 80 60 74 100 34 75 54 6 91 10 19 58 8 45
11 33 82 53 50 3 38 84 31 28 100 12
19 93 19 32 67 28 58 41 2 98 33 54 75 67 71 17 91 43 92 73
21 69 50 26 39 13 57 84 82 35 31 83 34 100 40 11 48 19 58 4 50 83
15 7 81 100 99 37 66 80 99 55 60 79 73 96 67 25
6 1 35 14 50 91 33
7 63 61 34 66 52 17 78
16 35 23 89 2 27 44 34 91 91 14 85 4 56 40 42 80
24 25 43 53 100 17 83 16 32 46 98 93 61 57 23 88 46 55 56 53 47 75 14 12 72
26 58 17 13 2 68 90 50 84 62 31 51 82 93 13 54 8 87 38 77 22 39 36 68 16 72 86
22 44 1 19 1 12 44 97 44 58 26 32 48 86 74 60 34 70 22 58 52 2 14
1 14
11 93 27 40 70 15 65 100 39 12 14 70
20 37 95 8 32 15 88 22 57 22 22 60 34 89 22 57 91 74 4 72 49
1 16
4 40 11 81 99
13 27 33 93 50 30 60 83 77 65 72 48 24 55
21 2 76 54 78 30 37 45 66 41 61 81 98 95 60 96 92 33 45 35 100 21
6 29 60 42 83 3 48
16 80 90 37 87 32 3 14 59 94 67 24 26 50 90 78 43
18 18 5 63 53 60 57 48 2 26 73 27 9 63 17 85 12 93 67
23 86 71 52 79 36 27 11 55 3 70 62 92 96 44 95 1 26 79 38 98 99 54 58
23 66 90 37 9 5 12 45 49 83 53 49 83 90 98 82 99 30 37 63 52 71 1 81
14 82 37 87 63 20 52 87 56 51 92 99 76 93 77
10 61 71 48 39 54 75 83 99 22 45
2 55 57
15 40 9 88 30 43 30 49 4 80 44 20 24 64 9 100
4 23 78 5 30
21 23 2 84 67 87 56 87 56 23 38 51 57 40 69 5 21 62 21 14 85 15
16 94 37 94 40 39 51 5 76 35 43 98 60 9 8 98 84
23 38 29 59 82 40 71 27 49 9 55 6 6 12 48 1 38 40 8 1 5 80 6 20
28 77 22 61 84 11 75 18 3 96 72 18 14 10 91 9 78 100 85 15 26 13 63 56 11 15 61 92 58
`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	pos++
	res := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d missing n", i+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		pos++
		if pos+n > len(fields) {
			return nil, fmt.Errorf("case %d missing values", i+1)
		}
		var input bytes.Buffer
		fmt.Fprintln(&input, 1)
		fmt.Fprintln(&input, n)
		for j := 0; j < n; j++ {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fields[pos+j])
		}
		input.WriteByte('\n')
		pos += n
		res = append(res, testCase{input: input.String()})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("testcase count mismatch: read %d tokens, have %d", pos, len(fields))
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, tc.input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func referenceSolve(input string) (string, error) {
	fields := strings.Fields(input)
	if len(fields) < 2 {
		return "", fmt.Errorf("invalid input")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return "", err
	}
	pos++
	if t != 1 {
		return "", fmt.Errorf("unexpected t %d", t)
	}
	n, err := strconv.Atoi(fields[pos])
	if err != nil {
		return "", err
	}
	pos++
	if pos+n > len(fields) {
		return "", fmt.Errorf("missing numbers")
	}
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		val, err := strconv.ParseInt(fields[pos+i], 10, 64)
		if err != nil {
			return "", err
		}
		arr[i] = val
	}
	return strconv.Itoa(solve(arr)), nil
}
