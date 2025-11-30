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

func runProg(bin string, input string) (string, error) {
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

const testcasesRaw = `3 2 2 3 1 2 2 3 3
4 3 3 4 1 2 2 3 3 4 4 5
5 4 4 5 1 2 3 2 4 3 5 4 6 5 7
6 5 5 6 1 2 3 4 2 5 3 6 4 7 5 8 6 9
7 6 6 7 1 2 3 4 5 2 6 3 7 4 8 5 9 6 10 7 11
8 7 7 8 1 2 3 4 5 6 2 7 3 8 4 9 5 10 6 11 7 12 8 13
9 8 8 9 1 2 3 4 5 6 7 2 8 3 9 4 10 5 11 6 12 7 13 8 14 9 15
10 1 9 10 1 2 3 4 5 6 7 8 2 9
2 2 2 1 2 10 1 11
3 3 2 3 1 2 11 3 12 1 13
4 4 4 1 2 3 2 12 3 13 4 14 1 15
5 5 3 4 5 1 2 2 13 3 14 4 15 5 16 1 17
6 6 2 3 4 5 6 1 2 14 3 15 4 16 5 17 6 18 1 19
7 7 1 2 3 4 5 6 7 2 15 3 16 4 17 5 18 6 19 7 20 1 21
8 8 8 1 2 3 4 5 6 7 2 16 3 17 4 18 5 19 6 20 7 21 8 22 1 23
9 1 8 9 1 2 3 4 5 6 7 2 17
10 2 8 9 10 1 2 3 4 5 6 7 2 18 3 19
2 3 1 2 2 19 1 20 2 21
3 4 2 3 1 2 20 3 21 1 22 2 23
4 5 1 2 3 4 2 21 3 22 4 23 1 24 2 25
5 6 2 3 4 5 1 2 22 3 23 4 24 5 25 1 26 2 27
6 7 5 6 1 2 3 4 2 23 3 24 4 25 5 26 6 27 1 28 2 29
7 8 3 4 5 6 7 1 2 2 24 3 25 4 26 5 27 6 28 7 29 1 30 2 31
8 1 1 2 3 4 5 6 7 8 2 25
9 2 8 9 1 2 3 4 5 6 7 2 26 3 27
10 3 7 8 9 10 1 2 3 4 5 6 2 27 3 28 4 29
2 4 2 1 2 28 1 29 2 30 1 31
3 5 2 3 1 2 29 3 30 1 31 2 32 3 33
4 6 2 3 4 1 2 30 3 31 4 32 1 33 2 34 3 35
5 7 1 2 3 4 5 2 31 3 32 4 33 5 34 1 35 2 36 3 37
6 8 2 3 4 5 6 1 2 32 3 33 4 34 5 35 6 36 1 37 2 38 3 39
7 1 5 6 7 1 2 3 4 2 33
8 2 2 3 4 5 6 7 8 1 2 34 3 35
9 3 8 9 1 2 3 4 5 6 7 2 35 3 36 4 37
10 4 6 7 8 9 10 1 2 3 4 5 2 36 3 37 4 38 5 39
2 5 1 2 2 37 1 38 2 39 1 40 2 41
3 6 2 3 1 2 38 3 39 1 40 2 41 3 42 1 43
4 7 3 4 1 2 2 39 3 40 4 41 1 42 2 43 3 44 4 45
5 8 5 1 2 3 4 2 40 3 41 4 42 5 43 1 44 2 45 3 46 4 47
6 1 5 6 1 2 3 4 2 41
7 2 7 1 2 3 4 5 6 2 42 3 43
8 3 3 4 5 6 7 8 1 2 2 43 3 44 4 45
9 4 8 9 1 2 3 4 5 6 7 2 44 3 45 4 46 5 47
10 5 5 6 7 8 9 10 1 2 3 4 2 45 3 46 4 47 5 48 6 49
2 6 2 1 2 46 1 47 2 48 1 49 2 50 1 51
3 7 2 3 1 2 47 3 48 1 49 2 50 3 51 1 52 2 53
4 8 4 1 2 3 2 48 3 49 4 50 1 51 2 52 3 53 4 54 1 55
5 1 4 5 1 2 3 2 49
6 2 2 3 4 5 6 1 2 50 3 51
7 3 2 3 4 5 6 7 1 2 51 3 52 4 53
8 4 4 5 6 7 8 1 2 3 2 52 3 53 4 54 5 55
9 5 8 9 1 2 3 4 5 6 7 2 53 3 54 4 55 5 56 6 57
10 6 4 5 6 7 8 9 10 1 2 3 2 54 3 55 4 56 5 57 6 58 7 59
2 7 1 2 2 55 1 56 2 57 1 58 2 59 1 60 2 61
3 8 2 3 1 2 56 3 57 1 58 2 59 3 60 1 61 2 62 3 63
4 1 1 2 3 4 2 57
5 2 3 4 5 1 2 2 58 3 59
6 3 5 6 1 2 3 4 2 59 3 60 4 61
7 4 4 5 6 7 1 2 3 2 60 3 61 4 62 5 63
8 5 5 6 7 8 1 2 3 4 2 61 3 62 4 63 5 64 6 65
9 6 8 9 1 2 3 4 5 6 7 2 62 3 63 4 64 5 65 6 66 7 67
10 7 3 4 5 6 7 8 9 10 1 2 2 63 3 64 4 65 5 66 6 67 7 68 8 69
2 8 2 1 2 64 1 65 2 66 1 67 2 68 1 69 2 70 1 71
3 1 2 3 1 2 65
4 2 2 3 4 1 2 66 3 67
5 3 2 3 4 5 1 2 67 3 68 4 69
6 4 2 3 4 5 6 1 2 68 3 69 4 70 5 71
7 5 6 7 1 2 3 4 5 2 69 3 70 4 71 5 72 6 73
8 6 6 7 8 1 2 3 4 5 2 70 3 71 4 72 5 73 6 74 7 75
9 7 8 9 1 2 3 4 5 6 7 2 71 3 72 4 73 5 74 6 75 7 76 8 77
10 8 2 3 4 5 6 7 8 9 10 1 2 72 3 73 4 74 5 75 6 76 7 77 8 78 9 79
2 1 1 2 2 73
3 2 2 3 1 2 74 3 75
4 3 3 4 1 2 2 75 3 76 4 77
5 4 1 2 3 4 5 2 76 3 77 4 78 5 79
6 5 5 6 1 2 3 4 2 77 3 78 4 79 5 80 6 81
7 6 1 2 3 4 5 6 7 2 78 3 79 4 80 5 81 6 82 7 83
8 7 7 8 1 2 3 4 5 6 2 79 3 80 4 81 5 82 6 83 7 84 8 85
9 8 8 9 1 2 3 4 5 6 7 2 80 3 81 4 82 5 83 6 84 7 85 8 86 9 87
10 1 1 2 3 4 5 6 7 8 9 10 2 81
2 2 2 1 2 82 1 83
3 3 2 3 1 2 83 3 84 1 85
4 4 4 1 2 3 2 84 3 85 4 86 1 87
5 5 5 1 2 3 4 2 85 3 86 4 87 5 88 1 89
6 6 2 3 4 5 6 1 2 86 3 87 4 88 5 89 6 90 1 91
7 7 3 4 5 6 7 1 2 2 87 3 88 4 89 5 90 6 91 7 92 1 93
8 8 8 1 2 3 4 5 6 7 2 88 3 89 4 90 5 91 6 92 7 93 8 94 1 95
9 1 8 9 1 2 3 4 5 6 7 2 89
10 2 10 1 2 3 4 5 6 7 8 9 2 90 3 91
2 3 1 2 2 91 1 92 2 93
3 4 2 3 1 2 92 3 93 1 94 2 95
4 5 1 2 3 4 2 93 3 94 4 95 1 96 2 97
5 6 4 5 1 2 3 2 94 3 95 4 96 5 97 1 98 2 99
6 7 5 6 1 2 3 4 2 95 3 96 4 97 5 98 6 99 1 100 2 101
7 8 5 6 7 1 2 3 4 2 96 3 97 4 98 5 99 6 100 7 101 1 102 2 103
8 1 1 2 3 4 5 6 7 8 2 97
9 2 8 9 1 2 3 4 5 6 7 2 98 3 99
10 3 9 10 1 2 3 4 5 6 7 8 2 99 3 100 4 101
2 4 2 1 2 100 1 101 2 102 1 103
3 5 2 3 1 2 101 3 102 1 103 2 104 3 105`

type testCase struct {
	n       int
	q       int
	a       []int
	queries [][2]int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("line %d malformed", idx+1)
		}
		pos := 0
		n, err := strconv.Atoi(parts[pos])
		if err != nil {
			return nil, fmt.Errorf("line %d parse n: %v", idx+1, err)
		}
		pos++
		qNum, err := strconv.Atoi(parts[pos])
		if err != nil {
			return nil, fmt.Errorf("line %d parse q: %v", idx+1, err)
		}
		pos++
		if len(parts) != 2+n+2*qNum {
			return nil, fmt.Errorf("line %d length mismatch", idx+1)
		}
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			val, err := strconv.Atoi(parts[pos])
			if err != nil {
				return nil, fmt.Errorf("line %d parse a[%d]: %v", idx+1, i, err)
			}
			pos++
			a[i] = val
		}
		queries := make([][2]int, qNum)
		for i := 0; i < qNum; i++ {
			first, err := strconv.Atoi(parts[pos])
			if err != nil {
				return nil, fmt.Errorf("line %d parse query %d first: %v", idx+1, i+1, err)
			}
			pos++
			second, err := strconv.Atoi(parts[pos])
			if err != nil {
				return nil, fmt.Errorf("line %d parse query %d second: %v", idx+1, i+1, err)
			}
			pos++
			queries[i] = [2]int{first, second}
		}
		cases = append(cases, testCase{n: n, q: qNum, a: a, queries: queries})
	}
	return cases, nil
}
func solveCase(n int, a []int, queries [][2]int) []int {
	// Mirrors 1719C.go logic.
	maxIdx := 1
	for i := 1; i <= n; i++ {
		if a[i] > a[maxIdx] {
			maxIdx = i
		}
	}
	wins := make([][]int, n+1)
	champion := 1
	for j := 2; j <= n; j++ {
		round := j - 1
		if a[champion] > a[j] {
			wins[champion] = append(wins[champion], round)
		} else {
			wins[j] = append(wins[j], round)
			champion = j
		}
	}
	res := make([]int, len(queries))
	for idx, q := range queries {
		i := q[0]
		k := q[1]
		limit := k
		if limit > n-1 {
			limit = n - 1
		}
		rounds := wins[i]
		cnt := sort.Search(len(rounds), func(x int) bool { return rounds[x] > limit })
		if i == maxIdx && k > n-1 {
			cnt += k - (n - 1)
		}
		res[idx] = cnt
	}
	return res
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	cases, err := parseTestcases()
	if err != nil {
		os.Exit(1)
	}
	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		inputBuilder := &strings.Builder{}
		fmt.Fprintf(inputBuilder, "1\n%d %d\n", tc.n, tc.q)
		for i := 1; i <= tc.n; i++ {
			if i > 1 {
				inputBuilder.WriteByte(' ')
			}
			fmt.Fprintf(inputBuilder, "%d", tc.a[i])
		}
		inputBuilder.WriteByte('\n')
		for _, qu := range tc.queries {
			fmt.Fprintf(inputBuilder, "%d %d\n", qu[0], qu[1])
		}
		out, err := runProg(bin, inputBuilder.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		expectedRes := solveCase(tc.n, tc.a, tc.queries)
		outTokens := strings.Fields(out)
		if len(outTokens) != len(expectedRes) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\n", idx+1, len(expectedRes), len(outTokens))
			os.Exit(1)
		}
		for i, tok := range outTokens {
			val, err := strconv.Atoi(tok)
			if err != nil || val != expectedRes[i] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %v\n", idx+1, expectedRes, outTokens)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
