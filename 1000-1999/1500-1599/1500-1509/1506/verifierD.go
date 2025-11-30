package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesD.txt so the verifier is self-contained.
const testcasesRaw = `8 76 70 17 48 78 61 81 75
3 78 2 61
9 71 30 25 92 61 70 71 61 51
5 30 82 20 67 50
1 86
6 10 14 61 28 8 10
7 26 49 86 53 7 86 55
4 27 76 44 10
8 21 80 39 88 91 43 36 28
2 5 16
1 46
2 30 7
7 52 75 82 78 62 70 77
10 45 15 47 74 93 80 20 91 13 87
10 64 91 86 80 92 41 55 28 2 39
6 6 88 18 23 25 86
2 19 82
9 76 10 69 41 95 74 40 14 77
4 55 93 25 76
5 38 96 40 68 5
8 68 72 55 13 69 78 56 8
1 20
5 26 37 71 14 32
9 79 37 36 35 59 22 38 7 15
10 47 45 35 47 2 61 39 92 1 41
7 96 5 3 77 10 45 64
3 53 75 93
6 84 60 73 6 52 26
10 77 2 33 62 59 44 64 75 92 90
8 95 42 28 46 19 89 89 79
10 12 11 10 25 69 52 68 92 42 5
9 48 29 30 6 96 87 80 85 54
4 96 65 7 34
9 97 15 67 85 67 23 40 15 90
7 40 35 100 87 75 100 40
6 1 10 35 5 73 3
8 62 63 1 51 87 68 62 18
9 59 87 22 17 75 73 51 60 93
9 73 57 35 42 30 46 92 19 22
9 77 29 44 22 57 98 44 33 14
7 48 55 84 45 75 70 24
6 56 76 53 61 13 45
7 37 17 20 55 39 65 75
1 31
8 77 88 27 75 83 45 86 26
8 89 69 98 90 10 2 88 4
8 83 75 56 56 11 46 58 38
10 29 32 61 36 55 10 28 34 52 39
9 46 8 57 75 33 3 17 88 11
7 20 93 92 26 72 91 33
8 26 33 77 28 56 83 8 24
2 67 89
1 61
1 28
9 43 92 42 28 3 41 49 17 28
4 9 68 12 13
5 19 28 77 95 24
5 68 21 48 27 10
1 19
4 63 94 30 18
10 39 4 84 23 2 62 12 96 57 65
4 27 56 45 73
4 81 72 35 9
2 26 36
3 58 41 33
8 7 8 88 62 43 53 88 80
3 94 75 23
1 91
5 64 21 31 48 95
10 61 18 35 76 46 40 12 46 21 32
8 11 41 72 100 68 3 4 61
9 43 38 39 28 87 69 53 86 29
8 38 68 61 37 31 9 77 17
1 79
5 44 34 11 69 96
8 45 81 30 72 99 13 20 47
2 5 76
9 22 33 38 31 58 86 58 68 72
6 16 14 25 54 65 55
8 93 38 61 90 12 6 38 86
5 10 68 95 31 96
10 54 2 49 28 46 63 93 39 35 49
8 59 31 5 2 73 48 31 71
10 55 86 23 41 83 40 42 67 70 39
7 75 42 22 54 5 63 61
3 82 72 74
9 10 57 64 80 3 41 79 8 61
10 15 3 71 42 21 66 99 28 32 49
1 1
6 31 15 13 21 57 76
1 96
9 86 14 1 33 79 1 61 26 2
3 23 39 87
8 23 75 32 90 13 83 50 12
6 90 21 46 37 64 83
6 26 72 65 96 34 23
10 40 25 62 21 44 11 35 76 14 98
7 44 21 37 98 26 24 81
5 12 61 8 19 59
5 83 16 38 4 41
7 10 14 78 33 86 2 61
10 46 60 49 27 38 72 77 7 10 50
6 74 47 62 33 2 4
1 1
5 20 68 42 22 95
9 32 8 95 57 82 23 61 65 79
8 45 64 34 4 23 56 53 84
2 52 20
10 65 29 69 54 72 9 59 46 35 1
10 59 61 50 76 49 98 51 11 21 95
5 80 20 35 99 33
8 63 7 65 46 81 51 87 69
6 14 54 100 74 82 74
5 51 11 83 96 12
3 57 10 95
8 35 47 21 37 9 78 7 93`

type testCase struct {
	arr []int
}

func solveCase(tc testCase) int {
	freq := make(map[int]int)
	maxCnt := 0
	for _, x := range tc.arr {
		freq[x]++
		if freq[x] > maxCnt {
			maxCnt = freq[x]
		}
	}
	n := len(tc.arr)
	if 2*maxCnt <= n {
		if n%2 == 0 {
			return 0
		}
		return 1
	}
	return 2*maxCnt - n
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	res := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", idx+1, err)
		}
		if len(fields)-1 != n {
			return nil, fmt.Errorf("case %d: expected %d numbers got %d", idx+1, n, len(fields)-1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[i+1])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse value %d: %v", idx+1, i+1, err)
			}
			arr[i] = v
		}
		res = append(res, testCase{arr: arr})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	return res, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i, tc := range cases {
		expected := strconv.Itoa(solveCase(tc))
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(len(tc.arr)))
		sb.WriteByte('\n')
		for idx, v := range tc.arr {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
