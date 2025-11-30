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

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `
2 91 28
1 54
4 43 70 60 54
1 27
4 50 99 75 90
1 98
5 49 62 1 46 39
4 54 69 96 95
5 78 29 63 29 35
4 63 4 50 44
4 93 22 60 17
5 69 4 51 76 73
1 11
4 18 60 24 7
3 49 42 28
4 42 44 98 49
3 97 54 33
1 61
1 96
5 7 45 29 84 9
1 97
1 32
2 3 80
2 31 17
4 86 15 73 28
4 90 33 99 48
2 78 78
1 100
2 40 14
5 4 40 74 87 49
4 92 26 10 76
2 14 90
3 88 77 16
5 6 45 69 55 85
3 9 65 83
3 2 54 63
1 56
3 82 59 91
2 56 23
5 84 35 79 69 100
5 99 33 48 31 26
4 42 85 3 65
3 30 81 25
2 11 70
3 12 30 96
1 53
3 83 68 71
1 87
1 18
3 22 71 17
4 92 27 88 32
1 34
5 36 94 83 9 49
3 90 77 72
1 98
3 84 18 7
5 24 53 9 29 80
5 96 55 55 61 80
5 92 50 13 12 21
4 65 25 31 34
2 60 68
2 36 95
2 38 77
4 65 10 41 18
1 26
1 80
4 77 10 62 41
5 65 36 13 41 85
5 49 50 43 39 74
3 57 4 65
3 23 21 8
1 43
2 5 59
5 85 83 53 31 17
2 16 88
1 82
1 6
1 10
2 44 62
3 3 2 86
3 48 8 69
2 97 1
1 1
3 39 10 92
2 10 8
2 53 49
1 63
2 92 10
2 35 3
3 69 15 14
3 16 87 59
2 73 83
3 44 21 45
5 41 15 89 17 61
4 84 73 20 7
4 23 11 23 5
3 51 22 47
5 89 60 85 41 38
3 37 14 9
5 18 22 41 67 20
1 81
1 100
`

type testCase struct {
	n   int
	arr []int
}

func solveCase(tc testCase) string {
	a := append([]int(nil), tc.arr...)
	sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
	used := make([]bool, tc.n+1)
	ok := true
	for _, v := range a {
		for v > tc.n || (v > 0 && used[v]) {
			v /= 2
		}
		if v == 0 {
			ok = false
			break
		}
		used[v] = true
	}
	if ok {
		for i := 1; i <= tc.n; i++ {
			if !used[i] {
				ok = false
				break
			}
		}
	}
	if ok {
		return "YES"
	}
	return "NO"
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, n+1, len(fields))
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse a[%d]: %v", idx+1, i, err)
			}
			arr[i] = v
		}
		cases = append(cases, testCase{n: n, arr: arr})
	}
	return cases, nil
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
		expected := solveCase(tc)

		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(tc.n))
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
