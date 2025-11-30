package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `3 10 98 9 33
2 8 98 58
8 7 27 13 63 4 50 56 78 98
1 8 35
4 10 14 41 4 3
1 9 2
7 4 55 93 4 68 29 98 57
8 9 30 45 30 87 29 98 59 38
1 7 72
2 3 81 93
5 2 96 43 93 92 65
7 9 86 25 39 37 76 64 65
7 10 5 62 32 96 52 54 86
3 6 71 90 100
6 2 57 85 66 14 100 21
9 7 48 63 94 4 61 6 40 91 79
10 10 51 83 22 22 65 30 2 99 26 70
9 4 52 66 45 74 46 59 35 85 71
10 1 50 95 66 17 67 100 72 27 55 8
8 6 73 71 26 65 53 63 46 54
6 1 69 70 80 79 43 59
10 1 30 82 23 71 75 24 12 71 33 5
2 2 3 58
1 5 32
5 2 80 24 45 38 9
3 3 33 68 22
5 5 59 90 42 64 61
2 1 40 50
6 7 25 34 14 33 94 66
4 10 56 3 29 3
7 3 5 93 21 58 91 65 87
7 9 29 81 89 67 58 29 68
1 7 87
10 6 85 81 55 8 95 39 17 28 7 40
2 2 40 39
3 7 73 33 17
1 9 5
10 4 73 59 22 100 91 80 66 5 49 26
6 2 27 74 87 56 76 25
8 2 86 50 38 65 64 3 42 79
7 5 3 21 26 42 73 18 44
7 4 35 87 13 49 71 45 88
9 8 99 69 31 9 93 6 11 18 22
3 9 28 35 98
6 10 65 33 48 44 44 15
5 4 78 100 92 63 18
10 9 99 14 42 6 53 10 49 19 17 44
2 10 76 49
2 10 71 29
10 2 35 47 38 73 69 15 59 36 14 6
5 1 79 86 2 12 53
2 1 25 31
10 7 21 15 58 22 88 31 21 96 14 56
7 9 38 71 33 92 62 41 13
4 6 6 4 2 38
10 6 58 51 41 52 9 9 41 77 59 15
5 4 80 100 70 89 61
6 5 24 70 27 40 26 32
6 2 36 12 97 58 12 84
10 6 30 50 40 6 42 24 41 75 39 32
6 2 70 79 75 77 12 32
4 1 32 52 10 35
9 2 94 10 3 82 2 38 97 46 64
8 3 13 65 100 42 10 66 86 23
3 3 19 41 40
2 9 78 38
3 4 19 70 93
1 6 80
9 4 23 39 56 69 21 7 92 86 32
5 2 88 58 56 71 33
9 8 69 59 2 51 44 22 34 63 4
7 10 3 8 89 46 75 18 76
3 3 34 36 51
10 7 23 79 12 30 63 1 23 68 41 65
8 4 31 41 64 88 62 29 92 53
6 9 79 94 84 36 83 29
1 2 98
9 6 21 66 99 27 40 39 89 39 71
6 3 90 90 95 60 77 11
2 10 66 74
7 3 20 33 55 28 73 93 97
1 8 88
7 6 50 66 22 70 94 6 68
2 5 81 13
5 2 18 100 79 85 88
2 8 31 49
7 7 22 42 57 17 80 63 28
2 7 77 69
7 2 85 38 36 32 49 96 72
1 4 68
8 10 3 4 81 78 32 34 27 23
5 3 70 26 35 40 75
5 8 22 70 46 63 54
2 4 74 50
4 5 14 4 16 73
1 9 38
3 2 65 48 74
5 7 65 87 46 98 68
6 1 16 57 92 58 45 40
9 7 44 94 88 74 64 15 83 49 49
`

type testCase struct {
	n   int
	m   int
	arr []int
}

func solveCase(n, m int, arr []int) int {
	cnt := make([]int, m)
	for _, v := range arr {
		cnt[v%m]++
	}
	ans := 0
	if cnt[0] > 0 {
		ans++
	}
	for r := 1; r*2 < m; r++ {
		a := cnt[r]
		b := cnt[m-r]
		if a == 0 && b == 0 {
			continue
		}
		diff := a - b
		if diff < 0 {
			diff = -diff
		}
		if diff <= 1 {
			ans++
		} else {
			ans += diff
		}
	}
	if m%2 == 0 && cnt[m/2] > 0 {
		ans++
	}
	return ans
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var res []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid n: %v", idx+1, err)
		}
		if len(parts) != n+2 {
			return nil, fmt.Errorf("line %d: expected %d numbers, got %d", idx+1, n+2, len(parts))
		}
		m, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid m: %v", idx+1, err)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(parts[2+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse value %d: %v", idx+1, i+1, err)
			}
			arr[i] = v
		}
		res = append(res, testCase{n: n, m: m, arr: arr})
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		var input strings.Builder
		fmt.Fprintf(&input, "1\n%d %d\n", tc.n, tc.m)
		for idx, v := range tc.arr {
			if idx > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		expected := strconv.Itoa(solveCase(tc.n, tc.m, tc.arr))
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
